package ubcell

import "image/color"

type CellBuffer struct {
	h       int
	w       int
	cells   []cell
	changes *changeBuffer
	draw    func(int, int, rune)
	ink     func(color.RGBA)
	done    func()
}

type cell struct {
	ch    rune
	style Style
}

type changeBuffer struct {
	start, end int
	underlying []int
}

func ChangeBuffer(h, w int) *changeBuffer {
	arr := make([]int, h*w)
	return &changeBuffer{0, 0, arr[:]}
}
func (c *changeBuffer) Push(change int) {
	c.end++
	if c.end >= len(c.underlying) {
		c.end = 0
	}
	c.underlying[c.end] = change

}

func (c *changeBuffer) Pop() int {
	if c.start >= c.end {
		return -1
	}
	c.start++
	if c.start >= len(c.underlying) {
		c.start = 0
	}
	return c.underlying[c.start]
}
func (c *changeBuffer) Len() int {
	return c.end - c.start
}

func (c *changeBuffer) Seperate() {
}

func (c *changeBuffer) Resize(h, w int) {
	underlying := make([]int, h*w)
	copy(underlying, c.underlying[c.start:c.end])
}

func NewCellBuffer(h, w int, drawer func(int, int, rune), ink func(color.RGBA), done func()) *CellBuffer {
	c := new(CellBuffer)
	c.h, c.w = h, w
	c.cells = make([]cell, h*w)
	c.draw = drawer
	c.ink = ink
	c.done = done
	c.changes = ChangeBuffer(h, w)
	return c

}

// SetContent sets the contents (primary rune, combining runes,
// and style) for a cell at a given location.
func (c *CellBuffer) SetContent(x int, y int, ch rune, style Style) {

	if x <= 0 && y <= 0 && x > c.w && y > c.h {
		return
	}

	ce := &c.cells[(c.w*y)+x]

	if ce.ch != ch || ce.style != style {
		ce.style = style
		ce.ch = ch
		c.changes.Push(x + (y * c.w))
	}

}

// GetContent returns the contents of a character cell
func (c *CellBuffer) GetContent(x, y int) (rune, Style) {
	if x >= 0 && y >= 0 && x < c.w && y < c.h {
		cel := &c.cells[(c.w*y)+x]
		return cel.ch, cel.style
	} else {
		return rune(0), Style{}
	}

}

// Size returns the (width, height) in cells of the buffer.
func (c *CellBuffer) Size() (int, int) {
	return c.h, c.w
}

func (c *CellBuffer) Resize(h, w int) {
	if c.h == h && c.w == w {
		return
	}
	newc := make([]cell, w*h)
	for y := 0; y < h && y < c.h; y++ {
		for x := 0; x < w && x < c.w; x++ {
			oc := &c.cells[(y*c.w)+x]
			nc := &newc[(y*w)+x]

			nc.style = oc.style
			nc.ch = oc.ch

		}

	}
	//
}
func (c *CellBuffer) Seperate(changes *changeBuffer) (foregrounds map[color.RGBA][]int, backgrounds map[color.RGBA][]int) {
	foregrounds = make(map[color.RGBA][]int)
	backgrounds = make(map[color.RGBA][]int)
	for i := 0; i < changes.Len(); i++ {
		cinc := changes.Pop()
		fg := c.cells[cinc].style.Foreground
		foregrounds[fg] = append(foregrounds[fg], cinc)
		bg := c.cells[cinc].style.Background
		backgrounds[bg] = append(backgrounds[bg], cinc)
	}
	return foregrounds, backgrounds
}
func (c *CellBuffer) Draw() {
	fgs, bgs := c.Seperate(c.changes)
	for color, list := range bgs {
		c.ink(color)
		for _, index := range list {
			c.draw(index%c.w, index/c.w, ' ')
		}
	}
	for color, list := range fgs {
		c.ink(color)
		for _, index := range list {
			c.draw(index%c.w, index/c.w, c.cells[index].ch)
		}
	}

	c.done()
}
