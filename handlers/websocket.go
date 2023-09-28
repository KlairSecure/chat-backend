package handlers

import (
	"fmt"
	"log"
	"net/http"
	"server/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var rooms = make(map[string]*models.Room)
var connections = make(map[string]*websocket.Conn)

// Handle WebSocket Messages
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	roomID := vars["roomID"]
	userID := vars["userID"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add the user's WebSocket connection to the connections map
	connections[userID] = conn

	// Check if the room exists
	room, exists := rooms[roomID]
	if !exists {
		log.Printf("Room %s does not exist\n", roomID)
		return
	}

	// Register user to the room
	room.Users[userID] = true

	// Handle WebSocket messages
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// Remove the user from the room and connections map when they disconnect
			delete(room.Users, userID)
			delete(connections, userID)
			return
		}

		// Handle chat message
		if messageType == websocket.TextMessage {
			message := string(p)

			// Include the user's ID in the message
			messageWithSender := fmt.Sprintf("[%s]: %s", userID, message)

			// Broadcast the message to all users in the room
			for recipientUserID := range room.Users {
				recipientConn, found := connections[recipientUserID]
				if !found {
					// Handle case where recipient is not connected
					log.Printf("User %s is not connected\n", recipientUserID)
					continue
				}
				if err := recipientConn.WriteMessage(websocket.TextMessage, []byte(messageWithSender)); err != nil {
					// Handle write error (e.g., user disconnected)
					log.Printf("Error sending message to user %s: %v\n", recipientUserID, err)
				}
			}
		}
	}
}
