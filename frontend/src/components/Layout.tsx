import React from 'react';
import { Sidebar } from './Sidebar';

interface LayoutProps {
  currentPage: string;
  onNavigate: (page: string) => void;
  children: React.ReactNode;
  profileID?: number;
}

export function Layout({ currentPage, onNavigate, children, profileID }: LayoutProps) {
  return (
    <div className="app-layout">
      <Sidebar currentPage={currentPage} onNavigate={onNavigate} profileID={profileID} />
      <main className="main-content">
        {children}
      </main>
    </div>
  );
}
