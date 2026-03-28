package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"personal-cockpit/models"
)

const (
	googleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
	calendarAPIURL = "https://www.googleapis.com/calendar/v3"
	oauthPort      = "9876"
	oauthRedirect  = "http://localhost:9876/oauth/callback"
	calendarScope  = "https://www.googleapis.com/auth/calendar"
)

type CalendarService struct {
	db           *sql.DB
	clientID     string
	clientSecret string
}

func NewCalendarService(db *sql.DB, clientID, clientSecret string) *CalendarService {
	return &CalendarService{db: db, clientID: clientID, clientSecret: clientSecret}
}

// IsConfigured returns true if OAuth credentials are present in config.
func (s *CalendarService) IsConfigured() bool {
	return s.clientID != "" && s.clientSecret != ""
}

// HasToken returns true if a profile already has a saved token.
func (s *CalendarService) HasToken(profileID int) bool {
	var tok string
	s.db.QueryRow(`SELECT COALESCE(google_access_token,'') FROM profiles WHERE id = ?`, profileID).Scan(&tok)
	return tok != ""
}

// StartOAuthFlow builds and returns the Google auth URL, then starts a local
// callback server that exchanges the code for a token automatically.
func (s *CalendarService) StartOAuthFlow(profileID int) (string, error) {
	if !s.IsConfigured() {
		return "", fmt.Errorf("credenciais Google não configuradas — adicione GOOGLE_CLIENT_ID e GOOGLE_CLIENT_SECRET no arquivo .env")
	}

	params := url.Values{
		"client_id":     {s.clientID},
		"redirect_uri":  {oauthRedirect},
		"response_type": {"code"},
		"scope":         {calendarScope},
		"access_type":   {"offline"},
		"prompt":        {"consent"},
	}
	authURL := googleAuthURL + "?" + params.Encode()

	go s.runCallbackServer(profileID)

	return authURL, nil
}

// runCallbackServer starts a one-shot HTTP server on port 9876 that captures
// the OAuth code from Google's redirect, exchanges it for tokens, and stops.
func (s *CalendarService) runCallbackServer(profileID int) {
	done := make(chan struct{}, 1)

	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":" + oauthPort, Handler: mux}

	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if code == "" {
			fmt.Fprint(w, htmlPage("❌ Código não recebido", "Tente novamente nas Configurações."))
		} else if err := s.exchangeCode(profileID, code); err != nil {
			fmt.Fprint(w, htmlPage("❌ Erro ao conectar", err.Error()))
		} else {
			fmt.Fprint(w, htmlPage("✅ Google Calendar conectado!", "Pode fechar esta aba e voltar ao app."))
		}

		done <- struct{}{}
	})

	go srv.ListenAndServe() //nolint
	<-done
	srv.Close()
}

func htmlPage(title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8">
<style>body{font-family:sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0;background:#f5f5f7}
.card{background:#fff;padding:40px;border-radius:16px;text-align:center;box-shadow:0 8px 32px rgba(0,0,0,.12)}</style></head>
<body><div class="card"><h2>%s</h2><p>%s</p><script>setTimeout(()=>window.close(),3000)</script></div></body></html>`,
		title, body)
}

// ── Token management ────────────────────────────────────────

type googleToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"` // unix timestamp
}

func (s *CalendarService) exchangeCode(profileID int, code string) error {
	resp, err := http.PostForm(googleTokenURL, url.Values{
		"code":          {code},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"redirect_uri":  {oauthRedirect},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return fmt.Errorf("erro na troca de código: %w", err)
	}
	defer resp.Body.Close()

	var tok googleToken
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return fmt.Errorf("erro ao decodificar token: %w", err)
	}
	if tok.AccessToken == "" {
		return fmt.Errorf("token vazio — verifique as credenciais OAuth")
	}

	tok.ExpiresAt = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second).Unix()
	return s.saveToken(profileID, tok)
}

func (s *CalendarService) saveToken(profileID int, tok googleToken) error {
	b, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`UPDATE profiles SET google_access_token = ? WHERE id = ?`, string(b), profileID)
	return err
}

func (s *CalendarService) getValidToken(profileID int) (string, error) {
	var raw string
	err := s.db.QueryRow(`SELECT COALESCE(google_access_token,'') FROM profiles WHERE id = ?`, profileID).Scan(&raw)
	if err != nil || raw == "" {
		return "", fmt.Errorf("Google Calendar não conectado — vá em Configurações → Google Calendar")
	}

	var tok googleToken
	if err := json.Unmarshal([]byte(raw), &tok); err != nil {
		return "", fmt.Errorf("token inválido — reconecte o Google Calendar")
	}

	// Refresh if expiring within 60 seconds
	if time.Now().Unix() >= tok.ExpiresAt-60 {
		if tok.RefreshToken == "" {
			return "", fmt.Errorf("token expirado — reconecte o Google Calendar nas Configurações")
		}
		refreshed, err := s.refreshToken(tok.RefreshToken)
		if err != nil {
			return "", err
		}
		refreshed.RefreshToken = tok.RefreshToken // preserve refresh token
		refreshed.ExpiresAt = time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second).Unix()
		_ = s.saveToken(profileID, *refreshed)
		return refreshed.AccessToken, nil
	}

	return tok.AccessToken, nil
}

func (s *CalendarService) refreshToken(refreshToken string) (*googleToken, error) {
	resp, err := http.PostForm(googleTokenURL, url.Values{
		"refresh_token": {refreshToken},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao renovar token: %w", err)
	}
	defer resp.Body.Close()

	var tok googleToken
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, fmt.Errorf("erro ao decodificar token renovado: %w", err)
	}
	if tok.AccessToken == "" {
		return nil, fmt.Errorf("token renovado vazio — reconecte o Google Calendar")
	}
	return &tok, nil
}

// DisconnectProfile removes the stored token for a profile.
func (s *CalendarService) DisconnectProfile(profileID int) error {
	_, err := s.db.Exec(`UPDATE profiles SET google_access_token = '' WHERE id = ?`, profileID)
	return err
}

// ── Calendar API ────────────────────────────────────────────

type gcalEvent struct {
	Summary     string       `json:"summary"`
	Description string       `json:"description,omitempty"`
	Location    string       `json:"location,omitempty"`
	Start       gcalDateTime `json:"start"`
	End         gcalDateTime `json:"end"`
}

type gcalDateTime struct {
	DateTime string `json:"dateTime,omitempty"`
	Date     string `json:"date,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

func (s *CalendarService) calendarRequest(method, path, token string, body interface{}) (map[string]interface{}, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(string(b))
	}

	req, err := http.NewRequest(method, calendarAPIURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição Calendar API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil, nil // DELETE success
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode >= 400 {
		msg := ""
		if e, ok := result["error"].(map[string]interface{}); ok {
			msg = fmt.Sprintf("%v", e["message"])
		}
		return nil, fmt.Errorf("Google Calendar API %d: %s", resp.StatusCode, msg)
	}

	return result, nil
}

// CreateEventFromEvent creates a Google Calendar event from an app Event.
// Returns the Google event ID.
func (s *CalendarService) CreateEventFromEvent(profileID int, event models.Event) (string, error) {
	token, err := s.getValidToken(profileID)
	if err != nil {
		return "", err
	}

	var ce gcalEvent
	ce.Summary = event.Title
	ce.Description = event.Description
	ce.Location = event.Location

	if event.AllDay {
		ce.Start = gcalDateTime{Date: event.StartDate.Format("2006-01-02")}
		ce.End = gcalDateTime{Date: event.EndDate.Format("2006-01-02")}
	} else {
		tz := "America/Sao_Paulo"
		ce.Start = gcalDateTime{DateTime: event.StartDate.UTC().Format(time.RFC3339), TimeZone: tz}
		ce.End = gcalDateTime{DateTime: event.EndDate.UTC().Format(time.RFC3339), TimeZone: tz}
	}

	result, err := s.calendarRequest("POST", "/calendars/primary/events", token, ce)
	if err != nil {
		return "", err
	}
	id, _ := result["id"].(string)
	return id, nil
}

// CreateEventFromTask creates a Google Calendar event from a Task (uses due_date as start, +1h as end).
func (s *CalendarService) CreateEventFromTask(profileID int, task models.Task) (string, error) {
	token, err := s.getValidToken(profileID)
	if err != nil {
		return "", err
	}
	if task.DueDate == nil {
		return "", fmt.Errorf("tarefa sem data de vencimento — defina uma data para sincronizar com o Calendar")
	}

	start := task.DueDate.UTC()
	end := start.Add(time.Hour)
	tz := "America/Sao_Paulo"

	ce := gcalEvent{
		Summary:     task.Title,
		Description: task.Description,
		Start:       gcalDateTime{DateTime: start.Format(time.RFC3339), TimeZone: tz},
		End:         gcalDateTime{DateTime: end.Format(time.RFC3339), TimeZone: tz},
	}

	result, err := s.calendarRequest("POST", "/calendars/primary/events", token, ce)
	if err != nil {
		return "", err
	}
	id, _ := result["id"].(string)
	return id, nil
}

// DeleteCalendarEvent removes an event from Google Calendar.
func (s *CalendarService) DeleteCalendarEvent(profileID int, googleEventID string) error {
	token, err := s.getValidToken(profileID)
	if err != nil {
		return err
	}
	_, err = s.calendarRequest("DELETE", "/calendars/primary/events/"+googleEventID, token, nil)
	return err
}
