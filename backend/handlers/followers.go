package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
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

	var user_id int64
	err = json.NewDecoder(r.Body).Decode(&user_id)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	userToFollowing, err := database.GetUserById(user_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	isFollowed, err := database.IsFollowed(user.ID, user_id)
	if err != nil {
		response := map[string]string{"error": "Failed to check if user is followed"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !isFollowed {
		if userToFollowing.Privacy == "public" {
			if err := database.AddFollower(user.ID, user_id); err != nil {
				response := map[string]string{"error": "Failed to follow user"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := map[string]string{"message": "User followed successfully"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			if err := database.CreateInvitation(user.ID, user_id); err != nil {
				response := map[string]string{"error": "Failed to follow user"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := map[string]string{"message": "User send invitation"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	invitation, err := database.CheckInvitation(user.ID, user_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if invitation {
		if err := database.DeleteInvitation(user.ID, user_id); err != nil {
			response := map[string]string{"error": "Failed to unfollow user"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := map[string]string{"message": "User unfollowed successfully"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.RemoveFollower(user.ID, user_id); err != nil {
		response := map[string]string{"error": "Failed to unfollow user"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"message": "User unfollowed successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func FollowersHandler(w http.ResponseWriter, r *http.Request) {
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

	var followers []structs.User
	if followed || info.Privacy == "public" || user_id == user.ID {
		followers, err = database.GetFollowers(user.ID)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve followers"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func FollowingHandler(w http.ResponseWriter, r *http.Request) {
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

	var following []structs.User
	if followed || info.Privacy == "public" || user_id == user.ID {
		following, err = database.GetFollowing(user.ID)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve following"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}
