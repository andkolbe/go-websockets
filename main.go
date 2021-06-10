package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
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

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config {
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true, // set to true so the back end can pass the cookie to the front end
	}))

	mux := routes()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe("127.0.0.1:8080", mux)
}


/*
Calling go run main.go results in
	Your program is compliled inside a temporary folder
	The compiled binary is executed
But the temporary folder is just for one execution. 
So the next time when you run your program via go run another folder is used
*/
