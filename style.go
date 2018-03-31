package ubcell

import "image/color"

type Style struct {
	Background color.RGBA
	Foreground color.RGBA
}

var (
	StyleDefault = Style{}
)
