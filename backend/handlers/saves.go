package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	structs "social-network/data"
	"social-network/database"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
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

	var ids struct {
		PostID  int64 `json:"post_id"`
		GroupID int64 `json:"group_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	group, err := database.GetGroupById(int64(ids.GroupID))
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(user.ID, ids.GroupID, ids.PostID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if post.GroupID != 0 {
		if group.Privacy == "private" {
			if member, err := database.IsMemberGroup(user.ID, post.GroupID); err != nil || !member {
				response := map[string]string{"error": "Failed to retrieve group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	} else if post.Privacy == "private" || post.Privacy == "almost_private" && post.Author != user.Username {
		if followed, err := database.IsFollowed(user.ID, post.UserID); err != nil || !followed {
			response := map[string]string{"error": "You are not authorized to view this post"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}
		if post.Privacy == "almost_private" {
			if authorized, err := database.IsAuthorized(user.ID, post.ID); err != nil || !authorized {
				response := map[string]string{"error": "You are not authorized to view this post"}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	if database.SavePost(user.ID, ids.PostID, ids.GroupID) != nil {
		response := map[string]string{"error": "Failed to save post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func GetSavedPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	saved_posts, err := database.GetSavedPosts(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve saved posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	saved_posts_group, err := database.GetSavedGroupPosts(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve saved posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var saved_posts_all = struct {
		Posts  []structs.Post `json:"saved_posts"`
		Groups []structs.Post `json:"saved_posts_group"`
	}{
		Posts:  saved_posts,
		Groups: saved_posts_group,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saved_posts_all)
}
