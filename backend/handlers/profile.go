package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	structs "social-network/data"
	"social-network/database"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	profileUserID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing user ID:", err)
		resp := map[string]string{"error": "Invalid user ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	profileInfo, err := database.GetProfileInfo(profileUserID, nil)
	if err != nil {
		fmt.Println("Error retrieving profile:", err)
		resp := map[string]string{"error": "Failed to retrieve profile"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if profileUserID == currentUser.UserID {
		profileInfo.Role = "owner"
	} else {
		profileInfo.Role = "user"
	}

	profileInfo.IsFollowing, err = database.IsUserFollowing(currentUser.UserID, profileUserID)
	if err != nil {
		fmt.Println("Error checking follow status:", err)
		resp := map[string]string{"error": "Failed to retrieve followings"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	profileInfo.IsFollower, err = database.IsUserFollowing(profileUserID, currentUser.UserID)
	if err != nil {
		fmt.Println("Error checking follow status:", err)
		resp := map[string]string{"error": "Failed to retrieve followers"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	profileInfo.IsPending, err = database.InvitationExists(currentUser.UserID, profileUserID, 0)
	if err != nil {
		fmt.Println("Error checking invitation:", err)
		resp := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if profileInfo.PrivacyLevel == "public" || profileInfo.IsFollowing || profileUserID == currentUser.UserID {
		ClientsMutex.Lock()
		profileInfo.IsOnline = structs.ConnectedClients[profileUserID] != nil
		ClientsMutex.Unlock()
	}

	if profileInfo.IsPending {
		profileInfo.AccountType = "Pending"
	} else if profileInfo.IsFollowing {
		profileInfo.AccountType = "Unfollow"
	} else if profileInfo.IsFollower {
		profileInfo.AccountType = "Follow back"
	} else {
		profileInfo.AccountType = "Follow"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profileInfo)
}

func ProfilePostsHandler(w http.ResponseWriter, r *http.Request) {
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

	profileUserID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing user ID:", err)
		resp := map[string]string{"error": "Invalid user ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	profileInfo, err := database.GetProfileInfo(profileUserID, nil)
	if err != nil {
		fmt.Println("Error retrieving profile:", err)
		resp := map[string]string{"error": "Failed to retrieve profile"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isFollowed, err := database.IsUserFollowing(currentUser.UserID, profileUserID)
	if err != nil {
		fmt.Println("Error checking follow status:", err)
		resp := map[string]string{"error": "Failed to retrieve followings"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var userPosts []structs.Post
	if isFollowed || profileInfo.PrivacyLevel == "public" || profileUserID == currentUser.UserID {
		userPosts, err = database.GetPostsByUser(profileUserID, currentUser.UserID, isFollowed)
		if err != nil {
			fmt.Println("Error retrieving posts:", err)
			resp := map[string]string{"error": "Failed to retrieve posts"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	if profileUserID == currentUser.UserID {
		profileInfo.Role = "owner"
	} else {
		profileInfo.Role = "user"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userPosts)
}
