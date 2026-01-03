package database

import (
	"time"

	structs "social-network/data"

	"github.com/gofrs/uuid"
)

func CreateUser(username, firstName, lastName, email, bio, avatarURL, coverURL, privacyLevel string, hashedPassword []byte, birthDate time.Time, sessionID uuid.UUID) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		`INSERT INTO users 
		(username, firstname, lastname, email, avatar, cover, privacy, date_of_birth, password, session_token, about)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		username, firstName, lastName, email, avatarURL, coverURL, privacyLevel, birthDate, hashedPassword, sessionID, bio,
	)

	return err
}

func FindUserByEmail(email string) (structs.User, error) {
	var user structs.User

	err := Database.QueryRow(
		"SELECT username, avatar, cover, password FROM users WHERE email = ?",
		email,
	).Scan(&user.Username, &user.AvatarURL, &user.CoverURL, &user.Password)

	return user, err
}

func UserExists(userID int64) (structs.User, error) {
	var user structs.User

	err := Database.QueryRow(
		"SELECT username FROM users WHERE id = ?",
		userID,
	).Scan(&user.Username)

	return user, err
}

func UpdateUserSession(email string, sessionID uuid.UUID) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"UPDATE users SET session_token = ? WHERE email = ?",
		sessionID, email,
	)

	return err
}

func GetUserBySession(token string) (structs.User, error) {
	var user structs.User

	err := Database.QueryRow(
		"SELECT id, username, avatar, session_token FROM users WHERE session_token = ?",
		token,
	).Scan(&user.UserID, &user.Username, &user.AvatarURL, &user.SessionID)

	return user, err
}

func ClearUserSession(userID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"UPDATE users SET session_token = ? WHERE id = ?",
		"",
		userID,
	)

	return err
}

func FindUserByID(userID int64) (structs.User, error) {
	var user structs.User

	err := Database.QueryRow(
		"SELECT id, username, avatar, privacy FROM users WHERE id = ?",
		userID,
	).Scan(
		&user.UserID,
		&user.Username,
		&user.AvatarURL,
		&user.PrivacyLevel,
	)

	return user, err
}
