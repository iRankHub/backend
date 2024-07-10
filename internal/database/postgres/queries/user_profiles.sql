-- name: GetUserProfile :one
SELECT * FROM UserProfiles
WHERE UserID = $1;

-- name: CreateUserProfile :one
INSERT INTO UserProfiles (UserID, Name, UserRole, Email, VerificationStatus, Address, Phone, Bio, ProfilePicture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE UserProfiles
SET Name = $2, UserRole = $3, Email = $4, VerificationStatus = $5, Address = $6, Phone = $7, Bio = $8, ProfilePicture = $9
WHERE UserID = $1
RETURNING *;

-- name: DeleteUserProfile :exec
DELETE FROM UserProfiles
WHERE UserID = $1;