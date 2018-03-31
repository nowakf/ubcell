package ubcell

import (
	"github.com/nowakf/pixel/pixelgl"
	"time"
)

type EventKey struct {
	t   time.Time
	mod pixelgl.ModKey
	key pixelgl.Key
	ch  rune
	s   string
}

func (ev *EventKey) When() time.Time {
	return ev.t
}
func (ev *EventKey) Rune() rune {
	return ev.ch
}
func (ev *EventKey) String() string {
	return ev.s
}
func (ev *EventKey) Key() pixelgl.Key {
	return ev.key
}
func (ev *EventKey) Modifiers() pixelgl.ModKey {
	return ev.mod
}
func (ev *EventKey) Name() string {
	return "not implemented"
}
func (ev *EventKey) NewEventKey(k pixelgl.Key, s string, mod pixelgl.ModKey) *EventKey {
	return &EventKey{t: time.Now(), key: k, s: s, mod: mod}
}
