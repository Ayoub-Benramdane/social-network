package database

import (
	"strings"
	"time"

	structs "social-network/data"
)

func CreatePostComment(content string, userID int64, post structs.Post, imageURL string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	result, err := Database.Exec(
		"INSERT INTO comments (content, user_id, post_id, image) VALUES (?, ?, ?, ?)",
		content, userID, post.PostID, imageURL,
	)
	if err != nil {
		return 0, err
	}

	_, err = Database.Exec(
		"UPDATE posts SET total_comments = total_comments + 1 WHERE id = ?",
		post.PostID,
	)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	return lastInsertID, err
}

func FetchPostComments(postID int64) ([]structs.Comment, error) {
	rows, err := Database.Query(
		`SELECT c.id, c.content, u.id, u.username, u.avatar, c.created_at, c.image
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.post_id = ?
		 ORDER BY c.created_at DESC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []structs.Comment

	for rows.Next() {
		var comment structs.Comment
		var createdAt time.Time

		err = rows.Scan(&comment.CommentID, &comment.Content, &comment.AuthorID, &comment.Username, &comment.AvatarURL, &createdAt, &comment.ImageURL,)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		comment.CreatedAt = TimeAgo(createdAt)
		comments = append(comments, comment)
	}

	return comments, nil
}

func CountUserComments(userID int64) (int64, error) {
	var total int64

	err := Database.QueryRow(
		"SELECT COUNT(*) FROM comments WHERE user_id = ?",
		userID,
	).Scan(&total)

	return total, err
}
