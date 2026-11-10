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

	selectable bool
	checkErr   error
}

func newBooking(t *testing.T, via Driver, d drivers.Driver) *Booking {
	b := &Booking{t: t, via: via, driver: d}
	b.GivenToday("9 November") // the examples' default: Monday 9 November, 00:00
	return b
}

// GivenToday sets "now" to 00:00 on the given day, e.g. "9 November".
func (b *Booking) GivenToday(day string) {
	b.t.Helper()
	b.driver.SetNow(parseDay(b.t, day))
}

// GivenNow sets "now" to a time on the given day: ("9 November", "10:07").
func (b *Booking) GivenNow(day, clock string) {
	b.t.Helper()
	at, err := time.Parse("15:04", clock)
	if err != nil {
		b.t.Fatalf("dsl: cannot read time %q, expected e.g. \"10:07\"", clock)
	}
	sinceMidnight := time.Duration(at.Hour())*time.Hour + time.Duration(at.Minute())*time.Minute
	b.driver.SetNow(parseDay(b.t, day).Add(sinceMidnight))
}

// WhenUserChooses picks a slot on a given day: ("9 December", "10:00–11:00").
// The slot may use an en dash or a hyphen.
func (b *Booking) WhenUserChooses(day, slot string) {
	b.t.Helper()
	date := parseDay(b.t, day).Format("2006-01-02")
	start, end := parseSlot(b.t, slot)
	b.selectable, b.checkErr = b.driver.CheckSlot(date, start, end)
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
