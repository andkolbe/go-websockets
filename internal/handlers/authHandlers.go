package handlers

import (
	"log"
	"net/http"

	"github.com/andkolbe/go-websockets/internal/forms"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/andkolbe/go-websockets/internal/render"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "password") // checks on our form
	if !form.Valid() {
		render.Template(w, "login.html", &models.TemplateData{})
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

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	user := models.User{
		Username:  r.Form.Get("username"),
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Password:  []byte(r.Form.Get("password")),
	}

	// create a new form
	form := forms.New(r.PostForm) // PostForm has all of the url values and their associated data
	form.Required("username", "first_name", "last_name", "email", "password")
	form.IsEmail("email")
	form.MinLength("password", 8) // add this to errors
	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = user
		render.Template(w, "register.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	id, err := m.DB.Register(user)
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// the user also logs in. Save the new user in the session
	m.App.Session.Put(r.Context(), "user_id", id) // add user_id into the session
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
