package handlers

// import (
// 	"strconv"
// 	"time"

// 	"github.com/andkolbe/go-websockets/internal/database"
// 	"github.com/andkolbe/go-websockets/internal/models"
// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gofiber/fiber/v2"
// 	"golang.org/x/crypto/bcrypt"
// )

// type Claims struct {
// 	jwt.StandardClaims
// }

// func PostLogin(c *fiber.Ctx) error {
// 	// parse data we received from the request
// 	var data map[string]string
// 	err := c.BodyParser(&data)
// 	if err != nil {
// 		return err
// 	}

// 	// Get the user from the database by their email
// 	var user models.User

// 	database.DB.Where("email = ?", data["email"]).First(&user)

// 	// if the user is not found in the database
// 	if user.ID == 0 {
// 		c.Status(404)
// 		return c.JSON(fiber.Map{
// 			"message": "user not found",
// 		})
// 	}

// 	// check their password against the database
// 	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
// 		c.Status(400)
// 		return c.JSON(fiber.Map{
// 			"message": "incorrect password",
// 		})
// 	}

// 	// claims is the same as a payload
// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 		Issuer:    strconv.Itoa(int(user.ID)),            // need to convert id to string for this package to work. Id of the user is the Issuer
// 		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // expire token in 24 hrs
// 	})

// 	token, err := jwtToken.SignedString([]byte("secret"))
// 	if err != nil {
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}

// 	// store jwt in a cookie
// 	cookie := fiber.Cookie{
// 		Name:     "jwt",
// 		Value:    token,
// 		Expires:  time.Now().Add(time.Hour * 24),
// 		HTTPOnly: true,
// 	}

// 	c.Cookie(&cookie)

// 	return c.JSON(fiber.Map{
// 		"jwt": token,
// 	})
// }

// func PostRegister(c *fiber.Ctx) error {
// 	// parse data we received from the request
// 	var data map[string]string
// 	err := c.BodyParser(&data)
// 	if err != nil {
// 		return err
// 	}

// 	// password validation (add minimum length later!)
// 	if data["password"] != data["password_confirm"] {
// 		c.Status(400)
// 		return c.JSON(fiber.Map{
// 			"message": "passwords do not match",
// 		})
// 	}

// 	// hash password
// 	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

// 	user := models.User{
// 		UserName:  data["username"],
// 		FirstName: data["first_name"],
// 		LastName:  data["last_name"],
// 		Email:     data["email"],
// 		Password:  password,
// 	}

// 	// put this user in the database
// 	database.DB.Create(&user)

// 	return c.JSON(user)
// }

// // returns authenticated user
// func User(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt") // "jwt" key matches what was sent from the back end

// 	// decode the jwt
// 	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(t *jwt.Token) (interface{}, error) {
// 		return []byte("secret"), nil
// 	})

// 	if err != nil || !token.Valid {
// 		c.Status(fiber.StatusUnauthorized)
// 		return c.JSON(fiber.Map{
// 			"message": "unauthenticated",
// 		})
// 	}

// 	// if the token is valid, get the id off of it
// 	claims := token.Claims.(*Claims)

// 	// hold the user model in a variable
// 	var user models.User

// 	// Return the first user that matches the correct id
// 	database.DB.Where("id = ?", claims.Issuer).First(&user) // Id of the user is the Issuer

// 	// return all of the user data
// 	return c.JSON(user)
// }

// func Logout(c *fiber.Ctx) error {
// 	cookie := fiber.Cookie{
// 		Name:     "jwt",
// 		Value:    "",
// 		Expires:  time.Now().Add(-time.Hour), // expire the cookie one hour ago
// 		HTTPOnly: true,
// 	}

// 	c.Cookie(&cookie)

// 	return c.JSON(fiber.Map{
// 		"message": "success",
// 	})
// }
