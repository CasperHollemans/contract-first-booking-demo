// BOOK-1 · Choose a date and time slot
// docs/stories/BOOK-1-choose-date-and-time-slot.md
//
// All examples assume it is Monday 9 November, 00:00, unless stated otherwise.
package specs

import (
	"testing"

	"github.com/casperhollemans/contract-first-booking-demo/acceptance/dsl"
)

func TestBOOK1_E1_ExactlyThirtyDaysAhead(t *testing.T) {
	dsl.ForEachDriver(t, func(user *dsl.Booking) {
		user.GivenToday("9 November")
		user.WhenUserChooses("9 December", "10:00–11:00")
		user.ThenSlotCanBeSelected()
	})
}

func TestBOOK1_E2_PastDate(t *testing.T) {
	dsl.ForEachDriver(t, func(user *dsl.Booking) {
		user.GivenToday("9 November")
		user.WhenUserChooses("8 November", "10:00–11:00")
		user.ThenSlotCannotBeSelected()
	})
}

func TestBOOK1_E3_StartsAtOpening(t *testing.T) {
	dsl.ForEachDriver(t, func(user *dsl.Booking) {
		user.WhenUserChoosesToday("08:00–08:30")
		user.ThenSlotCanBeSelected()
	})
}

func TestBOOK1_E4_EndsAtClosing(t *testing.T) {
	dsl.ForEachDriver(t, func(user *dsl.Booking) {
		user.WhenUserChoosesToday("18:30–19:00")
		user.ThenSlotCanBeSelected()
	})
}

func TestBOOK1_E5_TwentyNineMinutes(t *testing.T) {
	dsl.ForEachDriver(t, func(user *dsl.Booking) {
		user.WhenUserChoosesToday("10:00–10:29")
		user.ThenSlotCannotBeSelected()
	})
}

func TestBOOK1_E6_OffTheGrid(t *testing.T) {
	dsl.ForEachDriver(t, func(user *dsl.Booking) {
		user.WhenUserChoosesToday("10:07–10:45")
		user.ThenSlotCannotBeSelected()
	})
}
