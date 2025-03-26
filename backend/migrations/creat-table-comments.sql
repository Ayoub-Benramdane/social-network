-- +migrate Up
 CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
        );

-- func createTables() error {
-- _, err = DB.Exec(`
--         CREATE TABLE IF NOT EXISTS comments (
--             id INTEGER PRIMARY KEY AUTOINCREMENT,
--             post_id INTEGER NOT NULL,
--             user_id INTEGER NOT NULL,
--             content TEXT NOT NULL,
--             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--             FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
--             FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
--         );
--     `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'comments' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'comments' table created or already exists")
-- 	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS comments;