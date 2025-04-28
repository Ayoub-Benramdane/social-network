package database

import (
	structs "social-network/data"
	"time"
)

func SendMessage(sender_id, receiver_id, group_id int64, content, image string) error {
	if group_id != 0 {
		_, err := DB.Exec("INSERT INTO group_messages (sender_id, group_id, message, status) VALUES (?, ?, ?, ?)", sender_id, group_id, content, "unread")
		return err
	}
	_, err := DB.Exec("INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)", sender_id, receiver_id, content)
	return err
}

func GetConversation(user_id, receiver_id int64) ([]structs.Message, error) {
	rows, err := DB.Query("SELECT m.id, u.username, u.avatar, m.content, m.created_at FROM messages m JOIN users u ON u.id = m.sender_id WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?) ORDER BY m.created_at ASC", user_id, receiver_id, receiver_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []structs.Message
	for rows.Next() {
		var chat structs.Message
		var date time.Time
		if err := rows.Scan(&chat.ID, &chat.Username, &chat.Avatar, &chat.Content, &date); err != nil {
			return nil, err
		}
		chat.CreatedAt = TimeAgo(date)
		chats = append(chats, chat)
	}
	return chats, nil
}

func GetGroupConversation(group_id int64) ([]structs.Message, error) {
	rows, err := DB.Query("SELECT c,id, u.username, u.avatar, c.message, c.chat_image, c.created_at FROM group_chats c JOIN users u ON u.id = c.sender_id WHERE c.group_id = ? ORDER BY c.created_at ASC", group_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []structs.Message
	for rows.Next() {
		var chat structs.Message
		var date time.Time
		if err := rows.Scan(&chat.ID, &chat.Username, &chat.Avatar, &chat.Content, &chat.Image, &date); err != nil {
			return nil, err
		}

		chat.CreatedAt = TimeAgo(date)
		chats = append(chats, chat)
	}
	return chats, nil
}

func GetCountUserMessages(user_id int64) (int64, error) {
	var count int64
	var count2 int64
	err := DB.QueryRow("SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND status = ?", user_id, "unread").Scan(&count)
	if err != nil {
		return 0, err
	}
	rows, err := DB.Query("SELECT messages_not_read FROM group_status_messages WHERE user_id = ?", user_id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count1 int64
	for rows.Next() {
		if err := rows.Scan(&count1); err != nil {
			return 0, err
		}
		count2 += count1
	}
	return count + count2, nil
}

func GetCountConversationMessages(sender_id, user_id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM messages WHERE sender_id = ? AND receiver_id = ? AND status = ?", sender_id, user_id, "unread").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
