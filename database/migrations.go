package database

import (
	"fmt"
)

const CurrentSchemaVersion = 1

func (db *DB) RunMigrations() error {

	if err := db.createVersionTable(); err != nil {
		return err
	}

	currentVersion, err := db.getSchemaVersion()
	if err != nil {
		return err
	}

	fmt.Printf("📊 Schema atual: v%d | Schema necessário: v%d\n", currentVersion, CurrentSchemaVersion)

	if currentVersion >= CurrentSchemaVersion {
		fmt.Println("✅ Schema já está atualizado!")
		return nil
	}

	if err := db.executeMigrations(currentVersion); err != nil {
		return err
	}

	if err := db.setSchemaVersion(CurrentSchemaVersion); err != nil {
		return err
	}

	fmt.Printf("✅ Schema atualizado para v%d\n", CurrentSchemaVersion)
	return nil
}

func (db *DB) createVersionTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		description TEXT
	);
	`
	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) getSchemaVersion() (int, error) {
	var version int
	err := db.conn.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil {
		return 0, nil // Primeira vez, retorna 0
	}
	return version, nil
}

// setSchemaVersion registra nova versão
func (db *DB) setSchemaVersion(version int) error {
	query := "INSERT INTO schema_version (version, description) VALUES (?, ?)"
	_, err := db.conn.Exec(query, version, fmt.Sprintf("Schema v%d", version))
	return err
}

// executeMigrations executa migrations pendentes
func (db *DB) executeMigrations(currentVersion int) error {
	migrations := []Migration{
		// ═══════════════════════════════════════
		// VERSÃO 1 - Schema Inicial
		// ═══════════════════════════════════════
		{
			Version:     1,
			Description: "Criar tabelas iniciais",
			SQL: []string{
				createCategoriesTable,
				createTasksTable,
				createNotesTable,
				createEventsTable,
				createSettingsTable,
				createIndexes,
				createTriggers,
			},
		},
	}

	for _, migration := range migrations {
		if migration.Version > currentVersion {
			fmt.Printf("🔄 Executando migration v%d: %s\n", migration.Version, migration.Description)

			for i, sql := range migration.SQL {
				if _, err := db.conn.Exec(sql); err != nil {
					return fmt.Errorf("erro na migration v%d [passo %d]: %w", migration.Version, i+1, err)
				}
			}
		}
	}

	return nil
}

type Migration struct {
	Version     int
	Description string
	SQL         []string
}

// ═══════════════════════════════════════════════════════════
// MIGRATIONS - VERSÃO 1
// ═══════════════════════════════════════════════════════════

const createCategoriesTable = `
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    color TEXT DEFAULT '#3b82f6',
    type TEXT CHECK(type IN ('task', 'note', 'general')) DEFAULT 'general',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createTasksTable = `
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK(status IN ('pending', 'completed', 'cancelled')) DEFAULT 'pending',
    priority TEXT CHECK(priority IN ('low', 'medium', 'high')) DEFAULT 'medium',
    category_id INTEGER,
    due_date DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
`

const createNotesTable = `
CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT,
    category_id INTEGER,
    is_favorite INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
`

const createEventsTable = `
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    all_day INTEGER DEFAULT 0,
    color TEXT DEFAULT '#3b82f6',
    location TEXT,
    reminder_minutes INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CHECK(end_date >= start_date)
);
`

const createSettingsTable = `
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_category ON tasks(category_id);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
CREATE INDEX IF NOT EXISTS idx_notes_category ON notes(category_id);
CREATE INDEX IF NOT EXISTS idx_events_start_date ON events(start_date);
`

const createTriggers = `
CREATE TRIGGER IF NOT EXISTS update_task_timestamp 
AFTER UPDATE ON tasks
FOR EACH ROW
BEGIN
    UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS set_completed_at
AFTER UPDATE OF status ON tasks
FOR EACH ROW
WHEN NEW.status = 'completed' AND OLD.status != 'completed'
BEGIN
    UPDATE tasks SET completed_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_note_timestamp 
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_event_timestamp 
AFTER UPDATE ON events
FOR EACH ROW
BEGIN
    UPDATE events SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
`
