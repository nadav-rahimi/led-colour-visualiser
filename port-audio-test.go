package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gordonklaus/portaudio"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mjibson/go-dsp/fft"
	"math/cmplx"
	"math/rand"
	"strings"
)

// Constant and Variable Setup
const fCap = 4000
const usefulCap = 900
const totalHue = 200
const fCapHue = totalHue - 10
const bufferLength = 1024 * 2
const bufferLengthUseful = bufferLength / 2
const freqArrayL = 6
const interpNum = 2

var coloured_square *ui.Area
var current_colour_hex *uint32 = new(uint32)

// General Functions
func chk(err error) {
	if err != nil {
		panic(err)
	}
}

// UI Functions
func mkSolidBrush(color uint32, alpha float64) *ui.DrawBrush {
	brush := new(ui.DrawBrush)
	brush.Type = ui.DrawBrushTypeSolid
	component := uint8((color >> 16) & 0xFF)
	brush.R = float64(component) / 255
	component = uint8((color >> 8) & 0xFF)
	brush.G = float64(component) / 255
	component = uint8(color & 0xFF)
	brush.B = float64(component) / 255
	brush.A = alpha
	return brush
}

type areaHandler struct{}

func (areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	// fill the area with white
	brush := mkSolidBrush(*current_colour_hex, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(0, 0, p.AreaWidth, p.AreaHeight)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	// do nothing
	//*current_colour_hex = rand.Uint32()
	//coloured_square.QueueRedrawAll()
}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (areaHandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// do nothing
	return false
}

func setupUI() {
	mainwin := ui.NewWindow("LED Colour Visualiser", 480, 480, false)
	mainwin.SetMargined(true)
	mainwin.OnClosing(func(*ui.Window) bool {
		mainwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	mainwin.SetChild(hbox)
	coloured_square = ui.NewArea(areaHandler{})
	hbox.Append(coloured_square, true)

	visualise_button := ui.NewButton("Button")
	visualise_button.OnClicked(func(b *ui.Button) {
		go startPortAudio()
	})
	hbox.Append(visualise_button, false)

	mainwin.Show()
}

// Port Audio Functions
type boxColour colorful.Color
func (col boxColour) RGB255() (r, g, b uint32) {
	r = uint32(col.R*255.0 + 0.5)
	g = uint32(col.G*255.0 + 0.5)
	b = uint32(col.B*255.0 + 0.5)
	return r, g, b
}

func getHue(f float64) float64 {
	if f > usefulCap {
		return fCapHue + (totalHue-fCapHue)*float64(f)/float64(fCap)
	}
	return float64(f)/float64(usefulCap) * fCapHue
}

func startPortAudio()  {
	//Initialise portaudio and create the audio buffer
	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]float32, bufferLength)

	// Get the Virtual Audio Cable input device
	devices, err := portaudio.Devices()
	chk(err)

	var inpDev, outDev *portaudio.DeviceInfo
	for _, d := range devices {
		if (d.HostApi.Name == "MME") {
			if (strings.HasPrefix(d.Name, "Line 1")) {
				if (d.MaxInputChannels > 0) {
					inpDev = d
				}
			}
		}
	}
	fmt.Println("The input device is:  ", inpDev.Name)

	// Creating parameters
	p := portaudio.LowLatencyParameters(inpDev, outDev)
	p.FramesPerBuffer = len(buffer)
	var sampleRate = p.SampleRate
	var maxInfo = sampleRate / 2
	var fBinSize = int(maxInfo/bufferLengthUseful)

	// Create the stream
	stream, err := portaudio.OpenStream(p, buffer)
	chk(err)

	// Starting the stream
	chk(stream.Start())
	defer stream.Close()

	// Prepare variables for the stream
	buffer_64 := make([]float64, bufferLength)
	//buffer_fft_normalised := make([]float64, bufferLengthUseful)
	freqArray := make([]float64, freqArrayL)
	var freqCounter int = 0
	var oldFreq = new(int)
	var frequency = new(int)

	// Start processing the stream
	for {
		chk(stream.Read())
		//fmt.Println(buffer[:10])

		// Convert the buffer values to float64
		for i := 1; i < bufferLength; i++ {
			buffer_64[i] = float64(buffer[i])
		}

		// Perform the FFT on the buffer
		buffer_fft := fft.FFTReal(buffer_64)


		// Get the index of the frequency with the largest magnitude
		var max_v float64 = 0
		var max_v_i int = 0

		for i := 1; i < bufferLengthUseful; i++ {
			e := cmplx.Abs(buffer_fft[i])
			if e > max_v {
				max_v = e
				max_v_i = i
			}
		}

		// Calculate the new frequency
		*oldFreq = *frequency
		*frequency  = fBinSize * max_v_i
		if (*frequency > fCap) {
			*frequency = fCap
		}
		fmt.Println(*frequency)

		// Dampening
		freqArray[freqCounter] = float64(*frequency)
		freqCounter++
		freqCounter = freqCounter % freqArrayL
		var total float64 = 0
		for _, value:= range freqArray {
			total += value
		}

		// Interpolation
		hue := getHue(total/float64(len(freqArray)))
		old_hue := getHue(float64(*oldFreq))
		var interpInc float64 = (hue - old_hue) / interpNum

		for i := 1; i <= interpNum; i++ {
			r, g, b := boxColour(colorful.Hsv(old_hue + interpInc*float64(i), 1, 1)).RGB255()
			var newColour = r
			newColour = newColour << 8
			newColour = newColour + g
			newColour = newColour << 8
			newColour = newColour + b

			*current_colour_hex = newColour
			//fmt.Println(*current_colour_hex)
			coloured_square.QueueRedrawAll()
		}



	}
}

// Main
func main() {
	*current_colour_hex = rand.Uint32()
	ui.Main(setupUI)
}