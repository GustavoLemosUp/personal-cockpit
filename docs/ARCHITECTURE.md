# 🏗️ Arquitetura - Personal Cockpit

> Documentação técnica da arquitetura do Personal Cockpit

**Última atualização:** 26/12/2025

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura de Alto Nível](#arquitetura-de-alto-nível)
3. [Stack Tecnológica](#stack-tecnológica)
4. [Estrutura de Pastas](#estrutura-de-pastas)
5. [Camadas da Aplicação](#camadas-da-aplicação)
6. [Fluxo de Dados](#fluxo-de-dados)
7. [Comunicação Frontend ↔ Backend](#comunicação-frontend--backend)
8. [Persistência de Dados](#persistência-de-dados)
9. [Decisões Arquiteturais](#decisões-arquiteturais)

---

## 🎯 Visão Geral

Personal Cockpit é uma aplicação desktop **híbrida** construída com Wails, combinando:

- **Backend em Go**: Lógica de negócio, acesso a dados, operações de sistema
- **Frontend em React**: Interface do usuário moderna e responsiva
- **SQLite**: Banco de dados local embarcado
- **WebView2**: Renderização nativa da interface

### Características

- ✅ **Aplicação Desktop Nativa**: Executável único, sem necessidade de navegador
- ✅ **100% Local**: Não requer internet, dados armazenados localmente
- ✅ **Multiplataforma**: Windows, macOS, Linux (mesmo código)
- ✅ **Leve e Rápido**: ~15MB, inicialização < 2s
- ✅ **Type-Safe**: Go e TypeScript garantem segurança de tipos

---

## 🏛️ Arquitetura de Alto Nível

```
┌─────────────────────────────────────────────────────────────┐
│                    PERSONAL COCKPIT                         │
│                   (Aplicação Desktop)                       │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┴──────────────────────┐
        │                                             │
        ▼                                             ▼
┌──────────────────┐                         ┌──────────────────┐
│   FRONTEND       │◄────── Wails ─────────►│    BACKEND       │
│   (React/TS)     │      Bindings          │    (Go)          │
└──────────────────┘                         └──────────────────┘
        │                                             │
        │                                             │
        │ Componentes                                 │ Services
        │ Hooks                                       │ Models
        │ Context                                     │ Handlers
        │                                             │
        │                                             ▼
        │                                    ┌─────────────────┐
        │                                    │   DATABASE      │
        │                                    │   (SQLite)      │
        │                                    └─────────────────┘
        │                                             │
        │                                             │
        └─────────────────┬───────────────────────────┘
                          │
                          ▼
                 ┌─────────────────┐
                 │  SISTEMA OPERAC.│
                 │  (Files, OS API)│
                 └─────────────────┘
```

---

## 🛠️ Stack Tecnológica

### Backend

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| **Go** | 1.23+ | Linguagem principal do backend |
| **Wails** | 2.11.0 | Framework para desktop apps |
| **SQLite** | 3.x | Banco de dados embarcado |
| **go-sqlite3** | Latest | Driver SQLite para Go |

### Frontend

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| **React** | 18.x | Framework UI |
| **TypeScript** | 5.x | Type safety |
| **Vite** | 3.x | Build tool e dev server |
| **TailwindCSS** | 3.x | Styling (planejado) |
| **React Router** | 6.x | Navegação (planejado) |

### Build & Deploy

| Tecnologia | Uso |
|------------|-----|
| **Wails CLI** | Build e desenvolvimento |
| **WebView2** | Renderização no Windows |
| **WKWebView** | Renderização no macOS |
| **WebKitGTK** | Renderização no Linux |

---

## 📁 Estrutura de Pastas

```
personal-cockpit/
│
├── main.go                      # Entry point da aplicação
├── app.go                       # Struct principal, bindings
├── wails.json                   # Configuração do Wails
├── go.mod                       # Dependências Go
├── go.sum                       # Checksums das dependências
│
├── database/                    # Camada de persistência
│   ├── db.go                   # Conexão com SQLite
│   ├── migrations.go           # Criação/alteração de tabelas
│   └── queries.go              # Queries SQL reutilizáveis
│
├── models/                      # Estruturas de dados
│   ├── task.go                 # Model de tarefas
│   ├── note.go                 # Model de notas
│   ├── event.go                # Model de eventos
│   ├── category.go             # Model de categorias
│   └── settings.go             # Model de configurações
│
├── services/                    # Lógica de negócio
│   ├── task_service.go         # Regras de tarefas
│   ├── note_service.go         # Regras de notas
│   ├── event_service.go        # Regras de eventos
│   └── file_service.go         # Manipulação de arquivos
│
├── handlers/                    # Handlers HTTP (se necessário)
│   └── api_handlers.go         # Endpoints REST (futuro)
│
├── utils/                       # Utilitários
│   ├── logger.go               # Sistema de logs
│   ├── validator.go            # Validações
│   └── helpers.go              # Funções auxiliares
│
├── frontend/                    # Aplicação React
│   ├── src/
│   │   ├── main.tsx            # Entry point React
│   │   ├── App.tsx             # Componente raiz
│   │   │
│   │   ├── components/         # Componentes React
│   │   │   ├── Layout/
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   ├── Header.tsx
│   │   │   │   └── MainContent.tsx
│   │   │   │
│   │   │   ├── Dashboard/
│   │   │   │   ├── Dashboard.tsx
│   │   │   │   ├── StatsCard.tsx
│   │   │   │   └── QuickActions.tsx
│   │   │   │
│   │   │   ├── Tasks/
│   │   │   │   ├── TaskList.tsx
│   │   │   │   ├── TaskItem.tsx
│   │   │   │   ├── TaskForm.tsx
│   │   │   │   └── TaskFilters.tsx
│   │   │   │
│   │   │   ├── Notes/
│   │   │   │   ├── NoteList.tsx
│   │   │   │   ├── NoteEditor.tsx
│   │   │   │   └── NotePreview.tsx
│   │   │   │
│   │   │   ├── Calendar/
│   │   │   │   ├── Calendar.tsx
│   │   │   │   └── EventModal.tsx
│   │   │   │
│   │   │   └── Common/
│   │   │       ├── Button.tsx
│   │   │       ├── Input.tsx
│   │   │       ├── Modal.tsx
│   │   │       └── Toast.tsx
│   │   │
│   │   ├── hooks/              # Custom React Hooks
│   │   │   ├── useTasks.ts
│   │   │   ├── useNotes.ts
│   │   │   ├── useEvents.ts
│   │   │   └── useTheme.ts
│   │   │
│   │   ├── context/            # Context API
│   │   │   ├── ThemeContext.tsx
│   │   │   └── AppContext.tsx
│   │   │
│   │   ├── services/           # API calls (Wails bindings)
│   │   │   ├── taskService.ts
│   │   │   ├── noteService.ts
│   │   │   └── eventService.ts
│   │   │
│   │   ├── types/              # TypeScript types
│   │   │   ├── task.ts
│   │   │   ├── note.ts
│   │   │   └── event.ts
│   │   │
│   │   ├── utils/              # Utilitários frontend
│   │   │   ├── formatters.ts
│   │   │   └── validators.ts
│   │   │
│   │   └── styles/             # Estilos globais
│   │       ├── globals.css
│   │       └── themes.css
│   │
│   ├── public/                 # Assets estáticos
│   │   └── logo.png
│   │
│   ├── index.html              # HTML base
│   ├── package.json            # Dependências npm
│   ├── tsconfig.json           # Config TypeScript
│   └── vite.config.ts          # Config Vite
│
├── build/                       # Arquivos de build
│   ├── appicon.png             # Ícone do app
│   ├── darwin/                 # Build macOS
│   ├── windows/                # Build Windows
│   └── bin/                    # Executáveis compilados
│
└── docs/                        # Documentação
    ├── ROADMAP.md
    ├── ARCHITECTURE.md
    ├── DATABASE.md
    └── CONTRIBUTING.md
```

---

## 🔄 Camadas da Aplicação

### 1️⃣ Camada de Apresentação (Frontend)

**Responsabilidade:** Interface do usuário e interação

```tsx
// Exemplo de componente
import { CreateTask, GetTasks } from '../wailsjs/go/main/App';

function TaskList() {
  const [tasks, setTasks] = useState([]);

  useEffect(() => {
    loadTasks();
  }, []);

  async function loadTasks() {
    const result = await GetTasks();
    setTasks(result);
  }

  return (
    <div>
      {tasks.map(task => (
        <TaskItem key={task.id} task={task} />
      ))}
    </div>
  );
}
```

### 2️⃣ Camada de Lógica (Services)

**Responsabilidade:** Regras de negócio, validações

```go
// services/task_service.go
package services

type TaskService struct {
    db *database.DB
}

func (s *TaskService) CreateTask(task models.Task) error {
    // Validações
    if task.Title == "" {
        return errors.New("título é obrigatório")
    }
    
    // Lógica de negócio
    task.CreatedAt = time.Now()
    task.Status = "pending"
    
    // Persistência
    return s.db.CreateTask(task)
}
```

### 3️⃣ Camada de Dados (Database)

**Responsabilidade:** Acesso ao banco de dados

```go
// database/db.go
package database

func (db *DB) CreateTask(task models.Task) error {
    query := `
        INSERT INTO tasks (title, description, status, priority, due_date)
        VALUES (?, ?, ?, ?, ?)
    `
    _, err := db.conn.Exec(query, 
        task.Title, 
        task.Description, 
        task.Status, 
        task.Priority, 
        task.DueDate,
    )
    return err
}
```

### 4️⃣ Camada de Bindings (Wails)

**Responsabilidade:** Expor funções Go para JavaScript

```go
// app.go
package main

type App struct {
    ctx         context.Context
    taskService *services.TaskService
}

// CreateTask é exposto automaticamente para o frontend
func (a *App) CreateTask(task models.Task) error {
    return a.taskService.CreateTask(task)
}

// GetTasks é exposto automaticamente para o frontend
func (a *App) GetTasks() ([]models.Task, error) {
    return a.taskService.GetAllTasks()
}
```

---

## 🔁 Fluxo de Dados

### Criando uma Tarefa (exemplo completo)

```
1. USUÁRIO clica em "Nova Tarefa"
   │
   ▼
2. REACT exibe modal com formulário
   │
   ▼
3. USUÁRIO preenche e clica "Salvar"
   │
   ▼
4. REACT valida campos localmente
   │
   ▼
5. REACT chama: CreateTask(task)
   │
   │  ┌──────────────────────────────────┐
   │  │ Wails Bridge (automático)        │
   │  │ Serializa JSON → Go              │
   │  └──────────────────────────────────┘
   │
   ▼
6. GO (app.go) recebe a chamada
   │
   ▼
7. GO delega para TaskService
   │
   ▼
8. TaskService valida regras de negócio
   │
   ▼
9. TaskService chama Database.CreateTask()
   │
   ▼
10. DATABASE executa INSERT no SQLite
   │
   ▼
11. RETORNO: sucesso ou erro
   │
   │  ┌──────────────────────────────────┐
   │  │ Wails Bridge (automático)        │
   │  │ Go → JSON serializado            │
   │  └──────────────────────────────────┘
   │
   ▼
12. REACT recebe resposta
   │
   ▼
13. REACT atualiza UI (adiciona tarefa na lista)
   │
   ▼
14. USUÁRIO vê a nova tarefa
```

---

## 🔗 Comunicação Frontend ↔ Backend

### Wails Bindings Automáticos

O Wails **gera automaticamente** bindings TypeScript para todas as funções Go exportadas:

```go
// app.go (Go)
func (a *App) GetTasks() ([]models.Task, error) {
    return a.taskService.GetAllTasks()
}
```

↓ **Wails gera automaticamente** ↓

```typescript
// wailsjs/go/main/App.ts (gerado automaticamente)
export function GetTasks(): Promise<models.Task[]> {
  return window['go']['main']['App']['GetTasks']();
}
```

### Uso no React

```typescript
// frontend/src/hooks/useTasks.ts
import { GetTasks, CreateTask, DeleteTask } from '../../wailsjs/go/main/App';
import type { models } from '../../wailsjs/go/models';

export function useTasks() {
  const [tasks, setTasks] = useState<models.Task[]>([]);
  
  async function loadTasks() {
    const result = await GetTasks();
    setTasks(result);
  }
  
  async function addTask(task: models.Task) {
    await CreateTask(task);
    await loadTasks(); // Recarrega lista
  }
  
  return { tasks, loadTasks, addTask };
}
```

---

## 💾 Persistência de Dados

### SQLite Embarcado

- **Arquivo:** `cockpit.db` (na pasta do executável)
- **Driver:** `github.com/mattn/go-sqlite3`
- **Conexão:** Pool de conexões gerenciado
- **Transações:** Suporte completo

### Localização do Banco

```go
// Caminho padrão do banco de dados
// Windows: C:\Users\{user}\AppData\Roaming\Personal Cockpit\cockpit.db
// macOS: ~/Library/Application Support/Personal Cockpit/cockpit.db
// Linux: ~/.config/personal-cockpit/cockpit.db

func getDatabasePath() string {
    configDir, _ := os.UserConfigDir()
    appDir := filepath.Join(configDir, "Personal Cockpit")
    os.MkdirAll(appDir, 0755)
    return filepath.Join(appDir, "cockpit.db")
}
```

### Migrations

```go
// database/migrations.go
func (db *DB) RunMigrations() error {
    migrations := []string{
        createTasksTable,
        createNotesTable,
        createEventsTable,
        createCategoriesTable,
        createSettingsTable,
    }
    
    for _, migration := range migrations {
        if _, err := db.conn.Exec(migration); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 🎯 Decisões Arquiteturais

### Por que Wails?

✅ **Leveza:** App final ~15MB vs Electron ~120MB  
✅ **Performance:** Go é compilado, muito rápido  
✅ **Type Safety:** Go + TypeScript  
✅ **Bindings Automáticos:** Sem boilerplate  
✅ **WebView Nativo:** Usa o do sistema operacional  

### Por que SQLite?

✅ **Embarcado:** Sem necessidade de servidor  
✅ **Zero Configuração:** Funciona out-of-the-box  
✅ **Leve:** ~600KB  
✅ **Confiável:** Usado por bilhões de dispositivos  
✅ **ACID Compliant:** Transações seguras  

### Por que React?

✅ **Ecossistema Rico:** Milhares de bibliotecas  
✅ **Component Based:** Reutilização de código  
✅ **Virtual DOM:** Performance  
✅ **TypeScript Support:** Excelente  
✅ **Comunidade Grande:** Fácil encontrar ajuda  

### Por que TypeScript?

✅ **Type Safety:** Menos bugs  
✅ **IntelliSense:** Melhor DX  
✅ **Refactoring:** Mais seguro  
✅ **Documentação:** Tipos servem como docs  

---

## 🔐 Segurança

### Dados Locais

- ✅ Dados armazenados apenas localmente
- ✅ Nenhuma comunicação com internet (a não ser que explicitamente implementado)
- ✅ Sem telemetria ou analytics

### Futuras Implementações

- [ ] Criptografia de banco de dados (opcional)
- [ ] Senha para acesso ao app (opcional)
- [ ] Backup criptografado

---

## 📊 Performance

### Targets

- **Inicialização:** < 2 segundos
- **Operações CRUD:** < 100ms
- **Tamanho do executável:** < 20MB
- **Uso de memória:** < 100MB em idle
- **Frame rate UI:** 60fps

### Otimizações Planejadas

- [ ] Lazy loading de componentes React
- [ ] Virtual scrolling para listas grandes
- [ ] Índices no SQLite para queries frequentes
- [ ] Cache de queries comuns
- [ ] Debounce em buscas
- [ ] Paginação de resultados

---

## 🧪 Testes

### Backend (Go)

```bash
# Testes unitários
go test ./...

# Coverage
go test -cover ./...
```

### Frontend (React)

```bash
# Testes com Vitest
npm test

# Coverage
npm run test:coverage
```

---

## 📚 Referências

- [Wails Documentation](https://wails.io/docs/introduction)
- [Go SQLite3](https://github.com/mattn/go-sqlite3)
- [React Documentation](https://react.dev)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

---

**Última revisão:** 26/12/2025