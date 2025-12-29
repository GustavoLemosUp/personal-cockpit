import React, { useState, useEffect } from 'react';
import { 
  GetAllCategories, 
  CreateCategory, 
  DeleteCategory 
} from '../../../wailsjs/go/main/App';
import { models } from '../../../wailsjs/go/models';

export function CategoriesPage() {
  const [categories, setCategories] = useState<models.Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  
  // Form state
  const [name, setName] = useState('');
  const [color, setColor] = useState('#0071e3');
  const [type, setType] = useState('general');
  const [formError, setFormError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      const result = await GetAllCategories();
      setCategories(result || []);
    } catch (err) {
      console.error('Erro ao carregar categorias:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');
    setFormLoading(true);
    
    try {
      const category = new models.Category();
      category.name = name;
      category.color = color;
      category.type = type;
      
      await CreateCategory(category);
      
      setName('');
      setColor('#0071e3');
      setType('general');
      setShowModal(false);
      loadCategories();
    } catch (err: any) {
      setFormError(err.message || 'Erro ao criar categoria');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Deseja realmente deletar esta categoria?')) {
      return;
    }
    
    try {
      await DeleteCategory(id);
      loadCategories();
    } catch (err) {
      alert('Erro ao deletar categoria');
    }
  };

  const handleCloseModal = () => {
    setName('');
    setColor('#0071e3');
    setType('general');
    setFormError('');
    setShowModal(false);
  };

  if (loading) {
    return <div className="page">Carregando...</div>;
  }

  return (
    <div className="page">
      <div className="page-header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1 className="page-title">Categorias</h1>
            <p className="page-subtitle">Organize suas tarefas e notas</p>
          </div>
          <button 
            className="btn btn-primary"
            onClick={() => setShowModal(true)}
          >
            + Nova Categoria
          </button>
        </div>
      </div>

      <div className="page-content">
        {categories.length === 0 ? (
          <div className="empty-state">
            <p>🏷️</p>
            <p>Nenhuma categoria criada</p>
          </div>
        ) : (
          <div className="categories-grid">
            {categories.map((category) => (
              <div 
                key={category.id} 
                className="category-card"
                style={{ '--color': category.color } as React.CSSProperties}
              >
                <div className="category-header">
                  <h3 className="category-name">{category.name}</h3>
                  <span className="category-type">{category.type}</span>
                </div>
                
                <div className="category-footer">
                  <span className="category-date">
                    {new Date(category.created_at).toLocaleDateString('pt-BR')}
                  </span>
                  <div className="category-actions">
                    <button 
                      className="btn btn-danger"
                      onClick={() => handleDelete(category.id)}
                      style={{ fontSize: '0.875rem', padding: '4px 8px' }}
                    >
                      ×
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">Nova Categoria</h2>
              <button 
                type="button" 
                className="modal-close"
                onClick={handleCloseModal}
              >
                ×
              </button>
            </div>

            <form onSubmit={handleSubmit} className="modal-body">
              {formError && <div className="error-message">{formError}</div>}

              <div className="form-group">
                <label className="form-label">Nome *</label>
                <input
                  type="text"
                  className="form-input"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Nome da categoria"
                  required
                  autoFocus
                />
              </div>
              
              <div className="form-group">
                <label className="form-label">Tipo</label>
                <select
                  className="form-select"
                  value={type}
                  onChange={(e) => setType(e.target.value)}
                >
                  <option value="general">Geral</option>
                  <option value="task">Tarefa</option>
                  <option value="note">Nota</option>
                </select>
              </div>
              
              <div className="form-group">
                <label className="form-label">Cor</label>
                <div className="color-picker-group">
                  <input
                    type="color"
                    className="color-input"
                    value={color}
                    onChange={(e) => setColor(e.target.value)}
                  />
                  <input
                    type="text"
                    className="form-input"
                    value={color}
                    onChange={(e) => setColor(e.target.value)}
                    placeholder="#0071e3"
                    style={{ flex: 1 }}
                  />
                </div>
              </div>
              
              <div className="modal-footer">
                <button 
                  type="button" 
                  className="btn btn-secondary"
                  onClick={handleCloseModal}
                  disabled={formLoading}
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={formLoading}
                >
                  {formLoading ? 'Criando...' : 'Criar Categoria'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}