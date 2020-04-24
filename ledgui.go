package lcv

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"math/rand"
)

// Box wrapper which surrounds the square changing colour
var coloured_square *ui.Area

// Random colour generated to be fed to the area handler
var rand_color = rand.Uint32()

// The square which changes colour
var colored_area = areaHandler{area_color: &rand_color}

// The audio analyser which the ui uses
var aA = newAudioAnalyser(colored_area.changeColourUINT32, "starboy")

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
	coloured_square = ui.NewArea(colored_area)
	hbox.Append(coloured_square, true)

	visualise_button := ui.NewButton("start")
	visualise_button.OnClicked(func(b *ui.Button) {
		go aA.StartAnalysis()
	})
	hbox.Append(visualise_button, false)

	stop_button := ui.NewButton("stop")
	stop_button.OnClicked(func(b *ui.Button) {
		aA.u.stopSig <- true
	})
	hbox.Append(stop_button, false)

	mainwin.Show()
}
