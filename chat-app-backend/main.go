package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"chat-app-backend/models"
	"chat-app-backend/websockets"

	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
	var err error
	// Use Render-provided DATABASE_URL environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Create users table if not exists
	models.CreateUsersTable(db.DB)

	// Initialize WebSocket pool
	pool := websockets.NewPool()
	go pool.Start()

	// CORS Middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Update this if deploying frontend separately
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Define Routes
	http.HandleFunc("/login", loginHandler(db))
	http.HandleFunc("/register", RegisterUser(db))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWS(pool, w, r)
	})

	// **Render-Specific Port Handling**
	port := os.Getenv("PORT") // Get the port assigned by Render
	if port == "" {
		port = "8080" // Fallback for local development
	}

	fmt.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(http.DefaultServeMux)))
}
