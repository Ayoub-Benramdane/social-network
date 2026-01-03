package database

import (
	structs "social-network/data"
)

func TogglePostLike(userID int64, post structs.Post) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	if post.IsLiked {
		_, err := Database.Exec(
			"DELETE FROM post_likes WHERE user_id = ? AND post_id = ?",
			userID, post.PostID,
		)
		if err != nil {
			return 0, err
		}
	} else {
		_, err := Database.Exec(
			"INSERT INTO post_likes (user_id, post_id) VALUES (?, ?)",
			userID, post.PostID,
		)
		if err != nil {
			return 0, err
		}
	}

	var totalLikes int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM post_likes WHERE post_id = ?",
		post.PostID,
	).Scan(&totalLikes)
	if err != nil {
		return 0, err
	}

	_, err = Database.Exec(
		"UPDATE posts SET total_likes = ? WHERE id = ?",
		totalLikes, post.PostID,
	)

	return totalLikes, err
}

func IsPostLikedByUser(postID, userID int64) (bool, error) {
	var likeCount int
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?",
		postID, userID,
	).Scan(&likeCount)

	return likeCount > 0, err
}

func CountLikesByUser(userID int64) (int64, error) {
	var totalLikes int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM post_likes WHERE user_id = ?",
		userID,
	).Scan(&totalLikes)

	return totalLikes, err
}

func GetPostLikedUsers(postID int64) ([]structs.User, error) {
	var likedUsers []structs.User

	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar
		 FROM users u
		 JOIN post_likes pl ON u.id = pl.user_id
		 WHERE pl.post_id = ?`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		likedUsers = append(likedUsers, user)
	}

	return likedUsers, nil
}
