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
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID, LeagueID, CoordinatorID, NumberOfPreliminaryRounds, NumberOfEliminationRounds, JudgesPerDebatePreliminary, JudgesPerDebateElimination, TournamentFee, ImageUrl)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetTournamentByID :one
SELECT t.*, tf.FormatName, tf.Description AS FormatDescription, tf.SpeakersPerTeam,
       l.Name AS LeagueName, l.LeagueType, l.Details AS LeagueDetails,
       u.Name AS CoordinatorName, u.UserID AS CoordinatorID
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
JOIN Users u ON t.CoordinatorID = u.UserID
WHERE t.TournamentID = $1 AND t.deleted_at IS NULL;

-- name: GetActiveTournaments :many
SELECT * FROM Tournaments
WHERE StartDate > CURRENT_TIMESTAMP
  AND deleted_at IS NULL
ORDER BY StartDate;

-- name: ListTournamentsPaginated :many
SELECT
    t.*,
    tf.FormatName,
    l.Name AS LeagueName,
    u.Name AS CoordinatorName,
    u.UserID AS CoordinatorID,
    COUNT(DISTINCT CASE WHEN ti.InviteeRole = 'school' AND ti.Status = 'accepted' THEN ti.InvitationID END) AS AcceptedSchoolsCount,
    COUNT(DISTINCT tm.TeamID) AS TeamsCount
FROM
    Tournaments t
JOIN
    TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN
    Leagues l ON t.LeagueID = l.LeagueID
JOIN
    Users u ON t.CoordinatorID = u.UserID
LEFT JOIN
    TournamentInvitations ti ON t.TournamentID = ti.TournamentID
LEFT JOIN
    Teams tm ON t.TournamentID = tm.TournamentID
WHERE
    t.deleted_at IS NULL
GROUP BY
    t.TournamentID, tf.FormatName, l.Name, u.Name, u.UserID
ORDER BY
    t.StartDate DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentDetails :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6, LeagueID = $7,
    NumberOfPreliminaryRounds = $8, NumberOfEliminationRounds = $9,
    JudgesPerDebatePreliminary = $10, JudgesPerDebateElimination = $11, TournamentFee = $12, ImageUrl = $13
WHERE TournamentID = $1
RETURNING *;

-- name: DeleteTournamentByID :exec
UPDATE Tournaments
SET deleted_at = CURRENT_TIMESTAMP
WHERE TournamentID = $1;

-- name: CreateInvitation :one
INSERT INTO TournamentInvitations (TournamentID, InviteeID, InviteeRole, Status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetInvitationByID :one
SELECT * FROM TournamentInvitations WHERE InvitationID = $1;

-- name: GetInvitationsByTournament :many
SELECT
    ti.InvitationID,
    ti.Status,
    ti.InviteeID,
    CASE
        WHEN ti.InviteeRole = 'school' THEN s.SchoolName
        WHEN ti.InviteeRole = 'volunteer' THEN CONCAT(v.FirstName, ' ', v.LastName)
        WHEN ti.InviteeRole = 'student' THEN CONCAT(st.FirstName, ' ', st.LastName)
    END as InviteeName,
    ti.InviteeRole,
    ti.created_at,
    ti.updated_at
FROM
    TournamentInvitations ti
LEFT JOIN
    Schools s ON ti.InviteeID = s.iDebateSchoolID
LEFT JOIN
    Volunteers v ON ti.InviteeID = v.iDebateVolunteerID
LEFT JOIN
    Students st ON ti.InviteeID = st.iDebateStudentID
WHERE
    ti.TournamentID = $1
ORDER BY
    ti.created_at DESC;

-- name: UpdateInvitationStatus :one
UPDATE TournamentInvitations
SET Status = $2, updated_at = CURRENT_TIMESTAMP
WHERE InvitationID = $1
RETURNING *;

-- name: BulkUpdateInvitationStatus :many
UPDATE TournamentInvitations
SET Status = $2, updated_at = CURRENT_TIMESTAMP
WHERE InvitationID = ANY($1::int[])
RETURNING *;

-- name: DeleteInvitation :exec
DELETE FROM TournamentInvitations WHERE InvitationID = $1;

-- name: UpdateReminderSentAt :one
UPDATE TournamentInvitations
SET ReminderSentAt = $2
WHERE InvitationID = $1
RETURNING *;

-- name: GetInvitationsByUser :many
SELECT
    ti.*,
    CASE
        WHEN ti.InviteeRole = 'school' THEN s.SchoolName
        WHEN ti.InviteeRole = 'volunteer' THEN CONCAT(v.FirstName, ' ', v.LastName)
        WHEN ti.InviteeRole = 'student' THEN CONCAT(st.FirstName, ' ', st.LastName)
    END AS InviteeName
FROM TournamentInvitations ti
LEFT JOIN Schools s ON ti.InviteeRole = 'school' AND ti.InviteeID = s.iDebateSchoolID
LEFT JOIN Volunteers v ON ti.InviteeRole = 'volunteer' AND ti.InviteeID = v.iDebateVolunteerID
LEFT JOIN Students st ON ti.InviteeRole = 'student' AND ti.InviteeID = st.iDebateStudentID
WHERE
    (ti.InviteeRole = 'school' AND s.ContactPersonID = $1) OR
    (ti.InviteeRole = 'volunteer' AND v.UserID = $1) OR
    (ti.InviteeRole = 'student' AND st.UserID = $1)
ORDER BY ti.created_at DESC;

-- name: GetPendingInvitations :many
SELECT
    ti.*,
    CASE
        WHEN ti.InviteeRole = 'school' THEN s.SchoolName
        WHEN ti.InviteeRole = 'volunteer' THEN CONCAT(v.FirstName, ' ', v.LastName)
        WHEN ti.InviteeRole = 'student' THEN CONCAT(st.FirstName, ' ', st.LastName)
    END as InviteeName,
    CASE
        WHEN ti.InviteeRole = 'school' THEN s.ContactEmail
        WHEN ti.InviteeRole = 'volunteer' THEN u.Email
        WHEN ti.InviteeRole = 'student' THEN st.Email
    END as InviteeEmail,
    t.Name as TournamentName,
    t.StartDate as TournamentStartDate,
    t.EndDate as TournamentEndDate,
    t.Location as TournamentLocation
FROM
    TournamentInvitations ti
JOIN
    Tournaments t ON ti.TournamentID = t.TournamentID
LEFT JOIN
    Schools s ON ti.InviteeRole = 'school' AND ti.InviteeID = s.iDebateSchoolID
LEFT JOIN
    Volunteers v ON ti.InviteeRole = 'volunteer' AND ti.InviteeID = v.iDebateVolunteerID
LEFT JOIN
    Students st ON ti.InviteeRole = 'student' AND ti.InviteeID = st.iDebateStudentID
LEFT JOIN
    Users u ON (ti.InviteeRole = 'volunteer' AND v.UserID = u.UserID)
WHERE
    ti.TournamentID = $1 AND ti.Status = 'pending'
ORDER BY
    ti.created_at DESC;

-- name: GetTournamentRegistrations :many
SELECT
    DATE(ti.updated_at) AS registration_date,
    COUNT(*) AS registration_count
FROM
    TournamentInvitations ti
WHERE
    ti.Status = 'accepted'
GROUP BY
    DATE(ti.updated_at)
ORDER BY
    registration_date DESC
LIMIT 30; -- Limiting to last 30 days, adjust as needed

-- name: GetTournamentStats :one
WITH CurrentStats AS (
    SELECT
        COUNT(*) AS total_tournaments,
        COUNT(CASE WHEN StartDate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days' THEN 1 END) AS upcoming_tournaments
    FROM Tournaments
    WHERE deleted_at IS NULL
),
HistoricalStats AS (
    SELECT yesterday_total_count, yesterday_upcoming_count
    FROM Tournaments
    WHERE TournamentID = (SELECT MIN(TournamentID) FROM Tournaments)
)
SELECT
    cs.total_tournaments,
    cs.upcoming_tournaments,
    hs.yesterday_total_count,
    hs.yesterday_upcoming_count
FROM CurrentStats cs, HistoricalStats hs;

