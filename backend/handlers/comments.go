package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	structs "social-network/data"
	"social-network/database"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	var newComment structs.Comment
	newComment.Content = strings.TrimSpace(r.FormValue("content"))
	newComment.PostID, err = strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing post ID:", err)
		resp := map[string]string{"error": "Invalid post ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var imageURL string
	imageFile, imageHeader, err := r.FormFile("commentImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Error retrieving image:", err)
		resp := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if imageFile != nil {
		imageURL, err = SaveImage(imageFile, imageHeader, "../frontend/public/comments/")
		if err != nil {
			fmt.Println("Error saving image:", err)
			resp := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		parts := strings.Split(imageURL, "/public")
		imageURL = parts[1]
	}

	if newComment.Content == "" || len(newComment.Content) > 100 {
		fmt.Println("Comment content is required and must be less than 100 characters")
		resp := map[string]string{"error": "Comment content is required and must be less than 100 characters"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	postData, err := database.GetPost(currentUser.UserID, newComment.PostID)
	if err != nil {
		fmt.Println("Failed to retrieve post", err)
		resp := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var commentID int64
	notificationType := "comment"

	commentID, err = database.CreatePostComment(newComment.Content, currentUser.UserID, postData, imageURL)
	if err != nil {
		fmt.Println("Failed to create comment", err)
		resp := map[string]string{"error": "Failed to create comment"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if postData.AuthorID != currentUser.UserID {
		if err := database.CreateNotification(currentUser.UserID, postData.AuthorID, postData.PostID, postData.GroupID, 0, notificationType); err != nil {
			fmt.Println("Failed to create notification", err)
			resp := map[string]string{"error": "Failed to create notification"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	responseComment := structs.Comment{
		CommentID: commentID,
		PostID:    newComment.PostID,
		Content:   html.EscapeString(newComment.Content),
		AuthorID:  currentUser.UserID,
		Username:  currentUser.Username,
		AvatarURL: currentUser.AvatarURL,
		CreatedAt: "Just Now",
		ImageURL:  imageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseComment)
}
