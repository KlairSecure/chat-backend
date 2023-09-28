package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"server/handlers"
)

func main() {
	r := mux.NewRouter()

	// Handle WebSocket connections
	r.HandleFunc("/ws/{roomID}/{userID}", handlers.HandleWebSocket)

	// Handle room creation
	r.HandleFunc("/rooms", handlers.CreateRoom).Methods("POST")

	// Handle getting active room IDs with names
	r.HandleFunc("/rooms", handlers.GetActiveRooms).Methods("GET")

	http.Handle("/", r)

	// Start the server
	fmt.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
