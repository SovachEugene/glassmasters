package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Get DSN from environment variable or fallback to default
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalf("Database DSN is not set in environment variables")
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	log.Println("Successfully connected to the database")
}

// GetDB provides access to the database connection instance
func GetDB() *sql.DB {
	if db == nil {
		log.Panic("Database connection is not initialized")
	}
	return db
}

// CloseDB safely closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error while closing the database: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}
}
