import { useState, useEffect } from 'react';
import {
  CreateTask,
  GetAllCategories,
  CreateCategory,
  GetAllProfiles,
  AddAssignee,
  GoogleCalendarConnected,
  SyncTaskToCalendar,
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

interface CreateTaskModalProps {
  boardID: number;
  columnID: number | null;
  profileID: number;
  onCreated: () => void;
  onClose: () => void;
}

export function CreateTaskModal({ boardID, columnID, profileID, onCreated, onClose }: CreateTaskModalProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState('medium');
  const [dueDate, setDueDate] = useState('');
  const [categoryId, setCategoryId] = useState<number | undefined>();
  const [categories, setCategories] = useState<models.Category[]>([]);
  const [profiles, setProfiles] = useState<models.ProfileSummary[]>([]);
  const [assigneeIDs, setAssigneeIDs] = useState<number[]>([profileID]); // current user pre-selected
  const [gcalConnected, setGcalConnected] = useState(false);
  const [syncCalendar, setSyncCalendar] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // New category inline
  const [showNewCat, setShowNewCat] = useState(false);
  const [newCatName, setNewCatName] = useState('');
  const [newCatColor, setNewCatColor] = useState('#3b82f6');
  const [savingCat, setSavingCat] = useState(false);

  useEffect(() => {
    Promise.all([
      GetAllCategories(),
      GetAllProfiles(),
      GoogleCalendarConnected(profileID),
    ]).then(([cats, profs, gcal]) => {
      setCategories((cats || []).filter((c: models.Category) => c.type === 'task' || c.type === 'general'));
      setProfiles(profs || []);
      setGcalConnected(gcal as boolean);
    });
    document.body.classList.add('modal-open');
    return () => document.body.classList.remove('modal-open');
  }, []);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  const toggleAssignee = (id: number) => {
    setAssigneeIDs(prev =>
      prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]
    );
  };

  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCatName.trim()) return;
    setSavingCat(true);
    try {
      const cat = new models.Category();
      cat.name = newCatName.trim();
      cat.color = newCatColor;
      cat.type = 'task';
      const newID = await CreateCategory(cat);
      const updatedCats = await GetAllCategories();
      const filtered = (updatedCats || []).filter((c: models.Category) => c.type === 'task' || c.type === 'general');
      setCategories(filtered);
      setCategoryId(newID);
      setNewCatName('');
      setShowNewCat(false);
    } finally {
      setSavingCat(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const task = new models.Task();
      task.title = title;
      task.description = description;
      task.status = 'pending';
      task.priority = priority;
      task.category_id = categoryId;
      task.board_id = boardID;
      task.column_id = columnID ?? undefined;
      task.due_date = dueDate ? new Date(dueDate) : undefined;
      const taskID = await CreateTask(task);

      // Add assignees
      await Promise.all(assigneeIDs.map(pid => AddAssignee(taskID, pid)));

      // Sync to Google Calendar if requested
      if (syncCalendar && gcalConnected) {
        try {
          await SyncTaskToCalendar(taskID, profileID);
        } catch (calErr: any) {
          console.warn('Tarefa criada mas falhou ao sincronizar com Calendar:', calErr);
        }
      }

      onCreated();
    } catch (err: any) {
      setError(err.message || 'Erro ao criar tarefa');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Nova Tarefa</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="modal-body">
          {error && <div className="error-message">{error}</div>}

          <div className="form-group">
            <label className="form-label">Título *</label>
            <input
              className="form-input"
              value={title}
              onChange={e => setTitle(e.target.value)}
              required
              autoFocus
              placeholder="O que precisa ser feito?"
            />
          </div>

          <div className="form-group">
            <label className="form-label">Descrição</label>
            <textarea
              className="form-textarea"
              value={description}
              onChange={e => setDescription(e.target.value)}
              rows={3}
              placeholder="Detalhes (opcional)"
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label className="form-label">Prioridade</label>
              <select className="form-select" value={priority} onChange={e => setPriority(e.target.value)}>
                <option value="low">🟢 Baixa</option>
                <option value="medium">🟡 Média</option>
                <option value="high">🔴 Alta</option>
              </select>
            </div>
            <div className="form-group">
              <label className="form-label">Vencimento</label>
              <input className="form-input" type="date" value={dueDate} onChange={e => setDueDate(e.target.value)} />
            </div>
          </div>

          {/* Category */}
          <div className="form-group">
            <div className="create-task-label-row">
              <label className="form-label">Categoria</label>
              <button
                type="button"
                className="btn-inline-link"
                onClick={() => setShowNewCat(v => !v)}
              >
                {showNewCat ? 'Cancelar' : '+ Nova categoria'}
              </button>
            </div>

            {showNewCat ? (
              <div className="new-category-inline">
                <input
                  className="form-input"
                  value={newCatName}
                  onChange={e => setNewCatName(e.target.value)}
                  placeholder="Nome da categoria"
                  autoFocus
                />
                <input
                  type="color"
                  value={newCatColor}
                  onChange={e => setNewCatColor(e.target.value)}
                  className="color-picker"
                  title="Cor"
                />
                <button
                  type="button"
                  className="btn btn-primary btn-sm"
                  onClick={handleCreateCategory}
                  disabled={savingCat || !newCatName.trim()}
                >
                  {savingCat ? '...' : 'Criar'}
                </button>
              </div>
            ) : (
              <select
                className="form-select"
                value={categoryId ?? ''}
                onChange={e => setCategoryId(e.target.value ? Number(e.target.value) : undefined)}
              >
                <option value="">Sem categoria</option>
                {categories.map(c => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>
            )}
          </div>

          {/* Assignees */}
          {profiles.length > 0 && (
            <div className="form-group">
              <label className="form-label">Responsáveis</label>
              <div className="assignee-chips">
                {profiles.map(p => {
                  const selected = assigneeIDs.includes(p.id);
                  return (
                    <button
                      key={p.id}
                      type="button"
                      className={`assignee-chip ${selected ? 'selected' : ''}`}
                      onClick={() => toggleAssignee(p.id)}
                      title={p.name}
                    >
                      <div className="profile-avatar small">
                        {p.avatar_url
                          ? <img src={p.avatar_url} alt={p.name} />
                          : <span>{p.name.charAt(0).toUpperCase()}</span>
                        }
                      </div>
                      <span>{p.name}</span>
                      {selected && <span className="chip-check">✓</span>}
                    </button>
                  );
                })}
              </div>
            </div>
          )}

          {/* Google Calendar sync */}
          {gcalConnected && (
            <div className="form-group">
              <div className="checkbox-group gcal-sync-check">
                <input type="checkbox" id="taskSyncCal" checked={syncCalendar} onChange={e => setSyncCalendar(e.target.checked)} />
                <label htmlFor="taskSyncCal">📅 Adicionar ao Google Calendar</label>
              </div>
              {syncCalendar && !dueDate && (
                <p style={{ fontSize: '0.8125rem', color: 'var(--warning)', marginTop: 4 }}>
                  ⚠ Defina uma data de vencimento para sincronizar com o Calendar
                </p>
              )}
            </div>
          )}

          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={onClose}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Criando...' : 'Criar Tarefa'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
