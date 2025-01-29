// Connection represents a single websocket connection
package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

// Connection holds the websocket connection and events of interest
type Connection struct {
	ID               string
	ApplicationName  string
	Conn             *websocket.Conn
	SubscribedEvents []string
}

// AddEvent subscribes the connection to a specific event
func (connection *Connection) AddEvent(eventName string) {
	eventNameWithApplicationName := fmt.Sprintf("%s:%s", connection.ApplicationName, eventName)

	if connection.HasEvent(eventNameWithApplicationName) {
		return
	}
	connection.SubscribedEvents = append(connection.SubscribedEvents, eventNameWithApplicationName)
}

// RemoveEvent unsubscribes the connection from a specific event
func (connection *Connection) RemoveEvent(eventName string) {
	eventNameWithApplicationName := fmt.Sprintf("%s:%s", connection.ApplicationName, eventName)

	for i, e := range connection.SubscribedEvents {
		if e == eventNameWithApplicationName {
			connection.SubscribedEvents = append(connection.SubscribedEvents[:i], connection.SubscribedEvents[i+1:]...)
			return
		}
	}
}

// HasEvent checks if the connection is already subscribed to an event
func (connection *Connection) HasEvent(eventName string) bool {
	eventNameWithApplicationName := fmt.Sprintf("%s:%s", connection.ApplicationName, eventName)

	for _, e := range connection.SubscribedEvents {
		if e == eventNameWithApplicationName {
			return true
		}
	}
	return false
}

// Send writes JSON-encoded data to the websocket
func (connection *Connection) Send(payload any) {
	message, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		log.Println("Failed to encode message:", marshalErr)
		return
	}

	writeErr := connection.Conn.WriteMessage(websocket.TextMessage, message)
	if writeErr != nil {
		log.Println("Failed to send message:", writeErr)
	}
}

// ReadMessage listens for incoming messages and processes them
func (connection *Connection) ReadMessage() {
	for {
		_, rawMessage, readErr := connection.Conn.ReadMessage()
		if readErr != nil {
			log.Println("Read error:", readErr)
			return
		}

		var action WebsocketActionDTO
		if unmarshalErr := json.Unmarshal(rawMessage, &action); unmarshalErr != nil {
			log.Println("Decode error:", unmarshalErr)
			continue
		}

		handleWebsocketAction(connection, action)
	}
}

// SendToEvent sends data to all connections subscribed to the specified event
func (connection *Connection) SendToEvent(actionType string, eventName string, payload any) {
	eventNameWithApplicationName := fmt.Sprintf("%s:%s", connection.ApplicationName, eventName)

	for i := 0; i < len(activeConnections); i++ {
		otherConnection := activeConnections[i]
		if otherConnection.HasEvent(eventNameWithApplicationName) {
			otherConnection.Send(&WebsocketActionDTO{
				Action: actionType,
				Data: &EventDataDTO{
					ConnectionID: connection.ID,
					Event:        eventName,
					Data:         payload,
				},
			})
		}
	}
}

// handleWebsocketAction routes the incoming action to the appropriate handler
func handleWebsocketAction(connection *Connection, action WebsocketActionDTO) {
	switch strings.ToLower(action.Action) {
	case "id":
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   connection.ID,
		})
	case "ping":
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "pong",
		})
	case "application":
		applicationName, ok := action.Data.(string)

		if !ok {
			errMsg := "Invalid data type for application"
			log.Println(errMsg)
			connection.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
			return
		}

		connection.ApplicationName = applicationName
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "application name set to " + applicationName,
		})
	case "add event":
		eventName, ok := action.Data.(string)
		if !ok {
			errMsg := "Invalid data type for add event"
			log.Println(errMsg)
			connection.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
			return
		}
		connection.AddEvent(eventName)
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "added",
		})
	case "remove event":
		eventName, ok := action.Data.(string)
		if !ok {
			errMsg := "Invalid data type for remove event"
			log.Println(errMsg)
			connection.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
			return
		}
		connection.RemoveEvent(eventName)
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "removed",
		})
	case "events":
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   connection.SubscribedEvents,
		})
	case "send to event":
		dataMap := action.Data.(map[string]any)
		eventName, eventOk := dataMap["event"].(string)
		if !eventOk {
			errMsg := "Invalid data type for send to event"
			log.Println(errMsg)
			connection.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
			return
		}
		connection.SendToEvent(action.Action, eventName, dataMap["data"])
	default:
		errMsg := fmt.Sprintf("Unknown action: %s", action.Action)
		connection.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   errMsg,
		})
		log.Println(errMsg)
	}
}
