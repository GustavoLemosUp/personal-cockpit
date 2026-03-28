package models

import "time"

type Profile struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	GoogleID          string     `json:"google_id"`
	Role              string     `json:"role"` // "admin" | "member"
	AvatarURL         string     `json:"avatar_url"`
	GoogleAccessToken string     `json:"google_access_token,omitempty"`
	GoogleCalendarID  string     `json:"google_calendar_id"`
	PinHash           string     `json:"-"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
}

// ProfileSummary is the safe public representation (no tokens or pin).
type ProfileSummary struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Role             string    `json:"role"`
	AvatarURL        string    `json:"avatar_url"`
	GoogleCalendarID string    `json:"google_calendar_id"`
	HasGoogleToken   bool      `json:"has_google_token"`
	HasPin           bool      `json:"has_pin"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateProfileInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
	Pin       string `json:"pin,omitempty"`
}
