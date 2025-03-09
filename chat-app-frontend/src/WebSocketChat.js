import React, { useState, useEffect } from "react";

const WebSocketChat = () => {
  const [socket, setSocket] = useState(null);
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;
  
    const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
  
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prev) => [...prev, msg]);
    };
  
    setSocket(ws);
  
    return () => ws.close();
  }, []);
  

  const sendMessage = () => {
    if (socket && message) {
      const data = JSON.stringify(message); // ❌ This is wrong (message is a plain string)
      socket.send(data);
      setMessage("");
    }
  };
  

  return (
    <div>
      <h2>Chat Room</h2>
      <div>
        {messages.map((msg, index) => (
          <p key={index}>{msg}</p>
        ))}
      </div>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button onClick={sendMessage}>Send</button>
    </div>
  );
};

export default WebSocketChat;
