package handlers

import (
	"encoding/json"
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
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	following, err := database.GetFollowing(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve followings"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	not_following, err := database.GetNotFollowing(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve not following"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts, err := database.GetPosts(user.ID, following)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	best_categories, err := database.GetBestCategories()
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve best categories"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	my_groups, err := database.GetGroups(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	other_groups, err := database.GetOtherGroups(user.ID)
	if err != nil {
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
		response := map[string]string{"error": "Failed to retrieve invitations"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	Events, err := database.GetEvents()
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve events"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var home = struct {
		User               *structs.User        `json:"user"`
		Posts              []structs.Post       `json:"posts"`
		BestCategories     []structs.Category   `json:"best_categories"`
		MyGroups           []structs.Group      `json:"groups"`
		OtherGroups        []structs.Group      `json:"other_groups"`
		Following          []structs.User       `json:"following"`
		NotFollowing       []structs.User       `json:"not_following"`
		InvitationsFriends []structs.Invitation `json:"invitations_friends"`
		InvitationsGroups  []structs.Invitation `json:"invitations_groups"`
		Events             []structs.Event      `json:"events"`
	}{
		User:               user,
		Posts:              posts,
		BestCategories:     best_categories,
		MyGroups:           my_groups,
		OtherGroups:        other_groups,
		Following:          following,
		NotFollowing:       not_following,
		InvitationsFriends: invitations_friends,
		InvitationsGroups:  invitations_groups,
		Events:             Events,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}
