package database

import (
	"fmt"
	"sync"
	"time"

	structs "social-network/data"
)

var (
	ConnectedClients = structs.ConnectedClients
	wsMutex          sync.Mutex
)

func sendWebSocketNotification(userID int64, payload map[string]interface{}) {
	if clients, exists := ConnectedClients[userID]; exists {
		for _, client := range clients {
			if err := client.WriteJSON(payload); err != nil {
				fmt.Println("websocket send error:", err)
				return
			}
		}
	}
}

func CreateNotification(actorUserID, receiverUserID, postID, groupID, eventID int64, notificationType string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(`
		INSERT INTO notifications 
		(user_id, notified_id, post_id, group_id, event_id, type_notification)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		actorUserID, receiverUserID, postID, groupID, eventID, notificationType,
	)

	wsMutex.Lock()
	sendWebSocketNotification(receiverUserID, map[string]interface{}{
		"type": "notifications",
	})
	wsMutex.Unlock()

	return err
}

func GetNotifications(receiverUserID, offset int64) ([]structs.Notification, error) {
	var notifications []structs.Notification

	rows, err := Database.Query(`
		SELECT n.id, u.id, u.username, u.avatar,
		       n.post_id, n.group_id, n.event_id,
		       n.type_notification, n.read, n.created_at
		FROM notifications n
		JOIN users u ON u.id = n.user_id
		WHERE n.notified_id = ?
		ORDER BY n.created_at DESC
		LIMIT 20 OFFSET ?
	`, receiverUserID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification structs.Notification
		var createdAt time.Time

		if err := rows.Scan(&notification.NotificationID, &notification.UserID, &notification.Username, &notification.AvatarURL, &notification.PostID,
			&notification.GroupID, &notification.EventID, &notification.NotificationType, &notification.IsRead, &createdAt); err != nil {
			return nil, err
		}

		notification.CreatedAt = TimeAgo(createdAt)
		notification.Message = buildNotificationMessage(notification.NotificationType)
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func buildNotificationMessage(notificationType string) string {
	switch notificationType {
	case "like":
		return "liked your post"
	case "comment":
		return "commented on your post"
	case "save":
		return "saved your post"
	case "follow":
		return "started following you"
	case "follow_request":
		return "sent you a follow request"
	case "group":
		return "invited you to join a group"
	case "join_request":
		return "requested to join your group"
	case "join":
		return "joined your group"
	case "event":
		return "created an event"
	default:
		return "sent you a notification"
	}
}

func CountUnreadNotifications(userID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(`
		SELECT COUNT(*) 
		FROM notifications 
		WHERE notified_id = ? AND read = 0
	`, userID).Scan(&count)

	return count, err
}

func DeleteNotification(actorUserID, receiverUserID, postID, groupID, eventID int64, notificationType string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(`
		DELETE FROM notifications
		WHERE user_id = ?
		  AND notified_id = ?
		  AND post_id = ?
		  AND group_id = ?
		  AND event_id = ?
		  AND type_notification = ?
	`,
		actorUserID, receiverUserID, postID, groupID, eventID, notificationType,
	)

	return err
}

func MarkNotificationAsRead(userID, notificationID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(`
		UPDATE notifications
		SET read = 1
		WHERE notified_id = ? AND id = ?
	`, userID, notificationID)

	return err
}

func MarkAllNotificationsAsRead(userID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(`
		UPDATE notifications
		SET read = 1
		WHERE notified_id = ?
	`, userID)

	return err
}
