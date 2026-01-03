package database

import (
	"database/sql"
	"time"

	structs "social-network/data"
)

func CreateStory(imageURL string, userID int64) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	result, err := Database.Exec(
		"INSERT INTO stories (user_id, image) VALUES (?, ?)",
		userID, imageURL,
	)
	if err != nil {
		return 0, err
	}

	storyID, err := result.LastInsertId()
	return storyID, err
}

func CreateStoryStatus(storyID int64, followerIDs []int64) error {
	mu.Lock()
	defer mu.Unlock()

	for _, followerID := range followerIDs {
		_, err := Database.Exec(
			"INSERT INTO stories_status (story_id, user_id, read) VALUES (?, ?, ?)",
			storyID, followerID, false,
		)
		return err
	}
	return nil
}

func GetStories(currentUser structs.User, followingUsers []structs.User) ([]structs.Stories, error) {
	followingUsers = append(followingUsers, currentUser)

	var stories []structs.Stories
	var expiredStoryIDs []int64

	for _, followedUser := range followingUsers {
		var userStories structs.Stories
		userStories.User = followedUser

		rows, err := Database.Query(
			"SELECT id, image, created_at FROM stories WHERE user_id = ?",
			followedUser.UserID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var story structs.Story
			if err := rows.Scan(&story.StoryID, &story.ImageURL, &story.CreatedAt); err != nil {
				return nil, err
			}

			if time.Since(story.CreatedAt) > 24*time.Hour {
				expiredStoryIDs = append(expiredStoryIDs, story.StoryID)
				continue
			}

			userStories.Items = append(userStories.Items, story)
		}

		if len(userStories.Items) > 0 {
			GetStoryStatus(&userStories.Items, currentUser.UserID)
			userStories.HasUnseen = userStories.Items[len(userStories.Items)-1].IsRead
			stories = append(stories, userStories)
		}
	}

	for _, storyID := range expiredStoryIDs {
		if err := DeleteStory(storyID); err != nil {
			return nil, err
		}
	}

	return stories, nil
}

func GetStoryStatus(stories *[]structs.Story, userID int64) error {
	for i := 0; i < len(*stories); i++ {
		err := Database.QueryRow(
			"SELECT read FROM stories_status WHERE story_id = ? AND user_id = ?",
			(*stories)[i].StoryID, userID,
		).Scan(&(*stories)[i].IsRead)

		if err == sql.ErrNoRows {
			InsertStoryStatus((*stories)[i].StoryID, userID)
		} else if err != nil {
			return err
		}
	}
	return nil
}

func InsertStoryStatus(storyID int64, userID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"INSERT INTO stories_status (story_id, user_id, read) VALUES (?, ?, ?)",
		storyID, userID, false,
	)
	return err
}

func MarkStoryAsSeen(storyID int64, userID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"UPDATE stories_status SET read = ? WHERE story_id = ? AND user_id = ?",
		1, storyID, userID,
	)
	return err
}

func DeleteStory(storyID int64) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := Database.Exec(
		"DELETE FROM stories WHERE id = ?",
		storyID,
	)
	return err
}
