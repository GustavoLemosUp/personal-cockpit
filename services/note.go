package services

import (
	"database/sql"
	"fmt"
	"strings"

	"personal-cockpit/models"
)

type NoteService struct {
	db *sql.DB
}

func NewNoteService(db *sql.DB) *NoteService {
	return &NoteService{db: db}
}

func (s *NoteService) CreateNote(note models.Note) (int64, error) {
	var erros []string

	if note.Title == "" {
		erros = append(erros, "título")
	}

	if len(erros) > 0 {
		mensagem := "Campos obrigatórios:\n- " + strings.Join(erros, "\n- ")
		return 0, fmt.Errorf(mensagem)
	}

	query := `
		INSERT INTO notes (title, content, category_id, is_favorite)
		VALUES (?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		note.Title,
		note.Content,
		note.CategoryID,
		note.IsFavorite,
	)

	if err != nil {
		return 0, fmt.Errorf("erro ao criar nota: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter ID: %w", err)
	}

	return id, nil
}

func (s *NoteService) GetAllNotes() ([]models.Note, error) {
	query := `
		SELECT id, title, content, category_id, is_favorite, created_at, updated_at
		FROM notes
		ORDER BY updated_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar notas: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CategoryID,
			&note.IsFavorite,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler nota: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *NoteService) GetNoteByID(id int) (*models.Note, error) {
	query := `
		SELECT id, title, content, category_id, is_favorite, created_at, updated_at
		FROM notes
		WHERE id = ?
	`

	var note models.Note
	err := s.db.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CategoryID,
		&note.IsFavorite,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("nota não encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar nota: %w", err)
	}

	return &note, nil
}

func (s *NoteService) UpdateNote(note models.Note) error {
	if note.ID == 0 {
		return fmt.Errorf("ID da nota é obrigatório")
	}

	query := `
		UPDATE notes 
		SET title = ?, content = ?, category_id = ?, is_favorite = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		note.Title,
		note.Content,
		note.CategoryID,
		note.IsFavorite,
		note.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar nota: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("nota não encontrada")
	}

	return nil
}

func (s *NoteService) DeleteNote(id int) error {
	query := "DELETE FROM notes WHERE id = ?"

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar nota: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("nota não encontrada")
	}

	return nil
}

func (s *NoteService) ToggleFavorite(id int) error {
	note, err := s.GetNoteByID(id)
	if err != nil {
		return err
	}

	query := "UPDATE notes SET is_favorite = ? WHERE id = ?"
	_, err = s.db.Exec(query, !note.IsFavorite, id)

	return err
}

func (s *NoteService) GetFavoriteNotes() ([]models.Note, error) {
	query := `
		SELECT id, title, content, category_id, is_favorite, created_at, updated_at
		FROM notes
		WHERE is_favorite = 1
		ORDER BY updated_at DESC
	`
	// SQLite: 1 = true, 0 = false

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar notas: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CategoryID,
			&note.IsFavorite,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler nota: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *NoteService) SearchNotes(searchQuery string) ([]models.Note, error) {
	sqlQuery := `
		SELECT id, title, content, category_id, is_favorite, created_at, updated_at
		FROM notes
		WHERE title LIKE ? OR content LIKE ?
		ORDER BY updated_at DESC
	`

	searchTerm := "%" + searchQuery + "%"

	rows, err := s.db.Query(sqlQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar notas: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CategoryID,
			&note.IsFavorite,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler nota: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}
