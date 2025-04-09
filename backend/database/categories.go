package database

import structs "social-network/data"

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

func CheckCategory() *structs.Category {
	var cat structs.Category
	err := DB.QueryRow("SELECT * FROM categories").Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil
	}
	return &cat
}

func GetCategories() ([]structs.Category, error) {
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

func GetBestCategories() ([]structs.Category, error) {
	rows, err := DB.Query("SELECT c.name FROM categories c JOIN (SELECT ca.id category_id, COUNT(pc.category_id) as count FROM post_categories pc JOIN categories ca ON ca.id = pc.category_id GROUP BY ca.id ORDER BY count DESC LIMIT 5) p ON c.id = p.category_id")
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
