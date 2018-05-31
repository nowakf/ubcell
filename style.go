package ubcell

import (
	"golang.org/x/image/colornames"
	"image/color"
)

type Style struct {
	background color.RGBA
	foreground color.RGBA
}

func (s Style) Decompose() (color.RGBA, color.RGBA) {
	return s.foreground, s.background
}
func (s *Style) Foreground(c color.RGBA) *Style {
	s.foreground = c
	return s
}

func (s *Style) Background(c color.RGBA) *Style {
	s.background = c
	return s
}

var (
	StyleDefault = Style{colornames.Grey, colornames.White}
)
