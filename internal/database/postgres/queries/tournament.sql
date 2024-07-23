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
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID, LeagueID, NumberOfPreliminaryRounds, NumberOfEliminationRounds, JudgesPerDebatePreliminary, JudgesPerDebateElimination, TournamentFee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetTournamentByID :one
SELECT t.*, tf.*, l.*
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN TournamentCoordinators tc ON t.TournamentID = tc.TournamentID
WHERE t.TournamentID = $1 AND t.deleted_at IS NULL;

-- name: ListTournamentsPaginated :many
SELECT t.*, tf.*, l.*
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN TournamentCoordinators tc ON t.TournamentID = tc.TournamentID
WHERE t.deleted_at IS NULL
ORDER BY t.TournamentID
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentDetails :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6, LeagueID = $7, NumberOfPreliminaryRounds = $8, NumberOfEliminationRounds = $9, JudgesPerDebatePreliminary = $10, JudgesPerDebateElimination = $11, TournamentFee = $12
WHERE TournamentID = $1
RETURNING *;

-- name: DeleteTournamentByID :exec
UPDATE Tournaments
SET deleted_at = CURRENT_TIMESTAMP
WHERE TournamentID = $1;