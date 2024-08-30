-- name: GetUserByID :one
SELECT * FROM Users
WHERE UserID = $1 AND deleted_at IS NULL;
-- name: GetUserByEmail :one
SELECT * FROM Users
WHERE Email = $1 AND deleted_at IS NULL;

-- name: GetUserByEmailOrIDebateIDAndUpdateLoginAttempt :one
WITH updated_user AS (
    UPDATE Users u
    SET last_login_attempt = NOW()
    WHERE u.UserID IN (
        SELECT u.UserID
        FROM Users u
        LEFT JOIN Students s ON u.UserID = s.UserID
        LEFT JOIN Schools sch ON u.UserID = sch.ContactPersonID
        LEFT JOIN Volunteers v ON u.UserID = v.UserID
        WHERE (u.Email = $1
           OR s.iDebateStudentID = $1
           OR sch.iDebateSchoolID = $1
           OR v.iDebateVolunteerID = $1)
        AND u.deleted_at IS NULL
        LIMIT 1
    )
    RETURNING *
)
SELECT u.*,
       s.iDebateStudentID,
       sch.iDebateSchoolID,
       v.iDebateVolunteerID
FROM updated_user u
LEFT JOIN Students s ON u.UserID = s.UserID
LEFT JOIN Schools sch ON u.UserID = sch.ContactPersonID
LEFT JOIN Volunteers v ON u.UserID = v.UserID;

-- name: GetUserEmailAndNameByID :one
SELECT UserID, Email, Name, Password, UserRole FROM Users WHERE UserID = $1;

-- name: CreateUser :one
INSERT INTO Users (Name, Email, Password, UserRole, Status, Gender)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateUser :one
UPDATE Users
SET Name = $2, Email = $3, Password = $4, UserRole = $5, VerificationStatus = $6, Status = $7, Gender = $8
WHERE UserID = $1
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM Users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalUserCount :one
SELECT COUNT(*) FROM Users WHERE deleted_at IS NULL;

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

-- name: RejectAndGetUser :one
UPDATE Users
SET Status = 'rejected', deleted_at = CURRENT_TIMESTAMP
WHERE UserID = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE Users
SET Password = $2
WHERE UserID = $1;

-- name: DeleteUser :exec
UPDATE Users
SET Status = 'rejected', deleted_at = CURRENT_TIMESTAMP
WHERE UserID = $1;

-- name: UpdateAndEnableTwoFactor :exec
UPDATE Users
SET two_factor_secret = $2, two_factor_enabled = TRUE
WHERE UserID = $1;

-- name: IncrementAndGetFailedLoginAttempts :one
UPDATE Users
SET failed_login_attempts = failed_login_attempts + 1,
    last_login_attempt = NOW()
WHERE UserID = $1
RETURNING *;

-- name: UpdateLastLogout :exec
UPDATE Users
SET last_logout = $2
WHERE UserID = $1;

-- name: ResetFailedLoginAttempts :exec
UPDATE Users SET failed_login_attempts = 0 WHERE UserID = $1;

-- name: SetResetToken :exec
UPDATE Users SET reset_token = $2, reset_token_expires = $3 WHERE UserID = $1;
-- name: ClearResetToken :exec
UPDATE Users SET reset_token = NULL, reset_token_expires = NULL WHERE UserID = $1;
-- name: GetUserForWebAuthn :one
SELECT UserID, WebAuthnUserID, Email, Name FROM Users WHERE UserID = $1;
-- name: GetUserForWebAuthnByEmail :one
SELECT UserID, WebAuthnUserID, Email, Name FROM Users WHERE Email = $1;
-- name: GetWebAuthnCredentials :many
SELECT CredentialID, PublicKey, AttestationType, AAGUID, SignCount
FROM WebAuthnCredentials WHERE UserID = $1;
-- name: StoreWebAuthnSessionData :exec
INSERT INTO WebAuthnSessionData (UserID, SessionData)
VALUES ($1, $2)
ON CONFLICT (UserID) DO UPDATE SET SessionData = $2;
-- name: GetWebAuthnSessionData :one
SELECT SessionData FROM WebAuthnSessionData WHERE UserID = $1;
-- name: StoreWebAuthnCredential :exec
INSERT INTO WebAuthnCredentials (UserID, CredentialID, PublicKey, AttestationType, AAGUID, SignCount)
VALUES ($1, $2, $3, $4, $5, $6);
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

-- name: GetVolunteersAndAdmins :many
SELECT * FROM Users
WHERE UserRole IN ('volunteer', 'admin') AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalVolunteersAndAdminsCount :one
SELECT COUNT(*) FROM Users
WHERE UserRole IN ('volunteer', 'admin') AND deleted_at IS NULL;

-- name: UpdatePasswordAndClearResetCode :exec
WITH updated_users AS (
    UPDATE Users
    SET Password = $2, reset_token = NULL, reset_token_expires = NULL
    WHERE Users.UserID = $1
    RETURNING UserID
)
UPDATE UserProfiles
SET Password = $2
WHERE UserProfiles.UserID = (SELECT UserID FROM updated_users);

-- name: SetPasswordResetCodeAndGetUser :one
UPDATE Users
SET reset_token = $2, reset_token_expires = $3
WHERE UserID = $1
RETURNING UserID, Email, Name;

-- name: ValidateResetCodeAndGetUser :one
SELECT UserID, Email, Name, reset_token, reset_token_expires
FROM Users
WHERE UserID = $1 AND reset_token IS NOT NULL AND reset_token_expires > NOW();