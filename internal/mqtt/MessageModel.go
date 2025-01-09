package mqtt

import "time"

type Message struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Topic     string    `json:"topic"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
