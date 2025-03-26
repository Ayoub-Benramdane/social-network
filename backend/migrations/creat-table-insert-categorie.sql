-- +migrate Up
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);


-- _, err = DB.Exec(`
--         INSERT OR IGNORE INTO categories (name) VALUES 
--         ('Technology'),
--         ('Lifestyle'),
--         ('Travel'),
--         ('Food'),
-- 		('Sport'),
-- 		('Other')
--     `)
-- 	if err != nil {
-- 		log.Printf("Error inserting default categories: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("Default categories inserted or already exist")
-- 	}

-- +migrate Down
DROP TABLE IF EXISTS categories;