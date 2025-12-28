package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"personal-cockpit/models"
)

type EventService struct {
	db *sql.DB
}

func NewEventService(db *sql.DB) *EventService {
	return &EventService{db: db}
}

func (s *EventService) CreateEvent(event models.Event) (int64, error) {
	var erros []string

	if event.Title == "" {
		erros = append(erros, "título")
	}

	// IsZero() verifica se time.Time está vazio (valor zero)
	if event.StartDate.IsZero() {
		erros = append(erros, "data de início")
	}

	if event.EndDate.IsZero() {
		erros = append(erros, "data de término")
	}

	if !event.StartDate.IsZero() && !event.EndDate.IsZero() && event.EndDate.Before(event.StartDate) {
		erros = append(erros, "data de término deve ser após data de início")
	}

	if len(erros) > 0 {
		mensagem := "Campos obrigatórios:\n- " + strings.Join(erros, "\n- ")
		return 0, fmt.Errorf(mensagem)
	}

	query := `
		INSERT INTO events (title, description, start_date, end_date, all_day, color, location, reminder_minutes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.AllDay,
		event.Color,
		event.Location,
		event.ReminderMinutes,
	)

	if err != nil {
		return 0, fmt.Errorf("erro ao criar evento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter ID: %w", err)
	}

	return id, nil
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_date, end_date, all_day, color, location, 
		       reminder_minutes, created_at, updated_at
		FROM events
		ORDER BY start_date ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos: %w", err)
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.AllDay,
			&event.Color,
			&event.Location,
			&event.ReminderMinutes,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler evento: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *EventService) GetEventByID(id int) (*models.Event, error) {
	query := `
		SELECT id, title, description, start_date, end_date, all_day, color, location, 
		       reminder_minutes, created_at, updated_at
		FROM events
		WHERE id = ?
	`

	var event models.Event
	err := s.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.AllDay,
		&event.Color,
		&event.Location,
		&event.ReminderMinutes,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("evento não encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar evento: %w", err)
	}

	return &event, nil
}

func (s *EventService) UpdateEvent(event models.Event) error {
	if event.ID == 0 {
		return fmt.Errorf("ID do evento é obrigatório")
	}

	query := `
		UPDATE events 
		SET title = ?, description = ?, start_date = ?, end_date = ?, all_day = ?, 
		    color = ?, location = ?, reminder_minutes = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.AllDay,
		event.Color,
		event.Location,
		event.ReminderMinutes,
		event.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar evento: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("evento não encontrado")
	}

	return nil
}

func (s *EventService) DeleteEvent(id int) error {
	query := "DELETE FROM events WHERE id = ?"

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar evento: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("evento não encontrado")
	}

	return nil
}

// GetEventsByDateRange busca eventos entre duas datas
func (s *EventService) GetEventsByDateRange(startDate, endDate time.Time) ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_date, end_date, all_day, color, location, 
		       reminder_minutes, created_at, updated_at
		FROM events
		WHERE start_date >= ? AND start_date <= ?
		ORDER BY start_date ASC
	`

	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos: %w", err)
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.AllDay,
			&event.Color,
			&event.Location,
			&event.ReminderMinutes,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler evento: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

// GetTodayEvents retorna eventos de hoje
func (s *EventService) GetTodayEvents() ([]models.Event, error) {
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.GetEventsByDateRange(startOfDay, endOfDay)
}

func (s *EventService) GetUpcomingEvents() ([]models.Event, error) {
	now := time.Now()

	future := now.Add(7 * 24 * time.Hour)

	return s.GetEventsByDateRange(now, future)
}
