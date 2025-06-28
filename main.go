package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Configuração do banco de dados
var db *gorm.DB
var err error

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, change this for production environments
	},
}

func loadEnv() {
    // Load the .env file into the environment
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

// Estrutura de Mensagem
type Message struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Id_user   string `json:"Id_user"`
}

// Conectar ao banco de dados e configurar o GORM
func initDB() {
	// PostgreSQL connection string

	 // Access environment variables using os.Getenv
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", dbHost, dbUser, dbName, dbPassword, dbPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Auto migrate to create the table
	err = db.AutoMigrate(&Message{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}

	log.Println("✅ Connected to PostgreSQL database")
}

// Connected clients
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	// Add new client
	clients[ws] = true
	log.Println("🔌 New client connected")

	for {
		var msg Message
		// Read JSON message from client
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ Read error: %v", err)
			delete(clients, ws) // Remove the client on error
			break
		}

		// Add timestamp to message (formatted as a string)
		msg.Timestamp = time.Now().Format(time.RFC3339)

		// Log and send to broadcast channel
		log.Printf("📨 Received: %+v\n", msg)

		// Save the message to the database
		saveMessage(msg)

		// Send the message to the broadcast channel
		broadcast <- msg
	}
}

// Save the message to the database
func saveMessage(msg Message) {
	if err := db.Create(&msg).Error; err != nil {
		log.Printf("❌ Error saving message: %v", err)
	} else {
		log.Printf("✅ Message saved to database: %+v", msg)
	}
}

// Broadcast messages to all connected clients
func handleMessages() {
	for {
		msg := <-broadcast

		// Send to all clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("⚠️ Write error: %v", err)
				client.Close()
				delete(clients, client) // Remove client from the map if an error occurs
			}
		}
	}
}

// Get all messages from the database and return as JSON
func getMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message

	// Fetch all messages from the database
	err := db.Find(&messages).Error
	if err != nil {
		log.Printf("❌ Error retrieving messages: %v", err)
		http.Error(w, "Unable to retrieve messages", http.StatusInternalServerError)
		return
	}

	// Set response headers for CORS
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Explicitly set this header for CORS
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allow CORS methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allow specific headers
	w.WriteHeader(http.StatusOK)

	// Return messages as JSON
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("❌ Error encoding messages to JSON: %v", err)
		http.Error(w, "Unable to encode messages", http.StatusInternalServerError)
	}
}

func main() {
	// Initialize the database
	loadEnv()
	initDB()

	// Set up CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins (use specific origins in production)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Set this to true if you need credentials (cookies, auth headers)
	}).Handler

	// Set up HTTP handlers
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/getMessages", getMessages)

	// Apply CORS middleware to all routes
	http.Handle("/", corsHandler(http.DefaultServeMux))

	// Start the goroutine for message broadcasting
	go handleMessages()

	// Start the server
	fmt.Println("✅ Server listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
