package ubcell

import (
	"image/color"
	"io/ioutil"
	"os"
	"time"
	"unicode"

	"github.com/golang/freetype/truetype"
	"github.com/nowakf/pixel"
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/pixel/text"
)

func NewUBCellScreen(cfg Config) (Screen, error) {

	var err error

	u := new(ubcellScreen)

	u.win, err = u.window(cfg.GetWindowConfig())

	u.call = make(chan func(), 10)

	err = handle(err)

	u.cfg = cfg

	var atlas *text.Atlas

	atlas, err = u.makeAtlas(cfg.GetFontPath())

	err = handle(err)

	u.t = text.New(pixel.ZV, atlas)

	u.winc, u.hinc = u.glyphBounds(cfg, atlas)

	u.w, u.h = u.size(u.win, u.winc, u.hinc)

	u.cells = NewCellBuffer(u.h, u.w,
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
	return u, err

}

type ubcellScreen struct {
	t *text.Text

	win *pixelgl.Window

	h, w       int
	hinc, winc float64
	cells      *CellBuffer

	cfg Config

	backgroundC pixel.RGBA

	paused bool
	call   chan func()
}

func (u *ubcellScreen) size(win *pixelgl.Window, winc, hinc float64) (width int, height int) {

	var fw, fh float64
	fw, fh = win.Bounds().W(), win.Bounds().H()
	return int(fw / winc), int(fh / hinc)
}
func handle(err error) error {
	//eventually do something with the errors
	return err
}

func (u *ubcellScreen) window(cfg pixelgl.WindowConfig) (*pixelgl.Window, error) {

	win, err := pixelgl.NewWindow(cfg)

	win.SetSmooth(false)

	return win, handle(err)
}
func (u *ubcellScreen) makeAtlas(path string) (*text.Atlas, error) {
	f, err := os.Open(path)

	defer f.Close()

	err = handle(err)

	bytes, err := ioutil.ReadAll(f)

	err = handle(err)

	tt, err := truetype.Parse(bytes)

	err = handle(err)

	face := truetype.NewFace(tt, &truetype.Options{
		Size: u.cfg.GetFontSize(),
		DPI:  u.cfg.GetDPI(),
	})

	atlas := text.NewAtlas(
		face,
		text.RangeTable(unicode.Po),
		text.RangeTable(unicode.S),
		text.ASCII,
	)
	return atlas, err
}

func (u *ubcellScreen) loop() {
	fps := time.NewTicker(time.Second / 60)

	for !u.win.Closed() {
		select {
		case f := <-u.call:
			f()
			u.win.UpdateGraphics()
		case <-fps.C:
			u.win.UpdateInput()
		}
	}

	fps.Stop()

}
func (u *ubcellScreen) Call(f func(*pixelgl.Window)) {
	go func() {
		u.call <- func() {
			f(u.win)
		}
	}()
}
func (u *ubcellScreen) GetMatrix(x, y, w, h int) pixel.Matrix {
	return pixel.IM.Moved(u.win.Bounds().Center())
}

//things to do after the screen loads:
func (u *ubcellScreen) Init() error {
	go u.loop()
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

	xAdjust, yAdjust := cfg.GetAdjustXY()

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

func (u *ubcellScreen) Size() (int, int) {

	width, height := u.size(u.win, u.winc, u.hinc)
	return width, height
}
func (u *ubcellScreen) PollEvent() pixelgl.Event {
	return <-u.win.EventChannel

}

func (u *ubcellScreen) Show() {
	go func() {
		if u.call != nil {
			u.call <- u.cells.Draw
		}
	}()
}
