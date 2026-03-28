package models

import "time"

type Board struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ProfileID int       `json:"profile_id"`
	IsShared  bool      `json:"is_shared"`
	CreatedAt time.Time `json:"created_at"`
}

type KanbanColumn struct {
	ID        int       `json:"id"`
	BoardID   int       `json:"board_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Position  int       `json:"position"`
	WIPLimit  *int      `json:"wip_limit"`
	CreatedAt time.Time `json:"created_at"`
}

type Subtask struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type Label struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	BoardID int    `json:"board_id"`
}

type TaskComment struct {
	ID        int            `json:"id"`
	TaskID    int            `json:"task_id"`
	ProfileID int            `json:"profile_id"`
	Profile   *ProfileSummary `json:"profile,omitempty"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type TaskActivity struct {
	ID        int            `json:"id"`
	TaskID    int            `json:"task_id"`
	ProfileID *int           `json:"profile_id"`
	Profile   *ProfileSummary `json:"profile,omitempty"`
	Action    string         `json:"action"`
	FromValue string         `json:"from_value"`
	ToValue   string         `json:"to_value"`
	CreatedAt time.Time      `json:"created_at"`
}

// TaskDetail is the full task with all relations loaded.
type TaskDetail struct {
	Assignees    []ProfileSummary `json:"assignees"`
	Subtasks     []Subtask        `json:"subtasks"`
	Labels       []Label          `json:"labels"`
	Comments     []TaskComment    `json:"comments"`
	Dependencies []int            `json:"dependencies"` // IDs das tasks das quais esta depende
	Dependents   []int            `json:"dependents"`   // IDs das tasks que dependem desta
}

// BoardWithColumns is a board with its columns pre-loaded.
type BoardWithColumns struct {
	Board
	Columns []KanbanColumn `json:"columns"`
}

// MoveTaskInput represents a drag & drop action.
type MoveTaskInput struct {
	TaskID         int `json:"task_id"`
	TargetColumnID int `json:"target_column_id"`
	Position       int `json:"position"`
}
