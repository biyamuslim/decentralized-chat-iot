package mqtt

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// UserRepository handles database operations for the 'users' table
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create registers a new user with a hashed password
func (r *UserRepository) Create(user *User) error {
	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := r.DB.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, string(hashedPassword))
	return err
}

// Authenticate checks if a user exists with the given credentials
func (r *UserRepository) Authenticate(username, password string) (*User, error) {
	var user User
	err := r.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Return nil if passwords don't match
	}

	return &user, nil
}
