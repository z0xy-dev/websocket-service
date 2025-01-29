package server

import (
	"fmt"
	"net/http"
	v1 "websocketservice/ws/v1"
)

func Start() {
	http.HandleFunc("/ws/v1", v1.Handler)
	fmt.Println("Starting server on :3399")
	if err := http.ListenAndServe(":3399", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
