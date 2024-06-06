-- name: GetStudentRank :one
SELECT * FROM StudentRanks WHERE RankID = $1;

-- name: GetStudentRanks :many
SELECT * FROM StudentRanks;

-- name: CreateStudentRank :one
INSERT INTO StudentRanks (StudentID, TournamentID, RankValue, RankComments)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateStudentRank :one
UPDATE StudentRanks
SET StudentID = $2, TournamentID = $3, RankValue = $4, RankComments = $5
WHERE RankID = $1
RETURNING *;

-- name: DeleteStudentRank :exec
DELETE FROM StudentRanks WHERE RankID = $1;