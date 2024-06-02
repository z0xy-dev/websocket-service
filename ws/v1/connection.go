package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Connection struct {
	ID     string
	Conn   *websocket.Conn
	Events []string
}

func (c *Connection) AddEvent(event string) {
	if c.HasEvent(event) {
		return
	}

	c.Events = append(c.Events, event)
}

func (c *Connection) RemoveEvent(event string) {
	for i, e := range c.Events {
		if e == event {
			c.Events = append(c.Events[:i], c.Events[i+1:]...)
		}
	}
}

func (c *Connection) HasEvent(event string) bool {
	for _, e := range c.Events {
		if e == event {
			return true
		}
	}
	return false
}

func (c *Connection) Send(data any) {
	message, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		log.Println("Failed to send message:", jsonErr)
		return
	}

	err := c.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

func (c *Connection) ReadMessage() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var action WebsocketActionDTO
		if err := json.Unmarshal(message, &action); err != nil {
			log.Println("unmarshal:", err)
			continue
		}

		handleWebsocketAction(c, action)
	}
}

func (c *Connection) SendToEvent(action string, event string, data any) {
	for i := 0; i < len(Connections); i++ {
		connection := Connections[i]
		if connection.HasEvent(event) {
			connection.Send(&WebsocketActionDTO{
				Action: action,
				Data: &EventDataDTO{
					ConnectionID: c.ID,
					Event:        event,
					Data:         data,
				},
			})
		}
	}
}

func handleWebsocketAction(c *Connection, action WebsocketActionDTO) {
	switch action.Action {
	case "id":
		c.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   c.ID,
		})
	case "ping":
		c.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "pong",
		})
	case "add event":
		event, ok := action.Data.(string)
		if !ok {
			errMsg := "Invalid data type for add event"
			log.Println(errMsg)
			c.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
		}

		c.AddEvent(event)
		c.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "added",
		})
	case "remove event":
		event, ok := action.Data.(string)
		if !ok {
			errMsg := "Invalid data type for remove event"
			log.Println(errMsg)
			c.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
		}

		c.RemoveEvent(event)
		c.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   "removed",
		})
	case "events":
		c.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   c.Events,
		})
	case "send to event":
		data := action.Data.(map[string]any)
		event, eventOK := data["event"].(string)

		if !eventOK {
			errMsg := "Invalid data type for send to event"
			log.Println(errMsg)
			c.Send(&WebsocketActionDTO{
				Action: action.Action,
				Data:   errMsg,
			})
		}

		c.SendToEvent(action.Action, event, data["data"])
	default:
		errMsg := fmt.Sprintf("Unknown action: %s", action.Action)
		c.Send(&WebsocketActionDTO{
			Action: action.Action,
			Data:   errMsg,
		})
		log.Printf(errMsg)
	}
}
