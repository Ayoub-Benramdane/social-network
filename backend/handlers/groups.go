package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "social-network/data"
	"social-network/database"
)

func CreateGrpoupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "groups") {
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var group structs.Group
	group.Name = strings.TrimSpace(r.FormValue("name"))
	group.Description = strings.TrimSpace(r.FormValue("description"))
	group.PrivacyLevel = strings.TrimSpace(r.FormValue("privacy"))

	errors, valid := ValidateGroup(group.Name, group.Description, group.PrivacyLevel)
	if !valid {
		fmt.Println("Validation error", errors)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var imagePath string
	imageFile, imageHeader, err := r.FormFile("groupImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Error retrieving image:", err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if imageFile != nil {
		imagePath, err = SaveImage(imageFile, imageHeader, "../frontend/public/groups/")
		if err != nil {
			fmt.Println("Error saving image:", err)
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		imagePath = strings.Split(imagePath, "/public")[1]
	} else {
		imagePath = "/inconnu/group.jpeg"
	}

	var coverPath string
	coverFile, coverHeader, err := r.FormFile("groupCover")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Error retrieving cover:", err)
		response := map[string]string{"error": "Failed to retrieve cover"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if coverFile != nil {
		coverPath, err = SaveImage(coverFile, coverHeader, "../frontend/public/covers/")
		if err != nil {
			fmt.Println("Error saving cover:", err)
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		coverPath = strings.Split(coverPath, "/public")[1]
	} else {
		coverPath = "/inconnu/cover.jpg"
	}
	groupID, err := database.CreateGroup(user.UserID, group.Name, group.Description, imagePath, coverPath, group.PrivacyLevel)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint") {
		fmt.Println("Group name already exists")
		response := map[string]string{"error": "Group name already exists"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		fmt.Println("Failed to create group:", err)
		response := map[string]string{"error": "Failed to create group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.AddUserToGroup(user.UserID, groupID); err != nil {
		fmt.Println("Failed to add user to group:", err)
		response := map[string]string{"error": "Failed to add user to group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newGroup := structs.Group{
		GroupID:      groupID,
		AdminName:    user.Username,
		Name:         html.EscapeString(group.Name),
		Description:  html.EscapeString(group.Description),
		CreatedAt:    time.Now().Format("2006-01-02 15:04"),
		ImageURL:     imagePath,
		CoverURL:     coverPath,
		PrivacyLevel: group.PrivacyLevel,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGroup)
}

func AddMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodGet {
		groupID, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
		if err != nil {
			fmt.Println("Invalid group ID", err)
			response := map[string]string{"error": "Invalid group ID"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err = database.GetGroupByID(groupID)
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			response := map[string]string{"error": "Failed to retrieve group"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		users, err := database.GetInvitableUsers(user.UserID, groupID)
		if err != nil {
			log.Printf("Error retrieving users: %v", err)
			response := map[string]string{"error": "Failed to retrieve users"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	} else {
		type IDS struct {
			UserID  int64 `json:"user_id"`
			GroupID int64 `json:"group_id"`
		}

		var ids IDS
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			fmt.Println("Failed to decode request body", err)
			response := map[string]string{"error": "Failed to decode request body"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err = database.GetGroupByID(ids.GroupID)
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			response := map[string]string{"error": "Failed to retrieve group"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if _, err = database.FindUserByID(ids.UserID); err != nil {
			fmt.Println("Invalid user", err)
			response := map[string]string{"error": "Invalid user"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if isMember, err := database.IsUserGroupMember(ids.UserID, ids.GroupID); err != nil || isMember {
			fmt.Println("User is already a member of the group", err)
			response := map[string]string{"error": "User is already a member of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err = database.CreateInvitation(user.UserID, ids.UserID, ids.GroupID); err != nil {
			fmt.Println("Failed to send invitation to this user", err)
			response := map[string]string{"error": "Failed to send invitation to this user"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := database.CreateNotification(user.UserID, ids.UserID, 0, ids.GroupID, 0, "group"); err != nil {
			fmt.Println("Failed to create notification", err)
			response := map[string]string{"error": "Failed to create notification"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	Type := r.URL.Query().Get("type")

	var groups []structs.Group

	if Type == "suggested" {
		groups, err = database.GetSuggestedGroups(user.UserID)
	} else if Type == "pending" {
		groups, err = database.GetPendingGroups(user.UserID)
	} else if Type == "joined" {
		groups, err = database.GetUserGroups(*user)
	} else {
		fmt.Println("Invalid type groups")
		response := map[string]string{"error": "Invalid type groups"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}

	if err != nil {
		fmt.Println("Failed to retrieve groups:", err)
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	groupID, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid group ID", err)
		response := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	group, err := database.GetGroupByID(groupID)
	if err != nil {
		fmt.Println("Failed to retrieve group", err)
		response := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	member, err := database.IsUserGroupMember(user.UserID, groupID)
	if err != nil {
		fmt.Println("Failed to check if user is a member", err)
		response := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invited, err := database.HasGroupInvitation(user.UserID, groupID)
	if err != nil {
		fmt.Println("Failed to check if user has been invited to join the group", err)
		response := map[string]string{"error": "Failed to check if user has been invited to join the group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	requested, err := database.InvitationExists(user.UserID, group.AdminID, groupID)
	if err != nil {
		fmt.Println("Failed to check if user has requested to join the group", err)
		response := map[string]string{"error": "Failed to check if user has requested to join the group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if group.AdminName == user.Username {
		group.UserRole = "admin"
		group.ActionType = "Delete group"
	} else if member {
		group.UserRole = "member"
		group.ActionType = "Leave group"
	} else if requested {
		group.UserRole = "requested"
		group.ActionType = "Cancel request"
	} else if invited {
		group.InvitedBy, err = database.GetInvitedBy(user.UserID, groupID)
		if err != nil {
			fmt.Println("Failed to get invited by", err)
			response := map[string]string{"error": "Failed to get invited by"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		if group.InvitedBy == group.AdminID {
			group.IsOwner = true
		} else {
			group.IsOwner = false
		}

		group.UserRole = "invited"
		group.ActionType = "Accept invitation"
	} else {
		group.UserRole = "guest"
		group.ActionType = "Join group"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

func GroupDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrive user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	groupID, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid group ID", err)
		response := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	group, err := database.GetGroupByID(groupID)
	if err != nil {
		fmt.Println("Failed to retrieve group", err)
		response := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	member, err := database.IsUserGroupMember(user.UserID, groupID)
	if err != nil {
		fmt.Println("Failed to check if user is a member", err)
		response := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	Type := r.URL.Query().Get("type")

	if Type == "members" {
		if member || group.PrivacyLevel == "public" {
			members, err := database.FetchGroupMembers(user.UserID, groupID)
			if err != nil {
				fmt.Println("Failed to retrieve members", err)
				response := map[string]string{"error": "Failed to retrieve members"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(members)
		} else {
			fmt.Println("You are not a member of this group")
			response := map[string]string{"error": "You are not a member of this group"}
			json.NewEncoder(w).Encode(response)
		}
		return
	} else if Type == "invitations" {
		if group.AdminName == user.Username {
			invitations, err := database.GetGroupInvitationsByGroup(user.UserID, groupID)
			if err != nil {
				fmt.Println("Failed to retrieve invitations", err)
				response := map[string]string{"error": "Failed to retrieve invitations"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invitations)
		} else {
			fmt.Println("You are not the admin of this group")
			response := map[string]string{"error": "You are not the admin of this group"}
			json.NewEncoder(w).Encode(response)
		}
		return
	} else if Type == "posts" {
		if member || group.PrivacyLevel == "public" {
			posts, err := database.GetPostsGroup(groupID, user.UserID, group.PrivacyLevel)
			if err != nil {
				fmt.Println("Failed to retrieve posts", err)
				response := map[string]string{"error": "Failed to retrieve posts"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			for i := 0; i < len(posts); i++ {
				posts[i].SaveCount, err = database.CountSaves(posts[i].PostID, posts[i].GroupID)
				if err != nil {
					fmt.Println("Failed to count saves", err)
					response := map[string]string{"error": "Failed to count saves"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(posts)
		} else {
			fmt.Println("You are not a member of this group")
			response := map[string]string{"error": "You are not a member of this group"}
			json.NewEncoder(w).Encode(response)
		}
		return
	} else if Type == "events" {
		if member {
			events, err := database.GetGroupEvents(user.UserID, groupID)
			if err != nil {
				fmt.Println("Failed to retrieve events", err)
				response := map[string]string{"error": "Failed to retrieve events"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(events)
		} else {
			fmt.Println("You are not a member of this group")
			response := map[string]string{"error": "You are not a member of this group"}
			json.NewEncoder(w).Encode(response)
		}
		return
	} else {
		fmt.Println("Invalid type")
		response := map[string]string{"error": "Invalid type"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func ValidateGroup(name, description, privacy string) (map[string]string, bool) {
	errors := make(map[string]string)
	const maxName = 20
	const maxDescription = 300

	if name == "" {
		errors["name"] = "Name is required"
	} else if len(name) > maxName {
		errors["name"] = "Name must be less than " + strconv.Itoa(maxName) + " characters"
	}

	if description == "" {
		errors["description"] = "Description is required"
	} else if len(description) > maxDescription {
		errors["description"] = "Description must be less than " + strconv.Itoa(maxDescription) + " characters"
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
