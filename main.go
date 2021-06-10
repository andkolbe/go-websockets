package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/database"
	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/joho/godotenv"
)

var session *scs.SessionManager

func main() {

	err := run();
	if err != nil {
		log.Fatal(err)
	}

	mux := routes()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe("127.0.0.1:8080", mux)
}

func run() error {
	// .env files
	if err := godotenv.Load(); err != nil { log.Fatal("Error loading .env file") }
	dbConnect := os.Getenv("DBCONNECT")

	// enable sessions in the main package
	session = scs.New()
	session.Lifetime = 24 * time.Hour // active for 24 hours
	// stores the session in cookies by default. Can switch to Redis
	session.Cookie.Persist = true // cookie persists when the browser window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false // makes sure the cookies are encrypted and use https. CHANGE TO TRUE FOR PRODUCTION

	// connect to db
	database.Connect(dbConnect)
	log.Println("Connected to DB")

	// connect to ws
	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	return nil
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
