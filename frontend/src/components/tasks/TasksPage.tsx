import React, { useState, useEffect } from 'react';
import { GetAllTasks } from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';
import { TaskForm } from './TaskForm';
import { TaskList } from './TaskList';

export function TasksPage() {
  const [tasks, setTasks] = useState<models.Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);

  const loadTasks = async () => {
    try {
      const result = await GetAllTasks();
      setTasks(result || []);
    } catch (err) {
      console.error('Erro ao carregar tarefas:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTasks();
  }, []);

  if (loading) {
    return (
      <div className="page">
        <p>Carregando...</p>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="page-header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1 className="page-title">Tarefas</h1>
            <p className="page-subtitle">Gerencie suas tarefas diárias</p>
          </div>
          <button 
            className="btn btn-primary"
            onClick={() => setShowModal(true)}
          >
            + Nova Tarefa
          </button>
        </div>
      </div>

      <div className="page-content">
        <TaskList tasks={tasks} onTaskUpdated={loadTasks} />
      </div>

      {showModal && (
        <TaskForm 
          onTaskCreated={loadTasks} 
          onClose={() => setShowModal(false)}
        />
      )}
    </div>
  );
}