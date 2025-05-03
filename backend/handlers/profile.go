package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	user_id, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid user ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return

	}

	info, err := database.GetProfileInfo(user_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve profile"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	followed, err := database.IsFollowed(user.ID, user_id)
	if err != nil {
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve followings"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var posts []structs.Post
	if followed || info.Privacy == "public" || user_id == user.ID {
		posts, err = database.GetPostsByUser(user_id, user.ID, followed)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve posts"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if user_id == user.ID {
		info.Role = "owner"
	} else {
		info.Role = "user"
	}

	data := struct {
		ProfileInfo structs.User
		Posts       []structs.Post
	}{
		ProfileInfo: info,
		Posts:       posts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
