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

-- name: UpdateUserTwoFactorSecret :exec
UPDATE Users SET two_factor_secret = $2 WHERE UserID = $1;

-- name: EnableTwoFactor :exec
UPDATE Users SET two_factor_enabled = TRUE WHERE UserID = $1;

-- name: DisableTwoFactor :exec
UPDATE Users SET two_factor_enabled = FALSE WHERE UserID = $1;

-- name: IncrementFailedLoginAttempts :exec
UPDATE Users SET failed_login_attempts = failed_login_attempts + 1, last_login_attempt = NOW() WHERE UserID = $1;

-- name: ResetFailedLoginAttempts :exec
UPDATE Users SET failed_login_attempts = 0 WHERE UserID = $1;

-- name: SetResetToken :exec
UPDATE Users SET reset_token = $2, reset_token_expires = $3 WHERE UserID = $1;

-- name: ClearResetToken :exec
UPDATE Users SET reset_token = NULL, reset_token_expires = NULL WHERE UserID = $1;

-- name: SetBiometricToken :exec
UPDATE Users SET biometric_token = $2 WHERE UserID = $1;

-- name: GetUserByBiometricToken :one
SELECT * FROM Users WHERE biometric_token = $1 LIMIT 1;

-- name: GetUserByResetToken :one
SELECT * FROM Users WHERE reset_token = $1 AND reset_token_expires > NOW() LIMIT 1;

-- name: GetUserWithAuthDetails :one
SELECT * FROM Users WHERE UserID = $1;