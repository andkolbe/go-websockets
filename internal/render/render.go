package render

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
)

// must have this to use the jet templating engine
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(), // we don't have to restart our app every time we make a change to a jet template
)

// holds data send from handlers to templates
type TemplateData struct {
	CSRFToken string
}

// DefaultData adds default data which is accessible to all templates
func DefaultData(td TemplateData, r *http.Request, w http.ResponseWriter) TemplateData {
	td.CSRFToken = nosurf.Token(r)
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
