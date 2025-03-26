-- +migrate Up
CREATE TABLE IF NOT EXISTS comment_likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            comment_id INT NOT NULL,
            is_like BOOLEAN NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(comment_id, user_id),
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE
        );
-- func createTables() error {
-- _, err = DB.Exec(`
--         CREATE TABLE IF NOT EXISTS comment_likes (
--             id INTEGER PRIMARY KEY AUTOINCREMENT,
--             user_id INTEGER NOT NULL,
--             comment_id INT NOT NULL,
--             is_like BOOLEAN NOT NULL,
--             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
-- 			UNIQUE(comment_id, user_id),
--             FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
--             FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE
--         );
--     `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'comment_likes' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'comment_likes' table created or already exists")
-- 	}
-- }
-- +migrate Down
DROP TABLE IF EXISTS comment_likes;