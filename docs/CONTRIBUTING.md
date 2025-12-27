# 🤝 Guia de Contribuição - Personal Cockpit

> Como contribuir para o desenvolvimento do Personal Cockpit

**Bem-vindo!** Ficamos felizes que você queira contribuir com o projeto! 🎉

---

## 📋 Índice

1. [Como Começar](#como-começar)
2. [Configurando o Ambiente](#configurando-o-ambiente)
3. [Estrutura do Projeto](#estrutura-do-projeto)
4. [Workflow de Desenvolvimento](#workflow-de-desenvolvimento)
5. [Padrões de Código](#padrões-de-código)
6. [Commits Semânticos](#commits-semânticos)
7. [Pull Requests](#pull-requests)
8. [Testes](#testes)
9. [Reportando Bugs](#reportando-bugs)
10. [Sugerindo Features](#sugerindo-features)

---

## 🚀 Como Começar

### Encontrando algo para fazer

1. **Issues abertas:** Veja as [issues abertas](https://github.com/GustavoLemosUp/personal-cockpit/issues)
2. **Good First Issue:** Procure por issues marcadas como `good first issue`
3. **Help Wanted:** Issues marcadas como `help wanted` precisam de ajuda
4. **Roadmap:** Consulte o [ROADMAP.md](ROADMAP.md) para ver o que está planejado

### Tipos de contribuição

- 🐛 **Bug Fixes**: Corrigir bugs reportados
- ✨ **Features**: Implementar novas funcionalidades
- 📝 **Documentação**: Melhorar ou traduzir docs
- 🎨 **UI/UX**: Melhorias de interface
- ✅ **Testes**: Adicionar ou melhorar testes
- ♻️ **Refatoração**: Melhorar código existente

---

## 💻 Configurando o Ambiente

### Pré-requisitos

- **Go:** 1.23 ou superior ([Download](https://golang.org/dl/))
- **Node.js:** 18 ou superior ([Download](https://nodejs.org/))
- **Wails CLI:** v2.11.0 ou superior
- **Git:** Para versionamento
- **Editor:** VS Code recomendado

### Instalação do Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verificar instalação
wails version

# Verificar dependências
wails doctor
```

### Clone do Repositório

```bash
# Clone o repositório
git clone https://github.com/GustavoLemosUp/personal-cockpit.git
cd personal-cockpit

# Instale dependências do frontend
cd frontend
npm install
cd ..

# Rode em modo desenvolvimento
wails dev
```

### Configuração do VS Code (Recomendado)

**Extensões recomendadas:**

```json
{
  "recommendations": [
    "golang.go",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "bradlc.vscode-tailwindcss"
  ]
}
```

**settings.json:**

```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  }
}
```

---

## 📂 Estrutura do Projeto

```
personal-cockpit/
├── main.go              # Entry point
├── app.go               # Bindings principais
├── database/            # Camada de dados
├── models/              # Estruturas de dados
├── services/            # Lógica de negócio
├── utils/               # Utilitários
└── frontend/            # App React
    ├── src/
    │   ├── components/  # Componentes React
    │   ├── hooks/       # Custom hooks
    │   ├── services/    # API calls
    │   └── types/       # TypeScript types
    └── package.json
```

**Leia mais:** [ARCHITECTURE.md](ARCHITECTURE.md)

---

## 🔄 Workflow de Desenvolvimento

### 1. Criar uma Branch

```bash
# Sempre crie uma branch a partir da main
git checkout main
git pull origin main

# Crie uma branch com nome descritivo
git checkout -b feature/nome-da-feature
# ou
git checkout -b fix/nome-do-bug
```

### Nomenclatura de Branches

| Tipo | Prefixo | Exemplo |
|------|---------|---------|
| Nova feature | `feature/` | `feature/task-drag-drop` |
| Bug fix | `fix/` | `fix/task-delete-error` |
| Documentação | `docs/` | `docs/update-readme` |
| Refatoração | `refactor/` | `refactor/task-service` |
| Testes | `test/` | `test/add-task-tests` |

### 2. Desenvolver

```bash
# Rode o app em modo dev
wails dev

# O app vai recarregar automaticamente quando você salvar arquivos
```

### 3. Testar

```bash
# Testes do backend (Go)
go test ./...

# Testes do frontend (React)
cd frontend
npm test
```

### 4. Commit

```bash
# Adicione os arquivos alterados
git add .

# Faça commit com mensagem semântica
git commit -m "✨ feat: adiciona drag and drop em tarefas"
```

### 5. Push

```bash
# Envie para o GitHub
git push origin feature/nome-da-feature
```

### 6. Pull Request

1. Vá para o GitHub
2. Clique em "Compare & pull request"
3. Preencha o template de PR
4. Aguarde review

---

## 📝 Padrões de Código

### Backend (Go)

#### Formatação

```bash
# Formatar código
gofmt -w .

# Ou usar goimports (recomendado)
goimports -w .
```

#### Convenções

```go
// ✅ BOM - Nome de função exportada em PascalCase
func CreateTask(task Task) error {
    // ...
}

// ✅ BOM - Nome de função privada em camelCase
func validateTask(task Task) error {
    // ...
}

// ✅ BOM - Comentário antes de função exportada
// CreateTask cria uma nova tarefa no banco de dados.
// Retorna erro se a validação falhar.
func CreateTask(task Task) error {
    // ...
}

// ❌ RUIM - Função exportada sem comentário
func CreateTask(task Task) error {
    // ...
}
```

#### Error Handling

```go
// ✅ BOM - Sempre checar erros
task, err := s.GetTask(id)
if err != nil {
    return fmt.Errorf("failed to get task: %w", err)
}

// ❌ RUIM - Ignorar erros
task, _ := s.GetTask(id)
```

#### Estruturas

```go
// ✅ BOM - Usar tags JSON
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### Frontend (React/TypeScript)

#### Formatação

```bash
# Formatar com Prettier
npm run format

# Lint
npm run lint
```

#### Convenções

```tsx
// ✅ BOM - Componente em PascalCase
function TaskItem({ task }: TaskItemProps) {
    // ...
}

// ✅ BOM - Props interface
interface TaskItemProps {
    task: Task;
    onDelete?: (id: number) => void;
}

// ✅ BOM - Usar tipos em vez de any
function handleSubmit(data: FormData): void {
    // ...
}

// ❌ RUIM - Usar any
function handleSubmit(data: any) {
    // ...
}
```

#### Hooks

```tsx
// ✅ BOM - Custom hooks começam com "use"
function useTasks() {
    const [tasks, setTasks] = useState<Task[]>([]);
    // ...
    return { tasks, loadTasks, createTask };
}

// ✅ BOM - useEffect com dependencies corretas
useEffect(() => {
    loadTasks();
}, [loadTasks]); // Incluir dependências

// ❌ RUIM - useEffect sem dependencies
useEffect(() => {
    loadTasks();
}); // Vai executar a cada render!
```

#### Componentes

```tsx
// ✅ BOM - Componente funcional com TypeScript
interface ButtonProps {
    children: React.ReactNode;
    onClick?: () => void;
    variant?: 'primary' | 'secondary';
}

export function Button({ 
    children, 
    onClick, 
    variant = 'primary' 
}: ButtonProps) {
    return (
        <button 
            onClick={onClick}
            className={`btn btn-${variant}`}
        >
            {children}
        </button>
    );
}
```

### SQL

```sql
-- ✅ BOM - Maiúsculas para palavras-chave
SELECT id, title, status
FROM tasks
WHERE status = 'pending'
ORDER BY due_date ASC;

-- ✅ BOM - Indentação clara
SELECT 
    t.id,
    t.title,
    c.name AS category_name
FROM tasks t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.status = 'pending';
```

---

## 🎯 Commits Semânticos

Usamos **Conventional Commits** para padronizar mensagens de commit.

### Formato

```
<tipo>(<escopo>): <descrição>

[corpo opcional]

[footer opcional]
```

### Tipos

| Emoji | Tipo | Descrição | Exemplo |
|-------|------|-----------|---------|
| ✨ | `feat` | Nova funcionalidade | `✨ feat: adiciona filtro de tarefas por data` |
| 🐛 | `fix` | Correção de bug | `🐛 fix: corrige erro ao deletar tarefa` |
| 📝 | `docs` | Documentação | `📝 docs: atualiza guia de instalação` |
| 💄 | `style` | Mudanças de UI/CSS | `💄 style: melhora layout do dashboard` |
| ♻️ | `refactor` | Refatoração | `♻️ refactor: simplifica lógica de tarefas` |
| ⚡ | `perf` | Performance | `⚡ perf: otimiza query de busca` |
| ✅ | `test` | Testes | `✅ test: adiciona testes para TaskService` |
| 🔧 | `chore` | Configuração | `🔧 chore: atualiza dependências` |
| 🚀 | `build` | Build/Deploy | `🚀 build: configura CI/CD` |
| 🔥 | `remove` | Remoção de código | `🔥 remove: remove código não utilizado` |

### Exemplos

```bash
# Feature simples
git commit -m "✨ feat: adiciona botão de editar tarefa"

# Feature com escopo
git commit -m "✨ feat(tasks): implementa drag and drop"

# Bug fix
git commit -m "🐛 fix(database): corrige conexão SQLite"

# Documentação
git commit -m "📝 docs: adiciona seção de testes no README"

# Breaking change
git commit -m "💥 feat!: muda estrutura de Task no banco

BREAKING CHANGE: campo 'completed' renomeado para 'status'"
```

### Escopos comuns

- `tasks`: Sistema de tarefas
- `notes`: Sistema de notas
- `calendar`: Calendário
- `database`: Banco de dados
- `ui`: Interface geral
- `api`: Bindings Go ↔ React

---

## 🔀 Pull Requests

### Template de PR

Ao abrir um PR, preencha o template:

```markdown
## Descrição
[Descreva as mudanças feitas]

## Tipo de mudança
- [ ] Bug fix
- [ ] Nova feature
- [ ] Breaking change
- [ ] Documentação

## Como testar
1. [Passo 1]
2. [Passo 2]

## Checklist
- [ ] Código segue os padrões do projeto
- [ ] Comentei código complexo
- [ ] Atualizei a documentação
- [ ] Adicionei testes
- [ ] Todos os testes passam
- [ ] Build funciona sem erros

## Screenshots (se aplicável)
[Cole screenshots aqui]

## Issues relacionadas
Closes #123
```

### Review Process

1. **Automated Checks**: Testes e linting automáticos
2. **Code Review**: Pelo menos 1 aprovação necessária
3. **Merge**: Squash and merge na main

### Dicas para um bom PR

✅ **Faça PRs pequenos** - Mais fácil de revisar  
✅ **Um PR = Uma feature** - Não misture funcionalidades  
✅ **Escreva descrição clara** - Facilita o review  
✅ **Adicione screenshots** - Para mudanças visuais  
✅ **Atualize docs** - Se necessário  
✅ **Responda feedback** - Seja receptivo  

---

## ✅ Testes

### Backend (Go)

```go
// task_service_test.go
package services

import "testing"

func TestCreateTask(t *testing.T) {
    // Setup
    db := setupTestDB()
    service := NewTaskService(db)
    
    // Test
    task := Task{Title: "Test Task"}
    err := service.CreateTask(task)
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

**Rodar testes:**

```bash
# Todos os testes
go test ./...

# Com coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Frontend (React)

```tsx
// TaskItem.test.tsx
import { render, screen } from '@testing-library/react';
import { TaskItem } from './TaskItem';

test('renders task title', () => {
    const task = { id: 1, title: 'Test Task' };
    render(<TaskItem task={task} />);
    
    expect(screen.getByText('Test Task')).toBeInTheDocument();
});
```

**Rodar testes:**

```bash
cd frontend

# Todos os testes
npm test

# Com coverage
npm run test:coverage

# Watch mode
npm test -- --watch
```

---

## 🐛 Reportando Bugs

### Antes de reportar

1. Verifique se já existe uma issue sobre o bug
2. Teste na última versão
3. Reúna informações do erro

### Template de Bug Report

```markdown
## Descrição do Bug
[Descrição clara do que aconteceu]

## Como Reproduzir
1. Vá para '...'
2. Clique em '...'
3. Veja o erro

## Comportamento Esperado
[O que deveria acontecer]

## Screenshots
[Se aplicável]

## Ambiente
- OS: [Windows 10, macOS 13, Ubuntu 22.04]
- Versão do app: [1.0.0]
- Go version: [1.23.4]
- Node version: [18.12.0]

## Logs
```
[Cole logs de erro aqui]
```

## Informações Adicionais
[Qualquer outro contexto]
```

---

## 💡 Sugerindo Features

### Template de Feature Request

```markdown
## Descrição da Feature
[Descrição clara da funcionalidade]

## Problema que Resolve
[Qual problema esta feature resolve?]

## Solução Proposta
[Como você imagina que funcione?]

## Alternativas Consideradas
[Outras formas de resolver?]

## Informações Adicionais
[Mockups, exemplos, etc]
```

---

## 📚 Recursos Úteis

### Documentação

- [Wails Documentation](https://wails.io/docs/introduction)
- [React Documentation](https://react.dev)
- [Go Documentation](https://go.dev/doc/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

### Tutoriais Internos

- [ROADMAP.md](ROADMAP.md) - Planejamento do projeto
- [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitetura técnica
- [DATABASE.md](DATABASE.md) - Schema do banco

### Comunidade

- [Discord](https://discord.gg/wails) - Wails Community
- [GitHub Discussions](https://github.com/GustavoLemosUp/personal-cockpit/discussions)

---

## ❓ Dúvidas?

- Abra uma [Discussion](https://github.com/GustavoLemosUp/personal-cockpit/discussions)
- Ou entre em contato via [Issues](https://github.com/GustavoLemosUp/personal-cockpit/issues)

---

## 🙏 Agradecimentos

Obrigado por contribuir com o Personal Cockpit! Cada contribuição, grande ou pequena, é muito valiosa! ❤️

---

**Happy Coding! 🚀**