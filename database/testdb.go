package database

import (
	"database/sql"
	"fmt"
)

// NewInMemoryDB creates a fresh in-memory SQLite database with all migrations applied.
// Use this in tests to get an isolated, clean database for each test.
//
//	db, err := database.NewInMemoryDB()
//	svc := services.NewTaskService(db)
func NewInMemoryDB() (*sql.DB, error) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco em memória: %w", err)
	}

	db := &DB{conn: conn}

	if err := db.configurePragmas(); err != nil {
		conn.Close()
		return nil, err
	}

	if err := db.RunMigrations(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("erro ao executar migrations: %w", err)
	}

	return conn, nil
}
