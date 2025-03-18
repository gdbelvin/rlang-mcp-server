package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// HTTPTransport implements the types.Transport interface for HTTP with SSE
type HTTPTransport struct {
	addr           string
	clients        map[string]chan json.RawMessage
	clientsMutex   sync.Mutex
	nextClientID   int
	messageHandler func(json.RawMessage)
	closeHandler   func()
	errorHandler   func(error)
	server         *http.Server
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(addr string) *HTTPTransport {
	return &HTTPTransport{
		addr:         addr,
		clients:      make(map[string]chan json.RawMessage),
		nextClientID: 1,
	}
}

// Start starts the HTTP server
func (t *HTTPTransport) Start() error {
	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp/connect", t.handleConnect)
	mux.HandleFunc("/mcp/request", t.handleRequest)

	// Create HTTP server
	t.server = &http.Server{
		Addr:    t.addr,
		Handler: mux,
	}

	// Start the HTTP server in a goroutine
	go func() {
		fmt.Printf("Starting HTTP server on %s\n", t.addr)
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if t.errorHandler != nil {
				t.errorHandler(fmt.Errorf("HTTP server error: %w", err))
			} else {
				log.Printf("HTTP server error: %v", err)
			}
		}
	}()

	return nil
}

// Send sends a message to all connected clients
func (t *HTTPTransport) Send(message json.RawMessage) error {
	t.clientsMutex.Lock()
	defer t.clientsMutex.Unlock()

	// If no clients are connected, log a warning
	if len(t.clients) == 0 {
		return fmt.Errorf("no clients connected to receive message")
	}

	// Send the message to all connected clients
	for _, messageChan := range t.clients {
		select {
		case messageChan <- message:
			// Message sent successfully
		default:
			// Channel is full, log a warning
			if t.errorHandler != nil {
				t.errorHandler(fmt.Errorf("client message channel is full, dropping message"))
			} else {
				log.Printf("Client message channel is full, dropping message")
			}
		}
	}

	return nil
}

// OnMessage sets the callback for when a message is received
func (t *HTTPTransport) OnMessage(callback func(json.RawMessage)) {
	t.messageHandler = callback
}

// Close closes the HTTP server
func (t *HTTPTransport) Close() {
	if t.server != nil {
		if err := t.server.Close(); err != nil {
			if t.errorHandler != nil {
				t.errorHandler(fmt.Errorf("error closing HTTP server: %w", err))
			} else {
				log.Printf("Error closing HTTP server: %v", err)
			}
		}
	}

	// Close all client channels
	t.clientsMutex.Lock()
	for _, messageChan := range t.clients {
		close(messageChan)
	}
	t.clients = make(map[string]chan json.RawMessage)
	t.clientsMutex.Unlock()

	// Call the close handler
	if t.closeHandler != nil {
		t.closeHandler()
	}
}

// OnClose sets the callback for when the connection is closed
func (t *HTTPTransport) OnClose(callback func()) {
	t.closeHandler = callback
}

// OnError sets the callback for when an error occurs
func (t *HTTPTransport) OnError(callback func(error)) {
	t.errorHandler = callback
}

// handleConnect handles SSE connections
func (t *HTTPTransport) handleConnect(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a client ID and channel
	t.clientsMutex.Lock()
	clientID := fmt.Sprintf("client-%d", t.nextClientID)
	t.nextClientID++
	messageChan := make(chan json.RawMessage, 10)
	t.clients[clientID] = messageChan
	t.clientsMutex.Unlock()

	// Send a connected message
	fmt.Fprintf(w, "event: connected\ndata: {\"clientID\":\"%s\"}\n\n", clientID)
	w.(http.Flusher).Flush()

	// Clean up when the client disconnects
	defer func() {
		t.clientsMutex.Lock()
		delete(t.clients, clientID)
		close(messageChan)
		t.clientsMutex.Unlock()
		log.Printf("Client %s disconnected", clientID)
	}()

	// Keep the connection open and send messages as they arrive
	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				return
			}
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(message))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			log.Printf("Client %s connection closed", clientID)
			return
		}
	}
}

// handleRequest handles MCP requests
func (t *HTTPTransport) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the client ID from the request
	clientID := r.Header.Get("X-MCP-Client-ID")
	if clientID == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	// Check if the client exists
	t.clientsMutex.Lock()
	_, exists := t.clients[clientID]
	t.clientsMutex.Unlock()
	if !exists {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Call the message handler
	if t.messageHandler != nil {
		t.messageHandler(json.RawMessage(body))
	} else {
		if t.errorHandler != nil {
			t.errorHandler(fmt.Errorf("received message but no handler is set"))
		} else {
			log.Printf("Received message but no handler is set")
		}
	}

	// Send an acknowledgment response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
	})
}
