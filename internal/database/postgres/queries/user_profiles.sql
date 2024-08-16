-- name: GetUserProfile :one
SELECT * FROM UserProfiles
WHERE UserID = $1;

-- name: CreateUserProfile :one
INSERT INTO UserProfiles (UserID, Name, UserRole, Email, Password, VerificationStatus, Address, Phone, Bio, ProfilePicture, Gender)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE UserProfiles
SET Name = $2, UserRole = $3, Email = $4, VerificationStatus = $5, Address = $6, Phone = $7, Bio = $8, ProfilePicture = $9, Gender = $10
WHERE UserID = $1
RETURNING *;

-- name: DeleteUserProfile :exec
DELETE FROM UserProfiles
WHERE UserID = $1;

-- name: UpdateUserProfileBasicInfo :one
UPDATE UserProfiles
SET Name = $2, Email = $3, Gender = $4
WHERE UserID = $1
RETURNING *;