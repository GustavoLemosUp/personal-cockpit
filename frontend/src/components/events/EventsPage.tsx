import React, { useState, useEffect } from 'react';
import {
  GetAllEvents,
  CreateEvent,
  DeleteEvent,
  GoogleCalendarConnected,
  SyncEventToCalendar,
  DeleteFromCalendar,
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

interface EventsPageProps {
  profileID: number;
}

export function EventsPage({ profileID }: EventsPageProps) {
  const [events, setEvents] = useState<models.Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [gcalConnected, setGcalConnected] = useState(false);

  // Form state
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [allDay, setAllDay] = useState(false);
  const [location, setLocation] = useState('');
  const [syncCalendar, setSyncCalendar] = useState(false);
  const [formError, setFormError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    loadEvents();
    GoogleCalendarConnected(profileID).then(setGcalConnected);
  }, [profileID]);

  const loadEvents = async () => {
    try {
      const result = await GetAllEvents();
      setEvents(result || []);
    } catch (err) {
      console.error('Erro ao carregar eventos:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');
    setFormLoading(true);

    try {
      const event = new models.Event();
      event.title = title;
      event.description = description;
      event.start_date = new Date(startDate).toISOString();
      event.end_date = new Date(endDate).toISOString();
      event.all_day = allDay;
      event.location = location;
      event.color = '#0071e3';

      const eventID = await CreateEvent(event);

      if (syncCalendar && gcalConnected) {
        try {
          await SyncEventToCalendar(eventID, profileID);
        } catch (calErr: any) {
          console.warn('Evento criado mas falhou ao sincronizar com Calendar:', calErr);
        }
      }

      setTitle(''); setDescription(''); setStartDate(''); setEndDate('');
      setAllDay(false); setLocation(''); setSyncCalendar(false);
      setShowModal(false);
      loadEvents();
    } catch (err: any) {
      setFormError(err.message || 'Erro ao criar evento');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (event: models.Event) => {
    if (!window.confirm('Deseja realmente deletar este evento?')) return;
    try {
      if (event.google_event_id) {
        try { await DeleteFromCalendar(profileID, event.google_event_id); } catch {}
      }
      await DeleteEvent(event.id);
      loadEvents();
    } catch {
      alert('Erro ao deletar evento');
    }
  };

  const handleSyncExisting = async (event: models.Event) => {
    if (event.google_event_id) {
      alert('Este evento já está sincronizado com o Google Calendar.');
      return;
    }
    try {
      await SyncEventToCalendar(event.id, profileID);
      loadEvents();
    } catch (err: any) {
      alert(err.message || 'Erro ao sincronizar com o Calendar');
    }
  };

  const handleCloseModal = () => {
    setTitle(''); setDescription(''); setStartDate(''); setEndDate('');
    setAllDay(false); setLocation(''); setFormError(''); setSyncCalendar(false);
    setShowModal(false);
  };

  if (loading) return <div className="page">Carregando...</div>;

  return (
    <div className="page">
      <div className="page-header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1 className="page-title">Eventos</h1>
            <p className="page-subtitle">
              Gerencie sua agenda
              {gcalConnected && <span style={{ marginLeft: 8, color: 'var(--success)', fontSize: '0.8125rem' }}>● Google Calendar conectado</span>}
            </p>
          </div>
          <button className="btn btn-primary" onClick={() => setShowModal(true)}>
            + Novo Evento
          </button>
        </div>
      </div>

      <div className="page-content">
        {events.length === 0 ? (
          <div className="empty-state">
            <p>📅</p>
            <p>Nenhum evento agendado</p>
          </div>
        ) : (
          <div className="events-container">
            {events.map(event => (
              <div key={event.id} className="event-item">
                <div className="event-header">
                  <h3 className="event-title">{event.title}</h3>
                  <div style={{ display: 'flex', gap: 8 }}>
                    {gcalConnected && (
                      <button
                        className={`btn btn-sm ${event.google_event_id ? 'btn-secondary' : 'btn-primary'}`}
                        onClick={() => handleSyncExisting(event)}
                        title={event.google_event_id ? 'Já sincronizado' : 'Sincronizar com Google Calendar'}
                      >
                        {event.google_event_id ? '📅 Sincronizado' : '📅 Sincronizar'}
                      </button>
                    )}
                    <button className="btn btn-danger btn-sm" onClick={() => handleDelete(event)}>× Deletar</button>
                  </div>
                </div>

                {event.description && <p className="event-description">{event.description}</p>}

                <div className="event-meta">
                  <div className="event-meta-item">
                    <span>◷</span>
                    <span>
                      {new Date(event.start_date).toLocaleString('pt-BR')} –{' '}
                      {new Date(event.end_date).toLocaleString('pt-BR')}
                    </span>
                  </div>
                  {event.location && (
                    <div className="event-meta-item">
                      <span>📍</span>
                      <span>{event.location}</span>
                    </div>
                  )}
                  {event.all_day && <span className="event-badge">Dia inteiro</span>}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">Novo Evento</h2>
              <button type="button" className="modal-close" onClick={handleCloseModal}>×</button>
            </div>

            <form onSubmit={handleSubmit} className="modal-body">
              {formError && <div className="error-message">{formError}</div>}

              <div className="form-group">
                <label className="form-label">Título *</label>
                <input type="text" className="form-input" value={title} onChange={e => setTitle(e.target.value)} placeholder="Título do evento" required autoFocus />
              </div>

              <div className="form-group">
                <label className="form-label">Descrição</label>
                <textarea className="form-textarea" value={description} onChange={e => setDescription(e.target.value)} placeholder="Descrição do evento..." rows={3} />
              </div>

              <div className="datetime-row">
                <div className="form-group">
                  <label className="form-label">Início *</label>
                  <input type="datetime-local" className="form-input" value={startDate} onChange={e => setStartDate(e.target.value)} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Fim *</label>
                  <input type="datetime-local" className="form-input" value={endDate} onChange={e => setEndDate(e.target.value)} required />
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Local</label>
                <input type="text" className="form-input" value={location} onChange={e => setLocation(e.target.value)} placeholder="Local do evento" />
              </div>

              <div className="form-group">
                <div className="checkbox-group">
                  <input type="checkbox" id="allDay" checked={allDay} onChange={e => setAllDay(e.target.checked)} />
                  <label htmlFor="allDay">Evento de dia inteiro</label>
                </div>
              </div>

              {gcalConnected && (
                <div className="form-group">
                  <div className="checkbox-group gcal-sync-check">
                    <input type="checkbox" id="syncCal" checked={syncCalendar} onChange={e => setSyncCalendar(e.target.checked)} />
                    <label htmlFor="syncCal">📅 Adicionar ao Google Calendar</label>
                  </div>
                </div>
              )}

              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={handleCloseModal} disabled={formLoading}>Cancelar</button>
                <button type="submit" className="btn btn-primary" disabled={formLoading}>
                  {formLoading ? 'Criando...' : 'Criar Evento'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
