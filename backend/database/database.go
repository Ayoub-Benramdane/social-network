package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

var DB *sql.DB

// InitDB initializes the database and applies migrations
func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./data/social-network.db")
	if err != nil {
		return err
	}

	// Enable foreign key support
	_, err = DB.Exec(`
        PRAGMA foreign_keys = ON
    `)
	if err != nil {
		return err
	}

	// Apply migrations
	err = applyMigrations()
	if err != nil {
		return err
	}

	// Create Categories table
	return CreateCategoryies()
}

// applyMigrations applies all migrations to the database
func applyMigrations() error {
	m := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	// Run migrations
	_, err := migrate.Exec(DB, "sqlite3", m, migrate.Up)
	if err != nil {
		log.Printf("Error applying migrations: %v", err)
		return err
	}
	log.Println("Migrations applied successfully!")
	return nil
}
