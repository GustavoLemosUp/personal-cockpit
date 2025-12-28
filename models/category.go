package models

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
