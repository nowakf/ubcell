package ubcell

import "github.com/nowakf/pixel/pixelgl"

type Config interface {
	GetDPI() float64
	GetFontSize() float64
	GetFontPath() string
	GetAdjustXY() (float64, float64)
	GetWindowConfig() pixelgl.WindowConfig
}
