import React, { useEffect, useState } from 'react';
import { GetAllTasks, GetAllNotes, GetTodayEvents } from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';

export function Dashboard() {
  const [stats, setStats] = useState({
    totalTasks: 0,
    pendingTasks: 0,
    completedTasks: 0,
    totalNotes: 0,
    todayEvents: 0,
  });
  const [todayEvents, setTodayEvents] = useState<models.Event[]>([]);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const tasks = await GetAllTasks() || [];
      const notes = await GetAllNotes() || [];
      const events = await GetTodayEvents() || [];

      setStats({
        totalTasks: tasks.length,
        pendingTasks: tasks.filter((t: models.Task) => t.status === 'pending').length,
        completedTasks: tasks.filter((t: models.Task) => t.status === 'completed').length,
        totalNotes: notes.length,
        todayEvents: events.length,
      });

      setTodayEvents(events);
    } catch (err) {
      console.error('Erro ao carregar estatísticas:', err);
    }
  };

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">Dashboard</h1>
        <p className="page-subtitle">Visão geral do seu dia</p>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-icon">✓</div>
          <div className="stat-content">
            <p className="stat-label">Tarefas Pendentes</p>
            <p className="stat-value">{stats.pendingTasks}</p>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon">✓</div>
          <div className="stat-content">
            <p className="stat-label">Tarefas Concluídas</p>
            <p className="stat-value">{stats.completedTasks}</p>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon">◐</div>
          <div className="stat-content">
            <p className="stat-label">Notas</p>
            <p className="stat-value">{stats.totalNotes}</p>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon">◷</div>
          <div className="stat-content">
            <p className="stat-label">Eventos Hoje</p>
            <p className="stat-value">{stats.todayEvents}</p>
          </div>
        </div>
      </div>

      {/* Lista de Eventos de Hoje */}
      <div className="dashboard-section">
        <h2 className="section-title">Eventos de Hoje</h2>
        {todayEvents.length === 0 ? (
          <div className="empty-state-small">
            <p>Nenhum evento agendado para hoje</p>
          </div>
        ) : (
          <div className="today-events-list">
            {todayEvents.map((event) => (
              <div key={event.id} className="today-event-item">
                <div className="event-time">
                  {new Date(event.start_date).toLocaleTimeString('pt-BR', { 
                    hour: '2-digit', 
                    minute: '2-digit' 
                  })}
                </div>
                <div className="event-details">
                  <h4 className="event-name">{event.title}</h4>
                  {event.location && (
                    <p className="event-location">📍 {event.location}</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="info-section">
        <h2 className="section-title">Bem-vindo ao Personal Cockpit</h2>
        <p className="section-text">
          Organize suas tarefas, notas e eventos em um só lugar.
        </p>
      </div>
    </div>
  );
}