package main

import (
	"context"
	"fmt"

	"personal-cockpit/config"
	"personal-cockpit/database"
	"personal-cockpit/models"
	"personal-cockpit/services"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	db  *database.DB

	// Services
	taskService     *services.TaskService
	noteService     *services.NoteService
	eventService    *services.EventService
	categoryService *services.CategoryService
	profileService  *services.ProfileService
	settingsService *services.SettingsService
	boardService    *services.BoardService
	kanbanService   *services.KanbanService
	calendarService *services.CalendarService
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Inicializar banco de dados
	db, err := database.NewDB()
	if err != nil {
		fmt.Println("❌ Erro ao inicializar banco:", err)
		return
	}
	a.db = db

	// Inicializar todos os services
	conn := db.GetConnection()
	a.taskService = services.NewTaskService(conn)
	a.noteService = services.NewNoteService(conn)
	a.eventService = services.NewEventService(conn)
	a.categoryService = services.NewCategoryService(conn)
	a.profileService = services.NewProfileService(conn)
	a.settingsService = services.NewSettingsService(conn)
	a.boardService = services.NewBoardService(conn)
	a.kanbanService = services.NewKanbanService(conn)

	cfg := config.MustLoad()
	a.calendarService = services.NewCalendarService(conn, cfg.GoogleClientID, cfg.GoogleClientSecret)

	fmt.Println("✅ App inicializado com sucesso!")
}

func (a App) domReady(ctx context.Context) {

}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.db != nil {
		a.db.Close()
	}
	return false
}

func (a *App) shutdown(ctx context.Context) {

}

// ═══════════════════════════════════════════════════════════
// TASK METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) CreateTask(task models.Task) (int64, error) {
	return a.taskService.CreateTask(task)
}

func (a *App) GetAllTasks() ([]models.Task, error) {
	return a.taskService.GetAllTasks()
}

func (a *App) GetTaskByID(id int) (*models.Task, error) {
	return a.taskService.GetTaskByID(id)
}

func (a *App) UpdateTask(task models.Task) error {
	return a.taskService.UpdateTask(task)
}

func (a *App) DeleteTask(id int) error {
	return a.taskService.DeleteTask(id)
}

func (a *App) ToggleTaskStatus(id int) error {
	return a.taskService.ToggleTaskStatus(id)
}

func (a *App) GetPendingTasks() ([]models.Task, error) {
	return a.taskService.GetPendingTasks()
}

func (a *App) GetCompletedTasks() ([]models.Task, error) {
	return a.taskService.GetCompletedTasks()
}

func (a *App) GetTasksByFilter(filter models.TaskFilter) ([]models.Task, error) {
	return a.taskService.GetTasksByFilter(filter)
}

// ═══════════════════════════════════════════════════════════
// NOTE METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) CreateNote(note models.Note) (int64, error) {
	return a.noteService.CreateNote(note)
}

func (a *App) GetAllNotes() ([]models.Note, error) {
	return a.noteService.GetAllNotes()
}

func (a *App) GetNoteByID(id int) (*models.Note, error) {
	return a.noteService.GetNoteByID(id)
}

func (a *App) UpdateNote(note models.Note) error {
	return a.noteService.UpdateNote(note)
}

func (a *App) DeleteNote(id int) error {
	return a.noteService.DeleteNote(id)
}

func (a *App) ToggleNoteFavorite(id int) error {
	return a.noteService.ToggleFavorite(id)
}

func (a *App) GetFavoriteNotes() ([]models.Note, error) {
	return a.noteService.GetFavoriteNotes()
}

func (a *App) SearchNotes(query string) ([]models.Note, error) {
	return a.noteService.SearchNotes(query)
}

// ═══════════════════════════════════════════════════════════
// EVENT METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) CreateEvent(event models.Event) (int64, error) {
	return a.eventService.CreateEvent(event)
}

func (a *App) GetAllEvents() ([]models.Event, error) {
	return a.eventService.GetAllEvents()
}

func (a *App) GetEventByID(id int) (*models.Event, error) {
	return a.eventService.GetEventByID(id)
}

func (a *App) UpdateEvent(event models.Event) error {
	return a.eventService.UpdateEvent(event)
}

func (a *App) DeleteEvent(id int) error {
	return a.eventService.DeleteEvent(id)
}

func (a *App) GetTodayEvents() ([]models.Event, error) {
	return a.eventService.GetTodayEvents()
}

func (a *App) GetUpcomingEvents() ([]models.Event, error) {
	return a.eventService.GetUpcomingEvents()
}

// ═══════════════════════════════════════════════════════════
// CATEGORY METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) CreateCategory(category models.Category) (int64, error) {
	return a.categoryService.CreateCategory(category)
}

func (a *App) GetAllCategories() ([]models.Category, error) {
	return a.categoryService.GetAllCategories()
}

func (a *App) GetCategoryByID(id int) (*models.Category, error) {
	return a.categoryService.GetCategoryByID(id)
}

func (a *App) UpdateCategory(category models.Category) error {
	return a.categoryService.UpdateCategory(category)
}

func (a *App) DeleteCategory(id int) error {
	return a.categoryService.DeleteCategory(id)
}

func (a *App) GetTaskCategories() ([]models.Category, error) {
	return a.categoryService.GetTaskCategories()
}

func (a *App) GetNoteCategories() ([]models.Category, error) {
	return a.categoryService.GetNoteCategories()
}

// ═══════════════════════════════════════════════════════════
// APP INFO
// ═══════════════════════════════════════════════════════════

func (a *App) GetAppInfo() map[string]string {
	cfg := config.MustLoad()
	return map[string]string{
		"name":    cfg.AppName,
		"version": cfg.FullVersion(),
		"year":    cfg.CurrentYear,
	}
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Olá %s! Bem-vindo ao Personal Cockpit v%s", name, config.MustLoad().FullVersion())
}

// ═══════════════════════════════════════════════════════════
// THEME METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) SetTheme(theme string) error {
	// Cores para cada tema
	themes := map[string]struct {
		bg     int64
		text   int64
		border int64
	}{
		"light": {
			bg:     0xFFFFFF, // Branco
			text:   0x1D1D1F, // Preto
			border: 0xE5E5EA, // Cinza claro
		},
		"dark": {
			bg:     0x282828, // Cinza escuro
			text:   0xF5F5F7, // Branco
			border: 0x2C2C2E, // Cinza mais escuro
		},
	}

	selectedTheme, exists := themes[theme]
	if !exists {
		// Se não existir ou for "auto", usa detecção do sistema
		// Por padrão, vamos usar dark
		selectedTheme = themes["dark"]
	}

	// Nota: O Wails v2 não tem API direta para mudar a cor da titlebar em runtime
	// Isso precisa ser configurado no wails.json
	// Mas podemos retornar sucesso para o frontend saber que o tema foi "aplicado"

	fmt.Printf("✅ Tema alterado para: %s (bg: %x)\n", theme, selectedTheme.bg)
	return nil
}

// ═══════════════════════════════════════════════════════════
// PROFILE METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) CreateProfile(input models.CreateProfileInput) (*models.ProfileSummary, error) {
	return a.profileService.CreateProfile(input)
}

func (a *App) GetAllProfiles() ([]models.ProfileSummary, error) {
	return a.profileService.GetAllProfiles()
}

func (a *App) GetProfileByID(id int) (*models.ProfileSummary, error) {
	return a.profileService.GetProfileSummaryByID(id)
}

func (a *App) UpdateProfile(id int, input models.CreateProfileInput) error {
	return a.profileService.UpdateProfile(id, input)
}

func (a *App) DeleteProfile(id int) error {
	return a.profileService.DeleteProfile(id)
}

func (a *App) VerifyPin(profileID int, pin string) (bool, error) {
	return a.profileService.VerifyPin(profileID, pin)
}

func (a *App) SetPin(profileID int, pin string) error {
	return a.profileService.SetPin(profileID, pin)
}

func (a *App) SetCurrentProfile(id int) error {
	return a.profileService.SetCurrentProfile(id)
}

func (a *App) GetCurrentProfileID() int {
	return a.profileService.GetCurrentProfileID()
}

func (a *App) HasAnyProfile() (bool, error) {
	return a.profileService.HasAnyProfile()
}

// ═══════════════════════════════════════════════════════════
// BOARD METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) CreateBoard(name string, profileID int, isShared bool) (*models.BoardWithColumns, error) {
	return a.boardService.CreateBoard(name, profileID, isShared)
}

func (a *App) GetBoardsByProfile(profileID int, role string) ([]models.BoardWithColumns, error) {
	return a.boardService.GetBoardsByProfile(profileID, role)
}

func (a *App) GetBoardWithColumns(boardID int) (*models.BoardWithColumns, error) {
	return a.boardService.GetBoardWithColumns(boardID)
}

func (a *App) UpdateBoard(id int, name string, isShared bool) error {
	return a.boardService.UpdateBoard(id, name, isShared)
}

func (a *App) DeleteBoard(id int) error {
	return a.boardService.DeleteBoard(id)
}

// ─── Columns ───────────────────────────────────────────────

func (a *App) CreateColumn(boardID int, name, color string, wipLimit *int) (*models.KanbanColumn, error) {
	return a.boardService.CreateColumn(boardID, name, color, wipLimit)
}

func (a *App) UpdateColumn(id int, name, color string, wipLimit *int) error {
	return a.boardService.UpdateColumn(id, name, color, wipLimit)
}

func (a *App) DeleteColumn(id int) error {
	return a.boardService.DeleteColumn(id)
}

func (a *App) ReorderColumns(columnIDs []int) error {
	return a.boardService.ReorderColumns(columnIDs)
}

// ─── Task Kanban actions ────────────────────────────────────

func (a *App) MoveTask(taskID, targetColumnID, position int) error {
	return a.taskService.MoveTask(taskID, targetColumnID, position)
}

func (a *App) GetTasksByBoard(boardID int) ([]models.Task, error) {
	return a.taskService.GetTasksByBoard(boardID)
}

func (a *App) GetTasksByColumn(columnID int) ([]models.Task, error) {
	return a.taskService.GetTasksByColumn(columnID)
}

func (a *App) UpdateProgress(taskID, progress int) error {
	return a.taskService.UpdateProgress(taskID, progress)
}

// ═══════════════════════════════════════════════════════════
// KANBAN DETAIL METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) GetTaskDetail(taskID int) (*models.TaskDetail, error) {
	return a.kanbanService.GetTaskDetail(taskID)
}

func (a *App) CreateSubtask(taskID int, title string) (*models.Subtask, error) {
	return a.kanbanService.CreateSubtask(taskID, title)
}

func (a *App) ToggleSubtask(subtaskID int) error {
	return a.kanbanService.ToggleSubtask(subtaskID)
}

func (a *App) DeleteSubtask(subtaskID int) error {
	return a.kanbanService.DeleteSubtask(subtaskID)
}

func (a *App) AddAssignee(taskID, profileID int) error {
	return a.kanbanService.AddAssignee(taskID, profileID)
}

func (a *App) RemoveAssignee(taskID, profileID int) error {
	return a.kanbanService.RemoveAssignee(taskID, profileID)
}

func (a *App) AddDependency(taskID, dependsOnID int) error {
	return a.kanbanService.AddDependency(taskID, dependsOnID)
}

func (a *App) RemoveDependency(taskID, dependsOnID int) error {
	return a.kanbanService.RemoveDependency(taskID, dependsOnID)
}

func (a *App) GetLabelsByBoard(boardID int) ([]models.Label, error) {
	return a.kanbanService.GetLabelsByBoard(boardID)
}

func (a *App) CreateLabel(boardID int, name, color string) (*models.Label, error) {
	return a.kanbanService.CreateLabel(boardID, name, color)
}

func (a *App) DeleteLabel(id int) error {
	return a.kanbanService.DeleteLabel(id)
}

func (a *App) AddTaskLabel(taskID, labelID int) error {
	return a.kanbanService.AddTaskLabel(taskID, labelID)
}

func (a *App) RemoveTaskLabel(taskID, labelID int) error {
	return a.kanbanService.RemoveTaskLabel(taskID, labelID)
}

func (a *App) AddComment(taskID, profileID int, content string) (*models.TaskComment, error) {
	return a.kanbanService.AddComment(taskID, profileID, content)
}

func (a *App) DeleteComment(commentID int) error {
	return a.kanbanService.DeleteComment(commentID)
}

func (a *App) GetActivity(taskID int) ([]models.TaskActivity, error) {
	return a.kanbanService.GetActivity(taskID)
}

// ═══════════════════════════════════════════════════════════
// SETTINGS METHODS
// ═══════════════════════════════════════════════════════════

func (a *App) GetSetting(key string) (string, error) {
	return a.settingsService.Get(key)
}

func (a *App) SetSetting(key, value string) error {
	return a.settingsService.Set(key, value)
}

// ═══════════════════════════════════════════════════════════
// CALENDAR METHODS
// ═══════════════════════════════════════════════════════════

// GoogleCalendarConfigured returns true if GOOGLE_CLIENT_ID/SECRET are set in .env.
func (a *App) GoogleCalendarConfigured() bool {
	return a.calendarService.IsConfigured()
}

// GoogleCalendarConnected returns true if the profile has a stored token.
func (a *App) GoogleCalendarConnected(profileID int) bool {
	return a.calendarService.HasToken(profileID)
}

// StartGoogleAuth opens the browser with the OAuth consent page and starts
// a local callback server to capture the token automatically.
func (a *App) StartGoogleAuth(profileID int) error {
	authURL, err := a.calendarService.StartOAuthFlow(profileID)
	if err != nil {
		return err
	}
	wailsruntime.BrowserOpenURL(a.ctx, authURL)
	return nil
}

// DisconnectGoogleCalendar removes the saved token for a profile.
func (a *App) DisconnectGoogleCalendar(profileID int) error {
	return a.calendarService.DisconnectProfile(profileID)
}

// SyncEventToCalendar creates a Google Calendar event and saves the google_event_id.
func (a *App) SyncEventToCalendar(eventID, profileID int) error {
	event, err := a.eventService.GetEventByID(eventID)
	if err != nil {
		return err
	}
	gID, err := a.calendarService.CreateEventFromEvent(profileID, *event)
	if err != nil {
		return err
	}
	event.GoogleEventID = &gID
	return a.eventService.UpdateEvent(*event)
}

// SyncTaskToCalendar creates a Google Calendar event from a task and saves the google_event_id.
func (a *App) SyncTaskToCalendar(taskID, profileID int) error {
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	gID, err := a.calendarService.CreateEventFromTask(profileID, *task)
	if err != nil {
		return err
	}
	task.GoogleEventID = &gID
	return a.taskService.UpdateTask(*task)
}

// DeleteFromCalendar removes an event from Google Calendar by its google_event_id.
func (a *App) DeleteFromCalendar(profileID int, googleEventID string) error {
	return a.calendarService.DeleteCalendarEvent(profileID, googleEventID)
}
