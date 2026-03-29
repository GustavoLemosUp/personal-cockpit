package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	"personal-cockpit/models"
)

// WAStatus representa o estado da conexão com o WhatsApp.
type WAStatus string

const (
	WADisconnected WAStatus = "disconnected"
	WAConnecting   WAStatus = "connecting"
	WAConnected    WAStatus = "connected"
	WALoggedOut    WAStatus = "logged_out"
)

// WAStatusInfo é retornado ao frontend com o estado atual.
type WAStatusInfo struct {
	Status WAStatus `json:"status"`
	Phone  string   `json:"phone"`
}

// WhatsAppService gerencia conexão, mensagens e agendamentos.
type WhatsAppService struct {
	client        *whatsmeow.Client
	conn          *sql.DB
	dbPath        string
	status        WAStatus
	phone         string
	mu            sync.RWMutex
	schedulerStop chan struct{}
	onQR          func(qrBase64 string)
	onStatus      func(status WAStatus)
	onMessage     func(msg models.WAMessage)
}

func NewWhatsAppService(conn *sql.DB) (*WhatsAppService, error) {
	dataDir, err := getWhatsAppDataDir()
	if err != nil {
		return nil, err
	}
	return &WhatsAppService{
		conn:   conn,
		dbPath: filepath.Join(dataDir, "whatsapp.db"),
		status: WADisconnected,
	}, nil
}

func (s *WhatsAppService) SetOnQR(fn func(string))                  { s.onQR = fn }
func (s *WhatsAppService) SetOnStatus(fn func(WAStatus))            { s.onStatus = fn }
func (s *WhatsAppService) SetOnMessage(fn func(models.WAMessage))   { s.onMessage = fn }

// ─── Conexão ────────────────────────────────────────────────

func (s *WhatsAppService) Connect(ctx context.Context) error {
	if !s.setConnecting() {
		return nil
	}
	client, err := s.buildClient(ctx)
	if err != nil {
		s.setStatus(WADisconnected)
		return err
	}
	s.client = client

	if s.client.Store.ID == nil {
		s.startQRLoop(ctx)
	}

	if err := s.client.Connect(); err != nil {
		s.setStatus(WADisconnected)
		return fmt.Errorf("erro ao conectar: %w", err)
	}
	return nil
}

func (s *WhatsAppService) Disconnect() {
	s.stopScheduler()
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect()
	}
	s.setStatus(WADisconnected)
}

func (s *WhatsAppService) Logout(ctx context.Context) error {
	s.stopScheduler()
	if s.client == nil {
		return nil
	}
	if err := s.client.Logout(ctx); err != nil {
		return fmt.Errorf("erro ao fazer logout: %w", err)
	}
	s.mu.Lock()
	s.phone = ""
	s.mu.Unlock()
	s.setStatus(WALoggedOut)
	return nil
}

func (s *WhatsAppService) GetStatus() WAStatusInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return WAStatusInfo{Status: s.status, Phone: s.phone}
}

// ─── Mensagens ──────────────────────────────────────────────

// GetChats retorna todas as conversas ordenadas pela última mensagem.
func (s *WhatsAppService) GetChats() ([]models.WAChat, error) {
	rows, err := s.conn.Query(`
		SELECT jid, name, last_message, last_message_at, unread_count
		FROM wa_chats
		ORDER BY last_message_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.WAChat
	for rows.Next() {
		var c models.WAChat
		if err := rows.Scan(&c.JID, &c.Name, &c.LastMessage, &c.LastMessageAt, &c.UnreadCount); err != nil {
			continue
		}
		chats = append(chats, c)
	}
	return chats, nil
}

// GetMessages retorna mensagens de uma conversa com paginação.
// offset=0 retorna as mais recentes; offset>0 retorna mensagens mais antigas (scroll infinito).
func (s *WhatsAppService) GetMessages(chatJID string, limit, offset int) ([]models.WAMessage, error) {
	if limit <= 0 {
		limit = 10
	}
	// Busca em ordem DESC (mais recentes primeiro) depois inverte para exibição cronológica.
	rows, err := s.conn.Query(`
		SELECT id, chat_jid, sender_jid, content, is_from_me, timestamp
		FROM wa_messages
		WHERE chat_jid = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, chatJID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.WAMessage
	for rows.Next() {
		var m models.WAMessage
		var isFromMe int
		if err := rows.Scan(&m.ID, &m.ChatJID, &m.SenderJID, &m.Content, &isFromMe, &m.Timestamp); err != nil {
			continue
		}
		m.IsFromMe = isFromMe == 1
		msgs = append(msgs, m)
	}
	// Inverte para ordem cronológica (mais antigas primeiro).
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

// MarkAsRead zera o contador de não lidas de uma conversa.
func (s *WhatsAppService) MarkAsRead(chatJID string) error {
	_, err := s.conn.Exec(`UPDATE wa_chats SET unread_count = 0 WHERE jid = ?`, chatJID)
	return err
}

// SendTextMessage envia uma mensagem de texto e a persiste localmente.
func (s *WhatsAppService) SendTextMessage(ctx context.Context, chatJID, text string) error {
	if s.client == nil || !s.client.IsConnected() {
		return fmt.Errorf("WhatsApp não está conectado")
	}
	jid, err := types.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("JID inválido: %w", err)
	}
	msg := &waE2E.Message{Conversation: proto.String(text)}
	resp, err := s.client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("erro ao enviar: %w", err)
	}

	myJID := ""
	if s.client.Store.ID != nil {
		myJID = s.client.Store.ID.String()
	}
	s.storeMessage(resp.ID, chatJID, myJID, text, true, resp.Timestamp)
	s.upsertChat(chatJID, chatJID, text, resp.Timestamp)
	return nil
}

// ─── Agendamentos ───────────────────────────────────────────

// ScheduleMessage agenda uma mensagem para envio futuro.
func (s *WhatsAppService) ScheduleMessage(input models.WAScheduleInput) error {
	if strings.TrimSpace(input.Content) == "" {
		return fmt.Errorf("mensagem não pode ser vazia")
	}
	if input.ScheduledAt.Before(time.Now()) {
		return fmt.Errorf("horário agendado deve ser no futuro")
	}
	_, err := s.conn.Exec(`
		INSERT INTO wa_scheduled (chat_jid, chat_name, content, scheduled_at)
		VALUES (?, ?, ?, ?)
	`, input.ChatJID, input.ChatName, input.Content, input.ScheduledAt)
	return err
}

// GetScheduled retorna os agendamentos pendentes.
func (s *WhatsAppService) GetScheduled() ([]models.WAScheduled, error) {
	rows, err := s.conn.Query(`
		SELECT id, chat_jid, chat_name, content, scheduled_at, sent, created_at
		FROM wa_scheduled
		WHERE sent = 0
		ORDER BY scheduled_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.WAScheduled
	for rows.Next() {
		var m models.WAScheduled
		var sent int
		if err := rows.Scan(&m.ID, &m.ChatJID, &m.ChatName, &m.Content, &m.ScheduledAt, &sent, &m.CreatedAt); err != nil {
			continue
		}
		m.Sent = sent == 1
		list = append(list, m)
	}
	return list, nil
}

// DeleteScheduled cancela um agendamento.
func (s *WhatsAppService) DeleteScheduled(id int) error {
	_, err := s.conn.Exec(`DELETE FROM wa_scheduled WHERE id = ? AND sent = 0`, id)
	return err
}

// ─── Scheduler ──────────────────────────────────────────────

func (s *WhatsAppService) startScheduler() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.schedulerStop != nil {
		return
	}
	s.schedulerStop = make(chan struct{})
	go s.runScheduler()
}

func (s *WhatsAppService) stopScheduler() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.schedulerStop != nil {
		close(s.schedulerStop)
		s.schedulerStop = nil
	}
}

func (s *WhatsAppService) runScheduler() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.processPendingScheduled()
		case <-s.schedulerStop:
			return
		}
	}
}

func (s *WhatsAppService) processPendingScheduled() {
	rows, err := s.conn.Query(`
		SELECT id, chat_jid, content FROM wa_scheduled
		WHERE sent = 0 AND scheduled_at <= ?
	`, time.Now())
	if err != nil {
		return
	}
	defer rows.Close()

	type pending struct {
		id      int
		chatJID string
		content string
	}
	var items []pending
	for rows.Next() {
		var p pending
		if err := rows.Scan(&p.id, &p.chatJID, &p.content); err == nil {
			items = append(items, p)
		}
	}
	rows.Close()

	for _, p := range items {
		ctx := context.Background()
		if err := s.SendTextMessage(ctx, p.chatJID, p.content); err == nil {
			s.conn.Exec(`UPDATE wa_scheduled SET sent = 1 WHERE id = ?`, p.id)
		}
	}
}

// ─── Internos ───────────────────────────────────────────────

func (s *WhatsAppService) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		s.mu.Lock()
		s.status = WAConnected
		if s.client != nil && s.client.Store.ID != nil {
			s.phone = s.client.Store.ID.User
		}
		s.mu.Unlock()
		if s.onStatus != nil {
			s.onStatus(WAConnected)
		}
		s.startScheduler()

	case *events.Disconnected:
		_ = v
		s.stopScheduler()
		s.setStatus(WADisconnected)

	case *events.LoggedOut:
		_ = v
		s.stopScheduler()
		s.mu.Lock()
		s.phone = ""
		s.mu.Unlock()
		s.setStatus(WALoggedOut)

	case *events.Message:
		s.handleIncomingMessage(v)
	}
}

func (s *WhatsAppService) handleIncomingMessage(v *events.Message) {
	content := extractText(v.Message)
	if content == "" {
		return
	}

	chatJID := v.Info.Chat.String()
	senderJID := v.Info.Sender.String()
	name := v.Info.PushName
	if name == "" {
		name = v.Info.Sender.User
	}

	s.storeMessage(v.Info.ID, chatJID, senderJID, content, v.Info.IsFromMe, v.Info.Timestamp)
	s.upsertChat(chatJID, name, content, v.Info.Timestamp)

	if !v.Info.IsFromMe {
		s.incrementUnread(chatJID)
	}

	if s.onMessage != nil {
		s.onMessage(models.WAMessage{
			ID:        v.Info.ID,
			ChatJID:   chatJID,
			SenderJID: senderJID,
			Content:   content,
			IsFromMe:  v.Info.IsFromMe,
			Timestamp: v.Info.Timestamp,
		})
	}
}

func (s *WhatsAppService) storeMessage(id, chatJID, senderJID, content string, isFromMe bool, ts time.Time) {
	fromMeInt := 0
	if isFromMe {
		fromMeInt = 1
	}
	s.conn.Exec(`
		INSERT OR IGNORE INTO wa_messages (id, chat_jid, sender_jid, content, is_from_me, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, chatJID, senderJID, content, fromMeInt, ts)
}

func (s *WhatsAppService) upsertChat(jid, name, lastMessage string, lastAt time.Time) {
	s.conn.Exec(`
		INSERT INTO wa_chats (jid, name, last_message, last_message_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(jid) DO UPDATE SET
			name = excluded.name,
			last_message = excluded.last_message,
			last_message_at = excluded.last_message_at
	`, jid, name, lastMessage, lastAt)
}

func (s *WhatsAppService) incrementUnread(chatJID string) {
	s.conn.Exec(`UPDATE wa_chats SET unread_count = unread_count + 1 WHERE jid = ?`, chatJID)
}

func (s *WhatsAppService) setConnecting() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status == WAConnecting || s.status == WAConnected {
		return false
	}
	s.status = WAConnecting
	return true
}

func (s *WhatsAppService) buildClient(ctx context.Context) (*whatsmeow.Client, error) {
	container, err := sqlstore.New(ctx, "sqlite", "file:"+s.dbPath+"?_pragma=foreign_keys(1)", waLog.Noop)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco do WhatsApp: %w", err)
	}
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter dispositivo: %w", err)
	}
	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	client.AddEventHandler(s.handleEvent)
	return client, nil
}

func (s *WhatsAppService) startQRLoop(ctx context.Context) {
	qrChan, _ := s.client.GetQRChannel(ctx)
	go func() {
		for evt := range qrChan {
			if evt.Event == "code" && s.onQR != nil {
				if img, err := generateQRBase64(evt.Code); err == nil {
					s.onQR(img)
				}
			}
		}
	}()
}

func (s *WhatsAppService) setStatus(status WAStatus) {
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()
	if s.onStatus != nil {
		s.onStatus(status)
	}
}

func extractText(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	if t := msg.GetConversation(); t != "" {
		return t
	}
	if ext := msg.GetExtendedTextMessage(); ext != nil {
		return ext.GetText()
	}
	return ""
}

func generateQRBase64(code string) (string, error) {
	png, err := qrcode.Encode(code, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}

func getWhatsAppDataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, "Personal Cockpit")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return appDir, nil
}
