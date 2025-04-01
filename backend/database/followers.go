package database

import structs "social-network/data"

func AddFollower(follower_id, following_id int64) error {
	_, err := DB.Exec("INSERT INTO followers (follower_id, following_id) VALUES (?, ?)", follower_id, following_id)
	if err == nil {
		if err = CreateNotification(follower_id, 0, following_id, "follow"); err != nil {
			return err
		}
	}
	return err
}

func RemoveFollower(follower_id, following_id int64) error {
	_, err := DB.Exec("DELETE FROM followers WHERE follower_id = ? AND following_id = ?", follower_id, following_id)
	return err
}

func GetFollowers(user_id int64) ([]structs.User, error) {
	rows, err := DB.Query("SELECT u.id, u.username FROM users u JOIN followers f ON u.id = f.follower_id WHERE f.following_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var followers []structs.User
	for rows.Next() {
		var follower structs.User
		err = rows.Scan(&follower.ID, &follower.Username)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return followers, nil
}

func GetFollowing(user_id int64) ([]structs.User, error) {
	rows, err := DB.Query("SELECT u.id, u.username FROM users u JOIN followers f ON u.id = f.following_id WHERE f.follower_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var following []structs.User
	for rows.Next() {
		var follower structs.User
		err = rows.Scan(&follower.ID, &follower.Username)
		if err != nil {
			return nil, err
		}
		following = append(following, follower)
	}
	return following, nil
}

func GetNotFollowing(user_id int64) ([]structs.User, error) {
	rows, err := DB.Query("SELECT id, username FROM users WHERE id NOT IN (SELECT follower_id FROM followers WHERE following_id = ?) AND id != ?", user_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notFollowing []structs.User
	for rows.Next() {
		var user structs.User
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		notFollowing = append(notFollowing, user)
	}
	return notFollowing, nil
}