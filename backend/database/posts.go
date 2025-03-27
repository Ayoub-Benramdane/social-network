package database

import (
	structs "social-network/backend/data"
	"time"
)

func CreatePost(user_id int64, title, content string, categories []string) (int64, error) {
	result, err := DB.Exec("INSERT INTO posts (title, content, user_id, created_at) VALUES (?, ?, ?, ?)", title, content, user_id, time.Now())
	if err != nil {
		return 0, err
	}
	post_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	var category_id int64
	for _, category := range categories {
		err = DB.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&category_id)
		if err != nil {
			return 0, err
		}
		_, err = DB.Exec("INSERT INTO post_category (category_id, post_id) VALUES (?, ?)", category_id, post_id)
		if err != nil {
			return 0, err
		}
	}
	return post_id, nil
}

func GetPosts(id int64) ([]structs.Post, error) {
	rows, err := DB.Query("SELECT p.id, p.title, p.content, u.username, p.created_at, p.total_likes, p.total_comments FROM posts p JOIN users u WHERE u.id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post

	for rows.Next() {
		var post structs.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments)
		if err != nil {
			return nil, err
		}
		post.Categories, err = GetPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments, err = GetPostComments(post.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetCountUserPosts(id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", id).Scan(&count)
	return count, err
}