package models

import "time"

// WAChat representa uma conversa (contato ou grupo).
type WAChat struct {
	JID           string    `json:"jid"`
	Name          string    `json:"name"`
	LastMessage   string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `json:"unread_count"`
}

// WAMessage representa uma mensagem individual.
type WAMessage struct {
	ID        string    `json:"id"`
	ChatJID   string    `json:"chat_jid"`
	SenderJID string    `json:"sender_jid"`
	Content   string    `json:"content"`
	IsFromMe  bool      `json:"is_from_me"`
	Timestamp time.Time `json:"timestamp"`
}

// WAScheduled representa uma mensagem agendada.
type WAScheduled struct {
	ID          int       `json:"id"`
	ChatJID     string    `json:"chat_jid"`
	ChatName    string    `json:"chat_name"`
	Content     string    `json:"content"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Sent        bool      `json:"sent"`
	CreatedAt   time.Time `json:"created_at"`
}

// WAScheduleInput é o input para agendar uma mensagem.
type WAScheduleInput struct {
	ChatJID     string    `json:"chat_jid"`
	ChatName    string    `json:"chat_name"`
	Content     string    `json:"content"`
	ScheduledAt time.Time `json:"scheduled_at"`
}
