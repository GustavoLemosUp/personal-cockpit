import { useState, useEffect, useRef } from 'react';
import {
    ConnectWhatsApp, DisconnectWhatsApp, LogoutWhatsApp, GetWhatsAppStatus,
    GetWAChats, GetWAMessages, MarkWARead, SendWAMessage,
    ScheduleWAMessage, GetWAScheduled, DeleteWAScheduled,
    CreateTask,
} from '../../../wailsjs/go/main/App';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import '../../styles/whatsapp.css';

type WAStatus = 'disconnected' | 'connecting' | 'connected' | 'logged_out';
type Tab = 'chats' | 'scheduled';

interface WAChat {
    jid: string;
    name: string;
    last_message: string;
    last_message_at: string;
    unread_count: number;
}

interface WAMessage {
    id: string;
    chat_jid: string;
    sender_jid: string;
    content: string;
    is_from_me: boolean;
    timestamp: string;
}

interface WAScheduled {
    id: number;
    chat_jid: string;
    chat_name: string;
    content: string;
    scheduled_at: string;
    sent: boolean;
}

export function WhatsAppPage() {
    const [status, setStatus] = useState<WAStatus>('disconnected');
    const [phone, setPhone] = useState('');
    const [qrBase64, setQrBase64] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [tab, setTab] = useState<Tab>('chats');

    // Chat
    const [chats, setChats] = useState<WAChat[]>([]);
    const [selectedChat, setSelectedChat] = useState<WAChat | null>(null);
    const [messages, setMessages] = useState<WAMessage[]>([]);
    const [input, setInput] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);

    // Agendamentos
    const [scheduled, setScheduled] = useState<WAScheduled[]>([]);
    const [scheduleJID, setScheduleJID] = useState('');
    const [scheduleName, setScheduleName] = useState('');
    const [scheduleText, setScheduleText] = useState('');
    const [scheduleAt, setScheduleAt] = useState('');

    useEffect(() => {
        GetWhatsAppStatus().then(info => {
            setStatus(info.status as WAStatus);
            setPhone(info.phone ?? '');
            if (info.status === 'connected') loadChats();
        });

        const offQR = EventsOn('whatsapp:qr', (data: string) => {
            setQrBase64(data);
            setStatus('connecting');
            setLoading(false);
        });

        const offStatus = EventsOn('whatsapp:status', (s: string) => {
            setStatus(s as WAStatus);
            if (s === 'connected') {
                setQrBase64('');
                GetWhatsAppStatus().then(i => setPhone(i.phone ?? ''));
                loadChats();
            }
        });

        const offMsg = EventsOn('whatsapp:message', (_msg: WAMessage) => {
            loadChats();
            if (selectedChat && _msg.chat_jid === selectedChat.jid) {
                setMessages(prev => [...prev, _msg]);
            }
        });

        return () => { offQR(); offStatus(); offMsg(); };
    }, [selectedChat]);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    async function loadChats() {
        try {
            const list = await GetWAChats();
            setChats(list ?? []);
        } catch { /* silencioso */ }
    }

    async function loadScheduled() {
        try {
            const list = await GetWAScheduled();
            setScheduled(list ?? []);
        } catch { /* silencioso */ }
    }

    async function selectChat(chat: WAChat) {
        setSelectedChat(chat);
        const msgs = await GetWAMessages(chat.jid, 100);
        setMessages(msgs ?? []);
        await MarkWARead(chat.jid);
        setChats(prev => prev.map(c => c.jid === chat.jid ? { ...c, unread_count: 0 } : c));
    }

    async function handleSend() {
        if (!input.trim() || !selectedChat) return;
        const text = input.trim();
        setInput('');
        try {
            await SendWAMessage(selectedChat.jid, text);
            const msgs = await GetWAMessages(selectedChat.jid, 100);
            setMessages(msgs ?? []);
            loadChats();
        } catch (err: unknown) {
            setError(String(err));
        }
    }

    async function handleConnect() {
        setLoading(true);
        setError('');
        try { await ConnectWhatsApp(); }
        catch (err: unknown) { setError(String(err)); setLoading(false); }
    }

    async function handleDisconnect() {
        DisconnectWhatsApp();
        setQrBase64('');
        setSelectedChat(null);
        setChats([]);
    }

    async function handleLogout() {
        try { await LogoutWhatsApp(); setQrBase64(''); setSelectedChat(null); setChats([]); }
        catch (err: unknown) { setError(String(err)); }
    }

    async function handleCreateTask(msg: WAMessage) {
        await CreateTask({
            title: msg.content.slice(0, 120),
            description: `Via WhatsApp — ${formatTime(msg.timestamp)}`,
            status: 'pending',
            priority: 'medium',
        } as any);
    }

    async function handleSchedule() {
        if (!scheduleJID.trim() || !scheduleText.trim() || !scheduleAt) return;
        try {
            await ScheduleWAMessage({
                chat_jid: scheduleJID.trim(),
                chat_name: scheduleName.trim() || scheduleJID.trim(),
                content: scheduleText.trim(),
                scheduled_at: new Date(scheduleAt).toISOString(),
            } as any);
            setScheduleJID(''); setScheduleName(''); setScheduleText(''); setScheduleAt('');
            loadScheduled();
        } catch (err: unknown) { setError(String(err)); }
    }

    async function handleDeleteScheduled(id: number) {
        await DeleteWAScheduled(id);
        loadScheduled();
    }

    function formatTime(iso: string) {
        const d = new Date(iso);
        return d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' });
    }

    function formatDate(iso: string) {
        const d = new Date(iso);
        const today = new Date();
        if (d.toDateString() === today.toDateString()) return formatTime(iso);
        return d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' });
    }

    function formatPhone(jid: string) {
        const user = jid.split('@')[0];
        return user.startsWith('55') ? `+${user}` : `+${user}`;
    }

    const statusLabel: Record<WAStatus, string> = {
        disconnected: 'Desconectado',
        connecting: 'Aguardando leitura do QR code…',
        connected: phone ? `Conectado — +${phone}` : 'Conectado',
        logged_out: 'Sessão encerrada',
    };

    // ── Tela de conexão ──
    if (status !== 'connected') {
        return (
            <div className="whatsapp-page">
                <div className="page-header">
                    <h1 className="page-title">WhatsApp</h1>
                    <p className="page-subtitle">Conecte sua conta para usar o WhatsApp no Cockpit</p>
                </div>
                <div className="wa-connect-card">
                    <div className={`wa-status-badge wa-status-${status}`}>
                        <span className="wa-status-dot" />
                        <span>{statusLabel[status]}</span>
                    </div>

                    {qrBase64 && status === 'connecting' && (
                        <div className="wa-qr-section">
                            <p className="wa-qr-instruction">
                                No celular: <strong>WhatsApp → Dispositivos conectados → Conectar um dispositivo</strong>
                            </p>
                            <div className="wa-qr-wrapper">
                                <img src={`data:image/png;base64,${qrBase64}`} alt="QR Code" className="wa-qr-image" />
                            </div>
                            <p className="wa-qr-hint">O QR code expira em ~60s e é renovado automaticamente.</p>
                        </div>
                    )}

                    {error && <div className="wa-error"><strong>Erro:</strong> {error}</div>}

                    <div className="wa-actions">
                        {(status === 'disconnected' || status === 'logged_out') && (
                            <button className="btn btn-primary" onClick={handleConnect} disabled={loading}>
                                {loading ? 'Conectando…' : 'Conectar WhatsApp'}
                            </button>
                        )}
                    </div>
                </div>
            </div>
        );
    }

    // ── Interface de chat ──
    return (
        <div className="wa-app">

            {/* Sidebar esquerda */}
            <div className="wa-left">
                <div className="wa-left-header">
                    <div className="wa-left-status">
                        <span className="wa-dot-connected" />
                        <span className="wa-phone-label">+{phone}</span>
                    </div>
                    <div className="wa-left-actions">
                        <button className="wa-icon-btn" title="Desconectar" onClick={handleDisconnect}>⏏</button>
                        <button className="wa-icon-btn wa-icon-btn-danger" title="Sair da conta" onClick={handleLogout}>✕</button>
                    </div>
                </div>

                <div className="wa-tabs">
                    <button className={`wa-tab ${tab === 'chats' ? 'active' : ''}`} onClick={() => setTab('chats')}>
                        Conversas
                    </button>
                    <button className={`wa-tab ${tab === 'scheduled' ? 'active' : ''}`} onClick={() => { setTab('scheduled'); loadScheduled(); }}>
                        Agendados
                    </button>
                </div>

                {tab === 'chats' && (
                    <div className="wa-chat-list">
                        {chats.length === 0 && (
                            <p className="wa-empty">Nenhuma conversa ainda.<br />As mensagens aparecerão aqui conforme chegarem.</p>
                        )}
                        {chats.map(chat => (
                            <button
                                key={chat.jid}
                                className={`wa-chat-item ${selectedChat?.jid === chat.jid ? 'active' : ''}`}
                                onClick={() => selectChat(chat)}
                            >
                                <div className="wa-chat-avatar">
                                    {(chat.name || chat.jid)[0].toUpperCase()}
                                </div>
                                <div className="wa-chat-info">
                                    <div className="wa-chat-name">{chat.name || formatPhone(chat.jid)}</div>
                                    <div className="wa-chat-last">{chat.last_message}</div>
                                </div>
                                <div className="wa-chat-meta">
                                    <span className="wa-chat-time">{formatDate(chat.last_message_at)}</span>
                                    {chat.unread_count > 0 && (
                                        <span className="wa-unread">{chat.unread_count}</span>
                                    )}
                                </div>
                            </button>
                        ))}
                    </div>
                )}

                {tab === 'scheduled' && (
                    <div className="wa-scheduled-list">
                        <div className="wa-schedule-form">
                            <input className="wa-input-field" placeholder="JID (+5521…@s.whatsapp.net)" value={scheduleJID} onChange={e => setScheduleJID(e.target.value)} />
                            <input className="wa-input-field" placeholder="Nome (opcional)" value={scheduleName} onChange={e => setScheduleName(e.target.value)} />
                            <textarea className="wa-textarea" placeholder="Mensagem…" rows={2} value={scheduleText} onChange={e => setScheduleText(e.target.value)} />
                            <input className="wa-input-field" type="datetime-local" value={scheduleAt} onChange={e => setScheduleAt(e.target.value)} />
                            <button className="btn btn-primary btn-sm" onClick={handleSchedule}>Agendar</button>
                        </div>
                        {scheduled.length === 0 && <p className="wa-empty">Nenhum agendamento pendente.</p>}
                        {scheduled.map(s => (
                            <div key={s.id} className="wa-scheduled-item">
                                <div className="wa-scheduled-info">
                                    <span className="wa-scheduled-name">{s.chat_name || s.chat_jid}</span>
                                    <span className="wa-scheduled-text">{s.content}</span>
                                    <span className="wa-scheduled-time">
                                        {new Date(s.scheduled_at).toLocaleString('pt-BR')}
                                    </span>
                                </div>
                                <button className="wa-icon-btn wa-icon-btn-danger" onClick={() => handleDeleteScheduled(s.id)}>✕</button>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Painel de mensagens */}
            <div className="wa-right">
                {!selectedChat ? (
                    <div className="wa-welcome">
                        <div className="wa-welcome-icon">◉</div>
                        <p>Selecione uma conversa para começar</p>
                    </div>
                ) : (
                    <>
                        <div className="wa-chat-header">
                            <div className="wa-chat-avatar large">
                                {(selectedChat.name || selectedChat.jid)[0].toUpperCase()}
                            </div>
                            <div>
                                <div className="wa-chat-header-name">{selectedChat.name || formatPhone(selectedChat.jid)}</div>
                                <div className="wa-chat-header-jid">{selectedChat.jid}</div>
                            </div>
                        </div>

                        <div className="wa-messages">
                            {messages.map(msg => (
                                <div key={msg.id} className={`wa-msg ${msg.is_from_me ? 'from-me' : 'from-them'}`}>
                                    <div className="wa-msg-bubble">
                                        <p className="wa-msg-text">{msg.content}</p>
                                        <div className="wa-msg-footer">
                                            <span className="wa-msg-time">{formatTime(msg.timestamp)}</span>
                                            {!msg.is_from_me && (
                                                <button
                                                    className="wa-msg-action"
                                                    title="Criar tarefa desta mensagem"
                                                    onClick={() => handleCreateTask(msg)}
                                                >
                                                    + tarefa
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            ))}
                            <div ref={messagesEndRef} />
                        </div>

                        <div className="wa-input-bar">
                            <textarea
                                className="wa-input-msg"
                                placeholder="Escreva uma mensagem…"
                                rows={1}
                                value={input}
                                onChange={e => setInput(e.target.value)}
                                onKeyDown={e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend(); } }}
                            />
                            <button className="wa-send-btn" onClick={handleSend} disabled={!input.trim()}>
                                ➤
                            </button>
                        </div>
                    </>
                )}
            </div>

            {error && <div className="wa-toast-error">{error}<button onClick={() => setError('')}>✕</button></div>}
        </div>
    );
}
