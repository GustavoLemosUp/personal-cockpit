package main

import (
	"context"
	"fmt"

	"personal-cockpit/database"
	"personal-cockpit/models"
	"personal-cockpit/services"
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
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
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

	fmt.Println("✅ App inicializado com sucesso!")
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.db != nil {
		a.db.Close()
	}
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
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
	return map[string]string{
		"name":    AppName,
		"version": GetFullVersion(),
		"year":    CurrentYear,
	}
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Olá %s! Bem-vindo ao Personal Cockpit v%s", name, GetFullVersion())
}
