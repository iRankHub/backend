-- name: GetUserByID :one
SELECT * FROM Users
WHERE UserID = $1;

-- name: GetUserByEmail :one
SELECT * FROM Users
WHERE Email = $1;

-- name: CreateUser :one
INSERT INTO Users (Name, Email, Password, UserRole)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE Users
SET Name = $2, Email = $3, Password = $4, UserRole = $5, VerificationStatus = $6, ApprovalStatus = $7
WHERE UserID = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM Users
WHERE UserID = $1;