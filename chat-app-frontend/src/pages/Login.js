import { useState } from "react";
import { useNavigate } from "react-router-dom";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate(); // Used for redirection

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const backendURL = process.env.REACT_APP_BACKEND_URL;
      const response = await fetch(`${backendURL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        throw new Error("Invalid credentials");
      }

      const data = await response.json();

      // Ensure the backend sends the username along with the token
      if (!data.token || !data.username) {
        throw new Error("Missing username or token from server response");
      }

      // Save token and username to localStorage
      localStorage.setItem("token", data.token);
      localStorage.setItem("username", data.username);

      setMessage("✅ Login successful! Redirecting...");

      setTimeout(() => {
        navigate("/chat"); // Redirect to chat page
      }, 1500);
    } catch (error) {
      setMessage(`❌ ${error.message}`);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white shadow-md rounded-lg p-6 w-96">
        <h2 className="text-2xl font-semibold text-center text-gray-700 mb-4">
          Login
        </h2>
        <form onSubmit={handleLogin} className="space-y-4">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 transition"
          >
            Login
          </button>
        </form>
        {message && (
          <p className="text-center text-sm mt-2 text-gray-600">{message}</p>
        )}
        <p className="text-center text-gray-600 mt-4">
          Don't have an account?{" "}
          <a href="/register" className="text-blue-500 hover:underline">
            Register
          </a>
        </p>
      </div>
    </div>
  );
};

export default Login;
