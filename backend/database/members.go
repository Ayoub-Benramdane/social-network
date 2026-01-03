package database

import (
	structs "social-network/data"
)

func AddUserToGroup(userID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	if _, err := Database.Exec(
		"INSERT INTO group_members (user_id, group_id) VALUES (?, ?)",
		userID, groupID,
	); err != nil {
		return err
	}

	_, err := Database.Exec(
		"UPDATE groups SET members = members + 1 WHERE id = ?",
		groupID,
	)
	return err
}

func RemoveUserFromGroup(userID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	if _, err := Database.Exec(
		"DELETE FROM group_members WHERE user_id = ? AND group_id = ?",
		userID, groupID,
	); err != nil {
		return err
	}

	if _, err := Database.Exec(
		"DELETE FROM invitations WHERE invited_id = ? AND group_id = ?",
		userID, groupID,
	); err != nil {
		return err
	}

	_, err := Database.Exec(
		"UPDATE groups SET members = members - 1 WHERE id = ?",
		groupID,
	)
	return err
}

func IsUserGroupMember(userID, groupID int64) (bool, error) {
	var memberCount int
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID,
	).Scan(&memberCount)

	return memberCount > 0, err
}

func FetchGroupMembers(currentUserID, groupID int64) ([]structs.User, error) {
	var members []structs.User

	rows, err := Database.Query(
		`SELECT u.id, u.username, u.avatar, u.lastname, u.firstname
		 FROM users u
		 JOIN group_members gm ON u.id = gm.user_id
		 WHERE gm.group_id = ?`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member structs.User
		if err := rows.Scan(
			&member.UserID, &member.Username, &member.AvatarURL, &member.LastName, &member.FirstName); err != nil {
			return nil, err
		}

		if member.UserID == currentUserID {
			member.Role = "admin"
		} else {
			member.Role = "member"
		}

		members = append(members, member)
	}

	return members, nil
}

func GetInvitableUsers(userID, groupID int64) ([]structs.User, error) {
	var users []structs.User

	rows, err := Database.Query(
		`SELECT DISTINCT u.id, u.username, u.avatar
		 FROM users u
		 JOIN followers f ON u.id = f.follower_id OR u.id = f.followed_id
		 WHERE (f.follower_id = ? OR f.followed_id = ?)
		   AND u.id NOT IN (
			   SELECT user_id FROM group_members WHERE group_id = ?
			   UNION
			   SELECT recipient_id FROM invitations WHERE group_id = ?
		   )`,
		userID, userID, groupID, groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
