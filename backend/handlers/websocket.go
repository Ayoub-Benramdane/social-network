package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
	"strings"
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
		members, err := database.GetGroupMembers(group.ID)
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

	conn.WriteJSON(connections)

	for {
		var message structs.Message
		message.Type = r.FormValue("type")
		message.ReceiverID, err = strconv.ParseInt(r.FormValue("receiver_id"), 10, 64)
		if err != nil {
			fmt.Println(err)
			response := map[string]string{"error": "Invalid receiver ID"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		message.Content = r.FormValue("content")

		var imagePath string
		image, header, err := r.FormFile("chat_image")
		if err != nil && err.Error() != "http: no such file" {
			fmt.Println(err)
			response := map[string]string{"error": "Failed to retrieve image"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if image != nil {
			imagePath, err = SaveImage(image, header, "../frontend/public/chat/")
			if err != nil {
				response := map[string]string{"error": err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			newpath := strings.Split(imagePath, "/public")
			imagePath = newpath[1]
		}

		if message.Type == "message" {
			if (message.Content == "" || message.Image == "") && (message.ReceiverID == 0 || message.GroupID == 0) {
				response := map[string]string{"error": "Message content and image cannot be empty"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}
			if message.ReceiverID != 0 {
				if _, err := database.GetUserById(message.ReceiverID); err != nil {
					respose := map[string]string{"error": "Failed to retrieve user"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(respose)
					return
				}
				if database.SendMessage(user.ID, message.ReceiverID, 0, message.Content, imagePath) != nil {
					response := map[string]string{"error": "Failed to send message"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				SendWsMessage(message.ReceiverID, map[string]interface{}{"type": "message", "id": user.ID, "username": user.Username, "content": message.Content, "image": imagePath})

			} else if message.GroupID != 0 {
				if _, err := database.GetGroupById(message.GroupID); err != nil {
					response := map[string]string{"error": "Failed to retrieve group"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				if database.SendMessage(user.ID, 0, message.GroupID, message.Content, imagePath) != nil {
					response := map[string]string{"error": "Failed to send message"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				SendWsMessage(message.GroupID, map[string]interface{}{"type": "message", "id": user.ID, "username": user.Username, "content": message.Content, "image": imagePath})
			}
		} else if message.Type == "typing" {
			if message.ReceiverID != 0 {
				SendWsMessage(message.ReceiverID, map[string]interface{}{"type": "typing", "id": user.ID, "username": user.Username})
			} else if message.GroupID != 0 {
				SendWsMessage(message.GroupID, map[string]interface{}{"type": "typing", "id": user.ID, "username": user.Username})
			}
		} else if message.Type == "notification" {
			if message.ReceiverID != 0 {
				SendWsMessage(message.ReceiverID, map[string]interface{}{"type": "notification", "id": user.ID, "username": user.Username, "content": message.Content})
			} else if message.GroupID != 0 {
				SendWsMessage(message.GroupID, map[string]interface{}{"type": "notification", "id": user.ID, "username": user.Username, "content": message.Content})
			}
		}
	}

	RemoveClient(conn, user.ID)
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

func RemoveClient(conn *websocket.Conn, user_id int64) {
	Mutex.Lock()
	defer Mutex.Unlock()
	if clients, ok := Clients[user_id]; ok {
		for i, client := range clients {
			if client == conn {
				Clients[user_id] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}
}
