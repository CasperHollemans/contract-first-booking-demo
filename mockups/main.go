// Command mockups runs a dev server that renders the mockup pages and component library.
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				panic("dict requires an even number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					panic("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict
		},
	}

	templates = template.New("").Funcs(funcMap)
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
	log.Printf("  component library: http://localhost%s/library", addr)
	log.Printf("  book-1 mockup:     http://localhost%s/", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
