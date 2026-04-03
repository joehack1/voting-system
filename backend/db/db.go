package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    // Connection string
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"))

    // Connect
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    // Test connection
    err = DB.Ping()
    if err != nil {
        log.Fatal("DB ping failed:", err)
    }

    fmt.Println("✅ Connected to PostgreSQL")
}