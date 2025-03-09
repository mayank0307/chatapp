import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { createContext, useEffect, useState } from "react";
import Login from "./pages/Login";
import Chat from "./pages/Chat";
import Register from "./pages/Register";

export const WebSocketContext = createContext(null);

const WebSocketProvider = ({ children }) => {
  const [ws, setWs] = useState(null);
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return; // Don't connect if the user is not authenticated

    const backendURL = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";
    const wsProtocol = backendURL.startsWith("https") ? "wss" : "ws"; // Use "wss" for secure connections
    const wsURL = `${wsProtocol}://${backendURL.replace(/^https?:\/\//, '')}/ws?token=${token}`;

    console.log("🔗 Connecting to WebSocket:", wsURL);
    
    const socket = new WebSocket(wsURL);

    socket.onopen = () => console.log("✅ WebSocket Connected!");
    
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      console.log("📩 New Message:", msg);
      setMessages((prev) => [...prev, msg]);
    };

    socket.onclose = () => {
      console.warn("⚠️ WebSocket Disconnected! Reconnecting...");
      setTimeout(() => {
        setWs(new WebSocket(wsURL));
      }, 3000);
    };

    socket.onerror = (error) => console.error("❌ WebSocket Error:", error);

    setWs(socket);

    return () => socket.close();
  }, []);

  // Function to send messages
  const sendMessage = (message) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      const username = localStorage.getItem("username") || "Guest";
      const messageObj = { text: message, username };
      ws.send(JSON.stringify(messageObj));
    } else {
      console.error("❌ WebSocket is not connected!");
    }
  };

  return (
    <WebSocketContext.Provider value={{ messages, sendMessage }}>
      {children}
    </WebSocketContext.Provider>
  );
};

const App = () => {
  return (
    <WebSocketProvider>
      <Router>
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/chat" element={<Chat />} />
        </Routes>
      </Router>
    </WebSocketProvider>
  );
};

export default App;
