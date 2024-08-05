-- League Queries
-- name: CreateLeague :one
INSERT INTO Leagues (Name, LeagueType, Details)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLeagueByID :one
SELECT * FROM Leagues
WHERE LeagueID = $1 AND deleted_at IS NULL;

-- name: ListLeaguesPaginated :many
SELECT * FROM Leagues
WHERE deleted_at IS NULL
ORDER BY LeagueID
LIMIT $1 OFFSET $2;

-- name: UpdateLeague :one
UPDATE Leagues
SET Name = $2, LeagueType = $3, Details = $4
WHERE LeagueID = $1
RETURNING *;

-- name: DeleteLeagueByID :exec
UPDATE Leagues
SET deleted_at = CURRENT_TIMESTAMP
WHERE LeagueID = $1;

-- Tournament Format Queries
-- name: CreateTournamentFormat :one
INSERT INTO TournamentFormats (FormatName, Description, SpeakersPerTeam)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTournamentFormatByID :one
SELECT * FROM TournamentFormats
WHERE FormatID = $1 AND deleted_at IS NULL;

-- name: ListTournamentFormatsPaginated :many
SELECT * FROM TournamentFormats
WHERE deleted_at IS NULL
ORDER BY FormatID
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentFormatDetails :one
UPDATE TournamentFormats
SET FormatName = $2, Description = $3, SpeakersPerTeam = $4
WHERE FormatID = $1
RETURNING *;

-- name: DeleteTournamentFormatByID :exec
UPDATE TournamentFormats
SET deleted_at = CURRENT_TIMESTAMP
WHERE FormatID = $1;

-- Tournament Queries
-- name: CreateTournamentEntry :one
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID, LeagueID, CoordinatorID, NumberOfPreliminaryRounds, NumberOfEliminationRounds, JudgesPerDebatePreliminary, JudgesPerDebateElimination, TournamentFee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetTournamentByID :one
SELECT t.*, tf.FormatName, tf.Description AS FormatDescription, tf.SpeakersPerTeam,
       l.Name AS LeagueName, l.LeagueType, l.Details AS LeagueDetails
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
WHERE t.TournamentID = $1 AND t.deleted_at IS NULL;

-- name: GetActiveTournaments :many
SELECT * FROM Tournaments
WHERE StartDate > CURRENT_TIMESTAMP
  AND deleted_at IS NULL
ORDER BY StartDate;

-- name: ListTournamentsPaginated :many
SELECT t.*, tf.FormatName, l.Name AS LeagueName
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
WHERE t.deleted_at IS NULL
ORDER BY t.StartDate DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentDetails :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6, LeagueID = $7,
    NumberOfPreliminaryRounds = $8, NumberOfEliminationRounds = $9,
    JudgesPerDebatePreliminary = $10, JudgesPerDebateElimination = $11, TournamentFee = $12
WHERE TournamentID = $1
RETURNING *;

-- name: DeleteTournamentByID :exec
UPDATE Tournaments
SET deleted_at = CURRENT_TIMESTAMP
WHERE TournamentID = $1;

-- name: CreateInvitation :one
INSERT INTO TournamentInvitations (TournamentID, SchoolID, VolunteerID, StudentID, UserID, Status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInvitationByID :one
SELECT * FROM TournamentInvitations
WHERE InvitationID = $1;

-- name: GetInvitationsByUserID :many
SELECT * FROM TournamentInvitations
WHERE UserID = $1;

-- name: UpdateInvitationStatus :exec
UPDATE TournamentInvitations
SET Status = $2, RespondedAt = CURRENT_TIMESTAMP
WHERE InvitationID = $1;

-- name: UpdateInvitationStatusWithUserCheck :exec
UPDATE TournamentInvitations
SET Status = $2, RespondedAt = CURRENT_TIMESTAMP
WHERE InvitationID = $1 AND UserID = $3;

-- name: GetPendingInvitations :many
SELECT ti.*,
       s.SchoolName, s.ContactEmail, s.SchoolEmail,
       v.VolunteerID, v.FirstName as VolunteerFirstName, v.LastName as VolunteerLastName, u.Email as VolunteerEmail,
       st.StudentID, st.Email as StudentEmail, st.FirstName as StudentFirstName, st.LastName as StudentLastName
FROM TournamentInvitations ti
LEFT JOIN Schools s ON ti.SchoolID = s.SchoolID
LEFT JOIN Volunteers v ON ti.VolunteerID = v.VolunteerID
LEFT JOIN Users u ON ti.UserID = u.UserID
LEFT JOIN Students st ON ti.StudentID = st.StudentID
WHERE ti.Status = 'pending'
  AND ti.TournamentID = $1;

-- name: RegisterTeam :one
INSERT INTO Teams (Name, SchoolID, TournamentID, InvitationID)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: AddTeamMember :exec
INSERT INTO TeamMembers (TeamID, StudentID)
VALUES ($1, $2);

-- name: GetInvitationStatus :one
SELECT i.*,
       json_agg(json_build_object('team_id', t.TeamID, 'team_name', t.Name, 'number_of_speakers', COUNT(tm.StudentID))) as registered_teams
FROM TournamentInvitations i
LEFT JOIN Teams t ON i.InvitationID = t.InvitationID
LEFT JOIN TeamMembers tm ON t.TeamID = tm.TeamID
WHERE i.InvitationID = $1
GROUP BY i.InvitationID;

-- name: GetTeamsByInvitation :many
SELECT t.*, COUNT(tm.StudentID) as number_of_speakers
FROM Teams t
LEFT JOIN TeamMembers tm ON t.TeamID = tm.TeamID
WHERE t.InvitationID = $1
GROUP BY t.TeamID;

-- name: UpdateReminderSentAt :one
UPDATE TournamentInvitations
SET ReminderSentAt = $2
WHERE InvitationID = $1
RETURNING *;