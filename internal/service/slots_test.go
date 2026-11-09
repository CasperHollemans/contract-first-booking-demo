package service

import (
	"testing"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/internal/domain"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func TestSlotService_UsesTheClockForNow(t *testing.T) {
	tenFifteen := time.Date(2026, 11, 9, 10, 15, 0, 0, time.UTC)
	s := domain.Slot{Start: tenFifteen, End: tenFifteen.Add(45 * time.Minute)}

	before := NewSlotService(fixedClock{tenFifteen.Add(-8 * time.Minute)})
	if !before.IsSelectable(s) {
		t.Fatal("at 10:07, 10:15–11:00 should be selectable")
	}
	at := NewSlotService(fixedClock{tenFifteen})
	if at.IsSelectable(s) {
		t.Fatal("at 10:15, 10:15–11:00 should not be selectable")
	}
}
