package main

import (
	"database/sql"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	_ "github.com/lib/pq"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/andkolbe/go-websockets/internal/models"
)

var app config.AppConfig
var session *scs.SessionManager

// lets you store users in the session
func init() {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Birmingham")
}

func main() {

	// env vars
	dns := os.Getenv("HEROKU_POSTGRESQL_AMBER_URL")
	if dns == "" {
		log.Fatal("DATABASE_URL is not set.")
	}
	port := os.Getenv("PORT") // heroku will take this and use their own port number
	if port == "" {
		log.Fatal("$PORT is not set.")
	}

	// CHANGE THIS TO TRUE WHEN IN PRODUCTION
	app.InProduction = true

	log.Println("Connecting to database...")
	db, err := sql.Open("pgx", dns)
	if err != nil {
		log.Fatal("Cannot connect to db. Dying...")
	}
	log.Println("Connected to database!")

	// set some parameters on the db connection pool
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	defer db.Close() // whatever comes after the defer keyword, execute that as soon as the function is done

	// enable sessions in the main package
	log.Printf("Initializing session manager....")
	session = scs.New()
	session.Lifetime = 2 * time.Hour // active for 2 hours
	// stores the session in cookies by default. Can switch to Redis
	session.Cookie.Persist = true // cookie persists when the browser window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // makes sure the cookies are encrypted and use https. CHANGE TO TRUE FOR PRODUCTION
	app.Session = session

	// connect to ws
	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	// our app config (where we can put whatever we want) and our db (a pointer to a db driver) are available to all of our handlers
	// right now our db only holds postgres, but if we change or add more in the future, that can easily be refactored
	repo := handlers.NewRepo(&app, db)
	// pass the repo variable back to the handlers
	handlers.NewHandlers(repo)

	mux := routes()

	log.Println("Starting web server")

	_ = http.ListenAndServe(":"+port, mux)
	// _ = http.ListenAndServe("127.0.0.1:8080", mux)
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

/*
	Setting .env vars in Go
	type into cli: set PORT=portnumber
	Check if PORT has been set: echo %PORT%
	Do the same with DATABASE_URL

	In heroku format the CLEARDB_DATABASE_URL from mysql://alphanum-username:alphanum-password@us-cdbr-iron-east-01.cleardb.net/heroku_alphanum_name?reconnect=true
	to DATABASE_URL bf72ed7dc75b1c:982c4139@tcp(us-cdbr-east-04.cleardb.com)/heroku_b7277b2cda8d5c4

*/
