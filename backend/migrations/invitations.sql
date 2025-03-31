-- +migrate Up
CREATE TABLE
    IF NOT EXISTS invitations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        invited_id INTEGER NOT NULL,
        recipient_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (invited_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (recipient_id) REFERENCES users (id) ON DELETE CASCADE,
        NIQUE (invited_id, recipient_id)
    );

-- +migrate Down
DROP TABLE IF EXISTS invitations;