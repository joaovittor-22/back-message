# Go WebSocket Server with PostgreSQL Integration

This is a simple WebSocket server built using Go, GORM, and PostgreSQL. The server listens on a specified port, allows WebSocket connections, saves incoming messages to the PostgreSQL database, and broadcasts messages to all connected WebSocket clients.

## Features

- WebSocket support for real-time messaging.
- PostgreSQL database for storing messages.
- CORS support for cross-origin requests.
- Environment variable-based configuration.

## Requirements

- Go (version 1.18 or higher)
- PostgreSQL (installed and running locally or remotely)

## Setup and Running the Project

### 1. Clone the Repository

First, clone this repository to your local machine.

```bash
git clone https://github.com/yourusername/go-websocket-postgres.git
cd go-websocket-postgres

2. Install Dependencies

Run the following command to install all the required Go dependencies:

go mod tidy

This will download the necessary packages such as gorilla/websocket, gorm, and rs/cors.
3. Set Up the PostgreSQL Database

Make sure you have PostgreSQL installed and running. You need to create a PostgreSQL database and a user for this project:

4. Create .env File

Create a .env file in the root of your project directory to store environment variables. Here's an example .env file:

# PostgreSQL Database Configuration
DB_HOST=<database_host>
DB_PORT=<port_number>
DB_USER=<db_user>
DB_PASSWORD=<password>
DB_NAME=<db_name>

# Server Configuration
SERVER_PORT=8080

5. Install godotenv Package (If Not Already Installed)

This project uses the godotenv package to load environment variables. You can install it by running:

go get github.com/joho/godotenv

6. Run the Project

To run the application, use the following command:

go run main.go

The server will start and listen on port 8080 by default.
7. Access the Application

    WebSocket Endpoint: ws://localhost:8080/ws

        This endpoint allows WebSocket connections to send and receive real-time messages.

    HTTP API Endpoint: http://localhost:8080/getMessages

        This endpoint allows you to retrieve all stored messages in JSON format.

8. Test the WebSocket Connection

You can use tools like Postman or a WebSocket client to test the WebSocket connection:

    Connect to ws://localhost:8080/ws

    Send a JSON message (e.g., { "type": "text", "content": "Hello, World!", "Id_user": "user1" })

9. CORS and Front-End Integration

If you're accessing the server from a different domain (e.g., a React or Vue.js front-end), ensure you handle CORS correctly by using the Access-Control-Allow-Origin header. The CORS configuration is already included in the server code to allow all origins (*), but you can customize it for production.
Environment Variables
Variable	Description	Example
DB_HOST	The hostname or IP address of the PostgreSQL server.
DB_PORT	The port number for PostgreSQL.	default is: 5432
DB_USER	The username to access PostgreSQL.	postgres
DB_PASSWORD	The password for the PostgreSQL user.	
DB_NAME	The name of the database to connect to.
SERVER_PORT	The port on which the Go server will run.	Ex: 8080
GET /getMessages

Fetches all messages stored in the database and returns them as a JSON array.

Response:

[
  {
    "id": 1,
    "type": "text",
    "content": "Hello, World!",
    "timestamp": "2025-06-28T12:00:00Z",
    "Id_user": "user1"
  }
]

WS /ws

WebSocket endpoint for sending and receiving real-time messages.

    Message format:

{
  "type": "text",
  "content": "Hello, World!",
  "Id_user": "user1"
}

    The server broadcasts messages to all connected WebSocket clients.

Troubleshooting

    CORS Issues: If you encounter issues with cross-origin requests, make sure your front-end is correctly configured to handle CORS. The server includes basic CORS support, allowing any origin. You can modify the CORS configuration in the Go server for more restrictive policies.

    Database Connectivity: Ensure that the PostgreSQL database is running and that the correct connection string is provided in the .env file. If the database is hosted remotely, replace localhost with the appropriate host address.

---

### Key Sections of the README:
- **Setup Instructions**: Describes all the necessary steps, from cloning the repository to running the application.
- **Environment Variables Table**: Provides a table listing all environment variables and their descriptions, along with example values.
- **API Documentation**: Explains the available API endpoints and their expected behavior, including WebSocket messaging and the HTTP endpoint to retrieve messages.
- **Troubleshooting**: Addresses common issues like CORS and database connectivity problems.

This should be a complete guide for setting up, running, and testing your Go WebSocket server. Let me know if you need additional adjustments!

