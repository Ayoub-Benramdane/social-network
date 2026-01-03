package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-network/database"
)

func CheckLastActionTime(w http.ResponseWriter, r *http.Request, tableName string) bool {
	lastTimestamp, err := database.GetLastTime(tableName)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return true
	} else if err != nil {
		fmt.Println("Error fetching last action time:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch last action time",
		})
		return false
	}

	layout := "2006-01-02T15:04:05Z"
	parsedTime, err := time.Parse(layout, lastTimestamp)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse time",
		})
		return false
	}

	currentTime := time.Now()
	if currentTime.UnixMilli()-parsedTime.UnixMilli() <= 50 {
		fmt.Println("Action too soon")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Action too soon",
		})
		return false
	}

	return true
}
