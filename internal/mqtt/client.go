package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	mqttClient mqtt.Client
}

func NewClient(broker, clientID string) *Client {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to broker:", broker)
	return &Client{mqttClient: client}
}

func (c *Client) Disconnect() {
	c.mqttClient.Disconnect(250)
	fmt.Println("Client disconnected")
}
