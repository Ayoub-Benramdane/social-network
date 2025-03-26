-- +migrate Up
CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sender_id INTEGER NOT NULL,
	receiver_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

-- func createTables() error {
-- _, err = DB.Exec(`
--     CREATE TABLE IF NOT EXISTS messages (
-- 	id INTEGER PRIMARY KEY AUTOINCREMENT,
-- 	sender_id INTEGER NOT NULL,
-- 	receiver_id INTEGER NOT NULL,
-- 	content TEXT NOT NULL,
-- 	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
-- 	FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
-- 	FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
-- );
-- `)
-- 	if err != nil {
-- 		log.Printf("Error creating 'messages' table: %v", err)
-- 		return err
-- 	} else {
-- 		log.Println("'messages' table created or already exists")
-- 	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS messages;