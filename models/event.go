package models

import "time"

type Event struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	AllDay          bool      `json:"all_day"`
	Color           string    `json:"color"`
	Location        string    `json:"location"`
	ReminderMinutes *int      `json:"reminder_minutes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
