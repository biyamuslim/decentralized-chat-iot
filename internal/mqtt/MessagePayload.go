package mqtt

// MessagePayload represents the message structure used for communication.
type MessagePayload struct {
	ClientName string `json:"client_name"`
	Message    string `json:"message"`
}
