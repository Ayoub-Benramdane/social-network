package handlers

import (
	"encoding/json"
	"html"
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

	var comment structs.Comment
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

	if comment.Content == "" || len(comment.Content) > 100 {
		response := map[string]string{"error": "Comment content is required and must be less than 100 characters"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(user.ID, comment.PostID, 0)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := database.CreateComment(comment.Content, user.ID, post)
	if err != nil {
		response := map[string]string{"error": "Failed to create comment"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newComment := structs.Comment{
		ID:        id,
		PostID:    comment.PostID,
		Content:   html.EscapeString(comment.Content),
		Author:    user.Username,
		CreatedAt: "Just Now",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}

func CreateGroupCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var comment structs.Comment
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

	_, err = database.GetGroup(comment.GroupID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	member, err := database.IsMemberGroup(user.ID, comment.GroupID)
	if err != nil {
		response := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if !member {
		response := map[string]string{"error": "User is not a member of the group"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if comment.Content == "" || len(comment.Content) > 100 {
		response := map[string]string{"error": "Comment content is required and must be less than 100 characters"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(user.ID, comment.PostID, comment.GroupID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := database.CreateGroupComment(comment.Content, user.ID, comment.GroupID, post)
	if err != nil {
		response := map[string]string{"error": "Failed to create comment"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newComment := structs.Comment{
		ID:        id,
		PostID:    comment.PostID,
		Content:   html.EscapeString(comment.Content),
		Author:    user.Username,
		CreatedAt: "Just Now",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
