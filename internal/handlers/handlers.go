package handlers

import (
	// "log"
	// "net/http"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map {})
}

func Register(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map {})
}

func Chat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map {})
}


