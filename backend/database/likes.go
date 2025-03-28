package database

func LikePost(user_id, post_id int64) (int64, error) {
	var err error
	if !CheckLike(user_id, post_id) {
		_, err = DB.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", user_id, post_id)
	} else {
		_, err = DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", user_id, post_id)
	}
	if err != nil {
		return 0, err
	}
	var count int64
	err = DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post_id).Scan(&count)
	return count, err
}

func CheckLike(user_id, post_id int64) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", user_id, post_id).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func GetCountUserLikes(user_id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ?", user_id).Scan(&count)
	return count, err
}