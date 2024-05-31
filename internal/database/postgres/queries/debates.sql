-- name: GetDebate :one
SELECT * FROM Debates WHERE DebateID = $1;

-- name: GetDebates :many
SELECT * FROM Debates;

-- name: CreateDebate :one
INSERT INTO Debates (RoundID, TournamentID, Team1ID, Team2ID, StartTime, EndTime, RoomID, Status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateDebate :one
UPDATE Debates
SET RoundID = $2, TournamentID = $3, Team1ID = $4, Team2ID = $5, StartTime = $6, EndTime = $7, RoomID = $8, Status = $9
WHERE DebateID = $1
RETURNING *;

-- name: DeleteDebate :exec
DELETE FROM Debates WHERE DebateID = $1;