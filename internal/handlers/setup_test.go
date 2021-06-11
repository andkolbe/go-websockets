package handlers

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

// we need a lot of information before we can test our handlers
// responseWriter, request, access to session, all of our routes, middleware

var app config.AppConfig // holds app configuration
var testSession *scs.SessionManager

func TestMain(m *testing.M) {
	// // .env files
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// dbConnect := os.Getenv("DBCONNECT")

	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour 
	testSession.Cookie.Persist = true     
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false // not secure for testing

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
	mux.Get("/", Repo.LoginPage)
	// mux.Post("/", PostLogin)
	mux.Get("/register", Repo.RegisterPage)
	// mux.Post("/register", PostRegister)
	mux.Get("/chat", Repo.ChatRoomPage)
	// mux.Get("/user", User)
	// mux.Post("/logout", Logout)
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
		Path:     "/", // "/" means apply this cookie to the entire site
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return testSession.LoadAndSave(next)
}
