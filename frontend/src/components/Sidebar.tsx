import React, { useState } from 'react';
import { useTheme } from '../hooks/useTheme';
import logo from '../assets/images/logo.png';

interface SidebarProps {
  currentPage: string;
  onNavigate: (page: string) => void;
}

export function Sidebar({ currentPage, onNavigate }: SidebarProps) {
  const [collapsed, setCollapsed] = useState(false);
  const { theme, changeTheme } = useTheme();

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: '▦' },
    { id: 'tasks', label: 'Tarefas', icon: '✓' },
    { id: 'notes', label: 'Notas', icon: '◐' },
    { id: 'events', label: 'Eventos', icon: '◷' },
    { id: 'categories', label: 'Categorias', icon: '◈' },
  ];

  return (
    <aside 
      className={`sidebar ${collapsed ? 'collapsed' : ''}`}
      onMouseEnter={() => setCollapsed(false)}
      onMouseLeave={() => setCollapsed(true)}
    >
      <div className="sidebar-header">
        <div className="sidebar-logo">
          <img src={logo} alt="Personal Cockpit" className="logo-image" />
        </div>
        <div className="sidebar-title-wrapper">
          <h2 className="sidebar-title">Personal Cockpit</h2>
          <p className="sidebar-subtitle">v1.0.0</p>
        </div>
      </div>

      <nav className="sidebar-nav">
        <div className="nav-section">
          <p className="nav-section-title">Menu</p>
          {menuItems.map((item) => (
            <button
              key={item.id}
              className={`nav-item ${currentPage === item.id ? 'active' : ''}`}
              onClick={() => onNavigate(item.id)}
            >
              <span className="nav-icon">{item.icon}</span>
              <span className="nav-label">{item.label}</span>
            </button>
          ))}
        </div>
      </nav>

      {/* Tema no footer */}
      <div className="sidebar-footer">
        <div className="theme-switcher">
          <button 
            className={`theme-btn ${theme === 'light' ? 'active' : ''}`}
            onClick={() => changeTheme('light')}
            title="Claro"
          >
            ☀
          </button>
          <button 
            className={`theme-btn ${theme === 'auto' ? 'active' : ''}`}
            onClick={() => changeTheme('auto')}
            title="Auto"
          >
            ◐
          </button>
          <button 
            className={`theme-btn ${theme === 'dark' ? 'active' : ''}`}
            onClick={() => changeTheme('dark')}
            title="Escuro"
          >
            ☾
          </button>
        </div>
      </div>
    </aside>
  );
}