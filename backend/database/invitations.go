package database

import structs "social-network/data"

func CreateInvitation(user_id, invited_id int64) error {
	_, err := DB.Exec("INSERT INTO invitations (user_id, invited_id) VALUES (?, ?)", user_id, invited_id)
	return err
}

func GetInvitationsFriends(user_id int64) ([]structs.Invitation, error) {
	var invitations []structs.Invitation
	rows, err := DB.Query("SELECT i.id, u.id, u.username FROM invitations i JOIN users u ON i.invited_id = u.id WHERE recipient_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var invitation structs.Invitation
		err = rows.Scan(&invitation.ID, &invitation.SenderID, &invitation.Sender)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}

func GetInvitationsGroups(group_id int64) ([]structs.Invitation, error) {
	var invitations []structs.Invitation
	rows, err := DB.Query("SELECT i.id, u.id, u.username, u.avatar FROM invitations_groups i JOIN users u ON u.id = i.sender_id JOIN groups g ON i.group_id = g.id WHERE g.id = ?", group_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var invitation structs.Invitation
		err = rows.Scan(&invitation.ID, &invitation.SenderID, &invitation.Sender, &invitation.Avatar)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}

func DeleteInvitation(user_id, invited_id int64) error {
	_, err := DB.Exec("DELETE FROM invitations WHERE user_id = ? AND invited_id = ?", user_id, invited_id)
	return err
}

func DeleteInvitationGroup(user_id, invited_id int64) error {
	_, err := DB.Exec("DELETE FROM invitations_groups WHERE user_id = ? AND invited_id = ?", user_id, invited_id)
	return err
}

func CheckInvitation(user_id, invited_id int64) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM invitations WHERE user_id = ? AND invited_id = ?", user_id, invited_id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
