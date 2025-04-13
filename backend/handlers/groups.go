package handlers

import (
	"encoding/json"
	"html"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
	"strings"
)

func CreateGrpoupHandler(w http.ResponseWriter, r *http.Request) {
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

	var group structs.Group
	group.Name = r.FormValue("name")
	group.Description = r.FormValue("description")
	group.Privacy = r.FormValue("privacy")

	errors, valid := ValidateGroup(group.Name, group.Description, group.Privacy)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var imagePath string
	image, header, err := r.FormFile("groupImage")
	if err != nil && err.Error() != "http: no such file" {
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/groups/")
		if err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
	} else {
		imagePath = "/inconnu/Group.jpeg"
	}

	var coverPath string
	cover, header, err := r.FormFile("cover")
	if err != nil && err.Error() != "http: no such file" {
		response := map[string]string{"error": "Failed to retrieve cover"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if cover != nil {
		coverPath, err = SaveImage(cover, header, "../frontend/public/covers/")
		if err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(coverPath, "/public")
		coverPath = newpath[1]
	} else {
		coverPath = "/inconnu/cover.jpg"
	}

	id_group, err := database.CreateGroup(user.ID, group.Name, group.Description, imagePath, coverPath, group.Privacy)
	if err != nil {
		response := map[string]string{"error": "Failed to create group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.JoinGroup(user.ID, id_group); err != nil {
		response := map[string]string{"error": "Failed to add user to group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newGroup := structs.Group{
		ID:          id_group,
		Name:        html.EscapeString(group.Name),
		Description: html.EscapeString(group.Description),
		Image:       group.Image,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGroup)
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
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

	group_id, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	group, err := database.GetGroup(group_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	events, err := database.GetEventGroup(group_id)
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve events"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	member, err := database.IsMemberGroup(user.ID, group_id)
	if err != nil {
		response := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts := []structs.Post{}
	members := []structs.User{}
	invitations := []structs.Invitation{}
	if member || group.Privacy == "public" {
		posts, err = database.GetPostsGroup(group_id, user.ID)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve posts"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		members, err = database.GetGroupMembers(group_id)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve members"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if group.Admin == user.Username {
		invitations, err = database.GetInvitationsGroups(group_id)
		if err != nil {
			response := map[string]string{"error": "Failed to retrieve invitations"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	data := struct {
		Group       structs.Group
		Events      []structs.Event
		Members     []structs.User
		Posts       []structs.Post
		Invitations []structs.Invitation
	}{
		Group:       group,
		Events:      events,
		Members:     members,
		Posts:       posts,
		Invitations: invitations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ValidateGroup(title, content, privacy string) (map[string]string, bool) {
	errors := make(map[string]string)
	const maxTitle = 20
	const maxContent = 300

	if title == "" {
		errors["title"] = "Title is required"
	} else if len(title) > maxTitle {
		errors["title"] = "Title must be less than " + strconv.Itoa(maxTitle) + " characters"
	}

	if content == "" {
		errors["content"] = "Content is required"
	} else if len(content) > maxContent {
		errors["content"] = "Content must be less than " + strconv.Itoa(maxContent) + " characters"
	}

	if privacy == "" {
		errors["privacy"] = "Privacy is required"
	} else if privacy != "public" && privacy != "private" {
		errors["privacy"] = "Privacy must be public or private"
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
