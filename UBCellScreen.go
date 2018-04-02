package ubcell

import (
	"github.com/golang/freetype/truetype"
	"github.com/nowakf/pixel"
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/pixel/text"
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

	atlas := text.NewAtlas(face, text.ASCII)

	u.hinc = atlas.Glyph('S').Frame.H()
	u.winc = atlas.Glyph('S').Frame.W()
	u.h = int(fh / u.hinc)
	u.w = int(fw / u.winc)

	u.t = text.New(pixel.ZV, atlas)

	u.cells = NewCellBuffer(u.h, u.w,
		func(x, y int, r rune) {
			u.t.Add(r, pixel.V(u.winc*float64(x), u.hinc*float64(y)))
		},
		func(c color.RGBA) {
			//change the color
			u.t.Ink(c)
		},
		func() {
			u.t.Apply()
			u.t.Draw(p, pixel.IM)
		})

	//you can hand it a reference to the backing array here?
	return u, err
}

type ubcellScreen struct {
	t   *text.Text
	win *pixelgl.Window

	h, w       int
	hinc, winc float64
	cells      *CellBuffer
	sync.Mutex
}

func (u *ubcellScreen) Init() error {
	return nil
}
func (u *ubcellScreen) inputLoop() {}
func (u *ubcellScreen) mainLoop() {
	u.cells.Draw()
}
func (u *ubcellScreen) Fini() {
}
func (u *ubcellScreen) Clear()                      {}
func (u *ubcellScreen) Fill(r rune, col color.RGBA) {}
func (u *ubcellScreen) SetCell(x, y int, style Style, ch ...rune) {
	//this should never be called
	panic("called SetCell!")
}
func (u *ubcellScreen) GetContent(x, y int) (ch rune, style Style, width int) {
	return u.cells.GetContent(x, y), 1
}

func (u *ubcellScreen) SetContent(x, y int, ch rune, style Style) {
	u.cells.SetContent(x, y, ch, style)
}
func (u *ubcellScreen) showCursor()         {}
func (u *ubcellScreen) ShowCursor(x, y int) {}
func (u *ubcellScreen) HideCursor()         {}
func (u *ubcellScreen) Size() (int, int) {
	return u.h, u.w
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
