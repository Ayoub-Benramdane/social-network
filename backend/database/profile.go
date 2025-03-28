package database

import (
	structs "social-network/data"
)

func GetProfileInfo(user_id int64) (structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT id, username, firstname, lastname, email, date_of_birth, created_at, followers, following FROM users WHERE id = ?", user_id).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.DateOfBirth, &user.CreatedAt, &user.TotalFollowers, &user.TotalFollowing)
	if err != nil {
		return structs.User{}, err
	}
	user.TotalPosts, err = GetCountUserPosts(user_id)
	if err != nil {
		return structs.User{}, err
	}
	user.TotalLikes, err = GetCountUserLikes(user_id)
	if err != nil {
		return structs.User{}, err
	}
	user.TotalComments, err = GetCountUserComments(user_id)
	return user, err
}
