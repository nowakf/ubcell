package ubcell

type LineBuffer struct {
	w     int
	h     int
	lines []line
}

type line []rune

// SetContent sets the contents (primary rune, combining runes,
// and style) for a cell at a given location.
func (l *LineBuffer) SetContent(x int, y int,
	mainc rune, combc []rune, style Style) {
	if x >= 0 && y >= 0 && x < l.w && y < l.h {
		l.lines[y][x] = mainc
	}
}

// GetContent returns the contents of a character cell, including the
// primary rune, any combining character runes (which will usually be
// nil), the style, and the display width in cells.  (The width can be
// either 1, normally, or 2 for East Asian full-width characters.)
func (l *LineBuffer) GetContent(x, y int) (rune, []rune, Style, int) {
	var mainc rune
	var combc []rune
	var style Style
	var width int
	if x >= 0 && y >= 0 && x < l.w && y < l.h {
		mainc = l.lines[y][x]
	}
	return mainc, combc, style, width
}

// Size returns the (width, height) in cells of the buffer.
func (l *LineBuffer) Size() (int, int) {
	return l.h, l.w
}

// Resize is used to resize the cells array, with different dimensions,
// while preserving the original contents.  The cells will be invalidated
// so that they can be redrawn.
func (l *LineBuffer) Resize(w, h int) {

}

// Fill fills the entire cell buffer array with the specified character
// and style.  Normally choose ' ' to clear the screen.  This API doesn't
// support combining characters.
func (l *LineBuffer) Fill(r rune, style Style) {
}
