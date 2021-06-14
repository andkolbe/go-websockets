package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/driver"
	"github.com/andkolbe/go-websockets/internal/handlers"
	"github.com/andkolbe/go-websockets/internal/helpers"
	"github.com/andkolbe/go-websockets/internal/models"
	"github.com/andkolbe/go-websockets/internal/render"
	"github.com/joho/godotenv"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run();
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close() // db won't close until the main function stops running. CAn't put it in run() because that only runs once when we open the app


	fmt.Printf("Starting application on port %s", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// .env files
	if err := godotenv.Load(); err != nil { log.Fatal("Error loading .env file") }
	dbConnect := os.Getenv("DBCONNECT")

	// what I am going to put in the session 
	gob.Register(models.User{})

	// print these to the terminal
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// enable sessions in the main package
	session = scs.New()
	session.Lifetime = 24 * time.Hour // active for 24 hours
	// stores the session in cookies by default. Can switch to Redis
	session.Cookie.Persist = true // cookie persists when the browser window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // makes sure the cookies are encrypted and use https. CHANGE TO TRUE FOR PRODUCTION

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dbConnect)
	if err != nil {
		log.Fatal("Cannot connect to db. Dying...")
	}
	log.Println("Connected to database!")

	// connect to ws
	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	// create the template cache only once when the application loads
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc

	// our app config (where we can put whatever we want) and our db (a pointer to a db driver) are available to all of our handlers
	// right now our db only holds postgres, but if we change or add more in the future, that can easily be refactored
	repo := handlers.NewRepo(&app, db)
	// pass the repo variable back to the handlers
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
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
