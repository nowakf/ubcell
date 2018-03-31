package ubcell

import "time"

type EventResize struct{}

func (e *EventResize) When() time.Time {
	return time.Now()
}
