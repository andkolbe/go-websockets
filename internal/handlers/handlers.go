package handlers

import (
	"database/sql"
	"net/http"

	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/andkolbe/go-websockets/internal/repository"
	"github.com/andkolbe/go-websockets/internal/repository/dbrepo"
)

// repository pattern allows us to swap components out of our app with minimal changes required to the code base

// the repository used by the handlers
var Repo *Repository

// var app *config.AppConfig

// the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}
// creates a new repository
func NewRepo(a *config.AppConfig, db *sql.DB) *Repository {
	return &Repository {
		App: a,
		DB: dbrepo.NewPostgresRepo(db, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository {
		App: a,
		DB: dbrepo.NewTestingRepo(a),
	}
}

// when we call newRepo, we pass it the app config (a pointer to config.AppConfig), 
// and the database connection pool (a pointer to driver.DB, which holds to db connection pool)
// we then populate the Repository type with all of the information we receive as parameters
// and hand that back as a pointer to Repository

// sets the repository for the handlers on the main package
func NewHandlers(repo *Repository) {
	Repo = repo
}

// every web handler in Go must have a response writer and a pointer to a request

// giving the handlers a receiver links them together with the repository, so all of the handlers have access to the repository
// those handlers have access to everything inside of the app config and the database driver
// with receivers, Go doesn't have to mess around with classes or inheritance

func (m *Repository) LoginPage(w http.ResponseWriter, r *http.Request) {
	// if already logged in, take to chat room
	if m.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/auth/chat", http.StatusSeeOther)
		return
	}

	err := helpers.RenderPage(w, r, "login", nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

func (m *Repository) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// if already logged in, take to chat room
	if m.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/auth/chat", http.StatusSeeOther)
		return
	}
	
	err := helpers.RenderPage(w, r, "register", nil)
	if err != nil {
		printTemplateError(w, err)
	}
	
}

func (m *Repository) ChatRoomPage(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "chat", nil)
	if err != nil {
		printTemplateError(w, err)
	}
}