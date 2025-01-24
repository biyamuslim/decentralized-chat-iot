package mqtt

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// Helper function to set up an in-memory SQLite database
func setupClientTableTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}

	// Ensure database connection is ready
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Create the clients table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		client_name TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create clients table: %v", err)
	}

	return db
}

// Test client creation (saving to DB)
func TestCreateClient(t *testing.T) {
	db := setupClientTableTestDB(t)
	defer db.Close()

	repo := NewClientRepository(db)
	client := &MqttClient{
		UserID:     1,
		ClientName: "client1",
	}

	err := repo.Create(client.UserID, client)
	if err != nil {
		t.Fatalf("Failed to save client: %v", err)
	}

	// Verify client exists in DB
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM clients WHERE client_name = ?", client.ClientName).Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("Client not inserted correctly in database")
	}
}

// Test retrieving a client by user ID
func TestGetByUserID(t *testing.T) {
	db := setupClientTableTestDB(t)
	defer db.Close()

	repo := NewClientRepository(db)

	// Insert test clients
	testClients := []MqttClient{
		{UserID: 1, ClientName: "client1"},
		{UserID: 2, ClientName: "client2"},
	}

	for _, client := range testClients {
		if err := repo.Create(client.UserID, &client); err != nil {
			t.Fatalf("Failed to insert test client: %v", err)
		}
	}

	// Fetch client by user ID
	client, err := repo.GetByUserID(1)
	if err != nil {
		t.Fatalf("Failed to retrieve client by UserID: %v", err)
	}

	if client.ClientName != "client1" {
		t.Errorf("Expected client name 'client1', but got '%s'", client.ClientName)
	}
}

// Test retrieving all clients
func TestGetAllClients(t *testing.T) {
	db := setupClientTableTestDB(t)
	defer db.Close()

	repo := NewClientRepository(db)

	// Insert test clients
	testClients := []MqttClient{
		{UserID: 1, ClientName: "client1"},
		{UserID: 2, ClientName: "client2"},
	}

	for _, client := range testClients {
		if err := repo.Create(client.UserID, &client); err != nil {
			t.Fatalf("Failed to insert test client: %v", err)
		}
	}

	// Fetch all clients
	clients, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to retrieve all clients: %v", err)
	}

	if len(clients) != 2 {
		t.Errorf("Expected 2 clients, but got %d", len(clients))
	}

	// Check if the retrieved clients match expected content
	expectedClients := []string{"client1", "client2"}
	for i, client := range clients {
		if client.ClientName != expectedClients[i] {
			t.Errorf("Expected client name '%s', but got '%s'", expectedClients[i], client.ClientName)
		}
	}
}
