import { useEffect, useState } from 'react';
import { SetTheme } from '../../wailsjs/go/main/App';

export function useTheme() {
  const [theme, setTheme] = useState<'light' | 'dark' | 'auto'>('auto');

  useEffect(() => {
    // Carregar preferência salva
    const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | 'auto' | null;
    if (savedTheme) {
      setTheme(savedTheme);
      applyTheme(savedTheme);
    } else {
      // Auto por padrão
      applyTheme('auto');
    }
  }, []);

  const applyTheme = (newTheme: 'light' | 'dark' | 'auto') => {
    if (newTheme === 'auto') {
      // Remove atributo e deixa CSS media query decidir
      document.documentElement.removeAttribute('data-theme');
      
      // Detecta tema do sistema
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      SetTheme(isDark ? 'dark' : 'light').catch(console.error);
    } else {
      // Força tema específico
      document.documentElement.setAttribute('data-theme', newTheme);
      SetTheme(newTheme).catch(console.error);
    }
  };

  const changeTheme = (newTheme: 'light' | 'dark' | 'auto') => {
    setTheme(newTheme);
    localStorage.setItem('theme', newTheme);
    applyTheme(newTheme);
  };

  // Escuta mudanças no tema do sistema quando está em "auto"
  useEffect(() => {
    if (theme === 'auto') {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      
      const handleChange = (e: MediaQueryListEvent) => {
        SetTheme(e.matches ? 'dark' : 'light').catch(console.error);
      };

      mediaQuery.addEventListener('change', handleChange);
      return () => mediaQuery.removeEventListener('change', handleChange);
    }
  }, [theme]);

  return { theme, changeTheme };
}