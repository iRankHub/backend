-- name: GetUserByID :one
SELECT * FROM Users
WHERE UserID = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM Users
WHERE Email = $1 AND deleted_at IS NULL;

-- name: GetUserEmailAndNameByID :one
SELECT UserID, Email, Name, Password, UserRole FROM Users WHERE UserID = $1;

-- name: CreateUser :one
INSERT INTO Users (Name, Email, Password, UserRole, Status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE Users
SET Name = $2, Email = $3, Password = $4, UserRole = $5, VerificationStatus = $6, Status = $7
WHERE UserID = $1
RETURNING *;

-- name: UpdateUserStatus :exec
UPDATE Users
SET Status = $2
WHERE UserID = $1;

-- name: GetPendingUsers :many
SELECT * FROM Users
WHERE Status = 'pending' AND deleted_at IS NULL;

-- name: GetUsersByStatus :many
SELECT * FROM Users
WHERE Status = $1 AND deleted_at IS NULL;

-- name: UpdateUserPassword :exec
UPDATE Users
SET Password = $2
WHERE UserID = $1;

-- name: DeleteUser :exec
UPDATE Users
SET deleted_at = CURRENT_TIMESTAMP
WHERE UserID = $1;

-- name: UpdateUserTwoFactorSecret :exec
UPDATE Users SET two_factor_secret = $2 WHERE UserID = $1;

-- name: EnableTwoFactor :exec
UPDATE Users SET two_factor_enabled = TRUE WHERE UserID = $1;

-- name: DisableTwoFactor :exec
UPDATE Users
SET two_factor_enabled = FALSE, two_factor_secret = NULL
WHERE UserID = $1;

-- name: UpdateLastLoginAttempt :exec
UPDATE Users SET last_login_attempt = NOW() WHERE UserID = $1;

-- name: UpdateLastLogout :exec
UPDATE Users
SET last_logout = $2
WHERE UserID = $1;

-- name: IncrementFailedLoginAttempts :exec
UPDATE Users
SET failed_login_attempts = failed_login_attempts + 1,
    last_login_attempt = NOW()
WHERE UserID = $1;

-- name: ResetFailedLoginAttempts :exec
UPDATE Users SET failed_login_attempts = 0 WHERE UserID = $1;

-- name: SetResetToken :exec
UPDATE Users SET reset_token = $2, reset_token_expires = $3 WHERE UserID = $1;

-- name: ClearResetToken :exec
UPDATE Users SET reset_token = NULL, reset_token_expires = NULL WHERE UserID = $1;

-- name: SetBiometricToken :exec
UPDATE Users SET biometric_token = $2 WHERE UserID = $1;

-- name: GetUserByBiometricToken :one
SELECT * FROM Users
WHERE biometric_token = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByResetToken :one
SELECT * FROM Users
WHERE reset_token = $1 AND reset_token_expires > NOW() AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserWithAuthDetails :one
SELECT * FROM Users
WHERE UserID = $1 AND deleted_at IS NULL;

-- name: DeactivateAccount :exec
UPDATE Users
SET DeactivatedAt = CURRENT_TIMESTAMP
WHERE UserID = $1;

-- name: ReactivateAccount :exec
UPDATE Users
SET DeactivatedAt = NULL
WHERE UserID = $1;

-- name: GetAccountStatus :one
SELECT
    CASE
        WHEN DeactivatedAt IS NULL THEN 'active'
        ELSE 'deactivated'
    END AS status
FROM Users
WHERE UserID = $1;