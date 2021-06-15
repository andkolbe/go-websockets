package handlers

import (
	"log"
	"net/http"
	
	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/andkolbe/go-websockets/internal/models"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	id, hash, err := m.DB.Authenticate(r.Form.Get("username"), r.Form.Get("password"))
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		err := helpers.RenderPage(w, r, "login.jet.html", nil)
		if err != nil {
			printTemplateError(w, err)
		}
		return
	}

	// We authenticated. Get the user
	u, err := m.DB.GetUserByID(id)
	if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	m.App.Session.Put(r.Context(), "userID", id)
	m.App.Session.Put(r.Context(), "hashedPassword", hash)
	m.App.Session.Put(r.Context(), "flash", "You've been logged in successfully!")
	m.App.Session.Put(r.Context(), "user", u)

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
		Username:  r.Form.Get("username"), // r.Form.Get("username") matches the name="username" field on the html page
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Password:  []byte(r.Form.Get("password")),
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
