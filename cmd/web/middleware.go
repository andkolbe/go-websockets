package main

import (
	"fmt"
	"net/http"

	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/justinas/nosurf"
)

// adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// uses cookies to make sure the token it generates for us is on a per page basis
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", // "/" means apply this cookie to the entire site
		Secure:   app.InProduction, // CHANGE TO TRUE IN PRODUCTION
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// webservers are not state aware, so we need to add middleware that tells this application to remember state using sessions
// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth checks for authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			url := r.URL.Path
			http.Redirect(w, r, fmt.Sprintf("/?target=%s", url), http.StatusFound)
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