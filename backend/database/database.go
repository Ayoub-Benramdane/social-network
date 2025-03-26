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
	DB, err = sql.Open("sqlite3", "./backend/data/forum.db")
	if err != nil {
		return err
	}

	// Apply migrations
	if err := applyMigrations(); err != nil {
		return err
	}

	return nil
}

// applyMigrations applies all migrations to the database
func applyMigrations() error {
	m := &migrate.FileMigrationSource{
		Dir: "./backend/migrations",
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
