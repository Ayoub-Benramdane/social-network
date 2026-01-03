package database

import (
	"time"

	structs "social-network/data"
)

func CreateInvitation(invitedID, recipientID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO invitations (recipient_id, invited_id, group_id) VALUES (?, ?, ?)",
		recipientID, invitedID, groupID,
	)
	return err
}

func AcceptInvitation(invitationID, invitedID, recipientID, groupID int64) error {
	mu.Lock()
	defer mu.Unlock()

	var err error
	if groupID != 0 {
		_, err = Database.Exec(
			"INSERT INTO group_members (user_id, group_id) VALUES (?, ?)",
			recipientID, groupID,
		)
	} else {
		_, err = Database.Exec(
			"INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)",
			invitedID, recipientID,
		)
	}
	if err != nil {
		return err
	}

	return DeleteInvitation(invitationID)
}

func DeleteInvitation(invitationID int64) error {
	_, err := Database.Exec(
		"DELETE FROM invitations WHERE id = ?",
		invitationID,
	)
	return err
}

func GetFriendInvitations(userID int64) ([]structs.Invitation, error) {
	rows, err := Database.Query(
		`SELECT i.id, u.id, u.username, u.avatar
		 FROM invitations i
		 JOIN users u ON i.invited_id = u.id
		 WHERE i.recipient_id = ?
		 ORDER BY i.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []structs.Invitation

	for rows.Next() {
		var invitation structs.Invitation
		if err := rows.Scan(&invitation.InvitationID, &invitation.User.UserID, &invitation.User.Username, &invitation.User.AvatarURL); err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func GetGroupInvitations(userID int64) ([]structs.Invitation, error) {
	rows, err := Database.Query(
		`SELECT i.id, i.created_at,
		        u.id, u.username, u.avatar,
		        g.id, g.admin, g.name, g.members
		 FROM invitations i
		 JOIN users u ON u.id = i.invited_id
		 JOIN groups g ON i.group_id = g.id
		 WHERE i.recipient_id = ?
		 ORDER BY i.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []structs.Invitation

	for rows.Next() {
		var invitation structs.Invitation
		var createdAt time.Time

		if err := rows.Scan(&invitation.InvitationID, &createdAt, &invitation.User.UserID, &invitation.User.Username, &invitation.User.AvatarURL,
			&invitation.Group.GroupID, &invitation.Group.AdminID, &invitation.Group.Name, &invitation.Group.MemberCount); err != nil {
			return nil, err
		}

		if invitation.Group.AdminID == userID || invitation.Group.AdminID == invitation.User.UserID {
			invitation.IsOwner = true
		}

		var err error
		invitation.Group.AdminName, err = GetAdminUsername(invitation.Group.AdminID)
		if err != nil {
			return nil, err
		}

		invitation.CreatedAt = TimeAgo(createdAt)
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func GetGroupInvitationsByGroup(userID, groupID int64) ([]structs.Invitation, error) {
	rows, err := Database.Query(
		`SELECT i.id, i.created_at, u.id, u.username, u.avatar
		 FROM invitations i
		 JOIN users u ON u.id = i.invited_id
		 JOIN groups g ON i.group_id = g.id
		 WHERE i.recipient_id = ? AND g.id = ?
		 ORDER BY i.created_at DESC`,
		userID, groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []structs.Invitation

	for rows.Next() {
		var invitation structs.Invitation
		var createdAt time.Time

		if err := rows.Scan(
			&invitation.InvitationID,
			&createdAt,
			&invitation.User.UserID,
			&invitation.User.Username,
			&invitation.User.AvatarURL,
		); err != nil {
			return nil, err
		}

		invitation.IsOwner = true
		invitation.CreatedAt = TimeAgo(createdAt)
		invitation.Group.GroupID = groupID

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func GetAdminUsername(userID int64) (string, error) {
	var username string
	err := Database.QueryRow(
		"SELECT username FROM users WHERE id = ?",
		userID,
	).Scan(&username)
	return username, err
}

func InvitationExists(invitedID, recipientID, groupID int64) (bool, error) {
	var count int
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM invitations WHERE recipient_id = ? AND invited_id = ? AND group_id = ?",
		recipientID, invitedID, groupID,
	).Scan(&count)
	return count > 0, err
}

func HasGroupInvitation(recipientID, groupID int64) (bool, error) {
	var count int
	err := Database.QueryRow(
		"SELECT COUNT(*) FROM invitations WHERE recipient_id = ? AND group_id = ?",
		recipientID, groupID,
	).Scan(&count)
	return count > 0, err
}

func GetInvitationID(invitedID, recipientID, groupID int64) (int64, error) {
	var invitationID int64
	err := Database.QueryRow(
		"SELECT id FROM invitations WHERE recipient_id = ? AND invited_id = ? AND group_id = ?",
		recipientID, invitedID, groupID,
	).Scan(&invitationID)
	return invitationID, err
}

func AcceptAllFriendInvitations(userID int64) error {
	rows, err := Database.Query(
		"SELECT id, invited_id, group_id FROM invitations WHERE recipient_id = ?",
		userID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pending [][]int64

	for rows.Next() {
		var invitationID, invitedID, groupID int64
		if err := rows.Scan(&invitationID, &invitedID, &groupID); err != nil {
			return err
		}
		if groupID == 0 {
			pending = append(pending, []int64{invitationID, invitedID, userID, 0})
		}
	}

	for _, data := range pending {
		if err := AcceptInvitation(data[0], data[1], data[2], data[3]); err != nil {
			return err
		}
	}

	return nil
}

func GetInvitedBy(userID, groupID int64) (int64, error) {
	var invitedBy int64
	err := Database.QueryRow(
		"SELECT invited_id FROM invitations WHERE recipient_id = ? AND group_id = ?",
		userID, groupID,
	).Scan(&invitedBy)
	return invitedBy, err
}
