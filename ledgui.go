package lcv

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/lucasb-eyer/go-colorful"
)

var coloured_square *ui.Area
var Current_colour_hex *uint32 = new(uint32)

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
	brush := mkSolidBrush(*Current_colour_hex, 1.0)
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

func changeColour(hue float64) {
	*Current_colour_hex = boxColour(colorful.Hsv(hue, 1, 1)).UINT32()
	coloured_square.QueueRedrawAll()
}

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
	coloured_square = ui.NewArea(areaHandler{})
	hbox.Append(coloured_square, true)

	visualise_button := ui.NewButton("start")
	visualise_button.OnClicked(func(b *ui.Button) {
		go StartPortAudio()
	})
	hbox.Append(visualise_button, false)

	stop_button := ui.NewButton("stop")
	stop_button.OnClicked(func(b *ui.Button) {
		sig <- true
	})
	hbox.Append(stop_button, false)

	mainwin.Show()
}
