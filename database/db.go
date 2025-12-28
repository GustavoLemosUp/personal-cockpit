package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
	// Obter diretório de configuração do usuário
	// Windows: C:\Users\{user}\AppData\Roaming
	// macOS: ~/Library/Application Support
	// Linux: ~/.config
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Criar diretório do app se não existir
	appDir := filepath.Join(configDir, "Personal Cockpit")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	// Retornar caminho completo do banco
	return filepath.Join(appDir, "cockpit.db"), nil
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
