package lcv

import (
	"fmt"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gordonklaus/portaudio"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/nadav-rahimi/led-colour-visualiser/fftsingle"
	"log"
	"math/cmplx"
	"strings"
	"time"
)

// The Audio Analyser
type AudioAnalyser struct {
	// Parameters to initialise the analyser
	param *AudioAnalysisParams
	// The units the analyser uses during processing
	u *AudioAnalysisUnits
	// Contains the slices which the analyser logs to
	lg *AudioAnalysisLogs
	// The callback function which the analyser calls with the hue of the frequency
	cb func(uint32)
	// The handler for the udp stream
	udph *AudioAnalysisUDPHandler
}

// Parameters for setting up the analyser
type AudioAnalysisParams struct {
	// The maximum f the program will clamp to
	fCap float64
	// The upper range of frequencies the program considers useful.
	// After this barrier, the colour changes very slowly in relation
	// to change in f
	usefulCap float64
	// The total range of hues to use for the hsv colour spectrum
	// e.g. the first 200 hues
	totalHue float64
	// The hue colour at which the usefulCap is reached
	fCapHue float64
	// The length of the buffer used to store the audio data
	bufferLength int
	// This is the length the program uses to find the f with
	// the highest magnitude, this is half the buffer length because
	// the FFT is mirrored along the centre, thus only half the length
	// needs to be used
	bufferLengthUseful float64
	// The length of the array to use for damping
	freqArrayL int
	// Whether to enable damping
	damp bool
	// Whether to enable smoothing
	smooth bool
	// Smoothing alpha
	smoothA float64
	// Should a graph be created after visualisation is stopped stops
	creatVis bool
	// The name of the gradient to use for audio colouring
	gradName string
	// THe start of the name of the sound input device which portaudio reads from
	inputDeviceName string
}

// Stores values the analyser uses during computation
type AudioAnalysisUnits struct {
	// Array holding the past few max frequencies used for
	// damping
	farr []int
	// Counter for farr used to update it without constantly
	// shifting the array, check the dampFreqs function for
	// more reference
	c *int
	// The old frequency from one audio chunk before, used
	// to smooth the processing
	old_freq *int
	// The frequency with the highest magnitude
	f *int
	// For each audio chunk, this is the index in the bfft
	// where the frequency with the loudest magnitude is
	// found
	index int
	// This represents the difference in frequency between
	// each index of the bfft array
	fBinSize int
	// Buffer to hold the calculated FFT of the audio stream
	bfft []complex64
	// Enables or disables the use of custom gradients
	gtUsed bool
	// The gradient table used for custom gradients
	aaGT *GradientTable
	// Stop signal
	stopSig chan bool
	// States whether the analyser is running
	isRunning bool
}

// The slices which the analyser logs to for graphing
type AudioAnalysisLogs struct {
	// Buffer to hold the original calculated frequency for each audio chunk
	freqLog []int
	// Buffer to hold the damped frequency for each audio chunk
	dampLog []int
	// Buffer to hold the smoothed frequency for each audio chunk
	smthLog []int
}

// Handles the sending of the colour data to the arduino through a udp stream
type AudioAnalysisUDPHandler struct {
	// The UDP client
	client *udpC
	// Decides if the client should attempt to send a message over the UDP stream
	shouldsend bool
	// States whether the UDP client is running or not
	running bool
}

// From the fft array, the index of the f with the highest magnitude is returned
func (aa AudioAnalyser) maxFreqInd() int {
	var max_v float64 = 0
	var index int = 0

	for i := 1; i < int(aa.param.bufferLengthUseful); i++ {
		e := cmplx.Abs(complex128(aa.u.bfft[i]))
		if e > max_v {
			max_v = e
			index = i
		}
	}

	return index
}

// Converts the frequency calculated to a uint32 colour
func (aa AudioAnalyser) colourUINT32() uint32 {
	var h float64
	if float64(*aa.u.f) > aa.param.usefulCap {
		h = aa.param.fCapHue + (aa.param.totalHue-aa.param.fCapHue)*(float64(*aa.u.f)/aa.param.fCap)
	} else {
		h = float64(*aa.u.f) / aa.param.usefulCap * aa.param.fCapHue
	}

	if aa.u.gtUsed {
		return boxColour(aa.u.aaGT.GetInterpolatedColorFor(h / aa.param.totalHue)).UINT32()
	}
	return boxColour(colorful.Hsv(h, 1, 1)).UINT32()
}

// Takes in the current f and damps it based on past frequencies
func (aa AudioAnalyser) dampFreqs() {
	(aa.u.farr)[*aa.u.c] = *aa.u.f
	*aa.u.c++
	// Modulus of the counter is taken so farr can be updated without shifting
	*aa.u.c = *aa.u.c % aa.param.freqArrayL

	var total int = 0
	for _, value := range aa.u.farr {
		total += value
	}

	*aa.u.f = total / aa.param.freqArrayL
}

// Smooths the frequencies, alternative damping method, the larger alpha the more
// the old freq is weighted, alpha is [0, 1]
func (aa AudioAnalyser) smoothFreqs(alpha float64) {
	(*aa.u.f) = int(alpha*float64(*aa.u.old_freq) + (1-alpha)*float64(*aa.u.f))
}

// Updates the f with the value of the f with the highest magnitude
func (aa AudioAnalyser) updateFreq() {
	*aa.u.f = aa.u.fBinSize * aa.u.index
	// After the cap range our ears dont hear a difference so no use to visualise the cap
	if float64(*aa.u.f) > aa.param.fCap {
		*aa.u.f = int(aa.param.fCap)
	}
}

// Begins analysing audio from an audio stream
func (aa AudioAnalyser) StartAnalysis() {
	aa.u.isRunning = true

	// Check the gradient table exists if one is to be used
	if len(aa.param.gradName) > 0 {
		var err error
		aa.u.aaGT, err = getGradientTable(aa.param.gradName)
		chk(err)
		aa.u.gtUsed = true
	}

	//Initialise portaudio and create the audio buffer
	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]float32, aa.param.bufferLength)

	// Get the input device
	devices, err := portaudio.Devices()
	chk(err)

	var inpDev, outDev *portaudio.DeviceInfo
	for _, d := range devices {
		if d.HostApi.Name == "MME" {
			if strings.HasPrefix(d.Name, aa.param.inputDeviceName) {
				if d.MaxInputChannels > 0 {
					inpDev = d
				}
			}
		}
	}
	log.Println("The input device is:  ", inpDev.Name)

	// Creating parameters
	p := portaudio.LowLatencyParameters(inpDev, outDev)
	p.FramesPerBuffer = len(buffer)
	var sampleRate = p.SampleRate
	var maxInfo = sampleRate / 2
	aa.u.fBinSize = int(maxInfo / aa.param.bufferLengthUseful)

	// Create the stream
	stream, err := portaudio.OpenStream(p, buffer)
	chk(err)

	// Starting the stream
	chk(stream.Start())
	defer stream.Close()

	// Prepare variables for the stream
	aa.u.farr = make([]int, aa.param.freqArrayL)
	aa.u.c = new(int)
	aa.u.old_freq = new(int)
	aa.u.f = new(int)

	// Variables setup to record the data
	aa.lg.freqLog = make([]int, 1)
	aa.lg.dampLog = make([]int, 1)
	aa.lg.smthLog = make([]int, 1)

	// Variables to break the loop
	var breakLoop bool = false

	startTime := time.Now()
	// Start processing the stream
	for {
		chk(stream.Read())

		// Perform the FFT on the buffer
		aa.u.bfft = fftsingle.FFTReal(buffer)

		// Get the index of the f with the largest magnitude
		aa.u.index = aa.maxFreqInd()

		// Calculate the new frequency
		*aa.u.old_freq = *aa.u.f
		aa.updateFreq()
		log.Print("Frequency: ", *aa.u.f)
		aa.lg.freqLog = append(aa.lg.freqLog, *aa.u.f)

		// Dampening and Smoothing
		if aa.param.smooth {
			aa.smoothFreqs(aa.param.smoothA)
			aa.smoothFreqs(0.3)
			log.Print("Smoothed Frequency: ", *aa.u.f)
			aa.lg.smthLog = append(aa.lg.smthLog, *aa.u.f)
		}
		if aa.param.damp {
			aa.dampFreqs()
			log.Print("Damped Frequency: ", *aa.u.f)
			aa.lg.dampLog = append(aa.lg.dampLog, *aa.u.f)
		}

		// Calling the callback function with the colour value
		aa.cb(aa.colourUINT32())

		// Sending the value through the UDP stream
		if aa.udph.shouldsend {
			aa.udph.client.sendMsg(fmt.Sprint(aa.colourUINT32()))
		}

		// The analyser is stopped through the sig channel
		select {
		case <-aa.u.stopSig:
			breakLoop = true
		default:
		}

		if breakLoop {
			break
		}
	}
	aa.u.isRunning = false
	endTime := time.Now()
	chk(stream.Stop())
	if aa.param.creatVis {
		names := []string{"Original F", "Smoothed F", "Damped F"}
		// Start and end times are taken to find the elapsed time and scale the width of the graph generated
		createGraph(names, endTime.Sub(startTime), &aa.lg.freqLog, &aa.lg.smthLog, &aa.lg.dampLog)
	}
}

// Stops analysis of the audio stream
func (aa AudioAnalyser) StopAnalysis() {
	if aa.u.isRunning {
		aa.u.stopSig <- true
		aa.u.isRunning = false
	}
}

// Generates a new analyser object with default configuration
func newAudioAnalyser(f func(uint32), g string) *AudioAnalyser {
	return &AudioAnalyser{
		param: &AudioAnalysisParams{
			fCap:               2500,
			usefulCap:          1200,
			totalHue:           320,
			fCapHue:            310,
			bufferLength:       1024 * 2,
			bufferLengthUseful: 1024,
			freqArrayL:         4,
			damp:               true,
			smooth:             true,
			smoothA:            0.73,
			creatVis:           false,
			gradName:           g,
			inputDeviceName:    "Line 1",
		},
		u: &AudioAnalysisUnits{
			stopSig: make(chan bool),
		},
		lg: &AudioAnalysisLogs{},
		cb: f,
		udph: &AudioAnalysisUDPHandler{
			client:     newUdpC("6969"),
			shouldsend: false,
		},
	}
}

// Get portaudio input devices available
func getInputDevices() []string {
	portaudio.Initialize()
	defer portaudio.Terminate()

	devices, err := portaudio.Devices()
	chk(err)

	names := make([]string, 0)
	for _, d := range devices {
		if d.HostApi.Name == "MME" {
			if d.MaxInputChannels > 0 {
				names = append(names, d.Name)
			}
		}
	}

	return names
}
