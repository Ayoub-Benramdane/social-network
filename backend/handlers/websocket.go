package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/database"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients  = make(map[int64][]*websocket.Conn)
	Mutex    sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response := map[string]string{"error": "Failed to upgrade connection"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
	defer conn.Close()

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	Mutex.Lock()
	Clients[user.ID] = append(Clients[user.ID], conn)

	connections, err := database.GetConnections(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve connections"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	groups, err := database.GetGroups(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	for i, connection := range connections {
		if _, exist := Clients[connection.ID]; exist {
			connections[i].Online = true
		}
	}

	for _, group := range groups {
		members, err := database.GetGroupMembers(group.ID, user.ID)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve group member"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		for i, member := range members {
			if _, exist := Clients[member.ID]; exist {
				members[i].Online = true
			}
		}
	}

	SendWsMessage(user.ID, map[string]interface{}{"type": "online", "id": user.ID, "username": user.Username})
	Mutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			Mutex.Lock()
			for i, client := range Clients[user.ID] {
				if client == conn {
					Clients[user.ID] = append(Clients[user.ID][:i], Clients[user.ID][i+1:]...)
					break
				}
			}
			Mutex.Unlock()
			break
		}
	}
}

func SendWsMessage(user_id int64, message map[string]interface{}) {
	Mutex.Lock()
	defer Mutex.Unlock()
	if clients, ok := Clients[user_id]; ok {
		for _, client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(Clients, user_id)
			}
		}
	}
}
