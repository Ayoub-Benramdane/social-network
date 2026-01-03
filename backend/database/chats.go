package database

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	structs "social-network/data"
)

var mu sync.Mutex

func FetchUserConnections(userID int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT DISTINCT u.id, u.username, u.firstname, u.lastname, u.avatar, u.privacy
		 FROM users u
		 JOIN messages m ON (u.id = m.sender_id OR u.id = m.receiver_id)
		 WHERE (m.sender_id = ? OR m.receiver_id = ?)
		 AND m.group_id = 0
		 GROUP BY u.id
		 ORDER BY MAX(m.created_at) DESC`,
		userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []structs.User

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.AvatarURL,
			&user.PrivacyLevel,
		); err != nil {
			return nil, err
		}

		user.IsFollowing, err = IsUserFollowing(userID, user.UserID)
		if err != nil {
			return nil, err
		}

		user.IsOnline = structs.ConnectedClients[user.UserID] != nil

		if user.UserID != userID {
			user.MessageCount, err = CountConversationUnreadMessages(user.UserID, userID, 0)
			if err != nil {
				return nil, err
			}
			connections = append(connections, user)
		}
	}

	return connections, nil
}

func CreateMessage(senderID, receiverID, groupID int64, content, image string) error {
	fmt.Println("sender, receiver, content", senderID, receiverID, content)
	mu.Lock()
	defer mu.Unlock()

	unreadCount := 0
	err := Database.QueryRow(
		`SELECT messages_not_read
		 FROM messages
		 WHERE sender_id = ? AND receiver_id = ? AND group_id = ?
		 ORDER BY created_at DESC
		 LIMIT 1`,
		senderID, receiverID, groupID,
	).Scan(&unreadCount)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}

	if senderID == receiverID {
		unreadCount = -1
	}

	_, err = Database.Exec(
		`INSERT INTO messages (sender_id, receiver_id, group_id, content, messages_not_read)
		 VALUES (?, ?, ?, ?, ?)`,
		senderID, receiverID, groupID, content, unreadCount+1,
	)

	return err
}

func FetchConversation(userID, otherUserID, offset int64) ([]structs.Message, error) {
	rows, err := Database.Query(
		`SELECT m.id, u.username, u.avatar, m.content, m.created_at
		 FROM messages m
		 JOIN users u ON u.id = m.sender_id
		 WHERE ((m.sender_id = ? AND m.receiver_id = ?)
		    OR (m.sender_id = ? AND m.receiver_id = ?))
		 AND m.group_id = 0
		 ORDER BY m.created_at DESC
		 LIMIT ? OFFSET ?`,
		userID, otherUserID, otherUserID, userID, 20, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var msg structs.Message
		var createdAt time.Time

		if err := rows.Scan(
			&msg.MessageID,
			&msg.Username,
			&msg.AvatarURL,
			&msg.Content,
			&createdAt,
		); err != nil {
			return nil, err
		}

		msg.CreatedAt = TimeAgo(createdAt)
		messages = append(messages, msg)
	}

	return messages, nil
}

func FetchGroupConversation(groupID, userID, offset int64) ([]structs.Message, error) {
	rows, err := Database.Query(
		`SELECT c.id, u.username, u.avatar, c.content, c.sender_id, c.created_at
		 FROM messages c
		 JOIN users u ON u.id = c.sender_id
		 WHERE c.group_id = ? AND c.receiver_id = ?
		 ORDER BY c.created_at DESC
		 LIMIT ? OFFSET ?`,
		groupID, userID, 10, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var msg structs.Message
		var createdAt time.Time

		if err := rows.Scan(
			&msg.MessageID,
			&msg.Username,
			&msg.AvatarURL,
			&msg.Content,
			&msg.SenderID,
			&createdAt,
		); err != nil {
			return nil, err
		}

		msg.CurrentUserID = userID
		msg.CreatedAt = TimeAgo(createdAt)
		messages = append(messages, msg)
	}

	return messages, nil
}

func CountUserUnreadMessages(userID int64, groups []structs.Group) (int64, int64, error) {
	var privateCount int64
	var groupCount int64

	err := Database.QueryRow(
		`SELECT messages_not_read
		 FROM messages
		 WHERE receiver_id = ? AND group_id = ?
		 ORDER BY created_at DESC
		 LIMIT 1`,
		userID, 0,
	).Scan(&privateCount)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, 0, err
	}

	for _, group := range groups {
		var unread int64
		err = Database.QueryRow(
			`SELECT messages_not_read
			 FROM messages
			 WHERE receiver_id = ? AND group_id == ?
			 ORDER BY created_at DESC
			 LIMIT 1`,
			userID, group.GroupID,
		).Scan(&unread)

		if err != nil && err.Error() != "sql: no rows in result set" {
			return 0, 0, err
		}

		groupCount += unread
	}

	return privateCount, groupCount, nil
}

func CountConversationUnreadMessages(senderID, receiverID, groupID int64) (int64, error) {
	var total int64
	var rows *sql.Rows
	var err error

	if groupID == 0 {
		rows, err = Database.Query(
			`SELECT messages_not_read
			 FROM messages
			 WHERE sender_id = ? AND receiver_id = ? AND group_id = ?
			 ORDER BY created_at DESC
			 LIMIT 1`,
			senderID, receiverID, groupID,
		)
	} else {
		rows, err = Database.Query(
			`SELECT messages_not_read
			 FROM messages
			 WHERE receiver_id = ? AND group_id = ?
			 ORDER BY created_at DESC
			 LIMIT 1`,
			receiverID, groupID,
		)
	}

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var count int64
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		total += count
	}

	return total, nil
}

func MarkMessagesAsRead(senderID, receiverID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	unreadCount, err := CountConversationUnreadMessages(senderID, receiverID, groupID)
	if err != nil {
		return err
	}

	if unreadCount > 0 {
		if groupID == 0 {
			_, err = Database.Exec(
				"UPDATE messages SET messages_not_read = 0 WHERE receiver_id = ? AND sender_id = ? AND group_id = ?",
				receiverID, senderID, groupID,
			)
		} else {
			_, err = Database.Exec(
				"UPDATE messages SET messages_not_read = 0 WHERE receiver_id = ? AND group_id = ?",
				receiverID, groupID,
			)
		}
	}

	return err
}
