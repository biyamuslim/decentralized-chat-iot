// internal/mqtt/message_repository.go
package mqtt

import (
	"database/sql"
)

// MessageRepository handles database operations for 'messages' table
type MessageRepository struct {
	DB *sql.DB
}

// NewMessageRepository creates a new instance of MessageRepository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// Create inserts a new message into the database
func (r *MessageRepository) Create(message *Message) error {
	stmt, err := r.DB.Prepare("INSERT INTO messages (client_id, topic, message) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(message.ClientID, message.Topic, message.Message)
	return err
}

// GetByTopic fetches all messages for a specific topic
func (r *MessageRepository) GetByTopic(topic string) ([]Message, error) {
	rows, err := r.DB.Query("SELECT id, client_id, topic, message, timestamp FROM messages WHERE topic = ?", topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ClientID, &message.Topic, &message.Message, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
