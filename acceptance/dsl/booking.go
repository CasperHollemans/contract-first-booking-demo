// Package dsl turns example language ("9 December", "10:00–11:00") into
// driver calls. Specs only talk to this layer.
package dsl

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/acceptance/drivers"
)

// Examples don't name a year; the demo's examples are set in 2026
// (Monday 9 November 2026).
const exampleYear = 2026

// Example times are in the cafe's local time.
var cafeTime = time.UTC

// Booking is one user's session with the system, through one driver.
type Booking struct {
	t      *testing.T
	via    Driver // named in every failure, next to the subtest name
	driver drivers.Driver
	today  time.Time

	selectable bool
	checkErr   error
}

// ForEachDriver runs the example once through the API and once through the UI.
func ForEachDriver(t *testing.T, example func(b *Booking)) {
	t.Run("api", func(t *testing.T) {
		example(newBooking(t, API, drivers.NewAPIDriver(t)))
	})
	t.Run("ui", func(t *testing.T) {
		example(newBooking(t, UI, drivers.NewUIDriver(t)))
	})
}

func newBooking(t *testing.T, via Driver, d drivers.Driver) *Booking {
	b := &Booking{t: t, via: via, driver: d}
	b.GivenToday("9 November") // the examples' default: Monday 9 November, 00:00
	return b
}

// GivenToday sets "now" to 00:00 on the given day, e.g. "9 November".
func (b *Booking) GivenToday(day string) {
	b.t.Helper()
	b.today = parseDay(b.t, day)
	b.driver.SetNow(b.today)
}

// GivenNow sets "now" to a time on the given day: ("9 November", "10:07").
func (b *Booking) GivenNow(day, clock string) {
	b.t.Helper()
	b.today = parseDay(b.t, day)
	at, err := time.Parse("15:04", clock)
	if err != nil {
		b.t.Fatalf("dsl: cannot read time %q, expected e.g. \"10:07\"", clock)
	}
	b.driver.SetNow(b.today.Add(time.Duration(at.Hour())*time.Hour + time.Duration(at.Minute())*time.Minute))
}

// WhenUserChooses picks a slot on a given day: ("9 December", "10:00–11:00").
func (b *Booking) WhenUserChooses(day, slot string) {
	b.t.Helper()
	b.choose(parseDay(b.t, day), slot)
}

// WhenUserChoosesToday picks a slot on today's date: ("08:00–08:30").
func (b *Booking) WhenUserChoosesToday(slot string) {
	b.t.Helper()
	b.choose(b.today, slot)
}

func (b *Booking) choose(day time.Time, slot string) {
	b.t.Helper()
	start, end := parseSlot(b.t, slot)
	b.selectable, b.checkErr = b.driver.CheckSlot(day.Format("2006-01-02"), start, end)
}

func (b *Booking) ThenSlotCanBeSelected() {
	b.t.Helper()
	b.requireAnswer()
	if !b.selectable {
		b.t.Fatalf("%s: expected the slot to be selectable, but it cannot be selected", b.via)
	}
}

func (b *Booking) ThenSlotCannotBeSelected() {
	b.t.Helper()
	b.requireAnswer()
	if b.selectable {
		b.t.Fatalf("%s: expected the slot not to be selectable, but it can be selected", b.via)
	}
}

func (b *Booking) requireAnswer() {
	b.t.Helper()
	if b.checkErr != nil {
		b.t.Fatalf("%s: the system did not answer: %v", b.via, b.checkErr)
	}
}

// parseDay reads "9 December" as 00:00 on that day in exampleYear.
func parseDay(t *testing.T, day string) time.Time {
	t.Helper()
	d, err := time.ParseInLocation("2 January 2006", fmt.Sprintf("%s %d", day, exampleYear), cafeTime)
	if err != nil {
		t.Fatalf("dsl: cannot read day %q: %v", day, err)
	}
	return d
}

// parseSlot reads "10:00–11:00" (en dash or hyphen) into "10:00", "11:00".
// Times are passed through as written, so off-grid times like 10:07 reach the system.
func parseSlot(t *testing.T, slot string) (start, end string) {
	t.Helper()
	parts := strings.Split(strings.ReplaceAll(slot, "–", "-"), "-")
	if len(parts) != 2 {
		t.Fatalf("dsl: cannot read slot %q, expected e.g. \"10:00–11:00\"", slot)
	}
	for _, p := range parts {
		if _, err := time.Parse("15:04", strings.TrimSpace(p)); err != nil {
			t.Fatalf("dsl: cannot read time %q in slot %q", p, slot)
		}
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}
