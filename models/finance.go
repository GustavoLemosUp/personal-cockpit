package models

import "time"

// AccountType representa o tipo de conta financeira.
type AccountType string

const (
	AccountTypeChecking AccountType = "checking"  // Conta corrente
	AccountTypeSavings  AccountType = "savings"   // Poupança
	AccountTypeWallet   AccountType = "wallet"    // Carteira / dinheiro
	AccountTypeCredit   AccountType = "credit"    // Cartão de crédito
	AccountTypeInvest   AccountType = "invest"    // Investimento
)

// Account representa uma conta financeira.
type Account struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Type        AccountType `json:"type"`
	Balance     float64     `json:"balance"`
	Currency    string      `json:"currency"`
	Color       string      `json:"color"`
	Icon        string      `json:"icon"`
	Description string      `json:"description"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// TransactionType representa o tipo de transação.
type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"   // Receita
	TransactionTypeExpense  TransactionType = "expense"  // Despesa
	TransactionTypeTransfer TransactionType = "transfer" // Transferência entre contas
)

// Transaction representa um movimento financeiro.
type Transaction struct {
	ID                int             `json:"id"`
	AccountID         int             `json:"account_id"`
	AccountName       string          `json:"account_name,omitempty"`
	ToAccountID       *int            `json:"to_account_id,omitempty"` // só para transferências
	ToAccountName     string          `json:"to_account_name,omitempty"`
	Type              TransactionType `json:"type"`
	Category          string          `json:"category"`
	Description       string          `json:"description"`
	Amount            float64         `json:"amount"`
	Date              time.Time       `json:"date"`
	IsRecurring       bool            `json:"is_recurring"`
	RecurrenceRule    string          `json:"recurrence_rule"` // monthly, weekly, etc.
	Tags              string          `json:"tags"`
	CreatedAt         time.Time       `json:"created_at"`
}

// ScheduledTransaction representa um movimento financeiro agendado para o futuro.
type ScheduledTransaction struct {
	ID             int             `json:"id"`
	AccountID      int             `json:"account_id"`
	AccountName    string          `json:"account_name,omitempty"`
	Type           TransactionType `json:"type"`
	Category       string          `json:"category"`
	Description    string          `json:"description"`
	Amount         float64         `json:"amount"`
	ScheduledAt    time.Time       `json:"scheduled_at"`
	IsRecurring    bool            `json:"is_recurring"`
	RecurrenceRule string          `json:"recurrence_rule"`
	IsExecuted     bool            `json:"is_executed"`
	CreatedAt      time.Time       `json:"created_at"`
}

// FinanceSummary é o resumo financeiro geral.
type FinanceSummary struct {
	TotalBalance    float64              `json:"total_balance"`
	TotalIncome     float64              `json:"total_income"`
	TotalExpense    float64              `json:"total_expense"`
	Accounts        []AccountBalance     `json:"accounts"`
	RecentTransactions []Transaction    `json:"recent_transactions"`
}

// AccountBalance é o saldo de uma conta específica.
type AccountBalance struct {
	Account  Account `json:"account"`
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
}

// FinanceFilter é o filtro para busca de transações.
type FinanceFilter struct {
	AccountID  *int    `json:"account_id"`
	Type       string  `json:"type"`
	Category   string  `json:"category"`
	DateFrom   string  `json:"date_from"`
	DateTo     string  `json:"date_to"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}
