package mqtt

import "fmt"

func (c *Client) Publish(topic, message string) {
	fmt.Printf("Subscribing to topic: '%s'\n", topic)
	token := c.mqttClient.Publish(topic, 0, false, message)
	token.Wait()
	fmt.Println("Published message:", message)
}

func (c *Client) PublishAndSave(clientRepo *ClientRepository, messageRepo *MessageRepository, clientID int64, topic, message string) {
	// Publish the message to the MQTT broker
	c.Publish(topic, message)

	// Save the message to the database
	msg := &Message{
		ClientID: clientID,
		Topic:    topic,
		Message:  message,
	}
	err := messageRepo.Create(msg)
	if err != nil {
		fmt.Println("Error saving message to database:", err)
	}
	fmt.Println("Message saved to the database!")
}
