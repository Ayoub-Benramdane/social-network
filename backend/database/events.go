package database

import (
	"database/sql"
	"strings"
	"time"

	structs "social-network/data"
)

func CreateGroupEvent(userID int64, name, description, location string, startDate time.Time, endDate time.Time, groupID int64, imageURL string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	result, err := Database.Exec(
		`INSERT INTO group_events 
		(created_by, group_id, name, description, start_date, end_date, location, image) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, groupID, name, description, startDate, endDate, location, imageURL,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetUserEvents(userID int64, filterType string) ([]structs.Event, error) {
	var (
		rows            *sql.Rows
		err             error
		expiredEventIDs []int64
	)

	if filterType == "my-events" {
		rows, err = Database.Query(
			`SELECT DISTINCT 
				e.id, e.created_by, g.name, g.id, e.name, e.description,
				e.start_date, e.end_date, e.location, e.created_at, e.image
			FROM group_events e
			JOIN event_members em ON e.id = em.event_id
			JOIN groups g ON e.group_id = g.id
			WHERE em.user_id = ? AND e.end_date > CURRENT_TIMESTAMP
			ORDER BY e.start_date DESC`,
			userID,
		)
	} else {
		rows, err = Database.Query(
			`SELECT DISTINCT 
				e.id, e.created_by, g.name, g.id, e.name, e.description,
				e.start_date, e.end_date, e.location, e.created_at, e.image
			FROM group_events e
			JOIN groups g ON e.group_id = g.id
			JOIN group_members gm ON gm.group_id = g.id
			WHERE gm.user_id = ?
			  AND NOT EXISTS (
				  SELECT 1 FROM event_members em 
				  WHERE em.event_id = e.id AND em.user_id = ?
			  )
			  AND e.end_date > CURRENT_TIMESTAMP
			ORDER BY e.start_date DESC`,
			userID, userID,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []structs.Event

	for rows.Next() {
		var event structs.Event
		var createdAt time.Time

		err = rows.Scan(&event.EventID, &event.CreatorName, &event.GroupName, &event.GroupID, &event.Name, &event.Description, &event.StartDate, &event.EndDate, &event.Location, &createdAt, &event.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		if event.EndDate.Before(time.Now()) {
			expiredEventIDs = append(expiredEventIDs, event.EventID)
			continue
		}

		event.CreatedAt = TimeAgo(createdAt)

		isGroupMember, err := IsUserGroupMember(userID, event.GroupID)
		if err != nil {
			return nil, err
		}

		if !isGroupMember {
			continue
		}

		if isEventMember, _ := IsEventMember(userID, event.EventID); isEventMember {
			event.EventType = GetEventMemberType(userID, event.EventID)
		}

		events = append(events, event)
	}

	for _, eventID := range expiredEventIDs {
		if err := DeleteEventByID(eventID); err != nil {
			return nil, err
		}
	}

	return events, nil
}

func GetEventByID(eventID, userID int64) (structs.Event, error) {
	var event structs.Event
	var createdAt time.Time

	err := Database.QueryRow(
		`SELECT u.username, g.name, e.name, e.description,
		        e.start_date, e.end_date, e.location, e.created_at, e.image
		 FROM group_events e
		 JOIN users u ON u.id = e.created_by
		 JOIN groups g ON g.id = e.group_id
		 WHERE e.id = ?`,
		eventID,
	).Scan(
		&event.CreatorName,
		&event.GroupName,
		&event.Name,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.Location,
		&createdAt,
		&event.ImageURL,
	)

	if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
		return structs.Event{}, err
	}

	if event.EndDate.Before(time.Now()) {
		return structs.Event{}, DeleteEventByID(eventID)
	}

	event.CreatedAt = TimeAgo(createdAt)

	if isMember, _ := IsEventMember(userID, eventID); isMember {
		event.EventType = GetEventMemberType(userID, eventID)
	}

	return event, nil
}

func GetGroupEvents(userID, groupID int64) ([]structs.Event, error) {
	rows, err := Database.Query(
		`SELECT e.id, u.username, e.name, e.description,
		        e.start_date, e.end_date, e.location, e.created_at, e.image
		 FROM group_events e
		 JOIN users u ON u.id = e.created_by
		 WHERE e.group_id = ?
		 ORDER BY e.created_at DESC`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		events          []structs.Event
		expiredEventIDs []int64
	)

	for rows.Next() {
		var event structs.Event
		var createdAt time.Time

		err = rows.Scan(
			&event.EventID,
			&event.CreatorName,
			&event.Name,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.Location,
			&createdAt,
			&event.ImageURL,
		)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		if event.EndDate.Before(time.Now()) {
			expiredEventIDs = append(expiredEventIDs, event.EventID)
			continue
		}

		event.GroupID = groupID
		event.CreatedAt = TimeAgo(createdAt)

		if isMember, _ := IsEventMember(userID, event.EventID); isMember {
			event.EventType = GetEventMemberType(userID, event.EventID)
		}

		events = append(events, event)
	}

	for _, eventID := range expiredEventIDs {
		if err := DeleteEventByID(eventID); err != nil {
			return nil, err
		}
	}

	return events, nil
}

func IsEventMember(userID, eventID int64) (bool, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM event_members WHERE user_id = ? AND event_id = ?",
		userID, eventID,
	).Scan(&count)
	return count > 0, err
}

func GetEventMemberType(userID, eventID int64) string {
	var memberType string
	_ = Database.QueryRow(
		"SELECT type FROM event_members WHERE user_id = ? AND event_id = ?",
		userID, eventID,
	).Scan(&memberType)
	return memberType
}

func CountUserEvents(userID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM event_members WHERE user_id = ?",
		userID,
	).Scan(&count)
	return count, err
}

func JoinEvent(userID, eventID int64, memberType string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO event_members (user_id, event_id, type) VALUES (?, ?, ?)",
		userID, eventID, memberType,
	)
	return err
}

func UpdateEventMemberType(eventID int64, memberType string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"UPDATE event_members SET type = ? WHERE event_id = ?",
		memberType, eventID,
	)
	return err
}

func DeleteEventByID(eventID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"DELETE FROM group_events WHERE id = ?",
		eventID,
	)
	return err
}
