package database

import (
	structs "social-network/data"
)

func CreatePost(user_id int64, title, content, category, image, privacy string) (int64, error) {
	result, err := DB.Exec("INSERT INTO posts (title, content, category, user_id, image, privacy) VALUES (?, ?, ?, ?, ?, ?)", title, content, category, user_id, image, privacy)
	if err != nil {
		return 0, err
	}
	post_id, err := result.LastInsertId()
	return post_id, nil
}

func GetPosts(id int64, followers []structs.User) ([]structs.Post, error) {
	var posts []structs.Post
	for _, follower := range followers {
		rows, err := DB.Query("SELECT DISTINCT posts.id, posts.title, posts.content, posts.category, posts.image, users.username, posts.created_at, posts.total_likes, posts.total_comments FROM posts JOIN users ON posts.user_id = users.id WHERE posts.user_id = ? OR posts.privacy = ? OR (posts.privacy = ? AND posts.user_id = ?)", id, "public", "public", follower.ID)

		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post structs.Post
			err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Image, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments)
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

func GetPost(post_id int64) (structs.Post, error) {
	var post structs.Post
	err := DB.QueryRow("SELECT posts.id, posts.title, posts.content, posts.category, posts.image, users.username, posts.created_at, post.total_likes, post.total_comments FROM posts JOIN users ON posts.user_id = users.id WHERE posts.id = ?", post_id).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Image, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments)
	if err != nil {
		return structs.Post{}, err
	}
	post.Comments, err = GetPostComments(post_id)
	return post, nil
}

func GetCountUserPosts(id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", id).Scan(&count)
	return count, err
}
