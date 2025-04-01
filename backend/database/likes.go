package database

import structs "social-network/data"

func LikePost(user_id int64, post structs.Post) (int64, error) {
	var err error
	var notification bool
	if !CheckLike(user_id, post.ID) {
		notification = true
		_, err = DB.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", user_id, post.ID)
	} else {
		_, err = DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", user_id, post.ID)
	}
	if err != nil {
		return 0, err
	}

	if notification {
		if err = CreateNotification(user_id, post.ID, post.UserID, "like"); err != nil {
			return 0, err
		}
	} else if err = DeleteNotification(user_id, post.ID, post.UserID, "like"); err != nil {
		return 0, err
	}

	var count int64
	err = DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post.ID).Scan(&count)
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
