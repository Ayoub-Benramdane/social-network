-- +migrate Up
CREATE TABLE IF NOT EXISTS categories (
	            id INTEGER PRIMARY KEY AUTOINCREMENT,
	            name TEXT UNIQUE NOT NULL
	        );

-- func createTables() error {
--  _, err = DB.Exec(`
--          CREATE TABLE IF NOT EXISTS categories (
--              id INTEGER PRIMARY KEY AUTOINCREMENT,
--              name TEXT UNIQUE NOT NULL
--          );
--      `)
--  	if err != nil {
--  		log.Printf("Error creating 'categories' table: %v", err)
--  		return err
--  	} else {
--  		log.Println("'categories' table created or already exists")
--  	}
-- }

-- +migrate Down
DROP TABLE IF EXISTS categories;