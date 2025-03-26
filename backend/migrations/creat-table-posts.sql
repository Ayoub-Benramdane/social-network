-- +migrate Up
CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        privacy TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (session_token) ON DELETE CASCADE
    );

-- +migrate Down
DROP TABLE IF EXISTS posts;