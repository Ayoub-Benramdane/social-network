package database

import (
	structs "social-network/data"
	"time"
)

func CreateComment(content string, user_id int64, post structs.Post) (int64, error) {
	result, err := DB.Exec("INSERT INTO comments (content, user_id, post_id) VALUES (?, ?, ?)", content, user_id, post.ID)
	if err != nil {
		return 0, err
	}

	if err = CreateNotification(user_id, post.ID, post.UserID, "comment"); err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	return lastID, err
}

func CreateGroupComment(content string, user_id, group_id int64, post structs.Post) (int64, error) {
	result, err := DB.Exec("INSERT INTO group_comments (content, user_id, post_id, group_id) VALUES (?, ?, ?, ?)", content, user_id, post.ID, group_id)
	if err != nil {
		return 0, err
	}

	if err = CreateNotification(user_id, post.ID, post.UserID, "comment"); err != nil {
		return 0, err
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
