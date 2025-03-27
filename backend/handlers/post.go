package handlers

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	structs "social-network/backend/data"
	"social-network/backend/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var post struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"category"`
	}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := database.CreatePost(user.ID, post.Title, post.Content, post.Categories)
	if err != nil {
		log.Printf("Database error: %v", err)
		response := map[string]string{"error": "Failed to create post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newPost := structs.Post{
		ID:            id,
		Author:        user.Username,
		Title:         html.EscapeString(post.Title),
		Content:       html.EscapeString(post.Content),
		Categories:    post.Categories,
		CreatedAt:     "Just Now",
		TotalLikes:    0,
		TotalComments: 0,
		Comments:      nil,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts, err := database.GetPosts(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}