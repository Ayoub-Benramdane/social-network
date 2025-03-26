 -- +migrate Up
 CREATE TABLE IF NOT EXISTS post_likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            post_id INTEGER NOT NULL,
            is_like BOOLEAN NOT NULL,	
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			UNIQUE (user_id, post_id)
        );

-- func createTables() error {
-- _, err = DB.Exec(`
--         CREATE TABLE IF NOT EXISTS post_likes (
--             id INTEGER PRIMARY KEY AUTOINCREMENT,
--             user_id INTEGER NOT NULL,
--             post_id INTEGER NOT NULL,
--             is_like BOOLEAN NOT NULL,	
--             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--             FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
--             FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
-- 			UNIQUE (user_id, post_id)
--         );
--     `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'post_likes' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'post_likes' table created or already exists")
-- 	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS post_likes;