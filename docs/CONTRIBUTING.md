# Guia de Contribuição — Personal Cockpit

> Como contribuir para o desenvolvimento do Personal Cockpit

---

## Índice

1. [Onde Contribuir](#onde-contribuir)
2. [Configurando o Ambiente](#configurando-o-ambiente)
3. [Estrutura Real do Projeto](#estrutura-real-do-projeto)
4. [Áreas do App e Como Contribuir](#áreas-do-app-e-como-contribuir)
5. [Workflow de Desenvolvimento](#workflow-de-desenvolvimento)
6. [Padrões de Código](#padrões-de-código)
7. [Commits Semânticos](#commits-semânticos)
8. [Pull Requests](#pull-requests)
9. [Reportando Bugs](#reportando-bugs)
10. [Sugerindo Features](#sugerindo-features)

---

## Onde Contribuir

### Encontrando algo para fazer

- [Issues abertas](https://github.com/GustavoLemosUp/personal-cockpit/issues) — bugs e melhorias mapeadas
- Issues com `good first issue` — bom ponto de entrada para novos contribuidores
- [ROADMAP.md](ROADMAP.md) — o que está planejado para as próximas versões
- Qualquer bug que você encontrar no uso do app

### Tipos de contribuição

| Tipo | Exemplos |
|------|---------|
| Bug fix | Corrigir comportamento incorreto em qualquer módulo |
| Nova feature | Implementar algo do roadmap |
| UI/UX | Melhorar layout, acessibilidade, responsividade |
| Documentação | Melhorar ou corrigir este guia e os outros docs |
| Testes | Adicionar testes para services ou componentes |
| Novo módulo | Implementar um módulo completo (ex: Lista de Compras) |

---

## Configurando o Ambiente

### Pré-requisitos

| Ferramenta | Versão mínima | Download |
|------------|---------------|----------|
| Go | 1.24 | [golang.org/dl](https://golang.org/dl/) |
| Node.js | 18 | [nodejs.org](https://nodejs.org/) |
| Wails CLI | 2.11.0 | abaixo |
| Git | qualquer | [git-scm.com](https://git-scm.com) |

```bash
# Instalar Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verificar se tudo está em ordem
wails doctor
```

### Clonando e Rodando

```bash
# 1. Fork o repositório no GitHub, depois clone seu fork
git clone https://github.com/SEU-USUARIO/personal-cockpit.git
cd personal-cockpit

# 2. Instale dependências do frontend
cd frontend && npm install && cd ..

# 3. Rode em modo desenvolvimento
wails dev
```

O app abre com hot-reload. Mudanças no frontend (React/TS) recarregam automaticamente. Mudanças no backend (Go) recompilam e reiniciam o binário.

### VS Code (recomendado)

**Extensões:**
```json
{
  "recommendations": [
    "golang.go",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode"
  ]
}
```

**settings.json:**
```json
{
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  }
}
```

---

## Estrutura Real do Projeto

Esta é a estrutura **atual** do código (não planejada, não futura):

```
personal-cockpit/
├── main.go              # Inicializa Wails e passa App como binding
├── app.go               # Struct App — todos os métodos expostos ao frontend
├── version.go           # Versão atual
│
├── database/
│   ├── db.go            # Conexão SQLite e configuração de pragmas
│   └── migrations.go    # Schema e versões do banco
│
├── models/
│   ├── task.go          # Task, TaskFilter
│   ├── note.go          # Note
│   ├── event.go         # Event
│   └── category.go      # Category
│
├── services/
│   ├── task.go          # TaskService — lógica de negócio de tarefas
│   ├── note.go          # NoteService — lógica de notas
│   ├── event.go         # EventService — lógica de eventos
│   └── category.go      # CategoryService — lógica de categorias
│
└── frontend/src/
    ├── App.tsx                   # Switch/roteamento de páginas
    ├── components/
    │   ├── Layout.tsx
    │   ├── Sidebar.tsx
    │   ├── Dashboard.tsx
    │   ├── tasks/
    │   │   ├── TasksPage.tsx
    │   │   ├── TaskForm.tsx
    │   │   └── TaskList.tsx
    │   ├── notes/NotesPage.tsx
    │   ├── events/EventsPage.tsx
    │   └── categories/CategoriesPage.tsx
    ├── hooks/useTheme.ts
    └── styles/
        ├── global.css
        ├── layout.css
        ├── dashboard.css
        ├── tasks.css
        ├── notes.css
        ├── events.css
        └── categories.css
```

> **Não existe** pasta `utils/`, `handlers/`, `context/`, `services/` (no frontend), `types/` ou `queries.go`. Não crie esses arquivos só porque parecem fazer sentido — abra uma issue primeiro para discutir se fazem falta.

---

## Áreas do App e Como Contribuir

### Módulos já implementados (v1.0)

Cada módulo segue o mesmo padrão: model → migration → service → binding → componente → estilo.

---

#### Tarefas

**Escopo:** to-do list com prioridades, categorias e datas.

**Arquivos relevantes:**
- `models/task.go` — campos da tarefa
- `services/task.go` — CRUD, toggleStatus, GetTasksByFilter
- `app.go` — métodos CreateTask, GetAllTasks, UpdateTask, etc.
- `frontend/src/components/tasks/` — TasksPage, TaskForm, TaskList
- `frontend/src/styles/tasks.css`

**Como contribuir:**
- Filtros avançados (por data, por prioridade múltipla)
- Drag & drop para reordenar
- Notificação de vencimento
- Recorrência de tarefas

---

#### Notas

**Escopo:** notas livres com título, conteúdo, categorias e favoritos.

**Arquivos relevantes:**
- `models/note.go`
- `services/note.go` — CRUD, SearchNotes, ToggleFavorite
- `frontend/src/components/notes/NotesPage.tsx`
- `frontend/src/styles/notes.css`

**Como contribuir:**
- Editor Markdown com preview
- Auto-save ao digitar
- Ordenação por data de modificação
- Notas fixadas no topo

---

#### Eventos / Calendário

**Escopo:** agenda pessoal com data, hora, cor e localização.

**Arquivos relevantes:**
- `models/event.go`
- `services/event.go` — CRUD, GetTodayEvents, GetUpcomingEvents
- `frontend/src/components/events/EventsPage.tsx`
- `frontend/src/styles/events.css`

**Como contribuir:**
- Visualização em calendário mensal
- Eventos recorrentes
- Notificação de lembrete (via OS)
- Exportar para .ics

---

#### Categorias

**Escopo:** categorias com cor para organizar tarefas e notas.

**Arquivos relevantes:**
- `models/category.go`
- `services/category.go`
- `frontend/src/components/categories/CategoriesPage.tsx`
- `frontend/src/styles/categories.css`

---

### Adicionando um Módulo Novo

Para implementar um módulo do roadmap (ex: Lista de Compras), siga este padrão:

**1. Model** — `models/shopping.go`
```go
package models

type ShoppingList struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    // ...
}
```

**2. Migration** — adicione em `database/migrations.go`
```go
// Versão 2 — shopping
const v2CreateShoppingLists = `
CREATE TABLE IF NOT EXISTS shopping_lists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
```

**3. Service** — `services/shopping.go`
```go
type ShoppingService struct { db *database.DB }

func (s *ShoppingService) CreateList(list models.ShoppingList) error {
    // validações e INSERT
}
```

**4. Binding** — adicione em `app.go`
```go
func (a *App) CreateShoppingList(list models.ShoppingList) error {
    return a.shoppingService.CreateList(list)
}
```

**5. Componentes** — `frontend/src/components/shopping/`

**6. Estilo** — `frontend/src/styles/shopping.css`

**7. Roteamento** — adicione case em `App.tsx` e link em `Sidebar.tsx`

---

## Workflow de Desenvolvimento

### 1. Criar Branch

```bash
git checkout main
git pull origin main
git checkout -b feature/nome-descritivo
```

**Nomenclatura:**

| Tipo | Prefixo | Exemplo |
|------|---------|---------|
| Nova feature | `feature/` | `feature/lista-compras` |
| Bug fix | `fix/` | `fix/task-delete-erro` |
| Documentação | `docs/` | `docs/atualiza-readme` |
| Refatoração | `refactor/` | `refactor/note-service` |

### 2. Desenvolver

```bash
wails dev  # hot-reload ativo
```

### 3. Testar

```bash
# Backend Go
go test ./...
go test -v ./services/...

# Frontend
cd frontend
npm run build  # garante que compila sem erro
```

### 4. Commit e Push

```bash
git add arquivo1.go frontend/src/components/...
git commit -m "feat(shopping): adiciona modelo e service de lista de compras"
git push origin feature/lista-compras
```

---

## Padrões de Código

### Go (Backend)

```go
// Funções exportadas em PascalCase, com comentário
// CreateTask cria uma nova tarefa após validação.
func (s *TaskService) CreateTask(task models.Task) error {
    if strings.TrimSpace(task.Title) == "" {
        return errors.New("título é obrigatório")
    }
    // ...
}

// Sempre cheque erros — nunca use _
task, err := s.db.GetTask(id)
if err != nil {
    return fmt.Errorf("GetTask(%d): %w", id, err)
}

// Structs com tags JSON
type Task struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
}
```

```bash
# Formate antes de commitar
gofmt -w .
# ou
goimports -w .
```

### TypeScript/React (Frontend)

```tsx
// Componentes em PascalCase com interface de props explícita
interface TaskItemProps {
    task: models.Task;
    onDelete: (id: number) => void;
}

export function TaskItem({ task, onDelete }: TaskItemProps) {
    return <div>{task.title}</div>;
}

// Sem `any` — use tipos explícitos ou unknown
async function handleCreate(data: models.Task): Promise<void> { ... }

// useEffect sempre com dependências declaradas
useEffect(() => {
    loadTasks();
}, []); // [] = só na montagem
```

### SQL

```sql
-- Palavras-chave em maiúsculas, indentação clara
SELECT
    t.id,
    t.title,
    c.name AS category_name
FROM tasks t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.status = 'pending'
ORDER BY t.due_date ASC;
```

---

## Commits Semânticos

Usamos [Conventional Commits](https://www.conventionalcommits.org/).

### Formato

```
<tipo>(<escopo>): <descrição curta>
```

### Tipos e Escopos

| Tipo | Quando usar |
|------|------------|
| `feat` | Nova funcionalidade |
| `fix` | Correção de bug |
| `docs` | Apenas documentação |
| `style` | CSS, layout (sem mudança de lógica) |
| `refactor` | Reestruturação sem nova feature ou fix |
| `test` | Testes |
| `chore` | Configs, dependências, build |

**Escopos comuns:**

| Escopo | Área |
|--------|------|
| `tasks` | Módulo de tarefas |
| `notes` | Módulo de notas |
| `events` | Módulo de eventos |
| `categories` | Categorias |
| `dashboard` | Dashboard |
| `db` | Banco de dados |
| `shopping` | Lista de compras (futuro) |
| `habits` | Rastreador de hábitos (futuro) |
| `finance` | Controle financeiro (futuro) |

### Exemplos

```bash
git commit -m "feat(tasks): adiciona recorrência semanal em tarefas"
git commit -m "fix(notes): corrige busca que ignorava acentos"
git commit -m "style(dashboard): ajusta espaçamento dos cards de stats"
git commit -m "docs: atualiza guia de contribuição com novo padrão"
git commit -m "chore: atualiza dependências do frontend"
```

---

## Pull Requests

### Antes de abrir

- [ ] O build passa sem erros (`wails build`)
- [ ] O código segue os padrões do projeto
- [ ] Documentação atualizada se necessário

### Template de PR

```markdown
## O que muda
Breve descrição das mudanças.

## Tipo
- [ ] Bug fix
- [ ] Nova feature
- [ ] Documentação
- [ ] Refatoração

## Como testar
1. Rode `wails dev`
2. Acesse [módulo afetado]
3. [Descreva o que testar]

## Screenshots (se mudança visual)
[Cole aqui]

## Issues relacionadas
Closes #123
```

### Boas práticas

- Um PR = uma mudança coesa. Não misture bug fix com nova feature
- PRs pequenos são revisados mais rápido
- Adicione screenshots para qualquer mudança visual
- Responda ao feedback antes de pedir re-review

---

## Reportando Bugs

### Template

```markdown
**Descrição do bug**
O que aconteceu vs o que deveria acontecer.

**Como reproduzir**
1. Vá para ...
2. Clique em ...
3. Veja o erro

**Ambiente**
- OS: [Windows 11 / macOS 14 / Ubuntu 24]
- Versão do app: [v1.0.0]

**Logs ou screenshots**
[Cole aqui se disponível]
```

---

## Sugerindo Features

Antes de implementar algo grande, abra uma issue para discutir:

```markdown
**Qual problema resolve?**
[Descreva a dor ou limitação atual]

**Solução proposta**
[Como você imagina que funcione]

**Área do app**
[tasks / notes / events / nova área]

**Você quer implementar?**
[ ] Sim  [ ] Não  [ ] Talvez
```

Issues bem descritas têm mais chance de virar feature aceita.

---

## Recursos Úteis

| Recurso | Link |
|---------|------|
| Wails Docs | [wails.io/docs](https://wails.io/docs/introduction) |
| React Docs | [react.dev](https://react.dev) |
| Go Docs | [go.dev/doc](https://go.dev/doc/) |
| TypeScript Handbook | [typescriptlang.org/docs](https://www.typescriptlang.org/docs/) |
| SQLite Docs | [sqlite.org/docs](https://www.sqlite.org/docs.html) |
| Arquitetura do projeto | [ARCHITECTURE.md](ARCHITECTURE.md) |
| Schema do banco | [DATABASE.md](DATABASE.md) |
| Roadmap | [ROADMAP.md](ROADMAP.md) |

---

Obrigado por contribuir!
