package handlers

import (
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/gofiber/fiber/v2")


func Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

func PostLogin(c *fiber.Ctx) error {
	return c.SendString("Login")
}

func Register(c *fiber.Ctx) error {
	return c.SendString("Register")
}

func PostRegister(c *fiber.Ctx) error {
	user := models.User{
		
	}

	return c.JSON(user)
}

func Chat(c *fiber.Ctx) error {
	return c.SendString("Chat")
}