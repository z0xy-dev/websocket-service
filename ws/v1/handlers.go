package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// activeConnections holds all current connected clients
var activeConnections []*Connection
var connectionsMutex sync.Mutex

// Handler upgrades the HTTP connection to a WebSocket connection and manages it
func Handler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new WebSocket connection object
	connection := &Connection{
		ID:     uuid.New().String(),
		Conn:   wsConn,
		Events: []string{},
	}

	// Add the new WebSocket connection to the global activeConnections slice
	connectionsMutex.Lock()
	activeConnections = append(activeConnections, connection)
	connectionsMutex.Unlock()

	// Start listening for incoming messages from the connection
	go connection.ReadMessage()

	// Remove any closed connections
	removeClosedConnections()
}

// removeClosedConnections removes connections that are no longer active
func removeClosedConnections() {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	for i := len(activeConnections) - 1; i >= 0; i-- {
		connection := activeConnections[i]
		// Send a ping message to check if the connection is still active
		err := connection.Conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			// If the write fails, remove the connection from the slice
			activeConnections = append(activeConnections[:i], activeConnections[i+1:]...)
		}
	}
}
