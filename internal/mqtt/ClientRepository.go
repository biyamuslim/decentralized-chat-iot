package mqtt

import (
	"database/sql"
)

// ClientRepository handles database operations for 'clients' table
type ClientRepository struct {
	DB *sql.DB
}

// NewClientRepository creates a new instance of ClientRepository
func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{DB: db}
}

// Create inserts a new client into the database for a specific user
func (r *ClientRepository) Create(userID int64, client *MqttClient) error {
	stmt, err := r.DB.Prepare("INSERT INTO clients (user_id, client_name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, client.ClientName)
	return err
}

// GetByID fetches a client by ID
func (r *ClientRepository) GetByUserID(UserID int64) (*MqttClient, error) {
	var client MqttClient
	err := r.DB.QueryRow("SELECT id, client_name FROM clients WHERE user_id = ?", UserID).Scan(&client.ID, &client.ClientName)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetAll fetches all clients
func (r *ClientRepository) GetAll() ([]MqttClient, error) {
	rows, err := r.DB.Query("SELECT id, client_name FROM clients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []MqttClient
	for rows.Next() {
		var client MqttClient
		if err := rows.Scan(&client.ID, &client.ClientName); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}
