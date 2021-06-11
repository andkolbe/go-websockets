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
	IsAuthenticated int
}

// DefaultData adds default data which is accessible to all templates
func DefaultData(td TemplateData, r *http.Request, w http.ResponseWriter) TemplateData {
	td.CSRFToken = nosurf.Token(r)

	// DefaultData has access to the session because it has the request
	// if app.Session.Exists(r.Context(), "user_id") { // when someone logs in, we put "user_id" in the session
	// 	td.IsAuthenticated = 1 // 1 means the user is logged in. 0 means the user is logged out
	// }
	return td
}

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data jet.VarMap) error {
	// add default template data
	var td TemplateData

	// add default data
	td = DefaultData(td, r, w)

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
