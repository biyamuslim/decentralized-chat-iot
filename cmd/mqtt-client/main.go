package main

import (
	"bufio"
	"decentralized-chat-iot/internal/mqtt"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Initialization
	mqtt.InitDB()
	userRepo := mqtt.NewUserRepository(mqtt.GetDB())
	clientRepo := mqtt.NewClientRepository(mqtt.GetDB())
	messageRepo := mqtt.NewMessageRepository(mqtt.GetDB())

	var user *mqtt.User
	var client *mqtt.MqttClient

	// Handle authentication
	for {
		fmt.Println("Welcome! Choose an option:")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		if choice == 1 {
			user = handleLogin(userRepo)
			if user != nil {
				client, _ = clientRepo.GetByUserID(user.ID)
				if client != nil {
					break
				}
				fmt.Println("No client found for this user. Please register.")
			}
		} else if choice == 2 {
			user = handleRegistration(userRepo, clientRepo)
			if user != nil {
				client, _ = clientRepo.GetByUserID(user.ID)
				break
			}
		} else {
			fmt.Println("Invalid choice. Please try again.")
		}
	}

	if client == nil {
		log.Fatal("Unable to initialize MQTT client without valid client data.")
	}

	// Initialize the MQTT client using ClientID
	mqttClient := mqtt.NewClient("tcp://localhost:1883", strconv.FormatInt(client.ID, 10))

	// Subscribe to a topic in a goroutine
	go mqttClient.Subscribe("test/topic")

	// Start input loop to send messages dynamically
	fmt.Println("Type your messages below. Press Ctrl+C to exit.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		if message == "" {
			continue // Skip empty messages
		}

		// Publish and save the message using the client ID
		mqttClient.PublishAndSave(clientRepo, messageRepo, client.ID, "test/topic", message)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}

	select {}
}

func handleLogin(userRepo *mqtt.UserRepository) *mqtt.User {
	fmt.Print("Enter username: ")
	var username string
	fmt.Scan(&username)

	fmt.Print("Enter password: ")
	var password string
	fmt.Scan(&password)

	user, err := userRepo.Authenticate(username, password)
	if err != nil {
		fmt.Println("Login failed. Please try again.")
		return nil
	}

	fmt.Println("Login successful!")
	return user
}

func handleRegistration(userRepo *mqtt.UserRepository, clientRepo *mqtt.ClientRepository) *mqtt.User {
	fmt.Print("Enter new username: ")
	var username string
	fmt.Scan(&username)

	fmt.Print("Enter new password: ")
	var password string
	fmt.Scan(&password)

	newUser := &mqtt.User{Username: username, Password: password}
	err := userRepo.Create(newUser)
	if err != nil {
		fmt.Println("Registration failed. Username might already exist.")
		return nil
	}

	// Get the user ID of the newly registered user
	user, err := userRepo.Authenticate(username, password)
	if err != nil {
		fmt.Println("Registration error. Please try again.")
		return nil
	}

	fmt.Print("Enter a name for your client: ")
	var clientName string
	fmt.Scan(&clientName)

	client := &mqtt.MqttClient{ClientName: clientName}
	err = clientRepo.Create(user.ID, client)
	if err != nil {
		fmt.Println("Failed to create a client for the new user.")
		return nil
	}

	fmt.Println("Registration successful!")
	return user
}
