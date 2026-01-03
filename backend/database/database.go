package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

var Database *sql.DB

func InitializeDatabase() error {
	var err error

	Database, err = sql.Open("sqlite3", "./data/social-network.db")
	if err != nil {
		return err
	}

	_, err = Database.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return err
	}

	if err := runMigrations(); err != nil {
		return err
	}

	return SeedCategories()
}

func runMigrations() error {
	migrationSource := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	_, err := migrate.Exec(Database, "sqlite3", migrationSource, migrate.Up)
	if err != nil {
		log.Printf("Migration error: %v", err)
		return err
	}

	log.Println("Database migrations applied successfully")
	return nil
}
