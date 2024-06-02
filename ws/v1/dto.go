package v1

type WebsocketActionDTO struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type ValueDTO[T any] struct {
	Value T `json:"value"`
}

type EventDataDTO struct {
	ConnectionID string `json:"connection_id"`
	Event        string `json:"event"`
	Data         any    `json:"data"`
}

type SendToEventDataDTO struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}
