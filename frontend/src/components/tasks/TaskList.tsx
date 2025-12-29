import React from 'react';
import { DeleteTask, ToggleTaskStatus } from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

interface TaskListProps {
  tasks: models.Task[];
  onTaskUpdated: () => void;
}

export function TaskList({ tasks, onTaskUpdated }: TaskListProps) {
  const handleToggle = async (id: number) => {
    try {
      await ToggleTaskStatus(id);
      onTaskUpdated();
    } catch (err) {
      alert('Erro ao atualizar tarefa');
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Deseja realmente deletar esta tarefa?')) {
      return;
    }

    try {
      await DeleteTask(id);
      onTaskUpdated();
    } catch (err) {
      alert('Erro ao deletar tarefa');
    }
  };

  const getPriorityEmoji = (priority: string) => {
    switch (priority) {
      case 'high':
        return '🔴';
      case 'medium':
        return '🟡';
      case 'low':
        return '🟢';
      default:
        return '⚪';
    }
  };

  if (tasks.length === 0) {
    return (
      <div className="empty-state">
        <p>📭 Nenhuma tarefa ainda</p>
        <p>Adicione uma nova tarefa acima!</p>
      </div>
    );
  }

  return (
    <div className="task-list">
      <h2>📋 Minhas Tarefas ({tasks.length})</h2>

      {tasks.map((task) => (
        <div
          key={task.id}
          className={`task-item ${task.status === 'completed' ? 'completed' : ''}`}
        >
          <div className="task-header">
            <input
              type="checkbox"
              checked={task.status === 'completed'}
              onChange={() => handleToggle(task.id)}
            />
            <span className="task-title">{task.title}</span>
            <span className="task-priority">{getPriorityEmoji(task.priority)}</span>
          </div>

          {task.description && (
            <div className="task-description">{task.description}</div>
          )}

          <div className="task-footer">
            <span className="task-date">
              {new Date(task.created_at).toLocaleDateString('pt-BR')}
            </span>
            <button
              className="btn-delete"
              onClick={() => handleDelete(task.id)}
            >
              🗑️ Deletar
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}