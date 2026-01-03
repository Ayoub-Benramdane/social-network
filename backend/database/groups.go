package database

import (
	"strings"
	"time"

	structs "social-network/data"
)

func CreateGroup(adminID int64, name, description, imageURL, coverURL, privacy string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	result, err := Database.Exec(
		"INSERT INTO groups (name, description, image, cover, admin, privacy) VALUES (?, ?, ?, ?, ?, ?)",
		name, description, imageURL, coverURL, adminID, privacy,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func DeleteGroup(groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	if _, err := Database.Exec("DELETE FROM groups WHERE id = ?", groupID); err != nil {
		return err
	}
	if _, err := Database.Exec("DELETE FROM posts WHERE group_id = ?", groupID); err != nil {
		return err
	}
	_, err := Database.Exec("DELETE FROM invitations WHERE group_id = ?", groupID)
	return err
}

func GetUserGroups(user structs.User) ([]structs.Group, error) {
	rows, err := Database.Query(
		`SELECT g.id, g.name, g.description, g.cover, g.created_at,
		        g.admin, g.privacy, u.username, g.members, g.image
		 FROM groups g
		 LEFT JOIN group_members m ON g.id = m.group_id
		 LEFT JOIN users u ON u.id = g.admin
		 LEFT JOIN messages ms ON (u.id = ms.sender_id OR u.id = ms.receiver_id)
			AND ms.group_id = g.id
		 WHERE m.user_id = ?
		 GROUP BY g.id
		 ORDER BY MAX(ms.created_at) DESC`,
		user.UserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []structs.Group

	for rows.Next() {
		var group structs.Group
		var createdAt time.Time

		err = rows.Scan(&group.GroupID, &group.Name, &group.Description, &group.CoverURL, &createdAt, &group.AdminID, &group.PrivacyLevel, &group.AdminName, &group.MemberCount, &group.ImageURL)
		if err != nil && !strings.Contains(err.Error(), `name "image": converting NULL to string`) {
			return nil, err
		}

		group.MessageCount, err = CountConversationUnreadMessages(0, user.UserID, group.GroupID)
		if err != nil {
			return nil, err
		}

		group.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		group.PostCount, err = GetCountGroupPosts(group.GroupID)
		if err != nil {
			return nil, err
		}

		if user.Username == group.AdminName {
			group.UserRole = "admin"
		} else {
			group.UserRole = "member"
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func GetSuggestedGroups(userID int64) ([]structs.Group, error) {
	rows, err := Database.Query(
		`SELECT g.id, g.name, g.description, g.image, g.cover,
		        g.admin, g.privacy, g.created_at, u.username, g.members
		 FROM groups g
		 JOIN users u ON u.id = g.admin
		 WHERE g.id NOT IN (
			 SELECT group_id FROM group_members WHERE user_id = ?
			 UNION
			 SELECT group_id FROM invitations WHERE invited_id = ?
		 )
		 ORDER BY g.created_at DESC`,
		userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []structs.Group

	for rows.Next() {
		var group structs.Group
		var createdAt time.Time

		if err := rows.Scan(&group.GroupID, &group.Name, &group.Description, &group.ImageURL, &group.CoverURL, &group.AdminID, &group.PrivacyLevel, &createdAt, &group.AdminName, &group.MemberCount); err != nil {
			return nil, err
		}

		group.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		group.PostCount, err = GetCountGroupPosts(group.GroupID)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func GetPendingGroups(userID int64) ([]structs.Group, error) {
	rows, err := Database.Query(
		`SELECT g.id, i.recipient_id, g.name, g.description,
		        g.image, g.cover, g.admin, g.privacy,
		        g.created_at, u.username, g.members
		 FROM groups g
		 JOIN users u ON u.id = g.admin
		 JOIN invitations i ON g.id = i.group_id
		 WHERE i.invited_id = ?
		 ORDER BY g.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []structs.Group

	for rows.Next() {
		var group structs.Group
		var createdAt time.Time
		var recipientID int64

		if err := rows.Scan(&group.GroupID, &recipientID, &group.Name, &group.Description, &group.ImageURL, &group.CoverURL, &group.AdminID, &group.PrivacyLevel, &createdAt, &group.AdminName, &group.MemberCount); err != nil {
			return nil, err
		}

		isAdmin, err := IsGroupAdmin(recipientID, group.GroupID)
		if err != nil || !isAdmin {
			continue
		}

		group.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		group.PostCount, err = GetCountGroupPosts(group.GroupID)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func IsGroupAdmin(userID, groupID int64) (bool, error) {
	var isAdmin bool
	err := Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM groups WHERE id = ? AND admin = ?)",
		groupID, userID,
	).Scan(&isAdmin)

	return isAdmin, err
}

func GetGroupByID(groupID int64) (structs.Group, error) {
	var group structs.Group
	var createdAt time.Time

	err := Database.QueryRow(
		`SELECT g.id, g.name, g.description, g.image, g.cover,
		        g.admin, g.created_at, u.username, g.members, g.privacy
		 FROM groups g
		 JOIN users u ON u.id = g.admin
		 WHERE g.id = ?`,
		groupID,
	).Scan(&group.GroupID, &group.Name, &group.Description, &group.ImageURL, &group.CoverURL, &group.AdminID, &createdAt, &group.AdminName, &group.MemberCount, &group.PrivacyLevel)
	if err != nil {
		return group, err
	}

	group.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	group.PostCount, err = GetCountUserPosts(group.AdminID, groupID)

	return group, err
}

func CountUserGroups(userID int64) (int64, error) {
	var count int64
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM group_members WHERE user_id = ?",
		userID,
	).Scan(&count)
	return count, err
}

func GetGroupMemberIDs(groupID int64) ([]int64, error) {
	rows, err := Database.Query(
		"SELECT u.id FROM group_members m JOIN users u ON u.id = m.user_id WHERE m.group_id = ?",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberIDs []int64

	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, userID)
	}

	return memberIDs, nil
}
