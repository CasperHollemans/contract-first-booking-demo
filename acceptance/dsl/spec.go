package dsl

import (
	"flag"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/casperhollemans/contract-first-booking-demo/acceptance/drivers"
)

// Driver names a way of reaching the system. Each spec lists the drivers
// that prove it; choose the cheapest one that can.
type Driver string

const (
	// Domain: the service layer with the fake clock, no HTTP. Rules and logic.
	Domain Driver = "Domain"
	// API: the real server over HTTP. Errors, wiring, anything concurrent.
	API Driver = "API"
	// UI: the real server through headless Chrome. The happy path only.
	UI Driver = "UI"
)

// newDriver starts a fresh system (and fake clock) for one driver run.
func newDriver(t *testing.T, d Driver) drivers.Driver {
	t.Helper()
	switch d {
	case Domain:
		return drivers.NewDomainDriver()
	case API:
		return drivers.NewAPIDriver(t)
	case UI:
		return drivers.NewUIDriver(t)
	default:
		t.Fatalf("dsl: unknown driver %q", d)
		return nil
	}
}

// Spec is one behaviour of the contract. It runs the steps once per listed
// driver, each with a fresh Booking, system and clock:
//
//	Spec(t, "When …, expect …", func(b *Booking) { … }, Domain, API, UI)
//
// A spec with no drivers fails: it would prove nothing.
func Spec(t *testing.T, behaviour string, steps func(b *Booking), ds ...Driver) {
	t.Helper()
	t.Run(behaviour, func(t *testing.T) {
		t.Helper()
		if len(ds) == 0 {
			t.Fatal("this spec never ran: add Domain, API or UI")
		}
		for _, d := range ds {
			t.Run(string(d), func(t *testing.T) {
				ran.record(storyOf(t.Name()), d)
				steps(newBooking(t, d, newDriver(t, d)))
			})
		}
	})
}

// storyOf reads the story from the top-level test name:
// "TestBOOK1_DateRange/…" → "BOOK1".
func storyOf(testName string) string {
	top, _, _ := strings.Cut(testName, "/")
	story, _, _ := strings.Cut(strings.TrimPrefix(top, "Test"), "_")
	return story
}

// ran records, per story, which drivers ran a spec.
var ran = &runLog{byStory: map[string]map[Driver]bool{}}

type runLog struct {
	mu      sync.Mutex
	byStory map[string]map[Driver]bool
}

func (l *runLog) record(story string, d Driver) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.byStory[story] == nil {
		l.byStory[story] = map[Driver]bool{}
	}
	l.byStory[story][d] = true
}

// StoriesWithoutUI lists the stories that ran specs but none through the UI.
// Every story must prove at least its happy path in the browser.
//
// It returns nil when the run was filtered with -run or -skip: a partial run
// can't tell a missing UI spec from one that was filtered out.
func StoriesWithoutUI() []string {
	if filtered() {
		return nil
	}
	ran.mu.Lock()
	defer ran.mu.Unlock()
	var missing []string
	for story, ds := range ran.byStory {
		if !ds[UI] {
			missing = append(missing, story)
		}
	}
	sort.Strings(missing)
	return missing
}

func filtered() bool {
	for _, name := range []string{"test.run", "test.skip"} {
		if f := flag.Lookup(name); f != nil && f.Value.String() != "" {
			return true
		}
	}
	return false
}
