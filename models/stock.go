package models

import "time"

// Product representa um produto/item do estoque doméstico.
type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Unit         string    `json:"unit"`         // kg, un, L, cx, etc.
	MinQuantity  float64   `json:"min_quantity"` // alerta de estoque mínimo
	CurrentStock float64   `json:"current_stock"`
	Price        float64   `json:"price"`        // preço de referência
	Notes        string    `json:"notes"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// StockMovementType representa entrada ou saída do estoque.
type StockMovementType string

const (
	StockMovementIn  StockMovementType = "in"  // compra/entrada
	StockMovementOut StockMovementType = "out" // consumo/saída
)

// StockMovement registra uma movimentação de estoque.
type StockMovement struct {
	ID          int               `json:"id"`
	ProductID   int               `json:"product_id"`
	ProductName string            `json:"product_name,omitempty"`
	Type        StockMovementType `json:"type"`
	Quantity    float64           `json:"quantity"`
	Price       float64           `json:"price"` // preço na movimentação
	Notes       string            `json:"notes"`
	Date        time.Time         `json:"date"`
	CreatedAt   time.Time         `json:"created_at"`
}

// ShoppingList representa uma lista de compras.
type ShoppingList struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Month       string         `json:"month"` // YYYY-MM
	Description string         `json:"description"`
	IsCompleted bool           `json:"is_completed"`
	TotalBudget float64        `json:"total_budget"`
	TotalSpent  float64        `json:"total_spent"`
	Items       []ShoppingItem `json:"items,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ShoppingItem representa um item de uma lista de compras.
type ShoppingItem struct {
	ID             int       `json:"id"`
	ShoppingListID int       `json:"shopping_list_id"`
	ProductID      *int      `json:"product_id,omitempty"`
	Name           string    `json:"name"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	EstimatedPrice float64   `json:"estimated_price"`
	ActualPrice    float64   `json:"actual_price"`
	Category       string    `json:"category"`
	IsBought       bool      `json:"is_bought"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

// StockAlert representa um produto com estoque abaixo do mínimo.
type StockAlert struct {
	Product      Product `json:"product"`
	Deficit      float64 `json:"deficit"`
}
