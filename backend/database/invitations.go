package database

import structs "social-network/data"

func CreateInvitation(user_id, invited_id int64) error {
	_, err := DB.Exec("INSERT INTO invitations (user_id, invited_id) VALUES (?, ?)", user_id, invited_id)
	if err != nil {
		return err
	}
	return nil
}

func GetInvitationsFriends(user_id int64) ([]structs.Invitation, error) {
	var invitations []structs.Invitation
	rows, err := DB.Query("SELECT i.id, u.username FROM invitations i JOIN users u ON i.invited_id = u.id WHERE recipient_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var invitation structs.Invitation
		err = rows.Scan(&invitation.ID, &invitation.Sender)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}

func GetInvitationsGroups(user_id int64) ([]structs.Invitation, error) {
	var invitations []structs.Invitation
	rows, err := DB.Query("SELECT i.id, g.name FROM invitations i JOIN groups g ON i.invited_id = g.id WHERE recipient_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var invitation structs.Invitation
		err = rows.Scan(&invitation.ID, &invitation.Sender)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}

func DeleteInvitation(user_id, invited_id int64) error {
	_, err := DB.Exec("DELETE FROM invitations WHERE user_id = ? AND invited_id = ?", user_id, invited_id)
	if err != nil {
		return err
	}
	return nil
}