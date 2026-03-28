package services

import (
	"testing"

	"personal-cockpit/database"
	"personal-cockpit/models"
)

func setupNoteService(t *testing.T) *NoteService {
	t.Helper()
	db, err := database.NewInMemoryDB()
	if err != nil {
		t.Fatalf("falha ao criar banco de teste: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewNoteService(db)
}

func TestCreateNote(t *testing.T) {
	svc := setupNoteService(t)

	id, err := svc.CreateNote(models.Note{
		Title:   "Minha primeira nota",
		Content: "Conteúdo da nota",
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if id <= 0 {
		t.Errorf("esperado ID > 0, got %d", id)
	}
}

func TestCreateNoteRequiresTitle(t *testing.T) {
	svc := setupNoteService(t)

	_, err := svc.CreateNote(models.Note{Content: "sem título"})
	if err == nil {
		t.Error("esperava erro de validação para nota sem título")
	}
}

func TestGetAllNotes(t *testing.T) {
	svc := setupNoteService(t)

	for _, title := range []string{"Nota 1", "Nota 2", "Nota 3"} {
		svc.CreateNote(models.Note{Title: title})
	}

	notes, err := svc.GetAllNotes()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(notes) != 3 {
		t.Errorf("esperado 3 notas, got %d", len(notes))
	}
}

func TestGetNoteByID(t *testing.T) {
	svc := setupNoteService(t)

	id, _ := svc.CreateNote(models.Note{Title: "Nota específica", Content: "conteúdo"})

	note, err := svc.GetNoteByID(int(id))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if note.Title != "Nota específica" {
		t.Errorf("esperado 'Nota específica', got '%s'", note.Title)
	}
}

func TestGetNoteByIDNotFound(t *testing.T) {
	svc := setupNoteService(t)

	_, err := svc.GetNoteByID(9999)
	if err == nil {
		t.Error("esperava erro para ID inexistente")
	}
}

func TestUpdateNote(t *testing.T) {
	svc := setupNoteService(t)

	id, _ := svc.CreateNote(models.Note{Title: "Original", Content: "conteúdo original"})

	err := svc.UpdateNote(models.Note{
		ID:      int(id),
		Title:   "Atualizada",
		Content: "conteúdo atualizado",
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	note, _ := svc.GetNoteByID(int(id))
	if note.Title != "Atualizada" {
		t.Errorf("esperado 'Atualizada', got '%s'", note.Title)
	}
}

func TestDeleteNote(t *testing.T) {
	svc := setupNoteService(t)

	id, _ := svc.CreateNote(models.Note{Title: "Para deletar"})

	if err := svc.DeleteNote(int(id)); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	_, err := svc.GetNoteByID(int(id))
	if err == nil {
		t.Error("nota deveria ter sido deletada")
	}
}

func TestToggleFavorite(t *testing.T) {
	svc := setupNoteService(t)

	id, _ := svc.CreateNote(models.Note{Title: "Nota favorita"})

	svc.ToggleFavorite(int(id))
	note, _ := svc.GetNoteByID(int(id))
	if !note.IsFavorite {
		t.Error("nota deveria ser favorita após primeiro toggle")
	}

	svc.ToggleFavorite(int(id))
	note, _ = svc.GetNoteByID(int(id))
	if note.IsFavorite {
		t.Error("nota não deveria ser favorita após segundo toggle")
	}
}

func TestGetFavoriteNotes(t *testing.T) {
	svc := setupNoteService(t)

	id1, _ := svc.CreateNote(models.Note{Title: "Fav 1"})
	id2, _ := svc.CreateNote(models.Note{Title: "Fav 2"})
	svc.CreateNote(models.Note{Title: "Normal"})

	svc.ToggleFavorite(int(id1))
	svc.ToggleFavorite(int(id2))

	favorites, err := svc.GetFavoriteNotes()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(favorites) != 2 {
		t.Errorf("esperado 2 favoritas, got %d", len(favorites))
	}
}

func TestSearchNotes(t *testing.T) {
	svc := setupNoteService(t)

	svc.CreateNote(models.Note{Title: "Receita de bolo", Content: "farinha, ovos, açúcar"})
	svc.CreateNote(models.Note{Title: "Lista de compras", Content: "leite, pão, bolo"})
	svc.CreateNote(models.Note{Title: "Reunião", Content: "pauta da semana"})

	results, err := svc.SearchNotes("bolo")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("esperado 2 resultados para 'bolo', got %d", len(results))
	}
}

func TestSearchNotesNoResults(t *testing.T) {
	svc := setupNoteService(t)

	svc.CreateNote(models.Note{Title: "Nota qualquer"})

	results, err := svc.SearchNotes("xyzzy")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("esperado 0 resultados, got %d", len(results))
	}
}
