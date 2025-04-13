package database

import (
	structs "social-network/data"
)

func CreateGroup(admin int64, name, description, image, cover, privacy string) (int64, error) {
	result, err := DB.Exec("INSERT INTO groups (name, description, image, cover, admin, privacy) VALUES (?, ?, ?, ?, ?, ?, ?)", name, description, image, cover, admin, privacy)
	if err != nil {
		return 0, err
	}
	group_id, err := result.LastInsertId()
	return group_id, err
}

func GetGroups(user_id int64) ([]structs.Group, error) {
	var groups []structs.Group
	rows, err := DB.Query("SELECT g.id, g.name, g.description, g.image, g.cover, g.admin, g.members FROM groups g JOIN group_members m ON g.id = m.group_id WHERE m.user_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group structs.Group
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Image, &group.Cover, &group.Admin, &group.TotalMembers)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func GetOtherGroups(user_id int64) ([]structs.Group, error) {
	var groups []structs.Group
	rows, err := DB.Query("SELECT g.id, g.name, g.description, g.image, g.cover, g.admin, g.members FROM groups g WHERE g.id NOT IN (SELECT group_id FROM group_members WHERE user_id = ?)", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group structs.Group
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Image, &group.Cover, &group.Admin, &group.TotalMembers)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func GetGroup(group_id int64) (structs.Group, error) {
	var group structs.Group
	err := DB.QueryRow("SELECT id, name, description, image, cover, admin, privacy FROM groups WHERE id = ?", group_id).Scan(&group.ID, &group.Name, &group.Description, &group.Image, &group.Cover, &group.Admin, &group.Privacy)
	return group, err
}

func GetCountUserGroups(id int64) (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE user_id = ?", id).Scan(&count)
	return count, err
}
