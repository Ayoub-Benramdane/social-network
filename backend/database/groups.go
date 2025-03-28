package database

import (
	structs "social-network/data"
)

func CreateGroup(admin int64, name, description, image string) (int64, error) {
	result, err := DB.Exec("INSERT INTO groups (name, description, image, admin, ) VALUES (?, ?)", name, description)
	if err != nil {
		return 0, err
	}
	group_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return group_id, nil
}

func JoinGroup(user_id, group_id int64) error {
	_, err := DB.Exec("INSERT INTO members (user_id, group_id) VALUES (?, ?)", user_id, group_id)
	if err != nil {
		return err
	}
	return nil
}

func LeaveGroup(user_id, group_id int64) error {
	_, err := DB.Exec("DELETE FROM members WHERE user_id = ? AND group_id = ?", user_id, group_id)
	if err != nil {
		return err
	}
	return nil
}

func GetGroups(user_id int64) ([]structs.Group, error) {
	var groups []structs.Group
	rows, err := DB.Query("SELECT g.id, g.name, g.description, g.image, g.admin FROM groups g JOIN members m ON g.id = m.group_id WHERE m.user_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group structs.Group
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Image, &group.Admin)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func GetOtherGroups(user_id int64) ([]structs.Group, error) {
	var groups []structs.Group
	rows, err := DB.Query("SELECT id, name, description, image, admin FROM groups WHERE id NOT IN (SELECT group_id FROM members WHERE user_id = ?)", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group structs.Group
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Image, &group.Admin)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func GetGroup(group_id int64) (structs.Group, error) {
	var group structs.Group
	err := DB.QueryRow("SELECT id, name, description, image, admin FROM groups WHERE id = ?", group_id).Scan(&group.ID, &group.Name, &group.Description, &group.Image, &group.Admin)
	if err != nil {
		return structs.Group{}, err
	}
	return group, nil
}

func GetGroupMembers(group_id int64) ([]structs.User, error) {
	var members []structs.User
	rows, err := DB.Query("SELECT u.id, u.username, u.firstname, u.lastname, u.avatar FROM users u JOIN members m ON u.id = m.user_id WHERE m.group_id = ?", group_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member structs.User
		err = rows.Scan(&member.ID, &member.Username, &member.FirstName, &member.LastName, &member.ProfileImage)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func GetGroupPosts(group_id int64) ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := DB.Query("SELECT p.id, p.title, p.content, u.username, p.created_at, p.total_likes, p.total_comments FROM posts p JOIN users u ON p.user_id = u.id WHERE group_id = ?", group_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post structs.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt, &post.TotalLikes, &post.TotalComments)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
