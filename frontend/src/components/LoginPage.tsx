import { useState, useEffect } from 'react';
import {
  HasAnyProfile,
  GetAllProfiles,
  CreateProfile,
  VerifyPin,
  SetCurrentProfile,
} from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';

interface LoginPageProps {
  onLogin: (profileID: number) => void;
}

export function LoginPage({ onLogin }: LoginPageProps) {
  const [profiles, setProfiles] = useState<models.ProfileSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [view, setView] = useState<'select' | 'create' | 'pin'>('select');
  const [selectedProfile, setSelectedProfile] = useState<models.ProfileSummary | null>(null);

  // Create form
  const [newName, setNewName] = useState('');
  const [newEmail, setNewEmail] = useState('');
  const [newPin, setNewPin] = useState('');
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState('');

  // PIN form
  const [pin, setPin] = useState('');
  const [pinError, setPinError] = useState('');

  useEffect(() => {
    loadProfiles();
  }, []);

  const loadProfiles = async () => {
    try {
      const hasAny = await HasAnyProfile();
      if (!hasAny) {
        setView('create');
      } else {
        const list = await GetAllProfiles();
        setProfiles(list || []);
        setView('select');
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectProfile = (profile: models.ProfileSummary) => {
    if (profile.has_pin) {
      setSelectedProfile(profile);
      setPin('');
      setPinError('');
      setView('pin');
    } else {
      doLogin(profile.id);
    }
  };

  const handlePinSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProfile) return;
    try {
      const ok = await VerifyPin(selectedProfile.id, pin);
      if (ok) {
        doLogin(selectedProfile.id);
      } else {
        setPinError('PIN incorreto');
      }
    } catch {
      setPinError('Erro ao verificar PIN');
    }
  };

  const doLogin = async (profileID: number) => {
    await SetCurrentProfile(profileID);
    onLogin(profileID);
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateError('');
    setCreating(true);
    try {
      const input = new models.CreateProfileInput();
      input.name = newName;
      input.email = newEmail;
      input.pin = newPin;
      const profile = await CreateProfile(input);
      await doLogin(profile!.id);
    } catch (err: any) {
      setCreateError(err.message || 'Erro ao criar perfil');
    } finally {
      setCreating(false);
    }
  };

  if (loading) {
    return <div className="login-page"><p>Carregando...</p></div>;
  }

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-header">
          <h1 className="login-title">Personal Cockpit</h1>
          <p className="login-subtitle">Centro de comando da sua família</p>
        </div>

        {/* ── Seletor de perfis ── */}
        {view === 'select' && (
          <div className="login-body">
            <h2 className="login-section-title">Quem é você?</h2>
            <div className="profile-grid">
              {profiles.map((p) => (
                <button
                  key={p.id}
                  className="profile-card"
                  onClick={() => handleSelectProfile(p)}
                >
                  <div className="profile-avatar">
                    {p.avatar_url
                      ? <img src={p.avatar_url} alt={p.name} />
                      : <span>{p.name.charAt(0).toUpperCase()}</span>
                    }
                  </div>
                  <span className="profile-name">{p.name}</span>
                  {p.role === 'admin' && <span className="profile-badge">Admin</span>}
                  {p.has_pin && <span className="profile-lock">🔒</span>}
                </button>
              ))}
            </div>
            <button className="btn btn-secondary login-add-btn" onClick={() => setView('create')}>
              + Adicionar perfil
            </button>
          </div>
        )}

        {/* ── Criar perfil ── */}
        {view === 'create' && (
          <div className="login-body">
            <h2 className="login-section-title">
              {profiles.length === 0 ? 'Bem-vindo! Crie seu perfil de administrador' : 'Novo perfil'}
            </h2>
            <form onSubmit={handleCreate} className="login-form">
              {createError && <div className="error-message">{createError}</div>}
              <div className="form-group">
                <label className="form-label">Nome *</label>
                <input
                  className="form-input"
                  value={newName}
                  onChange={e => setNewName(e.target.value)}
                  placeholder="Seu nome"
                  required
                  autoFocus
                />
              </div>
              <div className="form-group">
                <label className="form-label">Email (Google)</label>
                <input
                  className="form-input"
                  type="email"
                  value={newEmail}
                  onChange={e => setNewEmail(e.target.value)}
                  placeholder="seu@gmail.com"
                />
              </div>
              <div className="form-group">
                <label className="form-label">PIN de acesso (opcional)</label>
                <input
                  className="form-input"
                  type="password"
                  value={newPin}
                  onChange={e => setNewPin(e.target.value)}
                  placeholder="4+ dígitos"
                  maxLength={8}
                />
              </div>
              <div className="login-form-actions">
                {profiles.length > 0 && (
                  <button type="button" className="btn btn-secondary" onClick={() => setView('select')}>
                    Cancelar
                  </button>
                )}
                <button type="submit" className="btn btn-primary" disabled={creating}>
                  {creating ? 'Criando...' : 'Entrar'}
                </button>
              </div>
            </form>
          </div>
        )}

        {/* ── PIN ── */}
        {view === 'pin' && selectedProfile && (
          <div className="login-body">
            <div className="pin-profile">
              <div className="profile-avatar large">
                {selectedProfile.avatar_url
                  ? <img src={selectedProfile.avatar_url} alt={selectedProfile.name} />
                  : <span>{selectedProfile.name.charAt(0).toUpperCase()}</span>
                }
              </div>
              <h2 className="login-section-title">{selectedProfile.name}</h2>
            </div>
            <form onSubmit={handlePinSubmit} className="login-form">
              {pinError && <div className="error-message">{pinError}</div>}
              <div className="form-group">
                <label className="form-label">PIN</label>
                <input
                  className="form-input pin-input"
                  type="password"
                  value={pin}
                  onChange={e => setPin(e.target.value)}
                  placeholder="••••"
                  autoFocus
                  maxLength={8}
                />
              </div>
              <div className="login-form-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setView('select')}>
                  Voltar
                </button>
                <button type="submit" className="btn btn-primary">Entrar</button>
              </div>
            </form>
          </div>
        )}
      </div>
    </div>
  );
}
