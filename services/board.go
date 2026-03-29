package services

import (
	"database/sql"
	"fmt"

	"personal-cockpit/models"
)

type BoardService struct {
	db *sql.DB
}

func NewBoardService(db *sql.DB) *BoardService {
	return &BoardService{db: db}
}

// CreateBoard creates a board and seeds it with default Kanban columns.
func (s *BoardService) CreateBoard(name string, profileID int, isShared bool) (*models.BoardWithColumns, error) {
	if name == "" {
		return nil, fmt.Errorf("nome do board é obrigatório")
	}
	result, err := s.db.Exec(`
		INSERT INTO boards (name, profile_id, is_shared) VALUES (?, ?, ?)
	`, name, profileID, isShared)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar board: %w", err)
	}
	boardID, _ := result.LastInsertId()

	defaults := []struct {
		name  string
		color string
	}{
		{"A fazer", "#6b7280"},
		{"Em progresso", "#3b82f6"},
		{"Revisão", "#f59e0b"},
		{"Concluído", "#10b981"},
	}
	for i, col := range defaults {
		_, err := s.db.Exec(`
			INSERT INTO kanban_columns (board_id, name, color, position) VALUES (?, ?, ?, ?)
		`, boardID, col.name, col.color, i)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar coluna padrão: %w", err)
		}
	}

	return s.GetBoardWithColumns(int(boardID))
}

func (s *BoardService) GetBoardWithColumns(boardID int) (*models.BoardWithColumns, error) {
	var b models.BoardWithColumns
	err := s.db.QueryRow(`
		SELECT id, name, profile_id, is_shared, created_at FROM boards WHERE id = ?
	`, boardID).Scan(&b.ID, &b.Name, &b.ProfileID, &b.IsShared, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("board não encontrado")
	}
	if err != nil {
		return nil, err
	}

	cols, err := s.GetColumnsByBoard(boardID)
	if err != nil {
		return nil, err
	}
	b.Columns = cols
	return &b, nil
}

// GetBoardsByProfile returns all boards visible to a profile.
// Admin sees all boards; members see only their own + shared boards.
func (s *BoardService) GetBoardsByProfile(profileID int, role string) ([]models.BoardWithColumns, error) {
	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = s.db.Query(`SELECT id FROM boards ORDER BY created_at ASC`)
	} else {
		rows, err = s.db.Query(`
			SELECT id FROM boards
			WHERE profile_id = ? OR is_shared = 1
			ORDER BY created_at ASC
		`, profileID)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar boards: %w", err)
	}
	defer rows.Close()

	var boards []models.BoardWithColumns
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		b, err := s.GetBoardWithColumns(id)
		if err != nil {
			return nil, err
		}
		boards = append(boards, *b)
	}
	return boards, nil
}

func (s *BoardService) UpdateBoard(id int, name string, isShared bool) error {
	result, err := s.db.Exec(`UPDATE boards SET name = ?, is_shared = ? WHERE id = ?`, name, isShared, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar board: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("board não encontrado")
	}
	return nil
}

func (s *BoardService) DeleteBoard(id int) error {
	_, err := s.db.Exec(`DELETE FROM boards WHERE id = ?`, id)
	return err
}

// ─── Columns ───────────────────────────────────────────────

func (s *BoardService) GetColumnsByBoard(boardID int) ([]models.KanbanColumn, error) {
	rows, err := s.db.Query(`
		SELECT id, board_id, name, color, position, wip_limit, created_at
		FROM kanban_columns WHERE board_id = ? ORDER BY position ASC
	`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []models.KanbanColumn
	for rows.Next() {
		var c models.KanbanColumn
		if err := rows.Scan(&c.ID, &c.BoardID, &c.Name, &c.Color, &c.Position, &c.WIPLimit, &c.CreatedAt); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, nil
}

func (s *BoardService) CreateColumn(boardID int, name, color string, wipLimit *int) (*models.KanbanColumn, error) {
	if name == "" {
		return nil, fmt.Errorf("nome da coluna é obrigatório")
	}
	var maxPos int
	s.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM kanban_columns WHERE board_id = ?`, boardID).Scan(&maxPos)

	result, err := s.db.Exec(`
		INSERT INTO kanban_columns (board_id, name, color, position, wip_limit)
		VALUES (?, ?, ?, ?, ?)
	`, boardID, name, color, maxPos+1, wipLimit)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar coluna: %w", err)
	}
	id, _ := result.LastInsertId()

	var c models.KanbanColumn
	s.db.QueryRow(`
		SELECT id, board_id, name, color, position, wip_limit, created_at
		FROM kanban_columns WHERE id = ?
	`, id).Scan(&c.ID, &c.BoardID, &c.Name, &c.Color, &c.Position, &c.WIPLimit, &c.CreatedAt)
	return &c, nil
}

func (s *BoardService) UpdateColumn(id int, name, color string, wipLimit *int) error {
	_, err := s.db.Exec(`
		UPDATE kanban_columns SET name = ?, color = ?, wip_limit = ? WHERE id = ?
	`, name, color, wipLimit, id)
	return err
}

func (s *BoardService) DeleteColumn(id int) error {
	// Unlink tasks from this column before deleting
	s.db.Exec(`UPDATE tasks SET column_id = NULL WHERE column_id = ?`, id)
	_, err := s.db.Exec(`DELETE FROM kanban_columns WHERE id = ?`, id)
	return err
}

// ReorderColumns updates the position of each column given an ordered list of IDs.
func (s *BoardService) ReorderColumns(columnIDs []int) error {
	for i, id := range columnIDs {
		if _, err := s.db.Exec(`UPDATE kanban_columns SET position = ? WHERE id = ?`, i, id); err != nil {
			return fmt.Errorf("erro ao reordenar colunas: %w", err)
		}
	}
	return nil
}
