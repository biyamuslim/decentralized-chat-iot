package mqtt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DecryptMessage decrypts the given message using AES.
func DecryptMessage(key, encryptedMessage string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func (c *Client) Subscribe(topic string) {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	token := c.mqttClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		encryptedMessage := string(msg.Payload())

		// Decrypt the message
		decryptedJSON, err := DecryptMessage(encryptionKey, encryptedMessage)
		if err != nil {
			fmt.Printf("Failed to decrypt message on topic '%s': %v\n", msg.Topic(), err)
			return
		}

		// Parse JSON to extract sender's client name
		var payload MessagePayload
		err = json.Unmarshal([]byte(decryptedJSON), &payload)
		if err != nil {
			fmt.Printf("Failed to parse message JSON on topic '%s': %v\n", msg.Topic(), err)
			return
		}

		// Display the message with the correct sender's client name
		fmt.Printf("Received message on topic '%s' from %s: %s\n", msg.Topic(), payload.ClientName, payload.Message)
	})

	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Error subscribing to topic '%s': %v\n", topic, token.Error())
	} else {
		fmt.Printf("Successfully subscribed to topic: %s\n", topic)
	}
}
