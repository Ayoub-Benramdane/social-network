package handlers

import (
	"encoding/json"
	"net/http"
	structs "social-network/data"
	"social-network/database"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var comment struct {
		PostID  int64  `json:"postId"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
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

	id, err := database.CreateComment(comment.Content, user.ID, comment.PostID)
	if err != nil {
		response := map[string]string{"error": "Failed to create comment"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newComment := structs.Comment{
		ID:        id,
		Content:   comment.Content,
		Author:    user.Username,
		CreatedAt: "Just Now",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
