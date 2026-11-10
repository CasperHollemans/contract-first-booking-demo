package specs

import (
	"testing"

	. "github.com/casperhollemans/contract-first-booking-demo/acceptance/dsl"
)

// BOOK-1 · Choose a date and time slot
// As a user, I want to choose a date and time slot for a room,
// so that I can make use of the room for a prespecified time period.

func TestBOOK1_DateRange(t *testing.T) {
	Spec(t, "When the user chooses exactly 30 days ahead, expect the slot can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("9 December", "10:00–11:00")
		b.ThenSlotCanBeSelected()
	}, Domain, API, UI)

	Spec(t, "When the user chooses 31 days ahead, expect the slot cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 December", "10:00–11:00")
		b.ThenSlotCannotBeSelected()
	}, Domain)

	Spec(t, "When the user chooses a past date, expect the slot cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("8 November", "10:00–11:00")
		b.ThenSlotCannotBeSelected()
	}, Domain, API)

	Spec(t, "When the user chooses a weekend day, expect the slot can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("14 November", "10:00–11:00")
		b.ThenSlotCanBeSelected()
	}, Domain)
}

func TestBOOK1_OpeningHours(t *testing.T) {
	Spec(t, "When the slot starts at opening, expect it can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "08:00–08:30")
		b.ThenSlotCanBeSelected()
	}, Domain)

	Spec(t, "When the slot ends exactly at closing, expect it can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "18:30–19:00")
		b.ThenSlotCanBeSelected()
	}, Domain)

	Spec(t, "When the slot covers the whole day, expect it can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "08:00–19:00")
		b.ThenSlotCanBeSelected()
	}, Domain)

	Spec(t, "When the slot starts before opening, expect it cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "07:45–08:30")
		b.ThenSlotCannotBeSelected()
	}, Domain)

	Spec(t, "When the slot ends after closing, expect it cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "18:45–19:15")
		b.ThenSlotCannotBeSelected()
	}, Domain)
}

func TestBOOK1_LengthAndGrid(t *testing.T) {
	Spec(t, "When the slot is 30 minutes, expect it can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "10:00–10:30")
		b.ThenSlotCanBeSelected()
	}, Domain)

	Spec(t, "When the slot is 45 minutes, expect it can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "10:00–10:45")
		b.ThenSlotCanBeSelected()
	}, Domain)

	Spec(t, "When the slot is 29 minutes, expect it cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "10:00–10:29")
		b.ThenSlotCannotBeSelected()
	}, Domain)

	Spec(t, "When the slot is 15 minutes, expect it cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "10:00–10:15")
		b.ThenSlotCannotBeSelected()
	}, Domain)

	Spec(t, "When the start is off the 15-minute grid, expect it cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "10:07–10:45")
		b.ThenSlotCannotBeSelected()
	}, Domain)

	Spec(t, "When the end is off the 15-minute grid, expect it cannot be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("10 November", "10:15–10:50")
		b.ThenSlotCannotBeSelected()
	}, Domain)
}

func TestBOOK1_Today(t *testing.T) {
	Spec(t, "When the slot later today hasn't started, expect it can be selected", func(b *Booking) {
		b.GivenNow("9 November", "10:07")
		b.WhenUserChooses("9 November", "10:15–11:00")
		b.ThenSlotCanBeSelected()
	}, Domain)

	Spec(t, "When the slot today has already started, expect it cannot be selected", func(b *Booking) {
		b.GivenNow("9 November", "10:07")
		b.WhenUserChooses("9 November", "10:00–11:00")
		b.ThenSlotCannotBeSelected()
	}, Domain)

	Spec(t, "When the slot today starts right now, expect it cannot be selected", func(b *Booking) {
		b.GivenNow("9 November", "10:15")
		b.WhenUserChooses("9 November", "10:15–11:00")
		b.ThenSlotCannotBeSelected()
	}, Domain)
}
