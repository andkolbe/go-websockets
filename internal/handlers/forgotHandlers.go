package handlers

import (
	"math/rand"
	"net/smtp"

	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)


func Forgot(c *fiber.Ctx) error {
	// parse the data from the request
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// we have to generate a token for them so they can keep going
	token := RandStringRunes(12)

	passwordReset := models.PasswordReset {
		Email: data["email"], // the email that we got from the request
		Token: token,
	}

	// query database
	database.DB.Create(&passwordReset)

	// send reset email
	from := "admin@go_websockets.com"
	to := []string {
		data["email"],
	}
	url := "http://localhost:3000/reset/" + token
	message := []byte("Click <a href=\"" + url + "\">here</a> to reset your password")
	err := smtp.SendMail("0.0.0.0:1025", nil, from, to, message)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map {
		"message": "success",
	})
}

func Reset(c *fiber.Ctx) error {
	// parse the data from the request
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// password validation
	if data["password"] != data["password_confirm"] {
		c.Status(400) 
		return c.JSON(fiber.Map {
			"message": "Passwords do not match",
		})
	}

	var passwordReset = models.PasswordReset{}

	// if the user clicks on the password reset button a lot, only change the password with the token on the last attempt
	if err := database.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		c.Status(400) 
		return c.JSON(fiber.Map {
			"message": "Invalid token!",
		})
	}

	// hash the new password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// update users table with new password
	database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	return c.JSON(fiber.Map {
		"message": "success",
	})
}

// we have to generate a token ourselves
func RandStringRunes(n int) string {
	// rune is an integer with 32 characters 
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	
	b := make([]rune, n) // make a slice of runes however long the number that was passed in as a parameter
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}