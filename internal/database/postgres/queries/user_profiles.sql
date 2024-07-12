-- name: GetUserProfile :one
SELECT * FROM UserProfiles
WHERE UserID = $1;

-- name: CreateUserProfile :one
INSERT INTO UserProfiles (UserID, Name, UserRole, Email, Password, VerificationStatus, Address, Phone, Bio, ProfilePicture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE UserProfiles
SET Name = $2, UserRole = $3, Email = $4, Password = $5, VerificationStatus = $6, Address = $7, Phone = $8, Bio = $9, ProfilePicture = $10
WHERE UserID = $1
RETURNING *;

-- name: DeleteUserProfile :exec
WITH deleted_profile AS (
    DELETE FROM UserProfiles
    WHERE UserProfiles.UserID = $1
    RETURNING UserProfiles.UserID
)
UPDATE Users
SET deleted_at = CURRENT_TIMESTAMP
WHERE Users.UserID = (SELECT UserID FROM deleted_profile);