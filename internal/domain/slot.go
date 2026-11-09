// Package domain holds the booking rules. Pure Go: no HTTP, no storage.
package domain

import "time"

const (
	// MaxDaysAhead is how far ahead a slot can be booked, inclusive.
	MaxDaysAhead = 30
	// Opening hours, as offsets from the start of the day. Same every day.
	OpensAt  = 8 * time.Hour
	ClosesAt = 19 * time.Hour
	// Grid is the step that start and end times fall on.
	Grid = 15 * time.Minute
	// MinLength is the shortest slot. Longer slots grow in Grid steps,
	// which onGrid already guarantees.
	MinLength = 30 * time.Minute
)

// Slot is a time slot on one day, in the cafe's time zone.
type Slot struct {
	Start time.Time
	End   time.Time
}

// SelectableAt reports whether the slot can be selected at the given moment.
func (s Slot) SelectableAt(now time.Time) bool {
	return s.withinBookingWindow(now) && s.withinOpeningHours() && s.onGrid() && s.longEnough() &&
		s.notStartedAt(now)
}

// notStartedAt: from the start time onwards, the slot can no longer be selected.
func (s Slot) notStartedAt(now time.Time) bool {
	return now.Before(s.Start)
}

func (s Slot) longEnough() bool {
	return s.End.Sub(s.Start) >= MinLength
}

// onGrid: start and end are whole quarters of an hour.
func (s Slot) onGrid() bool {
	return isOnGrid(s.Start) && isOnGrid(s.End)
}

func isOnGrid(t time.Time) bool {
	return t.Sub(startOfDay(t))%Grid == 0
}

// withinOpeningHours: starts at or after opening, ends at or before closing.
func (s Slot) withinOpeningHours() bool {
	day := startOfDay(s.Start)
	return !s.Start.Before(day.Add(OpensAt)) && !s.End.After(day.Add(ClosesAt))
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
