import { useState, useContext } from "react";
import { WebSocketContext } from "../App";

const Chat = () => {
  const { messages, sendMessage } = useContext(WebSocketContext);
  const [message, setMessage] = useState("");

  const handleSendMessage = () => {
    if (message.trim() !== "") {
      sendMessage(message);
      setMessage(""); // Clear input after sending
    }
  };

  return (
    <div className="flex flex-col min-h-screen items-center justify-center bg-gray-100">
      <h2 className="text-2xl font-bold mb-4">Chat Room</h2>
      
      {/* Messages List */}
      <div className="w-full max-w-md bg-white shadow-md rounded-lg p-4 h-64 overflow-auto border border-gray-300">
        {messages.length > 0 ? (
          messages.map((msg, index) => (
            <p key={index} className="mb-2">
              <strong className="text-blue-600">{msg.username}:</strong> {msg.text}
            </p>
          ))
        ) : (
          <p className="text-gray-500">No messages yet...</p>
        )}
      </div>

      {/* Message Input Box */}
      <div className="bg-white p-4 shadow-md flex items-center border-t w-full max-w-md mt-4">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type a message..."
          className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button
          onClick={handleSendMessage}
          className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition"
        >
          Send
        </button>
      </div>
    </div>
  );
};

export default Chat;
