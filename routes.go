package main

import (
	"net/http"

	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Login))
	// mux.Post("/", http.HandlerFunc(handlers.PostLogin))
	mux.Get("/register", http.HandlerFunc(handlers.Register))
	// mux.Post("/register", http.HandlerFunc(handlers.PostRegister))
	mux.Get("/chat", http.HandlerFunc(handlers.Chat))
	// mux.Get("/user", http.HandlerFunc(handlers.User))
	// mux.Post("/logout", http.HandlerFunc(handlers.Logout))
	// mux.Post("/forgot", http.HandlerFunc(handlers.Forgot))
	// mux.Post("/reset", http.HandlerFunc(handlers.Reset))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))
	
	return mux
}