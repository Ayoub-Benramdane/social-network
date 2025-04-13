package database

import structs "social-network/data"

func GetConnections(user_id int64) ([]structs.User, error) {
	rows, err := DB.Query("SELECT DISTINCT u.id, u.username, u.firstname, u.lastname, u.avatar FROM users u JOIN follows f ON (u.id = f.follower_id OR u.id = f.following_id) WHERE (f.follower_id  = ? OR f.following_id = ?)", user_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var connections []structs.User
	for rows.Next() {
		var connection structs.User
		err = rows.Scan(&connection.ID, &connection.Username, &connection.FirstName, &connection.LastName, &connection.Avatar)
		if err != nil {
			return nil, err
		}
		connection.TotalMessages, err = GetCountConversationMessages(connection.ID, user_id)
		if err != nil {
			return nil, err
		}
		if connection.ID != user_id {
			connections = append(connections, connection)
		}
	}
	return connections, nil
}
