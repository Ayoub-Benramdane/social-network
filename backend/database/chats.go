package database

import structs "social-network/data"

func GetConversation(user_id, receiver_id int64) ([]structs.Message, error) {
	rows, err := DB.Query("SELECT id, sender_username, receiver_username, content, status, created_at FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", user_id, receiver_id, receiver_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []structs.Message
	for rows.Next() {
		var chat structs.Message
		if err := rows.Scan(&chat.ID, &chat.SenderUsername, &chat.ReceiverUsername, &chat.Content, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
