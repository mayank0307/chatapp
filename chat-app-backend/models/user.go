package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// User struct with username
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"` // Add Username field
	Email    string `db:"email"`
	Password string `db:"password"` // Store hashed password
}


// CreateUsersTable initializes the users table
func CreateUsersTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating users table:", err)
	}
}