package handlers

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

// we need a lot of information before we can test our handlers
// responseWriter, request, access to session, all of our routes, middleware

var app config.AppConfig // holds app configuration
var session *scs.SessionManager

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

	helpers.SetViews("./../../views") // allows us to use jet in our testing

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
	mux.Use(RecoverPanic)

	mux.Get("/", Repo.LoginPage)
	mux.Post("/", Repo.Login)
	mux.Get("/register", Repo.RegisterPage)
	mux.Post("/register", Repo.Register)
	mux.Get("/logout", Repo.Logout)
	mux.Get("/ws", WsEndPoint)

	mux.Get("/auth/chat", Repo.ChatRoomPage)

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

// Auth checks for authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "userID") {
			
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// RecoverPanic recovers from a panic
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Check if there has been a panic
			if err := recover(); err != nil {
				// return a 500 Internal Server response
				helpers.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}