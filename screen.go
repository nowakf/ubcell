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

	GetContent(x, y int) (ch rune, style *Style)

	SetContent(x int, y int, ch rune, style *Style)

	ShowCursor(x int, y int)

	HideCursor()

	Size() (int, int)

	PollEvent() pixelgl.Event

	Show()
}

func NewScreen(p *pixelgl.Window, path string) (Screen, error) {
	if s, e := NewUBCellScreen(p, path); s != nil {
		return s, nil

	} else {
		return nil, e
	}
}
