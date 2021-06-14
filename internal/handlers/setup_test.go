package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

// we need a lot of information before we can test our handlers
// responseWriter, request, access to session, all of our routes, middleware

var app config.AppConfig // holds app configuration
var session *scs.SessionManager
var pathToTemplates = "./../../views"
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	// // .env files
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// dbConnect := os.Getenv("DBCONNECT")

	gob.Register(models.User{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour 
	session.Cookie.Persist = true     
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false // not secure for testing

	app.Session = session

	// create template cache
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	// // connect to ws
	// log.Println("Starting channel listener")
	// go handlers.ListenToWSChannel()

	os.Exit(m.Run()) // exit all of the tests when they are finished
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.LoginPage)
	mux.Post("/", Repo.Login)
	mux.Get("/register", Repo.RegisterPage)
	// mux.Post("/register", PostRegister)
	mux.Get("/chat", Repo.ChatRoomPage)
	// mux.Get("/user", User)
	mux.Post("/logout", Repo.Logout)
	// mux.Post("/forgot", Forgot)
	// mux.Post("/reset", Reset)
	mux.Get("/ws", WsEndPoint)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", 
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// creates a test template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// create a template cache that holds all our html templates in a map
	myCache := map[string]*template.Template{} // map with an index of type string and its contents are a pointer to template.Template

	// go to the templates folder, and get all of the files that start with anything but end with .page.html
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err

		}

		// go to the templates folder, and get all of the files that end with .layout.html
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err

		}

		// if a .layout.html match is found, the length will be greater than 0
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err

			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

