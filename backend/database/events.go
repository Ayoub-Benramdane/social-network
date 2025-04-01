package database

import structs "social-network/data"

func CreateEvent(user_id int64, event structs.Event) (int64, error) {
	result, err := DB.Exec("INSERT INTO events (created_by, user_id, group_id, name, description, start_date, end_date, location) VALUES (?, ?, ?, ?, ?, ?, ?)", user_id, user_id, event.GroupID, event.Name, event.Description, event.StartDate, event.EndDate, event.Location)
	if err != nil {
		return 0, err
	}

	members, err := GetGroupMembers(event.GroupID)
	if err != nil {
		return 0, err
	}

	for _, member := range members {
		if member.ID != user_id {
			if err = CreateNotification(user_id, 0, member.ID, "event_created"); err != nil {
				return 0, err
			}
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetEvents(user_id int64) ([]structs.Event, error) {
	rows, err := DB.Query("SELECT e.id, u.username, g.name, e.name, e.description, e.date, e.location FROM events e JOIN users u ON e.created_by = u.id JOIN groups g ON e.group_id = g.id WHERE e.user_id = ? ORDER BY e.date DESC LIMIT 5", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []structs.Event
	for rows.Next() {
		var event structs.Event
		err = rows.Scan(&event.ID, &event.Creator, &event.Name, &event.StartDate, &event.EndDate, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEvent(id int64) (structs.Event, error) {
	var event structs.Event
	err := DB.QueryRow("SELECT u.username, g.name, e.name, e.description, e.start_date, e.end_date, e.location FROM events e JOIN users u ON u.id = e.created_by JOIN groups g ON g.id = e.group_id WHERE e.id = ?", id).Scan(&event.Creator, &event.Group, &event.Name, &event.Description, &event.StartDate, &event.EndDate, &event.Location)
	if err != nil {
		return structs.Event{}, err
	}
	return event, nil
}

func GetEventGroup(group_id int64) ([]structs.Event, error) {
	rows, err := DB.Query("SELECT e.id, u.username, e.name, e.description, e.start_date, e.end_date, e.location FROM events e JOIN users u ON u.id = e.created_by WHERE  e.group_id = ?", group_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []structs.Event
	for rows.Next() {
		var event structs.Event
		err = rows.Scan(&event.ID, &event.Creator, &event.Name, &event.Description, &event.StartDate, &event.EndDate, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
