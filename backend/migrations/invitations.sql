-- +migrate Up
CREATE TABLE
    IF NOT EXISTS invitations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        recipient_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (recipient_id) REFERENCES users (id) ON DELETE CASCADE,
        NIQUE (sender_id, recipient_id)
    );

-- +migrate Down
DROP TABLE IF EXISTS invitations;