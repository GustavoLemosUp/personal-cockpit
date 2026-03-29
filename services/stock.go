package services

import (
	"database/sql"
	"time"

	"personal-cockpit/models"
)

type StockService struct {
	db *sql.DB
}

func NewStockService(db *sql.DB) *StockService {
	return &StockService{db: db}
}

// ─── Products ─────────────────────────────────────────────

func (s *StockService) CreateProduct(p models.Product) (int64, error) {
	p.IsActive = true
	res, err := s.db.Exec(`
		INSERT INTO stock_products (name, category, unit, min_quantity, current_stock, price, notes, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.Category, p.Unit, p.MinQuantity, p.CurrentStock, p.Price, p.Notes, boolToInt(p.IsActive))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *StockService) GetProducts() ([]models.Product, error) {
	rows, err := s.db.Query(`
		SELECT id, name, category, unit, min_quantity, current_stock, price, notes, is_active, created_at, updated_at
		FROM stock_products WHERE is_active=1 ORDER BY category ASC, name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (s *StockService) UpdateProduct(p models.Product) error {
	_, err := s.db.Exec(`
		UPDATE stock_products SET name=?, category=?, unit=?, min_quantity=?, price=?, notes=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`, p.Name, p.Category, p.Unit, p.MinQuantity, p.Price, p.Notes, p.ID)
	return err
}

func (s *StockService) DeleteProduct(id int) error {
	_, err := s.db.Exec(`UPDATE stock_products SET is_active=0 WHERE id=?`, id)
	return err
}

func (s *StockService) GetLowStock() ([]models.StockAlert, error) {
	rows, err := s.db.Query(`
		SELECT id, name, category, unit, min_quantity, current_stock, price, notes, is_active, created_at, updated_at
		FROM stock_products
		WHERE is_active=1 AND current_stock < min_quantity
		ORDER BY (min_quantity - current_stock) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products, err := scanProducts(rows)
	if err != nil {
		return nil, err
	}
	var alerts []models.StockAlert
	for _, p := range products {
		alerts = append(alerts, models.StockAlert{
			Product: p,
			Deficit: p.MinQuantity - p.CurrentStock,
		})
	}
	return alerts, nil
}

// ─── Stock Movements ──────────────────────────────────────

func (s *StockService) AddMovement(m models.StockMovement) (int64, error) {
	if m.Date.IsZero() {
		m.Date = time.Now()
	}
	res, err := s.db.Exec(`
		INSERT INTO stock_movements (product_id, type, quantity, price, notes, date)
		VALUES (?, ?, ?, ?, ?, ?)`,
		m.ProductID, m.Type, m.Quantity, m.Price, m.Notes, m.Date)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Atualiza estoque atual
	if m.Type == models.StockMovementIn {
		s.db.Exec(`UPDATE stock_products SET current_stock = current_stock + ?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, m.Quantity, m.ProductID)
	} else {
		s.db.Exec(`UPDATE stock_products SET current_stock = MAX(0, current_stock - ?), updated_at=CURRENT_TIMESTAMP WHERE id=?`, m.Quantity, m.ProductID)
	}
	return id, nil
}

func (s *StockService) GetMovements(productID int, limit int) ([]models.StockMovement, error) {
	if limit <= 0 {
		limit = 30
	}
	rows, err := s.db.Query(`
		SELECT m.id, m.product_id, p.name, m.type, m.quantity, m.price, m.notes, m.date, m.created_at
		FROM stock_movements m JOIN stock_products p ON m.product_id = p.id
		WHERE m.product_id = ?
		ORDER BY m.date DESC LIMIT ?`, productID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.StockMovement
	for rows.Next() {
		var m models.StockMovement
		if err := rows.Scan(&m.ID, &m.ProductID, &m.ProductName, &m.Type, &m.Quantity, &m.Price, &m.Notes, &m.Date, &m.CreatedAt); err != nil {
			continue
		}
		list = append(list, m)
	}
	return list, nil
}

// ─── Shopping Lists ───────────────────────────────────────

func (s *StockService) CreateShoppingList(sl models.ShoppingList) (int64, error) {
	if sl.Month == "" {
		now := time.Now()
		sl.Month = now.Format("2006-01")
	}
	res, err := s.db.Exec(`
		INSERT INTO shopping_lists (name, month, description, total_budget)
		VALUES (?, ?, ?, ?)`,
		sl.Name, sl.Month, sl.Description, sl.TotalBudget)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *StockService) GetShoppingLists() ([]models.ShoppingList, error) {
	rows, err := s.db.Query(`
		SELECT id, name, month, description, is_completed, total_budget, total_spent, created_at, updated_at
		FROM shopping_lists ORDER BY month DESC, created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ShoppingList
	for rows.Next() {
		var sl models.ShoppingList
		var isDone int
		if err := rows.Scan(&sl.ID, &sl.Name, &sl.Month, &sl.Description, &isDone,
			&sl.TotalBudget, &sl.TotalSpent, &sl.CreatedAt, &sl.UpdatedAt); err != nil {
			continue
		}
		sl.IsCompleted = isDone == 1
		list = append(list, sl)
	}
	return list, nil
}

func (s *StockService) GetShoppingListWithItems(id int) (*models.ShoppingList, error) {
	var sl models.ShoppingList
	var isDone int
	err := s.db.QueryRow(`
		SELECT id, name, month, description, is_completed, total_budget, total_spent, created_at, updated_at
		FROM shopping_lists WHERE id=?`, id).
		Scan(&sl.ID, &sl.Name, &sl.Month, &sl.Description, &isDone,
			&sl.TotalBudget, &sl.TotalSpent, &sl.CreatedAt, &sl.UpdatedAt)
	if err != nil {
		return nil, err
	}
	sl.IsCompleted = isDone == 1

	items, err := s.GetShoppingItems(id)
	if err != nil {
		return nil, err
	}
	sl.Items = items
	return &sl, nil
}

func (s *StockService) DeleteShoppingList(id int) error {
	_, err := s.db.Exec(`DELETE FROM shopping_lists WHERE id=?`, id)
	return err
}

// ─── Shopping Items ───────────────────────────────────────

func (s *StockService) AddShoppingItem(item models.ShoppingItem) (int64, error) {
	res, err := s.db.Exec(`
		INSERT INTO shopping_items (shopping_list_id, product_id, name, quantity, unit, estimated_price, category, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ShoppingListID, nullableInt(item.ProductID), item.Name, item.Quantity,
		item.Unit, item.EstimatedPrice, item.Category, item.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *StockService) GetShoppingItems(listID int) ([]models.ShoppingItem, error) {
	rows, err := s.db.Query(`
		SELECT id, shopping_list_id, product_id, name, quantity, unit,
			   estimated_price, actual_price, category, is_bought, notes, created_at
		FROM shopping_items WHERE shopping_list_id=?
		ORDER BY category ASC, name ASC`, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ShoppingItem
	for rows.Next() {
		var item models.ShoppingItem
		var isBought int
		var productID sql.NullInt64
		if err := rows.Scan(&item.ID, &item.ShoppingListID, &productID, &item.Name,
			&item.Quantity, &item.Unit, &item.EstimatedPrice, &item.ActualPrice,
			&item.Category, &isBought, &item.Notes, &item.CreatedAt); err != nil {
			continue
		}
		item.IsBought = isBought == 1
		if productID.Valid {
			v := int(productID.Int64)
			item.ProductID = &v
		}
		list = append(list, item)
	}
	return list, nil
}

func (s *StockService) ToggleShoppingItem(id int) error {
	_, err := s.db.Exec(`UPDATE shopping_items SET is_bought = NOT is_bought WHERE id=?`, id)
	if err != nil {
		return err
	}
	// Recalcula total gasto na lista
	var listID int
	s.db.QueryRow(`SELECT shopping_list_id FROM shopping_items WHERE id=?`, id).Scan(&listID)
	s.db.Exec(`
		UPDATE shopping_lists SET total_spent = (
			SELECT COALESCE(SUM(CASE WHEN is_bought=1 THEN actual_price ELSE 0 END),0)
			FROM shopping_items WHERE shopping_list_id=?
		) WHERE id=?`, listID, listID)
	return nil
}

func (s *StockService) UpdateShoppingItem(item models.ShoppingItem) error {
	_, err := s.db.Exec(`
		UPDATE shopping_items SET name=?, quantity=?, unit=?, estimated_price=?, actual_price=?, category=?, notes=?
		WHERE id=?`,
		item.Name, item.Quantity, item.Unit, item.EstimatedPrice, item.ActualPrice,
		item.Category, item.Notes, item.ID)
	return err
}

func (s *StockService) DeleteShoppingItem(id int) error {
	_, err := s.db.Exec(`DELETE FROM shopping_items WHERE id=?`, id)
	return err
}

// ─── Helpers ──────────────────────────────────────────────

func scanProducts(rows *sql.Rows) ([]models.Product, error) {
	var list []models.Product
	for rows.Next() {
		var p models.Product
		var isActive int
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Unit, &p.MinQuantity,
			&p.CurrentStock, &p.Price, &p.Notes, &isActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			continue
		}
		p.IsActive = isActive == 1
		list = append(list, p)
	}
	return list, nil
}
