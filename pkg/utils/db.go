package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitPostgres() *sql.DB {
	// Getting environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Missing one or more required DB environment variables")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	// Ping with timeout
	pingErr := tryPingDB(db, 5*time.Second)
	if pingErr != nil {
		log.Fatalf("DB ping failed: %v", pingErr)
	}

	return db
}

// tryPingDB tries pinging the DB with a timeout
func tryPingDB(db *sql.DB, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- db.Ping()
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("ping timed out after %s", timeout)
	}
}
