package v1

// WebsocketActionDTO holds details about an action to be processed.
type WebsocketActionDTO struct {
	// Action specifies the type of action.
	Action string `json:"action"`
	// Data contains any additional information for the action.
	Data any `json:"data"`
}

// EventDataDTO holds data associated with a specific event.
type EventDataDTO struct {
	// ConnectionID identifies the connection.
	ConnectionID string `json:"connection_id"`
	// Event specifies the event type.
	Event string `json:"event"`
	// Data contains any additional information for the event.
	Data any `json:"data"`
}
