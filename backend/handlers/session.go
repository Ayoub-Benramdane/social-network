package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	structs "social-network/data"
	"social-network/database"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Println("Not found", r.URL.Path)
		resp := map[string]string{"error": "Not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	} else if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
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

func CheckTheUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed")
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user")
		resp := map[string]string{"error": "User not logged in"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	profileInfo, err := database.GetProfileInfo(currentUser.UserID, nil)
	if err != nil {
		fmt.Println("Failed to retrieve user")
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profileInfo)
}

func GetUserFromSession(r *http.Request) (*structs.User, error) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Failed to retrieve cookie")
		return nil, err
	}

	foundUser, err := database.GetUserBySession(sessionCookie.Value)
	if err != nil {
		fmt.Println("Failed to retrieve user")
		return nil, err
	}

	return &foundUser, nil
}
