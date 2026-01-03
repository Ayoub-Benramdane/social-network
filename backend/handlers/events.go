package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "social-network/data"
	"social-network/database"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "group_events") {
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

	var newEvent structs.Event
	newEvent.Name = strings.TrimSpace(r.FormValue("name"))
	newEvent.Description = strings.TrimSpace(r.FormValue("description"))
	newEvent.Location = strings.TrimSpace(r.FormValue("location"))
	newEvent.GroupID, err = strconv.ParseInt(r.FormValue("group_id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing group ID:", err)
		resp := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if newEvent.Name == "" || newEvent.Description == "" || newEvent.Location == "" {
		fmt.Println("All fields are required!")
		resp := map[string]string{"error": "All fields are required!"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	newEvent.StartDate, err = time.Parse("2006-01-02T15:04", r.FormValue("start_date"))
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		resp := map[string]string{"error": "Invalid start date"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	newEvent.EndDate, err = time.Parse("2006-01-02T15:04", r.FormValue("end_date"))
	if err != nil {
		fmt.Println("Error parsing end date:", err)
		resp := map[string]string{"error": "Invalid end date"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	groupData, err := database.GetGroupByID(newEvent.GroupID)
	if err != nil {
		fmt.Println("Error retrieving group:", err)
		resp := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isMember, err := database.IsUserGroupMember(currentUser.UserID, newEvent.GroupID)
	if err != nil || !isMember {
		fmt.Println("User is not a member of the group", err)
		resp := map[string]string{"error": "User is not a member of the group"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	errors, valid := CheckEventInputs(newEvent.Name, newEvent.Description, newEvent.Location, newEvent.StartDate, newEvent.EndDate)
	if !valid {
		fmt.Println("Validation error", errors)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var imageURL string
	imageFile, imageHeader, err := r.FormFile("eventImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Error retrieving image:", err)
		resp := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if imageFile != nil {
		imageURL, err = SaveImage(imageFile, imageHeader, "../frontend/public/events/")
		if err != nil {
			fmt.Println("Error saving image:", err)
			resp := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		parts := strings.Split(imageURL, "/public")
		imageURL = parts[1]
	} else {
		imageURL = "/inconnu/event.jpg"
	}

	eventID, err := database.CreateGroupEvent(currentUser.UserID, newEvent.Name, newEvent.Description, newEvent.Location, newEvent.StartDate, newEvent.EndDate, newEvent.GroupID, imageURL)
	if err != nil {
		fmt.Println("Error creating event:", err)
		resp := map[string]string{"error": "Failed to create event"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	groupMembers, err := database.FetchGroupMembers(currentUser.UserID, newEvent.GroupID)
	if err != nil {
		fmt.Println("Error retrieving group members:", err)
		resp := map[string]string{"error": "Failed to retrieve group members"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	for _, member := range groupMembers {
		if member.UserID != currentUser.UserID {
			if err = database.CreateNotification(
				currentUser.UserID,
				member.UserID,
				0,
				newEvent.GroupID,
				eventID,
				"event",
			); err != nil {
				fmt.Println("Error creating notification:", err)
				resp := map[string]string{"error": "Failed to create notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}

	responseEvent := structs.Event{
		EventID:     eventID,
		Name:        newEvent.Name,
		Description: newEvent.Description,
		Location:    newEvent.Location,
		StartDate:   newEvent.StartDate,
		EndDate:     newEvent.EndDate,
		GroupID:     newEvent.GroupID,
		GroupName:   groupData.Name,
		ImageURL:    imageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseEvent)
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
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

	groupID, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing group ID:", err)
		resp := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	eventID, err := strconv.ParseInt(r.URL.Query().Get("event_id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing event ID:", err)
		resp := map[string]string{"error": "Invalid event ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = database.GetGroupByID(groupID)
	if err != nil {
		fmt.Println("Error retrieving group:", err)
		resp := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isMember, err := database.IsUserGroupMember(currentUser.UserID, groupID)
	if err != nil || !isMember {
		fmt.Println("User is not a member of the group", err)
		resp := map[string]string{"error": "User is not a member of the group"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	eventData, err := database.GetEventByID(eventID, currentUser.UserID)
	if err != nil {
		fmt.Println("Error retrieving event:", err)
		resp := map[string]string{"error": "Failed to retrieve event"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	eventData.EventID = eventID
	eventData.GroupID = groupID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventData)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
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

	eventType := r.URL.Query().Get("type")
	if eventType != "my-events" && eventType != "discover" {
		fmt.Println("Invalid type events")
		resp := map[string]string{"error": "Invalid type events"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	eventsList, err := database.GetUserEvents(currentUser.UserID, eventType)
	if err != nil {
		fmt.Println("Error retrieving events:", err)
		resp := map[string]string{"error": "Failed to retrieve events"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventsList)
}

func JoinToEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	type EventRequest struct {
		GroupID int64  `json:"group_id"`
		EventID int64  `json:"event_id"`
		Type    string `json:"type"`
	}

	var request EventRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("Error decoding request body:", err)
		resp := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if request.Type != "going" && request.Type != "not_going" {
		fmt.Println("Invalid event type:", request.Type)
		resp := map[string]string{"error": "Invalid event type"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if request.Type == "not_going" {
		request.Type = "NOT GOING"
	} else {
		request.Type = "GOING"
	}

	_, err = database.GetGroupByID(request.GroupID)
	if err != nil {
		fmt.Println("Error retrieving group:", err)
		resp := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isMember, err := database.IsUserGroupMember(currentUser.UserID, request.GroupID)
	if err != nil || !isMember {
		fmt.Println("User is not a member of the group", err)
		resp := map[string]string{"error": "User is not a member of the group"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = database.GetEventByID(request.EventID, currentUser.UserID)
	if err != nil {
		fmt.Println("Error retrieving event:", err)
		resp := map[string]string{"error": "Failed to retrieve event"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isEventMember, err := database.IsEventMember(currentUser.UserID, request.EventID)
	if err != nil {
		fmt.Println("Error checking user", err)
		resp := map[string]string{"error": "Error checking user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !isEventMember {
		if err = database.JoinEvent(currentUser.UserID, request.EventID, request.Type); err != nil {
			fmt.Println("Error joining to event:", err)
			resp := map[string]string{"error": "Failed to join to event"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Left event successfully"})
	} else {
		if err = database.UpdateEventMemberType(request.EventID, request.Type); err != nil {
			fmt.Println("Error change type event:", err)
			resp := map[string]string{"error": "Failed to change type event"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Joined to event successfully"})
	}
}

func CheckEventInputs(eventName, eventDesc, eventPlace string, startAt, endAt time.Time) (map[string]string, bool) {
	fieldErrors := make(map[string]string)

	const nameLimit = 30
	const descLimit = 500
	const placeLimit = 100

	if eventName == "" {
		fieldErrors["name"] = "Name is required"
	} else if len(eventName) > nameLimit {
		fieldErrors["name"] = "Name must be at most " + strconv.Itoa(nameLimit) + " characters"
	}

	if eventDesc == "" {
		fieldErrors["description"] = "Description is required"
	} else if len(eventDesc) > descLimit {
		fieldErrors["description"] = "Description must be at most " + strconv.Itoa(descLimit) + " characters"
	}

	if eventPlace == "" {
		fieldErrors["location"] = "Location is required"
	} else if len(eventPlace) > placeLimit {
		fieldErrors["location"] = "Location must be at most " + strconv.Itoa(placeLimit) + " characters"
	}

	if startAt.IsZero() {
		fieldErrors["start_date"] = "Start date is required"
	}

	if endAt.IsZero() {
		fieldErrors["end_date"] = "End date is required"
	}

	if endAt.Before(startAt) {
		fieldErrors["end_date"] = "End date must be after start date"
	}

	if startAt.Before(time.Now()) {
		fieldErrors["start_date"] = "Start date must be in the future"
	}

	if len(fieldErrors) > 0 {
		return fieldErrors, false
	}

	return nil, true
}
