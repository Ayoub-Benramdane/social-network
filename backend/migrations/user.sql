-- +migrate Up
CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        firstname TEXT NOT NULL,
        lastname TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        date_of_birth DATE NOT NULL,
        password TEXT NOT NULL,
        avater TEXT,
        about TEXT,
        session_token TEXT NOT NULL,
        followers INTEGER DEFAULT 0,
        following INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

-- +migrate Down
DROP TABLE IF EXISTS users;