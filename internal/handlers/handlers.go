package handlers

import (
	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)


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
	// parse data we received from the request
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	
	// password validation (add minimum length later!)
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map {
			"message": "passwords do not match",
		})
	}

	// hash password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
 
	user := models.User {
		UserName: data["username"],
		FirstName: data["first_name"],
		LastName: data["last_name"],
		Email: data["email"],
		Password: password,
	}

	// put this user in the database
	database.DB.Create(&user)

	return c.JSON(user)
}

func Chat(c *fiber.Ctx) error {
	return c.SendString("Chat")
}