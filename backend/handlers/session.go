package handlers

import (
	"encoding/json"
	"net/http"
	structs "social-network/data"
	"social-network/database"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}
	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		json.NewEncoder(w).Encode(false)
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "guest",
			MaxAge: -1,
		})
		return
	}
	json.NewEncoder(w).Encode(true)
}

func GetUserFromSession(r *http.Request) (*structs.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	user, err := database.GetUserConnected(cookie.Value)
	return &user, err
}
