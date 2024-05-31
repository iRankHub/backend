-- name: GetTeam :one
SELECT * FROM Teams WHERE TeamID = $1;

-- name: GetTeams :many
SELECT * FROM Teams;

-- name: CreateTeam :one
INSERT INTO Teams (Name, SchoolID, TournamentID)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateTeam :one
UPDATE Teams
SET Name = $2, SchoolID = $3, TournamentID = $4
WHERE TeamID = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM Teams WHERE TeamID = $1;