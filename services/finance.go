package services

import (
	"database/sql"
	"fmt"
	"time"

	"personal-cockpit/models"
)

type FinanceService struct {
	db *sql.DB
}

func NewFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{db: db}
}

// ─── Accounts ──────────────────────────────────────────────

func (s *FinanceService) CreateAccount(a models.Account) (int64, error) {
	if a.Currency == "" {
		a.Currency = "BRL"
	}
	if a.Color == "" {
		a.Color = "#6b7280"
	}
	res, err := s.db.Exec(`
		INSERT INTO finance_accounts (name, type, balance, currency, color, icon, description, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Name, a.Type, a.Balance, a.Currency, a.Color, a.Icon, a.Description, boolToInt(a.IsActive))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *FinanceService) GetAccounts() ([]models.Account, error) {
	rows, err := s.db.Query(`
		SELECT id, name, type, balance, currency, color, icon, description, is_active, created_at, updated_at
		FROM finance_accounts WHERE is_active = 1 ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Account
	for rows.Next() {
		var a models.Account
		var isActive int
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.Currency,
			&a.Color, &a.Icon, &a.Description, &isActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
			continue
		}
		a.IsActive = isActive == 1
		list = append(list, a)
	}
	return list, nil
}

func (s *FinanceService) UpdateAccount(a models.Account) error {
	_, err := s.db.Exec(`
		UPDATE finance_accounts SET name=?, type=?, color=?, icon=?, description=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`, a.Name, a.Type, a.Color, a.Icon, a.Description, a.ID)
	return err
}

func (s *FinanceService) DeleteAccount(id int) error {
	_, err := s.db.Exec(`UPDATE finance_accounts SET is_active=0 WHERE id=?`, id)
	return err
}

// ─── Transactions ──────────────────────────────────────────

func (s *FinanceService) CreateTransaction(t models.Transaction) (int64, error) {
	if t.Date.IsZero() {
		t.Date = time.Now()
	}
	res, err := s.db.Exec(`
		INSERT INTO finance_transactions (account_id, to_account_id, type, category, description, amount, date, is_recurring, recurrence_rule, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.AccountID, nullableInt(t.ToAccountID), t.Type, t.Category, t.Description,
		t.Amount, t.Date, boolToInt(t.IsRecurring), t.RecurrenceRule, t.Tags)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Atualiza saldo da conta
	if err := s.applyBalanceChange(t); err != nil {
		return id, fmt.Errorf("transação criada mas saldo não atualizado: %w", err)
	}
	return id, nil
}

func (s *FinanceService) applyBalanceChange(t models.Transaction) error {
	switch t.Type {
	case models.TransactionTypeIncome:
		_, err := s.db.Exec(`UPDATE finance_accounts SET balance = balance + ?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, t.Amount, t.AccountID)
		return err
	case models.TransactionTypeExpense:
		_, err := s.db.Exec(`UPDATE finance_accounts SET balance = balance - ?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, t.Amount, t.AccountID)
		return err
	case models.TransactionTypeTransfer:
		_, err := s.db.Exec(`UPDATE finance_accounts SET balance = balance - ?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, t.Amount, t.AccountID)
		if err != nil {
			return err
		}
		if t.ToAccountID != nil {
			_, err = s.db.Exec(`UPDATE finance_accounts SET balance = balance + ?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, t.Amount, *t.ToAccountID)
		}
		return err
	}
	return nil
}

func (s *FinanceService) GetTransactions(filter models.FinanceFilter) ([]models.Transaction, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	query := `
		SELECT t.id, t.account_id, a.name, t.to_account_id,
			   t.type, t.category, t.description, t.amount, t.date,
			   t.is_recurring, t.recurrence_rule, t.tags, t.created_at
		FROM finance_transactions t
		JOIN finance_accounts a ON t.account_id = a.id
		WHERE 1=1`
	args := []interface{}{}

	if filter.AccountID != nil {
		query += ` AND t.account_id = ?`
		args = append(args, *filter.AccountID)
	}
	if filter.Type != "" {
		query += ` AND t.type = ?`
		args = append(args, filter.Type)
	}
	if filter.Category != "" {
		query += ` AND t.category = ?`
		args = append(args, filter.Category)
	}
	if filter.DateFrom != "" {
		query += ` AND t.date >= ?`
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		query += ` AND t.date <= ?`
		args = append(args, filter.DateTo)
	}
	query += ` ORDER BY t.date DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var isRec int
		var toAccID sql.NullInt64
		if err := rows.Scan(&t.ID, &t.AccountID, &t.AccountName, &toAccID,
			&t.Type, &t.Category, &t.Description, &t.Amount, &t.Date,
			&isRec, &t.RecurrenceRule, &t.Tags, &t.CreatedAt); err != nil {
			continue
		}
		t.IsRecurring = isRec == 1
		if toAccID.Valid {
			v := int(toAccID.Int64)
			t.ToAccountID = &v
		}
		list = append(list, t)
	}
	return list, nil
}

func (s *FinanceService) DeleteTransaction(id int) error {
	// Reverte o efeito no saldo antes de deletar
	var t models.Transaction
	var toAccID sql.NullInt64
	err := s.db.QueryRow(`SELECT account_id, to_account_id, type, amount FROM finance_transactions WHERE id=?`, id).
		Scan(&t.AccountID, &toAccID, &t.Type, &t.Amount)
	if err != nil {
		return err
	}
	if toAccID.Valid {
		v := int(toAccID.Int64)
		t.ToAccountID = &v
	}
	// Reverte: aplica operação inversa
	switch t.Type {
	case models.TransactionTypeIncome:
		s.db.Exec(`UPDATE finance_accounts SET balance = balance - ? WHERE id=?`, t.Amount, t.AccountID)
	case models.TransactionTypeExpense:
		s.db.Exec(`UPDATE finance_accounts SET balance = balance + ? WHERE id=?`, t.Amount, t.AccountID)
	case models.TransactionTypeTransfer:
		s.db.Exec(`UPDATE finance_accounts SET balance = balance + ? WHERE id=?`, t.Amount, t.AccountID)
		if t.ToAccountID != nil {
			s.db.Exec(`UPDATE finance_accounts SET balance = balance - ? WHERE id=?`, t.Amount, *t.ToAccountID)
		}
	}
	_, err = s.db.Exec(`DELETE FROM finance_transactions WHERE id=?`, id)
	return err
}

// ─── Scheduled Transactions ────────────────────────────────

func (s *FinanceService) CreateScheduled(st models.ScheduledTransaction) (int64, error) {
	res, err := s.db.Exec(`
		INSERT INTO finance_scheduled (account_id, type, category, description, amount, scheduled_at, is_recurring, recurrence_rule)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		st.AccountID, st.Type, st.Category, st.Description, st.Amount,
		st.ScheduledAt, boolToInt(st.IsRecurring), st.RecurrenceRule)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *FinanceService) GetScheduled() ([]models.ScheduledTransaction, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.account_id, a.name, s.type, s.category, s.description,
			   s.amount, s.scheduled_at, s.is_recurring, s.recurrence_rule, s.is_executed, s.created_at
		FROM finance_scheduled s
		JOIN finance_accounts a ON s.account_id = a.id
		WHERE s.is_executed = 0
		ORDER BY s.scheduled_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ScheduledTransaction
	for rows.Next() {
		var st models.ScheduledTransaction
		var isRec, isExec int
		if err := rows.Scan(&st.ID, &st.AccountID, &st.AccountName, &st.Type,
			&st.Category, &st.Description, &st.Amount, &st.ScheduledAt,
			&isRec, &st.RecurrenceRule, &isExec, &st.CreatedAt); err != nil {
			continue
		}
		st.IsRecurring = isRec == 1
		st.IsExecuted = isExec == 1
		list = append(list, st)
	}
	return list, nil
}

func (s *FinanceService) ExecuteScheduled(id int) error {
	var st models.ScheduledTransaction
	err := s.db.QueryRow(`SELECT account_id, type, category, description, amount FROM finance_scheduled WHERE id=?`, id).
		Scan(&st.AccountID, &st.Type, &st.Category, &st.Description, &st.Amount)
	if err != nil {
		return err
	}
	now := time.Now()
	if _, err := s.CreateTransaction(models.Transaction{
		AccountID:   st.AccountID,
		Type:        st.Type,
		Category:    st.Category,
		Description: st.Description,
		Amount:      st.Amount,
		Date:        now,
	}); err != nil {
		return err
	}
	_, err = s.db.Exec(`UPDATE finance_scheduled SET is_executed=1 WHERE id=?`, id)
	return err
}

func (s *FinanceService) DeleteScheduled(id int) error {
	_, err := s.db.Exec(`DELETE FROM finance_scheduled WHERE id=?`, id)
	return err
}

// ─── Summary ───────────────────────────────────────────────

func (s *FinanceService) GetSummary() (*models.FinanceSummary, error) {
	accounts, err := s.GetAccounts()
	if err != nil {
		return nil, err
	}

	var total float64
	var balances []models.AccountBalance
	for _, acc := range accounts {
		total += acc.Balance
		balances = append(balances, models.AccountBalance{Account: acc})
	}

	// Receitas e despesas do mês atual
	now := time.Now()
	monthStart := fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())
	monthEnd := fmt.Sprintf("%d-%02d-31", now.Year(), now.Month())

	var income, expense float64
	s.db.QueryRow(`SELECT COALESCE(SUM(amount),0) FROM finance_transactions WHERE type='income' AND date >= ? AND date <= ?`, monthStart, monthEnd).Scan(&income)
	s.db.QueryRow(`SELECT COALESCE(SUM(amount),0) FROM finance_transactions WHERE type='expense' AND date >= ? AND date <= ?`, monthStart, monthEnd).Scan(&expense)

	recent, _ := s.GetTransactions(models.FinanceFilter{Limit: 10})

	return &models.FinanceSummary{
		TotalBalance:       total,
		TotalIncome:        income,
		TotalExpense:       expense,
		Accounts:           balances,
		RecentTransactions: recent,
	}, nil
}

// ─── Helpers ──────────────────────────────────────────────

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullableInt(v *int) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
