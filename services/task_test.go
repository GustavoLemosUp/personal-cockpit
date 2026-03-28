package services

import (
	"testing"

	"personal-cockpit/database"
	"personal-cockpit/models"
)

func setupTaskService(t *testing.T) *TaskService {
	t.Helper()
	db, err := database.NewInMemoryDB()
	if err != nil {
		t.Fatalf("falha ao criar banco de teste: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewTaskService(db)
}

func TestCreateTask(t *testing.T) {
	svc := setupTaskService(t)

	id, err := svc.CreateTask(models.Task{
		Title:       "Comprar café",
		Description: "Café moído, 500g",
		Status:      "pending",
		Priority:    "high",
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if id <= 0 {
		t.Errorf("esperado ID > 0, got %d", id)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	svc := setupTaskService(t)

	tests := []struct {
		name string
		task models.Task
	}{
		{"sem título", models.Task{Description: "desc", Priority: "high", Status: "pending"}},
		{"sem descrição", models.Task{Title: "titulo", Priority: "high", Status: "pending"}},
		{"sem prioridade", models.Task{Title: "titulo", Description: "desc", Status: "pending"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateTask(tt.task)
			if err == nil {
				t.Error("esperava erro de validação, got nil")
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	svc := setupTaskService(t)

	for _, title := range []string{"Tarefa A", "Tarefa B", "Tarefa C"} {
		svc.CreateTask(models.Task{Title: title, Description: "desc", Status: "pending", Priority: "medium"})
	}

	tasks, err := svc.GetAllTasks()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("esperado 3 tarefas, got %d", len(tasks))
	}
}

func TestGetTaskByID(t *testing.T) {
	svc := setupTaskService(t)

	id, _ := svc.CreateTask(models.Task{
		Title:       "Tarefa específica",
		Description: "desc",
		Status:      "pending",
		Priority:    "low",
	})

	task, err := svc.GetTaskByID(int(id))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if task.Title != "Tarefa específica" {
		t.Errorf("esperado 'Tarefa específica', got '%s'", task.Title)
	}
}

func TestGetTaskByIDNotFound(t *testing.T) {
	svc := setupTaskService(t)

	_, err := svc.GetTaskByID(9999)
	if err == nil {
		t.Error("esperava erro para ID inexistente")
	}
}

func TestUpdateTask(t *testing.T) {
	svc := setupTaskService(t)

	id, _ := svc.CreateTask(models.Task{
		Title:       "Original",
		Description: "desc",
		Status:      "pending",
		Priority:    "low",
	})

	err := svc.UpdateTask(models.Task{
		ID:          int(id),
		Title:       "Atualizado",
		Description: "nova desc",
		Status:      "pending",
		Priority:    "high",
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	task, _ := svc.GetTaskByID(int(id))
	if task.Title != "Atualizado" {
		t.Errorf("esperado 'Atualizado', got '%s'", task.Title)
	}
	if task.Priority != "high" {
		t.Errorf("esperado 'high', got '%s'", task.Priority)
	}
}

func TestDeleteTask(t *testing.T) {
	svc := setupTaskService(t)

	id, _ := svc.CreateTask(models.Task{
		Title:       "Para deletar",
		Description: "desc",
		Status:      "pending",
		Priority:    "low",
	})

	if err := svc.DeleteTask(int(id)); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	_, err := svc.GetTaskByID(int(id))
	if err == nil {
		t.Error("tarefa deveria ter sido deletada")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	svc := setupTaskService(t)

	err := svc.DeleteTask(9999)
	if err == nil {
		t.Error("esperava erro para ID inexistente")
	}
}

func TestToggleTaskStatus(t *testing.T) {
	svc := setupTaskService(t)

	id, _ := svc.CreateTask(models.Task{
		Title:       "Toggle test",
		Description: "desc",
		Status:      "pending",
		Priority:    "medium",
	})

	if err := svc.ToggleTaskStatus(int(id)); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	task, _ := svc.GetTaskByID(int(id))
	if task.Status != "completed" {
		t.Errorf("esperado 'completed', got '%s'", task.Status)
	}

	svc.ToggleTaskStatus(int(id))
	task, _ = svc.GetTaskByID(int(id))
	if task.Status != "pending" {
		t.Errorf("esperado 'pending' após segundo toggle, got '%s'", task.Status)
	}
}

func TestGetPendingAndCompletedTasks(t *testing.T) {
	svc := setupTaskService(t)

	id1, _ := svc.CreateTask(models.Task{Title: "P1", Description: "d", Status: "pending", Priority: "low"})
	id2, _ := svc.CreateTask(models.Task{Title: "P2", Description: "d", Status: "pending", Priority: "low"})
	svc.CreateTask(models.Task{Title: "C1", Description: "d", Status: "pending", Priority: "low"})

	svc.ToggleTaskStatus(int(id1))
	svc.ToggleTaskStatus(int(id2))

	pending, _ := svc.GetPendingTasks()
	completed, _ := svc.GetCompletedTasks()

	if len(pending) != 1 {
		t.Errorf("esperado 1 pendente, got %d", len(pending))
	}
	if len(completed) != 2 {
		t.Errorf("esperado 2 concluídas, got %d", len(completed))
	}
}

func TestGetTasksByFilter(t *testing.T) {
	svc := setupTaskService(t)

	svc.CreateTask(models.Task{Title: "Alta 1", Description: "d", Status: "pending", Priority: "high"})
	svc.CreateTask(models.Task{Title: "Alta 2", Description: "d", Status: "pending", Priority: "high"})
	svc.CreateTask(models.Task{Title: "Baixa 1", Description: "d", Status: "pending", Priority: "low"})

	tasks, err := svc.GetTasksByFilter(models.TaskFilter{Priority: "high"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("esperado 2 tarefas de alta prioridade, got %d", len(tasks))
	}
}
