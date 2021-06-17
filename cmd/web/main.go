package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/driver"
	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

var app config.AppConfig
var session *scs.SessionManager

// lets you store users in the session
func init() {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Birmingham")
}

func main() {

	// .env files
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// redisURL := os.Getenv("REDIS_URL")
	// redisTLSURL := os.Getenv("REDIS_TLS_URL")
	mySQLConnect := os.Getenv("MYSQLCONNECT")

	// CHANGE THIS TO TRUE WHEN IN PRODUCTION
	app.InProduction = false

	// Establish a redigo connection pool.
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	// enable sessions in the main package
	// add redis store
	log.Printf("Initializing session manager....")
	session = scs.New()
	session.Store = redisstore.New(pool)
	session.Lifetime = 2 * time.Hour // active for 2 hours
	// stores the session in cookies by default. Can switch to Redis
	// session.Cookie.Persist = true // cookie persists when the browser window is closed
	// session.Cookie.SameSite = http.SameSiteLaxMode
	// session.Cookie.Secure = app.InProduction // makes sure the cookies are encrypted and use https. CHANGE TO TRUE FOR PRODUCTION

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(mySQLConnect)
	if err != nil {
		log.Fatal("Cannot connect to db. Dying...")
	}
	log.Println("Connected to database!")

	// connect to ws
	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	// our app config (where we can put whatever we want) and our db (a pointer to a db driver) are available to all of our handlers
	// right now our db only holds postgres, but if we change or add more in the future, that can easily be refactored
	repo := handlers.NewRepo(&app, db)
	// pass the repo variable back to the handlers
	handlers.NewHandlers(repo)

	defer db.SQL.Close() // db won't close until the main function stops running

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
	Change
_ = http.ListenAndServe(":8080", mux)
_ = http.ListenAndServe("127.0.0.1:8080", mux)
*/