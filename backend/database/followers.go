package database

import (
	structs "social-network/data"
)

func FollowUser(followerID, followingID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)",
		followerID, followingID,
	)
	return err
}

func UnfollowUser(followerID, followingID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"DELETE FROM followers WHERE follower_id = ? AND followed_id = ?",
		followerID, followingID,
	)
	return err
}

func GetUserFollowers(userID int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar
		 FROM users u
		 JOIN followers f ON u.id = f.follower_id
		 WHERE f.followed_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []structs.User

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	return followers, nil
}

func GetUserFollowerIDs(userID int64) ([]int64, error) {
	rows, err := Database.Query(
		`SELECT u.id
		 FROM users u
		 JOIN followers f ON u.id = f.follower_id
		 WHERE f.followed_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followerIDs []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, id)
	}

	return followerIDs, nil
}

func GetUserFollowing(userID int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar
		 FROM users u
		 JOIN followers f ON u.id = f.followed_id
		 WHERE f.follower_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []structs.User

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		following = append(following, user)
	}

	return following, nil
}

func CountUserFollowing(userID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM followers WHERE follower_id = ?",
		userID,
	).Scan(&count)
	return count, err
}

func GetSuggestedUsers(userID int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT id, username, avatar, lastname, firstname
		 FROM users
		 WHERE id NOT IN (
			 SELECT followed_id FROM followers WHERE follower_id = ?
		 )
		 AND id NOT IN (
			 SELECT recipient_id FROM invitations WHERE invited_id = ?
		 )
		 AND id != ?`,
		userID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.User

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL, &user.LastName, &user.FirstName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetReceivedFollowRequests(userID int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar
		 FROM users u
		 JOIN invitations i ON u.id = i.invited_id
		 WHERE i.recipient_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.User

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetPendingFollowRequests(userID int64) ([]structs.User, error) {
	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar
		 FROM users u
		 JOIN invitations i ON u.id = i.recipient_id
		 WHERE i.group_id = 0 AND i.invited_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.User

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func IsUserFollowing(followerID, followingID int64) (bool, error) {
	var count int
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ?",
		followerID, followingID,
	).Scan(&count)
	return count > 0, err
}

func CountUserFollowers(userID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM followers WHERE followed_id = ?",
		userID,
	).Scan(&count)
	return count, err
}
