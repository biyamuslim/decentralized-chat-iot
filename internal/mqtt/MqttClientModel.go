package mqtt

type MqttClient struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`     // Corresponds to the user_id field in the table
	ClientName string `json:"client_name"` // Corresponds to the client_name field in the table
}
