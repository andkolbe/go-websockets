package handlers

// import "net/http"

// we need a lot of information before we can test our handlers
// responseWriter, request, access to session

// func getRoutes() http.Handler {
// 	// .env files
// 	if err := godotenv.Load(); err != nil { log.Fatal("Error loading .env file") }
// 	dbConnect := os.Getenv("DBCONNECT")

// 	database.Connect(dbConnect)
// 	log.Println("Connected to DB")

// 	log.Println("Starting channel listener")
// 	go handlers.ListenToWSChannel()

// 	return nil
// }