package mqtt

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to set up an in-memory SQLite database
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}

	// Create a users table
	createTableSQL := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

// Test user registration (valid case)
func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{Username: "testuser", Password: "password123"}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user exists in DB
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("User not inserted correctly in database")
	}
}

// Test duplicate user registration
func TestCreateDuplicateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{Username: "duplicateUser", Password: "password123"}

	// Insert user twice
	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	err = repo.Create(user) // Should fail due to UNIQUE constraint
	if err == nil {
		t.Errorf("Expected error when inserting duplicate user, but got nil")
	}
}

// Test user authentication (correct credentials)
func TestAuthenticateUser_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create a user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "testuser", string(hashedPassword))
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Authenticate with correct password
	user, err := repo.Authenticate("testuser", "password123")
	if err != nil || user == nil {
		t.Errorf("Expected authentication to succeed but failed: %v", err)
	}
}

// Test user authentication (incorrect password)
func TestAuthenticateUser_Fail_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create a user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "testuser", string(hashedPassword))
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Attempt authentication with wrong password
	_, err = repo.Authenticate("testuser", "wrongpassword")
	if err == nil {
		t.Errorf("Expected authentication to fail with wrong password, but it succeeded")
	}
}

// Test user authentication (non-existent user)
func TestAuthenticateUser_Fail_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Attempt authentication with a non-existent user
	_, err := repo.Authenticate("unknownuser", "password123")
	if err == nil {
		t.Errorf("Expected authentication to fail for non-existent user, but it succeeded")
	}
}
