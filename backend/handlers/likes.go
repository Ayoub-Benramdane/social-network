package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/database"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "User not logged in"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var ids struct {
		GroupID int64 `json:"group_id"`
		PostID  int64 `json:"post_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		response := map[string]string{"error": "Invalid request"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(user.ID, ids.PostID, ids.GroupID)
	if err != nil {
		response := map[string]string{"error": "Post not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if ids.GroupID == 0 {
		count, err := database.LikePost(user.ID, post)
		if err != nil {
			response := map[string]string{"error": "Error liking post"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		post.TotalLikes = count
		post.IsLiked = !post.IsLiked
	} else {
		count, err := database.LikeGroupPost(user.ID, ids.GroupID, post)
		if err != nil {
			response := map[string]string{"error": "Error liking post"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		post.TotalLikes = count
		post.IsLiked = !post.IsLiked
	}
	fmt.Println(post)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// func LikeGroupHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		response := map[string]string{"error": "Method not allowed"}
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	user, err := GetUserFromSession(r)
// 	if err != nil || user == nil {
// 		response := map[string]string{"error": "User not logged in"}
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	var ids struct {
// 		GroupID int64 `json:"group_id"`
// 		PostID  int64 `json:"post_id"`
// 	}

// 	err = json.NewDecoder(r.Body).Decode(&ids)
// 	if err != nil {
// 		response := map[string]string{"error": "Invalid request"}
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	post, err := database.GetPost(user.ID, ids.PostID, ids.GroupID)
// 	if err != nil {
// 		response := map[string]string{"error": "Post not found"}
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(post)
// }
