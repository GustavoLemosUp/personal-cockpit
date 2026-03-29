import { useState, useEffect } from 'react';
import {
    GetAccounts, CreateAccount, DeleteAccount,
    GetTransactions, CreateTransaction, DeleteTransaction,
    GetScheduledTransactions, CreateScheduledTransaction, ExecuteScheduledTransaction, DeleteScheduledTransaction,
    GetFinanceSummary,
} from '../../../wailsjs/go/main/App';

type FinanceTab = 'overview' | 'transactions' | 'accounts' | 'scheduled';

interface Account {
    id: number; name: string; type: string; balance: number;
    currency: string; color: string; icon: string; description: string;
}

interface Transaction {
    id: number; account_id: number; account_name: string;
    type: string; category: string; description: string;
    amount: number; date: string; tags: string;
}

interface ScheduledTx {
    id: number; account_id: number; account_name: string;
    type: string; category: string; description: string;
    amount: number; scheduled_at: string; is_recurring: boolean; recurrence_rule: string;
}

interface Summary {
    total_balance: number; total_income: number; total_expense: number;
    accounts: { account: Account }[];
    recent_transactions: Transaction[];
}

const ACCOUNT_TYPES = [
    { value: 'checking', label: 'Conta Corrente', icon: '🏦' },
    { value: 'savings',  label: 'Poupança',       icon: '🐷' },
    { value: 'wallet',   label: 'Carteira',        icon: '👛' },
    { value: 'credit',   label: 'Cartão de Crédito', icon: '💳' },
    { value: 'invest',   label: 'Investimentos',   icon: '📈' },
];

const EXPENSE_CATS = ['Alimentação', 'Transporte', 'Moradia', 'Saúde', 'Educação', 'Lazer', 'Vestuário', 'Serviços', 'Outros'];
const INCOME_CATS  = ['Salário', 'Freelance', 'Investimentos', 'Aluguel', 'Presente', 'Outros'];

const fmtBRL = (v: number) => v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
const fmtDate = (iso: string) => new Date(iso).toLocaleDateString('pt-BR');

export function FinancePage() {
    const [tab, setTab] = useState<FinanceTab>('overview');
    const [summary, setSummary] = useState<Summary | null>(null);
    const [accounts, setAccounts] = useState<Account[]>([]);
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [scheduled, setScheduled] = useState<ScheduledTx[]>([]);
    const [error, setError] = useState('');

    // Filters
    const [filterType, setFilterType] = useState('');
    const [filterAccount, setFilterAccount] = useState<number | undefined>();

    // New account form
    const [showAccForm, setShowAccForm] = useState(false);
    const [newAcc, setNewAcc] = useState({ name: '', type: 'checking', balance: 0, color: '#6b7280', description: '' });

    // New transaction form
    const [showTxForm, setShowTxForm] = useState(false);
    const [newTx, setNewTx] = useState({ account_id: 0, type: 'expense', category: '', description: '', amount: 0, date: todayStr() });

    // New scheduled form
    const [showSchForm, setShowSchForm] = useState(false);
    const [newSch, setNewSch] = useState({ account_id: 0, type: 'expense', category: '', description: '', amount: 0, scheduled_at: '', is_recurring: false, recurrence_rule: '' });

    useEffect(() => { load(); }, []);

    async function load() {
        try {
            const [sum, accs] = await Promise.all([GetFinanceSummary(), GetAccounts()]);
            setSummary(sum as any);
            setAccounts(accs ?? []);
        } catch (e) { setError(String(e)); }
    }

    async function loadTransactions() {
        try {
            const list = await GetTransactions({
                type: filterType,
                account_id: filterAccount,
                limit: 50,
                offset: 0,
            } as any);
            setTransactions(list ?? []);
        } catch (e) { setError(String(e)); }
    }

    async function loadScheduled() {
        try {
            const list = await GetScheduledTransactions();
            setScheduled(list ?? []);
        } catch (e) { setError(String(e)); }
    }

    useEffect(() => {
        if (tab === 'transactions') loadTransactions();
        if (tab === 'scheduled') loadScheduled();
    }, [tab, filterType, filterAccount]);

    async function handleCreateAccount() {
        if (!newAcc.name.trim()) return;
        try {
            await CreateAccount(newAcc as any);
            setShowAccForm(false);
            setNewAcc({ name: '', type: 'checking', balance: 0, color: '#6b7280', description: '' });
            load();
        } catch (e) { setError(String(e)); }
    }

    async function handleDeleteAccount(id: number) {
        if (!confirm('Desativar esta conta?')) return;
        await DeleteAccount(id);
        load();
    }

    async function handleCreateTx() {
        if (!newTx.description.trim() || newTx.amount <= 0 || !newTx.account_id) return;
        try {
            await CreateTransaction({ ...newTx, date: new Date(newTx.date).toISOString() } as any);
            setShowTxForm(false);
            setNewTx({ account_id: 0, type: 'expense', category: '', description: '', amount: 0, date: todayStr() });
            load();
            if (tab === 'transactions') loadTransactions();
        } catch (e) { setError(String(e)); }
    }

    async function handleDeleteTx(id: number) {
        await DeleteTransaction(id);
        loadTransactions();
        load();
    }

    async function handleCreateSch() {
        if (!newSch.description.trim() || newSch.amount <= 0 || !newSch.account_id || !newSch.scheduled_at) return;
        try {
            await CreateScheduledTransaction({ ...newSch, scheduled_at: new Date(newSch.scheduled_at).toISOString() } as any);
            setShowSchForm(false);
            setNewSch({ account_id: 0, type: 'expense', category: '', description: '', amount: 0, scheduled_at: '', is_recurring: false, recurrence_rule: '' });
            loadScheduled();
        } catch (e) { setError(String(e)); }
    }

    async function handleExecuteSch(id: number) {
        await ExecuteScheduledTransaction(id);
        loadScheduled();
        load();
    }

    async function handleDeleteSch(id: number) {
        await DeleteScheduledTransaction(id);
        loadScheduled();
    }

    const accountTypeLabel = (type: string) => ACCOUNT_TYPES.find(t => t.value === type)?.label ?? type;
    const accountTypeIcon  = (type: string) => ACCOUNT_TYPES.find(t => t.value === type)?.icon ?? '💰';

    return (
        <div className="finance-page">
            <div className="page-header">
                <h1 className="page-title">Finanças</h1>
                <p className="page-subtitle">Gestão financeira pessoal</p>
            </div>

            <div className="finance-tabs">
                {(['overview','transactions','accounts','scheduled'] as FinanceTab[]).map(t => (
                    <button key={t} className={`finance-tab ${tab === t ? 'active' : ''}`} onClick={() => setTab(t)}>
                        {{ overview: 'Visão Geral', transactions: 'Transações', accounts: 'Contas', scheduled: 'Agendados' }[t]}
                    </button>
                ))}
            </div>

            {error && <div className="finance-error">{error}<button onClick={() => setError('')}>✕</button></div>}

            {/* ── Overview ── */}
            {tab === 'overview' && summary && (
                <div className="finance-overview">
                    <div className="finance-summary-cards">
                        <div className="finance-summary-card balance">
                            <span className="fs-label">Patrimônio Total</span>
                            <span className="fs-value">{fmtBRL(summary.total_balance)}</span>
                        </div>
                        <div className="finance-summary-card income">
                            <span className="fs-label">Receitas do mês</span>
                            <span className="fs-value">{fmtBRL(summary.total_income)}</span>
                        </div>
                        <div className="finance-summary-card expense">
                            <span className="fs-label">Despesas do mês</span>
                            <span className="fs-value">{fmtBRL(summary.total_expense)}</span>
                        </div>
                        <div className={`finance-summary-card ${summary.total_income - summary.total_expense >= 0 ? 'income' : 'expense'}`}>
                            <span className="fs-label">Saldo do mês</span>
                            <span className="fs-value">{fmtBRL(summary.total_income - summary.total_expense)}</span>
                        </div>
                    </div>

                    <div className="finance-accounts-row">
                        {summary.accounts.map(({ account: a }) => (
                            <div key={a.id} className="finance-account-card" style={{ borderColor: a.color }}>
                                <div className="fac-icon" style={{ background: a.color + '22', color: a.color }}>
                                    {accountTypeIcon(a.type)}
                                </div>
                                <div className="fac-info">
                                    <span className="fac-name">{a.name}</span>
                                    <span className="fac-type">{accountTypeLabel(a.type)}</span>
                                </div>
                                <span className={`fac-balance ${a.balance >= 0 ? '' : 'negative'}`}>{fmtBRL(a.balance)}</span>
                            </div>
                        ))}
                        <button className="finance-account-card add-card" onClick={() => { setTab('accounts'); setShowAccForm(true); }}>
                            + Nova Conta
                        </button>
                    </div>

                    <div className="finance-recent">
                        <div className="finance-section-header">
                            <h3>Transações Recentes</h3>
                            <button className="btn btn-sm btn-secondary" onClick={() => setTab('transactions')}>Ver todas</button>
                        </div>
                        <div className="finance-tx-list">
                            {summary.recent_transactions?.map(tx => (
                                <TxRow key={tx.id} tx={tx} onDelete={() => {}} />
                            ))}
                            {(!summary.recent_transactions || summary.recent_transactions.length === 0) && (
                                <p className="finance-empty">Nenhuma transação ainda.</p>
                            )}
                        </div>
                    </div>
                </div>
            )}

            {/* ── Transactions ── */}
            {tab === 'transactions' && (
                <div className="finance-section">
                    <div className="finance-section-header">
                        <div className="finance-filters">
                            <select className="finance-select" value={filterType} onChange={e => setFilterType(e.target.value)}>
                                <option value="">Todos os tipos</option>
                                <option value="income">Receitas</option>
                                <option value="expense">Despesas</option>
                                <option value="transfer">Transferências</option>
                            </select>
                            <select className="finance-select" value={filterAccount ?? ''} onChange={e => setFilterAccount(e.target.value ? +e.target.value : undefined)}>
                                <option value="">Todas as contas</option>
                                {accounts.map(a => <option key={a.id} value={a.id}>{a.name}</option>)}
                            </select>
                        </div>
                        <button className="btn btn-primary btn-sm" onClick={() => setShowTxForm(v => !v)}>+ Nova Transação</button>
                    </div>

                    {showTxForm && (
                        <div className="finance-form">
                            <div className="finance-form-row">
                                <select className="finance-select" value={newTx.type} onChange={e => setNewTx(v => ({ ...v, type: e.target.value }))}>
                                    <option value="expense">Despesa</option>
                                    <option value="income">Receita</option>
                                    <option value="transfer">Transferência</option>
                                </select>
                                <select className="finance-select" value={newTx.account_id} onChange={e => setNewTx(v => ({ ...v, account_id: +e.target.value }))}>
                                    <option value={0}>Selecione a conta</option>
                                    {accounts.map(a => <option key={a.id} value={a.id}>{a.name}</option>)}
                                </select>
                            </div>
                            <div className="finance-form-row">
                                <input className="finance-input" placeholder="Descrição" value={newTx.description} onChange={e => setNewTx(v => ({ ...v, description: e.target.value }))} />
                                <select className="finance-select" value={newTx.category} onChange={e => setNewTx(v => ({ ...v, category: e.target.value }))}>
                                    <option value="">Categoria</option>
                                    {(newTx.type === 'income' ? INCOME_CATS : EXPENSE_CATS).map(c => <option key={c} value={c}>{c}</option>)}
                                </select>
                            </div>
                            <div className="finance-form-row">
                                <input className="finance-input" type="number" min="0" step="0.01" placeholder="Valor (R$)" value={newTx.amount || ''} onChange={e => setNewTx(v => ({ ...v, amount: +e.target.value }))} />
                                <input className="finance-input" type="date" value={newTx.date} onChange={e => setNewTx(v => ({ ...v, date: e.target.value }))} />
                                <button className="btn btn-primary btn-sm" onClick={handleCreateTx}>Salvar</button>
                            </div>
                        </div>
                    )}

                    <div className="finance-tx-list">
                        {transactions.map(tx => <TxRow key={tx.id} tx={tx} onDelete={handleDeleteTx} />)}
                        {transactions.length === 0 && <p className="finance-empty">Nenhuma transação encontrada.</p>}
                    </div>
                </div>
            )}

            {/* ── Accounts ── */}
            {tab === 'accounts' && (
                <div className="finance-section">
                    <div className="finance-section-header">
                        <h3>Contas</h3>
                        <button className="btn btn-primary btn-sm" onClick={() => setShowAccForm(v => !v)}>+ Nova Conta</button>
                    </div>

                    {showAccForm && (
                        <div className="finance-form">
                            <div className="finance-form-row">
                                <input className="finance-input" placeholder="Nome da conta" value={newAcc.name} onChange={e => setNewAcc(v => ({ ...v, name: e.target.value }))} />
                                <select className="finance-select" value={newAcc.type} onChange={e => setNewAcc(v => ({ ...v, type: e.target.value }))}>
                                    {ACCOUNT_TYPES.map(t => <option key={t.value} value={t.value}>{t.icon} {t.label}</option>)}
                                </select>
                            </div>
                            <div className="finance-form-row">
                                <input className="finance-input" type="number" step="0.01" placeholder="Saldo inicial" value={newAcc.balance || ''} onChange={e => setNewAcc(v => ({ ...v, balance: +e.target.value }))} />
                                <input className="finance-input" type="color" value={newAcc.color} onChange={e => setNewAcc(v => ({ ...v, color: e.target.value }))} title="Cor" style={{ width: 48, padding: '2px' }} />
                                <input className="finance-input" placeholder="Descrição (opcional)" value={newAcc.description} onChange={e => setNewAcc(v => ({ ...v, description: e.target.value }))} />
                                <button className="btn btn-primary btn-sm" onClick={handleCreateAccount}>Salvar</button>
                            </div>
                        </div>
                    )}

                    <div className="finance-accounts-grid">
                        {accounts.map(a => (
                            <div key={a.id} className="finance-account-detail" style={{ borderLeftColor: a.color }}>
                                <div className="fad-header">
                                    <span className="fad-icon">{accountTypeIcon(a.type)}</span>
                                    <div className="fad-info">
                                        <span className="fad-name">{a.name}</span>
                                        <span className="fad-type">{accountTypeLabel(a.type)}</span>
                                    </div>
                                    <button className="icon-btn danger" onClick={() => handleDeleteAccount(a.id)} title="Remover">✕</button>
                                </div>
                                <div className={`fad-balance ${a.balance >= 0 ? 'positive' : 'negative'}`}>{fmtBRL(a.balance)}</div>
                                {a.description && <p className="fad-desc">{a.description}</p>}
                            </div>
                        ))}
                        {accounts.length === 0 && <p className="finance-empty">Nenhuma conta cadastrada.</p>}
                    </div>
                </div>
            )}

            {/* ── Scheduled ── */}
            {tab === 'scheduled' && (
                <div className="finance-section">
                    <div className="finance-section-header">
                        <h3>Movimentos Agendados</h3>
                        <button className="btn btn-primary btn-sm" onClick={() => setShowSchForm(v => !v)}>+ Novo Agendamento</button>
                    </div>

                    {showSchForm && (
                        <div className="finance-form">
                            <div className="finance-form-row">
                                <select className="finance-select" value={newSch.type} onChange={e => setNewSch(v => ({ ...v, type: e.target.value }))}>
                                    <option value="expense">Despesa</option>
                                    <option value="income">Receita</option>
                                </select>
                                <select className="finance-select" value={newSch.account_id} onChange={e => setNewSch(v => ({ ...v, account_id: +e.target.value }))}>
                                    <option value={0}>Selecione a conta</option>
                                    {accounts.map(a => <option key={a.id} value={a.id}>{a.name}</option>)}
                                </select>
                            </div>
                            <div className="finance-form-row">
                                <input className="finance-input" placeholder="Descrição" value={newSch.description} onChange={e => setNewSch(v => ({ ...v, description: e.target.value }))} />
                                <select className="finance-select" value={newSch.category} onChange={e => setNewSch(v => ({ ...v, category: e.target.value }))}>
                                    <option value="">Categoria</option>
                                    {(newSch.type === 'income' ? INCOME_CATS : EXPENSE_CATS).map(c => <option key={c} value={c}>{c}</option>)}
                                </select>
                            </div>
                            <div className="finance-form-row">
                                <input className="finance-input" type="number" min="0" step="0.01" placeholder="Valor (R$)" value={newSch.amount || ''} onChange={e => setNewSch(v => ({ ...v, amount: +e.target.value }))} />
                                <input className="finance-input" type="datetime-local" value={newSch.scheduled_at} onChange={e => setNewSch(v => ({ ...v, scheduled_at: e.target.value }))} />
                                <label className="finance-check">
                                    <input type="checkbox" checked={newSch.is_recurring} onChange={e => setNewSch(v => ({ ...v, is_recurring: e.target.checked }))} />
                                    Recorrente
                                </label>
                                {newSch.is_recurring && (
                                    <select className="finance-select" value={newSch.recurrence_rule} onChange={e => setNewSch(v => ({ ...v, recurrence_rule: e.target.value }))}>
                                        <option value="monthly">Mensal</option>
                                        <option value="weekly">Semanal</option>
                                        <option value="yearly">Anual</option>
                                    </select>
                                )}
                                <button className="btn btn-primary btn-sm" onClick={handleCreateSch}>Salvar</button>
                            </div>
                        </div>
                    )}

                    <div className="finance-tx-list">
                        {scheduled.map(s => (
                            <div key={s.id} className={`finance-tx-row ${s.type}`}>
                                <div className="ftx-left">
                                    <span className="ftx-type-badge">{s.type === 'income' ? '↑' : '↓'}</span>
                                    <div>
                                        <span className="ftx-desc">{s.description}</span>
                                        <span className="ftx-meta">{s.account_name} • {s.category} • {new Date(s.scheduled_at).toLocaleString('pt-BR')}</span>
                                        {s.is_recurring && <span className="ftx-recurring">↻ {s.recurrence_rule}</span>}
                                    </div>
                                </div>
                                <div className="ftx-right">
                                    <span className={`ftx-amount ${s.type}`}>{fmtBRL(s.amount)}</span>
                                    <button className="btn btn-sm btn-secondary" onClick={() => handleExecuteSch(s.id)} title="Executar agora">✓</button>
                                    <button className="icon-btn danger" onClick={() => handleDeleteSch(s.id)}>✕</button>
                                </div>
                            </div>
                        ))}
                        {scheduled.length === 0 && <p className="finance-empty">Nenhum agendamento pendente.</p>}
                    </div>
                </div>
            )}
        </div>
    );
}

function TxRow({ tx, onDelete }: { tx: Transaction; onDelete: (id: number) => void }) {
    return (
        <div className={`finance-tx-row ${tx.type}`}>
            <div className="ftx-left">
                <span className="ftx-type-badge">{tx.type === 'income' ? '↑' : tx.type === 'expense' ? '↓' : '⇄'}</span>
                <div>
                    <span className="ftx-desc">{tx.description}</span>
                    <span className="ftx-meta">{tx.account_name}{tx.category ? ` • ${tx.category}` : ''} • {fmtDate(tx.date)}</span>
                </div>
            </div>
            <div className="ftx-right">
                <span className={`ftx-amount ${tx.type}`}>{fmtBRL(tx.amount)}</span>
                {onDelete && <button className="icon-btn danger" onClick={() => onDelete(tx.id)}>✕</button>}
            </div>
        </div>
    );
}

function todayStr() {
    return new Date().toISOString().slice(0, 10);
}
