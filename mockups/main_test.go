package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func page(title string) *fstest.MapFile {
	return &fstest.MapFile{Data: []byte(`<!doctype html><html><head><title>` + title + `</title></head><body>` +
		`{{template "button" "x"}}</body></html>`)}
}

func testRoot() fstest.MapFS {
	return fstest.MapFS{
		"web/templates/components/button.html": {Data: []byte(`{{define "button"}}<button>{{.}}</button>{{end}}`)},
		"web/static/tokens.css":                {Data: []byte(`:root{}`)},
		"mockups/library.html":                 page("Library"),
		"mockups/site/book.html":               page("Book a room"),
		"mockups/site/book.option-b.html":      page("Book: Option B"),
		"mockups/site/book.option-a.html":      page("Book: Option A"),
		"mockups/site/thank-you.html":          page("Thanks"),
		"mockups/prototypes/scratch.html":      page("Scratch"),
	}
}

func TestDiscover_GroupsVariantsUnderTheirPage(t *testing.T) {
	c, err := discover(mustSub(t, testRoot(), "mockups"))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Site) != 2 || c.Site[0].Label != "Book" || c.Site[1].Label != "Thank you" {
		t.Fatalf("site pages = %+v, want Book and Thank you", c.Site)
	}
	var got []string
	for _, v := range c.Site[0].Variants {
		got = append(got, v.Label+" "+v.URL+" "+v.Title)
	}
	want := []string{
		"Current /site/book Book a room",
		"Option A /site/book.option-a Book: Option A",
		"Option B /site/book.option-b Book: Option B",
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("variants =\n%v\nwant\n%v", got, want)
	}
	if len(c.Prototypes) != 1 || c.Prototypes[0].Variants[0].URL != "/prototypes/scratch" {
		t.Fatalf("prototypes = %+v, want scratch", c.Prototypes)
	}
}

func TestPage_RendersWithComponentsAndNavBar(t *testing.T) {
	body := get(t, "/site/book.option-a", http.StatusOK)
	for _, want := range []string{
		`<button>x</button>`,         // shared components work
		`data-mockup-nav`,            // the bar is there
		`href="/library"`,            // library is always reachable
		`href="/prototypes/scratch"`, // prototypes menu
		`href="/site/thank-you"`,     // site menu
		`href="/site/book.option-a" data-mk-variant aria-current="page"`, // current variant marked
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page is missing %q", want)
		}
	}
	if i, j := strings.Index(body, "data-mockup-nav"), strings.LastIndex(body, "</body>"); i > j {
		t.Error("nav bar is not inside <body>")
	}
}

func TestPage_WithoutVariantsHasNoSwitcher(t *testing.T) {
	body := get(t, "/site/thank-you", http.StatusOK)
	if strings.Contains(body, " data-mk-variant ") { // the attribute on a link, not the script's selector
		t.Error("a page without variants shows a variant switcher")
	}
}

func TestLibrary_HasNavBar(t *testing.T) {
	body := get(t, "/library", http.StatusOK)
	if !strings.Contains(body, `href="/library" aria-current="page"`) {
		t.Error("library link is not marked as current on the library")
	}
}

func TestIndex_ListsSitePrototypesAndLibrary(t *testing.T) {
	body := get(t, "/", http.StatusOK)
	for _, want := range []string{`href="/site/book"`, `href="/site/book.option-b"`, `href="/prototypes/scratch"`, `href="/library"`} {
		if !strings.Contains(body, want) {
			t.Errorf("index is missing %q", want)
		}
	}
}

func TestUnknownPagesAreNotFound(t *testing.T) {
	for _, url := range []string{"/site/nope", "/other/book", "/site/..%2Flibrary"} {
		get(t, url, http.StatusNotFound)
	}
}

func get(t *testing.T, url string, wantStatus int) string {
	t.Helper()
	rec := httptest.NewRecorder()
	newServer(testRoot()).ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
	if rec.Code != wantStatus {
		t.Fatalf("GET %s = %d, want %d: %s", url, rec.Code, wantStatus, rec.Body)
	}
	b, _ := io.ReadAll(rec.Body)
	return string(b)
}

func mustSub(t *testing.T, fsys fstest.MapFS, dir string) fstest.MapFS {
	t.Helper()
	sub := fstest.MapFS{}
	for name, f := range fsys {
		if strings.HasPrefix(name, dir+"/") {
			sub[strings.TrimPrefix(name, dir+"/")] = f
		}
	}
	return sub
}
