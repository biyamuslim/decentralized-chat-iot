package mqtt

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// Global variable to hold the database connection
var db *sql.DB

// GetDB returns the initialized database connection
func GetDB() *sql.DB {
	return db
}

// Initialize the SQLite database connection
func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get database path from .env
	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		log.Fatal("DATABASE_PATH not set in .env file")
	}

	// Open the SQLite database connection
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Check if the database is reachable
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("Successfully connected to the SQLite database!")
}
