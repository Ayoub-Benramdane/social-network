package database

import (
	"fmt"
	structs "social-network/data"
	"strings"
)

func CreatePost(user_id, category_id int64, title, content, image, privacy string) (int64, error) {
	fmt.Println(user_id, category_id, title, content, image, privacy)
	result, err := DB.Exec("INSERT INTO posts (title, content, category_id, user_id, image, privacy) VALUES (?, ?, ?, ?, ?, ?)", title, content, category_id, user_id, image, privacy)
	if err != nil {
		return 0, err
	}
	post_id, err := result.LastInsertId()
	return post_id, err
}

func GetPosts(id int64, followers []structs.User) ([]structs.Post, error) {
	var posts []structs.Post
	var userIds []int64
	userIds = append(userIds, id)
	for _, follower := range followers {
		userIds = append(userIds, follower.ID)
	}

	placeholders := make([]string, len(userIds))
	args := make([]interface{}, len(userIds)+4)

	args[0] = "public"
	args[1] = "almost_private"
	args[2] = id
	args[3] = "private"

	for i, uid := range userIds {
		placeholders[i] = "?"
		args[i+4] = uid
	}

	query := fmt.Sprintf(`
	SELECT DISTINCT posts.id, posts.title, posts.content, categories.name, users.username,
	       posts.created_at, posts.total_likes, posts.total_comments, posts.privacy, posts.image
	FROM posts
	JOIN categories ON categories.id = posts.category_id
	JOIN users ON posts.user_id = users.id
	LEFT JOIN post_privacy ON post_privacy.post_id = posts.id
	WHERE posts.privacy = ?
	   OR (posts.privacy = ? AND post_privacy.user_id = ?)
	   OR (posts.privacy = ? AND posts.user_id IN (%s))
	`, strings.Join(placeholders, ","))

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments, &post.Privacy, &post.Image)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			fmt.Println(err)
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

func GetPost(post_id int64) (structs.Post, error) {
	var post structs.Post
	err := DB.QueryRow("SELECT posts.id, posts.title, posts.content, categories.name, posts.image, users.username, posts.created_at, post.total_likes, post.total_comments FROM posts JOIN users ON posts.user_id = users.id JOIN categories ON categories.id = posts.category_id WHERE posts.id = ?", post_id).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Image, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments)
	if err != nil {
		return structs.Post{}, err
	}
	post.Comments, err = GetPostComments(post_id)
	return post, err
}

func GetCountUserPosts(id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", id).Scan(&count)
	return count, err
}
