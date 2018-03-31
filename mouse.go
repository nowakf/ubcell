package ubcell

import "time"

type EventMouse struct{}

func (e *EventMouse) When() time.Time {
	return time.Now()
}
