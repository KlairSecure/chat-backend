package handlers

import (
	"encoding/json"
	"net/http"
	"server/models"

	"server/utils"
)

// Create a new room
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Parse the request JSON to get the room name
	var requestData struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique room ID
	roomID := utils.GenerateUniqueID()

	// Create a new room and add it to the rooms map
	newRoom := &models.Room{
		ID:    roomID,
		Name:  requestData.Name,
		Users: make(map[string]bool),
	}

	// Register new room
	rooms[roomID] = newRoom

	// Return the room ID and name in the response
	responseData := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		ID:   newRoom.ID,
		Name: newRoom.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get active rooms (already created)
func GetActiveRooms(w http.ResponseWriter, r *http.Request) {
	// Create a slice to store active room IDs and names
	type responseData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	activeRooms := make([]responseData, 0, len(rooms))

	// Iterate over the rooms map and collect active rooms
	for _, room := range rooms {
		activeRooms = append(activeRooms, responseData{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	// Return a JSON response containing active room IDs and names
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(activeRooms); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
