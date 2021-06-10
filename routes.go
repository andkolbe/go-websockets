package main

import (
	"net/http"

	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	// middleware allows you process a request as it comes into your web app and perform some action on it
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Login)
	// mux.Post("/", handlers.PostLogin)
	mux.Get("/register", handlers.Register)
	// mux.Post("/register", handlers.PostRegister)
	mux.Get("/chat", handlers.Chat)
	// mux.Get("/user", handlers.User)
	// mux.Post("/logout", handlers.Logout)
	// mux.Post("/forgot", handlers.Forgot)
	// mux.Post("/reset", handlers.Reset)
	mux.Get("/ws", handlers.WsEndPoint)

	// if a user is disconnected, and then reconnects, they rejoin automatically
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
