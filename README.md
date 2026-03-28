# Personal Cockpit

> Centro de comando pessoal para gestão da sua casa, rotinas e desenvolvimento — tudo local, tudo seu.

[![Wails](https://img.shields.io/badge/Wails-v2.11.0-blue)](https://wails.io)
[![React](https://img.shields.io/badge/React-18-61dafb)](https://reactjs.org)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8)](https://golang.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-4.6-3178c6)](https://www.typescriptlang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## Sobre o Projeto

**Personal Cockpit** é um aplicativo desktop nativo construído com **Wails** (Go + React). Funciona como painel central para gestão pessoal e doméstica: tarefas, notas, agenda, listas de compras, controle financeiro, rastreamento de hábitos e muito mais — tudo armazenado localmente no seu computador, sem nuvem e sem dependências externas.

### Por que este projeto existe?

A maioria dos apps de produtividade é fragmentada (um app pra tarefas, outro pra finanças, outro pra hábitos), exige internet ou coleta seus dados. O Personal Cockpit reúne tudo em um só lugar, leve, rápido e 100% privado.

### Princípios

- **Local-first**: Seus dados ficam no seu computador, ponto final
- **Privacidade total**: Nenhuma telemetria, nenhum login, nenhuma sincronização em nuvem
- **Leve e rápido**: ~15MB, inicialização < 2s
- **Multiplataforma**: Windows, macOS e Linux com o mesmo código

---

## Estado Atual

### v1.0 — Fundação (Entregue)

A base do app está completa e funcional:

| Módulo | Status | Descrição |
|--------|--------|-----------|
| **Dashboard** | Entregue | Visão geral com estatísticas, tarefas pendentes e eventos do dia |
| **Tarefas** | Entregue | CRUD completo com prioridades, categorias, status e datas |
| **Notas** | Entregue | Notas com título, conteúdo, favoritos e busca |
| **Calendário** | Entregue | Eventos com data/hora, cor, local e lembretes |
| **Categorias** | Entregue | Organização de tarefas e notas por categoria |
| **Tema Escuro/Claro** | Entregue | Preferência persistida no banco |

### Próximas Versões

**v2.0 — Gestão de Casa e Rotinas**
- Lista de Compras
- Rastreador de Hábitos
- Pomodoro Timer
- Projetos com subtarefas

**v3.0 — Controle Financeiro e Mais**
- Controle de despesas e receitas
- Orçamento mensal e alertas
- Diário/Journal
- Backup e exportação de dados

Ver planejamento completo: [docs/ROADMAP.md](docs/ROADMAP.md)

---

## Tecnologias

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| **Wails** | 2.11.0 | Framework desktop (Go + WebView) |
| **Go** | 1.24+ | Backend, lógica de negócio, SQLite |
| **React** | 18.2 | Interface do usuário |
| **TypeScript** | 4.6 | Tipagem no frontend |
| **Vite** | 3.x | Build tool |
| **SQLite** | 3.x | Banco de dados local (via modernc.org/sqlite) |

Arquitetura detalhada: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## Começando

### Pré-requisitos

- [Go](https://golang.org/dl/) 1.24+
- [Node.js](https://nodejs.org/) 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2.11.0+

```bash
# Instalar Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verificar instalação
wails doctor
```

### Rodando em Desenvolvimento

```bash
# Clone o repositório
git clone https://github.com/GustavoLemosUp/personal-cockpit.git
cd personal-cockpit

# Instale dependências do frontend
cd frontend && npm install && cd ..

# Inicie o app em modo desenvolvimento
wails dev
```

O app abre com hot-reload: mudanças no frontend recarregam automaticamente. Mudanças no backend recompilam o binário Go.

### Build para Produção

```bash
wails build
# Executável gerado em: build/bin/
```

---

## Estrutura do Projeto

```
personal-cockpit/
├── main.go              # Entry point
├── app.go               # Struct App — bindings expostos ao frontend
├── version.go           # Versão do app
├── wails.json           # Configuração do Wails
│
├── database/
│   ├── db.go            # Conexão SQLite e configuração de pragmas
│   └── migrations.go    # Schema e versionamento do banco
│
├── models/
│   ├── task.go          # Struct Task
│   ├── note.go          # Struct Note
│   ├── event.go         # Struct Event
│   └── category.go      # Struct Category
│
├── services/
│   ├── task.go          # Lógica de negócio de tarefas
│   ├── note.go          # Lógica de negócio de notas
│   ├── event.go         # Lógica de negócio de eventos
│   └── category.go      # Lógica de negócio de categorias
│
├── frontend/
│   └── src/
│       ├── App.tsx               # Roteamento entre páginas
│       ├── components/
│       │   ├── Layout.tsx        # Layout com sidebar
│       │   ├── Sidebar.tsx       # Menu de navegação
│       │   ├── Dashboard.tsx     # Página inicial
│       │   ├── tasks/            # Componentes de tarefas
│       │   ├── notes/            # Componentes de notas
│       │   ├── events/           # Componentes de eventos
│       │   └── categories/       # Componentes de categorias
│       ├── hooks/
│       │   └── useTheme.ts       # Hook de tema
│       └── styles/               # CSS modular por página
│
├── build/               # Artefatos de build (Windows, macOS)
└── docs/                # Documentação técnica
    ├── ARCHITECTURE.md
    ├── DATABASE.md
    ├── ROADMAP.md
    └── CONTRIBUTING.md
```

---

## Como Contribuir

Contribuições são bem-vindas! Veja o guia completo: [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md)

### Quick Start para contribuidores

```bash
# Fork o projeto e clone seu fork
git clone https://github.com/SEU-USUARIO/personal-cockpit.git

# Crie uma branch descritiva
git checkout -b feature/lista-de-compras

# Faça suas mudanças, então commit
git commit -m "feat(shopping): adiciona módulo de lista de compras"

# Push e abra um Pull Request
git push origin feature/lista-de-compras
```

---

## Armazenamento de Dados

Os dados ficam num banco SQLite local:

| Sistema | Localização |
|---------|------------|
| Windows | `C:\Users\{user}\AppData\Roaming\Personal Cockpit\cockpit.db` |
| macOS | `~/Library/Application Support/Personal Cockpit/cockpit.db` |
| Linux | `~/.config/Personal Cockpit/cockpit.db` |

Schema completo: [docs/DATABASE.md](docs/DATABASE.md)

---

## Licença

MIT — veja [LICENSE](LICENSE).

---

## Autor

**Gustavo Lemos** — [@GustavoLemosUp](https://github.com/GustavoLemosUp)

---

<p align="center">Feito com Go, React e muita determinação.</p>
