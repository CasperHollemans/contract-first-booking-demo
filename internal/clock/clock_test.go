package clock

import (
	"testing"
	"time"
)

func TestSystem_ReturnsTheRealTime(t *testing.T) {
	before := time.Now()
	got := System{}.Now()
	if got.Before(before) || got.After(time.Now()) {
		t.Fatalf("System.Now() = %v, want the current time", got)
	}
}
