package lcv

import (
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/lucasb-eyer/go-colorful"
	"io/ioutil"
	"math/rand"
	"sort"
	"strings"
)

// Box wrapper which surrounds the square changing colour
var coloured_square *ui.Area

// Random colour generated to be fed to the area handler
var rand_color = rand.Uint32()

// The handler for drawing a gradient area
var gh = &gradientareahandler{numcolours: 5}

// Checkbox which determines whether the custom gradient should be used
var cgbox *ui.Checkbox

// The combobox which houses all the hardcoded gradients
var gradientcbox *ui.Combobox

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

// Gradient handler struct which handles the drawing of blended gradients
type gradientareahandler struct {
	isreference bool
	gt          *GradientTable
	cboxes      [5]*ui.ColorButton
	sliders     [5]*ui.Slider
	numcolours  int
}

// For a the slider slice for "n" enabled slices, this returns the indexes
// of the sorted slices as to not sort the slice itself
func nMinPosition(s [5]*ui.Slider, n int) []int {
	num_slice := make([]int, 0)
	num_slice_ref := make([]int, 0)
	for i := 0; i < n; i++ {
		num_slice = append(num_slice, s[i].Value())
		num_slice_ref = append(num_slice_ref, s[i].Value())
	}
	sort.Ints(num_slice)

	indexes := make([]int, len(num_slice))
	for i := 0; i < len(num_slice); i++ {
		indexes[i] = intpos(num_slice_ref, num_slice[i])
	}

	return indexes
}

// Calculates the new gradient table for the area handler based on its
// sliders and comboboxes
func (gh *gradientareahandler) CalculateGradientTable() {
	if gh.cboxes[0] != nil {
		switch gh.numcolours {
		case 2:
			indxs := nMinPosition(gh.sliders, 2)

			r, g, b, _ := gh.cboxes[indxs[0]].Color()
			r2, g2, b2, _ := gh.cboxes[indxs[1]].Color()
			gh.gt = &GradientTable{
				{colorful.Color{r, g, b}, float64(gh.sliders[indxs[0]].Value()) / 10000},
				{colorful.Color{r2, g2, b2}, float64(gh.sliders[indxs[1]].Value()) / 10000},
			}
		case 3:
			indxs := nMinPosition(gh.sliders, 3)

			r, g, b, _ := gh.cboxes[indxs[0]].Color()
			fmt.Println(r, g, b)
			r2, g2, b2, _ := gh.cboxes[indxs[1]].Color()
			r3, g3, b3, _ := gh.cboxes[indxs[2]].Color()
			gh.gt = &GradientTable{
				{colorful.Color{r, g, b}, float64(gh.sliders[indxs[0]].Value()) / 10000},
				{colorful.Color{r2, g2, b2}, float64(gh.sliders[indxs[1]].Value()) / 10000},
				{colorful.Color{r3, g3, b3}, float64(gh.sliders[indxs[2]].Value()) / 10000},
			}
		case 4:
			indxs := nMinPosition(gh.sliders, 4)

			r, g, b, _ := gh.cboxes[indxs[0]].Color()
			r2, g2, b2, _ := gh.cboxes[indxs[1]].Color()
			r3, g3, b3, _ := gh.cboxes[indxs[2]].Color()
			r4, g4, b4, _ := gh.cboxes[indxs[3]].Color()
			gh.gt = &GradientTable{
				{colorful.Color{r, g, b}, float64(gh.sliders[indxs[0]].Value()) / 10000},
				{colorful.Color{r2, g2, b2}, float64(gh.sliders[indxs[1]].Value()) / 10000},
				{colorful.Color{r3, g3, b3}, float64(gh.sliders[indxs[2]].Value()) / 10000},
				{colorful.Color{r4, g4, b4}, float64(gh.sliders[indxs[3]].Value()) / 10000},
			}
		case 5:
			indxs := nMinPosition(gh.sliders, 5)

			r, g, b, _ := gh.cboxes[indxs[0]].Color()
			r2, g2, b2, _ := gh.cboxes[indxs[1]].Color()
			r3, g3, b3, _ := gh.cboxes[indxs[2]].Color()
			r4, g4, b4, _ := gh.cboxes[indxs[3]].Color()
			r5, g5, b5, _ := gh.cboxes[indxs[4]].Color()
			gh.gt = &GradientTable{
				{colorful.Color{r, g, b}, float64(gh.sliders[indxs[0]].Value()) / 10000},
				{colorful.Color{r2, g2, b2}, float64(gh.sliders[indxs[1]].Value()) / 10000},
				{colorful.Color{r3, g3, b3}, float64(gh.sliders[indxs[2]].Value()) / 10000},
				{colorful.Color{r4, g4, b4}, float64(gh.sliders[indxs[3]].Value()) / 10000},
				{colorful.Color{r5, g5, b5}, float64(gh.sliders[indxs[4]].Value()) / 10000},
			}
		}
	}
}

func (gh gradientareahandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	//fmt.Println("keep this print comment here")
	gh.CalculateGradientTable()

	for x := p.AreaWidth - 1; x >= 0; x-- {
		c := boxColour(gh.gt.GetInterpolatedColorFor(float64(x) / float64(p.AreaWidth))).UINT32()
		brush := mkSolidBrush(c, 1.0)

		path := ui.DrawNewPath(ui.DrawFillModeAlternate)
		path.AddRectangle(0, 0, float64(x+1), p.AreaHeight)
		path.End()
		p.Context.Fill(path, brush)
		path.Free()
	}

	if cgbox.Checked() {
		gh.CalculateGradientTable()
		aA.u.aaGT = gh.gt
		aA.u.gtUsed = true
	}
}

func (gradientareahandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	// do nothing
	//*current_colour_hex = rand.Uint32()
	//coloured_square.QueueRedrawAll()
}

func (gradientareahandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (gradientareahandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (gradientareahandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// do nothing
	return false
}

// Constructs the widgets for the visualisation page of the gui
func makeVisualisationPage() ui.Control {
	// Create the hbox for the visualiser
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	// Vbox for the program settings
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	// Adding the visualiser to the hbox
	coloured_square = ui.NewArea(colored_area)
	hbox.Append(coloured_square, true)

	// execution controls label
	vbox.Append(ui.NewLabel("main controls:"), false)

	// Button to start visualisation
	visualise_button := ui.NewButton("start")
	vbox.Append(visualise_button, false)

	// Button to stop visualisation
	stop_button := ui.NewButton("stop")
	vbox.Append(stop_button, false)

	// Gradient combobox
	vbox.Append(ui.NewLabel("gradients:"), false)
	gradientcbox = ui.NewCombobox()
	for _, name := range gradientList() {
		gradientcbox.Append(name)
	}
	gradientcbox.SetSelected(stringpos(gradientList(), "default"))
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
	optionshbox := ui.NewVerticalBox()
	hbox.SetPadded(true)
	vbox.Append(optionshbox, false)

	// Smoothing Checkbox
	smoothbox := ui.NewCheckbox("smoothing")
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
	dampbox := ui.NewCheckbox("dampening")
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

	// Custom Gradient Checkbox
	cgbox = ui.NewCheckbox("custom gradient")
	if aA.u.gtUsed {
		cgbox.SetChecked(true)
	}
	cgbox.OnToggled(func(c *ui.Checkbox) {
		if c.Checked() {
			gh.CalculateGradientTable()
			aA.u.aaGT = gh.gt
			aA.u.gtUsed = true
		} else {
			name := getGradientName(gradientcbox.Selected())
			aA.u.aaGT = gradients[name]
			if name != "default" {
				aA.u.gtUsed = true
			} else {
				aA.u.gtUsed = false
			}
		}
	})
	optionshbox.Append(cgbox, false)

	// Option to decide whether to send data over a udp connection to a localhost websocket
	vbox.Append(ui.NewLabel("connection options:"), false)
	udpledcntrl := ui.NewCheckbox("send data to led lights") // TODO underline or bold the heading text
	if aA.udph.shouldsend {
		udpledcntrl.SetChecked(true)
	}
	udpledcntrl.OnToggled(func(c *ui.Checkbox) {
		if c.Checked() {
			if !aA.udph.running {
				aA.udph.client.start()
				aA.udph.running = true
			}
			aA.udph.shouldsend = true
		} else {
			if aA.udph.running {
				aA.udph.client.closeConnection()
				aA.udph.running = false
			}
			aA.udph.shouldsend = false
		}
	})
	vbox.Append(udpledcntrl, false)

	// Defined here so the devicebox variable is in scope meaning it can be disabled on start of analysis
	visualise_button.OnClicked(func(b *ui.Button) {
		devicecbox.Disable()
		go aA.StartAnalysis()
	})
	stop_button.OnClicked(func(b *ui.Button) {
		aA.StopAnalysis()
		devicecbox.Enable()
		aA.udph.client.closeConnection()
		udpledcntrl.SetChecked(false)
	})

	return hbox
}

// Constructs the widgets for creating custom gradients
func makeGradientPage(mainwin *ui.Window) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	// Reference visualisation for the starboy gradient
	rh := &gradientareahandler{
		gt:          gradients["starboy"],
		numcolours:  5,
		isreference: true,
	}
	referencevis := ui.NewArea(rh)

	// The visualisation of the user created gradient
	gradientvis := ui.NewArea(gh)

	// Spinbox which controls the number of enabled colour boxes
	// when creating the gradient
	numberbox := ui.NewHorizontalBox()
	enabledspinbox := ui.NewSpinbox(2, 5)
	enabledspinbox.SetValue(3)
	gh.numcolours = enabledspinbox.Value()
	enabledspinbox.OnChanged(func(s *ui.Spinbox) {
		gh.numcolours = enabledspinbox.Value()
		for i, v := range gh.cboxes {
			if i+1 <= gh.numcolours {
				v.Enable()
				gh.sliders[i].Enable()
			} else {
				v.Disable()
				gh.sliders[i].Disable()
			}
		}

		gradientvis.QueueRedrawAll()
	})
	numberbox.Append(ui.NewLabel("number of colours to take into account:"), true)
	numberbox.Append(enabledspinbox, true)

	vbox.Append(referencevis, true)
	vbox.Append(gradientvis, true)

	// Creation of the colour boxes and their respective sliders which control blending position
	colourchoosingbox := ui.NewHorizontalBox()
	colour_picker_box := ui.NewVerticalBox()
	colour_mover_box := ui.NewVerticalBox()
	cbox1 := ui.NewColorButton()
	// Sliders set to 10000 but then when drawing are scaled back down to one as per the gradient
	// table struct. This allows the slider to have 10000 digits as opposed to simply 0 and 1
	slider1 := ui.NewSlider(0, 10000)
	cbox2 := ui.NewColorButton()
	slider2 := ui.NewSlider(0, 10000)
	cbox3 := ui.NewColorButton()
	slider3 := ui.NewSlider(0, 10000)
	cbox4 := ui.NewColorButton()
	slider4 := ui.NewSlider(0, 10000)
	cbox5 := ui.NewColorButton()
	slider5 := ui.NewSlider(0, 10000)

	cbox1.OnChanged(func(colorButton *ui.ColorButton) {
		gradientvis.QueueRedrawAll()
	})
	cbox2.OnChanged(func(colorButton *ui.ColorButton) {
		gradientvis.QueueRedrawAll()
	})
	cbox3.OnChanged(func(colorButton *ui.ColorButton) {
		gradientvis.QueueRedrawAll()
	})
	cbox4.OnChanged(func(colorButton *ui.ColorButton) {
		gradientvis.QueueRedrawAll()
	})
	cbox5.OnChanged(func(colorButton *ui.ColorButton) {
		gradientvis.QueueRedrawAll()
	})
	slider1.OnChanged(func(slider *ui.Slider) {
		gradientvis.QueueRedrawAll()
	})
	slider2.OnChanged(func(slider *ui.Slider) {
		gradientvis.QueueRedrawAll()
	})
	slider3.OnChanged(func(slider *ui.Slider) {
		gradientvis.QueueRedrawAll()
	})
	slider4.OnChanged(func(slider *ui.Slider) {
		gradientvis.QueueRedrawAll()
	})
	slider5.OnChanged(func(slider *ui.Slider) {
		gradientvis.QueueRedrawAll()
	})

	colour_picker_box.Append(cbox1, false)
	colour_mover_box.Append(slider1, false)
	colour_picker_box.Append(cbox2, false)
	colour_mover_box.Append(slider2, false)
	colour_picker_box.Append(cbox3, false)
	colour_mover_box.Append(slider3, false)
	colour_picker_box.Append(cbox4, false)
	colour_mover_box.Append(slider4, false)
	colour_picker_box.Append(cbox5, false)
	colour_mover_box.Append(slider5, false)

	// Populates the gradient handler with the correct colour boxes and sliders
	gh.cboxes = [5]*ui.ColorButton{cbox1, cbox2, cbox3, cbox4, cbox5}
	gh.sliders = [5]*ui.Slider{slider1, slider2, slider3, slider4, slider5}

	// Performs initial pass to enable/disable the correct number of sliders
	for i, v := range gh.cboxes {
		if i+1 <= gh.numcolours {
			v.Enable()
			gh.sliders[i].Enable()
		} else {
			v.Disable()
			gh.sliders[i].Disable()
		}
	}

	colour_picker_box.SetPadded(true)
	colour_mover_box.SetPadded(false)
	colourchoosingbox.Append(colour_picker_box, true)
	colourchoosingbox.Append(colour_mover_box, true)

	vbox.Append(colourchoosingbox, false)
	vbox.Append(numberbox, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	// Gradient saving/loading section
	gradsavingbox := ui.NewHorizontalBox()
	gradsavingbox.SetPadded(true)

	// Button to load a gradient from the file
	loadbtn := ui.NewButton("  load gradient from file  ")
	loadbtn.OnClicked(func(b *ui.Button) {
		filename := ui.OpenFile(mainwin)

		if filename != "" {
			file, err := ioutil.ReadFile(filename)
			chk(err)

			var unmarshaledtable *GradientTable
			err = json.Unmarshal(file, &unmarshaledtable)
			chk(err)

			// Based on the length of the gradient table, widgets which are disabled and reset
			// so only the data from the gradient is shown on the blended area handler
			enabledspinbox.SetValue(len(*unmarshaledtable))
			gh.numcolours = len(*unmarshaledtable)
			for i := 0; i < len(*unmarshaledtable); i++ {
				r, g, b := (*unmarshaledtable)[i].Col.RGB255()
				// RGB values must be scaled between 0 and 1
				gh.cboxes[i].SetColor(float64(r)/255, float64(g)/255, float64(b)/255, float64(1))
				gh.cboxes[i].Enable()
				gh.sliders[i].SetValue(int((*unmarshaledtable)[i].Pos * 10000))
				gh.sliders[i].Enable()
			}
			if len(*unmarshaledtable) < 5 {
				for i := len(*unmarshaledtable); i < 5; i++ {
					gh.cboxes[i].Disable()
					gh.sliders[i].Disable()
					gh.cboxes[i].SetColor(0, 0, 0, 1)
					gh.sliders[i].SetValue(0)
				}
			}

			gradientvis.QueueRedrawAll()
		}

	})
	savebtn := ui.NewButton("  save gradient to file  ")
	savebtn.OnClicked(func(b *ui.Button) {
		filename := ui.SaveFile(mainwin)

		i := strings.Index(filename, ".")
		if i != -1 {
			filename = strings.Split(filename, ".")[0]
		}
		filename += ".json"

		// Calculates the gradient table to ensure no logic error occurs
		gh.CalculateGradientTable()

		file, err := json.MarshalIndent(*gh.gt, "", " ")
		chk(err)
		err = ioutil.WriteFile(filename, file, 0644)
		chk(err)

	})
	gradsavingbox.Append(ui.NewLabel(""), true)
	gradsavingbox.Append(loadbtn, false)
	gradsavingbox.Append(savebtn, false)

	vbox.Append(gradsavingbox, false)

	return vbox
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

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Visualisation", makeVisualisationPage())
	tab.SetMargined(0, true)

	// The main window is passed for the open and save file dialogs
	tab.Append("Gradient Creator", makeGradientPage(mainwin))
	tab.SetMargined(1, true)

	mainwin.Show()
}
