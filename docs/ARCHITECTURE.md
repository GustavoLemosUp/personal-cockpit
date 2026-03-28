# Arquitetura — Personal Cockpit

> Documentação técnica da arquitetura atual do app

**Última atualização:** Março 2026

---

## Visão Geral

Personal Cockpit é uma aplicação desktop **híbrida** construída com Wails. O frontend React roda dentro de um WebView nativo; o backend Go gerencia dados, lógica de negócio e expõe funções para o frontend via bindings automáticos do Wails.

```
┌──────────────────────────────────────────────┐
│              PERSONAL COCKPIT                │
│           (executável único ~15MB)           │
└──────────────────────────────────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌──────────────────────┐
│  FRONTEND       │◄──►│  BACKEND             │
│  React/TS       │    │  Go                  │
│  (WebView)      │    │  (goroutine nativa)  │
└─────────────────┘    └──────────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  SQLite         │
                       │  (arquivo .db)  │
                       └─────────────────┘
```

---

## Stack

### Backend

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| Go | 1.24 | Linguagem do backend |
| Wails | 2.11.0 | Framework desktop |
| modernc.org/sqlite | 1.42+ | Driver SQLite puro Go |

> **Por que `modernc.org/sqlite` e não `go-sqlite3`?** O driver da modernc é puro Go — não exige CGO nem compilador C instalado. Isso simplifica o build em todas as plataformas.

### Frontend

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| React | 18.2 | UI |
| TypeScript | 4.6 | Tipagem |
| Vite | 3.x | Build e dev server |
| CSS modules | — | Estilos (sem Tailwind) |

### Render Engine

| OS | Engine |
|----|--------|
| Windows | WebView2 (Edge Chromium) |
| macOS | WKWebView |
| Linux | WebKitGTK |

---

## Estrutura de Pastas (código real)

```
personal-cockpit/
│
├── main.go                    # Entry point — inicializa Wails e App
├── app.go                     # Struct App com todos os bindings expostos
├── version.go                 # Constante da versão atual
├── wails.json                 # Config do app (nome, ícone, janela)
├── go.mod / go.sum            # Dependências Go
│
├── database/
│   ├── db.go                  # Conexão SQLite, pragmas e configuração
│   └── migrations.go          # Schema versionado (tabelas, índices, triggers)
│
├── models/
│   ├── task.go                # Struct Task + TaskFilter
│   ├── note.go                # Struct Note
│   ├── event.go               # Struct Event
│   └── category.go            # Struct Category
│
├── services/
│   ├── task.go                # TaskService — CRUD e validações
│   ├── note.go                # NoteService — CRUD, favoritos, busca
│   ├── event.go               # EventService — CRUD, filtros por data
│   └── category.go            # CategoryService — CRUD, filtros por tipo
│
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    └── src/
        ├── main.tsx                  # Entry point React
        ├── App.tsx                   # Roteamento entre páginas (switch/state)
        ├── style.css                 # Import central de estilos
        │
        ├── components/
        │   ├── Layout.tsx            # Wrapper: Sidebar + conteúdo principal
        │   ├── Sidebar.tsx           # Menu de navegação lateral
        │   ├── Dashboard.tsx         # Página inicial com stats
        │   ├── tasks/
        │   │   ├── TasksPage.tsx     # Página de tarefas
        │   │   ├── TaskForm.tsx      # Modal de criação/edição
        │   │   └── TaskList.tsx      # Lista de itens de tarefa
        │   ├── notes/
        │   │   └── NotesPage.tsx     # Página de notas
        │   ├── events/
        │   │   └── EventsPage.tsx    # Página de eventos
        │   └── categories/
        │       └── CategoriesPage.tsx # Página de categorias
        │
        ├── hooks/
        │   └── useTheme.ts           # Hook para tema claro/escuro
        │
        ├── styles/
        │   ├── global.css            # Variáveis, reset, utilitários
        │   ├── layout.css            # Layout e sidebar
        │   ├── dashboard.css         # Estilos do dashboard
        │   ├── tasks.css             # Estilos de tarefas
        │   ├── notes.css             # Estilos de notas
        │   ├── events.css            # Estilos de eventos
        │   └── categories.css        # Estilos de categorias
        │
        └── assets/
            └── fonts/                # Nunito (woff2)
```

> **Nota:** Os bindings TypeScript (`wailsjs/`) são **gerados automaticamente** pelo Wails a partir do código Go durante `wails dev` ou `wails build`. Não edite esses arquivos manualmente.

---

## Camadas da Aplicação

### 1. Apresentação (Frontend React)

Responsável por renderizar a UI e chamar as funções Go via bindings.

```tsx
// Exemplo: TasksPage.tsx
import { GetAllTasks, CreateTask } from '../../wailsjs/go/main/App';
import type { main } from '../../wailsjs/go/models';

function TasksPage() {
  const [tasks, setTasks] = useState<main.Task[]>([]);

  useEffect(() => {
    GetAllTasks().then(setTasks);
  }, []);

  async function handleCreate(task: main.Task) {
    await CreateTask(task);
    const updated = await GetAllTasks();
    setTasks(updated);
  }
  // ...
}
```

### 2. Bindings (app.go)

A struct `App` expõe métodos Go para o frontend. Todo método público em `App` vira uma função assíncrona no frontend automaticamente.

```go
// app.go
func (a *App) CreateTask(task models.Task) error {
    return a.taskService.CreateTask(task)
}

func (a *App) GetAllTasks() ([]models.Task, error) {
    return a.taskService.GetAllTasks()
}
```

### 3. Services (lógica de negócio)

Validações, regras de negócio e chamadas ao banco.

```go
// services/task.go
func (s *TaskService) CreateTask(task models.Task) error {
    if strings.TrimSpace(task.Title) == "" {
        return errors.New("título é obrigatório")
    }
    if task.Priority == "" {
        task.Priority = "medium"
    }
    task.Status = "pending"
    return s.db.CreateTask(task)
}
```

### 4. Database (acesso ao SQLite)

SQL direto, sem ORM. O banco é inicializado com pragmas para WAL mode e foreign keys.

```go
// database/db.go
func Open(dbPath string) (*DB, error) {
    conn, err := sql.Open("sqlite", dbPath)
    // ...
    conn.Exec("PRAGMA journal_mode = WAL")
    conn.Exec("PRAGMA foreign_keys = ON")
    conn.Exec("PRAGMA synchronous = NORMAL")
    // ...
}
```

---

## Fluxo de Dados: Criando uma Tarefa

```
Usuário preenche formulário
        │
        ▼
TaskForm.tsx valida campos
        │
        ▼
Chama CreateTask(task) — binding Wails
        │
        │  [serialização automática JSON → Go]
        ▼
app.go: App.CreateTask(task)
        │
        ▼
services/task.go: TaskService.CreateTask(task)
  - Valida título
  - Define status="pending"
  - Chama db.CreateTask()
        │
        ▼
database: INSERT INTO tasks ...
        │
        │  [resposta Go → JSON]
        ▼
Frontend recebe retorno (nil ou erro)
        │
        ▼
UI atualiza lista de tarefas
```

---

## Comunicação Frontend ↔ Backend

O Wails gera automaticamente interfaces TypeScript a partir das structs e métodos Go:

```go
// models/task.go (Go)
type Task struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Status      string `json:"status"`
    Priority    string `json:"priority"`
    // ...
}
```

Gera automaticamente:

```typescript
// wailsjs/go/models.ts (gerado — não editar)
export namespace main {
  export class Task {
    id: number;
    title: string;
    status: string;
    priority: string;
    // ...
  }
}
```

```typescript
// wailsjs/go/main/App.ts (gerado — não editar)
export function CreateTask(arg1: main.Task): Promise<void>;
export function GetAllTasks(): Promise<Array<main.Task>>;
```

---

## Banco de Dados

**Localização do arquivo:**
- Windows: `C:\Users\{user}\AppData\Roaming\Personal Cockpit\cockpit.db`
- macOS: `~/Library/Application Support/Personal Cockpit/cockpit.db`
- Linux: `~/.config/Personal Cockpit/cockpit.db`

**Tabelas atuais:** `tasks`, `notes`, `events`, `categories`, `settings`, `schema_version`

**Migrations:** versionadas na tabela `schema_version`. Cada versão é aplicada apenas uma vez, em ordem. Para adicionar uma tabela nova, crie uma nova migration em `database/migrations.go`.

Ver schema completo: [DATABASE.md](DATABASE.md)

---

## Adicionando um Novo Módulo

Para adicionar um módulo novo (ex: lista de compras), o padrão é:

1. **Model** — crie `models/shopping.go` com a struct e tags JSON
2. **Migration** — adicione a nova versão em `database/migrations.go`
3. **Service** — crie `services/shopping.go` com CRUD e validações
4. **Binding** — adicione os métodos em `app.go` (serão expostos automaticamente)
5. **Frontend** — crie `frontend/src/components/shopping/` com os componentes
6. **Estilo** — crie `frontend/src/styles/shopping.css`
7. **Roteamento** — adicione o case em `App.tsx` e o link em `Sidebar.tsx`

---

## Decisões Técnicas

### Por que Wails e não Electron?

Electron embarca um Chromium completo → ~120MB. Wails usa o WebView nativo do SO → ~15MB. Para um app pessoal offline, peso importa.

### Por que SQLite sem ORM?

SQL direto é mais explícito, mais performático e sem dependências extras. A quantidade de queries é pequena e bem controlada.

### Por que CSS puro e não Tailwind?

O app tem estilos bem específicos e escopo pequeno. CSS modular por página é suficiente e evita adicionar um toolchain a mais. Pode mudar no futuro se o app crescer muito.

### Por que TypeScript 4.6 e não 5.x?

Restrição do Wails v2.11.0 no template gerado. Pode ser atualizado junto com uma atualização do Wails.

---

## Performance

| Métrica | Target |
|---------|--------|
| Inicialização | < 2s |
| Operações CRUD | < 100ms |
| Executável | < 20MB |
| Memória em idle | < 100MB |

---

**Revisão:** Março 2026
