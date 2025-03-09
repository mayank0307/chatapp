package websockets

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Message struct to include sender's username
type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}
// Client represents a WebSocket connection
type Client struct {
	Conn  *websocket.Conn
	Pool  *Pool
	Email string // ✅ Store email in Client struct
	Send  chan Message
}

// ReadMessages listens for messages from the client
func (c *Client) ReadMessages() {
	defer func() {
		c.Pool.Unregister <- c // Remove client on disconnect
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("❌ JSON Decode Error:", err)
			break
		}

		// Use the username sent by the client (instead of extracting from email)
		msgToSend := Message{
			Username: msg.Username,  // ✅ Use the username provided by frontend
			Text:     msg.Text,
		}

		fmt.Println("📩 Message received from", msg.Username, ":", msg.Text)

		// Send to broadcast channel
		c.Pool.Broadcast <- msgToSend
	}
}


// WriteMessages listens on the Send channel and writes messages to the WebSocket
func (c *Client) WriteMessages() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("❌ Error sending message:", err)
			break
		}
	}
}
