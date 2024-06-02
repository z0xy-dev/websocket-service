package main

import (
	"WebsocketService/ws/v1"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws/v1", v1.Handler)
	fmt.Println("Starting server on :3399")
	if err := http.ListenAndServe(":3399", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
