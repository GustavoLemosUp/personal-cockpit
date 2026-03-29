package database

import (
	"fmt"
)

const CurrentSchemaVersion = 5

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
		// ═══════════════════════════════════════
		// ═══════════════════════════════════════
		// VERSÃO 2 - Multi-perfil + Kanban
		// ═══════════════════════════════════════
		{
			Version:     2,
			Description: "Perfis familiares, boards Kanban e campos extras em tasks",
			SQL: []string{
				createProfilesTable,
				createBoardsTable,
				createKanbanColumnsTable,
				createSubtasksTable,
				createTaskAssigneesTable,
				createTaskDependenciesTable,
				createLabelsTable,
				createTaskLabelsTable,
				createTaskCommentsTable,
				createTaskActivityTable,
				alterTasksAddBoardID,
				alterTasksAddColumnID,
				alterTasksAddPosition,
				alterTasksAddStartDate,
				alterTasksAddEstimatedHours,
				alterTasksAddProgress,
				alterTasksAddGoogleEventID,
				createV2Indexes,
			},
		},
		// ═══════════════════════════════════════
		// VERSÃO 3 - Google Calendar para Eventos
		// ═══════════════════════════════════════
		{
			Version:     3,
			Description: "Adiciona google_event_id na tabela events",
			SQL: []string{
				alterEventsAddGoogleEventID,
			},
		},
		// ═══════════════════════════════════════
		// VERSÃO 4 - WhatsApp
		// ═══════════════════════════════════════
		{
			Version:     4,
			Description: "Tabelas de conversas, mensagens e agendamentos WhatsApp",
			SQL: []string{
				createWAChatsTable,
				createWAMessagesTable,
				createWAScheduledTable,
				createWAIndexes,
			},
		},
		// ═══════════════════════════════════════
		// VERSÃO 5 - Finanças + Estoque
		// ═══════════════════════════════════════
		{
			Version:     5,
			Description: "Gestão financeira pessoal e controle de estoque doméstico",
			SQL: []string{
				createFinanceAccountsTable,
				createFinanceTransactionsTable,
				createFinanceScheduledTable,
				createStockProductsTable,
				createStockMovementsTable,
				createShoppingListsTable,
				createShoppingItemsTable,
				createFinanceIndexes,
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

// ═══════════════════════════════════════════════════════════
// MIGRATIONS - VERSÃO 2 (WhatsApp)
// ═══════════════════════════════════════════════════════════

const createWAChatsTable = `
CREATE TABLE IF NOT EXISTS wa_chats (
    jid TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    last_message TEXT DEFAULT '',
    last_message_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    unread_count INTEGER DEFAULT 0
);
`

const createWAMessagesTable = `
CREATE TABLE IF NOT EXISTS wa_messages (
    id TEXT PRIMARY KEY,
    chat_jid TEXT NOT NULL,
    sender_jid TEXT DEFAULT '',
    content TEXT NOT NULL,
    is_from_me INTEGER DEFAULT 0,
    timestamp DATETIME NOT NULL,
    FOREIGN KEY (chat_jid) REFERENCES wa_chats(jid) ON DELETE CASCADE
);
`

const createWAScheduledTable = `
CREATE TABLE IF NOT EXISTS wa_scheduled (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_jid TEXT NOT NULL,
    chat_name TEXT DEFAULT '',
    content TEXT NOT NULL,
    scheduled_at DATETIME NOT NULL,
    sent INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createWAIndexes = `
CREATE INDEX IF NOT EXISTS idx_wa_messages_chat ON wa_messages(chat_jid, timestamp);
CREATE INDEX IF NOT EXISTS idx_wa_scheduled_pending ON wa_scheduled(scheduled_at, sent);
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

// ═══════════════════════════════════════════════════════════
// MIGRATIONS - VERSÃO 2
// ═══════════════════════════════════════════════════════════

const createProfilesTable = `
CREATE TABLE IF NOT EXISTS profiles (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    name                 TEXT NOT NULL,
    email                TEXT UNIQUE,
    google_id            TEXT UNIQUE,
    role                 TEXT CHECK(role IN ('admin', 'member')) DEFAULT 'member',
    avatar_url           TEXT,
    google_access_token  TEXT,
    google_calendar_id   TEXT DEFAULT 'primary',
    pin_hash             TEXT,
    is_active            INTEGER DEFAULT 1,
    created_at           DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createBoardsTable = `
CREATE TABLE IF NOT EXISTS boards (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    profile_id  INTEGER REFERENCES profiles(id) ON DELETE CASCADE,
    is_shared   INTEGER DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createKanbanColumnsTable = `
CREATE TABLE IF NOT EXISTS kanban_columns (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    board_id    INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    color       TEXT DEFAULT '#3b82f6',
    position    INTEGER NOT NULL DEFAULT 0,
    wip_limit   INTEGER,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createSubtasksTable = `
CREATE TABLE IF NOT EXISTS subtasks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id     INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    completed   INTEGER DEFAULT 0,
    position    INTEGER NOT NULL DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createTaskAssigneesTable = `
CREATE TABLE IF NOT EXISTS task_assignees (
    task_id     INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    profile_id  INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, profile_id)
);
`

const createTaskDependenciesTable = `
CREATE TABLE IF NOT EXISTS task_dependencies (
    task_id       INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    depends_on_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, depends_on_id),
    CHECK(task_id != depends_on_id)
);
`

const createLabelsTable = `
CREATE TABLE IF NOT EXISTS labels (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL,
    color    TEXT DEFAULT '#6b7280',
    board_id INTEGER REFERENCES boards(id) ON DELETE CASCADE
);
`

const createTaskLabelsTable = `
CREATE TABLE IF NOT EXISTS task_labels (
    task_id  INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label_id INTEGER NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);
`

const createTaskCommentsTable = `
CREATE TABLE IF NOT EXISTS task_comments (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id     INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    profile_id  INTEGER NOT NULL REFERENCES profiles(id),
    content     TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createTaskActivityTable = `
CREATE TABLE IF NOT EXISTS task_activity (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id     INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    profile_id  INTEGER REFERENCES profiles(id),
    action      TEXT NOT NULL,
    from_value  TEXT,
    to_value    TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const alterTasksAddBoardID = `ALTER TABLE tasks ADD COLUMN board_id INTEGER REFERENCES boards(id);`
const alterTasksAddColumnID = `ALTER TABLE tasks ADD COLUMN column_id INTEGER REFERENCES kanban_columns(id);`
const alterTasksAddPosition = `ALTER TABLE tasks ADD COLUMN position INTEGER DEFAULT 0;`
const alterTasksAddStartDate = `ALTER TABLE tasks ADD COLUMN start_date DATETIME;`
const alterTasksAddEstimatedHours = `ALTER TABLE tasks ADD COLUMN estimated_hours REAL;`
const alterTasksAddProgress = `ALTER TABLE tasks ADD COLUMN progress INTEGER DEFAULT 0;`
const alterTasksAddGoogleEventID = `ALTER TABLE tasks ADD COLUMN google_event_id TEXT;`

// ═══════════════════════════════════════════════════════════
// MIGRATIONS - VERSÃO 3
// ═══════════════════════════════════════════════════════════

const alterEventsAddGoogleEventID = `ALTER TABLE events ADD COLUMN google_event_id TEXT;`

const createV2Indexes = `
CREATE INDEX IF NOT EXISTS idx_boards_profile ON boards(profile_id);
CREATE INDEX IF NOT EXISTS idx_kanban_columns_board ON kanban_columns(board_id);
CREATE INDEX IF NOT EXISTS idx_subtasks_task ON subtasks(task_id);
CREATE INDEX IF NOT EXISTS idx_task_assignees_task ON task_assignees(task_id);
CREATE INDEX IF NOT EXISTS idx_task_assignees_profile ON task_assignees(profile_id);
CREATE INDEX IF NOT EXISTS idx_task_comments_task ON task_comments(task_id);
CREATE INDEX IF NOT EXISTS idx_task_activity_task ON task_activity(task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_board ON tasks(board_id);
CREATE INDEX IF NOT EXISTS idx_tasks_column ON tasks(column_id);
`

// ═══════════════════════════════════════════════════════════
// MIGRATIONS - VERSÃO 5 — Finanças + Estoque
// ═══════════════════════════════════════════════════════════

const createFinanceAccountsTable = `
CREATE TABLE IF NOT EXISTS finance_accounts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    type        TEXT NOT NULL DEFAULT 'checking',
    balance     REAL NOT NULL DEFAULT 0,
    currency    TEXT NOT NULL DEFAULT 'BRL',
    color       TEXT NOT NULL DEFAULT '#6b7280',
    icon        TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createFinanceTransactionsTable = `
CREATE TABLE IF NOT EXISTS finance_transactions (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id       INTEGER NOT NULL REFERENCES finance_accounts(id) ON DELETE CASCADE,
    to_account_id    INTEGER REFERENCES finance_accounts(id),
    type             TEXT NOT NULL,
    category         TEXT NOT NULL DEFAULT '',
    description      TEXT NOT NULL DEFAULT '',
    amount           REAL NOT NULL,
    date             DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_recurring     INTEGER NOT NULL DEFAULT 0,
    recurrence_rule  TEXT NOT NULL DEFAULT '',
    tags             TEXT NOT NULL DEFAULT '',
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createFinanceScheduledTable = `
CREATE TABLE IF NOT EXISTS finance_scheduled (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id       INTEGER NOT NULL REFERENCES finance_accounts(id) ON DELETE CASCADE,
    type             TEXT NOT NULL,
    category         TEXT NOT NULL DEFAULT '',
    description      TEXT NOT NULL DEFAULT '',
    amount           REAL NOT NULL,
    scheduled_at     DATETIME NOT NULL,
    is_recurring     INTEGER NOT NULL DEFAULT 0,
    recurrence_rule  TEXT NOT NULL DEFAULT '',
    is_executed      INTEGER NOT NULL DEFAULT 0,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createStockProductsTable = `
CREATE TABLE IF NOT EXISTS stock_products (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          TEXT NOT NULL,
    category      TEXT NOT NULL DEFAULT '',
    unit          TEXT NOT NULL DEFAULT 'un',
    min_quantity  REAL NOT NULL DEFAULT 0,
    current_stock REAL NOT NULL DEFAULT 0,
    price         REAL NOT NULL DEFAULT 0,
    notes         TEXT NOT NULL DEFAULT '',
    is_active     INTEGER NOT NULL DEFAULT 1,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createStockMovementsTable = `
CREATE TABLE IF NOT EXISTS stock_movements (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id  INTEGER NOT NULL REFERENCES stock_products(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    quantity    REAL NOT NULL,
    price       REAL NOT NULL DEFAULT 0,
    notes       TEXT NOT NULL DEFAULT '',
    date        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createShoppingListsTable = `
CREATE TABLE IF NOT EXISTS shopping_lists (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT NOT NULL,
    month        TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    is_completed INTEGER NOT NULL DEFAULT 0,
    total_budget REAL NOT NULL DEFAULT 0,
    total_spent  REAL NOT NULL DEFAULT 0,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createShoppingItemsTable = `
CREATE TABLE IF NOT EXISTS shopping_items (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    shopping_list_id INTEGER NOT NULL REFERENCES shopping_lists(id) ON DELETE CASCADE,
    product_id       INTEGER REFERENCES stock_products(id),
    name             TEXT NOT NULL,
    quantity         REAL NOT NULL DEFAULT 1,
    unit             TEXT NOT NULL DEFAULT 'un',
    estimated_price  REAL NOT NULL DEFAULT 0,
    actual_price     REAL NOT NULL DEFAULT 0,
    category         TEXT NOT NULL DEFAULT '',
    is_bought        INTEGER NOT NULL DEFAULT 0,
    notes            TEXT NOT NULL DEFAULT '',
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createFinanceIndexes = `
CREATE INDEX IF NOT EXISTS idx_finance_tx_account ON finance_transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_finance_tx_date ON finance_transactions(date);
CREATE INDEX IF NOT EXISTS idx_finance_tx_type ON finance_transactions(type);
CREATE INDEX IF NOT EXISTS idx_finance_scheduled_account ON finance_scheduled(account_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_product ON stock_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_shopping_items_list ON shopping_items(shopping_list_id);
`
