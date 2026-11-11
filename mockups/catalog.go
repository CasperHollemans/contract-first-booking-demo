package main

import (
	"html"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

// The mockups are organised like the site itself:
//
//	site/book.html            the Book page (current)
//	site/book.option-a.html   a variant of the Book page
//	prototypes/*.html         loose experiments, kept out of the site menu
//	library.html              the component library
type Catalog struct {
	Site       []Page
	Prototypes []Page
}

// Page is one page of the site (or one prototype) with its variants.
type Page struct {
	Section  string // "site" or "prototypes"
	Name     string // "book"
	Label    string // "Book"
	Variants []Variant
}

// Variant is one version of a page. The variant without a suffix is "Current".
type Variant struct {
	Name  string // "" for current, "option-a" for book.option-a.html
	Label string // "Current", "Option A"
	Title string // the page's <title>
	File  string // "site/book.option-a.html"
	URL   string // "/site/book.option-a"
}

var sections = []string{"site", "prototypes"}

// discover reads the mockup pages from fsys, which is rooted at mockups/.
func discover(fsys fs.FS) (Catalog, error) {
	var c Catalog
	for _, section := range sections {
		pages, err := discoverSection(fsys, section)
		if err != nil {
			return Catalog{}, err
		}
		if section == "site" {
			c.Site = pages
		} else {
			c.Prototypes = pages
		}
	}
	return c, nil
}

func discoverSection(fsys fs.FS, section string) ([]Page, error) {
	files, err := fs.Glob(fsys, section+"/*.html")
	if err != nil {
		return nil, err
	}
	byName := map[string]*Page{}
	for _, file := range files {
		base := strings.TrimSuffix(path.Base(file), ".html")
		name, variant, _ := strings.Cut(base, ".")
		p, ok := byName[name]
		if !ok {
			p = &Page{Section: section, Name: name, Label: pageLabel(name)}
			byName[name] = p
		}
		v := Variant{Name: variant, Label: variantLabel(variant), File: file, URL: "/" + section + "/" + base}
		if b, err := fs.ReadFile(fsys, file); err == nil {
			v.Title = titleOf(b)
		}
		p.Variants = append(p.Variants, v)
	}

	pages := make([]Page, 0, len(byName))
	for _, p := range byName {
		sort.Slice(p.Variants, func(i, j int) bool {
			a, b := p.Variants[i], p.Variants[j]
			if (a.Name == "") != (b.Name == "") {
				return a.Name == "" // Current first
			}
			return a.Name < b.Name
		})
		pages = append(pages, *p)
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].Name < pages[j].Name })
	return pages, nil
}

// find returns the page and variant served at /{section}/{base}.
func (c Catalog) find(section, base string) (Page, Variant, bool) {
	pages := c.Site
	if section == "prototypes" {
		pages = c.Prototypes
	}
	for _, p := range pages {
		for _, v := range p.Variants {
			if v.URL == "/"+section+"/"+base {
				return p, v, true
			}
		}
	}
	return Page{}, Variant{}, false
}

// pageLabel: "thank-you" → "Thank you".
func pageLabel(name string) string {
	s := strings.ReplaceAll(name, "-", " ")
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// variantLabel: "" → "Current", "option-a" → "Option A".
func variantLabel(name string) string {
	if name == "" {
		return "Current"
	}
	words := strings.Split(name, "-")
	for i, w := range words {
		if w != "" {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

var titleRE = regexp.MustCompile(`(?is)<title>(.*?)</title>`)

func titleOf(b []byte) string {
	m := titleRE.FindSubmatch(b)
	if m == nil {
		return ""
	}
	return strings.TrimSpace(html.UnescapeString(string(m[1])))
}
