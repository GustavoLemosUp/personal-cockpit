package models

import "time"

type Task struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	CategoryID     *int       `json:"category_id"`
	BoardID        *int       `json:"board_id"`
	ColumnID       *int       `json:"column_id"`
	Position       int        `json:"position"`
	DueDate        *time.Time `json:"due_date"`
	StartDate      *time.Time `json:"start_date"`
	EstimatedHours *float64   `json:"estimated_hours"`
	Progress       int        `json:"progress"`
	CompletedAt    *time.Time `json:"completed_at"`
	GoogleEventID  *string    `json:"google_event_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type TaskFilter struct {
	Status     string
	Priority   string
	CategoryID *int
	BoardID    *int
	ColumnID   *int
}
