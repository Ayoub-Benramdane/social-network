package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
	"strings"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var comment structs.Comment
	comment.Content = r.FormValue("content")
	comment.PostID, err = strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid post ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	comment.GroupID, err = strconv.ParseInt(r.FormValue("group_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid post ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var imagePath string
	image, header, err := r.FormFile("commentImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/comments/")
		if err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
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

	var id int64
	if comment.GroupID == 0 {
		id, err = database.CreateComment(comment.Content, user.ID, post, imagePath)
		if err != nil {
			response := map[string]string{"error": "Failed to create comment"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		id, err = database.CreateGroupComment(comment.Content, user.ID, comment.GroupID, post, imagePath)
		if err != nil {
			response := map[string]string{"error": "Failed to create comment"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	newComment := structs.Comment{
		ID:        id,
		PostID:    comment.PostID,
		GroupID:   comment.GroupID,
		Content:   html.EscapeString(comment.Content),
		Author:    user.Username,
		CreatedAt: "Just Now",
		Image:     imagePath,
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

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var comment structs.Comment
	comment.Content = r.FormValue("content")
	comment.PostID, err = strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid post ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	comment.GroupID, err = strconv.ParseInt(r.FormValue("group_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid post ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var imagePath string
	image, header, err := r.FormFile("commentImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/comments/")
		if err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
	}

	if comment.Content == "" || len(comment.Content) > 100 {
		response := map[string]string{"error": "Comment content is required and must be less than 100 characters"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetGroupById(comment.GroupID)
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

	id, err := database.CreateGroupComment(comment.Content, user.ID, comment.GroupID, post, imagePath)
	if err != nil {
		response := map[string]string{"error": "Failed to create comment"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newComment := structs.Comment{
		ID:        id,
		PostID:    comment.PostID,
		GroupID:   comment.GroupID,
		Content:   html.EscapeString(comment.Content),
		Author:    user.Username,
		CreatedAt: "Just Now",
		Image:     imagePath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
