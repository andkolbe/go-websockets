package helpers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/CloudyKit/jet/v6"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/justinas/nosurf"
)

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
	User            models.User
	Flash           string
	Warning         string
	Error           string
}

// DefaultData adds default data which is accessible to all templates
func DefaultData(td TemplateData, r *http.Request) TemplateData {
	td.CSRFToken = nosurf.Token(r)
	// NEED TO FIGURE OUT WHAT THE PROBLEM IS
	// td.IsAuthenticated = IsAuthenticated(r)
	// // if logged in, store user id in template data
	// if td.IsAuthenticated {
	// 	u := app.Session.Get(r.Context(), "user").(models.User)
	// 	td.User = u
	// }

	// td.Flash = app.Session.PopString(r.Context(), "flash")
	// td.Warning = app.Session.PopString(r.Context(), "warning")
	// td.Error = app.Session.PopString(r.Context(), "error")

	return td
}

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data jet.VarMap) error {
	// add default template data
	var td TemplateData

	// add default data
	td = DefaultData(td, r)

	view, err := views.GetTemplate(fmt.Sprintf("%s.jet.html", tmpl))
	if err != nil {
		log.Println(err)
		return err
	}
	if err = view.Execute(w, data, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}
