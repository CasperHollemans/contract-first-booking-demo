// Command mockups runs a dev server that renders the mockup pages and component library.
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/casperhollemans/contract-first-booking-demo/web"
)

var templates *template.Template

func init() {
	templates = template.New("").Funcs(web.FuncMap())
	templates = template.Must(templates.ParseGlob("web/templates/components/*.html"))
	templates = template.Must(templates.ParseGlob("mockups/*.html"))
	templates = template.Must(templates.ParseGlob("mockups/pages/*.html"))
}

func main() {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Component library
	mux.HandleFunc("GET /library", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "library.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Mockup pages
	mux.HandleFunc("GET /option-a", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "book-1-option-a.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /option-b", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "book-1-option-b.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "book-1.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := ":3000"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("mockup server listening on %s", addr)
	log.Printf("  component library:       http://localhost%s/library", addr)
	log.Printf("  option A (week grid):    http://localhost%s/option-a", addr)
	log.Printf("  option B (two-step):     http://localhost%s/option-b", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
