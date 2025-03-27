package database

func GetCountUserLikes(user_id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ?", user_id).Scan(&count)
	return count, err
}