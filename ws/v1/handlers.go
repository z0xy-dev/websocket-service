package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Connections holds all current connected clients
var Connections []*Connection
var connectionsMutex sync.Mutex

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Add the new WebSocket connection to the global clients slice
	connection := &Connection{
		ID:   uuid.New().String(),
		Conn: conn,
	}
	connectionsMutex.Lock()
	Connections = append(Connections, connection)
	connectionsMutex.Unlock()

	// Start listening for incoming messages from the connection
	go connection.ReadMessage()

	RemoveClosedConnections()
}

func RemoveClosedConnections() {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	for i := len(Connections) - 1; i >= 0; i-- {
		connection := Connections[i]
		err := connection.Conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			// If the write fails, remove the connection from the slice
			Connections = append(Connections[:i], Connections[i+1:]...)
		}
	}
}
