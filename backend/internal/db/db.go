package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Connect opens and returns a *sql.DB connection using environment variables.
func Connect() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pwd := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "contract_employee")
	sslmode := getEnv("DB_SSLMODE", "disable")

	log.Printf("DB config host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbName, sslmode)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pwd, dbName, sslmode,
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, err
	}

	var currentDB, currentSchema, searchPath string
	err = dbConn.QueryRow("SELECT current_database(), current_schema(), current_setting('search_path')").Scan(&currentDB, &currentSchema, &searchPath)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to database=%s schema=%s search_path=%s", currentDB, currentSchema, searchPath)

	return dbConn, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
