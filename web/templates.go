// Package web provides the shared template functions and template loader.
package web

import "html/template"

// FuncMap returns the shared template function map.
func FuncMap() template.FuncMap {
	return template.FuncMap{
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
}

// LoadTemplates parses the shared component and page templates from web/.
func LoadTemplates() *template.Template {
	t := template.New("").Funcs(FuncMap())
	t = template.Must(t.ParseGlob("web/templates/components/*.html"))
	t = template.Must(t.ParseGlob("web/templates/pages/*.html"))
	return t
}
