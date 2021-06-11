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

	form := forms.New(r.PostForm)
	form.Required("username", "password")
	if !form.Valid() {
		helpers.RenderPage(w, r, "login.jet.html", nil)
		return
	}

	id, _, err := m.DB.Login(r.Form.Get("username"), r.Form.Get("password"))
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id) // add user_id into the session
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// destroy the session
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


// form.MinLength("password", 8) add this to register