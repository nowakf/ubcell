package ubcell

import (
	"github.com/nowakf/pixel/pixelgl"
	"image/color"
)

type Screen interface {
	Init() error

	Fini()

	Clear()

	Fill(rune, color.RGBA)

	GetContent(x, y int) (ch rune, style Style)

	SetContent(x int, y int, ch rune, style Style)

	ShowCursor(x int, y int)

	HideCursor()

	Size() (width, height int)

	PollEvent() pixelgl.Event

	PostEvent()

	Show()
}

func NewScreen(p *pixelgl.Window, cfg Config) (Screen, error) {
	if s, e := NewUBCellScreen(p, cfg); s != nil {
		return s, nil

	} else {
		return nil, e
	}
}

type Config interface {
	FontSize() float64
	FontPath() string
	AdjustXY() (float64, float64)
}
