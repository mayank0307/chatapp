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
	var dbURL string
	fmt.Println("starting database...")
	// Check if DATABASE_URL is set (for Render deployment)
	if os.Getenv("DATABASE_URL") != "" {
		dbURL = os.Getenv("DATABASE_URL") // Use Render-provided database URL
	} else {
		dbURL = "postgres://mayank:@localhost:5432/chatdb?sslmode=disable" // Local development database
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
		handlers.AllowedOrigins([]string{"https://chatapp-1-t7sw.onrender.com"}), // Allow your frontend
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Define Routes
	http.HandleFunc("/login", loginHandler(db))
	http.HandleFunc("/register", RegisterUser(db))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWS(pool, w, r)
	})

	// Use PORT from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "10001"
	}

	fmt.Println("Server running on port:  " + port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(http.DefaultServeMux)))
}
