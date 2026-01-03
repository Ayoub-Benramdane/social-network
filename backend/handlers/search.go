package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	structs "social-network/data"
	"social-network/database"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		resp := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	searchQuery := strings.TrimSpace(r.URL.Query().Get("query"))
	if searchQuery == "" {
		fmt.Println("Empty search query")
		resp := map[string]string{"message": "Empty search query"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = SaveLastSearch(currentUser.UserID, searchQuery)
	if err != nil {
		fmt.Println("Failed to update last search", err)
		resp := map[string]string{"error": "Failed to update last search"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var searchUsers bool
	var searchGroups bool
	var searchEvents bool
	var searchPosts bool

	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println("Invalid offset", err)
		resp := map[string]string{"error": "Invalid offset"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	searchType := r.URL.Query().Get("type")
	if searchType == "all" {
		searchUsers = true
		searchGroups = true
		searchEvents = true
		searchPosts = true
	} else if searchType == "users" {
		searchUsers = true
	} else if searchType == "groups" {
		searchGroups = true
	} else if searchType == "events" {
		searchEvents = true
	} else if searchType == "posts" {
		searchPosts = true
	} else {
		fmt.Println("Invalid search type")
		resp := map[string]string{"error": "Invalid search type"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var foundUsers []structs.User
	var foundGroups []structs.Group
	var foundEvents []structs.Event
	var foundPosts []structs.Post

	if searchUsers {
		foundUsers, err = database.SearchUsers(searchQuery, offset)
		if err != nil {
			fmt.Println("Error searching users:", err)
			resp := map[string]string{"error": "Failed to search users"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	if searchGroups {
		foundGroups, err = database.SearchGroups(searchQuery, offset)
		if err != nil {
			fmt.Println("Error searching groups:", err)
			resp := map[string]string{"error": "Failed to search groups"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	if searchEvents {
		foundEvents, err = database.SearchEvents(searchQuery, offset)
		if err != nil {
			fmt.Println("Error searching events:", err)
			resp := map[string]string{"error": "Failed to search events"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	if searchPosts {
		foundPosts, err = database.SearchPosts(currentUser.UserID, searchQuery, offset)
		if err != nil {
			fmt.Println("Error searching posts:", err)
			resp := map[string]string{"error": "Failed to search posts"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	responseData := struct {
		Users  []structs.User  `json:"users"`
		Groups []structs.Group `json:"groups"`
		Events []structs.Event `json:"events"`
		Posts  []structs.Post  `json:"posts"`
	}{
		Users:  foundUsers,
		Groups: foundGroups,
		Events: foundEvents,
		Posts:  foundPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func SaveLastSearch(userID int64, searchQuery string) error {
	searchCount, err := database.GetCountSearchUser(userID)
	if err != nil {
		return err
	}

	if searchCount < 3 {
		return database.InsertSearch(userID, searchQuery)
	}

	firstSearchID, err := database.GetIDFirstSearch(userID)
	if err != nil {
		return err
	}

	return database.UpdateFirstSearch(firstSearchID, searchQuery)
}
