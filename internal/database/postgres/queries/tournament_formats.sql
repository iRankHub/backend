-- name: GetTournamentFormat :one
SELECT * FROM TournamentFormats WHERE FormatID = $1;

-- name: GetTournamentFormats :many
SELECT * FROM TournamentFormats;

-- name: CreateTournamentFormat :one
INSERT INTO TournamentFormats (FormatName, Description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTournamentFormat :one
UPDATE TournamentFormats
SET FormatName = $2, Description = $3
WHERE FormatID = $1
RETURNING *;

-- name: DeleteTournamentFormat :exec
DELETE FROM TournamentFormats WHERE FormatID = $1;