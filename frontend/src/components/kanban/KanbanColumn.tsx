import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { models } from '../../../wailsjs/go/models';
import { TaskCard } from './TaskCard';

interface KanbanColumnProps {
  column: models.KanbanColumn;
  tasks: models.Task[];
  atLimit: boolean;
  onAddTask: () => void;
  onTaskClick: (id: number) => void;
  onDeleteTask: (id: number) => void;
}

export function KanbanColumn({ column, tasks, atLimit, onAddTask, onTaskClick, onDeleteTask }: KanbanColumnProps) {
  const { setNodeRef, isOver } = useDroppable({ id: `column-${column.id}` });

  return (
    <div
      className={`kanban-column ${isOver ? 'kanban-column--over' : ''} ${atLimit ? 'kanban-column--limit' : ''}`}
    >
      <div className="kanban-column-header" style={{ borderTopColor: column.color }}>
        <div className="kanban-column-title">
          <span className="kanban-column-dot" style={{ background: column.color }} />
          <h3>{column.name}</h3>
          <span className="kanban-column-count">{tasks.length}</span>
          {column.wip_limit != null && (
            <span className={`kanban-wip ${atLimit ? 'kanban-wip--full' : ''}`}>
              / {column.wip_limit}
            </span>
          )}
        </div>
        <button className="kanban-add-btn" onClick={onAddTask} disabled={atLimit} title="Adicionar tarefa">
          +
        </button>
      </div>

      <div ref={setNodeRef} className="kanban-column-body">
        <SortableContext
          items={tasks.map(t => `task-${t.id}`)}
          strategy={verticalListSortingStrategy}
        >
          {tasks.map(task => (
            <TaskCard
              key={task.id}
              task={task}
              onClick={() => onTaskClick(task.id)}
              onDelete={() => onDeleteTask(task.id)}
            />
          ))}
        </SortableContext>

        {tasks.length === 0 && (
          <div className="kanban-column-empty">
            <p>Arraste uma tarefa aqui</p>
          </div>
        )}
      </div>
    </div>
  );
}
