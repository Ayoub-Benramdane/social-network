package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	structs "social-network/data"
	"social-network/database"
)

func CreateStoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "stories") {
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

	var savedImagePath string
	imageFile, imageHeader, err := r.FormFile("storyImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Failed to retrieve image", err)
		resp := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if imageFile != nil {
		savedImagePath, err = SaveImage(imageFile, imageHeader, "../frontend/public/stories/")
		if err != nil {
			fmt.Println("image path", savedImagePath)
			fmt.Println("Failed to save image", err)
			resp := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		parts := strings.Split(savedImagePath, "/public")
		savedImagePath = parts[1]
	} else {
		fmt.Println("No image provided")
		resp := map[string]string{"error": "No image provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	storyID, err := database.CreateStory(savedImagePath, currentUser.UserID)
	if err != nil {
		fmt.Println("Failed to create story", err)
		resp := map[string]string{"error": "Failed to create story"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	followerIDs, err := database.GetUserFollowerIDs(currentUser.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve followers", err)
		resp := map[string]string{"error": "Failed to retrieve followers"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	followerIDs = append(followerIDs, currentUser.UserID)

	if err := database.CreateStoryStatus(storyID, followerIDs); err != nil {
		fmt.Println("Failed to update story status", err)
		resp := map[string]string{"error": "Failed to update story status"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	storyList := []structs.Story{}
	storyList = append(storyList, structs.Story{
		StoryID:   storyID,
		ImageURL: savedImagePath,
		IsRead:   false,
		CreatedAt: time.Now(),
	})

	responseData := structs.Stories{}
	responseData.Items = storyList
	responseData.User = *currentUser

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func SeenStory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "stories") {
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

	var storyID int64
	err = json.NewDecoder(r.Body).Decode(&storyID)
	if err != nil {
		fmt.Println("Failed to decode request body", err)
		resp := map[string]string{"error": "Failed to decode request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := database.MarkStoryAsSeen(storyID, currentUser.UserID); err != nil {
		fmt.Println("Failed to seen story", err)
		resp := map[string]string{"error": "Failed to seen story"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Story seen"})
}
