-- +migrate Up
CREATE TABLE
    IF NOT EXISTS follows (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        follower_id INTEGER NOT NULL,
        following_id INTEGER NOT NULL,
        FOREIGN KEY (followig_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
        UNIQUE (follower_id, following_id)
    );

-- +migrate Down
DROP TABLE IF EXISTS follows;