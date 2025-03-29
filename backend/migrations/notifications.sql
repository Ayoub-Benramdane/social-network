-- +migrate Up
CREATE TABLE
    IF NOT EXISTS notifications (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        notified_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        type_notification TEXT NOT NULL,
        read BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

-- +migrate Down
DROP TABLE IF EXISTS notifications;