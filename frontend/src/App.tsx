import { useState, useEffect } from 'react';
import './App.css';
import { GetCurrentProfileID, HasAnyProfile } from '../wailsjs/go/main/App';
import { Layout } from './components/Layout';
import { LoginPage } from './components/LoginPage';
import { Dashboard } from './components/Dashboard';
import { TasksPage } from './components/tasks/TasksPage';
import { NotesPage } from './components/notes/NotesPage';
import { EventsPage } from './components/events/EventsPage';
import { CategoriesPage } from './components/categories/CategoriesPage';
import { KanbanPage } from './components/kanban/KanbanPage';
import { SettingsPage } from './components/settings/SettingsPage';
import { WhatsAppPage } from './components/whatsapp/WhatsAppPage';

function App() {
  const [currentPage, setCurrentPage] = useState('dashboard');
  const [profileID, setProfileID] = useState<number | null>(null);
  const [loadingAuth, setLoadingAuth] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const hasAny = await HasAnyProfile();
      if (!hasAny) {
        setProfileID(null);
      } else {
        const id = await GetCurrentProfileID();
        setProfileID(id > 0 ? id : null);
      }
    } catch {
      setProfileID(null);
    } finally {
      setLoadingAuth(false);
    }
  };

  const handleLogin = (id: number) => {
    setProfileID(id);
    setCurrentPage('dashboard');
  };

  const handleLogout = () => {
    setProfileID(null);
  };

  if (loadingAuth) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
        <p>Carregando...</p>
      </div>
    );
  }

  if (!profileID) {
    return <LoginPage onLogin={handleLogin} />;
  }

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':   return <Dashboard />;
      case 'tasks':       return <TasksPage />;
      case 'kanban':      return <KanbanPage profileID={profileID} />;
      case 'notes':       return <NotesPage />;
      case 'events':      return <EventsPage profileID={profileID} />;
      case 'categories':  return <CategoriesPage />;
      case 'settings':    return <SettingsPage profileID={profileID} onLogout={handleLogout} />;
      case 'whatsapp':    return <WhatsAppPage />;
      default:            return <Dashboard />;
    }
  };

  return (
    <Layout currentPage={currentPage} onNavigate={setCurrentPage} profileID={profileID}>
      {renderPage()}
    </Layout>
  );
}

export default App;
