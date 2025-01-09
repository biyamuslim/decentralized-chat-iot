package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (c *Client) Subscribe(topic string) {
	token := c.mqttClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message on topic '%s': %s\n", msg.Topic(), string(msg.Payload()))
	})
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Error subscribing to topic '%s': %v\n", topic, token.Error())
	} else {
		fmt.Printf("Successfully subscribed to topic: %s\n", topic)
	}
}
