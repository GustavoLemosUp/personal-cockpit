package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"personal-cockpit/config"
	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter caminho do banco: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	db := &DB{conn: conn}

	if err := db.configurePragmas(); err != nil {
		return nil, err
	}

	if err := db.RunMigrations(); err != nil {
		return nil, fmt.Errorf("erro ao executar migrations: %w", err)
	}

	return db, nil
}

func getDatabasePath() (string, error) {
	cfg := config.MustLoad()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, cfg.DBDirName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, cfg.DBFile), nil
}

func (db *DB) configurePragmas() error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA foreign_keys = ON",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -2000",
		"PRAGMA temp_store = MEMORY",
		"PRAGMA busy_timeout = 5000",
	}

	for _, pragma := range pragmas {
		if _, err := db.conn.Exec(pragma); err != nil {
			return fmt.Errorf("erro ao configurar pragma: %w", err)
		}
	}

	return nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}
