package database

func AddGroupMember(user_id, group_id int64) error {
	_, err := DB.Exec("INSERT INTO members (user_id, group_id) VALUES (?, ?)", user_id, group_id)
	return err
}
