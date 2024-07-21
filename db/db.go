package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := inicializeDatabase(db); err != nil {
		return nil, err
	}

	slog.Info("Database inicialized")

	return &Database{db: db}, nil
}

func (db *Database) Close() {
	db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func inicializeDatabase(db *sql.DB) error {
	_, currentFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(currentFile), "migrations")

	upFile := filepath.Join(migrationsDir, "up.sql")
	sqlBytes, err := os.ReadFile(upFile)
	if err != nil {
		slog.Error("Error to read up file", "error", err)
		return fmt.Errorf("failed to read up file: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		slog.Error("Failed to execute up", "error", err)
		return fmt.Errorf("failed to execute up: %w", err)
	}

	slog.Info("Database up successfully")
	return nil
}
