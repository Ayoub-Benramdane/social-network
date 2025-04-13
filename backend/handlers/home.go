package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	structs "social-network/data"
	"social-network/database"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	user_info, err := database.GetProfileInfo(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	following, err := database.GetFollowing(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve followings"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	suggested_users, err := database.GetNotFollowing(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve not following"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts, err := database.GetPosts(user.ID, following)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	best_categories, err := database.GetBestCategories()
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve best categories"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	my_groups, err := database.GetGroups(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	suggested_groups, err := database.GetOtherGroups(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitations_friends, err := database.GetInvitationsFriends(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve invitations"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitations_groups, err := database.GetInvitationsGroups(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve invitations"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	Events, err := database.GetEvents(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve events"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	connections, err := database.GetConnections(user.ID)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve connections"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var home = struct {
		User               structs.User         `json:"user"`
		Posts              []structs.Post       `json:"posts"`
		BestCategories     []structs.Category   `json:"best_categories"`
		MyGroups           []structs.Group      `json:"my_groups"`
		SuggestedGroups    []structs.Group      `json:"suggested_groups"`
		SuggestedUsers     []structs.User       `json:"suggested_users"`
		InvitationsFriends []structs.Invitation `json:"invitations_friends"`
		InvitationsGroups  []structs.Invitation `json:"invitations_groups"`
		Events             []structs.Event      `json:"events"`
		Connections        []structs.User       `json:"connections"`
	}{
		User:               user_info,
		Posts:              posts,
		BestCategories:     best_categories,
		MyGroups:           my_groups,
		SuggestedGroups:    suggested_groups,
		SuggestedUsers:     suggested_users,
		InvitationsFriends: invitations_friends,
		InvitationsGroups:  invitations_groups,
		Events:             Events,
		Connections:        connections,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}
