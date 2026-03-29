package services

import (
	"testing"
	"time"

	"personal-cockpit/database"
	"personal-cockpit/models"
)

func setupEventService(t *testing.T) *EventService {
	t.Helper()
	db, err := database.NewInMemoryDB()
	if err != nil {
		t.Fatalf("falha ao criar banco de teste: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewEventService(db)
}

func makeEvent(title string, start, end time.Time) models.Event {
	return models.Event{
		Title:     title,
		StartDate: start,
		EndDate:   end,
		Color:     "#3b82f6",
	}
}

func TestCreateEvent(t *testing.T) {
	svc := setupEventService(t)

	now := time.Now()
	id, err := svc.CreateEvent(makeEvent("Reunião", now, now.Add(time.Hour)))

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if id <= 0 {
		t.Errorf("esperado ID > 0, got %d", id)
	}
}

func TestCreateEventValidation(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	tests := []struct {
		name  string
		event models.Event
	}{
		{"sem título", models.Event{StartDate: now, EndDate: now.Add(time.Hour)}},
		{"sem data início", models.Event{Title: "Evento", EndDate: now.Add(time.Hour)}},
		{"sem data fim", models.Event{Title: "Evento", StartDate: now}},
		{"fim antes do início", makeEvent("Inválido", now, now.Add(-time.Hour))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateEvent(tt.event)
			if err == nil {
				t.Error("esperava erro de validação, got nil")
			}
		})
	}
}

func TestGetAllEvents(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	for _, title := range []string{"Evento A", "Evento B"} {
		svc.CreateEvent(makeEvent(title, now, now.Add(time.Hour)))
	}

	events, err := svc.GetAllEvents()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("esperado 2 eventos, got %d", len(events))
	}
}

func TestGetEventByID(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	id, _ := svc.CreateEvent(makeEvent("Evento específico", now, now.Add(time.Hour)))

	event, err := svc.GetEventByID(int(id))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if event.Title != "Evento específico" {
		t.Errorf("esperado 'Evento específico', got '%s'", event.Title)
	}
}

func TestGetEventByIDNotFound(t *testing.T) {
	svc := setupEventService(t)

	_, err := svc.GetEventByID(9999)
	if err == nil {
		t.Error("esperava erro para ID inexistente")
	}
}

func TestUpdateEvent(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	id, _ := svc.CreateEvent(makeEvent("Original", now, now.Add(time.Hour)))

	err := svc.UpdateEvent(models.Event{
		ID:        int(id),
		Title:     "Atualizado",
		StartDate: now,
		EndDate:   now.Add(2 * time.Hour),
		Color:     "#ff0000",
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	event, _ := svc.GetEventByID(int(id))
	if event.Title != "Atualizado" {
		t.Errorf("esperado 'Atualizado', got '%s'", event.Title)
	}
}

func TestDeleteEvent(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	id, _ := svc.CreateEvent(makeEvent("Para deletar", now, now.Add(time.Hour)))

	if err := svc.DeleteEvent(int(id)); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	_, err := svc.GetEventByID(int(id))
	if err == nil {
		t.Error("evento deveria ter sido deletado")
	}
}

func TestGetTodayEvents(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	// Evento de hoje
	svc.CreateEvent(makeEvent("Hoje", now, now.Add(time.Hour)))

	// Evento de amanhã
	amanha := now.Add(25 * time.Hour)
	svc.CreateEvent(makeEvent("Amanhã", amanha, amanha.Add(time.Hour)))

	events, err := svc.GetTodayEvents()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("esperado 1 evento hoje, got %d", len(events))
	}
	if events[0].Title != "Hoje" {
		t.Errorf("esperado 'Hoje', got '%s'", events[0].Title)
	}
}

func TestGetUpcomingEvents(t *testing.T) {
	svc := setupEventService(t)
	now := time.Now()

	// Dentro de 7 dias
	svc.CreateEvent(makeEvent("Em 3 dias", now.Add(72*time.Hour), now.Add(73*time.Hour)))
	svc.CreateEvent(makeEvent("Em 6 dias", now.Add(144*time.Hour), now.Add(145*time.Hour)))

	// Além de 7 dias
	longe := now.Add(10 * 24 * time.Hour)
	svc.CreateEvent(makeEvent("Em 10 dias", longe, longe.Add(time.Hour)))

	events, err := svc.GetUpcomingEvents()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("esperado 2 eventos nos próximos 7 dias, got %d", len(events))
	}
}
