package lcv

import (
	"./fftsingle"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gordonklaus/portaudio"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math/cmplx"
	"strings"
)

// Constant and Variable Setup
const (
	// The maximum frequency the program will clamp to
	fCap = 2500
	// The upper range of frequencies the program considers useful.
	// After this barrier, the colour changes very slowly in relation
	// to change in frequency
	usefulCap = 900
	// The total range of hues to use for the hsv colour spectrum
	// e.g. the first 200 hues
	totalHue = 320
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
	freqArrayL = 10
	// Whether to enable damping
	damp bool = true
	// Whether to enable smoothing
	smooth bool = true
	// Smoothing alpha
	smoothA float64 = 0.73
)

var sig = make(chan bool)

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
	freqArray := make([]int, freqArrayL)
	var freqCounter = new(int)
	var old_freq = new(int)
	var frequency = new(int)

	//// Variables setup to record the data
	//freqLog := make([]int, 1)
	//dampLog := make([]int, 1)
	//smthLog := make([]int, 1)

	// Start processing the stream
	for {
		chk(stream.Read())

		// Perform the FFT on the buffer
		buffer_fft := fftsingle.FFTReal(buffer)

		// Get the index of the frequency with the largest magnitude
		var index int = maxFreqInd(&buffer_fft)

		// Calculate the new frequency
		*old_freq = *frequency
		updateFreq(frequency, fBinSize, index)
		log.Print("Frequency: ", *frequency)

		// Dampening and Smoothing
		if smooth {
			smoothFreqs(frequency, old_freq, smoothA)
			log.Print("Smoothed Frequency: ", *frequency)
		}
		if damp {
			dampFreqs(frequency, &freqArray, freqCounter)
			log.Print("Damped Frequency: ", *frequency)
		}

		// Changing the colour
		hue := getHue(float64(*frequency))
		changeColour(hue)

		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())

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

// Smooths the frequencies, alternative damping method
func smoothFreqs(f, of *int, alpha float64) {
	(*f) = int(alpha*float64(*of) + (1-alpha)*float64(*f))
}

// From the fft array, the index of the frequency with the highest magnitude is returned
func maxFreqInd(b *[]complex64) int {
	var max_v float64 = 0
	var index int = 0

	for i := 1; i < bufferLengthUseful; i++ {
		e := cmplx.Abs(complex128((*b)[i]))
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
