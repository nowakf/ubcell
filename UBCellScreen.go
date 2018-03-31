package ubcell

import (
	"github.com/golang/freetype/truetype"
	"github.com/nowakf/pixel/pixelgl"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"os"
	"sync"
)

const dpi = 72

func NewUBCellScreen(p *pixelgl.Window) (Screen, error) {

	u := new(ubcellScreen)

	f, err := os.Open("./assets/fonts/DejaVuSansMono.ttf")
	defer f.Close()

	if err != nil {
		return u, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return u, err
	}

	fw, fh := p.Bounds().W(), p.Bounds().H()
	tt, err := truetype.Parse(bytes)
	if err != nil {
		return u, err
	}
	face := truetype.NewFace(tt, &truetype.Options{
		Size: 12,
		DPI:  dpi,
	})
	bounds, _, _ := face.GlyphBounds('S')
	glyphHeight := bounds.Max.Y - bounds.Min.Y
	glyphWidth := bounds.Max.X - bounds.Min.X
	u.h = int(float64(fh) / float64(glyphHeight))
	u.w = int(float64(fw) / float64(glyphWidth))
	//you can hand it a reference to the backing array here?
	return u, err
}

type ubcellScreen struct {
	face  *font.Face
	lines LineBuffer
	h, w  int
	sync.Mutex
}

func (u *ubcellScreen) Init() error {
	return nil
}
func (u *ubcellScreen) inputLoop() {}
func (u *ubcellScreen) mainLoop()  {}
func (u *ubcellScreen) draw()      {}
func (u *ubcellScreen) Fini() {
}
func (u *ubcellScreen) Clear()                      {}
func (u *ubcellScreen) Fill(r rune, col color.RGBA) {}
func (u *ubcellScreen) SetCell(x, y int, style Style, ch ...rune) {
	//this should never be called
	panic("called SetCell!")
}
func (u *ubcellScreen) GetContent(x, y int) (ch rune, style Style, width int) {
	return ' ', StyleDefault, y
}

func (u *ubcellScreen) SetContent(x, y int, ch rune, style Style) {

}
func (u *ubcellScreen) showCursor()         {}
func (u *ubcellScreen) ShowCursor(x, y int) {}
func (u *ubcellScreen) HideCursor()         {}
func (u *ubcellScreen) Size() (int, int) {
	return u.lines.Size()
}
func (u *ubcellScreen) PollEvent() Event {
	return nil
}
func (u *ubcellScreen) PostEvent(ev Event) error {
	return nil
}
func (u *ubcellScreen) PostEventWait(ev Event) {}
func (u *ubcellScreen) EnableMouse()           {}
func (u *ubcellScreen) DisableMouse()          {}
func (u *ubcellScreen) HasMouse() bool {
	return false
}
func (u *ubcellScreen) Show() {}
func (u *ubcellScreen) Sync() {}
