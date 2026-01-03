package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"social-network/database"
)

func GetConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed")
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user")
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userConnections, err := database.FetchUserConnections(currentUser.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve connections")
		resp := map[string]string{"error": "Failed to retrieve connections"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userConnections)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	receiverID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid receiver ID", err)
		resp := map[string]string{"error": "Invalid receiver ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	messageOffset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println("Invalid offset", err)
		resp := map[string]string{"error": "Invalid offset"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = database.UserExists(receiverID)
	if err != nil {
		fmt.Println("Failed to retrieve recipient", err)
		resp := map[string]string{"error": "Failed to retrieve recipient"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ClientsMutex.Lock()
	err = database.MarkMessagesAsRead(receiverID, currentUser.UserID, 0)
	if err != nil {
		fmt.Println("Failed to mark messages as read", err)
		resp := map[string]string{"error": "Failed to mark messages as read"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	SendWsMessage(currentUser.UserID, map[string]interface{}{"type": "read_messages"})
	ClientsMutex.Unlock()

	privateChats, err := database.FetchConversation(currentUser.UserID, receiverID, messageOffset)
	if err != nil {
		fmt.Println("Failed to retrieve chats", err)
		resp := map[string]string{"error": "Failed to retrieve chats"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	slices.Reverse(privateChats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(privateChats)
}

func ChatGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	groupID, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid group ID", err)
		resp := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	messageOffset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println("Invalid offset", err)
		resp := map[string]string{"error": "Invalid offset"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = database.GetGroupByID(groupID)
	if err != nil {
		fmt.Println("Failed to retrieve group", err)
		resp := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isMember, err := database.IsUserGroupMember(currentUser.UserID, groupID)
	if err != nil {
		fmt.Println("Failed to check if user is a member", err)
		resp := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	} else if !isMember {
		fmt.Println("User is not a member of the group")
		resp := map[string]string{"error": "User is not a member of the group"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = database.MarkMessagesAsRead(0, currentUser.UserID, groupID)
	if err != nil {
		fmt.Println("Failed to mark messages as read", err)
		resp := map[string]string{"error": "Failed to mark messages as read"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ClientsMutex.Lock()
	SendWsMessage(currentUser.UserID, map[string]interface{}{"type": "read_messages"})
	ClientsMutex.Unlock()

	groupChats, err := database.FetchGroupConversation(groupID, currentUser.UserID, messageOffset)
	if err != nil {
		fmt.Println("Failed to retrieve chats", err)
		resp := map[string]string{"error": "Failed to retrieve chats"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	slices.Reverse(groupChats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupChats)
}

func ReadMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "messages") {
		return
	}

	var payload struct {
		TargetUserID  int64 `json:"user_id"`
		TargetGroupID int64 `json:"group_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Invalid request body", err)
		resp := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if payload.TargetUserID == 0 && payload.TargetGroupID == 0 {
		fmt.Println("Invalid message ID")
		resp := map[string]string{"error": "Invalid message ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ClientsMutex.Lock()
	if err = database.MarkMessagesAsRead(payload.TargetUserID, currentUser.UserID, payload.TargetGroupID); err != nil {
		fmt.Println("Failed to mark messages as read", err)
		resp := map[string]string{"error": "Failed to mark messages as read"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	SendWsMessage(currentUser.UserID, map[string]interface{}{"type": "read_messages"})
	ClientsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("success")
}
