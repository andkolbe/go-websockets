package main

import (
	"github.com/andkolbe/go-websockets/internal/handlers"	
	"github.com/gofiber/fiber/v2"
)

func routes(app *fiber.App) {
	app.Get("/", handlers.Home)
}