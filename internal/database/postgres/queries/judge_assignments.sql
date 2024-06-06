-- name: GetJudgeAssignment :one
SELECT * FROM JudgeAssignments WHERE AssignmentID = $1;

-- name: GetJudgeAssignments :many
SELECT * FROM JudgeAssignments;

-- name: CreateJudgeAssignment :one
INSERT INTO JudgeAssignments (VolunteerID, TournamentID, DebateID)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateJudgeAssignment :one
UPDATE JudgeAssignments
SET VolunteerID = $2, TournamentID = $3, DebateID = $4
WHERE AssignmentID = $1
RETURNING *;

-- name: DeleteJudgeAssignment :exec
DELETE FROM JudgeAssignments WHERE AssignmentID = $1;