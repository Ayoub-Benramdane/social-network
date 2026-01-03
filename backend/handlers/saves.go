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
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "saves") {
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

	var requestPost structs.Post
	err = json.NewDecoder(r.Body).Decode(&requestPost)
	if err != nil {
		fmt.Println("Invalid request body", err)
		resp := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var relatedGroup structs.Group
	if requestPost.GroupID != 0 {
		relatedGroup, err = database.GetGroupByID(int64(requestPost.GroupID))
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			resp := map[string]string{"error": "Failed to retrieve groups"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	fullPost, err := database.GetPost(currentUser.UserID, requestPost.PostID)
	if err != nil {
		fmt.Println("Failed to retrieve post", err)
		resp := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if fullPost.GroupID != 0 {
		if relatedGroup.PrivacyLevel == "private" {
			isMember, err := database.IsUserGroupMember(currentUser.UserID, fullPost.GroupID)
			if err != nil || !isMember {
				fmt.Println("Failed to check if user is a member", err)
				resp := map[string]string{"error": "Failed to check if user is a member"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	} else if (fullPost.PrivacyLevel == "almost_private" || fullPost.PrivacyLevel == "private") &&
		fullPost.AuthorName != currentUser.Username {

		isFollowing, err := database.IsUserFollowing(currentUser.UserID, fullPost.AuthorID)
		if err != nil || !isFollowing {
			fmt.Println("You are not authorized to view this post", err)
			resp := map[string]string{"error": "You are not authorized to view this post"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if fullPost.PrivacyLevel == "private" {
			isAuthorized, err := database.IsAuthorized(currentUser.UserID, fullPost.PostID)
			if err != nil || !isAuthorized {
				fmt.Println("You are not authorized to view this post", err)
				resp := map[string]string{"error": "You are not authorized to view this post"}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}

	alreadySaved, err := database.IsSaved(currentUser.UserID, fullPost.PostID)
	if err != nil {
		fmt.Println("Failed to check if post is saved", err)
		resp := map[string]string{"error": "Failed to check if post is saved"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !alreadySaved {
		if err := database.SavePost(currentUser.UserID, fullPost.PostID, fullPost.GroupID); err != nil {
			fmt.Println("Failed to save post", err)
			resp := map[string]string{"error": "Failed to save post"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if currentUser.UserID != fullPost.AuthorID {
			if err := database.CreateNotification(currentUser.UserID, fullPost.AuthorID, fullPost.PostID, fullPost.GroupID, 0, "save"); err != nil {
				fmt.Println("Failed to create notification", err)
				resp := map[string]string{"error": "Failed to create notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	} else {
		if err := database.UnsavePost(currentUser.UserID, fullPost.PostID, fullPost.GroupID); err != nil {
			fmt.Println("Failed to unsave post", err)
			resp := map[string]string{"error": "Failed to unsave post"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := database.DeleteNotification(currentUser.UserID, fullPost.AuthorID, fullPost.PostID, fullPost.GroupID, 0, "save"); err != nil {
			fmt.Println("Failed to delete notification", err)
			resp := map[string]string{"error": "Failed to delete notification"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	fullPost.IsSaved = !alreadySaved
	fullPost.SaveCount, err = database.CountSaves(fullPost.PostID, fullPost.GroupID)
	if err != nil {
		fmt.Println("Failed to count saves", err)
		resp := map[string]string{"error": "Failed to count saves"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullPost)
}

func GetSavedPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	requestType := r.URL.Query().Get("type")

	var savedPosts []structs.Post
	if requestType == "post" {
		savedPosts, err = database.GetSavedPosts(currentUser.UserID, 0)
	} else if requestType == "group" {
		savedPosts, err = database.GetSavedPosts(currentUser.UserID, 1)
	} else {
		fmt.Println("Invalid type parameter")
		resp := map[string]string{"error": "Invalid type parameter"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err != nil {
		fmt.Println("Failed to retrieve saved posts", err)
		resp := map[string]string{"error": "Failed to retrieve saved posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	for i := 0; i < len(savedPosts); i++ {
		savedPosts[i].SaveCount, err = database.CountSaves(savedPosts[i].PostID, savedPosts[i].GroupID)
		if err != nil {
			fmt.Println("Failed to count saves", err)
			resp := map[string]string{"error": "Failed to count saves"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedPosts)
}
