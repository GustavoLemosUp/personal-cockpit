# 🗄️ Database Schema - Personal Cockpit

> Documentação completa do banco de dados SQLite

**Última atualização:** 26/12/2025

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Diagrama ER](#diagrama-er)
3. [Tabelas](#tabelas)
4. [Relacionamentos](#relacionamentos)
5. [Índices](#índices)
6. [Migrations](#migrations)
7. [Queries Comuns](#queries-comuns)

---

## 🎯 Visão Geral

### Tecnologia

- **SGBD:** SQLite 3
- **Driver Go:** `github.com/mattn/go-sqlite3`
- **Localização:** `~/.config/personal-cockpit/cockpit.db` (varia por OS)
- **Encoding:** UTF-8
- **Journal Mode:** WAL (Write-Ahead Logging)

### Características

- ✅ **ACID Compliant:** Transações seguras
- ✅ **Zero Configuration:** Sem servidor necessário
- ✅ **Embarcado:** Arquivo único
- ✅ **Leve:** ~600KB
- ✅ **Type Safety:** Strict mode habilitado

---

## 📊 Diagrama ER

```
┌─────────────────┐
│   categories    │
│─────────────────│
│ id (PK)         │
│ name            │
│ color           │
│ type            │
│ created_at      │
└─────────────────┘
         │
         │ 1
         │
         │ N
         ▼
┌─────────────────┐         ┌─────────────────┐
│     tasks       │         │      notes      │
│─────────────────│         │─────────────────│
│ id (PK)         │         │ id (PK)         │
│ title           │         │ title           │
│ description     │         │ content         │
│ status          │         │ category_id(FK) │
│ priority        │         │ created_at      │
│ category_id(FK) │         │ updated_at      │
│ due_date        │         └─────────────────┘
│ completed_at    │
│ created_at      │         ┌─────────────────┐
│ updated_at      │         │     events      │
└─────────────────┘         │─────────────────│
                            │ id (PK)         │
┌─────────────────┐         │ title           │
│    settings     │         │ description     │
│─────────────────│         │ start_date      │
│ key (PK)        │         │ end_date        │
│ value           │         │ all_day         │
│ updated_at      │         │ color           │
└─────────────────┘         │ created_at      │
                            │ updated_at      │
                            └─────────────────┘

┌─────────────────┐
│  task_files     │
│─────────────────│
│ id (PK)         │
│ task_id (FK)    │
│ file_path       │
│ file_name       │
│ file_size       │
│ mime_type       │
│ created_at      │
└─────────────────┘
```

---

## 📑 Tabelas

### 1. `categories`

Categorias para tarefas, notas e outros itens.

```sql
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    color TEXT DEFAULT '#3b82f6',
    type TEXT CHECK(type IN ('task', 'note', 'general')) DEFAULT 'general',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Campos

| Campo | Tipo | Descrição | Constraints |
|-------|------|-----------|-------------|
| `id` | INTEGER | ID único | PK, AUTO_INCREMENT |
| `name` | TEXT | Nome da categoria | NOT NULL, UNIQUE |
| `color` | TEXT | Cor em hex (#RRGGBB) | DEFAULT '#3b82f6' |
| `type` | TEXT | Tipo (task/note/general) | CHECK constraint |
| `created_at` | DATETIME | Data de criação | DEFAULT NOW |

#### Exemplo de dados

```sql
INSERT INTO categories (name, color, type) VALUES
    ('Trabalho', '#ef4444', 'task'),
    ('Pessoal', '#10b981', 'task'),
    ('Estudos', '#3b82f6', 'task'),
    ('Ideias', '#f59e0b', 'note');
```

---

### 2. `tasks`

Tarefas do sistema To-Do.

```sql
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
```

#### Campos

| Campo | Tipo | Descrição | Constraints |
|-------|------|-----------|-------------|
| `id` | INTEGER | ID único | PK, AUTO_INCREMENT |
| `title` | TEXT | Título da tarefa | NOT NULL |
| `description` | TEXT | Descrição detalhada | NULL |
| `status` | TEXT | Status atual | CHECK (pending/completed/cancelled) |
| `priority` | TEXT | Prioridade | CHECK (low/medium/high) |
| `category_id` | INTEGER | ID da categoria | FK → categories.id |
| `due_date` | DATETIME | Data de vencimento | NULL |
| `completed_at` | DATETIME | Data de conclusão | NULL |
| `created_at` | DATETIME | Data de criação | DEFAULT NOW |
| `updated_at` | DATETIME | Última atualização | DEFAULT NOW |

#### Status possíveis

- `pending`: Tarefa pendente
- `completed`: Tarefa concluída
- `cancelled`: Tarefa cancelada

#### Prioridades

- `low`: Baixa prioridade (🟢)
- `medium`: Média prioridade (🟡)
- `high`: Alta prioridade (🔴)

#### Triggers

```sql
-- Atualizar updated_at automaticamente
CREATE TRIGGER update_task_timestamp 
AFTER UPDATE ON tasks
BEGIN
    UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Setar completed_at quando status muda para completed
CREATE TRIGGER set_completed_at
AFTER UPDATE OF status ON tasks
WHEN NEW.status = 'completed' AND OLD.status != 'completed'
BEGIN
    UPDATE tasks SET completed_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

---

### 3. `notes`

Notas rápidas e anotações.

```sql
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
```

#### Campos

| Campo | Tipo | Descrição | Constraints |
|-------|------|-----------|-------------|
| `id` | INTEGER | ID único | PK, AUTO_INCREMENT |
| `title` | TEXT | Título da nota | NOT NULL |
| `content` | TEXT | Conteúdo (Markdown) | NULL |
| `category_id` | INTEGER | ID da categoria | FK → categories.id |
| `is_favorite` | INTEGER | Favorito (0/1) | DEFAULT 0 |
| `created_at` | DATETIME | Data de criação | DEFAULT NOW |
| `updated_at` | DATETIME | Última edição | DEFAULT NOW |

#### Triggers

```sql
CREATE TRIGGER update_note_timestamp 
AFTER UPDATE ON notes
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

---

### 4. `events`

Eventos do calendário.

```sql
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
```

#### Campos

| Campo | Tipo | Descrição | Constraints |
|-------|------|-----------|-------------|
| `id` | INTEGER | ID único | PK, AUTO_INCREMENT |
| `title` | TEXT | Título do evento | NOT NULL |
| `description` | TEXT | Descrição | NULL |
| `start_date` | DATETIME | Data/hora início | NOT NULL |
| `end_date` | DATETIME | Data/hora fim | NOT NULL, >= start_date |
| `all_day` | INTEGER | Evento de dia inteiro (0/1) | DEFAULT 0 |
| `color` | TEXT | Cor do evento | DEFAULT '#3b82f6' |
| `location` | TEXT | Local do evento | NULL |
| `reminder_minutes` | INTEGER | Lembrete (minutos antes) | NULL |
| `created_at` | DATETIME | Data de criação | DEFAULT NOW |
| `updated_at` | DATETIME | Última atualização | DEFAULT NOW |

#### Triggers

```sql
CREATE TRIGGER update_event_timestamp 
AFTER UPDATE ON events
BEGIN
    UPDATE events SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

---

### 5. `settings`

Configurações do aplicativo.

```sql
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Campos

| Campo | Tipo | Descrição | Constraints |
|-------|------|-----------|-------------|
| `key` | TEXT | Chave da configuração | PK |
| `value` | TEXT | Valor (JSON string) | NOT NULL |
| `updated_at` | DATETIME | Última atualização | DEFAULT NOW |

#### Configurações padrão

```sql
INSERT INTO settings (key, value) VALUES
    ('theme', 'light'),
    ('language', 'pt-BR'),
    ('notification_sound', 'true'),
    ('auto_backup', 'true'),
    ('backup_frequency', 'daily');
```

---

### 6. `task_files` (Fase 2)

Arquivos anexados a tarefas.

```sql
CREATE TABLE IF NOT EXISTS task_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    mime_type TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
```

#### Campos

| Campo | Tipo | Descrição | Constraints |
|-------|------|-----------|-------------|
| `id` | INTEGER | ID único | PK, AUTO_INCREMENT |
| `task_id` | INTEGER | ID da tarefa | FK → tasks.id |
| `file_path` | TEXT | Caminho do arquivo | NOT NULL |
| `file_name` | TEXT | Nome original | NOT NULL |
| `file_size` | INTEGER | Tamanho em bytes | NOT NULL |
| `mime_type` | TEXT | Tipo MIME | NULL |
| `created_at` | DATETIME | Data de upload | DEFAULT NOW |

---

## 🔗 Relacionamentos

### 1:N Relationships

```
categories (1) ──── (N) tasks
categories (1) ──── (N) notes
tasks (1) ──── (N) task_files
```

### Integridade Referencial

| Tabela Filho | Tabela Pai | On Delete |
|--------------|------------|-----------|
| tasks | categories | SET NULL |
| notes | categories | SET NULL |
| task_files | tasks | CASCADE |

---

## 📇 Índices

### Índices de Performance

```sql
-- Índice para busca de tarefas por status
CREATE INDEX idx_tasks_status ON tasks(status);

-- Índice para busca de tarefas por categoria
CREATE INDEX idx_tasks_category ON tasks(category_id);

-- Índice para busca de tarefas por data de vencimento
CREATE INDEX idx_tasks_due_date ON tasks(due_date);

-- Índice para busca de notas por categoria
CREATE INDEX idx_notes_category ON notes(category_id);

-- Índice para busca de eventos por data de início
CREATE INDEX idx_events_start_date ON events(start_date);

-- Índice full-text para busca em notas (opcional)
CREATE VIRTUAL TABLE notes_fts USING fts5(title, content);
```

---

## 🔄 Migrations

### Estrutura de Migrations

```go
// database/migrations.go
package database

const (
    // Version 1 - Initial schema
    migration_v1_categories = `...`
    migration_v1_tasks = `...`
    migration_v1_notes = `...`
    migration_v1_events = `...`
    migration_v1_settings = `...`
    
    // Version 2 - Add task_files
    migration_v2_task_files = `...`
)

func (db *DB) GetCurrentVersion() int {
    // Retorna versão atual do schema
}

func (db *DB) Migrate() error {
    version := db.GetCurrentVersion()
    
    migrations := []struct{
        version int
        sql string
    }{
        {1, migration_v1_categories},
        {1, migration_v1_tasks},
        // ...
    }
    
    for _, m := range migrations {
        if m.version > version {
            if err := db.Exec(m.sql); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

### Schema Version Table

```sql
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO schema_version (version) VALUES (1);
```

---

## 🔍 Queries Comuns

### Tarefas

```sql
-- Buscar tarefas pendentes
SELECT * FROM tasks 
WHERE status = 'pending' 
ORDER BY due_date ASC;

-- Buscar tarefas de hoje
SELECT * FROM tasks 
WHERE DATE(due_date) = DATE('now') 
AND status = 'pending';

-- Buscar tarefas por categoria
SELECT t.*, c.name as category_name, c.color as category_color
FROM tasks t
LEFT JOIN categories c ON t.category_id = c.id
WHERE c.name = 'Trabalho';

-- Tarefas vencidas
SELECT * FROM tasks 
WHERE due_date < datetime('now') 
AND status = 'pending'
ORDER BY due_date ASC;

-- Estatísticas de tarefas
SELECT 
    status,
    COUNT(*) as count,
    AVG(julianday(completed_at) - julianday(created_at)) as avg_completion_days
FROM tasks
WHERE created_at >= date('now', '-30 days')
GROUP BY status;
```

### Notas

```sql
-- Buscar notas recentes
SELECT * FROM notes 
ORDER BY updated_at DESC 
LIMIT 10;

-- Busca full-text em notas
SELECT * FROM notes 
WHERE title LIKE '%keyword%' 
OR content LIKE '%keyword%';

-- Notas favoritas
SELECT * FROM notes 
WHERE is_favorite = 1 
ORDER BY updated_at DESC;
```

### Eventos

```sql
-- Eventos de hoje
SELECT * FROM events 
WHERE DATE(start_date) = DATE('now')
ORDER BY start_date ASC;

-- Eventos da semana
SELECT * FROM events 
WHERE start_date BETWEEN date('now') AND date('now', '+7 days')
ORDER BY start_date ASC;

-- Próximos eventos
SELECT * FROM events 
WHERE start_date > datetime('now')
ORDER BY start_date ASC 
LIMIT 5;
```

### Categorias

```sql
-- Categorias com contagem de tarefas
SELECT 
    c.id,
    c.name,
    c.color,
    COUNT(t.id) as task_count
FROM categories c
LEFT JOIN tasks t ON c.id = t.category_id
WHERE c.type = 'task'
GROUP BY c.id;
```

---

## 🔧 Configuração SQLite

### Pragmas Recomendados

```sql
-- Journal mode para melhor performance em concorrência
PRAGMA journal_mode = WAL;

-- Foreign keys habilitadas
PRAGMA foreign_keys = ON;

-- Encoding UTF-8
PRAGMA encoding = "UTF-8";

-- Cache size (em páginas, -2000 = 2MB)
PRAGMA cache_size = -2000;

-- Synchronous mode (NORMAL é bom compromisso)
PRAGMA synchronous = NORMAL;

-- Temp store em memória
PRAGMA temp_store = MEMORY;
```

### Inicialização em Go

```go
func NewDB() (*DB, error) {
    dbPath := getDatabasePath()
    
    conn, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Configurar pragmas
    pragmas := []string{
        "PRAGMA journal_mode = WAL",
        "PRAGMA foreign_keys = ON",
        "PRAGMA synchronous = NORMAL",
        "PRAGMA cache_size = -2000",
        "PRAGMA temp_store = MEMORY",
    }
    
    for _, pragma := range pragmas {
        if _, err := conn.Exec(pragma); err != nil {
            return nil, err
        }
    }
    
    db := &DB{conn: conn}
    
    // Rodar migrations
    if err := db.Migrate(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

---

## 📊 Estimativas de Tamanho

### Tamanho Médio por Registro

| Tabela | Tamanho Médio | 1000 Registros |
|--------|---------------|----------------|
| tasks | ~500 bytes | ~500 KB |
| notes | ~1 KB | ~1 MB |
| events | ~300 bytes | ~300 KB |
| categories | ~100 bytes | ~100 KB |
| task_files | ~200 bytes | ~200 KB |

### Crescimento Estimado

- **Uso leve** (10 tarefas/dia): ~2 MB/ano
- **Uso médio** (50 tarefas/dia): ~10 MB/ano
- **Uso pesado** (200 tarefas/dia): ~40 MB/ano

---

## 🔐 Backup e Manutenção

### Backup

```go
// Função de backup
func (db *DB) Backup(destPath string) error {
    // SQLite permite backup enquanto em uso
    backupDB, err := sql.Open("sqlite3", destPath)
    if err != nil {
        return err
    }
    defer backupDB.Close()
    
    // Usar API de backup do SQLite
    // ...
}
```

### Vacuum

```sql
-- Limpar espaço não utilizado
VACUUM;

-- Analisar para otimizar query planner
ANALYZE;
```

---

## 📚 Referências

- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [SQLite Best Practices](https://www.sqlite.org/bestpractice.html)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)

---

**Última revisão:** 26/12/2025