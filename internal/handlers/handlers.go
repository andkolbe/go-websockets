package handlers

import "github.com/gofiber/fiber/v2"

func Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

func Register(c *fiber.Ctx) error {
	return c.SendString("Register")
}

func Chat(c *fiber.Ctx) error {
	return c.SendString("Chat")
}
