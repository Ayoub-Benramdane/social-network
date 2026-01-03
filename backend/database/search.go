package database

import (
	structs "social-network/data"
	"strings"
	"time"
)

func SearchUsers(searchQuery string, offset int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar, u.firstname, u.lastname, u.privacy
		 FROM users u
		 WHERE u.username LIKE ?
		 LIMIT 5 OFFSET ?`,
		searchQuery+"%", offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foundUsers []structs.User

	for rows.Next() {
		var user structs.User
		err = rows.Scan(&user.UserID, &user.Username, &user.AvatarURL, &user.FirstName, &user.LastName, &user.PrivacyLevel)
		if err != nil {
			return nil, err
		}
		foundUsers = append(foundUsers, user)
	}

	return foundUsers, nil
}

func SearchGroups(searchQuery string, offset int64) ([]structs.Group, error) {
	rows, err := Database.Query(
		`SELECT g.id, g.name, g.image
		 FROM groups g
		 WHERE g.name LIKE ?
		 LIMIT 5 OFFSET ?`,
		searchQuery+"%", offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foundGroups []structs.Group

	for rows.Next() {
		var group structs.Group
		err = rows.Scan(&group.GroupID, &group.Name, &group.ImageURL)
		if err != nil {
			return nil, err
		}
		foundGroups = append(foundGroups, group)
	}

	return foundGroups, nil
}

func SearchEvents(searchQuery string, offset int64) ([]structs.Event, error) {
	rows, err := Database.Query(
		`SELECT e.id, e.name, e.start_date, e.end_date, e.image
		 FROM group_events e
		 WHERE e.name LIKE ?
		 LIMIT 5 OFFSET ?`,
		searchQuery+"%", offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var validEvents []structs.Event
	var expiredEventIDs []int64

	for rows.Next() {
		var event structs.Event
		var startDate time.Time

		err = rows.Scan(&event.EventID, &event.Name, &startDate, &event.EndDate, &event.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		if event.EndDate.After(time.Now()) {
			expiredEventIDs = append(expiredEventIDs, event.EventID)
			continue
		}

		validEvents = append(validEvents, event)
	}

	for _, eventID := range expiredEventIDs {
		if err := DeleteEventByID(eventID); err != nil {
			return nil, err
		}
	}

	return validEvents, nil
}

func SearchPosts(userID int64, searchQuery string, offset int64) ([]structs.Post, error) {
	rows, err := Database.Query(
		`SELECT posts.id, posts.title, posts.content,
		        categories.name, categories.color, categories.background,
		        users.id, users.username, users.Avatar,
		        posts.created_at, posts.total_likes, posts.total_comments,
		        posts.privacy, posts.image
		 FROM posts
		 JOIN categories ON categories.id = posts.category_id
		 JOIN users ON posts.user_id = users.id
		 WHERE posts.title LIKE ?
		 LIMIT 5 OFFSET ?`,
		searchQuery+"%", offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foundPosts []structs.Post

	for rows.Next() {
		var post structs.Post
		var createdAt time.Time

		err = rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CategoryName,
			&post.CategoryColor, &post.CategoryBackground, &post.AuthorID, &post.AuthorName,
			&post.AuthorAvatar, &createdAt, &post.LikeCount, &post.CommentCount, &post.PrivacyLevel, &post.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		post.CreatedAt = TimeAgo(createdAt)

		post.IsLiked, err = IsPostLikedByUser(post.PostID, userID)
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

		post.IsSaved, err = IsSaved(userID, post.PostID)
		if err != nil {
			return nil, err
		}

		foundPosts = append(foundPosts, post)
	}

	return foundPosts, nil
}

func InsertSearch(userID int64, searchQuery string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO searches (user_id, content) VALUES (?, ?)",
		userID, searchQuery,
	)
	return err
}

func UpdateFirstSearch(searchID int64, searchQuery string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"UPDATE searches SET content = ? WHERE id = ?",
		searchQuery, searchID,
	)
	return err
}

func GetCountSearchUser(userID int64) (int64, error) {
	var searchCount int64

	err := Database.QueryRow(
		"SELECT COUNT(*) FROM searches WHERE user_id = ?",
		userID,
	).Scan(&searchCount)

	return searchCount, err
}

func GetIDFirstSearch(userID int64) (int64, error) {
	var searchID int64

	err := Database.QueryRow(
		"SELECT id FROM searches WHERE user_id = ? ORDER BY created_at DESC",
		userID,
	).Scan(&searchID)

	return searchID, err
}
