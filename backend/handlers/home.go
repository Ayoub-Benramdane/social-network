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

	followers, err := database.GetFollowers(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve followers"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	followers, err = database.GetFollowers(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve followers"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts, err := database.GetPosts(user.ID, followers)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve posts"}
		w.WriteHeader(http.StatusInternalServerError)
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

	var home = struct {
		User         *structs.User   `json:"user"`
		Posts        []structs.Post  `json:"posts"`
		MyGroups     []structs.Group `json:"groups"`
		OtherGroups  []structs.Group `json:"other_groups"`
		Followers    []structs.User  `json:"followers"`
		Following    []structs.User  `json:"following"`
		NotFollowing []structs.User  `json:"not_following"`
	}{
		User:         user,
		Posts:        posts,
		MyGroups:     my_groups,
		OtherGroups:  other_groups,
		Followers:    followers,
		Following:    following,
		NotFollowing: not_following,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}
