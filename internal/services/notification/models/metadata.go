package models

import "time"

// AuthMetadata contains metadata specific to authentication notifications
type AuthMetadata struct {
	IPAddress    string    `json:"ip_address,omitempty"`
	DeviceInfo   string    `json:"device_info,omitempty"`
	Location     string    `json:"location,omitempty"`
	AttemptCount int       `json:"attempt_count,omitempty"`
	LastAttempt  time.Time `json:"last_attempt,omitempty"`
}

// TournamentMetadata contains metadata specific to tournament notifications
type TournamentMetadata struct {
	TournamentID   string    `json:"tournament_id"`
	TournamentName string    `json:"tournament_name"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Location       string    `json:"location"`
	Format         string    `json:"format"`
	League         string    `json:"league"`
	Fee            float64   `json:"fee,omitempty"`
	Currency       string    `json:"currency,omitempty"`
	Coordinator    string    `json:"coordinator,omitempty"`
}

// DebateMetadata contains metadata specific to debate notifications
type DebateMetadata struct {
	DebateID      string    `json:"debate_id"`
	TournamentID  string    `json:"tournament_id"`
	RoundNumber   int       `json:"round_number"`
	IsElimination bool      `json:"is_elimination"`
	Room          string    `json:"room"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Team1         string    `json:"team1"`
	Team2         string    `json:"team2"`
	JudgePanel    []string  `json:"judge_panel"`
	HeadJudge     string    `json:"head_judge,omitempty"`
	Motion        string    `json:"motion,omitempty"`
}

// UserMetadata contains metadata specific to user notifications
type UserMetadata struct {
	Changes        map[string]string `json:"changes,omitempty"`
	PreviousRole   string            `json:"previous_role,omitempty"`
	NewRole        string            `json:"new_role,omitempty"`
	Reason         string            `json:"reason,omitempty"`
	ApprovedBy     string            `json:"approved_by,omitempty"`
	ApprovedAt     time.Time         `json:"approved_at,omitempty"`
	ExpirationDate *time.Time        `json:"expiration_date,omitempty"`
}

// ReportMetadata contains metadata specific to report notifications
type ReportMetadata struct {
	ReportID    string            `json:"report_id"`
	ReportType  string            `json:"report_type"`
	Period      string            `json:"period"`
	GeneratedAt time.Time         `json:"generated_at"`
	GeneratedBy string            `json:"generated_by"`
	Size        string            `json:"size,omitempty"`
	DownloadURL string            `json:"download_url,omitempty"`
	Summary     map[string]string `json:"summary,omitempty"`
	KeyMetrics  map[string]any    `json:"key_metrics,omitempty"`
	ExpiresAt   time.Time         `json:"expires_at"`
	FileSize    string            `json:"file_size,omitempty"`
}

// MetadataStore represents metadata that will be stored in the database
type MetadataStore struct {
	ID             string     `db:"id"`
	NotificationID string     `db:"notification_id"`
	UserID         string     `db:"user_id"`
	Category       Category   `db:"category"`
	Type           Type       `db:"type"`
	Status         Status     `db:"status"`
	IsRead         bool       `db:"is_read"`
	ReadAt         *time.Time `db:"read_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	ExpiresAt      time.Time  `db:"expires_at"`
}
