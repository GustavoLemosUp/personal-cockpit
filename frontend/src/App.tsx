import { useState } from 'react';
import './App.css';
import { Layout } from './components/Layout';
import { Dashboard } from './components/Dashboard';
import { TasksPage } from './components/tasks/TasksPage';
import { NotesPage } from './components/notes/NotesPage';
import { EventsPage } from './components/events/EventsPage';  // ← ADICIONAR
import { CategoriesPage } from './components/categories/CategoriesPage';  // ← ADICIONAR

function App() {
  const [currentPage, setCurrentPage] = useState('dashboard');

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <Dashboard />;
      case 'tasks':
        return <TasksPage />;
      case 'notes':
        return <NotesPage />;
      case 'events':
        return <EventsPage />;  // ← USAR
      case 'categories':
        return <CategoriesPage />;  // ← USAR
      default:
        return <Dashboard />;
    }
  };

  return (
    <Layout currentPage={currentPage} onNavigate={setCurrentPage}>
      {renderPage()}
    </Layout>
  );
}

export default App;