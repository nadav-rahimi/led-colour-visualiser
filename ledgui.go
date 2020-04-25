package lcv

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"math/rand"
	"strings"
)

// Box wrapper which surrounds the square changing colour
var coloured_square *ui.Area

// Random colour generated to be fed to the area handler
var rand_color = rand.Uint32()

// The square which changes colour
var colored_area = areaHandler{area_color: &rand_color}

// The audio analyser which the ui uses
var aA = newAudioAnalyser(colored_area.changeColourUINT32, "")

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

// Custom areaHandler interface
type areaHandler struct {
	area_color *uint32
}

func (ah areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	// fill the area with white
	brush := mkSolidBrush(*ah.area_color, 1.0)
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

// Changes the colour of the canvas area based on a uint32 value
func (ah areaHandler) changeColourUINT32(c uint32) {
	*ah.area_color = c
	coloured_square.QueueRedrawAll()
}

// Initialises and constructs the UI window
func SetupUI() {
	// Create the main UI window
	mainwin := ui.NewWindow("LED Colour Visualiser", 625, 480, false)
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

	// Create the hbox for the visualiser
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	mainwin.SetChild(hbox)

	// Vbox for the program settings
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	// Adding the visualiser to the hbox
	coloured_square = ui.NewArea(colored_area)
	hbox.Append(coloured_square, true)

	// execution controls label
	vbox.Append(ui.NewLabel("start/stop controls:"), false)

	// Button to start visualisation
	visualise_button := ui.NewButton("start")
	vbox.Append(visualise_button, false)

	// Button to stop visualisation
	stop_button := ui.NewButton("stop")
	vbox.Append(stop_button, false)

	// Gradient combobox
	vbox.Append(ui.NewLabel("choose gradient:"), false)
	gradientcbox := ui.NewCombobox()
	for _, name := range gradientList() {
		gradientcbox.Append(name)
	}
	gradientcbox.SetSelected(pos(gradientList(), "default"))
	gradientcbox.OnSelected(func(c *ui.Combobox) {
		if getGradientName(gradientcbox.Selected()) == "default" {
			aA.u.gtUsed = false
		} else {
			aA.u.gtUsed = true
			aA.u.aaGT = gradients[getGradientName(gradientcbox.Selected())]
		}
	})
	vbox.Append(gradientcbox, false)

	// Audio device combobox
	vbox.Append(ui.NewLabel("audio device:"), false)
	devicecbox := ui.NewCombobox()
	for i, name := range getInputDevices() {
		devicecbox.Append(name)
		if strings.HasPrefix(name, "Line 1") {
			devicecbox.SetSelected(i)
		}
	}
	devicecbox.OnSelected(func(c *ui.Combobox) {
		aA.param.inputDeviceName = getInputDevices()[devicecbox.Selected()]
	})
	aA.param.inputDeviceName = getInputDevices()[devicecbox.Selected()]
	vbox.Append(devicecbox, false)

	// Options hbox
	vbox.Append(ui.NewLabel("visualisation options:"), false)
	optionshbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(optionshbox, false)

	// Smoothing Checkbox
	smoothbox := ui.NewCheckbox("Smoothing")
	if aA.param.smooth {
		smoothbox.SetChecked(true)
	}
	smoothbox.OnToggled(func(c *ui.Checkbox) {
		if c.Checked() {
			aA.param.smooth = true
		} else {
			aA.param.smooth = false
		}
	})
	optionshbox.Append(smoothbox, false)

	// Dampening Checkbox
	dampbox := ui.NewCheckbox("Dampening")
	if aA.param.damp {
		dampbox.SetChecked(true)
	}
	dampbox.OnToggled(func(c *ui.Checkbox) {
		if c.Checked() {
			aA.param.damp = true
		} else {
			aA.param.damp = false
		}
	})
	optionshbox.Append(dampbox, false)

	// Defined here so the devicebox variable is in scope meaning it can be disabled on start of analysis
	visualise_button.OnClicked(func(b *ui.Button) {
		devicecbox.Disable()
		go aA.StartAnalysis()
	})
	stop_button.OnClicked(func(b *ui.Button) {
		aA.StopAnalysis()
		devicecbox.Enable()
	})

	mainwin.Show()
}
