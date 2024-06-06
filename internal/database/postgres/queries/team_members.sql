-- name: GetTeamMember :one
SELECT * FROM TeamMembers WHERE TeamID = $1 AND StudentID = $2;

-- name: GetTeamMembers :many
SELECT * FROM TeamMembers WHERE TeamID = $1;

-- name: AddTeamMember :exec
INSERT INTO TeamMembers (TeamID, StudentID)
VALUES ($1, $2);

-- name: RemoveTeamMember :exec
DELETE FROM TeamMembers WHERE TeamID = $1 AND StudentID = $2;