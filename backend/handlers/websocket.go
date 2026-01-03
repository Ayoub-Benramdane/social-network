package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	structs "social-network/data"
	"social-network/database"

	"github.com/gorilla/websocket"
)

var (
	ConnectedClients = structs.ConnectedClients
	ClientsMutex     sync.Mutex
	wsUpgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		resp := map[string]string{"error": "Failed to upgrade connection"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println(err)
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ClientsMutex.Lock()
	ConnectedClients[currentUser.UserID] = append(ConnectedClients[currentUser.UserID], connection)
	ClientsMutex.Unlock()

	NotifyUsers(currentUser.UserID, "online")
	ListenForMessages(connection, currentUser.UserID, w, r)
}

func ListenForMessages(connection *websocket.Conn, currentUserID int64, w http.ResponseWriter, r *http.Request) {
	defer func() {
		RemoveClient(connection, currentUserID)
		NotifyUsers(currentUserID, "offline")
		connection.Close()
	}()

	currentUser, err := database.GetProfileInfo(currentUserID, nil)
	if err != nil {
		fmt.Println(err)
	}

	for {
		var incomingMessage structs.Message
		err = connection.ReadJSON(&incomingMessage)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}

		if currentUserID == incomingMessage.SenderID {
			continue
		}

		messageType := ""
		if incomingMessage.MessageType == "message" {
			messageType = "message"
		} else if incomingMessage.MessageType == "typing" {
			messageType = "typing"
		}

		if !CheckLastActionTime(w, r, "messages") {
			continue
		}

		var targetUserIDs []int64

		if incomingMessage.GroupID != 0 {
			groupData, err := database.GetGroupByID(incomingMessage.GroupID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if isMember, err := database.IsUserGroupMember(currentUserID, incomingMessage.GroupID); err != nil || !isMember {
				fmt.Println(err)
				continue
			}

			targetUserIDs, err = database.GetGroupMemberIDs(incomingMessage.GroupID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, memberID := range targetUserIDs {
				if _, err := database.FindUserByID(memberID); err != nil {
					fmt.Println("Error getting user by ID:", err)
					return
				}

				isMember, err := database.IsUserGroupMember(memberID, incomingMessage.GroupID)
				if err != nil {
					fmt.Println("Error checking group membership:", err)
					continue
				}
				if !isMember {
					continue
				}

				ClientsMutex.Lock()
				err = database.CreateMessage(currentUserID, memberID, incomingMessage.GroupID, incomingMessage.Content, "")
				if err != nil && currentUserID != memberID {
					fmt.Println("Error sending message:", err)
					ClientsMutex.Unlock()
					continue
				}

				SendWsMessage(memberID, map[string]interface{}{
					"type":         messageType,
					"message_id":   time.Now(),
					"name":         groupData.Name,
					"user_id":      currentUser.UserID,
					"group_id":     groupData.GroupID,
					"username":     currentUser.Username,
					"avatar":       currentUser.AvatarURL,
					"content":      incomingMessage.Content,
					"current_user": memberID,
					"created_at":   "Just now",
				})
				ClientsMutex.Unlock()
			}
		} else {
			targetUser, err := database.FindUserByID(incomingMessage.SenderID)
			if err != nil {
				fmt.Println("Error getting user by ID:", err)
				return
			}

			isFollowed, err := database.IsUserFollowing(currentUserID, incomingMessage.SenderID)
			if err != nil || (!isFollowed && targetUser.PrivacyLevel == "private") {
				fmt.Println("User is not followed or privacy restricted:", err)
				continue
			}

			ClientsMutex.Lock()
			err = database.CreateMessage(currentUserID, incomingMessage.SenderID, 0, incomingMessage.Content, "")
			if err != nil {
				fmt.Println("Error sending message:", err)
				ClientsMutex.Unlock()
				continue
			}

			SendWsMessage(currentUserID, map[string]interface{}{
				"type":         messageType,
				"message_id":   time.Now(),
				"user_id":      currentUser.UserID,
				"username":     currentUser.Username,
				"avatar":       currentUser.AvatarURL,
				"content":      incomingMessage.Content,
				"current_user": currentUserID,
				"created_at":   "Just now",
			})

			SendWsMessage(incomingMessage.SenderID, map[string]interface{}{
				"type":         messageType,
				"message_id":   time.Now(),
				"user_id":      currentUser.UserID,
				"username":     currentUser.Username,
				"avatar":       currentUser.AvatarURL,
				"content":      incomingMessage.Content,
				"current_user": incomingMessage.SenderID,
				"created_at":   "Just now",
			})
			ClientsMutex.Unlock()
		}
	}
}

func NotifyUsers(userID int64, status string) {
	connections, err := database.FetchUserConnections(userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, connection := range connections {
		if connection.UserID == userID {
			continue
		}

		ClientsMutex.Lock()
		if _, ok := ConnectedClients[connection.UserID]; ok {
			if status == "online" {
				SendWsMessage(connection.UserID, map[string]interface{}{
					"type":    "new_connection",
					"user_id": userID,
					"online":  true,
				})
			} else if status == "offline" {
				SendWsMessage(connection.UserID, map[string]interface{}{
					"type":    "disconnection",
					"user_id": userID,
					"online":  false,
				})
			}
		}
		ClientsMutex.Unlock()
	}
}

func SendWsMessage(userID int64, payload map[string]interface{}) {
	if userConnections, ok := ConnectedClients[userID]; ok {
		for _, clientConn := range userConnections {
			if err := clientConn.WriteJSON(payload); err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}

func RemoveClient(connection *websocket.Conn, userID int64) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	if userConnections, ok := ConnectedClients[userID]; ok {
		for index, clientConn := range userConnections {
			if clientConn == connection {
				ConnectedClients[userID] = append(userConnections[:index], userConnections[index+1:]...)
				break
			}
		}

		if len(ConnectedClients[userID]) == 0 {
			delete(ConnectedClients, userID)
		}
	}
}
