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
	mux.Use(RecoverPanic)

	mux.Get("/", handlers.Repo.LoginPage)
	mux.Post("/", handlers.Repo.Login)
	mux.Get("/register", handlers.Repo.RegisterPage)
	mux.Post("/register", handlers.Repo.Register)
	mux.Get("/logout", handlers.Repo.Logout)
	mux.Get("/ws", handlers.WsEndPoint)

	// auth routes
	mux.Route("/auth", func(mux chi.Router) {
		// all admin routes are protected
		mux.Use(Auth)
		mux.Get("/chat", handlers.Repo.ChatRoomPage)
	})

	// if a user is disconnected, and then reconnects, they rejoin automatically
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
