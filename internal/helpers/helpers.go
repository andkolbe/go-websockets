package helpers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// must have this to use the jet templating engine
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(), // we don't have to restart our app every time we make a change to a jet template
)

// allows us to use the views variable in our handlers tests
func SetViews(path string) {
	views = jet.NewSet(
		jet.NewOSFileSystemLoader(path),
	)
}

// holds data send from handlers to templates
type TemplateData struct {
	CSRFToken       string
	IsAuthenticated bool
	Flash           string
	Warning         string
	Error           string
}

// DefaultData adds default data which is accessible to all templates
func DefaultData(td TemplateData, r *http.Request) TemplateData {
	td.CSRFToken = nosurf.Token(r)

	return td
}

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data jet.VarMap) error {
	// add default template data
	var td TemplateData

	// add default data
	td = DefaultData(td, r)

	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	err = view.Execute(w, data, td)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// return true or false if the user is authenticated or not
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
