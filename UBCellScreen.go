package ubcell

import (
	"image/color"
	"io/ioutil"
	"os"
	"sync"
	"unicode"

	"github.com/golang/freetype/truetype"
	"github.com/nowakf/pixel"
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/pixel/text"
)

func NewUBCellScreen(cfg Config) (Screen, error) {
	u := new(ubcellScreen)

	win, err := u.window()

	if err != nil {
		return u, err
	}

	f, err := os.Open(cfg.FontPath())

	defer f.Close()

	if err != nil {
		return u, err
	}

	bytes, err := ioutil.ReadAll(f)

	if err != nil {
		return u, err
	}

	tt, err := truetype.Parse(bytes)

	if err != nil {
		return u, err
	}

	face := truetype.NewFace(tt, &truetype.Options{
		Size: cfg.FontSize(),
		DPI:  cfg.DPI(),
	})

	atlas := text.NewAtlas(
		face,
		text.RangeTable(unicode.Po),
		text.RangeTable(unicode.S),
		text.ASCII,
	)
	winc, hinc := u.glyphBounds(cfg, atlas)

	w, h := u.size(win, winc, hinc)

	t := text.New(pixel.ZV, atlas)

	cells := NewCellBuffer(
		h, w,
		func(x, y int, r rune) {
			u.t.Add(r, pixel.V(u.winc*float64(x), u.win.Bounds().H()-u.hinc*float64(y)-u.hinc))
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
			u.t.Draw(u.win, pixel.IM)
			u.t.Clear()
		})

	u = &ubcellScreen{
		win:   win,
		cfg:   cfg,
		cells: cells,
		t:     t,
		h:     h,
		w:     w,
		hinc:  hinc,
		winc:  winc,
	}

	go u.loop()

	return u, err
}

type ubcellScreen struct {
	t     *text.Text
	atlas *text.Atlas

	win *pixelgl.Window

	h, w       int
	hinc, winc float64
	cells      *CellBuffer

	cfg Config

	backgroundC pixel.RGBA

	sync.Mutex
	paused bool
}

func (u *ubcellScreen) window() (*pixelgl.Window, error) {

	u.Lock()
	defer u.Unlock()

	cfg := pixelgl.WindowConfig{
		Title:     "testing",
		Bounds:    pixel.R(0, 0, 1024, 1024),
		Resizable: true,
		//Monitor:   pixelgl.PrimaryMonitor(),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)

	win.SetSmooth(false)

	return win, err
}
func (u *ubcellScreen) loop() {
	for !u.win.Closed() {
		if !u.paused {
			u.Lock()
			u.win.Update()
			u.Unlock()
		} else {
			u.win.UpdateInput()
		}
	}

}
func (u *ubcellScreen) Init() error {

	return nil
}
func (u *ubcellScreen) PostEvent() {

	u.win.PostEmpty()
}
func (u *ubcellScreen) Fini() {

}
func (u *ubcellScreen) Clear() {

	w, h := u.size(u.win, u.winc, u.hinc)
	u.cells.Resize(w, h)
	u.Show()

}

func (u *ubcellScreen) glyphBounds(cfg Config, atlas *text.Atlas) (width float64, height float64) {

	xAdjust, yAdjust := cfg.AdjustXY()

	hinc := atlas.Glyph('█').Frame.H() + xAdjust
	winc := atlas.Glyph('█').Frame.W() + yAdjust
	return winc, hinc
}

func (u *ubcellScreen) Fill(r rune, col color.RGBA) {

}

func (u *ubcellScreen) GetContent(x, y int) (ch rune, style Style) {

	if u.cells != nil {
		ch, style = u.cells.GetContent(x, y)
	} else {
		ch = '!'
		style = StyleDefault
	}
	return ch, style
}

func (u *ubcellScreen) SetContent(x, y int, ch rune, style Style) {
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
func (u *ubcellScreen) ShowCursor(x, y int) {

}
func (u *ubcellScreen) HideCursor() {

}

func (u *ubcellScreen) size(win *pixelgl.Window, winc, hinc float64) (width int, height int) {

	var fw, fh float64
	u.Lock()
	fw, fh = win.Bounds().W(), win.Bounds().H()
	u.Unlock()
	return int(fw / winc), int(fh / hinc)
}
func (u *ubcellScreen) Size() (int, int) {

	width, height := u.size(u.win, u.winc, u.hinc)
	return width, height
}
func (u *ubcellScreen) PollEvent() pixelgl.Event {
	u.paused = true
	ev := u.win.PollEvent()
	u.paused = false
	switch ev.(type) {
	case *pixelgl.CursorEvent:

		//we need to convert it
		//resize should be here, too
		return ev
	default:
		return ev

	}

}
func (u *ubcellScreen) Show() {
	//this should block until it's done...
	u.Lock()
	u.cells.Draw()
	u.Unlock()
}
