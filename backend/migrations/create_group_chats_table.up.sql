-- +migrate Up
CREATE TABLE
    IF NOT EXISTS group_chats (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        group_id INTEGER NOT NULL,
        sender_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        status TEXT DEFAULT 'unread',
        FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
        FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE
    );