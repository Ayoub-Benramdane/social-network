package handlers

import (
	"net/http"
	"social-network/backend/database"
	structs "social-network/backend/data"
)

func GetUserFromSession(r *http.Request) (*structs.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
