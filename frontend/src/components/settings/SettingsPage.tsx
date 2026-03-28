import { useState, useEffect, useRef } from 'react';
import {
  GetProfileByID,
  GetAllProfiles,
  CreateProfile,
  DeleteProfile,
  SetPin,
  SetCurrentProfile,
  GoogleCalendarConfigured,
  GoogleCalendarConnected,
  StartGoogleAuth,
  DisconnectGoogleCalendar,
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

interface SettingsPageProps {
  profileID: number;
  onLogout: () => void;
}

export function SettingsPage({ profileID, onLogout }: SettingsPageProps) {
  const [currentProfile, setCurrentProfile] = useState<models.ProfileSummary | null>(null);
  const [profiles, setProfiles] = useState<models.ProfileSummary[]>([]);
  const [tab, setTab] = useState<'members' | 'google' | 'pin'>('members');
  const [loading, setLoading] = useState(true);

  // Add member
  const [newName, setNewName] = useState('');
  const [newEmail, setNewEmail] = useState('');
  const [addError, setAddError] = useState('');

  // PIN
  const [newPin, setNewPin] = useState('');
  const [pinMsg, setPinMsg] = useState('');

  // Google Calendar
  const [gcalConfigured, setGcalConfigured] = useState(false);
  const [gcalConnected, setGcalConnected] = useState(false);
  const [gcalLoading, setGcalLoading] = useState(false);
  const [gcalMsg, setGcalMsg] = useState('');
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    loadData();
    return () => { if (pollRef.current) clearInterval(pollRef.current); };
  }, [profileID]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [p, all, configured, connected] = await Promise.all([
        GetProfileByID(profileID),
        GetAllProfiles(),
        GoogleCalendarConfigured(),
        GoogleCalendarConnected(profileID),
      ]);
      setCurrentProfile(p);
      setProfiles(all || []);
      setGcalConfigured(configured);
      setGcalConnected(connected);
    } finally {
      setLoading(false);
    }
  };

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddError('');
    try {
      const input = new models.CreateProfileInput();
      input.name = newName;
      input.email = newEmail;
      input.role = 'member';
      await CreateProfile(input);
      setNewName('');
      setNewEmail('');
      await loadData();
    } catch (err: any) {
      setAddError(err.message || 'Erro ao adicionar membro');
    }
  };

  const handleRemoveMember = async (id: number) => {
    if (!window.confirm('Remover este membro?')) return;
    try {
      await DeleteProfile(id);
      await loadData();
    } catch (err: any) {
      alert(err.message || 'Erro ao remover membro');
    }
  };

  const handleSetPin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await SetPin(profileID, newPin);
      setPinMsg(newPin ? 'PIN definido com sucesso!' : 'PIN removido.');
      setNewPin('');
    } catch {
      setPinMsg('Erro ao definir PIN.');
    }
  };

  const handleLogout = async () => {
    await SetCurrentProfile(0);
    onLogout();
  };

  // ── Google Calendar ──────────────────────────────────────

  const handleConnectGoogle = async () => {
    setGcalLoading(true);
    setGcalMsg('Abrindo navegador para autorização...');
    try {
      await StartGoogleAuth(profileID);
      setGcalMsg('Autorize no navegador. O app detectará automaticamente quando concluir.');
      // Poll every 2s to detect when token arrives
      pollRef.current = setInterval(async () => {
        const connected = await GoogleCalendarConnected(profileID);
        if (connected) {
          setGcalConnected(true);
          setGcalMsg('✅ Google Calendar conectado com sucesso!');
          setGcalLoading(false);
          if (pollRef.current) clearInterval(pollRef.current);
        }
      }, 2000);
      // Stop polling after 5 minutes
      setTimeout(() => {
        if (pollRef.current) { clearInterval(pollRef.current); setGcalLoading(false); }
      }, 5 * 60 * 1000);
    } catch (err: any) {
      setGcalMsg(err.message || 'Erro ao iniciar autenticação');
      setGcalLoading(false);
    }
  };

  const handleDisconnectGoogle = async () => {
    if (!window.confirm('Desconectar o Google Calendar deste perfil?')) return;
    try {
      await DisconnectGoogleCalendar(profileID);
      setGcalConnected(false);
      setGcalMsg('Google Calendar desconectado.');
    } catch (err: any) {
      setGcalMsg(err.message || 'Erro ao desconectar');
    }
  };

  if (loading) return <div className="page"><p>Carregando...</p></div>;

  const isAdmin = currentProfile?.role === 'admin';

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">Configurações</h1>
        <p className="page-subtitle">Perfil: {currentProfile?.name}</p>
      </div>

      <div className="settings-tabs">
        {isAdmin && (
          <button className={`tab-btn ${tab === 'members' ? 'active' : ''}`} onClick={() => setTab('members')}>
            Membros
          </button>
        )}
        <button className={`tab-btn ${tab === 'google' ? 'active' : ''}`} onClick={() => setTab('google')}>
          Google Calendar
          {gcalConnected && <span style={{ marginLeft: 6, color: 'var(--success)' }}>●</span>}
        </button>
        <button className={`tab-btn ${tab === 'pin' ? 'active' : ''}`} onClick={() => setTab('pin')}>
          PIN de acesso
        </button>
      </div>

      <div className="settings-body">
        {/* Members tab */}
        {tab === 'members' && isAdmin && (
          <div className="settings-section">
            <h2 className="settings-section-title">Membros da Família</h2>
            <div className="members-list">
              {profiles.map(p => (
                <div key={p.id} className="member-item">
                  <div className="profile-avatar">
                    {p.avatar_url ? <img src={p.avatar_url} alt={p.name} /> : <span>{p.name.charAt(0)}</span>}
                  </div>
                  <div className="member-info">
                    <strong>{p.name}</strong>
                    {p.email && <span className="member-email">{p.email}</span>}
                    <span className={`profile-badge ${p.role}`}>{p.role === 'admin' ? 'Admin' : 'Membro'}</span>
                  </div>
                  {p.id !== profileID && (
                    <button className="btn btn-danger btn-sm" onClick={() => handleRemoveMember(p.id)}>Remover</button>
                  )}
                </div>
              ))}
            </div>
            <div className="settings-subsection">
              <h3>Adicionar membro</h3>
              {addError && <div className="error-message">{addError}</div>}
              <form onSubmit={handleAddMember} className="add-member-form">
                <input className="form-input" value={newName} onChange={e => setNewName(e.target.value)} placeholder="Nome" required />
                <input className="form-input" type="email" value={newEmail} onChange={e => setNewEmail(e.target.value)} placeholder="Email Google (opcional)" />
                <button type="submit" className="btn btn-primary">Adicionar</button>
              </form>
              <p className="settings-hint">O novo membro poderá selecionar o perfil na próxima abertura do app.</p>
            </div>
          </div>
        )}

        {/* Google Calendar tab */}
        {tab === 'google' && (
          <div className="settings-section">
            <h2 className="settings-section-title">Google Calendar</h2>

            {!gcalConfigured ? (
              <div className="gcal-setup-notice">
                <div className="gcal-setup-icon">📋</div>
                <h3>Configuração necessária</h3>
                <p>Para sincronizar com o Google Calendar, adicione suas credenciais OAuth2 no arquivo <code>.env</code>:</p>
                <div className="gcal-code-block">
                  <code>GOOGLE_CLIENT_ID=seu_client_id_aqui</code>
                  <code>GOOGLE_CLIENT_SECRET=seu_client_secret_aqui</code>
                </div>
                <div className="gcal-steps">
                  <p><strong>Como obter as credenciais:</strong></p>
                  <ol>
                    <li>Acesse <strong>console.cloud.google.com</strong></li>
                    <li>Crie um projeto → Habilite a <strong>Google Calendar API</strong></li>
                    <li>Credenciais → Criar → ID do cliente OAuth → Tipo: <strong>Aplicativo para computador</strong></li>
                    <li>Em URIs de redirecionamento, adicione: <code>http://localhost:9876/oauth/callback</code></li>
                    <li>Copie o Client ID e Secret para o arquivo <code>.env</code></li>
                    <li>Reinicie o app</li>
                  </ol>
                </div>
              </div>
            ) : gcalConnected ? (
              <div className="gcal-connected">
                <div className="gcal-status-row">
                  <span className="gcal-status-dot connected" />
                  <div>
                    <strong>Conectado</strong>
                    <p className="settings-hint">Eventos e tarefas podem ser sincronizados com sua agenda do Google.</p>
                  </div>
                </div>
                {gcalMsg && <div className="success-message">{gcalMsg}</div>}
                <button className="btn btn-danger" onClick={handleDisconnectGoogle}>
                  Desconectar Google Calendar
                </button>
              </div>
            ) : (
              <div className="gcal-connect">
                <div className="gcal-status-row">
                  <span className="gcal-status-dot" />
                  <div>
                    <strong>Não conectado</strong>
                    <p className="settings-hint">Conecte sua conta para sincronizar eventos e tarefas com o Google Calendar.</p>
                  </div>
                </div>
                {gcalMsg && (
                  <div className={gcalMsg.startsWith('✅') ? 'success-message' : 'error-message'}>
                    {gcalMsg}
                  </div>
                )}
                <button className="btn btn-primary" onClick={handleConnectGoogle} disabled={gcalLoading}>
                  {gcalLoading ? 'Aguardando autorização...' : '🔗 Conectar Google Calendar'}
                </button>
              </div>
            )}
          </div>
        )}

        {/* PIN tab */}
        {tab === 'pin' && (
          <div className="settings-section">
            <h2 className="settings-section-title">PIN de acesso</h2>
            <p className="settings-hint">Defina um PIN para proteger seu perfil na tela de login. Deixe em branco para remover.</p>
            {pinMsg && <div className="success-message">{pinMsg}</div>}
            <form onSubmit={handleSetPin} className="login-form">
              <div className="form-group">
                <label className="form-label">Novo PIN</label>
                <input className="form-input" type="password" value={newPin} onChange={e => setNewPin(e.target.value)} placeholder="Deixe em branco para remover" maxLength={8} />
              </div>
              <button type="submit" className="btn btn-primary">Salvar PIN</button>
            </form>
          </div>
        )}
      </div>

      <div className="settings-logout">
        <button className="btn btn-secondary" onClick={handleLogout}>Trocar perfil</button>
      </div>
    </div>
  );
}
