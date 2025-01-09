package mqtt

import "time"

// Client represents the structure of the 'clients' table
type MqttClient struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`     // Corresponds to the user_id field in the table
	ClientName string `json:"client_name"` // Corresponds to the client_name field in the table
}

// Message represents the structure of the 'messages' table
type Message struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Topic     string    `json:"topic"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// User represents the structure of the 'users' table
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
