package handlers

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	structs "social-network/data"
	"social-network/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var post struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		Image     string `json:"image"`
		Category string `json:"category"`
		Privacy   string `json:"privacy"`
	}

	err := json.NewDecoder(r.Body).Decode(&post)
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

	id, err := database.CreatePost(user.ID, post.Title, post.Content, post.Category, post.Image, post.Privacy)
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
		Image:         post.Image,
		CreatedAt:     "Just Now",
		Privacy:       post.Privacy,
		TotalLikes:    0,
		TotalComments: 0,
		Comments:      nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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

	var post_id int64
	err = json.NewDecoder(r.Body).Decode(&post_id)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(post_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if post.Privacy == "private" || post.Privacy == "almost_private" && post.Author != user.Username {
		response := map[string]string{"error": "You are not authorized to view this post"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
