package ubcell

import (
	"image/color"

	"github.com/nowakf/pixel"
	"github.com/nowakf/pixel/pixelgl"
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

	Call(func(*pixelgl.Window))

	GetMatrix(x, y, w, h int) pixel.Matrix

	PollEvent() pixelgl.Event
	//to be merged:

	PostEvent()

	Show()
}

func NewScreen(cfg Config) (Screen, error) {
	if s, e := NewUBCellScreen(cfg); s != nil {
		return s, e

	} else {
		return nil, e
	}
}
