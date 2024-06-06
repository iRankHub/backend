-- name: GetBallot :one
SELECT * FROM Ballots WHERE BallotID = $1;

-- name: GetBallots :many
SELECT * FROM Ballots;

-- name: CreateBallot :one
INSERT INTO Ballots (DebateID, JudgeID, Team1DebaterAScore, Team1DebaterAComments, Team1DebaterBScore, Team1DebaterBComments, Team1DebaterCScore, Team1DebaterCComments, Team1TotalScore, Team2DebaterAScore, Team2DebaterAComments, Team2DebaterBScore, Team2DebaterBComments, Team2DebaterCScore, Team2DebaterCComments, Team2TotalScore)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: UpdateBallot :one
UPDATE Ballots
SET DebateID = $2, JudgeID = $3, Team1DebaterAScore = $4, Team1DebaterAComments = $5, Team1DebaterBScore = $6, Team1DebaterBComments = $7, Team1DebaterCScore = $8, Team1DebaterCComments = $9, Team1TotalScore = $10, Team2DebaterAScore = $11, Team2DebaterAComments = $12, Team2DebaterBScore = $13, Team2DebaterBComments = $14, Team2DebaterCScore = $15, Team2DebaterCComments = $16, Team2TotalScore = $17
WHERE BallotID = $1
RETURNING *;

-- name: DeleteBallot :exec
DELETE FROM Ballots WHERE BallotID = $1;