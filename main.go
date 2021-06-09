package main

import (
	"log"
	"os"

	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// .env files
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConnect := os.Getenv("DBCONNECT")

	database.Connect(dbConnect)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true, // set to true so the back end can pass the cookie to the front end
	}))

	routes(app)

	app.Listen("127.0.0.1:3000")
}