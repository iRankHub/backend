-- League Queries
-- name: CreateLeague :one
INSERT INTO Leagues (Name, LeagueType)
VALUES ($1, $2)
RETURNING *;

-- name: CreateLocalLeagueDetails :exec
INSERT INTO LocalLeagueDetails (LeagueID, Province, District)
VALUES ($1, $2, $3);

-- name: CreateInternationalLeagueDetails :exec
INSERT INTO InternationalLeagueDetails (LeagueID, Continent, Country)
VALUES ($1, $2, $3);

-- name: GetLeagueByID :one
SELECT l.*,
       COALESCE(lld.Province, ild.Continent) AS detail1,
       COALESCE(lld.District, ild.Country) AS detail2
FROM Leagues l
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
WHERE l.LeagueID = $1;

-- name: ListLeaguesPaginated :many
SELECT l.*,
       COALESCE(lld.Province, ild.Continent) AS detail1,
       COALESCE(lld.District, ild.Country) AS detail2
FROM Leagues l
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
ORDER BY l.LeagueID
LIMIT $1 OFFSET $2;

-- name: UpdateLeagueDetails :one
UPDATE Leagues
SET Name = $2, LeagueType = $3
WHERE LeagueID = $1
RETURNING *;

-- name: UpdateLocalLeagueDetailsInfo :exec
UPDATE LocalLeagueDetails
SET Province = $2, District = $3
WHERE LeagueID = $1;

-- name: UpdateInternationalLeagueDetailsInfo :exec
UPDATE InternationalLeagueDetails
SET Continent = $2, Country = $3
WHERE LeagueID = $1;

-- name: DeleteLeagueByID :exec
DELETE FROM Leagues WHERE LeagueID = $1;

-- Tournament Format Queries
-- name: CreateTournamentFormat :one
INSERT INTO TournamentFormats (FormatName, Description, SpeakersPerTeam)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTournamentFormatByID :one
SELECT * FROM TournamentFormats WHERE FormatID = $1;

-- name: ListTournamentFormatsPaginated :many
SELECT * FROM TournamentFormats
ORDER BY FormatID
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentFormatDetails :one
UPDATE TournamentFormats
SET FormatName = $2, Description = $3, SpeakersPerTeam = $4
WHERE FormatID = $1
RETURNING *;

-- name: DeleteTournamentFormatByID :exec
DELETE FROM TournamentFormats WHERE FormatID = $1;

-- Tournament Queries
-- name: CreateTournamentEntry :one
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID, LeagueID, NumberOfPreliminaryRounds, NumberOfEliminationRounds, JudgesPerDebatePreliminary, JudgesPerDebateElimination, TournamentFee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetTournamentByID :one
SELECT t.*, tf.*, l.*,
       COALESCE(lld.Province, ild.Continent) AS league_detail1,
       COALESCE(lld.District, ild.Country) AS league_detail2,
       tc.VolunteerID as CoordinatorID
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
LEFT JOIN TournamentCoordinators tc ON t.TournamentID = tc.TournamentID
WHERE t.TournamentID = $1;

-- name: ListTournamentsPaginated :many
SELECT t.*, tf.*, l.*,
       COALESCE(lld.Province, ild.Continent) AS league_detail1,
       COALESCE(lld.District, ild.Country) AS league_detail2,
       tc.VolunteerID as CoordinatorID
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
LEFT JOIN TournamentCoordinators tc ON t.TournamentID = tc.TournamentID
ORDER BY t.TournamentID
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentDetails :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6, LeagueID = $7, NumberOfPreliminaryRounds = $8, NumberOfEliminationRounds = $9, JudgesPerDebatePreliminary = $10, JudgesPerDebateElimination = $11, TournamentFee = $12
WHERE TournamentID = $1
RETURNING *;

-- name: DeleteTournamentByID :exec
DELETE FROM Tournaments WHERE TournamentID = $1;