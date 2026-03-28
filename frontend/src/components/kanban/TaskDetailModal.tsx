import { useState, useEffect } from 'react';
import {
  GetTaskByID,
  GetTaskDetail,
  UpdateTask,
  UpdateProgress,
  CreateSubtask,
  ToggleSubtask,
  DeleteSubtask,
  AddAssignee,
  RemoveAssignee,
  AddDependency,
  RemoveDependency,
  AddComment,
  DeleteComment,
  GetActivity,
  GetAllProfiles,
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

interface TaskDetailModalProps {
  taskID: number;
  profileID: number;
  onClose: () => void;
  onUpdated: () => void;
}

export function TaskDetailModal({ taskID, profileID, onClose, onUpdated }: TaskDetailModalProps) {
  const [task, setTask] = useState<models.Task | null>(null);
  const [detail, setDetail] = useState<models.TaskDetail | null>(null);
  const [activity, setActivity] = useState<models.TaskActivity[]>([]);
  const [allProfiles, setAllProfiles] = useState<models.ProfileSummary[]>([]);
  const [tab, setTab] = useState<'subtasks' | 'assignees' | 'deps' | 'comments' | 'activity'>('subtasks');
  const [loading, setLoading] = useState(true);

  // edit fields
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [progress, setProgress] = useState(0);
  const [saving, setSaving] = useState(false);

  // new subtask
  const [newSubtask, setNewSubtask] = useState('');
  // new comment
  const [newComment, setNewComment] = useState('');

  useEffect(() => {
    loadAll();
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [taskID]);

  const loadAll = async () => {
    setLoading(true);
    try {
      const [t, d, act, profiles] = await Promise.all([
        GetTaskByID(taskID),
        GetTaskDetail(taskID),
        GetActivity(taskID),
        GetAllProfiles(),
      ]);
      setTask(t);
      setDetail(d);
      setActivity(act || []);
      setAllProfiles(profiles || []);
      setTitle(t?.title || '');
      setDescription(t?.description || '');
      setProgress(t?.progress || 0);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!task) return;
    setSaving(true);
    try {
      await UpdateTask({ ...task, title, description, progress } as models.Task);
      onUpdated();
    } finally {
      setSaving(false);
    }
  };

  const handleProgressChange = async (val: number) => {
    setProgress(val);
    await UpdateProgress(taskID, val);
    onUpdated();
  };

  // ── Subtasks ──
  const handleAddSubtask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSubtask.trim()) return;
    await CreateSubtask(taskID, newSubtask.trim());
    setNewSubtask('');
    await loadAll();
  };

  const handleToggleSubtask = async (id: number) => {
    await ToggleSubtask(id);
    await loadAll();
  };

  const handleDeleteSubtask = async (id: number) => {
    await DeleteSubtask(id);
    await loadAll();
  };

  // ── Assignees ──
  const isAssigned = (pid: number) => detail?.assignees?.some(a => a.id === pid) ?? false;

  const handleToggleAssignee = async (pid: number) => {
    if (isAssigned(pid)) {
      await RemoveAssignee(taskID, pid);
    } else {
      await AddAssignee(taskID, pid);
    }
    await loadAll();
  };

  // ── Dependencies ──
  const handleAddDep = async (depID: number) => {
    await AddDependency(taskID, depID);
    await loadAll();
  };

  const handleRemoveDep = async (depID: number) => {
    await RemoveDependency(taskID, depID);
    await loadAll();
  };

  // ── Comments ──
  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;
    await AddComment(taskID, profileID, newComment.trim());
    setNewComment('');
    await loadAll();
  };

  const subtasksDone = detail?.subtasks?.filter(s => s.completed).length ?? 0;
  const subtasksTotal = detail?.subtasks?.length ?? 0;

  if (loading || !task) {
    return (
      <div className="modal-overlay">
        <div className="modal-content modal-content--large">
          <p>Carregando...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content modal-content--large task-detail-modal" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Detalhes da Tarefa</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>

        <div className="task-detail-body">
          {/* Left: main info */}
          <div className="task-detail-main">
            <div className="form-group">
              <input
                className="form-input task-detail-title"
                value={title}
                onChange={e => setTitle(e.target.value)}
                placeholder="Título"
              />
            </div>
            <div className="form-group">
              <textarea
                className="form-textarea"
                value={description}
                onChange={e => setDescription(e.target.value)}
                rows={4}
                placeholder="Descrição..."
              />
            </div>

            {/* Progress */}
            <div className="form-group">
              <label className="form-label">Progresso: {progress}%</label>
              <input
                type="range" min={0} max={100} step={5}
                value={progress}
                onChange={e => handleProgressChange(Number(e.target.value))}
                className="progress-slider"
              />
            </div>

            {/* Subtask counter */}
            {subtasksTotal > 0 && (
              <div className="task-detail-subtask-summary">
                <div className="progress-bar-bg">
                  <div className="progress-bar-fill" style={{ width: `${Math.round(subtasksDone / subtasksTotal * 100)}%` }} />
                </div>
                <span>{subtasksDone}/{subtasksTotal} subtarefas</span>
              </div>
            )}

            <div className="task-detail-save">
              <button className="btn btn-primary" onClick={handleSave} disabled={saving}>
                {saving ? 'Salvando...' : 'Salvar'}
              </button>
            </div>
          </div>

          {/* Right: tabs */}
          <div className="task-detail-side">
            <div className="task-detail-tabs">
              {(['subtasks', 'assignees', 'deps', 'comments', 'activity'] as const).map(t => (
                <button key={t} className={`tab-btn ${tab === t ? 'active' : ''}`} onClick={() => setTab(t)}>
                  {{ subtasks: 'Subtarefas', assignees: 'Responsáveis', deps: 'Dependências', comments: 'Comentários', activity: 'Histórico' }[t]}
                </button>
              ))}
            </div>

            <div className="task-detail-tab-body">
              {/* Subtasks */}
              {tab === 'subtasks' && (
                <div>
                  <form onSubmit={handleAddSubtask} className="inline-add-form">
                    <input className="form-input" value={newSubtask} onChange={e => setNewSubtask(e.target.value)} placeholder="Nova subtarefa..." />
                    <button type="submit" className="btn btn-primary btn-sm">+</button>
                  </form>
                  <ul className="subtask-list">
                    {detail?.subtasks?.map(st => (
                      <li key={st.id} className={`subtask-item ${st.completed ? 'done' : ''}`}>
                        <input type="checkbox" checked={st.completed} onChange={() => handleToggleSubtask(st.id)} />
                        <span>{st.title}</span>
                        <button className="btn-icon" onClick={() => handleDeleteSubtask(st.id)}>×</button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Assignees */}
              {tab === 'assignees' && (
                <div className="assignee-list">
                  {allProfiles.map(p => (
                    <div key={p.id} className={`assignee-item ${isAssigned(p.id) ? 'assigned' : ''}`} onClick={() => handleToggleAssignee(p.id)}>
                      <div className="profile-avatar small">
                        {p.avatar_url ? <img src={p.avatar_url} alt={p.name} /> : <span>{p.name.charAt(0)}</span>}
                      </div>
                      <span>{p.name}</span>
                      {isAssigned(p.id) && <span className="check">✓</span>}
                    </div>
                  ))}
                </div>
              )}

              {/* Dependencies */}
              {tab === 'deps' && (
                <div>
                  <p className="tab-section-label">Esta tarefa depende de:</p>
                  {detail?.dependencies?.map(depID => {
                    return (
                      <div key={depID} className="dep-item">
                        <span>#{depID}</span>
                        <button className="btn-icon" onClick={() => handleRemoveDep(depID)}>×</button>
                      </div>
                    );
                  })}
                  <select className="form-select" defaultValue="" onChange={e => { if (e.target.value) handleAddDep(Number(e.target.value)); e.target.value = ''; }}>
                    <option value="" disabled>Adicionar dependência...</option>
                    {/* In a real scenario we'd load other tasks here */}
                  </select>
                  {detail?.dependents && detail.dependents.length > 0 && (
                    <>
                      <p className="tab-section-label">Tasks que dependem desta:</p>
                      {detail.dependents.map(id => <div key={id} className="dep-item dep-item--dependent">#{id}</div>)}
                    </>
                  )}
                </div>
              )}

              {/* Comments */}
              {tab === 'comments' && (
                <div>
                  <div className="comment-list">
                    {detail?.comments?.map(c => (
                      <div key={c.id} className="comment-item">
                        <div className="comment-header">
                          <strong>{c.profile?.name || 'Sistema'}</strong>
                          <span className="comment-date">{new Date(c.created_at).toLocaleString('pt-BR')}</span>
                          <button className="btn-icon" onClick={() => { DeleteComment(c.id).then(loadAll); }}>×</button>
                        </div>
                        <p>{c.content}</p>
                      </div>
                    ))}
                  </div>
                  <form onSubmit={handleAddComment} className="comment-form">
                    <textarea className="form-textarea" value={newComment} onChange={e => setNewComment(e.target.value)} rows={2} placeholder="Escrever comentário..." />
                    <button type="submit" className="btn btn-primary btn-sm">Enviar</button>
                  </form>
                </div>
              )}

              {/* Activity */}
              {tab === 'activity' && (
                <ul className="activity-list">
                  {activity.map(a => (
                    <li key={a.id} className="activity-item">
                      <span className="activity-profile">{a.profile?.name || 'Sistema'}</span>
                      <span className="activity-action">{a.action}</span>
                      {a.from_value && <span className="activity-from">{a.from_value}</span>}
                      {a.to_value && <><span>→</span><span className="activity-to">{a.to_value}</span></>}
                      <span className="activity-date">{new Date(a.created_at).toLocaleString('pt-BR')}</span>
                    </li>
                  ))}
                  {activity.length === 0 && <li className="empty-state-small"><p>Sem histórico</p></li>}
                </ul>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
