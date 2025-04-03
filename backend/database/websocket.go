package database

import structs "social-network/data"

func GetConnections(user_id int64) ([]structs.User, error) {
	// rows, err := DB.Query("SELECT id, username, firstname, lastname, email, date_of_birth, created_at, followers, following FROM users WHERE id IN (SELECT user_id FROM connections WHERE connection_id = ?)", user_id)
	return []structs.User{}, nil
}