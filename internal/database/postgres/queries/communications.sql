-- name: GetCommunication :one
SELECT * FROM Communications WHERE CommunicationID = $1;

-- name: GetCommunications :many
SELECT * FROM Communications;

-- name: CreateCommunication :one
INSERT INTO Communications (UserID, SchoolID, Type, Content, Timestamp)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateCommunication :one
UPDATE Communications
SET UserID = $2, SchoolID = $3, Type = $4, Content = $5, Timestamp = $6
WHERE CommunicationID = $1
RETURNING *;

-- name: DeleteCommunication :exec
DELETE FROM Communications WHERE CommunicationID = $1;