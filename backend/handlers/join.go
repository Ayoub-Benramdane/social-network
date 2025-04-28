package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/database"
	"strconv"
)

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
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

	var group_id int64
	group_id, err = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	// err = json.NewDecoder(r.Body).Decode(&group_id)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	group, err := database.GetGroupById(group_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	isMember, err := database.IsMemberGroup(group.ID, user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !isMember {
		if group.Privacy == "public" {
			if database.JoinGroup(user.ID, group.ID) != nil {
				response := map[string]string{"error": "Failed to join group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			fmt.Println("join group")
			response := map[string]string{"message": "Successfully joined group"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			if database.CreateInvitationGroup(user.ID, group.AdminID, group.ID) != nil {
				response := map[string]string{"error": "Failed to join group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			fmt.Println("invite")
			response := map[string]string{"message": "User sent invitation to join group"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	invitation, err := database.CheckInvitationGroup(user.ID, group.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if invitation {
		if err := database.DeleteInvitationGroup(user.ID, group.ID); err != nil {
			response := map[string]string{"error": "Failed to delete invitation"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := map[string]string{"message": "Successfully deleted invitation"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	if database.LeaveGroup(user.ID, group.ID) != nil {
		response := map[string]string{"error": "Failed to leave group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"message": "Successfully left group"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
