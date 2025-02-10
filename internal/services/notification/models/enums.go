package models

// Category represents the main classification of notifications
type Category string

const (
	AuthCategory       Category = "auth"
	UserCategory       Category = "user"
	TournamentCategory Category = "tournament"
	DebateCategory     Category = "debate"
	ReportCategory     Category = "report"
)

// Type represents specific notification types within each category
type Type string

const (
	// Auth notifications
	AccountCreation Type = "account_creation"
	AccountApproval Type = "account_approval"
	PasswordReset   Type = "password_reset"
	SecurityAlert   Type = "security_alert"
	TwoFactorAuth   Type = "two_factor_auth"

	// User notifications
	ProfileUpdate  Type = "profile_update"
	RoleAssignment Type = "role_assignment"
	StatusChange   Type = "status_change"

	// Tournament notifications
	TournamentInvite       Type = "tournament_invite"
	TournamentRegistration Type = "tournament_registration"
	TournamentSchedule     Type = "tournament_schedule"
	TournamentPayment      Type = "tournament_payment"
	CoordinatorAssignment  Type = "coordinator_assignment"

	// Debate notifications
	RoundAssignment  Type = "round_assignment"
	JudgeAssignment  Type = "judge_assignment"
	BallotSubmission Type = "ballot_submission"
	DebateResults    Type = "debate_results"
	RoomChange       Type = "room_change"

	// Report notifications
	ReportGeneration  Type = "report_generation"
	PerformanceReport Type = "performance_report"
	AnalyticsReport   Type = "analytics_report"
	AuditReport       Type = "audit_report"
)

// DeliveryMethod defines how the notification should be delivered
type DeliveryMethod string

const (
	EmailDelivery DeliveryMethod = "email"
	InAppDelivery DeliveryMethod = "in_app"
	PushDelivery  DeliveryMethod = "push"
)

// Priority defines the urgency level of notifications
type Priority string

const (
	LowPriority    Priority = "low"
	MediumPriority Priority = "medium"
	HighPriority   Priority = "high"
	UrgentPriority Priority = "urgent"
)

// Status represents the current state of a notification
type Status string

const (
	StatusPending   Status = "pending"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
	StatusExpired   Status = "expired"
)

// UserRole defines the role of the notification recipient
type UserRole string

const (
	AdminRole       UserRole = "admin"
	VolunteerRole   UserRole = "volunteer"
	SchoolRole      UserRole = "school"
	StudentRole     UserRole = "student"
	UnspecifiedRole UserRole = "unspecified"
)

// ActionType defines the types of actions that can be taken on notifications
type ActionType string

const (
	ActionView     ActionType = "view"
	ActionAccept   ActionType = "accept"
	ActionReject   ActionType = "reject"
	ActionSubmit   ActionType = "submit"
	ActionDownload ActionType = "download"
)
