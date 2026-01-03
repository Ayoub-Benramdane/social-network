package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	structs "social-network/data"
	"social-network/database"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing offset:", err)
		response := map[string]string{"error": "Invalid offset"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	following, err := database.GetUserFollowing(user.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve followings", err)
		response := map[string]string{"error": "Failed to retrieve followings"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	userInfo, err := database.GetProfileInfo(user.UserID, following)
	if err != nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	suggestedUsers, err := database.GetSuggestedUsers(user.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve not following", err)
		response := map[string]string{"error": "Failed to retrieve not following"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts, err := database.GetPosts(user.UserID, offset, following)
	if err != nil {
		fmt.Println("Failed to retrieve posts", err)
		response := map[string]string{"error": "Failed to retrieve posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	bestCategories, err := database.FetchTopCategories()
	if err != nil {
		fmt.Println("Failed to retrieve best categories", err)
		response := map[string]string{"error": "Failed to retrieve best categories"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	myGroups, err := database.GetUserGroups(*user)
	if err != nil {
		fmt.Println("Failed to retrieve my groups", err)
		response := map[string]string{"error": "Failed to retrieve my groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	suggestedGroups, err := database.GetSuggestedGroups(user.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve suggested groups", err)
		response := map[string]string{"error": "Failed to retrieve suggested groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	connections, err := database.FetchUserConnections(user.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve connections", err)
		response := map[string]string{"error": "Failed to retrieve connections"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	home := struct {
		User           structs.User       `json:"user"`
		Posts          []structs.Post     `json:"posts"`
		BestCategories []structs.Category `json:"best_categories"`
		MyGroups       []structs.Group    `json:"my_groups"`
		DiscoverGroups []structs.Group    `json:"discover_groups"`
		SuggestedUsers []structs.User     `json:"suggested_users"`
		Connections    []structs.User     `json:"connections"`
	}{
		User:           userInfo,
		Posts:          posts,
		BestCategories: bestCategories,
		MyGroups:       myGroups,
		DiscoverGroups: suggestedGroups,
		SuggestedUsers: suggestedUsers,
		Connections:    connections,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}
