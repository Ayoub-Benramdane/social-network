package database

import (
	structs "social-network/data"
	"strings"
	"time"
)

func SavePost(userID, postID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO saves (user_id, post_id, group_id) VALUES (?, ?, ?)",
		userID, postID, groupID,
	)
	return err
}

func UnsavePost(userID, postID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"DELETE FROM saves WHERE user_id = ? AND post_id = ? AND group_id = ?",
		userID, postID, groupID,
	)
	return err
}

func GetSavedPosts(userID, requestedGroupID int64) ([]structs.Post, error) {
	rows, err := Database.Query(`
		SELECT p.id, p.group_id,
		       u.username, u.avatar,
		       p.title, p.content,
		       c.name, c.color, c.background,
		       p.created_at,
		       p.total_likes, p.total_comments,
		       p.privacy, p.image
		FROM saves s
		JOIN posts p ON s.post_id = p.id
		JOIN categories c ON c.id = p.category_id
		JOIN users u ON u.id = p.user_id
		WHERE s.user_id = ?
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var savedPosts []structs.Post

	for rows.Next() {
		var post structs.Post
		var createdAt time.Time

		err := rows.Scan(&post.PostID, &post.GroupID, &post.AuthorName, &post.AuthorAvatar,
			&post.Title, &post.Content, &post.CategoryName, &post.CategoryColor, &post.CategoryBackground,
			&createdAt, &post.LikeCount, &post.CommentCount, &post.PrivacyLevel, &post.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		post.CreatedAt = TimeAgo(createdAt)

		post.IsLiked, err = IsPostLikedByUser(post.PostID, userID)
		if err != nil {
			return nil, err
		}

		if requestedGroupID != 0 && post.GroupID != 0 {
			savedPosts = append(savedPosts, post)
		} else if requestedGroupID == 0 && post.GroupID == 0 {
			savedPosts = append(savedPosts, post)
		}
	}

	return savedPosts, nil
}

func IsSaved(userID, postID int64) (bool, error) {
	var saveCount int

	err := Database.QueryRow(
		"SELECT COUNT(*) FROM saves WHERE user_id = ? AND post_id = ?",
		userID, postID,
	).Scan(&saveCount)

	return saveCount > 0, err
}

func CountSaves(postID, groupID int64) (int64, error) {
	var actualCount int64
	var storedCount int64

	err := Database.QueryRow(
		"SELECT COUNT(*) FROM saves WHERE post_id = ? AND group_id = ?",
		postID, groupID,
	).Scan(&actualCount)
	if err != nil {
		return 0, err
	}

	err = Database.QueryRow(
		"SELECT total_saves FROM posts WHERE id = ?",
		postID,
	).Scan(&storedCount)
	if err != nil {
		return 0, err
	}

	if actualCount != storedCount {
		mu.Lock()
		defer mu.Unlock()

		_, err = Database.Exec(
			"UPDATE posts SET total_saves = ? WHERE id = ?",
			actualCount, postID,
		)
	}

	return actualCount, err
}
