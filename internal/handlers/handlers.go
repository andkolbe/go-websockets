package handlers

import (
	"strconv"
	"time"

	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

func PostLogin(c *fiber.Ctx) error {
	// parse data we received from the request
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	// Get the user from the database by their email
	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	// if the user is not found in the database
	if user.ID == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// check their password against the database
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// claims is the same as a payload
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),            // need to convert id to string for this package to work. Id of the user is the Issuer
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // expire token in 24 hrs
	})

	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// store jwt in a cookie
	cookie := fiber.Cookie {
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"jwt": token,
	})
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
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	// hash password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		UserName:  data["username"],
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	// put this user in the database
	database.DB.Create(&user)

	return c.JSON(user)
}

func Chat(c *fiber.Ctx) error {
	return c.SendString("Chat")
}
