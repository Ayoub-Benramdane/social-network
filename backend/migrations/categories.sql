-- +migrate Up
CREATE TABLE
	IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

-- +migrate Down
DROP TABLE IF EXISTS categories;