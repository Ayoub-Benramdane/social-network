package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	structs "social-network/data"
)

func CreatePost(authorID, groupID, categoryID int64, title, content, imageURL, privacyLevel string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	result, err := Database.Exec(
		"INSERT INTO posts (title, content, category_id, user_id, group_id, image, privacy) VALUES (?, ?, ?, ?, ?, ?, ?)",
		title, content, categoryID, authorID, groupID, imageURL, privacyLevel,
	)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	return postID, err
}

func GetPosts(currentUserID, offset int64, followedUsers []structs.User) ([]structs.Post, error) {
	var posts []structs.Post
	var allowedUserIDs []int64

	allowedUserIDs = append(allowedUserIDs, currentUserID)
	for _, user := range followedUsers {
		allowedUserIDs = append(allowedUserIDs, user.UserID)
	}

	placeholders := make([]string, len(allowedUserIDs))
	args := make([]interface{}, len(allowedUserIDs)+8)

	args[0] = "public"
	args[1] = currentUserID
	args[2] = "private"
	args[3] = currentUserID
	args[4] = "almost_private"

	for i, userID := range allowedUserIDs {
		placeholders[i] = "?"
		args[i+5] = userID
	}

	args[len(args)-3] = int64(0)
	args[len(args)-2] = int64(10)
	args[len(args)-1] = offset

	query := fmt.Sprintf(`
		SELECT DISTINCT posts.id, posts.title, posts.content,
		       categories.name, categories.color, categories.background,
		       users.id, users.username, users.avatar,
		       posts.created_at, posts.total_likes, posts.total_comments,
		       posts.privacy, posts.image
		FROM posts
		JOIN categories ON categories.id = posts.category_id
		JOIN users ON posts.user_id = users.id
		LEFT JOIN post_privacy ON post_privacy.post_id = posts.id
		WHERE (
			posts.privacy = ?
			OR posts.user_id = ?
			OR (posts.privacy = ? AND post_privacy.user_id = ?)
			OR (posts.privacy = ? AND posts.user_id IN (%s))
		)
		AND posts.group_id = ?
		ORDER BY posts.created_at DESC
		LIMIT ? OFFSET ?
	`, strings.Join(placeholders, ","))

	rows, err := Database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Post
		var createdAt time.Time

		err = rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CategoryName, &post.CategoryColor,
			&post.CategoryBackground, &post.AuthorID, &post.AuthorName, &post.AuthorAvatar, &createdAt,
			&post.LikeCount, &post.CommentCount, &post.PrivacyLevel, &post.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		post.CreatedAt = TimeAgo(createdAt)

		post.IsLiked, err = IsPostLikedByUser(post.PostID, currentUserID)
		if err != nil {
			return nil, err
		}

		post.LikedBy, err = GetPostLikedUsers(post.PostID)
		if err != nil {
			return nil, err
		}

		post.SaveCount, err = CountSaves(post.PostID, 0)
		if err != nil {
			return nil, err
		}

		post.IsSaved, err = IsSaved(currentUserID, post.PostID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsByUser(targetUserID, viewerUserID int64, isFollowed bool) ([]structs.Post, error) {
	var posts []structs.Post
	var rows *sql.Rows
	var err error

	if targetUserID == viewerUserID || isFollowed {
		rows, err = Database.Query(`
			SELECT DISTINCT posts.id, posts.title, posts.content,
			       categories.name, categories.color, categories.background,
			       users.username, users.avatar,
			       posts.created_at, posts.total_likes, posts.total_comments,
			       posts.privacy, posts.image
			FROM posts
			JOIN categories ON categories.id = posts.category_id
			JOIN users ON posts.user_id = users.id
			LEFT JOIN post_privacy ON post_privacy.post_id = posts.id
			WHERE posts.user_id = ?
			  AND posts.group_id = 0
			  AND (
			       posts.privacy = ?
			    OR posts.privacy = ?
			    OR (posts.privacy = ? AND post_privacy.user_id = ?)
			  )
			ORDER BY posts.created_at DESC
		`, targetUserID, "public", "almost_private", "private", viewerUserID)
	} else {
		rows, err = Database.Query(`
			SELECT DISTINCT posts.id, posts.title, posts.content,
			       categories.name, categories.color, categories.background,
			       users.username, users.avatar,
			       posts.created_at, posts.total_likes, posts.total_comments,
			       posts.privacy, posts.image
			FROM posts
			JOIN categories ON categories.id = posts.category_id
			JOIN users ON posts.user_id = users.id
			WHERE posts.user_id = ?
			  AND posts.group_id = 0
			  AND posts.privacy = ?
			ORDER BY posts.created_at DESC
		`, targetUserID, "public")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Post
		var createdAt time.Time

		err = rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CategoryName,
			&post.CategoryColor, &post.CategoryBackground, &post.AuthorName, &post.AuthorAvatar,
			&createdAt, &post.LikeCount, &post.CommentCount, &post.PrivacyLevel, &post.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		post.CreatedAt = TimeAgo(createdAt)

		post.IsLiked, err = IsPostLikedByUser(post.PostID, viewerUserID)
		if err != nil {
			return nil, err
		}

		post.LikedBy, err = GetPostLikedUsers(post.PostID)
		if err != nil {
			return nil, err
		}

		post.SaveCount, err = CountSaves(post.PostID, post.GroupID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsGroup(groupID, currentUserID int64, privacyLevel string) ([]structs.Post, error) {
	var posts []structs.Post

	rows, err := Database.Query(`
		SELECT p.id, p.title, p.content,
		       categories.name, categories.color, categories.background,
		       users.username, users.avatar,
		       p.created_at, p.total_likes, p.total_comments, p.image
		FROM posts p
		JOIN categories ON categories.id = p.category_id
		JOIN users ON p.user_id = users.id
		WHERE p.group_id = ?
		ORDER BY p.created_at DESC
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Post
		var createdAt time.Time

		err = rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CategoryName,
			&post.CategoryColor, &post.CategoryBackground, &post.AuthorName, &post.AuthorAvatar,
			&createdAt, &post.LikeCount, &post.CommentCount, &post.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		post.CreatedAt = TimeAgo(createdAt)
		post.PrivacyLevel = privacyLevel

		post.IsLiked, err = IsPostLikedByUser(post.PostID, currentUserID)
		if err != nil {
			return nil, err
		}

		post.LikedBy, err = GetPostLikedUsers(post.PostID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsByCategory(categoryID, currentUserID int64) ([]structs.Post, error) {
	var posts []structs.Post

	rows, err := Database.Query(`
		SELECT p.id, p.title, p.content,
		       categories.name, categories.color, categories.background,
		       users.username, users.avatar,
		       p.created_at, p.total_likes, p.total_comments, p.image
		FROM posts p
		JOIN categories ON categories.id = p.category_id
		JOIN users ON p.user_id = users.id
		WHERE p.category_id = ?
		ORDER BY p.created_at DESC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Post
		var createdAt time.Time

		err = rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CategoryName, &post.CategoryColor,
			&post.CategoryBackground, &post.AuthorName, &post.AuthorAvatar, &createdAt, &post.LikeCount, &post.CommentCount, &post.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		post.CreatedAt = TimeAgo(createdAt)

		post.IsLiked, err = IsPostLikedByUser(post.PostID, currentUserID)
		if err != nil {
			return nil, err
		}

		post.LikedBy, err = GetPostLikedUsers(post.PostID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPost(currentUserID, postID int64) (structs.Post, error) {
	var post structs.Post
	var createdAt time.Time

	err := Database.QueryRow(`
		SELECT posts.id, posts.title, posts.content,
		       categories.name, categories.color, categories.background,
		       users.id, users.username, users.avatar,
		       posts.created_at, posts.total_likes, posts.total_comments,
		       posts.privacy, posts.image
		FROM posts
		JOIN users ON posts.user_id = users.id
		JOIN categories ON categories.id = posts.category_id
		WHERE posts.id = ?
	`, postID).Scan(&post.PostID, &post.Title, &post.Content, &post.CategoryName,
		&post.CategoryColor, &post.CategoryBackground, &post.AuthorID, &post.AuthorName,
		&post.AuthorAvatar, &createdAt, &post.LikeCount, &post.CommentCount, &post.PrivacyLevel, &post.ImageURL)
	if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
		return structs.Post{}, err
	}

	post.CreatedAt = TimeAgo(createdAt)

	post.IsLiked, err = IsPostLikedByUser(post.PostID, currentUserID)
	if err != nil {
		return structs.Post{}, err
	}

	post.IsSaved, err = IsSaved(currentUserID, post.PostID)
	if err != nil {
		return structs.Post{}, err
	}

	post.LikedBy, err = GetPostLikedUsers(post.PostID)
	return post, err
}

func GetCountUserPosts(authorID, groupID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM posts WHERE user_id = ? AND group_id = ?",
		authorID, groupID,
	).Scan(&count)
	return count, err
}

func GetCountGroupPosts(groupID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM posts WHERE group_id = ?",
		groupID,
	).Scan(&count)
	return count, err
}

func IsAuthorized(userID, postID int64) (bool, error) {
	var count int
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM post_privacy WHERE post_id = ? AND user_id = ?",
		postID, userID,
	).Scan(&count)
	return count > 0, err
}

func AddAlmostPrivateUser(userID, postID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO post_privacy (user_id, post_id) VALUES (?, ?)",
		userID, postID,
	)
	return err
}

func GetLastTime(tableName string) (string, error) {
	var lastTime string
	query := fmt.Sprintf(
		"SELECT created_at FROM %s ORDER BY created_at DESC LIMIT 1",
		tableName,
	)

	err := Database.QueryRow(query).Scan(&lastTime)
	return lastTime, err
}
