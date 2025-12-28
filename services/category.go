package services

import (
	"database/sql"
	"fmt"
	"strings"

	"personal-cockpit/models"
)

// CategoryService gerencia operações de categorias
type CategoryService struct {
	db *sql.DB
}

// NewCategoryService cria novo serviço de categorias
func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{db: db}
}

// CreateCategory cria uma nova categoria
func (s *CategoryService) CreateCategory(category models.Category) (int64, error) {
	// Validações
	var erros []string

	if category.Name == "" {
		erros = append(erros, "nome")
	}

	if len(erros) > 0 {
		mensagem := "Campos obrigatórios:\n- " + strings.Join(erros, "\n- ")
		return 0, fmt.Errorf(mensagem)
	}

	// Query SQL
	query := `
		INSERT INTO categories (name, color, type)
		VALUES (?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		category.Name,
		category.Color,
		category.Type,
	)

	if err != nil {
		return 0, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter ID: %w", err)
	}

	return id, nil
}

// GetAllCategories retorna todas as categorias
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	query := `
		SELECT id, name, color, type, created_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Color,
			&category.Type,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler categoria: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryByID busca categoria por ID
func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	query := `
		SELECT id, name, color, type, created_at
		FROM categories
		WHERE id = ?
	`

	var category models.Category
	err := s.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Color,
		&category.Type,
		&category.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("categoria não encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	return &category, nil
}

// UpdateCategory atualiza uma categoria
func (s *CategoryService) UpdateCategory(category models.Category) error {
	if category.ID == 0 {
		return fmt.Errorf("ID da categoria é obrigatório")
	}

	query := `
		UPDATE categories 
		SET name = ?, color = ?, type = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		category.Name,
		category.Color,
		category.Type,
		category.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("categoria não encontrada")
	}

	return nil
}

// DeleteCategory deleta uma categoria
func (s *CategoryService) DeleteCategory(id int) error {
	query := "DELETE FROM categories WHERE id = ?"

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("categoria não encontrada")
	}

	return nil
}

// GetCategoriesByType busca categorias por tipo
func (s *CategoryService) GetCategoriesByType(categoryType string) ([]models.Category, error) {
	query := `
		SELECT id, name, color, type, created_at
		FROM categories
		WHERE type = ?
		ORDER BY name ASC
	`

	rows, err := s.db.Query(query, categoryType)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Color,
			&category.Type,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler categoria: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetTaskCategories retorna apenas categorias de tarefas
func (s *CategoryService) GetTaskCategories() ([]models.Category, error) {
	return s.GetCategoriesByType("task")
}

// GetNoteCategories retorna apenas categorias de notas
func (s *CategoryService) GetNoteCategories() ([]models.Category, error) {
	return s.GetCategoriesByType("note")
}
