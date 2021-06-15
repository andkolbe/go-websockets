package handlers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/CloudyKit/jet/v6"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/driver"
	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/andkolbe/go-websockets/internal/models"
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
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository {
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
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

func (m *Repository) LoginPage(w http.ResponseWriter, r *http.Request) {
	// if already logged in, take to chat room
	if m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	err := helpers.RenderPage(w, r, "login.jet.html", nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

func (m *Repository) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var emptyUser models.User
	data := make(jet.VarMap)
	data.Set("user", emptyUser)
	
	err := helpers.RenderPage(w, r, "register.jet.html", data)
	if err != nil {
		printTemplateError(w, err)
	}
	
}

func (m *Repository) ChatRoomPage(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "chat.jet.html", nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

func printTemplateError(w http.ResponseWriter, err error) {
	_, _ = fmt.Fprintf(w,`<small><span class='text-danger'>Error executing template: %s</span></small>`, err)
}

// ClientError will display error page for client error i.e. bad request
func ClientError(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		show404(w, r)
	case http.StatusInternalServerError:
		show500(w, r)
	default:
		http.Error(w, http.StatusText(status), status)
	}
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)
	show500(w, r)
}

func show404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/404.html")
}

func show500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}

