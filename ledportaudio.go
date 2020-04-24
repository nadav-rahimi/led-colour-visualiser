package lcv

import (
	"./fftsingle"
	"fmt"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gordonklaus/portaudio"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math/cmplx"
	"strings"
	"time"
)

// The Audio Analyser
type AudioAnalyser struct {
	param *AudioAnalysisParams
	u     *AudioAnalysisUnits
	lg    *AudioAnalysisLogs
	cb    func(uint32)
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
}

// Stores values the analyser uses during computation
type AudioAnalysisUnits struct {
	farr     []int
	c        *int
	old_freq *int
	f        *int
	index    int
	fBinSize int
	bfft     []complex64
	gtUsed   bool
	aaGT     GradientTable
	// Stop signal
	stopSig chan bool
}

// The slices which the analyser logs to for graphing
type AudioAnalysisLogs struct {
	freqLog []int
	dampLog []int
	smthLog []int
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
	*aa.u.c = *aa.u.c % aa.param.freqArrayL

	var total int = 0
	for _, value := range aa.u.farr {
		total += value
	}

	*aa.u.f = total / aa.param.freqArrayL
}

// Smooths the frequencies, alternative damping method
func (aa AudioAnalyser) smoothFreqs(alpha float64) {
	(*aa.u.f) = int(alpha*float64(*aa.u.old_freq) + (1-alpha)*float64(*aa.u.f))
}

// Updates the f with the value of the f with the highest magnitude
func (aa AudioAnalyser) updateFreq() {
	*aa.u.f = aa.u.fBinSize * aa.u.index
	if float64(*aa.u.f) > aa.param.fCap {
		*aa.u.f = int(aa.param.fCap)
	}
}

// Begins analysing audio from an audio stream
func (aa AudioAnalyser) StartAnalysis() {
	// Check the gradient table exists if one is to be used
	fmt.Println(aa.param.gradName)
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

	// TODO: function for retrieving the devices
	// Get the Virtual Audio Cable input device
	devices, err := portaudio.Devices()
	chk(err)

	var inpDev, outDev *portaudio.DeviceInfo
	for _, d := range devices {
		if d.HostApi.Name == "MME" {
			if strings.HasPrefix(d.Name, "Line 1") {
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

		// Make sig part of port audio
		select {
		case <-aa.u.stopSig:
			breakLoop = true
		default:
		}

		if breakLoop {
			break
		}
	}
	endTime := time.Now()
	chk(stream.Stop())
	if aa.param.creatVis {
		names := []string{"Original F", "Smoothed F", "Damped F"}
		createGraph(names, endTime.Sub(startTime), &aa.lg.freqLog, &aa.lg.smthLog, &aa.lg.dampLog)
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
			freqArrayL:         7,
			damp:               true,
			smooth:             true,
			smoothA:            0.73,
			creatVis:           true,
			gradName:           g,
		},
		u: &AudioAnalysisUnits{
			stopSig: make(chan bool),
		},
		lg: &AudioAnalysisLogs{},
		cb: f,
	}
}
