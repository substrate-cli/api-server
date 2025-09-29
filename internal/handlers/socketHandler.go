package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// More permissive origin checking for CSR and extensions
		origin := r.Header.Get("Origin")
		log.Printf("WebSocket connection from origin: %s", origin)

		// Allow connections from extensions (chrome-extension://)
		// Allow local development (localhost, 127.0.0.1)
		// Allow null origin (for file:// protocol or direct connections)
		if origin == "" || origin == "null" {
			return true
		}

		// You can add specific allowed origins here
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173", // Vite default
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}

		// Allow chrome extensions
		if strings.HasPrefix(origin, "chrome-extension://") {
			return true
		}

		// Allow moz extensions
		if strings.HasPrefix(origin, "moz-extension://") {
			return true
		}

		log.Printf("Origin not allowed: %s", origin)
		return true // Change to false in production for security
	},
	HandshakeTimeout: 30 * time.Second, // Increased timeout
	ReadBufferSize:   4096,
	WriteBufferSize:  4096,
	// Enable compression for better performance
	EnableCompression: true,
}

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func broadcastMessage(message []byte) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Broadcasting to %d clients", len(clients))

	for client := range clients {
		// Set write deadline to prevent blocking
		client.SetWriteDeadline(time.Now().Add(10 * time.Second))
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	// Log all headers for debugging
	log.Printf("WebSocket connection attempt from: %s", r.RemoteAddr)
	log.Printf("User-Agent: %s", r.Header.Get("User-Agent"))
	log.Printf("Origin: %s", r.Header.Get("Origin"))
	log.Printf("Sec-WebSocket-Protocol: %s", r.Header.Get("Sec-WebSocket-Protocol"))

	// Set CORS headers before upgrade (important for CSR)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	// Handle preflight requests (important for CORS)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}

	log.Println("WebSocket connected successfully")

	// Add to clients map
	mu.Lock()
	clients[conn] = true
	clientCount := len(clients)
	mu.Unlock()

	log.Printf("Client added. Total clients: %d", clientCount)

	// Configure connection settings for different environments
	// Longer timeouts for extensions and CSR apps
	readTimeout := 120 * time.Second
	writeTimeout := 30 * time.Second
	pingInterval := 45 * time.Second

	conn.SetReadDeadline(time.Now().Add(readTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	// Send a welcome message
	welcomeMsg := []byte(`{"type":"connection","message":"Connected to WebSocket server"}`)
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := conn.WriteMessage(websocket.TextMessage, welcomeMsg); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}

	// Start ping ticker for keep-alive
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	// Channel to signal when to stop the ping routine
	done := make(chan struct{})

	// Handle messages and connection lifecycle
	defer func() {
		close(done)

		mu.Lock()
		delete(clients, conn)
		clientCount := len(clients)
		mu.Unlock()

		conn.Close()
		log.Printf("WebSocket disconnected and cleaned up. Remaining clients: %d", clientCount)
	}()

	// Goroutine for handling pings
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Ping routine panic recovered: %v", r)
			}
		}()

		for {
			select {
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Ping failed: %v", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Main message reading loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Check for different types of close errors
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNoStatusReceived) {
				log.Printf("WebSocket closed: %v", err)
			} else if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNoStatusReceived) {
				log.Printf("WebSocket unexpected close: %v", err)
			} else {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		log.Printf("Received message type %d: %s", messageType, string(message))

		// Handle different message types
		switch messageType {
		case websocket.TextMessage:
			// Echo the message back or handle it as needed
			response := []byte(`{"type":"echo","data":"` + string(message) + `"}`)
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
				log.Printf("Echo failed: %v", err)
				return
			}
		case websocket.BinaryMessage:
			// Handle binary messages if needed
			log.Printf("Received binary message of length: %d", len(message))
		}

		// Reset read deadline after successful message
		conn.SetReadDeadline(time.Now().Add(readTimeout))
	}
}
