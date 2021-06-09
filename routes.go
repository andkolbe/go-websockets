package main

import (
	"github.com/andkolbe/go-websockets/internal/handlers"	
	"github.com/gofiber/fiber/v2"
)

func routes(app *fiber.App) {
	app.Get("/", handlers.Login)
	app.Post("/", handlers.PostLogin)
	app.Get("/register", handlers.Register)
	app.Post("/register", handlers.PostRegister)
	app.Get("/chat", handlers.Chat)
	app.Get("/api/user", handlers.User)
	app.Post("/api/logout", handlers.Logout)
	// app.Post("/api/forgot", handlers.Forgot)
	// app.Post("/api/reset", handlers.Reset)
}