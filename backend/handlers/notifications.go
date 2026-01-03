package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"social-network/database"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
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

	queryOffset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing offset:", err)
		resp := map[string]string{"error": "Invalid offset"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	notificationList, err := database.GetNotifications(currentUser.UserID, queryOffset)
	if err != nil {
		fmt.Println("Failed to retrieve notifications:", err)
		resp := map[string]string{"error": "Failed to retrieve notifications"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notificationList)
}

func MarkNotificationAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "notifications") {
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

	notificationID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing notification ID:", err)
		resp := map[string]string{"error": "Invalid notification ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := database.MarkNotificationAsRead(currentUser.UserID, notificationID); err != nil {
		resp := map[string]string{"error": "Failed to mark notification as read"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := map[string]string{"message": "Notification marked as read"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func MarkNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		resp := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !CheckLastActionTime(w, r, "notifications") {
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

	if err := database.MarkAllNotificationsAsRead(currentUser.UserID); err != nil {
		resp := map[string]string{"error": "Failed to mark notifications as read"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := map[string]string{"message": "Notifications marked as read"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
