import React, { useState, useEffect } from 'react';
import { 
  GetAllNotes, 
  CreateNote, 
  DeleteNote, 
  ToggleNoteFavorite 
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

export function NotesPage() {
  const [notes, setNotes] = useState<models.Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  
  // Form state
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [formError, setFormError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    loadNotes();
  }, []);

  const loadNotes = async () => {
    try {
      const result = await GetAllNotes();
      setNotes(result || []);
    } catch (err) {
      console.error('Erro ao carregar notas:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');
    setFormLoading(true);
    
    try {
      const note = new models.Note();
      note.title = title;
      note.content = content;
      
      await CreateNote(note);
      
      setTitle('');
      setContent('');
      setShowModal(false);
      loadNotes();
    } catch (err: any) {
      setFormError(err.message || 'Erro ao criar nota');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Deseja realmente deletar esta nota?')) {
      return;
    }
    
    try {
      await DeleteNote(id);
      loadNotes();
    } catch (err) {
      alert('Erro ao deletar nota');
    }
  };

  const handleToggleFavorite = async (id: number) => {
    try {
      await ToggleNoteFavorite(id);
      loadNotes();
    } catch (err) {
      alert('Erro ao favoritar nota');
    }
  };

  const handleCloseModal = () => {
    setTitle('');
    setContent('');
    setFormError('');
    setShowModal(false);
  };

  if (loading) {
    return <div className="page">Carregando...</div>;
  }

  return (
    <div className="page">
      <div className="page-header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1 className="page-title">Notas</h1>
            <p className="page-subtitle">Suas anotações e ideias</p>
          </div>
          <button 
            className="btn btn-primary"
            onClick={() => setShowModal(true)}
          >
            + Nova Nota
          </button>
        </div>
      </div>

      <div className="page-content">
        {notes.length === 0 ? (
          <div className="empty-state">
            <p>📝</p>
            <p>Nenhuma nota ainda. Crie sua primeira nota!</p>
          </div>
        ) : (
          <div className="notes-grid">
            {notes.map((note) => (
              <div key={note.id} className={`note-card ${note.is_favorite ? 'favorite' : ''}`}>
                <div className="note-card-header">
                  <h3 className="note-card-title">{note.title}</h3>
                </div>
                
                {note.content && (
                  <p className="note-card-content">{note.content}</p>
                )}
                
                <div className="note-card-footer">
                  <span className="note-card-date">
                    {new Date(note.created_at).toLocaleDateString('pt-BR')}
                  </span>
                  <div className="note-card-actions">
                    <button onClick={() => handleToggleFavorite(note.id)}>
                      {note.is_favorite ? '★' : '☆'}
                    </button>
                    <button onClick={() => handleDelete(note.id)}>
                      ×
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">Nova Nota</h2>
              <button 
                type="button" 
                className="modal-close"
                onClick={handleCloseModal}
              >
                ×
              </button>
            </div>

            <form onSubmit={handleSubmit} className="modal-body">
              {formError && <div className="error-message">{formError}</div>}

              <div className="form-group">
                <label className="form-label">Título *</label>
                <input
                  type="text"
                  className="form-input"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Título da nota"
                  required
                  autoFocus
                />
              </div>
              
              <div className="form-group">
                <label className="form-label">Conteúdo</label>
                <textarea
                  className="form-textarea"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  placeholder="Escreva sua nota..."
                  rows={8}
                />
              </div>
              
              <div className="modal-footer">
                <button 
                  type="button" 
                  className="btn btn-secondary"
                  onClick={handleCloseModal}
                  disabled={formLoading}
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={formLoading}
                >
                  {formLoading ? 'Salvando...' : 'Salvar Nota'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}