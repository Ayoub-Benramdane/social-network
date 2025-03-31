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
	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		NewPostGet(w, r, user)
	case http.MethodPost:
		NewPostPost(w, r, user)
	default:
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func NewPostGet(w http.ResponseWriter, r *http.Request, user *structs.User) {
	categories, err := database.GetCategories()
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		response := map[string]string{"error": "Failed to retrieve categories"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	users, err := database.GetFollowers(user.ID)
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		response := map[string]string{"error": "Failed to retrieve users"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	data := struct {
		Categories []structs.Category
		Users      []structs.User
	}{
		Categories: categories,
		Users:      users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

func NewPostPost(w http.ResponseWriter, r *http.Request, user *structs.User) {
	var post structs.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if post.Title == "" || len(post.Title) > 20 {
		response := map[string]string{"error": "Post title is required and must be less than 20 characters"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if post.Content == "" || len(post.Content) > 500 {
		response := map[string]string{"error": "Post content is required and must be less than 500 characters"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if post.Category == "" {
		response := map[string]string{"error": "Post category is required"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if post.Privacy == "" {
		response := map[string]string{"error": "Post privacy is required"}
		w.WriteHeader(http.StatusBadRequest)
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
