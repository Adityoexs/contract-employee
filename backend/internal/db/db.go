package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Connect opens and returns a *sql.DB connection using environment variables.
func Connect() (*sql.DB, error) {
	host   := getEnv("DB_HOST",     "localhost")
	port   := getEnv("DB_PORT",     "5432")
	user   := getEnv("DB_USER",     "postgres")
	pwd    := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME",     "contract_employee")
	sslmode := getEnv("DB_SSLMODE", "disable")

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

	return dbConn, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
