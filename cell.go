package ubcell

import "image/color"
import "sort"

type CellBuffer struct {
	h       int
	w       int
	cells   []cell
	changes sort.IntSlice

	draw  func(int, int, rune)
	ink   func(color.RGBA)
	clear func()
	done  func()
}

type cell struct {
	ch    rune
	style Style
}

func NewCellBuffer(h, w int, drawer func(int, int, rune), ink func(color.RGBA), clear func(), done func()) *CellBuffer {
	c := new(CellBuffer)
	c.h, c.w = h, w
	c.cells = make([]cell, h*w)
	c.draw = drawer
	c.ink = ink
	c.clear = clear
	c.done = done
	c.changes = make(sort.IntSlice, 0)
	return c

}

// SetContent sets the contents (primary rune, combining runes,
// and style) for a cell at a given location.
func (c *CellBuffer) SetContent(x int, y int, ch rune, style Style) {

	if x < 0 || y < 0 || x >= c.w || y >= c.h {
		return
	}
	loc := (c.w * y) + x

	ce := &c.cells[loc]

	ce.style = style
	ce.ch = ch
	ind := c.changes.Search(loc)
	if ind < len(c.changes) {
		c.changes[ind] = loc
	} else {
		c.changes = append(c.changes, loc)
	}

}

// GetContent returns the contents of a character cell
func (c *CellBuffer) GetContent(x, y int) (rune, Style) {
	if x < 0 || y < 0 || x <= c.w || y <= c.h {
		cel := &c.cells[(c.w*y)+x]
		return cel.ch, cel.style
	} else {
		return '?', StyleDefault
	}

}

// Size returns the (height, width) in cells of the buffer.
func (c *CellBuffer) Size() (int, int) {
	return c.h, c.w
}

func (c *CellBuffer) Resize(w, h int) {
	if c.h == h && c.w == w {
		return
	}

	newc := make([]cell, w*h)
	//	for y := 0; y < h && y < c.h; y++ {
	//		for x := 0; x < w && x < c.w; x++ {
	//			oc := &c.cells[(y*c.w)+x]
	//			nc := &newc[(y*w)+x]
	//
	//			nc.style = oc.style
	//			nc.ch = oc.ch
	//
	//		}
	//
	//	}
	c.w = w
	c.h = h
	c.cells = newc
	//
}
func (c *CellBuffer) Seperate(changes sort.IntSlice) (foregrounds map[color.RGBA][]int, backgrounds map[color.RGBA][]int) {

	foregrounds = make(map[color.RGBA][]int)
	backgrounds = make(map[color.RGBA][]int)

	changes.Sort()

	for i := 0; i < changes.Len(); i++ {
		cell := c.cells[changes[i]]
		fg, bg := cell.style.Decompose()
		foregrounds[fg] = append(foregrounds[fg], changes[i])
		backgrounds[bg] = append(backgrounds[bg], changes[i])
	}
	c.changes = make(sort.IntSlice, 0)
	return foregrounds, backgrounds
}

func (c *CellBuffer) Draw() {

	fgs, bgs := c.Seperate(c.changes)

	if len(fgs)+len(bgs) == 0 {
		return
	}

	for color, list := range bgs {
		c.ink(color)
		for _, index := range list {
			c.draw(index%c.w, index/c.w, '█')
		}
		c.done()
	}

	for color, list := range fgs {
		c.ink(color)
		for _, index := range list {
			c.draw(index%c.w, index/c.w, c.cells[index].ch)
		}
		c.done()
	}

}

func (c *CellBuffer) Clear() {
	c.clear()
}

func (c *CellBuffer) Init() error {
	panic("not implemented")
}

func (c *CellBuffer) Fini() {
	panic("not implemented")
}

func (c *CellBuffer) Fill(rune, color.RGBA) {
	panic("not implemented")
}

func (c *CellBuffer) ShowCursor(x int, y int) {
	panic("not implemented")
}

func (c *CellBuffer) HideCursor() {
	panic("not implemented")
}

func (c *CellBuffer) PollEvent() interface{} {
	panic("not implemented")
}

func (c *CellBuffer) Show() {
	panic("not implemented")
}
