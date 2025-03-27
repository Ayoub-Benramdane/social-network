package database

import (
	structs "social-network/backend/data"
	"time"
)

func CreatePost(user_id int64, title, content, image, privacy string) (int64, error) {
	result, err := DB.Exec("INSERT INTO posts (title, content, user_id, created_at, privacy) VALUES (?, ?, ?, ?, ?)", title, content, user_id, time.Now(), privacy)
	if err != nil {
		return 0, err
	}
	post_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if image != "" {
		_, err = DB.Exec("INSERT INTO images (post_id, image) VALUES (?, ?)", post_id, image)
		if err != nil {
			return 0, err
		}
	}
	return post_id, nil
}

func GetPosts(id int64, followers []structs.User) ([]structs.Post, error) {
	var posts []structs.Post
	for _, follower := range followers {
		rows, err := DB.Query("SELECT DISTINCT posts.id, posts.title, posts.content, users.username, posts.created_at, post.total_likes, post.total_comments FROM posts JOIN users ON posts.user_id = users.id WHERE posts.user_id = ? OR post.privacy = ? OR (post.privacy = ? AND post.user_id = ?)", id, "public", "public", follower.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post structs.Post
			err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments)
			if err != nil {
				return nil, err
			}
			post.Comments, err = GetPostComments(post.ID)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func GetCountUserPosts(id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", id).Scan(&count)
	return count, err
}
