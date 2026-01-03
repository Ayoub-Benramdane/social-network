package database

import (
	structs "social-network/data"
)

func GetProfileInfo(profileUserID int64, followedUsers []structs.User) (structs.User, error) {
	var profile structs.User

	err := Database.QueryRow(`
		SELECT id, username, firstname, lastname, email,
		       avatar, cover, about, privacy,
		       date_of_birth, created_at,
		       followers, following
		FROM users
		WHERE id = ?
	`, profileUserID).Scan(&profile.UserID, &profile.Username, &profile.FirstName, &profile.LastName,
		&profile.Email, &profile.AvatarURL, &profile.CoverURL, &profile.Bio, &profile.PrivacyLevel,
		&profile.BirthDate, &profile.CreatedAt, &profile.FollowerCount, &profile.FollowingCount)
	if err != nil {
		return profile, err
	}

	profile.PostCount, err = GetCountUserPosts(profileUserID, 0)
	if err != nil {
		return profile, err
	}

	profile.GroupCount, err = CountUserGroups(profileUserID)
	if err != nil {
		return profile, err
	}

	profile.EventCount, err = CountUserEvents(profileUserID)
	if err != nil {
		return profile, err
	}

	profile.LikeCount, err = CountLikesByUser(profileUserID)
	if err != nil {
		return profile, err
	}

	profile.CommentCount, err = CountUserComments(profileUserID)
	if err != nil {
		return profile, err
	}

	userGroups, err := GetUserGroups(profile)
	if err != nil {
		return profile, err
	}

	profile.ChatMessageCount, profile.GroupMessageCount, err =
		CountUserUnreadMessages(profileUserID, userGroups)
	if err != nil {
		return profile, err
	}

	profile.FollowerCount, err = CountUserFollowers(profileUserID)
	if err != nil {
		return profile, err
	}

	profile.FollowingCount, err = CountUserFollowing(profileUserID)
	if err != nil {
		return profile, err
	}

	profile.MessageCount =
		profile.ChatMessageCount + profile.GroupMessageCount

	profile.NotificationCount, err =
		CountUnreadNotifications(profileUserID)
	if err != nil {
		return profile, err
	}

	profile.UserStories, err =
		GetStories(profile, followedUsers)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func UpdateProfile(userID int64, username, firstName, lastName, bio, privacyLevel string) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := Database.Exec(`
		UPDATE users SET username = ?, firstname = ?, lastname = ?, privacy = ?, about = ? WHERE id = ?
	`, username, firstName, lastName, privacyLevel, bio, userID)

	return err
}
