package domain

import (
	"testing"
	"time"
)

// Helpers: all times are in the cafe's time zone (UTC for the demo).
func at(day, hhmm string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04", day+" "+hhmm, time.UTC)
	if err != nil {
		panic(err)
	}
	return t
}

func slot(day, start, end string) Slot {
	return Slot{Start: at(day, start), End: at(day, end)}
}

const monday = "2026-11-09"

type selectCase struct {
	name string
	now  time.Time
	slot Slot
	want bool
}

func runSelectCases(t *testing.T, cases []selectCase) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.slot.SelectableAt(c.now); got != c.want {
				t.Fatalf("SelectableAt(%s) for %s–%s = %v, want %v",
					c.now.Format("2006-01-02 15:04"), c.slot.Start.Format("2006-01-02 15:04"), c.slot.End.Format("15:04"), got, c.want)
			}
		})
	}
}

// Rule: a date from today up to and including 30 days ahead.
func TestSlot_DateRange(t *testing.T) {
	midnight := at(monday, "00:00")
	runSelectCases(t, []selectCase{
		{"yesterday is not selectable", midnight, slot("2026-11-08", "10:00", "11:00"), false},
		{"today is selectable", midnight, slot(monday, "10:00", "11:00"), true},
		{"29 days ahead is selectable", midnight, slot("2026-12-08", "10:00", "11:00"), true},
		{"exactly 30 days ahead is selectable", midnight, slot("2026-12-09", "10:00", "11:00"), true},
		{"31 days ahead is not selectable", midnight, slot("2026-12-10", "10:00", "11:00"), false},
		{"30 days counts calendar days, not 24h from now", at(monday, "18:00"), slot("2026-12-09", "18:30", "19:00"), true},
	})
}

// Rule: within opening hours 08:00–19:00, every day; may end exactly at 19:00.
func TestSlot_OpeningHours(t *testing.T) {
	midnight := at(monday, "00:00")
	runSelectCases(t, []selectCase{
		{"starting before opening is not selectable", midnight, slot(monday, "07:45", "08:30"), false},
		{"starting at opening is selectable", midnight, slot(monday, "08:00", "08:30"), true},
		{"ending at closing is selectable", midnight, slot(monday, "18:30", "19:00"), true},
		{"ending after closing is not selectable", midnight, slot(monday, "18:45", "19:15"), false},
		{"starting at closing is not selectable", midnight, slot(monday, "19:00", "19:30"), false},
		{"the whole day is selectable", midnight, slot(monday, "08:00", "19:00"), true},
		{"a Saturday is selectable", midnight, slot("2026-11-14", "10:00", "11:00"), true},
		{"a Sunday is selectable", midnight, slot("2026-11-15", "10:00", "11:00"), true},
	})
}

// Rule: start and end fall on the 15-minute grid (:00, :15, :30, :45).
func TestSlot_Grid(t *testing.T) {
	midnight := at(monday, "00:00")
	runSelectCases(t, []selectCase{
		{"every quarter is selectable as a start", midnight, slot(monday, "10:15", "11:00"), true},
		{"start at :30 is selectable", midnight, slot(monday, "10:30", "11:00"), true},
		{"start at :45 is selectable", midnight, slot(monday, "10:45", "11:15"), true},
		{"start off the grid is not selectable", midnight, slot(monday, "10:07", "10:45"), false},
		{"end off the grid is not selectable", midnight, slot(monday, "10:15", "10:50"), false},
		{"start one minute past a quarter is not selectable", midnight, slot(monday, "10:01", "11:00"), false},
		{"end one minute before a quarter is not selectable", midnight, slot(monday, "10:00", "10:59"), false},
		{"seconds off the grid are not selectable", midnight,
			Slot{Start: at(monday, "10:00").Add(30 * time.Second), End: at(monday, "11:00")}, false},
	})
}

// Rule: at least 30 minutes, in 15-minute increments.
func TestSlot_Length(t *testing.T) {
	midnight := at(monday, "00:00")
	runSelectCases(t, []selectCase{
		{"15 minutes is not selectable", midnight, slot(monday, "10:00", "10:15"), false},
		{"29 minutes is not selectable", midnight, slot(monday, "10:00", "10:29"), false},
		{"30 minutes is selectable", midnight, slot(monday, "10:00", "10:30"), true},
		{"45 minutes is selectable", midnight, slot(monday, "10:00", "10:45"), true},
		{"zero length is not selectable", midnight, slot(monday, "10:00", "10:00"), false},
		{"end before start is not selectable", midnight, slot(monday, "11:00", "10:00"), false},
	})
}

// Rule: on today's date, only before the slot starts.
func TestSlot_NotStartedYet(t *testing.T) {
	runSelectCases(t, []selectCase{
		{"at 10:07, 10:15 is selectable", at(monday, "10:07"), slot(monday, "10:15", "11:00"), true},
		{"at 10:07, 10:00 has started and is not selectable", at(monday, "10:07"), slot(monday, "10:00", "11:00"), false},
		{"at exactly 10:15, 10:15 is not selectable", at(monday, "10:15"), slot(monday, "10:15", "11:00"), false},
		{"one second before 10:15, 10:15 is selectable", at(monday, "10:15").Add(-time.Second), slot(monday, "10:15", "11:00"), true},
		{"after closing, nothing today is selectable", at(monday, "19:30"), slot(monday, "18:30", "19:00"), false},
	})
}

func TestSlot_ValidSlotLaterTodayIsSelectable(t *testing.T) {
	now := at(monday, "00:00")
	if !slot(monday, "10:00", "11:00").SelectableAt(now) {
		t.Fatal("expected 10:00–11:00 today to be selectable at 00:00")
	}
}
