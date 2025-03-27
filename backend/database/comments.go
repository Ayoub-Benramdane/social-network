package database

import (
	structs "social-network/backend/data"
	"time"
)

func CreateComment(content string, user_id, post_id int64) (int64, error) {
	result, errInsert := DB.Exec("INSERT INTO comments (content, user_id, post_id, created_at) VALUES (?, ?, ?, ?)", content, user_id, post_id, time.Now())
	if errInsert != nil {
		return 0, errInsert
	}
	lastID, err := result.LastInsertId()
	return lastID, err
}

func GetPostComments(post_id int64) ([]structs.Comment, error) {
	rows, err := DB.Query("SELECT c.id, c.content, u.username, c.created_at FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = ?", post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		var date time.Time
		err = rows.Scan(&comment.ID, &comment.Content, &comment.Author, &date)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = TimeAgo(date)
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetCountUserComments(user_id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ?", user_id).Scan(&count)
	return count, err
}