-- name: GetTournament :one
SELECT * FROM Tournaments WHERE TournamentID = $1;

-- name: GetTournaments :many
SELECT * FROM Tournaments;

-- name: CreateTournament :one
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTournament :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6
WHERE TournamentID = $1
RETURNING *;

-- name: DeleteTournament :exec
DELETE FROM Tournaments WHERE TournamentID = $1;