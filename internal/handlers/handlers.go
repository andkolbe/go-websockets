package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(), // we don't have to restart our app every time we make a change to a jet template
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "login.jet.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "register.jet.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func Chat(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "chat.jet.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
