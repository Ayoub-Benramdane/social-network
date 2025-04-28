package handlers

import (
	"encoding/json"
	"net/http"
	structs "social-network/data"
	"social-network/database"
)

func GetInvitationsGroups(w http.ResponseWriter, r *http.Request) {
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

	invitations, err := database.GetInvitationsGroups(user.ID)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve invitations"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitations)
}

func AcceptInvitationHandler(w http.ResponseWriter, r *http.Request) {
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

	var invitation structs.Invitation
	err = json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetUserById(invitation.SenderID)
	if err != nil {
		response := map[string]string{"error": "Invalid sender"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if invitation.SenderID == user.ID {
		response := map[string]string{"error": "Cannot accept your own invitation"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetGroupById(invitation.GroupID)
	if err != nil {
		response := map[string]string{"error": "Invalid group"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetInvitationById(invitation.ID, invitation.GroupID)
	if err != nil {
		response := map[string]string{"error": "Invalid invitation"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if database.AcceptInvitation(invitation.ID, invitation.SenderID, user.ID, invitation.GroupID) != nil {
		response := map[string]string{"error": "Failed to accept invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func DeclineInvitationHandler(w http.ResponseWriter, r *http.Request) {
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

	var invitation structs.Invitation
	err = json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetUserById(invitation.SenderID)
	if err != nil {
		response := map[string]string{"error": "Invalid sender"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if invitation.SenderID == user.ID {
		response := map[string]string{"error": "Cannot accept your own invitation"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetGroupById(invitation.GroupID)
	if err != nil {
		response := map[string]string{"error": "Invalid group"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = database.GetInvitationById(invitation.ID, invitation.GroupID)
	if err != nil {
		response := map[string]string{"error": "Invalid invitation"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if database.DeclineInvitation(invitation.ID, invitation.GroupID) != nil {
		response := map[string]string{"error": "Failed to decline invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}
