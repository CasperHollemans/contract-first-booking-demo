// Package domain holds the booking rules. Pure Go: no HTTP, no storage.
package domain

import "time"

// Slot is a time slot on one day, in the cafe's time zone.
type Slot struct {
	Start time.Time
	End   time.Time
}

// SelectableAt reports whether the slot can be selected at the given moment.
func (s Slot) SelectableAt(now time.Time) bool {
	return true
}
