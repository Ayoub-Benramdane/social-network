package database

import (
	structs "social-network/data"
)

func SeedCategories() error {
	if existingCategory := FetchAnyCategory(); existingCategory == nil {
		categoryNames := []string{
			"Sport", "General", "Tech", "Gaming", "Movies", "Music",
			"Health", "Travel", "Food", "Fashion", "Education",
			"Science", "Art", "Finance", "Lifestyle", "History",
		}

		categoryColors := []string{
			"#FF5733", "#33C3FF", "#8E44AD", "#E74C3C", "#F39C12", "#1ABC9C",
			"#2ECC71", "#3498Database", "#E67E22", "#9B59B6", "#34495E", "#16A085",
			"#D35400", "#C0392B", "#7F8C8D", "#BDC3C7",
		}

		categoryBackgrounds := []string{
			"#FFE5E0", "#E0F7FF", "#F3E5F5", "#FFE0E0", "#FFF3E0", "#E0FFF5",
			"#E0FFE0", "#E0EFFF", "#FFF0E0", "#F5E0FF", "#E0E0F5", "#E0FFF0",
			"#FFEDE0", "#FFE0E0", "#F0F0F0", "#F5F5F5",
		}

		for i, name := range categoryNames {
			_, err := Database.Exec(
				"INSERT INTO categories (name, color, background) VALUES (?, ?, ?)",
				name,
				categoryColors[i],
				categoryBackgrounds[i],
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func FetchAnyCategory() *structs.Category {
	var category structs.Category

	err := Database.QueryRow(
		"SELECT * FROM categories",
	).Scan(
		&category.CategoryID,
		&category.Name,
		&category.Color,
		&category.Background,
	)

	if err != nil {
		return nil
	}

	return &category
}

func FetchAllCategories() ([]structs.Category, error) {
	rows, err := Database.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []structs.Category

	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.CategoryID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func FetchTopCategories() ([]structs.Category, error) {
	rows, err := Database.Query(
		`SELECT c.id, c.name, COUNT(p.category_id)
		 FROM categories c
		 LEFT JOIN posts p ON p.category_id = c.id
		 GROUP BY c.id
		 ORDER BY COUNT(p.category_id) DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []structs.Category

	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(
			&category.CategoryID,
			&category.Name,
			&category.ItemCount,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func FetchCategoryByID(categoryID int64) (*structs.Category, error) {
	var category structs.Category

	err := Database.QueryRow(
		"SELECT id, name, color, background FROM categories WHERE id = ?",
		categoryID,
	).Scan(
		&category.CategoryID,
		&category.Name,
		&category.Color,
		&category.Background,
	)

	return &category, err
}
