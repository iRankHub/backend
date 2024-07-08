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
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.Province ELSE ild.Continent END as detail1,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.District ELSE ild.Country END as detail2
FROM Leagues l
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
WHERE l.LeagueID = $1;

-- name: ListLeaguesPaginated :many
SELECT l.*,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.Province ELSE ild.Continent END as detail1,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.District ELSE ild.Country END as detail2
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

-- name: CreateTournamentFormat :one
INSERT INTO TournamentFormats (FormatName, Description)
VALUES ($1, $2)
RETURNING *;

-- name: GetTournamentFormatByID :one
SELECT * FROM TournamentFormats WHERE FormatID = $1;

-- name: ListTournamentFormatsPaginated :many
SELECT * FROM TournamentFormats
ORDER BY FormatID
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentFormatDetails :one
UPDATE TournamentFormats
SET FormatName = $2, Description = $3
WHERE FormatID = $1
RETURNING *;

-- name: DeleteTournamentFormatByID :exec
DELETE FROM TournamentFormats WHERE FormatID = $1;

-- name: CreateTournamentEntry :one
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID, LeagueID)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTournamentByID :one
SELECT t.*, tf.*, l.*,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.Province ELSE ild.Continent END as league_detail1,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.District ELSE ild.Country END as league_detail2
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
WHERE t.TournamentID = $1;

-- name: ListTournamentsPaginated :many
SELECT t.*, tf.*, l.*,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.Province ELSE ild.Continent END as league_detail1,
       CASE WHEN l.LeagueType = 'LEAGUE_TYPE_LOCAL' THEN lld.District ELSE ild.Country END as league_detail2
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
ORDER BY t.TournamentID
LIMIT $1 OFFSET $2;

-- name: UpdateTournamentDetails :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6, LeagueID = $7
WHERE TournamentID = $1
RETURNING *;

-- name: DeleteTournamentByID :exec
DELETE FROM Tournaments WHERE TournamentID = $1;