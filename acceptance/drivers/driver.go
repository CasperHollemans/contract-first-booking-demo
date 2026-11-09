// Package drivers is the only acceptance-test layer that knows how to reach
// the system: it starts the real server on a real port (with a fake clock)
// and talks to it over HTTP (API driver) or through headless Chrome (UI driver).
package drivers

import "time"

// Driver is what the DSL needs from any way of reaching the system.
// Dates are "2006-01-02", times are "15:04".
type Driver interface {
	// SetNow sets the system's current time.
	SetNow(now time.Time)
	// CheckSlot asks whether the slot can be selected.
	CheckSlot(date, start, end string) (selectable bool, err error)
}
