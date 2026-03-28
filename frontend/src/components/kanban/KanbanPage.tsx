import { useState, useEffect } from 'react';
import {
  DndContext,
  DragEndEvent,
  DragOverEvent,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
  closestCorners,
  DragOverlay,
} from '@dnd-kit/core';
import {
  GetBoardsByProfile,
  GetProfileByID,
  CreateBoard,
  GetTasksByBoard,
  MoveTask,
  DeleteTask,
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';
import { KanbanColumn } from './KanbanColumn';
import { TaskCard } from './TaskCard';
import { TaskDetailModal } from './TaskDetailModal';
import { CreateTaskModal } from './CreateTaskModal';

interface KanbanPageProps {
  profileID: number;
}

function parseTaskID(dndID: string | number): number {
  return Number(String(dndID).replace('task-', ''));
}

function parseColumnID(dndID: string | number): number {
  return Number(String(dndID).replace('column-', ''));
}

export function KanbanPage({ profileID }: KanbanPageProps) {
  const [boards, setBoards] = useState<models.BoardWithColumns[]>([]);
  const [activeBoardID, setActiveBoardID] = useState<number | null>(null);
  const [tasks, setTasks] = useState<models.Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState<models.ProfileSummary | null>(null);

  const [draggingTask, setDraggingTask] = useState<models.Task | null>(null);
  const [selectedTaskID, setSelectedTaskID] = useState<number | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createColumnID, setCreateColumnID] = useState<number | null>(null);

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

  useEffect(() => {
    loadData();
  }, [profileID]);

  const loadData = async () => {
    setLoading(true);
    try {
      const p = await GetProfileByID(profileID);
      setProfile(p);
      const boardList = await GetBoardsByProfile(profileID, p?.role || 'member');
      setBoards(boardList || []);

      if (boardList && boardList.length > 0) {
        const bid = boardList[0].id;
        setActiveBoardID(bid);
        await loadTasks(bid);
      } else {
        const newBoard = await CreateBoard('Meu Board', profileID, false);
        setBoards([newBoard!]);
        setActiveBoardID(newBoard!.id);
        setTasks([]);
      }
    } catch (err) {
      console.error('Erro ao carregar boards:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadTasks = async (boardID: number) => {
    const result = await GetTasksByBoard(boardID);
    setTasks(result || []);
  };

  const switchBoard = async (boardID: number) => {
    setActiveBoardID(boardID);
    await loadTasks(boardID);
  };

  // ── Drag & Drop ──────────────────────────────────────────

  const onDragStart = (event: DragStartEvent) => {
    const id = parseTaskID(event.active.id);
    setDraggingTask(tasks.find(t => t.id === id) || null);
  };

  // Move task to new column visually as the user drags
  const onDragOver = (event: DragOverEvent) => {
    const { active, over } = event;
    if (!over) return;

    const activeStr = String(active.id);
    const overStr = String(over.id);
    if (!activeStr.startsWith('task-')) return;

    const taskID = parseTaskID(activeStr);

    let targetColumnID: number;
    if (overStr.startsWith('column-')) {
      targetColumnID = parseColumnID(overStr);
    } else if (overStr.startsWith('task-')) {
      const overTask = tasks.find(t => t.id === parseTaskID(overStr));
      if (!overTask?.column_id) return;
      targetColumnID = overTask.column_id;
    } else {
      return;
    }

    setTasks(prev => {
      const task = prev.find(t => t.id === taskID);
      if (!task || task.column_id === targetColumnID) return prev;
      return prev.map(t =>
        t.id === taskID
          ? ({ ...t, column_id: targetColumnID } as models.Task)
          : t
      );
    });
  };

  const onDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    const task = draggingTask;
    setDraggingTask(null);

    if (!over || !task) {
      // User dropped outside — revert
      if (activeBoardID) await loadTasks(activeBoardID);
      return;
    }

    // Find the column the task ended up in (already updated by onDragOver)
    const movedTask = tasks.find(t => t.id === task.id);
    if (!movedTask?.column_id) return;

    // No column change? Nothing to persist
    if (movedTask.column_id === task.column_id) return;

    const position = tasks.filter(t => t.column_id === movedTask.column_id).length - 1;

    try {
      await MoveTask(task.id, movedTask.column_id, Math.max(0, position));
    } catch (err) {
      console.error('Erro ao mover tarefa:', err);
      if (activeBoardID) await loadTasks(activeBoardID);
    }
  };

  // ── Task actions ─────────────────────────────────────────

  const handleDeleteTask = async (taskID: number) => {
    if (!window.confirm('Deseja deletar esta tarefa?')) return;
    try {
      await DeleteTask(taskID);
      setTasks(prev => prev.filter(t => t.id !== taskID));
    } catch {
      alert('Erro ao deletar tarefa');
    }
  };

  const handleTaskCreated = async () => {
    if (activeBoardID) await loadTasks(activeBoardID);
    setShowCreateModal(false);
    setCreateColumnID(null);
  };

  const activeBoard = boards.find(b => b.id === activeBoardID);

  if (loading) {
    return <div className="page"><p>Carregando board...</p></div>;
  }

  return (
    <div className="page kanban-page">
      <div className="page-header">
        <div className="kanban-header">
          <div>
            <h1 className="page-title">Kanban</h1>
            <p className="page-subtitle">
              {profile?.role === 'admin' ? 'Visão de todos os boards' : 'Seus boards'}
            </p>
          </div>

          <div className="board-selector">
            {boards.map(b => (
              <button
                key={b.id}
                className={`board-tab ${b.id === activeBoardID ? 'active' : ''}`}
                onClick={() => switchBoard(b.id)}
              >
                {b.name}
                {b.is_shared && <span className="board-shared-badge">Família</span>}
              </button>
            ))}
            <button className="board-tab board-tab-add" onClick={async () => {
              const name = prompt('Nome do novo board:');
              if (!name) return;
              const nb = await CreateBoard(name, profileID, false);
              if (nb) { setBoards(prev => [...prev, nb]); switchBoard(nb.id); }
            }}>+ Board</button>
          </div>
        </div>
      </div>

      {activeBoard && (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={onDragStart}
          onDragOver={onDragOver}
          onDragEnd={onDragEnd}
        >
          <div className="kanban-board">
            {activeBoard.columns?.map(col => {
              const colTasks = tasks
                .filter(t => t.column_id === col.id)
                .sort((a, b) => a.position - b.position);
              const atLimit = col.wip_limit != null && colTasks.length >= col.wip_limit!;

              return (
                <KanbanColumn
                  key={col.id}
                  column={col}
                  tasks={colTasks}
                  atLimit={atLimit}
                  onAddTask={() => { setCreateColumnID(col.id); setShowCreateModal(true); }}
                  onTaskClick={(id) => setSelectedTaskID(id)}
                  onDeleteTask={handleDeleteTask}
                />
              );
            })}
          </div>

          <DragOverlay>
            {draggingTask && <TaskCard task={draggingTask} isDragging />}
          </DragOverlay>
        </DndContext>
      )}

      {showCreateModal && activeBoardID && (
        <CreateTaskModal
          boardID={activeBoardID}
          columnID={createColumnID}
          profileID={profileID}
          onCreated={handleTaskCreated}
          onClose={() => { setShowCreateModal(false); setCreateColumnID(null); }}
        />
      )}

      {selectedTaskID && (
        <TaskDetailModal
          taskID={selectedTaskID}
          profileID={profileID}
          onClose={() => setSelectedTaskID(null)}
          onUpdated={() => activeBoardID && loadTasks(activeBoardID)}
        />
      )}
    </div>
  );
}
