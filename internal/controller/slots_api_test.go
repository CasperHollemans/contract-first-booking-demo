package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/internal/service"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

// newTestMux serves the slot routes at Monday 9 November 2026, 00:00 UTC.
func newTestMux() *http.ServeMux {
	svc := service.NewSlotService(fixedClock{time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)})
	mux := http.NewServeMux()
	NewSlotController(svc, time.UTC).Register(mux)
	return mux
}

func get(mux http.Handler, url string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, url, nil))
	return rec
}

func TestSlotsCheck_ReturnsSelectable(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"/slots/check?date=2026-12-09&start=10:00&end=11:00", true},
		{"/slots/check?date=2026-11-08&start=10:00&end=11:00", false},
	}
	for _, c := range cases {
		rec := get(newTestMux(), c.url)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d, want 200", c.url, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Fatalf("%s: Content-Type %q, want application/json", c.url, ct)
		}
		var body struct {
			Selectable *bool `json:"selectable"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Selectable == nil {
			t.Fatalf("%s: body %q is not {\"selectable\": bool}", c.url, rec.Body.String())
		}
		if *body.Selectable != c.want {
			t.Fatalf("%s: selectable = %v, want %v", c.url, *body.Selectable, c.want)
		}
	}
}

func TestSlotsCheck_MalformedInputIsBadRequest(t *testing.T) {
	for _, url := range []string{
		"/slots/check",
		"/slots/check?date=9-12-2026&start=10:00&end=11:00",
		"/slots/check?date=2026-12-09&start=10am&end=11:00",
		"/slots/check?date=2026-12-09&start=10:00",
	} {
		if rec := get(newTestMux(), url); rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: status %d, want 400", url, rec.Code)
		}
	}
}
