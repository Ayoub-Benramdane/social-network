package database

import (
	structs "social-network/data"
	"time"

	"github.com/gofrs/uuid"
)

func RegisterUser(Username, FirstName, LastName, Email, Image string, HashedPassword []byte, DateOfBirth time.Time, SessionToken uuid.UUID) error {
	_, err := DB.Exec("INSERT INTO users (username, firstname, lastname, email, avatar, date_of_birth, password, session_token) VALUES (?, ?, ?, ?, ?, ?, ?)", Username, FirstName, LastName, Email, Image,  DateOfBirth, HashedPassword, SessionToken)
	return err
}

func GetUserByEmail(email string) (structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT username, avatar, password, session_token FROM users WHERE email = ?", email).Scan(&user.Username, &user.Avatar, &user.Password, &user.SessionToken)
	return user, err
}

func CheckUser(user_id int64) (structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT username, password, session_token FROM users WHERE id = ?", user_id).Scan(&user.Username, &user.Password, &user.SessionToken)
	return user, err
}

func UpdateSession(Email string, sessionToken uuid.UUID) error {
	_, err := DB.Exec("UPDATE users SET session_token = ? WHERE email = ?", sessionToken, Email)
	return err
}

func GetUserConnected(token string) (structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT id, username, session_token FROM users WHERE token = ?", token).Scan(&user.ID, &user.Username, &user.SessionToken)
	return user, err
}

func DeleteSession(user_id int64) error {
	_, err := DB.Exec("UPDATE users SET session_token = NULL WHERE id = ?", user_id)
	return err
}
