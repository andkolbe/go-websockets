package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// .env files
	if err := godotenv.Load(); err != nil { log.Fatal("Error loading .env file") }
	dbConnect := os.Getenv("DBCONNECT")

	database.Connect(dbConnect)
	log.Println("Connected to DB")

	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe("127.0.0.1:8080", mux)
}


/*
Calling go run main.go results in
	Your program is compliled inside a temporary folder
	The compiled binary is executed
But the temporary folder is just for one execution. 
So the next time when you run your program via go run another folder is used
	Change	
_ = http.ListenAndServe(":8080", mux)
_ = http.ListenAndServe("127.0.0.1:8080", mux)
*/
