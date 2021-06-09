package main

import (
	"log"
	"os"

	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/gofiber/fiber/v2"
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

	routes(app)

	app.Listen(":3000")
}