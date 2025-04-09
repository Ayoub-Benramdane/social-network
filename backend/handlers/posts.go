package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
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
}

func NewPostPost(w http.ResponseWriter, r *http.Request, user *structs.User) {
	var post structs.Post
	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	post.Privacy = r.FormValue("privacy")
	post.Category = r.FormValue("category")
	fmt.Println(post)

	errors, valid := ValidatePost(post.Title, post.Content, post.Category, post.Privacy)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var imagePath string
	image, header, err := r.FormFile("postImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "./data/images/")
		if err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	id, err := database.CreatePost(user.ID, post.Title, post.Content, post.Category, imagePath, post.Privacy)
	if err != nil {
		fmt.Println(err)

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
		fmt.Println(err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(5)
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

func ValidatePost(title, content, category, privacy string) (map[string]string, bool) {
	errors := make(map[string]string)
	const maxTitle = 20
	const maxContent = 300
	const maxCategory = 20

	if title == "" {
		errors["title"] = "Title is required"
	} else if len(title) > maxTitle {
		errors["title"] = "Title must be less than " + strconv.Itoa(maxTitle) + " characters"
	}

	if content == "" {
		errors["content"] = "Content is required"
	} else if len(content) > maxContent {
		errors["content"] = "Content must be less than " + strconv.Itoa(maxContent) + " characters"
	}

	if category == "" {
		errors["category"] = "Category is required"
	} else if len(category) > maxCategory {
		errors["category"] = "Category must be less than " + strconv.Itoa(maxCategory) + " characters"
	}

	if privacy == "" {
		errors["privacy"] = "Privacy is required"
	} else if privacy != "public" && privacy != "private" && privacy != "almost_private" {
		errors["privacy"] = "Privacy must be public, private, or almost_private"
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
