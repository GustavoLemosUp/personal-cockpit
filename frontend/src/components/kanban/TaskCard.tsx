import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { models } from '../../../wailsjs/go/models';

interface TaskCardProps {
  task: models.Task;
  onClick?: () => void;
  onDelete?: () => void;
  isDragging?: boolean;
}

const PRIORITY_COLOR: Record<string, string> = {
  high:   '#ef4444',
  medium: '#f59e0b',
  low:    '#10b981',
};

function dueDateClass(task: models.Task): string {
  if (!task.due_date) return '';
  const due = new Date(task.due_date as any);
  const now = new Date();
  const diff = (due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24);
  if (diff < 0) return 'task-card--overdue';
  if (diff <= 2) return 'task-card--due-soon';
  return '';
}

export function TaskCard({ task, onClick, onDelete, isDragging }: TaskCardProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging: sorting } = useSortable({
    id: `task-${task.id}`,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: sorting ? 0.4 : 1,
  };

  const ddClass = dueDateClass(task);

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`task-card ${ddClass} ${isDragging ? 'task-card--dragging' : ''}`}
      onClick={onClick}
      {...attributes}
      {...listeners}
    >
      <div className="task-card-header">
        <div
          className="task-card-priority"
          style={{ background: PRIORITY_COLOR[task.priority] || '#6b7280' }}
        />
        <span className="task-card-title">{task.title}</span>
        <button
          className="task-card-delete"
          onClick={e => { e.stopPropagation(); onDelete?.(); }}
          title="Deletar"
        >×</button>
      </div>

      {task.description && (
        <p className="task-card-desc">{task.description}</p>
      )}

      <div className="task-card-footer">
        {task.due_date && (
          <span className="task-card-date">
            {new Date(task.due_date as any).toLocaleDateString('pt-BR')}
          </span>
        )}
        {task.progress != null && task.progress > 0 && (
          <div className="task-card-progress">
            <div className="task-card-progress-bar" style={{ width: `${task.progress}%` }} />
          </div>
        )}
        {task.google_event_id && (
          <span className="task-card-gcal" title="No Google Calendar">📅</span>
        )}
      </div>
    </div>
  );
}
