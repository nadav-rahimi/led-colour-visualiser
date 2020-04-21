package lcv

import (
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gordonklaus/portaudio"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mjibson/go-dsp/fft"
	"log"
	"math/cmplx"
	"strings"
)

// Constant and Variable Setup
const (
	// The maximum frequency the program will clamp to
	fCap = 4000
	// The upper range of frequencies the program considers useful.
	// After this barrier, the colour changes very slowly in relation
	// to change in frequency
	usefulCap = 1000
	// The total range of hues to use for the hsv colour spectrum
	// e.g. the first 200 hues
	totalHue = 300
	// The hue colour at which the usefulCap is reached
	fCapHue = totalHue - 10
	// The length of the buffer used to store the audio data
	bufferLength = 1024 * 2
	// This is the length the program uses to find the frequency with
	// the highest magnitude, this is half the buffer length because
	// the FFT is mirrored along the centre, thus only half the length
	// needs to be used
	bufferLengthUseful = bufferLength / 2
	// The length of the array to use for damping
	freqArrayL = 5
	// The number of interpolation points to use
	interpNum = 1
)

var (
	// Whether to enable damping
	damp bool = true
	// Whether to enable interpolation
	interp bool = true
)

// Port Audio Functions
func StartPortAudio() {
	//Initialise portaudio and create the audio buffer
	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]float32, bufferLength)

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
	//log.Println("The input device is:  ", inpDev.Name)

	// Creating parameters
	p := portaudio.LowLatencyParameters(inpDev, outDev)
	p.FramesPerBuffer = len(buffer)
	var sampleRate = p.SampleRate
	var maxInfo = sampleRate / 2
	var fBinSize = int(maxInfo / bufferLengthUseful)

	// Create the stream
	stream, err := portaudio.OpenStream(p, buffer)
	chk(err)

	// Starting the stream
	chk(stream.Start())
	defer stream.Close()

	// Prepare variables for the stream
	buffer_64 := make([]float64, bufferLength)
	freqArray := make([]int, freqArrayL)
	var freqCounter = new(int)
	var frequency = new(int)
	var oldFreq = new(int)

	// Start processing the stream
	for {
		chk(stream.Read())

		// Convert the buffer values to float64
		for i := 1; i < bufferLength; i++ {
			buffer_64[i] = float64(buffer[i])
		}

		// Perform the FFT on the buffer
		buffer_fft := fft.FFTReal(buffer_64)

		// Get the index of the frequency with the largest magnitude
		var index int = maxFreqInd(&buffer_fft)

		// Calculate the new frequency
		*oldFreq = *frequency
		updateFreq(frequency, fBinSize, index)
		log.Print("Frequency: ", *frequency)

		// Dampening
		if damp {
			dampFreqs(frequency, &freqArray, freqCounter)
		}

		// Interpolation
		hue := getHue(float64(*frequency))
		old_hue := getHue(float64(*oldFreq))

		if interp {
			interpolate(hue, old_hue)
		} else {
			changeColour(hue)
		}

	}
}

type boxColour colorful.Color

// Converts the box colour type to a uint32 value
func (col boxColour) UINT32() uint32 {
	var r = uint32(col.R*255.0 + 0.5)
	var g = uint32(col.G*255.0 + 0.5)
	var b = uint32(col.B*255.0 + 0.5)

	r = r << 8
	r = r + g
	r = r << 8
	r = r + b

	return r
}

// Calculates the hue of the visualised frequency on the hsv scale
func getHue(f float64) float64 {
	if f > usefulCap {
		return fCapHue + (totalHue-fCapHue)*float64(f)/float64(fCap)
	}
	return float64(f) / float64(usefulCap) * fCapHue
}

// Takes in the current frequency and damps it based on past frequencies
func dampFreqs(f *int, farr *[]int, c *int) {
	(*farr)[*c] = *f
	*c++
	*c = *c % freqArrayL

	var total int = 0
	for _, value := range *farr {
		total += value
	}

	*f = total / freqArrayL
}

// From the fft array, the index of the frequency with the highest magnitude is returned
func maxFreqInd(b *[]complex128) int {
	var max_v float64 = 0
	var index int = 0

	for i := 1; i < bufferLengthUseful; i++ {
		e := cmplx.Abs((*b)[i])
		if e > max_v {
			max_v = e
			index = i
		}
	}

	return index
}

// Updates the frequency with the value of the frequency with the highest magnitude
func updateFreq(f *int, binsize, i int) {
	*f = binsize * i
	if *f > fCap {
		*f = fCap
	}
}

// Interpolates the colour switching as opposed to directly switching colours
func interpolate(nh, oh float64) {
	var interpInc float64 = (nh - oh) / (interpNum + 1)

	for i := 1; i <= (interpNum + 1); i++ {
		changeColour(oh + interpInc*float64(i))
	}
}
