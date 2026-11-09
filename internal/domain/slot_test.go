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

func TestSlot_ValidSlotLaterTodayIsSelectable(t *testing.T) {
	now := at(monday, "00:00")
	if !slot(monday, "10:00", "11:00").SelectableAt(now) {
		t.Fatal("expected 10:00–11:00 today to be selectable at 00:00")
	}
}
