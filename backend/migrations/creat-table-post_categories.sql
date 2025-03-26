 -- +migrate Up
 CREATE TABLE IF NOT EXISTS post_categories (
            post_id INTEGER NOT NULL,
            category_id INTEGER NOT NULL,
            PRIMARY KEY (post_id, category_id),
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
            FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
        );

-- func createTables() error {
-- _, err = DB.Exec(`
--         CREATE TABLE IF NOT EXISTS post_categories (
--             post_id INTEGER NOT NULL,
--             category_id INTEGER NOT NULL,
--             PRIMARY KEY (post_id, category_id),
--             FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
--             FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
--         );
--     `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'post_categories' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'post_categories' table created or already exists")
-- 	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS categories;