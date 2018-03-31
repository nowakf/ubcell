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

	SetCell(x int, y int, style Style, ch ...rune)

	GetContent(x, y int) (ch rune, style Style, width int)

	// SetContent sets the contents of the given cell location.  If
	// the coordinates are out of range, then the operation is ignored.
	//
	// The first rune is the primary non-zero width rune.  The array
	// that follows is a possible list of combining characters to append,
	// and will usually be nil (no combining characters.)
	//
	// The results are not displayd until Show() or Sync() is called.
	//
	// Note that wide (East Asian full width) runes occupy two cells,
	// and attempts to place character at next cell to the right will have
	// undefined effects.  Wide runes that are printed in the
	// last column will be replaced with a single width space on output.
	SetContent(x int, y int, ch rune, style Style)

	// SetStyle sets the default style to use when clearing the screen
	// or when StyleDefault is specified.  If it is also StyleDefault,
	// then whatever system/terminal default is relevant will be used.

	// ShowCursor is used to display the cursor at a given location.
	// If the coordinates -1, -1 are given or are otherwise outside the
	// dimensions of the screen, the cursor will be hidden.
	ShowCursor(x int, y int)

	// HideCursor is used to hide the cursor.  Its an alias for
	// ShowCursor(-1, -1).
	HideCursor()

	// Size returns the screen size as width, height.  This changes in
	// response to a call to Clear or Flush.
	Size() (int, int)

	// PollEvent waits for events to arrive.  Main application loops
	// must spin on this to prevent the application from stalling.
	// Furthermore, this will return nil if the Screen is finalized.
	PollEvent() Event

	// PostEvent tries to post an event into the event stream.  This
	// can fail if the event queue is full.  In that case, the event
	// is dropped, and ErrEventQFull is returned.
	PostEvent(ev Event) error

	// PostEventWait is like PostEvent, but if the queue is full, it
	// blocks until there is space in the queue, making delivery
	// reliable.  However, it is VERY important that this function
	// never be called from within whatever event loop is polling
	// with PollEvent(), otherwise a deadlock may arise.
	//
	// For this reason, when using this function, the use of a
	// Goroutine is recommended to ensure no deadlock can occur.
	PostEventWait(ev Event)

	// EnableMouse enables the mouse.  (If your terminal supports it.)
	EnableMouse()

	// DisableMouse disables the mouse.
	DisableMouse()

	// HasMouse returns true if the terminal (apparently) supports a
	// mouse.  Note that the a return value of true doesn't guarantee that
	// a mouse/pointing device is present; a false return definitely
	// indicates no mouse support is available.
	HasMouse() bool

	// Show makes all the content changes made using SetContent() visible
	// on the display.
	//
	// It does so in the most efficient and least visually disruptive
	// manner possible.
	Show()

	// Sync works like Show(), but it updates every visible cell on the
	// physical display, assuming that it is not synchronized with any
	// internal model.  This may be both expensive and visually jarring,
	// so it should only be used when believed to actually be necessary.
	//
	// Typically this is called as a result of a user-requested redraw
	// (e.g. to clear up on screen corruption caused by some other program),
	// or during a resize event.
	Sync()
}

func NewScreen(p *pixelgl.Window) (Screen, error) {
	// First we attempt to obtain a terminfo screen.  This should work
	// in most places if $TERM is set.
	if s, e := NewUBCellScreen(p); s != nil {
		return s, nil

	} else {
		return nil, e
	}
}
