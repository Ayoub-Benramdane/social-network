package database

import (
	structs "social-network/backend/data"
)

func CheckCategory() *structs.Category {
	var cat structs.Category
	err := DB.QueryRow("SELECT * FROM categories").Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil
	}
	return &cat
}

func CreateCategoryies() error {
	if cat := CheckCategory(); cat == nil {
		categories := []string{"Sport", "General", "Tech", "Gaming", "Movies", "Music", "Health", "Travel", "Food", "Fashion", "Education", "Science", "Art", "Finance", "Lifestyle", "History"}
		for _, category := range categories {
			_, err := DB.Exec("INSERT INTO categories (name) VALUES (?)", category)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetPostCategories(post_id int64) ([]string, error) {
	rows, err := DB.Query("SELECT c.name FROM categories c JOIN post_category pc ON c.id = pc.category_id WHERE pc.post_id = ?", post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []string
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetAllCategories() ([]structs.Category, error) {
	rows, err := DB.Query("SELECT name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []structs.Category
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}