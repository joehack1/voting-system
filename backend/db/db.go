package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    // Connection string
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"))

    // Connect
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
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

	fmt.Println("✅ Connected to PostgreSQL")
}

func ensureSchema() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS polls (
			id UUID PRIMARY KEY,
			question TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMP NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE
		);`,
		`CREATE TABLE IF NOT EXISTS options (
			id UUID PRIMARY KEY,
			poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
			value TEXT NOT NULL,
			votes_count INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS votes (
			poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
			option_id UUID NOT NULL REFERENCES options(id) ON DELETE CASCADE,
			ip_address TEXT NOT NULL,
			voted_at TIMESTAMP NOT NULL DEFAULT NOW(),
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
