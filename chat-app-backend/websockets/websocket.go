package websockets

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

// Upgrader config for WebSocket
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// ServeWS handles new WebSocket connections
func ServeWS(pool *Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	email := r.URL.Query().Get("email") // Get email from query params
	log.Println("📧 New client connected with email:", email)
	if email == "" {
		email = "guest@example.com" // Default if no email is provided
	}

	client := &Client{
		Conn:  conn,
		Pool:  pool,
		Email: email, // Store email in Client struct
		Send:  make(chan Message),
	}

	pool.Register <- client // Register the client in the pool

	go client.ReadMessages()
	go client.WriteMessages()
}

