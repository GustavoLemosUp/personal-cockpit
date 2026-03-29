import { useState, useEffect } from 'react';
import {
    GetProducts, CreateProduct, UpdateProduct, DeleteProduct, GetLowStock,
    AddStockMovement, GetStockMovements,
    GetShoppingLists, CreateShoppingList, DeleteShoppingList,
    GetShoppingListWithItems, AddShoppingItem, ToggleShoppingItem, DeleteShoppingItem,
} from '../../../wailsjs/go/main/App';

type StockTab = 'estoque' | 'compras';

interface Product {
    id: number; name: string; category: string; unit: string;
    min_quantity: number; current_stock: number; price: number; notes: string;
}

interface StockMovement {
    id: number; product_id: number; product_name: string;
    type: string; quantity: number; price: number; notes: string; date: string;
}

interface StockAlert { product: Product; deficit: number; }

interface ShoppingList {
    id: number; name: string; month: string; description: string;
    is_completed: boolean; total_budget: number; total_spent: number;
    items?: ShoppingItem[];
}

interface ShoppingItem {
    id: number; shopping_list_id: number; name: string; quantity: number;
    unit: string; estimated_price: number; actual_price: number;
    category: string; is_bought: boolean; notes: string;
}

const STOCK_CATS = ['Alimentos', 'Limpeza', 'Higiene', 'Bebidas', 'Frios', 'Hortifrúti', 'Padaria', 'Outros'];
const UNITS = ['un', 'kg', 'g', 'L', 'mL', 'cx', 'pct', 'dz'];

const fmtBRL = (v: number) => v > 0 ? v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) : '—';

export function StockPage() {
    const [tab, setTab] = useState<StockTab>('estoque');
    const [products, setProducts] = useState<Product[]>([]);
    const [alerts, setAlerts] = useState<StockAlert[]>([]);
    const [lists, setLists] = useState<ShoppingList[]>([]);
    const [selectedList, setSelectedList] = useState<ShoppingList | null>(null);
    const [movements, setMovements] = useState<StockMovement[]>([]);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [error, setError] = useState('');
    const [filterCat, setFilterCat] = useState('');

    // Forms
    const [showProdForm, setShowProdForm] = useState(false);
    const [newProd, setNewProd] = useState({ name: '', category: '', unit: 'un', min_quantity: 0, current_stock: 0, price: 0, notes: '' });
    const [showMoveForm, setShowMoveForm] = useState(false);
    const [newMove, setNewMove] = useState({ product_id: 0, type: 'in', quantity: 1, price: 0, notes: '' });

    const [showListForm, setShowListForm] = useState(false);
    const [newList, setNewList] = useState({ name: '', month: currMonth(), description: '', total_budget: 0 });
    const [showItemForm, setShowItemForm] = useState(false);
    const [newItem, setNewItem] = useState({ name: '', quantity: 1, unit: 'un', estimated_price: 0, category: '', notes: '' });

    useEffect(() => { loadProducts(); }, []);

    async function loadProducts() {
        try {
            const [prods, alts] = await Promise.all([GetProducts(), GetLowStock()]);
            setProducts(prods ?? []);
            setAlerts(alts ?? []);
        } catch (e) { setError(String(e)); }
    }

    async function loadLists() {
        try {
            const l = await GetShoppingLists();
            setLists(l ?? []);
        } catch (e) { setError(String(e)); }
    }

    useEffect(() => {
        if (tab === 'compras') loadLists();
    }, [tab]);

    async function selectList(id: number) {
        try {
            const l = await GetShoppingListWithItems(id);
            setSelectedList(l as any);
        } catch (e) { setError(String(e)); }
    }

    async function selectProduct(p: Product) {
        setSelectedProduct(p);
        setNewMove(v => ({ ...v, product_id: p.id }));
        const mvs = await GetStockMovements(p.id, 20);
        setMovements(mvs ?? []);
        setShowMoveForm(false);
    }

    async function handleCreateProd() {
        if (!newProd.name.trim()) return;
        try {
            await CreateProduct(newProd as any);
            setShowProdForm(false);
            setNewProd({ name: '', category: '', unit: 'un', min_quantity: 0, current_stock: 0, price: 0, notes: '' });
            loadProducts();
        } catch (e) { setError(String(e)); }
    }

    async function handleDeleteProd(id: number) {
        if (!confirm('Remover produto?')) return;
        await DeleteProduct(id);
        loadProducts();
        if (selectedProduct?.id === id) setSelectedProduct(null);
    }

    async function handleMove() {
        if (newMove.quantity <= 0 || !newMove.product_id) return;
        try {
            await AddStockMovement(newMove as any);
            setShowMoveForm(false);
            setNewMove(v => ({ ...v, quantity: 1, price: 0, notes: '' }));
            loadProducts();
            if (selectedProduct) {
                const mvs = await GetStockMovements(selectedProduct.id, 20);
                setMovements(mvs ?? []);
            }
        } catch (e) { setError(String(e)); }
    }

    async function handleCreateList() {
        if (!newList.name.trim()) return;
        try {
            await CreateShoppingList(newList as any);
            setShowListForm(false);
            setNewList({ name: '', month: currMonth(), description: '', total_budget: 0 });
            loadLists();
        } catch (e) { setError(String(e)); }
    }

    async function handleDeleteList(id: number) {
        if (!confirm('Excluir lista?')) return;
        await DeleteShoppingList(id);
        if (selectedList?.id === id) setSelectedList(null);
        loadLists();
    }

    async function handleAddItem() {
        if (!newItem.name.trim() || !selectedList) return;
        try {
            await AddShoppingItem({ ...newItem, shopping_list_id: selectedList.id } as any);
            setNewItem({ name: '', quantity: 1, unit: 'un', estimated_price: 0, category: '', notes: '' });
            selectList(selectedList.id);
        } catch (e) { setError(String(e)); }
    }

    async function handleToggleItem(id: number) {
        await ToggleShoppingItem(id);
        if (selectedList) selectList(selectedList.id);
    }

    async function handleDeleteItem(id: number) {
        await DeleteShoppingItem(id);
        if (selectedList) selectList(selectedList.id);
    }

    const filtered = filterCat ? products.filter(p => p.category === filterCat) : products;
    const categories = [...new Set(products.map(p => p.category).filter(Boolean))];

    return (
        <div className="stock-page">
            <div className="page-header">
                <h1 className="page-title">Estoque & Compras</h1>
                <p className="page-subtitle">Controle o estoque da sua casa e organize suas compras</p>
            </div>

            <div className="stock-tabs">
                <button className={`stock-tab ${tab === 'estoque' ? 'active' : ''}`} onClick={() => setTab('estoque')}>Estoque</button>
                <button className={`stock-tab ${tab === 'compras' ? 'active' : ''}`} onClick={() => setTab('compras')}>Listas de Compras</button>
            </div>

            {error && <div className="stock-error">{error}<button onClick={() => setError('')}>✕</button></div>}

            {/* ── Estoque ── */}
            {tab === 'estoque' && (
                <div className="stock-layout">
                    {/* Lista de produtos */}
                    <div className="stock-left">
                        {alerts.length > 0 && (
                            <div className="stock-alerts">
                                <span className="stock-alert-title">⚠️ Estoque baixo ({alerts.length})</span>
                                {alerts.map(a => (
                                    <div key={a.product.id} className="stock-alert-item" onClick={() => selectProduct(a.product)}>
                                        <span>{a.product.name}</span>
                                        <span className="alert-deficit">falta {a.deficit.toFixed(1)} {a.product.unit}</span>
                                    </div>
                                ))}
                            </div>
                        )}

                        <div className="stock-list-header">
                            <select className="stock-select" value={filterCat} onChange={e => setFilterCat(e.target.value)}>
                                <option value="">Todas categorias</option>
                                {categories.map(c => <option key={c} value={c}>{c}</option>)}
                            </select>
                            <button className="btn btn-primary btn-sm" onClick={() => setShowProdForm(v => !v)}>+ Produto</button>
                        </div>

                        {showProdForm && (
                            <div className="stock-form">
                                <input className="stock-input" placeholder="Nome do produto *" value={newProd.name} onChange={e => setNewProd(v => ({ ...v, name: e.target.value }))} />
                                <div className="stock-form-row">
                                    <select className="stock-select" value={newProd.category} onChange={e => setNewProd(v => ({ ...v, category: e.target.value }))}>
                                        <option value="">Categoria</option>
                                        {STOCK_CATS.map(c => <option key={c} value={c}>{c}</option>)}
                                    </select>
                                    <select className="stock-select" value={newProd.unit} onChange={e => setNewProd(v => ({ ...v, unit: e.target.value }))}>
                                        {UNITS.map(u => <option key={u} value={u}>{u}</option>)}
                                    </select>
                                </div>
                                <div className="stock-form-row">
                                    <input className="stock-input" type="number" min="0" step="0.01" placeholder="Qtd inicial" value={newProd.current_stock || ''} onChange={e => setNewProd(v => ({ ...v, current_stock: +e.target.value }))} />
                                    <input className="stock-input" type="number" min="0" step="0.01" placeholder="Qtd mínima" value={newProd.min_quantity || ''} onChange={e => setNewProd(v => ({ ...v, min_quantity: +e.target.value }))} />
                                    <input className="stock-input" type="number" min="0" step="0.01" placeholder="Preço ref." value={newProd.price || ''} onChange={e => setNewProd(v => ({ ...v, price: +e.target.value }))} />
                                </div>
                                <input className="stock-input" placeholder="Observações" value={newProd.notes} onChange={e => setNewProd(v => ({ ...v, notes: e.target.value }))} />
                                <button className="btn btn-primary btn-sm" onClick={handleCreateProd}>Salvar</button>
                            </div>
                        )}

                        <div className="stock-product-list">
                            {filtered.map(p => (
                                <div key={p.id} className={`stock-product-item ${selectedProduct?.id === p.id ? 'selected' : ''} ${p.current_stock < p.min_quantity ? 'low' : ''}`}
                                     onClick={() => selectProduct(p)}>
                                    <div className="spi-info">
                                        <span className="spi-name">{p.name}</span>
                                        <span className="spi-cat">{p.category}</span>
                                    </div>
                                    <div className="spi-stock">
                                        <span className={`spi-qty ${p.current_stock < p.min_quantity ? 'low' : ''}`}>
                                            {p.current_stock} {p.unit}
                                        </span>
                                        {p.min_quantity > 0 && <span className="spi-min">min: {p.min_quantity}</span>}
                                    </div>
                                </div>
                            ))}
                            {filtered.length === 0 && <p className="stock-empty">Nenhum produto cadastrado.</p>}
                        </div>
                    </div>

                    {/* Detalhe do produto */}
                    <div className="stock-right">
                        {!selectedProduct ? (
                            <div className="stock-welcome">
                                <span>📦</span>
                                <p>Selecione um produto para ver detalhes e movimentar estoque</p>
                            </div>
                        ) : (
                            <div className="stock-detail">
                                <div className="stock-detail-header">
                                    <div>
                                        <h3>{selectedProduct.name}</h3>
                                        <span className="sd-meta">{selectedProduct.category} • {selectedProduct.unit}</span>
                                    </div>
                                    <button className="icon-btn danger" onClick={() => handleDeleteProd(selectedProduct.id)}>✕</button>
                                </div>

                                <div className="stock-stat-row">
                                    <div className="stock-stat">
                                        <span>Em estoque</span>
                                        <strong className={selectedProduct.current_stock < selectedProduct.min_quantity ? 'low' : ''}>
                                            {selectedProduct.current_stock} {selectedProduct.unit}
                                        </strong>
                                    </div>
                                    <div className="stock-stat">
                                        <span>Mínimo</span>
                                        <strong>{selectedProduct.min_quantity} {selectedProduct.unit}</strong>
                                    </div>
                                    <div className="stock-stat">
                                        <span>Preço ref.</span>
                                        <strong>{fmtBRL(selectedProduct.price)}</strong>
                                    </div>
                                </div>

                                <div className="stock-move-section">
                                    <button className="btn btn-primary btn-sm" onClick={() => setShowMoveForm(v => !v)}>
                                        {showMoveForm ? 'Cancelar' : '+ Movimentação'}
                                    </button>
                                    {showMoveForm && (
                                        <div className="stock-form" style={{ marginTop: 8 }}>
                                            <div className="stock-form-row">
                                                <select className="stock-select" value={newMove.type} onChange={e => setNewMove(v => ({ ...v, type: e.target.value }))}>
                                                    <option value="in">Entrada (compra)</option>
                                                    <option value="out">Saída (consumo)</option>
                                                </select>
                                                <input className="stock-input" type="number" min="0.01" step="0.01" placeholder="Quantidade" value={newMove.quantity || ''} onChange={e => setNewMove(v => ({ ...v, quantity: +e.target.value }))} />
                                                <input className="stock-input" type="number" min="0" step="0.01" placeholder="Preço (opcional)" value={newMove.price || ''} onChange={e => setNewMove(v => ({ ...v, price: +e.target.value }))} />
                                            </div>
                                            <input className="stock-input" placeholder="Observação" value={newMove.notes} onChange={e => setNewMove(v => ({ ...v, notes: e.target.value }))} />
                                            <button className="btn btn-primary btn-sm" onClick={handleMove}>Registrar</button>
                                        </div>
                                    )}
                                </div>

                                <div className="stock-history">
                                    <h4>Histórico</h4>
                                    {movements.map(m => (
                                        <div key={m.id} className={`stock-move-row ${m.type}`}>
                                            <span className="smr-badge">{m.type === 'in' ? '↑' : '↓'}</span>
                                            <div className="smr-info">
                                                <span>{m.type === 'in' ? 'Entrada' : 'Saída'} — {m.quantity} {selectedProduct.unit}</span>
                                                <span className="smr-date">{new Date(m.date).toLocaleDateString('pt-BR')}{m.notes ? ` — ${m.notes}` : ''}</span>
                                            </div>
                                            {m.price > 0 && <span className="smr-price">{fmtBRL(m.price)}</span>}
                                        </div>
                                    ))}
                                    {movements.length === 0 && <p className="stock-empty" style={{ fontSize: '0.8rem' }}>Sem movimentações.</p>}
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* ── Listas de Compras ── */}
            {tab === 'compras' && (
                <div className="shopping-layout">
                    <div className="shopping-left">
                        <div className="stock-list-header">
                            <h3>Listas</h3>
                            <button className="btn btn-primary btn-sm" onClick={() => setShowListForm(v => !v)}>+ Nova Lista</button>
                        </div>

                        {showListForm && (
                            <div className="stock-form">
                                <input className="stock-input" placeholder="Nome da lista *" value={newList.name} onChange={e => setNewList(v => ({ ...v, name: e.target.value }))} />
                                <div className="stock-form-row">
                                    <input className="stock-input" type="month" value={newList.month} onChange={e => setNewList(v => ({ ...v, month: e.target.value }))} />
                                    <input className="stock-input" type="number" min="0" step="0.01" placeholder="Orçamento (R$)" value={newList.total_budget || ''} onChange={e => setNewList(v => ({ ...v, total_budget: +e.target.value }))} />
                                </div>
                                <input className="stock-input" placeholder="Descrição" value={newList.description} onChange={e => setNewList(v => ({ ...v, description: e.target.value }))} />
                                <button className="btn btn-primary btn-sm" onClick={handleCreateList}>Salvar</button>
                            </div>
                        )}

                        <div className="shopping-list-items">
                            {lists.map(l => (
                                <div key={l.id} className={`shopping-list-card ${selectedList?.id === l.id ? 'selected' : ''} ${l.is_completed ? 'done' : ''}`}
                                     onClick={() => selectList(l.id)}>
                                    <div className="slc-info">
                                        <span className="slc-name">{l.name}</span>
                                        <span className="slc-month">{l.month}</span>
                                    </div>
                                    <div className="slc-meta">
                                        {l.total_budget > 0 && <span className="slc-budget">orç: {fmtBRL(l.total_budget)}</span>}
                                        <button className="icon-btn danger" onClick={e => { e.stopPropagation(); handleDeleteList(l.id); }}>✕</button>
                                    </div>
                                </div>
                            ))}
                            {lists.length === 0 && <p className="stock-empty">Nenhuma lista criada.</p>}
                        </div>
                    </div>

                    <div className="shopping-right">
                        {!selectedList ? (
                            <div className="stock-welcome">
                                <span>🛒</span>
                                <p>Selecione uma lista para ver e gerenciar os itens</p>
                            </div>
                        ) : (
                            <div className="shopping-detail">
                                <div className="stock-detail-header">
                                    <div>
                                        <h3>{selectedList.name}</h3>
                                        <span className="sd-meta">{selectedList.month}</span>
                                    </div>
                                    <div className="shopping-budget">
                                        {selectedList.total_budget > 0 && (
                                            <>
                                                <span>Gasto: {fmtBRL(selectedList.total_spent)}</span>
                                                <span>/ Orçamento: {fmtBRL(selectedList.total_budget)}</span>
                                            </>
                                        )}
                                    </div>
                                </div>

                                <div className="shopping-add-item">
                                    <button className="btn btn-secondary btn-sm" onClick={() => setShowItemForm(v => !v)}>
                                        {showItemForm ? 'Cancelar' : '+ Adicionar Item'}
                                    </button>
                                    {showItemForm && (
                                        <div className="stock-form" style={{ marginTop: 8 }}>
                                            <div className="stock-form-row">
                                                <input className="stock-input" placeholder="Nome do item *" value={newItem.name} onChange={e => setNewItem(v => ({ ...v, name: e.target.value }))} />
                                                <select className="stock-select" value={newItem.category} onChange={e => setNewItem(v => ({ ...v, category: e.target.value }))}>
                                                    <option value="">Categoria</option>
                                                    {STOCK_CATS.map(c => <option key={c} value={c}>{c}</option>)}
                                                </select>
                                            </div>
                                            <div className="stock-form-row">
                                                <input className="stock-input" type="number" min="0.01" step="0.01" placeholder="Qtd" value={newItem.quantity || ''} onChange={e => setNewItem(v => ({ ...v, quantity: +e.target.value }))} />
                                                <select className="stock-select" value={newItem.unit} onChange={e => setNewItem(v => ({ ...v, unit: e.target.value }))}>
                                                    {UNITS.map(u => <option key={u} value={u}>{u}</option>)}
                                                </select>
                                                <input className="stock-input" type="number" min="0" step="0.01" placeholder="Preço est." value={newItem.estimated_price || ''} onChange={e => setNewItem(v => ({ ...v, estimated_price: +e.target.value }))} />
                                                <button className="btn btn-primary btn-sm" onClick={handleAddItem}>Adicionar</button>
                                            </div>
                                        </div>
                                    )}
                                </div>

                                {/* Agrupar por categoria */}
                                {(() => {
                                    const items = selectedList.items ?? [];
                                    const cats = [...new Set(items.map(i => i.category || 'Outros'))];
                                    return cats.map(cat => (
                                        <div key={cat} className="shopping-category-group">
                                            <h4 className="shopping-cat-label">{cat}</h4>
                                            {items.filter(i => (i.category || 'Outros') === cat).map(item => (
                                                <div key={item.id} className={`shopping-item ${item.is_bought ? 'bought' : ''}`}>
                                                    <button className="shopping-check" onClick={() => handleToggleItem(item.id)}>
                                                        {item.is_bought ? '✓' : '○'}
                                                    </button>
                                                    <div className="shopping-item-info">
                                                        <span className="si-name">{item.name}</span>
                                                        <span className="si-qty">{item.quantity} {item.unit}{item.estimated_price > 0 ? ` — est. ${fmtBRL(item.estimated_price)}` : ''}</span>
                                                    </div>
                                                    <button className="icon-btn danger" onClick={() => handleDeleteItem(item.id)}>✕</button>
                                                </div>
                                            ))}
                                        </div>
                                    ));
                                })()}
                                {(selectedList.items ?? []).length === 0 && <p className="stock-empty">Lista vazia. Adicione itens acima.</p>}
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}

function currMonth() {
    return new Date().toISOString().slice(0, 7);
}
