// Command mockups runs a dev server for the mockup pages and the component library.
//
// Pages are found by filename, so a new mockup is just a new file:
//
//	mockups/site/<page>.html             a page of the site
//	mockups/site/<page>.<variant>.html   a variant of that page
//	mockups/prototypes/<name>.html       a loose experiment
//
// Every page gets a navigation bar at the bottom. Pages are read on each
// request, so edits and new files show up without a restart.
//
// Run it from the repository root: go run ./mockups
package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/casperhollemans/contract-first-booking-demo/web"
)

//go:embed chrome/*.html
var chromeFiles embed.FS

var chrome = template.Must(template.ParseFS(chromeFiles, "chrome/*.html"))

func main() {
	addr := ":3000"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("mockup server listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, newServer(os.DirFS("."))))
}

// newServer serves the mockups from root, the repository root.
func newServer(root fs.FS) http.Handler {
	s := &server{root: root}
	mux := http.NewServeMux()
	static, _ := fs.Sub(root, "web/static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /library", s.library)
	for _, section := range sections {
		mux.HandleFunc("GET /"+section+"/{page}", s.page(section))
	}
	return mux
}

type server struct {
	root fs.FS
}

// navData is what the navigation bar needs to render.
type navData struct {
	Catalog Catalog
	Page    *Page   // nil on the library
	Current Variant // the variant being shown
	Library bool
}

func (s *server) catalog() (Catalog, error) {
	mockups, err := fs.Sub(s.root, "mockups")
	if err != nil {
		return Catalog{}, err
	}
	return discover(mockups)
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	c, err := s.catalog()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := chrome.ExecuteTemplate(w, "index.html", c); err != nil {
		log.Print(err)
	}
}

func (s *server) library(w http.ResponseWriter, r *http.Request) {
	c, err := s.catalog()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, "mockups/library.html", navData{Catalog: c, Library: true})
}

// page serves /{section}/{page}. Only pages found by discover are served.
func (s *server) page(section string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := s.catalog()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, v, ok := c.find(section, r.PathValue("page"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		s.render(w, "mockups/"+v.File, navData{Catalog: c, Page: &p, Current: v})
	}
}

// render executes a mockup with the shared components and adds the nav bar.
func (s *server) render(w http.ResponseWriter, file string, nav navData) {
	t := template.New("").Funcs(web.FuncMap())
	t, err := t.ParseFS(s.root, "web/templates/components/*.html", file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var page bytes.Buffer
	if err := t.ExecuteTemplate(&page, path.Base(file), nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var bar bytes.Buffer
	if err := chrome.ExecuteTemplate(&bar, "nav.html", nav); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(injectNav(page.Bytes(), bar.Bytes()))
}

// injectNav puts the bar just before the last </body>, or at the end.
func injectNav(page, bar []byte) []byte {
	i := bytes.LastIndex(bytes.ToLower(page), []byte("</body>"))
	if i < 0 {
		return append(page, bar...)
	}
	out := make([]byte, 0, len(page)+len(bar))
	out = append(out, page[:i]...)
	out = append(out, bar...)
	return append(out, page[i:]...)
}
