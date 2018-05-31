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
	"unicode"
)

const dpi = 72

const ADJUSTX = -2
const ADJUSTY = -2

func NewUBCellScreen(p *pixelgl.Window, path string) (Screen, error) {

	u := new(ubcellScreen)
	u.win = p

	f, err := os.Open(path)
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
		Size: 22,
		DPI:  dpi,
	})

	atlas := text.NewAtlas(
		face,
		text.RangeTable(unicode.Po),
		text.RangeTable(unicode.S),
		text.ASCII,
	)

	u.hinc = atlas.Glyph('█').Frame.H() + ADJUSTY
	u.winc = atlas.Glyph('█').Frame.W() + ADJUSTX

	u.h = int(fh / u.hinc)
	u.w = int(fw / u.winc)

	u.t = text.New(pixel.ZV, atlas)

	u.cells = NewCellBuffer(u.h, u.w,
		func(x, y int, r rune) {
			u.t.Add(r, pixel.V(u.winc*float64(x), fh-u.hinc*float64(y)-u.hinc))
		},
		func(c color.RGBA) {
			//change the color
			u.t.Ink(c)
		},
		func() {
			u.win.Clear(u.backgroundC)
		},
		func() {
			u.t.Apply()
			u.t.Draw(p, pixel.IM)
			u.t.Clear()
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

	backgroundC pixel.RGBA

	sync.Mutex
}

func (u *ubcellScreen) Init() error {
	return nil
}

func (u *ubcellScreen) Fini() {
}
func (u *ubcellScreen) Clear() {
	u.Lock()
	defer u.Unlock()
	u.win.Clear(u.backgroundC)
}
func (u *ubcellScreen) Fill(r rune, col color.RGBA) {
}

func (u *ubcellScreen) GetContent(x, y int) (ch rune, style *Style) {
	u.Lock()
	defer u.Unlock()
	ch, sty := u.cells.GetContent(x, y)
	return ch, sty
}

func (u *ubcellScreen) SetContent(x, y int, ch rune, style *Style) {
	u.Lock()
	defer u.Unlock()
	u.cells.SetContent(x, y, ch, style)
}

func (u *ubcellScreen) Cat(r rune) (names []string) {
	names = make([]string, 0)
	for name, table := range unicode.Categories {
		if unicode.Is(table, r) {
			names = append(names, name)
		}
	}
	return
}
func (u *ubcellScreen) ShowCursor(x, y int) {}
func (u *ubcellScreen) HideCursor()         {}
func (u *ubcellScreen) Size() (int, int) {
	return u.w, u.h
}
func (u *ubcellScreen) PollEvent() pixelgl.Event {
	return u.win.PollEvent()
}
func (u *ubcellScreen) PostEvent() error {
	u.win.PostEmpty()
	return nil
}
func (u *ubcellScreen) Show() {
	u.Lock()
	defer u.Unlock()
	u.cells.Draw()
}
