package render

import (
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// a map of functions that can be used in a template
var functions = template.FuncMap{}

var app *config.AppConfig

// sets the app var when this function is called
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// needs the response writer so it has somewhere to send the webpage
func Template(w http.ResponseWriter, html string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config. The template cache is initialied in our main function
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// pull the individual template out of the map
	// check if the template exists by adding the ok var
	t, ok := tc[html]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	// puts the parsed template that is being held in memory, into bytes
	buf := new(bytes.Buffer)
	// take the template, execute it, don't pass it any data, and store the value in the buf

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)
	// write to response writer
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser")
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// you can look up data in a map very quickly. Similar to arrays
	// the string name is the template path "about.page.html"
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// does the template match any layouts in our templates folder?
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		// add the template to the cache
		myCache[name] = ts
	}

	return myCache, nil

}
