package websockets

import "fmt"

// Pool manages all WebSocket connections
type Pool struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

// NewPool creates a new Pool instance
func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Start listens for client connections and messages
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("✅ New client joined. Total clients:", len(pool.Clients))
			go client.WriteMessages() // ✅ Start listening for outgoing messages

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			close(client.Send) // ✅ Close Send channel
			fmt.Println("❌ Client disconnected. Total clients:", len(pool.Clients))

		case message := <-pool.Broadcast:
			fmt.Println("📢 Broadcasting message from", message.Username, ":", message.Text) // ✅ Fix field names
			for client := range pool.Clients {
				err := client.Conn.WriteJSON(message) // ✅ Send full message object
				if err != nil {
					fmt.Println("❌ Error sending message:", err)
					client.Conn.Close()
					delete(pool.Clients, client)
				}
			}
		
		}
	}
}
