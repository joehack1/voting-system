package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	// Use a local SQLite database file by default.
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "voting.db"
	}

	// Connect
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Enable foreign key constraints for SQLite.
	if _, err := DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	// Ensure required tables exist
	if err := ensureSchema(); err != nil {
		log.Println("Schema initialization skipped:", err)
	}

	fmt.Printf("Connected to SQLite (%s)\n", dbPath)
}

func ensureSchema() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS polls (
			id TEXT PRIMARY KEY,
			question TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NULL,
			is_active BOOLEAN NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE IF NOT EXISTS options (
			id TEXT PRIMARY KEY,
			poll_id TEXT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
			value TEXT NOT NULL,
			votes_count INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS votes (
			poll_id TEXT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
			option_id TEXT NOT NULL REFERENCES options(id) ON DELETE CASCADE,
			ip_address TEXT NOT NULL,
			voted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (poll_id, ip_address)
		);`,
	}

	for _, stmt := range statements {
		if _, err := DB.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
