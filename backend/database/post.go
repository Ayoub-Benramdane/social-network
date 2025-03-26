package database


func CreatePost(title, content, author string, category []string) error {
	_, err := DB.Exec("INSERT INTO posts (title, content, author, category) VALUES (?, ?, ?, ?)", title, content, author, category)
	return err
}