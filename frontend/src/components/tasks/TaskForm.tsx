import React, { useState, useEffect } from 'react';
import { CreateTask, GetAllCategories } from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

interface TaskFormProps {
  onTaskCreated: () => void;
  onClose: () => void;
}

export function TaskForm({ onTaskCreated, onClose }: TaskFormProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState('medium');
  const [categoryId, setCategoryId] = useState<number | undefined>(undefined);
  const [categories, setCategories] = useState<models.Category[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Carregar categorias
  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      const result = await GetAllCategories();
      // Filtrar apenas categorias de tarefas ou gerais
      const taskCategories = (result || []).filter(
        (c: models.Category) => c.type === 'task' || c.type === 'general'
      );
      setCategories(taskCategories);
    } catch (err) {
      console.error('Erro ao carregar categorias:', err);
    }
  };

  // Bloqueia scroll do body
  useEffect(() => {
    document.body.classList.add('modal-open');
    return () => {
      document.body.classList.remove('modal-open');
    };
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const task = new models.Task();
      task.title = title;
      task.description = description;
      task.status = 'pending';
      task.priority = priority;
      task.category_id = categoryId;

      await CreateTask(task);

      setTitle('');
      setDescription('');
      setPriority('medium');
      setCategoryId(undefined);

      onTaskCreated();
      onClose();
    } catch (err: any) {
      setError(err.message || 'Erro ao criar tarefa');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setTitle('');
    setDescription('');
    setPriority('medium');
    setCategoryId(undefined);
    setError('');
    onClose();
  };

  // Fechar com ESC
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        handleCancel();
      }
    };

    window.addEventListener('keydown', handleEscape);
    return () => window.removeEventListener('keydown', handleEscape);
  }, []);

  return (
    <div className="modal-overlay" onClick={handleCancel}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Nova Tarefa</h2>
          <button 
            type="button" 
            className="modal-close"
            onClick={handleCancel}
          >
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="modal-body">
          {error && <div className="error-message">{error}</div>}

          <div className="form-group">
            <label className="form-label">Título *</label>
            <input
              type="text"
              className="form-input"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Digite o título da tarefa"
              required
              autoFocus
            />
          </div>

          <div className="form-group">
            <label className="form-label">Descrição</label>
            <textarea
              className="form-textarea"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Descreva a tarefa (opcional)"
              rows={4}
            />
          </div>

          <div className="form-group">
            <label className="form-label">Prioridade</label>
            <select 
              className="form-select"
              value={priority} 
              onChange={(e) => setPriority(e.target.value)}
            >
              <option value="low">Baixa</option>
              <option value="medium">Média</option>
              <option value="high">Alta</option>
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Categoria</label>
            <select 
              className="form-select"
              value={categoryId || ''} 
              onChange={(e) => setCategoryId(e.target.value ? Number(e.target.value) : undefined)}
            >
              <option value="">Sem categoria</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>

          <div className="modal-footer">
            <button 
              type="button" 
              className="btn btn-secondary"
              onClick={handleCancel}
              disabled={loading}
            >
              Cancelar
            </button>
            <button 
              type="submit" 
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? 'Criando...' : 'Criar Tarefa'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}