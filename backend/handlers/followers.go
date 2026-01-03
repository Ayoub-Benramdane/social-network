package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	structs "social-network/data"
	"social-network/database"
)

func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user"})
		return
	}

	targetUserID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid user ID", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	isFollowing, err := database.IsUserFollowing(currentUser.UserID, targetUserID)
	if err != nil {
		fmt.Println("Failed to check follow status", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve followings"})
		return
	}

	profileInfo, err := database.GetProfileInfo(targetUserID, nil)
	if err != nil {
		fmt.Println("Failed to retrieve profile", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve profile"})
		return
	}

	var followersList []structs.User
	if isFollowing || profileInfo.PrivacyLevel == "public" || targetUserID == currentUser.UserID {
		followersList, err = database.GetUserFollowers(targetUserID)
		if err != nil {
			fmt.Println("Failed to retrieve followers", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve followers"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followersList)
}

func GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user"})
		return
	}

	targetUserID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid user ID", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	profileInfo, err := database.GetProfileInfo(targetUserID, nil)
	if err != nil {
		fmt.Println("Failed to retrieve profile", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve profile"})
		return
	}

	isFollowing, err := database.IsUserFollowing(currentUser.UserID, targetUserID)
	if err != nil {
		fmt.Println("Failed to check follow status", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve followings"})
		return
	}

	var followingList []structs.User
	if isFollowing || profileInfo.PrivacyLevel == "public" || targetUserID == currentUser.UserID {
		followingList, err = database.GetUserFollowing(targetUserID)
		if err != nil {
			fmt.Println("Failed to retrieve following", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve following"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followingList)
}

func GetSuggestedUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user"})
		return
	}

	requestType := r.URL.Query().Get("type")

	var resultUsers []structs.User
	switch requestType {
	case "suggested":
		resultUsers, err = database.GetSuggestedUsers(currentUser.UserID)
	case "received":
		resultUsers, err = database.GetReceivedFollowRequests(currentUser.UserID)
	case "pending":
		resultUsers, err = database.GetPendingFollowRequests(currentUser.UserID)
	default:
		fmt.Println("Invalid type", requestType)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid type"})
		return
	}

	if err != nil {
		fmt.Println("Failed to retrieve users", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve suggested users"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultUsers)
}
