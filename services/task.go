package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"personal-cockpit/models"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(task models.Task) (int64, error) {
	var erros []string
	if task.Title == "" {
		erros = append(erros, "título")
	}
	if task.Priority == "" {
		erros = append(erros, "prioridade")
	}
	if len(erros) > 0 {
		return 0, errors.New("Campos obrigatórios:\n- " + strings.Join(erros, "\n- "))
	}

	result, err := s.db.Exec(`
		INSERT INTO tasks
		    (title, description, status, priority, category_id, due_date,
		     board_id, column_id, position, start_date, estimated_hours, progress, google_event_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Title, task.Description, task.Status, task.Priority, task.CategoryID, task.DueDate,
		task.BoardID, task.ColumnID, task.Position, task.StartDate, task.EstimatedHours,
		task.Progress, task.GoogleEventID,
	)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar tarefa: %w", err)
	}
	return result.LastInsertId()
}

const taskSelectColumns = `
	id, title, description, status, priority, category_id,
	board_id, column_id, position, due_date, start_date,
	estimated_hours, progress, completed_at, google_event_id,
	created_at, updated_at
`

func scanTask(row interface{ Scan(...any) error }) (models.Task, error) {
	var t models.Task
	err := row.Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.CategoryID,
		&t.BoardID, &t.ColumnID, &t.Position, &t.DueDate, &t.StartDate,
		&t.EstimatedHours, &t.Progress, &t.CompletedAt, &t.GoogleEventID,
		&t.CreatedAt, &t.UpdatedAt,
	)
	return t, err
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	rows, err := s.db.Query(`SELECT ` + taskSelectColumns + ` FROM tasks ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	row := s.db.QueryRow(`SELECT `+taskSelectColumns+` FROM tasks WHERE id = ?`, id)
	t, err := scanTask(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tarefa não encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	return &t, nil
}

func (s *TaskService) UpdateTask(task models.Task) error {
	if task.ID == 0 {
		return fmt.Errorf("ID da tarefa é obrigatório")
	}
	result, err := s.db.Exec(`
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?,
		    category_id = ?, due_date = ?, board_id = ?, column_id = ?,
		    position = ?, start_date = ?, estimated_hours = ?, progress = ?,
		    google_event_id = ?
		WHERE id = ?
	`,
		task.Title, task.Description, task.Status, task.Priority,
		task.CategoryID, task.DueDate, task.BoardID, task.ColumnID,
		task.Position, task.StartDate, task.EstimatedHours, task.Progress,
		task.GoogleEventID, task.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar tarefa: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tarefa não encontrada")
	}
	return nil
}

func (s *TaskService) DeleteTask(id int) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar tarefa: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tarefa não encontrada")
	}
	return nil
}

func (s *TaskService) ToggleTaskStatus(id int) error {
	task, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}
	newStatus := "pending"
	if task.Status == "pending" {
		newStatus = "completed"
	}
	_, err = s.db.Exec("UPDATE tasks SET status = ? WHERE id = ?", newStatus, id)
	return err
}

func (s *TaskService) GetTasksByFilter(filter models.TaskFilter) ([]models.Task, error) {
	query := `SELECT ` + taskSelectColumns + ` FROM tasks WHERE 1=1`
	args := []interface{}{}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.Priority != "" {
		query += " AND priority = ?"
		args = append(args, filter.Priority)
	}
	if filter.CategoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *filter.CategoryID)
	}
	if filter.BoardID != nil {
		query += " AND board_id = ?"
		args = append(args, *filter.BoardID)
	}
	if filter.ColumnID != nil {
		query += " AND column_id = ?"
		args = append(args, *filter.ColumnID)
	}

	query += " ORDER BY position ASC, created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func (s *TaskService) GetPendingTasks() ([]models.Task, error) {
	return s.GetTasksByFilter(models.TaskFilter{Status: "pending"})
}

func (s *TaskService) GetCompletedTasks() ([]models.Task, error) {
	return s.GetTasksByFilter(models.TaskFilter{Status: "completed"})
}

// GetTasksByColumn returns tasks for a specific Kanban column, ordered by position.
func (s *TaskService) GetTasksByColumn(columnID int) ([]models.Task, error) {
	return s.GetTasksByFilter(models.TaskFilter{ColumnID: &columnID})
}

// GetTasksByBoard returns all tasks for a board.
func (s *TaskService) GetTasksByBoard(boardID int) ([]models.Task, error) {
	return s.GetTasksByFilter(models.TaskFilter{BoardID: &boardID})
}

// MoveTask moves a task to another column at a given position, shifting others.
func (s *TaskService) MoveTask(taskID, targetColumnID, position int) error {
	// Shift tasks in target column to make room
	_, err := s.db.Exec(`
		UPDATE tasks SET position = position + 1
		WHERE column_id = ? AND position >= ?
	`, targetColumnID, position)
	if err != nil {
		return fmt.Errorf("erro ao reordenar coluna: %w", err)
	}

	_, err = s.db.Exec(`
		UPDATE tasks SET column_id = ?, position = ?, status = 'pending'
		WHERE id = ?
	`, targetColumnID, position, taskID)
	if err != nil {
		return fmt.Errorf("erro ao mover tarefa: %w", err)
	}
	return nil
}

// UpdateProgress updates only the progress field of a task.
func (s *TaskService) UpdateProgress(id, progress int) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progresso deve ser entre 0 e 100")
	}
	_, err := s.db.Exec("UPDATE tasks SET progress = ? WHERE id = ?", progress, id)
	return err
}

// --- helpers ---

func scanTasks(rows *sql.Rows) ([]models.Task, error) {
	var tasks []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler tarefa: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
