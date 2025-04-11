package handlers

import (
	"encoding/json"
	"net/http"
	structs "social-network/data"
	"social-network/database"
	"strconv"
	"strings"
	"time"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
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

	var event structs.Event
	event.Name = r.FormValue("name")
	event.Description = r.FormValue("description")
	event.Location = r.FormValue("location")
	event.GroupID, err = strconv.ParseInt(r.FormValue("group_id"), 10, 64)
	if err != nil {
		response := map[string]string{"error": "Invalid group ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	event.StartDate, err = time.Parse("2006-01-02", r.FormValue("start_date"))
	if err != nil {
		response := map[string]string{"error": "Invalid start date"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	event.EndDate, err = time.Parse("2006-01-02", r.FormValue("end_date"))
	if err != nil {
		response := map[string]string{"error": "Invalid end date"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	errors, valid := ValidateEvent(event.Name, event.Description, event.Location, event.StartDate, event.EndDate)
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
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/events/")
		if err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
	} else {
		imagePath = "/inconnu/events.jpg"
	}
	id, err := database.CreateEvent(user.ID, event.Name, event.Description, event.Location, event.StartDate, event.EndDate, event.GroupID, imagePath)
	if err != nil {
		response := map[string]string{"error": "Failed to create event"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return

	}

	newEvent := structs.Event{
		ID:          id,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		GroupID:     event.GroupID,
		Image:       imagePath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEvent)
}

func ValidateEvent(name, description, location string, startDate, endDate time.Time) (map[string]string, bool) {
	errors := make(map[string]string)
	const maxNameLength = 30
	const maxDescriptionLength = 500
	const maxLocationLength = 100

	if name == "" {
		errors["name"] = "Name is required"
	} else if len(name) > maxNameLength {
		errors["name"] = "Name must be at most " + strconv.Itoa(maxNameLength) + " characters"
	}

	if description == "" {
		errors["description"] = "Description is required"
	} else if len(description) > maxDescriptionLength {
		errors["description"] = "Description must be at most " + strconv.Itoa(maxDescriptionLength) + " characters"
	}

	if location == "" {
		errors["location"] = "Location is required"
	} else if len(location) > maxLocationLength {
		errors["location"] = "Location must be at most " + strconv.Itoa(maxLocationLength) + " characters"
	}

	if startDate.IsZero() {
		errors["start_date"] = "Start date is required"
	} else if startDate.After(endDate) {
		errors["start_date"] = "Start date must be before end date"
	}

	if endDate.IsZero() {
		errors["end_date"] = "End date is required"
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
