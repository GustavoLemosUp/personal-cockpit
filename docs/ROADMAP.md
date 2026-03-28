# Roadmap — Personal Cockpit

> Planejamento do desenvolvimento do Personal Cockpit

**Última atualização:** Março 2026
**Status atual:** v1.0 entregue — iniciando v2.0

---

## Visão do Produto

O Personal Cockpit existe para ser o centro de comando da sua vida doméstica e pessoal. Não é só um to-do list: é onde você gerencia sua casa, suas finanças, suas compras, seus hábitos e seu desenvolvimento pessoal — tudo em um app leve, sem internet e sem abrir mão dos seus dados.

**Áreas que o app cobre (atual e futura):**

| Área | v1.0 | v2.0 | v3.0 |
|------|:----:|:----:|:----:|
| Tarefas e to-do | Entregue | — | — |
| Notas | Entregue | — | — |
| Agenda / Calendário | Entregue | — | — |
| Categorias | Entregue | — | — |
| Lista de Compras | — | Planejado | — |
| Rastreador de Hábitos | — | Planejado | — |
| Pomodoro / Foco | — | Planejado | — |
| Projetos | — | Planejado | — |
| Controle Financeiro | — | — | Planejado |
| Diário / Journal | — | — | Planejado |
| Gestão de Casa (manutenção, checklists) | — | — | Planejado |
| Backup e Exportação | — | — | Planejado |

---

## v1.0 — Fundação (Entregue)

**Status:** Concluído

### O que foi entregue

#### Backend (Go + SQLite)
- Banco SQLite com migrations versionadas
- Models: Task, Note, Event, Category
- Services com CRUD completo e validações para todas as entidades
- Métodos expostos via Wails para o frontend:
  - Tarefas: criar, listar, buscar, atualizar, deletar, alternar status, filtrar
  - Notas: criar, listar, buscar, favoritar, atualizar, deletar
  - Eventos: criar, listar, listar por data, atualizar, deletar
  - Categorias: criar, listar, atualizar, deletar, filtrar por tipo
- Configuração de pragma SQLite (WAL, foreign keys, cache)
- Persistent settings via tabela `settings`

#### Frontend (React + TypeScript)
- Roteamento entre páginas via App.tsx
- Layout com Sidebar de navegação
- Dashboard com estatísticas e eventos do dia
- Página de Tarefas com formulário, listagem e filtros
- Página de Notas com busca e favoritos
- Página de Eventos com gestão de agenda
- Página de Categorias
- CSS modular por página, tema escuro como padrão
- Tema claro/escuro com persistência

---

## v2.0 — Gestão de Casa e Rotinas

**Status:** Planejado
**Período estimado:** Q2 2026

Esta versão transforma o app de um organizador pessoal para um verdadeiro gestor doméstico, adicionando as funcionalidades que tornam o dia a dia mais fácil.

### Lista de Compras

**O que resolve:** Acabar com o papel de supermercado e centralizar compras do mês.

```
Funcionalidades:
- Criar listas por categoria (Mercado, Farmácia, etc.)
- Adicionar itens com quantidade e unidade
- Marcar itens como comprados
- Histórico de itens frequentes (sugestão automática)
- Listas recorrentes (compras mensais fixas)
- Exportar lista para texto
```

**Backend:** novo modelo `ShoppingList`, `ShoppingItem`
**Frontend:** `components/shopping/`
**Banco:** tabelas `shopping_lists`, `shopping_items`

---

### Rastreador de Hábitos

**O que resolve:** Acompanhar rotinas diárias com visibilidade de progresso.

```
Funcionalidades:
- Criar hábitos com frequência (diário, dias da semana)
- Marcar conclusão por dia
- Streak atual e melhor streak
- Calendário de histórico (heat map)
- Progresso do mês em porcentagem
```

**Backend:** modelos `Habit`, `HabitLog`
**Frontend:** `components/habits/`
**Banco:** tabelas `habits`, `habit_logs`

---

### Pomodoro / Foco

**O que resolve:** Sessões de trabalho focado integradas com as tarefas do app.

```
Funcionalidades:
- Timer configurável (padrão 25/5/15min)
- Vincular sessão a uma tarefa
- Notificação visual/sonora ao fim do timer
- Histórico de sessões do dia/semana
- Total de tempo focado
```

**Frontend:** `components/pomodoro/` (componente overlay)
**Banco:** tabela `pomodoro_sessions`

---

### Projetos

**O que resolve:** Agrupar tarefas relacionadas com visão de progresso.

```
Funcionalidades:
- Criar projetos com descrição e deadline
- Associar tarefas existentes ao projeto
- Barra de progresso (% concluído)
- Subtarefas
- Visão kanban (futuro)
```

**Backend:** modelo `Project`
**Frontend:** `components/projects/`
**Banco:** tabela `projects`, relação `project_tasks`

---

## v3.0 — Controle Financeiro e Mais

**Status:** Planejado
**Período estimado:** Q3–Q4 2026

### Controle Financeiro

**O que resolve:** Acompanhar gastos, receitas e orçamento mensal sem planilhas.

```
Funcionalidades:
- Registrar despesas e receitas
- Categorias financeiras (Alimentação, Moradia, Lazer, etc.)
- Orçamento mensal por categoria
- Alertas ao atingir limite
- Relatório mensal com gráficos
- Visão de saldo atual
- Contas recorrentes (aluguel, assinaturas)
```

**Backend:** modelos `Transaction`, `FinancialCategory`, `Budget`
**Frontend:** `components/finance/`
**Banco:** tabelas `transactions`, `financial_categories`, `budgets`

---

### Gestão da Casa

**O que resolve:** Controlar manutenções, limpezas e tarefas domésticas recorrentes.

```
Funcionalidades:
- Checklists de limpeza (semanal, mensal)
- Registro de manutenções (data, custo, próxima revisão)
- Alertas de revisão (filtro de ar, revisão de carro, etc.)
- Inventário básico (eletrodomésticos, documentos)
```

---

### Diário / Journal

**O que resolve:** Registro pessoal diário com privacidade total.

```
Funcionalidades:
- Entradas diárias com editor de texto
- Tags de humor/estado
- Busca por período e conteúdo
- Criptografia opcional da entrada
```

---

### Backup e Exportação

```
Funcionalidades:
- Exportar todos os dados em JSON
- Exportar dados específicos (só finanças, só tarefas, etc.)
- Backup automático diário (para pasta local)
- Importar backup
- Exportar relatório financeiro em CSV
```

---

## Backlog / Ideias Futuras

- Sistema de plugins e extensões
- Integração com Google Calendar (importar eventos, não sincronizar)
- Modo offline-first com sync opcional via arquivo (sem cloud)
- Atalhos de teclado globais (abrir app, nova tarefa)
- CLI para adicionar tarefas por linha de comando

---

## Métricas de Qualidade (todos os releases)

| Métrica | Target |
|---------|--------|
| Inicialização | < 2 segundos |
| Operações CRUD | < 100ms |
| Tamanho do executável | < 20MB |
| Uso de memória (idle) | < 100MB |
| Zero crashes em uso normal | Obrigatório |

---

## Como Contribuir com o Roadmap

Tem uma ideia para o app? Abra uma issue no GitHub com:
- O problema que você quer resolver
- Como o app poderia ajudar
- Se você está disposto a implementar

Issues relevantes são priorizadas para o próximo milestone.

[Ver issues abertas](https://github.com/GustavoLemosUp/personal-cockpit/issues)

---

**Revisão:** Março 2026
