// Package domain holds the booking rules. Pure Go: no HTTP, no storage.
package domain

import "time"

// MaxDaysAhead is how far ahead a slot can be booked, inclusive.
const MaxDaysAhead = 30

// Slot is a time slot on one day, in the cafe's time zone.
type Slot struct {
	Start time.Time
	End   time.Time
}

// SelectableAt reports whether the slot can be selected at the given moment.
func (s Slot) SelectableAt(now time.Time) bool {
	return s.withinBookingWindow(now)
}

// withinBookingWindow: the slot's day is between today and today+30, inclusive,
// counted in calendar days in the cafe's time zone.
func (s Slot) withinBookingWindow(now time.Time) bool {
	day := startOfDay(s.Start)
	today := startOfDay(now.In(s.Start.Location()))
	lastDay := today.AddDate(0, 0, MaxDaysAhead)
	return !day.Before(today) && !day.After(lastDay)
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
