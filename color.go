package ubcell

import (
	"golang.org/x/image/colornames"
	"image/color"
)

func GetColor(name string) color.RGBA {
	if c, ok := colornames.Map[name]; ok {
		return c
	} else {
		return colornames.Mediumspringgreen
	}
}
