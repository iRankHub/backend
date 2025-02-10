-- Migration file for all indexes
-- Order of implementation: Core user tables first, then tournament/debate related,
-- followed by supporting features

-- ==========================================
-- Core User Management Indexes
-- These indexes support basic user operations, authentication, and profile management
-- ==========================================

-- Basic user lookup and authentication
CREATE INDEX IF NOT EXISTS idx_users_email ON Users(Email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_status ON Users(Status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_email_password ON Users(Email, Password) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_reset_token ON Users(reset_token) WHERE reset_token IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_role_status ON Users(UserRole, Status, deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_reset_token_expires ON Users(reset_token, reset_token_expires) WHERE reset_token IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_two_factor ON Users(UserID, two_factor_enabled) WHERE two_factor_enabled = true;
CREATE INDEX IF NOT EXISTS idx_users_login_attempts ON Users(UserID, failed_login_attempts) WHERE failed_login_attempts > 0;
CREATE INDEX IF NOT EXISTS idx_users_status_created ON Users(Status, created_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deactivated ON Users(DeactivatedAt) WHERE DeactivatedAt IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_timestamps ON Users(created_at, updated_at, last_login_attempt);
CREATE INDEX IF NOT EXISTS idx_users_active ON Users(UserID) WHERE deleted_at IS NULL AND DeactivatedAt IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_pending_verification ON Users(Status, VerificationStatus) WHERE Status = 'pending' AND VerificationStatus = false;

-- User Profile Management
CREATE INDEX IF NOT EXISTS idx_user_profiles_composite ON UserProfiles(UserID, UserRole);
CREATE INDEX IF NOT EXISTS idx_user_profiles_verification ON UserProfiles(UserID, VerificationStatus);

-- WebAuthn and Security
CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_lookup ON WebAuthnCredentials(UserID, CredentialID);
CREATE INDEX IF NOT EXISTS idx_webauthn_session_data ON WebAuthnSessionData(UserID, SessionData);

-- ==========================================
-- Tournament Management Indexes
-- These indexes optimize tournament operations and management
-- ==========================================

-- Basic tournament lookups
CREATE INDEX IF NOT EXISTS idx_tournaments_coordinator_id ON Tournaments(CoordinatorID);
CREATE INDEX IF NOT EXISTS idx_tournaments_startdate ON Tournaments(StartDate) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tournaments_dates ON Tournaments(StartDate, EndDate) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tournaments_coordinator ON Tournaments(CoordinatorID, deleted_at);
CREATE INDEX IF NOT EXISTS idx_tournaments_league_format ON Tournaments(LeagueID, FormatID);
CREATE INDEX IF NOT EXISTS idx_tournaments_active ON Tournaments(StartDate, deleted_at);
CREATE INDEX IF NOT EXISTS idx_tournaments_upcoming ON Tournaments(StartDate, deleted_at);
CREATE INDEX IF NOT EXISTS idx_tournaments_stats ON Tournaments(TournamentID, StartDate, deleted_at);
CREATE INDEX IF NOT EXISTS idx_tournaments_timestamps ON Tournaments(created_at, updated_at);

-- Tournament Format and League
CREATE INDEX IF NOT EXISTS idx_tournament_formats_active ON TournamentFormats(FormatID) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_leagues_active ON Leagues(LeagueID) WHERE deleted_at IS NULL;

-- ==========================================
-- Debate Management Indexes
-- These indexes support debate operations, scoring, and judging
-- ==========================================

-- Debate Operations
CREATE INDEX IF NOT EXISTS idx_debates_tournament_round ON Debates(TournamentID, RoundNumber, IsEliminationRound);
CREATE INDEX IF NOT EXISTS idx_debates_teams ON Debates(Team1ID, Team2ID);
CREATE INDEX IF NOT EXISTS idx_debates_elimination_tournament ON Debates(IsEliminationRound, TournamentID);
CREATE INDEX IF NOT EXISTS idx_debates_tournament_round_room ON Debates(TournamentID, RoundNumber, RoomID);
CREATE INDEX IF NOT EXISTS idx_debates_elimination_stats ON Debates(TournamentID, IsEliminationRound, RoundNumber);
CREATE INDEX IF NOT EXISTS idx_debates_roundid ON Debates(RoundID);
CREATE INDEX IF NOT EXISTS idx_debates_tournamentid ON Debates(TournamentID);
CREATE INDEX IF NOT EXISTS idx_debates_roomid ON Debates(RoomID);

-- Scoring and Ballots
CREATE INDEX IF NOT EXISTS idx_ballots_debate_status ON Ballots(DebateID, RecordingStatus);
CREATE INDEX IF NOT EXISTS idx_ballots_judge ON Ballots(JudgeID, DebateID);
CREATE INDEX IF NOT EXISTS idx_ballots_verdict ON Ballots(Verdict, RecordingStatus);
CREATE INDEX IF NOT EXISTS idx_ballots_debate_judge_status ON Ballots(DebateID, JudgeID, RecordingStatus);

-- Judge Assignments and Feedback
CREATE INDEX IF NOT EXISTS idx_judge_assignments_tournament_round ON JudgeAssignments(TournamentID, RoundNumber, IsElimination);
CREATE INDEX IF NOT EXISTS idx_judge_assignments_debate_judge ON JudgeAssignments(DebateID, JudgeID);
CREATE INDEX IF NOT EXISTS idx_judge_assignments_head_judge ON JudgeAssignments(DebateID, IsHeadJudge);
CREATE INDEX IF NOT EXISTS idx_debatejudges_debateid ON DebateJudges(DebateID);
CREATE INDEX IF NOT EXISTS idx_debatejudges_judgeid ON DebateJudges(JudgeID);

-- ==========================================
-- Team and Student Management Indexes
-- These indexes support team management and student participation
-- ==========================================

-- Team Management
CREATE INDEX IF NOT EXISTS idx_team_members_teamid ON TeamMembers(TeamID);
CREATE INDEX IF NOT EXISTS idx_teams_tournament ON Teams(TournamentID, TeamID);
CREATE INDEX IF NOT EXISTS idx_team_members_student ON TeamMembers(StudentID, TeamID);
CREATE INDEX IF NOT EXISTS idx_team_scores_debate ON TeamScores(DebateID, TeamID);
CREATE INDEX IF NOT EXISTS idx_team_scores_tournament ON TeamScores(TeamID, IsElimination);
CREATE INDEX IF NOT EXISTS idx_team_scores_tournament_elimination ON TeamScores(TeamID, DebateID, IsElimination);

-- Speaker Scores
CREATE INDEX IF NOT EXISTS idx_speaker_scores_ballot ON SpeakerScores(BallotID, SpeakerID);
CREATE INDEX IF NOT EXISTS idx_speaker_scores_points ON SpeakerScores(SpeakerPoints);
CREATE INDEX IF NOT EXISTS idx_speaker_scores_rank ON SpeakerScores(SpeakerRank);
CREATE INDEX IF NOT EXISTS idx_speaker_scores_stats ON SpeakerScores(SpeakerID, SpeakerPoints, SpeakerRank);

-- ==========================================
-- Participant Management Indexes
-- These indexes support student, volunteer, and school management
-- ==========================================

-- Student Management
CREATE INDEX IF NOT EXISTS idx_students_email ON Students(Email);
CREATE INDEX IF NOT EXISTS idx_students_schoolid ON Students(SchoolID);
CREATE INDEX IF NOT EXISTS idx_students_userid ON Students(UserID);
CREATE INDEX IF NOT EXISTS idx_students_composite ON Students(UserID, SchoolID, iDebateStudentID);
CREATE INDEX IF NOT EXISTS idx_students_school_stats ON Students(SchoolID, UserID);

-- Volunteer Management
CREATE INDEX IF NOT EXISTS idx_volunteers_userid ON Volunteers(UserID);
CREATE INDEX IF NOT EXISTS idx_volunteers_composite ON Volunteers(UserID, Role, iDebateVolunteerID);
CREATE INDEX IF NOT EXISTS idx_volunteers_stats ON Volunteers(UserID, Role, HasInternship, IsEnrolledInUniversity);

-- School Management
CREATE INDEX IF NOT EXISTS idx_schools_contactpersonid ON Schools(ContactPersonID);
CREATE INDEX IF NOT EXISTS idx_schools_contactemail ON Schools(ContactEmail);
CREATE INDEX IF NOT EXISTS idx_schools_composite ON Schools(ContactPersonID, iDebateSchoolID, SchoolType);

-- ==========================================
-- Invitation and Notification Indexes
-- These indexes support the invitation system and notifications
-- ==========================================

-- Invitation Management
CREATE INDEX IF NOT EXISTS idx_tournament_invitations_status ON TournamentInvitations(Status);
CREATE INDEX IF NOT EXISTS idx_tournament_invitations_tournament_id ON TournamentInvitations(TournamentID);
CREATE INDEX IF NOT EXISTS idx_invitations_composite ON TournamentInvitations(TournamentID, InviteeID, Status);
CREATE INDEX IF NOT EXISTS idx_invitations_role_status ON TournamentInvitations(InviteeRole, Status);
CREATE INDEX IF NOT EXISTS idx_invitations_reminder ON TournamentInvitations(ReminderSentAt) WHERE ReminderSentAt IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_invitations_pending ON TournamentInvitations(Status, created_at) WHERE Status = 'pending';
CREATE INDEX IF NOT EXISTS idx_invitations_registration ON TournamentInvitations(Status, updated_at) WHERE Status = 'accepted';
CREATE INDEX IF NOT EXISTS idx_invitations_timestamps ON TournamentInvitations(created_at, updated_at, ReminderSentAt);
CREATE INDEX IF NOT EXISTS idx_tournament_invitations_status_tracking ON TournamentInvitations(TournamentID, Status, InviteeRole);
CREATE INDEX IF NOT EXISTS idx_invitations_pending_reminders ON TournamentInvitations(TournamentID) WHERE Status = 'pending' AND ReminderSentAt IS NULL;

-- Notification Management
CREATE INDEX IF NOT EXISTS idx_notification_metadata_user_id ON notification_metadata(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_metadata_notification_id ON notification_metadata(notification_id);
CREATE INDEX IF NOT EXISTS idx_notification_metadata_category ON notification_metadata(category);
CREATE INDEX IF NOT EXISTS idx_notification_metadata_status ON notification_metadata(status);
CREATE INDEX IF NOT EXISTS idx_notification_metadata_expires_at ON notification_metadata(expires_at);
CREATE INDEX IF NOT EXISTS idx_notification_metadata_is_read ON notification_metadata(is_read);