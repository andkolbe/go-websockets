package handlers

import (
	"log"
	"net/http"

	"github.com/andkolbe/go-websockets/internal/forms"
	"github.com/andkolbe/go-websockets/internal/helpers"
)


func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("username", "password")
	if !form.Valid() {
		render.RenderPage(w, r, "login.jet.html", nil)
		return
	}

	id, _, err := m.DB.Login(username, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
