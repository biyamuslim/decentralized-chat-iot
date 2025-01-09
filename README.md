# Decentralized Chat System for IoT

A decentralized chat system leveraging MQTT, built using Go, designed for IoT applications.

## Table of Contents
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

---

## Getting Started
Follow the steps below to set up and run the decentralized chat system.

## Prerequisites
- Latest version of [Go](https://go.dev/dl/) installed on your system.
- SQLite3 installed (for database support).

## Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/decentralized-chat-iot.git
   cd decentralized-chat-iot
   ```

2. **Initialize Go Modules**
   ```bash
   go mod tidy
   ```

## Configuration

Create a `.env` file in the root directory and add the following keys:
```env
DATABASE_PATH=./app.db
MIGRATION_PATH=./migrations
ENCRYPTION_KEY="your-32-byte-secret-key-for-aes!"
```

## Running the Application

1. **Run Migrations**

   Install `goose` for database migrations:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

   Run migrations using the following command:
   ```bash
   goose -dir=migrations sqlite3 ./app.db up
   ```

2. **Build the Project**
   ```bash
   go build -o main ./cmd/mqtt-client
   ```

3. **Run the Application**

   Open multiple terminal instances and run the following command in each terminal:
   ```bash
   go run ./cmd/mqtt-client
   ```

   Upon running, you will see two options:
   1. **Login**
   2. **Register**

   - Use the **Register** option to create a client.
   - Repeat the process in other terminals to register more clients.

4. **Chat Across Terminals**
   After registering, you can start writing messages in one terminal. The messages will be broadcast to all other connected terminals.

## Usage

- **Register a Client**: Each terminal represents a client. Register using the "Register" option.
- **Login**: Use the login option to reconnect an existing client.
- **Broadcast Messages**: Any message sent from one terminal will appear in all other terminals.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
