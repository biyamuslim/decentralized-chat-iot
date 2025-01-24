package mqtt

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Helper function to set up an in-memory SQLite database
func setupMessageTableTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}

	// Ensure database connection is ready
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Create the messages table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id INTEGER NOT NULL,
		topic TEXT NOT NULL,
		message TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create messages table: %v", err)
	}

	return db
}

// Test message creation (saving to DB)
func TestCreateMessage(t *testing.T) {
	db := setupMessageTableTestDB(t)
	defer db.Close()

	repo := NewMessageRepository(db)
	encryptionKey := os.Getenv("ENCRYPTION_KEY") //encryption key
	encryptedMessage, err := EncryptMessage(encryptionKey, "Hello, world!")
	message := &Message{
		ClientID: 12345,
		Topic:    "test/topic",
		Message:  encryptedMessage,
	}

	err = repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}

	// Verify message exists in DB
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Messages WHERE topic = ?", message.Topic).Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("Message not inserted correctly in database")
	}
}

// Test retrieving messages by topic
func TestGetByTopic(t *testing.T) {
	db := setupMessageTableTestDB(t)
	defer db.Close()

	repo := NewMessageRepository(db)

	// Insert test messages
	testMessages := []Message{
		{ClientID: 1, Topic: "chat/general", Message: "Message 1"},
		{ClientID: 2, Topic: "chat/general", Message: "Message 2"},
		{ClientID: 3, Topic: "chat/private", Message: "Private message"},
	}

	for _, msg := range testMessages {
		if err := repo.Create(&msg); err != nil {
			t.Fatalf("Failed to insert test message: %v", err)
		}
	}

	// Fetch messages for topic "chat/general"
	messages, err := repo.GetByTopic("chat/general")
	if err != nil {
		t.Fatalf("Failed to retrieve messages: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, but got %d", len(messages))
	}

	// Check if the retrieved messages match expected content
	expectedMessages := []string{"Message 1", "Message 2"}
	for i, msg := range messages {
		if msg.Message != expectedMessages[i] {
			t.Errorf("Expected message '%s', but got '%s'", expectedMessages[i], msg.Message)
		}
	}
}

// Test timestamp field is set correctly
func TestMessageTimestamp(t *testing.T) {
	db := setupMessageTableTestDB(t)
	defer db.Close()

	repo := NewMessageRepository(db)
	message := &Message{
		ClientID: 456,
		Topic:    "test/time",
		Message:  "Timestamp check",
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}

	// Retrieve message and check timestamp
	var retrievedTimestamp time.Time
	err = db.QueryRow("SELECT timestamp FROM messages WHERE topic = ?", message.Topic).Scan(&retrievedTimestamp)
	if err != nil {
		t.Fatalf("Failed to retrieve message timestamp: %v", err)
	}

	if retrievedTimestamp.IsZero() {
		t.Errorf("Expected a valid timestamp, but got zero value")
	}
}
