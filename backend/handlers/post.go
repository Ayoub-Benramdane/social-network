package handlers

import (
	"database/sql"
	"encoding/json"
	"html"
	"log"
	"net/http"
	"social-network/backend/database"
	structs "social-network/backend/structs"
	"time"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		response := map[string]string{"error": "Unauthorized"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{"error": "Invalid email or password"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("Database error: %v", err)
			response := map[string]string{"error": "Internal server error"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	var post struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Author   string   `json:"author"`
		Categories []string `json:"category"`
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := database.CreatePost(user.ID, post.Title, post.Content, post.Author, post.Category)
	if err != nil {
		log.Printf("Database error: %v", err)
		response := map[string]string{"error": "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	newPost := structs.Post{
		ID:            id,
		Author:        user.Username,
		Title:         html.EscapeString(post.Title),
		Content:       html.EscapeString(post.Content),
		CreatedAt:     "Just Now",
		TotalLikes:    0,
		TotalDislikes: 0,
		UserID:        user.ID,
		Categories:    post.Categories,
		TotalComments: 0,
		Comments:      nil,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	// response := make(map[string]string)
	// if r.Method != http.MethodGet {
	// 	response = map[string]string{"error": "Method not allowed"}
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
}
