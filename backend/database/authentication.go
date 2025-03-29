package database

import (
	structs "social-network/data"
	"time"

	"github.com/gofrs/uuid"
)

func RegisterUser(Username, FirstName, LastName, Email string, hashedPassword []byte, DateOfBirth time.Time, sessionToken uuid.UUID) error {
	_, err := DB.Exec("INSERT INTO users (username, firstname, lastname, email, date_of_birth, password, session_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", Username, FirstName, LastName, Email, DateOfBirth, hashedPassword, sessionToken)
	return err
}

func GetUserByEmail(email string) (structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT password, session_token FROM users WHERE email = ?", email).Scan(&user.Password, &user.SessionToken)
	return user, err
}

func UpdateSession(Email string, sessionToken uuid.UUID) error {
	_, err := DB.Exec("UPDATE users SET session_token = ? WHERE email = ?", sessionToken, Email)
	return err
}

func GetUserConnected(token string) (structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT id, session_token FROM users WHERE token = ?", token).Scan(&user.ID, &user.SessionToken)
	return user, err
}

func DeleteSession(user_id int64) error {
	_, err := DB.Exec("UPDATE users SET session_token = NULL WHERE id = ?", user_id)
	return err
}