package controller

import (
	"net/http"
	"strings"
	"testing"
)

func TestPage_HasTheSlotForm(t *testing.T) {
	rec := get(newTestMux(), "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type %q, want text/html", ct)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`data-testid="date"`, `data-testid="start"`, `data-testid="end"`, `data-testid="check"`,
		`hx-get="/ui/slots/check"`, `src="/static/htmx.min.js"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page is missing %s", want)
		}
	}
}

func TestPage_ServesHtmx(t *testing.T) {
	rec := get(newTestMux(), "/static/htmx.min.js")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "htmx") {
		t.Fatalf("status %d; want htmx served", rec.Code)
	}
}

func TestFragment_ShowsSelectable(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{"/ui/slots/check?date=2026-12-09&start=10:00&end=11:00", `data-selectable="true"`},
		{"/ui/slots/check?date=2026-11-08&start=10:00&end=11:00", `data-selectable="false"`},
	}
	for _, c := range cases {
		rec := get(newTestMux(), c.url)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d, want 200", c.url, rec.Code)
		}
		body := rec.Body.String()
		if !strings.Contains(body, `data-testid="slot-result"`) || !strings.Contains(body, c.want) {
			t.Fatalf("%s: fragment %q, want slot-result with %s", c.url, body, c.want)
		}
	}
}

func TestFragment_MalformedInputIsBadRequest(t *testing.T) {
	if rec := get(newTestMux(), "/ui/slots/check?date=&start=&end="); rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want 400", rec.Code)
	}
}
