package services

import (
	"database/sql"
	"fmt"
	"strings"

	"personal-cockpit/models"
)

// TaskService gerencia operações de tarefas
type TaskService struct {
	db *sql.DB
}

// NewTaskService cria novo serviço de tarefas
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

// CreateTask cria uma nova tarefa
func (s *TaskService) CreateTask(task models.Task) (int64, error) {

	// Validações
	var erros []string

	if task.Title == "" {
		erros = append(erros, "título")
	}

	if task.Priority == "" {
		erros = append(erros, "prioridade")
	}

	if task.Description == "" {
		erros = append(erros, "descrição")
	}

	// Se tem erros, retornar todos de uma vez
	if len(erros) > 0 {
		mensagem := "Campos obrigatórios:\n- " + strings.Join(erros, "\n- ")
		return 0, fmt.Errorf(mensagem)
	}

	// Query SQL
	query := `
		INSERT INTO tasks (title, description, status, priority, category_id, due_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.CategoryID,
		task.DueDate,
	)

	if err != nil {
		return 0, fmt.Errorf("erro ao criar tarefa: %w", err)
	}

	// Retornar ID da tarefa criada
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter ID: %w", err)
	}

	return id, nil
}

// GetAllTasks retorna todas as tarefas
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, category_id, 
		       due_date, completed_at, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CategoryID,
			&task.DueDate,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler tarefa: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID busca tarefa por ID
func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, category_id, 
		       due_date, completed_at, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	var task models.Task
	err := s.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CategoryID,
		&task.DueDate,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tarefa não encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}

	return &task, nil
}

// UpdateTask atualiza uma tarefa
func (s *TaskService) UpdateTask(task models.Task) error {
	if task.ID == 0 {
		return fmt.Errorf("ID da tarefa é obrigatório")
	}

	query := `
		UPDATE tasks 
		SET title = ?, description = ?, status = ?, priority = ?, 
		    category_id = ?, due_date = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.CategoryID,
		task.DueDate,
		task.ID,
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

// DeleteTask deleta uma tarefa
func (s *TaskService) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"

	result, err := s.db.Exec(query, id)
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

// ToggleTaskStatus alterna status entre pending e completed
func (s *TaskService) ToggleTaskStatus(id int) error {
	// Buscar tarefa
	task, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}

	// Alternar status
	newStatus := "pending"
	if task.Status == "pending" {
		newStatus = "completed"
	}

	query := "UPDATE tasks SET status = ? WHERE id = ?"
	_, err = s.db.Exec(query, newStatus, id)

	return err
}

// GetTasksByFilter busca tarefas com filtros
func (s *TaskService) GetTasksByFilter(filter models.TaskFilter) ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, category_id, 
		       due_date, completed_at, created_at, updated_at
		FROM tasks
		WHERE 1=1
	`

	args := []interface{}{}

	// Filtro por status
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	// Filtro por prioridade
	if filter.Priority != "" {
		query += " AND priority = ?"
		args = append(args, filter.Priority)
	}

	// Filtro por categoria
	if filter.CategoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *filter.CategoryID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CategoryID,
			&task.DueDate,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler tarefa: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetPendingTasks retorna apenas tarefas pendentes
func (s *TaskService) GetPendingTasks() ([]models.Task, error) {
	filter := models.TaskFilter{Status: "pending"}
	return s.GetTasksByFilter(filter)
}

// GetCompletedTasks retorna apenas tarefas concluídas
func (s *TaskService) GetCompletedTasks() ([]models.Task, error) {
	filter := models.TaskFilter{Status: "completed"}
	return s.GetTasksByFilter(filter)
}
