-- name: GetUserProfile :one
SELECT * FROM UserProfiles
WHERE UserID = $1;

-- name: CreateUserProfile :one
INSERT INTO UserProfiles (UserID, Address, Phone, Bio, ProfilePicture)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE UserProfiles
SET Address = $2, Phone = $3, Bio = $4, ProfilePicture = $5
WHERE UserID = $1
RETURNING *;

-- name: DeleteUserProfile :exec
DELETE FROM UserProfiles
WHERE UserID = $1;