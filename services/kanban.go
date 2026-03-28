package services

import (
	"database/sql"
	"fmt"

	"personal-cockpit/models"
)

type KanbanService struct {
	db *sql.DB
}

func NewKanbanService(db *sql.DB) *KanbanService {
	return &KanbanService{db: db}
}

// ─── Subtasks ──────────────────────────────────────────────

func (s *KanbanService) GetSubtasks(taskID int) ([]models.Subtask, error) {
	rows, err := s.db.Query(`
		SELECT id, task_id, title, completed, position, created_at
		FROM subtasks WHERE task_id = ? ORDER BY position ASC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Subtask
	for rows.Next() {
		var st models.Subtask
		if err := rows.Scan(&st.ID, &st.TaskID, &st.Title, &st.Completed, &st.Position, &st.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, st)
	}
	return items, nil
}

func (s *KanbanService) CreateSubtask(taskID int, title string) (*models.Subtask, error) {
	if title == "" {
		return nil, fmt.Errorf("título da subtarefa é obrigatório")
	}
	var maxPos int
	s.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM subtasks WHERE task_id = ?`, taskID).Scan(&maxPos)

	result, err := s.db.Exec(`
		INSERT INTO subtasks (task_id, title, position) VALUES (?, ?, ?)
	`, taskID, title, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar subtarefa: %w", err)
	}
	id, _ := result.LastInsertId()

	var st models.Subtask
	s.db.QueryRow(`SELECT id, task_id, title, completed, position, created_at FROM subtasks WHERE id = ?`, id).
		Scan(&st.ID, &st.TaskID, &st.Title, &st.Completed, &st.Position, &st.CreatedAt)
	return &st, nil
}

func (s *KanbanService) ToggleSubtask(id int) error {
	_, err := s.db.Exec(`
		UPDATE subtasks SET completed = CASE WHEN completed = 1 THEN 0 ELSE 1 END WHERE id = ?
	`, id)
	return err
}

func (s *KanbanService) DeleteSubtask(id int) error {
	_, err := s.db.Exec(`DELETE FROM subtasks WHERE id = ?`, id)
	return err
}

// SubtaskProgress returns (completed, total) for a task.
func (s *KanbanService) SubtaskProgress(taskID int) (int, int, error) {
	var total, completed int
	s.db.QueryRow(`SELECT COUNT(*), SUM(CASE WHEN completed = 1 THEN 1 ELSE 0 END) FROM subtasks WHERE task_id = ?`, taskID).
		Scan(&total, &completed)
	return completed, total, nil
}

// ─── Assignees ─────────────────────────────────────────────

func (s *KanbanService) GetAssignees(taskID int) ([]models.ProfileSummary, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.name, p.email, p.role, p.avatar_url, p.google_calendar_id,
		       CASE WHEN p.google_access_token IS NOT NULL AND p.google_access_token != '' THEN 1 ELSE 0 END,
		       CASE WHEN p.pin_hash IS NOT NULL AND p.pin_hash != '' THEN 1 ELSE 0 END,
		       p.created_at
		FROM profiles p
		JOIN task_assignees ta ON ta.profile_id = p.id
		WHERE ta.task_id = ?
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.ProfileSummary
	for rows.Next() {
		var p models.ProfileSummary
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Role, &p.AvatarURL,
			&p.GoogleCalendarID, &p.HasGoogleToken, &p.HasPin, &p.CreatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (s *KanbanService) AddAssignee(taskID, profileID int) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO task_assignees (task_id, profile_id) VALUES (?, ?)
	`, taskID, profileID)
	return err
}

func (s *KanbanService) RemoveAssignee(taskID, profileID int) error {
	_, err := s.db.Exec(`DELETE FROM task_assignees WHERE task_id = ? AND profile_id = ?`, taskID, profileID)
	return err
}

// ─── Labels ────────────────────────────────────────────────

func (s *KanbanService) GetLabelsByBoard(boardID int) ([]models.Label, error) {
	rows, err := s.db.Query(`SELECT id, name, color, board_id FROM labels WHERE board_id = ?`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var l models.Label
		if err := rows.Scan(&l.ID, &l.Name, &l.Color, &l.BoardID); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (s *KanbanService) CreateLabel(boardID int, name, color string) (*models.Label, error) {
	result, err := s.db.Exec(`INSERT INTO labels (name, color, board_id) VALUES (?, ?, ?)`, name, color, boardID)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &models.Label{ID: int(id), Name: name, Color: color, BoardID: boardID}, nil
}

func (s *KanbanService) DeleteLabel(id int) error {
	_, err := s.db.Exec(`DELETE FROM labels WHERE id = ?`, id)
	return err
}

func (s *KanbanService) GetTaskLabels(taskID int) ([]models.Label, error) {
	rows, err := s.db.Query(`
		SELECT l.id, l.name, l.color, l.board_id
		FROM labels l JOIN task_labels tl ON tl.label_id = l.id
		WHERE tl.task_id = ?
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var l models.Label
		if err := rows.Scan(&l.ID, &l.Name, &l.Color, &l.BoardID); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (s *KanbanService) AddTaskLabel(taskID, labelID int) error {
	_, err := s.db.Exec(`INSERT OR IGNORE INTO task_labels (task_id, label_id) VALUES (?, ?)`, taskID, labelID)
	return err
}

func (s *KanbanService) RemoveTaskLabel(taskID, labelID int) error {
	_, err := s.db.Exec(`DELETE FROM task_labels WHERE task_id = ? AND label_id = ?`, taskID, labelID)
	return err
}

// ─── Dependencies ──────────────────────────────────────────

func (s *KanbanService) GetDependencies(taskID int) ([]int, error) {
	rows, err := s.db.Query(`SELECT depends_on_id FROM task_dependencies WHERE task_id = ?`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *KanbanService) AddDependency(taskID, dependsOnID int) error {
	_, err := s.db.Exec(`INSERT OR IGNORE INTO task_dependencies (task_id, depends_on_id) VALUES (?, ?)`, taskID, dependsOnID)
	return err
}

func (s *KanbanService) RemoveDependency(taskID, dependsOnID int) error {
	_, err := s.db.Exec(`DELETE FROM task_dependencies WHERE task_id = ? AND depends_on_id = ?`, taskID, dependsOnID)
	return err
}

// ─── Comments ──────────────────────────────────────────────

func (s *KanbanService) GetComments(taskID int) ([]models.TaskComment, error) {
	rows, err := s.db.Query(`
		SELECT tc.id, tc.task_id, tc.profile_id, tc.content, tc.created_at,
		       p.name, p.avatar_url
		FROM task_comments tc
		LEFT JOIN profiles p ON p.id = tc.profile_id
		WHERE tc.task_id = ? ORDER BY tc.created_at ASC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.TaskComment
	for rows.Next() {
		var c models.TaskComment
		var pName, pAvatar sql.NullString
		if err := rows.Scan(&c.ID, &c.TaskID, &c.ProfileID, &c.Content, &c.CreatedAt, &pName, &pAvatar); err != nil {
			return nil, err
		}
		if pName.Valid {
			c.Profile = &models.ProfileSummary{ID: c.ProfileID, Name: pName.String, AvatarURL: pAvatar.String}
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *KanbanService) AddComment(taskID, profileID int, content string) (*models.TaskComment, error) {
	if content == "" {
		return nil, fmt.Errorf("comentário não pode ser vazio")
	}
	result, err := s.db.Exec(`
		INSERT INTO task_comments (task_id, profile_id, content) VALUES (?, ?, ?)
	`, taskID, profileID, content)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	comments, err := s.GetComments(taskID)
	if err != nil {
		return nil, err
	}
	for _, c := range comments {
		if c.ID == int(id) {
			return &c, nil
		}
	}
	return nil, nil
}

func (s *KanbanService) DeleteComment(id int) error {
	_, err := s.db.Exec(`DELETE FROM task_comments WHERE id = ?`, id)
	return err
}

// ─── Activity ──────────────────────────────────────────────

func (s *KanbanService) LogActivity(taskID int, profileID *int, action, from, to string) {
	s.db.Exec(`
		INSERT INTO task_activity (task_id, profile_id, action, from_value, to_value)
		VALUES (?, ?, ?, ?, ?)
	`, taskID, profileID, action, from, to)
}

func (s *KanbanService) GetActivity(taskID int) ([]models.TaskActivity, error) {
	rows, err := s.db.Query(`
		SELECT ta.id, ta.task_id, ta.profile_id, ta.action, ta.from_value, ta.to_value, ta.created_at,
		       p.name, p.avatar_url
		FROM task_activity ta
		LEFT JOIN profiles p ON p.id = ta.profile_id
		WHERE ta.task_id = ? ORDER BY ta.created_at DESC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.TaskActivity
	for rows.Next() {
		var a models.TaskActivity
		var pName, pAvatar sql.NullString
		if err := rows.Scan(&a.ID, &a.TaskID, &a.ProfileID, &a.Action,
			&a.FromValue, &a.ToValue, &a.CreatedAt, &pName, &pAvatar); err != nil {
			return nil, err
		}
		if pName.Valid {
			pid := 0
			if a.ProfileID != nil {
				pid = *a.ProfileID
			}
			a.Profile = &models.ProfileSummary{ID: pid, Name: pName.String, AvatarURL: pAvatar.String}
		}
		activities = append(activities, a)
	}
	return activities, nil
}

// ─── Task Detail (all relations) ───────────────────────────

func (s *KanbanService) GetTaskDetail(taskID int) (*models.TaskDetail, error) {
	detail := &models.TaskDetail{}

	assignees, err := s.GetAssignees(taskID)
	if err != nil {
		return nil, err
	}
	detail.Assignees = assignees

	subtasks, err := s.GetSubtasks(taskID)
	if err != nil {
		return nil, err
	}
	detail.Subtasks = subtasks

	labels, err := s.GetTaskLabels(taskID)
	if err != nil {
		return nil, err
	}
	detail.Labels = labels

	comments, err := s.GetComments(taskID)
	if err != nil {
		return nil, err
	}
	detail.Comments = comments

	deps, err := s.GetDependencies(taskID)
	if err != nil {
		return nil, err
	}
	detail.Dependencies = deps

	// Dependents: tasks that depend on this one
	rows, err := s.db.Query(`SELECT task_id FROM task_dependencies WHERE depends_on_id = ?`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		detail.Dependents = append(detail.Dependents, id)
	}

	return detail, nil
}
