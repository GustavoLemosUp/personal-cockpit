# 🗺️ Roadmap - Personal Cockpit

> Planejamento completo do desenvolvimento do Personal Cockpit

**Última atualização:** 26/12/2025  
**Status atual:** 🟢 Fase 0 - Configuração concluída

---

## 📊 Visão Geral

```
┌─────────────┬─────────────┬─────────────┬─────────────┐
│   FASE 0    │   FASE 1    │   FASE 2    │   FASE 3    │
│  Setup ✅   │  MVP 🔄     │  Avançado   │   Futuro    │
│  1 semana   │  6 semanas  │  4 semanas  │   TBD       │
└─────────────┴─────────────┴─────────────┴─────────────┘
```

---

## ✅ FASE 0: Configuração do Ambiente (CONCLUÍDA)

**Período:** 26/12/2025  
**Status:** ✅ **Concluído**

### Objetivos
- [x] Instalar Go 1.23+
- [x] Instalar Node.js 18+
- [x] Instalar Wails CLI v2.11.0
- [x] Criar projeto base com template React + TypeScript
- [x] Configurar Git e GitHub
- [x] Estruturar documentação inicial

### Resultados
- ✅ Projeto inicializado com Wails
- ✅ Frontend React + Vite configurado
- ✅ Git configurado com .gitignore adequado
- ✅ Primeiro commit realizado
- ✅ Repositório no GitHub criado

---

## 🎯 FASE 1: MVP - Funcionalidades Essenciais

**Período:** 6 semanas (Janeiro - Fevereiro 2026)  
**Status:** 🔄 **Planejado**

### Semana 1-2: Banco de Dados e Backend Base

#### Objetivos
- [ ] Configurar SQLite local
- [ ] Criar schema do banco de dados
- [ ] Implementar camada de persistência
- [ ] Criar estruturas de dados (models)

#### Entregas
```go
// Tabelas a serem criadas:
- tasks          // Tarefas
- categories     // Categorias
- notes          // Notas
- events         // Eventos do calendário
- settings       // Configurações do app
```

#### Arquivos a criar
- `database/db.go` - Conexão com SQLite
- `database/migrations.go` - Criação de tabelas
- `models/task.go` - Model de tarefas
- `models/category.go` - Model de categorias
- `models/note.go` - Model de notas
- `models/event.go` - Model de eventos

#### Critérios de Sucesso
- ✅ Banco de dados SQLite criado automaticamente
- ✅ Todas as tabelas criadas com relacionamentos
- ✅ CRUD básico funcionando para todas entidades
- ✅ Testes unitários para camada de dados

---

### Semana 3: Sistema de Tarefas (Backend)

#### Objetivos
- [ ] Implementar CRUD completo de tarefas
- [ ] Sistema de categorias/tags
- [ ] Sistema de prioridades
- [ ] Filtros e buscas

#### Funcionalidades - Tarefas
```
✓ Criar tarefa
✓ Editar tarefa
✓ Deletar tarefa
✓ Marcar como concluída
✓ Definir prioridade (Alta/Média/Baixa)
✓ Adicionar categoria
✓ Definir data de vencimento
✓ Buscar por título/descrição
✓ Filtrar por status/prioridade/categoria
```

#### API (Bindings Go → React)
```go
func (a *App) CreateTask(task Task) error
func (a *App) GetTasks(filter TaskFilter) ([]Task, error)
func (a *App) UpdateTask(id int, task Task) error
func (a *App) DeleteTask(id int) error
func (a *App) ToggleTaskStatus(id int) error
func (a *App) GetTasksByCategory(categoryId int) ([]Task, error)
```

#### Critérios de Sucesso
- ✅ Todas as operações CRUD funcionando
- ✅ Filtros retornando resultados corretos
- ✅ Validações de dados implementadas
- ✅ Tratamento de erros adequado

---

### Semana 4-5: Interface de Tarefas (Frontend)

#### Objetivos
- [ ] Layout principal com sidebar
- [ ] Lista de tarefas com cards
- [ ] Formulário criar/editar tarefa
- [ ] Sistema de filtros e busca
- [ ] Animações e transições suaves

#### Componentes React a criar
```tsx
frontend/src/components/
├── Layout/
│   ├── Sidebar.tsx        // Menu lateral
│   ├── Header.tsx         // Cabeçalho
│   └── MainContent.tsx    // Container principal
├── Tasks/
│   ├── TaskList.tsx       // Lista de tarefas
│   ├── TaskItem.tsx       // Card individual
│   ├── TaskForm.tsx       // Formulário
│   ├── TaskFilters.tsx    // Filtros
│   └── TaskStats.tsx      // Estatísticas
└── Common/
    ├── Button.tsx
    ├── Input.tsx
    ├── Select.tsx
    └── Modal.tsx
```

#### Features de UI
```
✓ Drag and drop para reordenar tarefas
✓ Checkbox animado para marcar como concluída
✓ Badges coloridos para prioridades
✓ Ícones para categorias
✓ Contador de tarefas pendentes
✓ Filtro rápido por status
✓ Busca em tempo real
✓ Animação de loading
✓ Toast notifications para ações
```

#### Critérios de Sucesso
- ✅ Interface responsiva e fluida
- ✅ Todas as ações funcionando
- ✅ Feedback visual para o usuário
- ✅ Validação de formulários
- ✅ Acessibilidade (keyboard navigation)

---

### Semana 6: Notas Rápidas

#### Backend - Objetivos
- [ ] CRUD de notas
- [ ] Sistema de categorização
- [ ] Busca full-text
- [ ] Timestamps (criado/editado)

#### Frontend - Objetivos
- [ ] Editor de texto rico (ou Markdown)
- [ ] Lista de notas
- [ ] Preview de notas
- [ ] Categorias de notas

#### Componentes
```tsx
frontend/src/components/Notes/
├── NoteList.tsx       // Lista lateral
├── NoteEditor.tsx     // Editor principal
├── NotePreview.tsx    // Preview da nota
└── NoteCategories.tsx // Gerenciar categorias
```

#### Critérios de Sucesso
- ✅ Criar, editar, deletar notas
- ✅ Editor funcional (texto ou Markdown)
- ✅ Busca por conteúdo
- ✅ Auto-save (salvar automaticamente)

---

### Semana 7: Calendário Básico

#### Objetivos
- [ ] Implementar CRUD de eventos
- [ ] Visualização mensal
- [ ] Criar/editar eventos
- [ ] Notificações básicas

#### Features
```
✓ Visualização mensal (calendário)
✓ Adicionar evento com data/hora
✓ Editar evento
✓ Deletar evento
✓ Cores para tipos de eventos
✓ Lista de eventos do dia
```

#### Biblioteca Sugerida
- **react-big-calendar** ou **FullCalendar**

#### Critérios de Sucesso
- ✅ Calendário visual funcionando
- ✅ CRUD de eventos completo
- ✅ Navegação entre meses
- ✅ Eventos aparecem nas datas corretas

---

### Semana 8: Dashboard + Polish

#### Dashboard - Objetivos
- [ ] Visão geral de tudo
- [ ] Cards com estatísticas
- [ ] Próximas tarefas
- [ ] Eventos hoje
- [ ] Gráficos simples

#### Dashboard - Componentes
```tsx
frontend/src/components/Dashboard/
├── Dashboard.tsx          // Container principal
├── StatsCard.tsx         // Card de estatística
├── UpcomingTasks.tsx     // Próximas tarefas
├── TodayEvents.tsx       // Eventos de hoje
├── ProductivityChart.tsx // Gráfico de produtividade
└── QuickActions.tsx      // Ações rápidas
```

#### Informações no Dashboard
```
📊 Estatísticas:
- Total de tarefas
- Tarefas concluídas hoje
- Tarefas pendentes
- Eventos da semana

📋 Próximas Tarefas (5 mais próximas do vencimento)
📅 Eventos de Hoje
📈 Gráfico de tarefas concluídas (últimos 7 dias)
⚡ Ações Rápidas (+ Nova Tarefa, + Nota, + Evento)
```

#### Tema Claro/Escuro
- [ ] Sistema de temas
- [ ] Toggle de tema
- [ ] Persistir preferência
- [ ] Cores otimizadas para cada tema

#### Polish Geral
- [ ] Animações suaves
- [ ] Transições entre páginas
- [ ] Loading states
- [ ] Empty states (quando não há dados)
- [ ] Tratamento de erros com UI amigável
- [ ] Atalhos de teclado

#### Critérios de Sucesso - MVP Completo
- ✅ Dashboard mostrando resumo geral
- ✅ To-Do List totalmente funcional
- ✅ Notas funcionando
- ✅ Calendário básico operacional
- ✅ Tema claro/escuro
- ✅ App estável e usável
- ✅ **Primeira versão pronta para uso! 🎉**

---

## 🚀 FASE 2: Funcionalidades Avançadas

**Período:** 4 semanas (Março 2026)  
**Status:** 📋 **Planejado**

### Semana 9-10: Gerenciador de Arquivos

#### Objetivos
- [ ] Upload de arquivos via drag & drop
- [ ] Categorização de arquivos
- [ ] Sistema de tags
- [ ] Preview de PDFs e imagens
- [ ] Busca por nome/tipo

#### Funcionalidades
```
✓ Arrastar arquivos para o app
✓ Organizar em categorias
✓ Adicionar tags
✓ Ver preview (PDF, imagens)
✓ Abrir arquivo no programa padrão
✓ Deletar arquivos
✓ Buscar arquivos
```

---

### Semana 11: Sistema de Projetos

#### Objetivos
- [ ] Criar projetos
- [ ] Agrupar tarefas em projetos
- [ ] Timeline de projeto
- [ ] Progresso visual
- [ ] Subtarefas

#### Features
```
✓ Criar/editar/deletar projetos
✓ Adicionar tarefas ao projeto
✓ Ver progresso (% concluído)
✓ Subtarefas (tarefas dentro de tarefas)
✓ Deadline do projeto
✓ Membros do projeto (futuro)
```

---

### Semana 12: Rastreador de Hábitos

#### Objetivos
- [ ] Criar hábitos diários
- [ ] Marcar como concluído por dia
- [ ] Streaks (sequências)
- [ ] Gráficos de progresso
- [ ] Notificações de lembrete

#### Features
```
✓ Criar hábito (ex: "Exercício", "Ler")
✓ Marcar conclusão diária
✓ Ver histórico (calendário)
✓ Streak atual (dias consecutivos)
✓ Melhor streak
✓ Gráfico mensal
```

---

### Semana 13: Pomodoro Timer

#### Objetivos
- [ ] Timer configurável
- [ ] Sessões de foco
- [ ] Pausas automáticas
- [ ] Histórico de sessões
- [ ] Integração com tarefas

#### Features
```
✓ Configurar tempo de foco (25min padrão)
✓ Pausas curtas (5min)
✓ Pausas longas (15min)
✓ Notificação sonora/visual
✓ Vincular sessão a uma tarefa
✓ Histórico de sessões Pomodoro
✓ Total de tempo focado no dia/semana
```

---

## 💡 FASE 3: Funcionalidades Futuras

**Período:** TBD  
**Status:** 💭 **Ideias**

### Controle Financeiro Básico
```
- Registrar despesas/receitas
- Categorias financeiras
- Relatórios mensais
- Gráficos de gastos
- Orçamento mensal
- Alertas de limite
```

### Diário/Journal
```
- Entradas diárias
- Editor rico
- Anexar fotos
- Tags de humor/emoções
- Busca por período
- Criptografia (opcional)
```

### Backup e Sincronização
```
- Exportar todos os dados
- Importar dados
- Backup automático local
- Sincronização via cloud (opcional)
- Versionamento de backup
```

### Sistema de Plugins
```
- API para extensões
- Marketplace de plugins
- Temas customizados
- Integrações externas
```

---

## 📈 Métricas de Sucesso

### MVP (v1.0)
- [ ] App abre em < 2 segundos
- [ ] Operações CRUD em < 100ms
- [ ] Zero crashes em uso normal
- [ ] Interface responsiva (60fps)
- [ ] Executável < 20MB

### v2.0
- [ ] Suporte a 1000+ tarefas sem lag
- [ ] Upload de arquivos < 50MB
- [ ] Sincronização em < 5 segundos
- [ ] Taxa de satisfação > 90%

---

## 🎯 Próximos Passos Imediatos

### Esta Semana (26/12 - 01/01)
1. ✅ Finalizar documentação
2. [ ] Estudar SQLite em Go
3. [ ] Criar protótipo de tela (Figma/papel)
4. [ ] Definir paleta de cores e design system

### Próxima Semana (02/01 - 08/01)
1. [ ] Implementar camada de banco de dados
2. [ ] Criar migrations
3. [ ] Implementar models
4. [ ] Testes unitários da camada de dados

---

## 🔄 Processo de Desenvolvimento

### Workflow
```
1. Criar branch: feature/nome-da-feature
2. Desenvolver funcionalidade
3. Testar localmente
4. Commit com mensagem descritiva
5. Push para GitHub
6. Merge para main quando estável
```

### Commits Semânticos
```
✨ feat: Nova funcionalidade
🐛 fix: Correção de bug
📝 docs: Documentação
💄 style: UI/CSS
♻️ refactor: Refatoração
✅ test: Testes
🔧 chore: Configuração
```

---

## 📞 Suporte e Feedback

Encontrou um bug? Tem uma sugestão?  
Abra uma **issue** no GitHub: [Issues](https://github.com/GustavoLemosUp/personal-cockpit/issues)

---

**Última revisão:** 26/12/2025  
**Próxima revisão:** Semanalmente