package main

import (
	"fmt"
	"log"
	"net/http"

	"chat-app-backend/models"
	"chat-app-backend/websockets"

	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://mayank:@localhost:5432/chatdb?sslmode=disable")
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
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Allow frontend to access backend
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Define Routes
	http.HandleFunc("/login", loginHandler(db)) // ✅ Using loginHandler from auth.go
	http.HandleFunc("/register", RegisterUser(db)) // ✅ Using RegisterUser from auth.go
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWS(pool, w, r)
	})

	// Start server
	port := "8080"
	fmt.Println("Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(http.DefaultServeMux)))
}
