package services

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"personal-cockpit/models"
)

type ProfileService struct {
	db       *sql.DB
	settings *SettingsService
}

func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{
		db:       db,
		settings: NewSettingsService(db),
	}
}

// CreateProfile creates a new profile. The first profile ever created gets role=admin.
func (s *ProfileService) CreateProfile(input models.CreateProfileInput) (*models.ProfileSummary, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}

	role := input.Role
	if role == "" {
		count, _ := s.countProfiles()
		if count == 0 {
			role = "admin"
		} else {
			role = "member"
		}
	}

	var pinHash string
	if input.Pin != "" {
		pinHash = hashPin(input.Pin)
	}

	result, err := s.db.Exec(`
		INSERT INTO profiles (name, email, role, avatar_url, pin_hash)
		VALUES (?, ?, ?, ?, ?)
	`, input.Name, input.Email, role, input.AvatarURL, pinHash)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar perfil: %w", err)
	}

	id, _ := result.LastInsertId()
	return s.GetProfileSummaryByID(int(id))
}

// GetAllProfiles returns all active profiles as summaries (no tokens).
func (s *ProfileService) GetAllProfiles() ([]models.ProfileSummary, error) {
	rows, err := s.db.Query(`
		SELECT id, name, email, role, avatar_url, google_calendar_id,
		       CASE WHEN google_access_token != '' AND google_access_token IS NOT NULL THEN 1 ELSE 0 END,
		       CASE WHEN pin_hash != '' AND pin_hash IS NOT NULL THEN 1 ELSE 0 END,
		       created_at
		FROM profiles
		WHERE is_active = 1
		ORDER BY role DESC, created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar perfis: %w", err)
	}
	defer rows.Close()

	var profiles []models.ProfileSummary
	for rows.Next() {
		var p models.ProfileSummary
		err := rows.Scan(
			&p.ID, &p.Name, &p.Email, &p.Role, &p.AvatarURL,
			&p.GoogleCalendarID, &p.HasGoogleToken, &p.HasPin, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// GetProfileSummaryByID returns a profile summary by ID.
func (s *ProfileService) GetProfileSummaryByID(id int) (*models.ProfileSummary, error) {
	var p models.ProfileSummary
	err := s.db.QueryRow(`
		SELECT id, name, email, role, avatar_url, google_calendar_id,
		       CASE WHEN google_access_token != '' AND google_access_token IS NOT NULL THEN 1 ELSE 0 END,
		       CASE WHEN pin_hash != '' AND pin_hash IS NOT NULL THEN 1 ELSE 0 END,
		       created_at
		FROM profiles WHERE id = ?
	`, id).Scan(
		&p.ID, &p.Name, &p.Email, &p.Role, &p.AvatarURL,
		&p.GoogleCalendarID, &p.HasGoogleToken, &p.HasPin, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("perfil não encontrado")
	}
	return &p, err
}

// GetProfileByID returns the full profile (including token) for internal use.
func (s *ProfileService) GetProfileByID(id int) (*models.Profile, error) {
	var p models.Profile
	var email, googleID, avatarURL, token, calID, pinHash sql.NullString
	err := s.db.QueryRow(`
		SELECT id, name, email, google_id, role, avatar_url,
		       google_access_token, google_calendar_id, pin_hash, is_active, created_at
		FROM profiles WHERE id = ?
	`, id).Scan(
		&p.ID, &p.Name, &email, &googleID, &p.Role, &avatarURL,
		&token, &calID, &pinHash, &p.IsActive, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("perfil não encontrado")
	}
	if err != nil {
		return nil, err
	}
	p.Email = email.String
	p.GoogleID = googleID.String
	p.AvatarURL = avatarURL.String
	p.GoogleAccessToken = token.String
	p.GoogleCalendarID = calID.String
	p.PinHash = pinHash.String
	return &p, nil
}

// UpdateProfile updates name, email, avatar and role.
func (s *ProfileService) UpdateProfile(id int, input models.CreateProfileInput) error {
	_, err := s.db.Exec(`
		UPDATE profiles SET name = ?, email = ?, avatar_url = ?, role = ?
		WHERE id = ?
	`, input.Name, input.Email, input.AvatarURL, input.Role, id)
	return err
}

// DeleteProfile soft-deletes a profile.
func (s *ProfileService) DeleteProfile(id int) error {
	// Prevent deleting the last admin
	var adminCount int
	s.db.QueryRow(`SELECT COUNT(*) FROM profiles WHERE role = 'admin' AND is_active = 1`).Scan(&adminCount)

	var role string
	s.db.QueryRow(`SELECT role FROM profiles WHERE id = ?`, id).Scan(&role)
	if role == "admin" && adminCount <= 1 {
		return fmt.Errorf("não é possível remover o único administrador")
	}

	_, err := s.db.Exec(`UPDATE profiles SET is_active = 0 WHERE id = ?`, id)
	return err
}

// VerifyPin checks if the given PIN matches the stored hash.
func (s *ProfileService) VerifyPin(id int, pin string) (bool, error) {
	var pinHash sql.NullString
	err := s.db.QueryRow(`SELECT pin_hash FROM profiles WHERE id = ?`, id).Scan(&pinHash)
	if err != nil {
		return false, err
	}
	if !pinHash.Valid || pinHash.String == "" {
		return true, nil // no PIN set — always passes
	}
	return hashPin(pin) == pinHash.String, nil
}

// SetPin updates the PIN for a profile (empty pin removes it).
func (s *ProfileService) SetPin(id int, pin string) error {
	var h string
	if pin != "" {
		h = hashPin(pin)
	}
	_, err := s.db.Exec(`UPDATE profiles SET pin_hash = ? WHERE id = ?`, h, id)
	return err
}

// SetCurrentProfile persists the active profile ID in settings.
func (s *ProfileService) SetCurrentProfile(id int) error {
	return s.settings.Set("current_profile_id", fmt.Sprintf("%d", id))
}

// GetCurrentProfileID returns the last active profile ID (0 if none).
func (s *ProfileService) GetCurrentProfileID() int {
	val, _ := s.settings.Get("current_profile_id")
	var id int
	fmt.Sscanf(val, "%d", &id)
	return id
}

// SaveGoogleToken stores the OAuth2 token JSON for a profile.
func (s *ProfileService) SaveGoogleToken(profileID int, tokenJSON, email, calendarID string) error {
	_, err := s.db.Exec(`
		UPDATE profiles SET google_access_token = ?, email = ?, google_calendar_id = ?
		WHERE id = ?
	`, tokenJSON, email, calendarID, profileID)
	return err
}

// ClearGoogleToken removes Google tokens from a profile.
func (s *ProfileService) ClearGoogleToken(profileID int) error {
	_, err := s.db.Exec(`
		UPDATE profiles SET google_access_token = NULL WHERE id = ?
	`, profileID)
	return err
}

// HasAnyProfile returns true if at least one profile exists.
func (s *ProfileService) HasAnyProfile() (bool, error) {
	count, err := s.countProfiles()
	return count > 0, err
}

func (s *ProfileService) countProfiles() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM profiles WHERE is_active = 1`).Scan(&count)
	return count, err
}

func hashPin(pin string) string {
	h := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", h)
}
