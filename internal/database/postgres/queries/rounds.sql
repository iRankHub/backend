-- name: GetRound :one
SELECT * FROM Rounds WHERE RoundID = $1;

-- name: GetRounds :many
SELECT * FROM Rounds;

-- name: CreateRound :one
INSERT INTO Rounds (TournamentID, RoundNumber, IsEliminationRound)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRound :one
UPDATE Rounds
SET TournamentID = $2, RoundNumber = $3, IsEliminationRound = $4
WHERE RoundID = $1
RETURNING *;

-- name: DeleteRound :exec
DELETE FROM Rounds WHERE RoundID = $1;