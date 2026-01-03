package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	structs "social-network/data"
	"social-network/database"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "post_likes") {
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		resp := map[string]string{"error": "User not logged in"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var postData structs.Post
	err = json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		resp := map[string]string{"error": "Invalid request"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	postData, err = database.GetPost(currentUser.UserID, postData.PostID)
	if err != nil {
		fmt.Println("Error retrieving post:", err)
		resp := map[string]string{"error": "Post not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	notificationType := "like"
	likeCount, err := database.TogglePostLike(currentUser.UserID, postData)
	if err != nil {
		fmt.Println("Error liking post:", err)
		resp := map[string]string{"error": "Error liking post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	postData.LikeCount = likeCount
	postData.IsLiked = !postData.IsLiked

	if postData.AuthorID != currentUser.UserID {
		if postData.IsLiked {
			if err = database.CreateNotification(currentUser.UserID, postData.AuthorID, postData.PostID, postData.GroupID, 0, notificationType); err != nil {
				fmt.Println("Error creating notification:", err)
				resp := map[string]string{"error": "Error creating notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
		} else if err = database.DeleteNotification(currentUser.UserID, postData.AuthorID, postData.PostID, postData.GroupID, 0, notificationType); err != nil {
			fmt.Println("Error deleting notification:", err)
			resp := map[string]string{"error": "Error deleting notification"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	if postData.IsLiked {
		postData.LikedBy = append(postData.LikedBy, *currentUser)
	} else {
		filteredLikedBy := make([]structs.User, 0)
		for _, likedUser := range postData.LikedBy {
			if likedUser.UserID != currentUser.UserID {
				filteredLikedBy = append(filteredLikedBy, likedUser)
			}
		}
		postData.LikedBy = filteredLikedBy
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postData)
}
