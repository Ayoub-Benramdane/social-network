package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/database"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var receiver_id int64
	json.NewDecoder(r.Body).Decode(&receiver_id)
	_, err = database.CheckUser(receiver_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve recipient"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	chats, err := database.GetConversation(user.ID, receiver_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve chats"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func ChatGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	chats, err := database.GetGroupConversation(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve chats"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
