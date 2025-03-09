package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chat-app-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("mySuperSecretKey@123!") // Secure key (use environment variables in production)

// Structs for login and registration
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT Claims structure
type Claims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Generate JWT token
func generateToken(username, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 1 day
	claims := &Claims{
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// User login handler
// User login handler
func loginHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			log.Println("Invalid request format:", err)
			return
		}

		log.Println("Login Attempt - Email:", creds.Email)

		// Fetch user from database
		var user models.User
		err := db.Get(&user, "SELECT username, email, password FROM users WHERE email=$1", creds.Email)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			log.Println("User not found:", creds.Email)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}

		// Compare passwords
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			log.Println("Password mismatch for user:", creds.Email)
			return
		}

		// Generate JWT token
		token, err := generateToken(user.Username, user.Email)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			log.Println("JWT Token generation error:", err)
			return
		}

		// Send token and username in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":    token,         // ✅ Correct variable name
			"username": user.Username, // ✅ Correct variable name
		})
	}
}


// User registration handler
// User registration handler
func RegisterUser(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"message": "Invalid request format"}`, http.StatusBadRequest)
			log.Println("Invalid request format:", err)
			return
		}

		log.Println("Registration Attempt - Username:", req.Username, "Email:", req.Email)

		// Check if username already exists
		var existingUser models.User
		err := db.Get(&existingUser, "SELECT username FROM users WHERE username=$1", req.Username)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"message": "Username already exists"})
			log.Println("Username already exists:", req.Username)
			return
		}

		// Check if email already exists
		err = db.Get(&existingUser, "SELECT email FROM users WHERE email=$1", req.Email)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"message": "Email already exists"})
			log.Println("Email already exists:", req.Email)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Error hashing password"})
			log.Println("Error hashing password:", err)
			return
		}

		// Insert user into database
		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
			req.Username, req.Email, string(hashedPassword))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Error registering user"})
			log.Println("Database error during user registration:", err)
			return
		}

		// Success response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
		log.Println("User registered successfully:", req.Username)
	}
}
