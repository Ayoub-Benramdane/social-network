package database

import (
	structs "social-network/data"
)

func CreateNotification(user_id, notified_id int64, content, type_notification string) error {
	_, err := DB.Exec("INSERT INTO notifications (user_id, notified_id, content, type_notification) VALUES (?, ?, ?, ?)", user_id, notified_id, content, type_notification)
	if err != nil {
		return err
	}
	return nil
}

func GetNotifications(notified_id int64) ([]structs.Notification, error) {
	var notifications []structs.Notification
	rows, err := DB.Query("SELECT id, user_id, content, type_notification FROM notifications WHERE user_id = ? ORDER BY created_at DESC", notified_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification structs.Notification
		err = rows.Scan(&notification.ID, &notification.UserID, &notification.Content, &notification.TypeNotification)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func DeleteNotification(id int64) error {
	_, err := DB.Exec("DELETE FROM notifications WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
