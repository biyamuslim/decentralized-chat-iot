package mqtt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// EncryptMessage encrypts the given message using AES encryption.
func EncryptMessage(key, message string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(message))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *Client) Publish(topic, message string) {
	token := c.mqttClient.Publish(topic, 0, false, message)
	token.Wait()
}

func (c *Client) PublishAndSave(clientRepo *ClientRepository, messageRepo *MessageRepository, clientID int64, topic, message string) {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	// Get the sender's client name
	client, err := clientRepo.GetByUserID(clientID)
	if err != nil {
		fmt.Println("Error retrieving client details:", err)
		return
	}

	// Create a structured message
	payload := MessagePayload{
		ClientName: client.ClientName,
		Message:    message,
	}

	// Convert to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return
	}

	// Encrypt the JSON message
	encryptedMessage, err := EncryptMessage(encryptionKey, string(jsonPayload))
	if err != nil {
		fmt.Println("Error encrypting message:", err)
		return
	}

	// Publish the encrypted message
	c.Publish(topic, encryptedMessage)

	// Save the encrypted message
	msg := &Message{
		ClientID: clientID,
		Topic:    topic,
		Message:  encryptedMessage,
	}
	err = messageRepo.Create(msg)
	if err != nil {
		fmt.Println("Error saving message to database:", err)
	}
}
