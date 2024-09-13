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
WITH ApprovedCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE Status = 'approved' AND deleted_at IS NULL
),
RecentSignupsCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE created_at >= NOW() - INTERVAL '30 days' AND deleted_at IS NULL
)
SELECT
    u.*,
    CASE
        WHEN u.UserRole = 'student' THEN s.iDebateStudentID
        WHEN u.UserRole = 'volunteer' THEN v.iDebateVolunteerID
        WHEN u.UserRole = 'school' THEN sch.iDebateSchoolID
        WHEN u.UserRole = 'admin' THEN 'iDebate'
        ELSE NULL
    END AS iDebateID,
    CASE
        WHEN u.UserRole = 'school' THEN sch.SchoolName
        ELSE u.Name
    END AS DisplayName,
    (SELECT count FROM ApprovedCount) AS approved_users_count,
    (SELECT count FROM RecentSignupsCount) AS recent_signups_count
FROM Users u
LEFT JOIN Students s ON u.UserID = s.UserID
LEFT JOIN Volunteers v ON u.UserID = v.UserID
LEFT JOIN Schools sch ON u.UserID = sch.ContactPersonID
WHERE u.deleted_at IS NULL
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserStatistics :one
WITH AdminCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE UserRole = 'admin' AND deleted_at IS NULL
),
SchoolCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE UserRole = 'school' AND deleted_at IS NULL
),
StudentCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE UserRole = 'student' AND deleted_at IS NULL
),
VolunteerCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE UserRole = 'volunteer' AND deleted_at IS NULL
),
ApprovedCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE Status = 'approved' AND deleted_at IS NULL
),
NewRegistrationsCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE Status = 'pending' AND created_at >= NOW() - INTERVAL '30 days' AND deleted_at IS NULL
),
LastMonthNewUsersCount AS (
    SELECT COUNT(*) AS count
    FROM Users
    WHERE Status = 'pending' AND created_at >= NOW() - INTERVAL '60 days' AND created_at < NOW() - INTERVAL '30 days' AND deleted_at IS NULL
),
YesterdayApprovedCount AS (
    SELECT yesterday_approved_count
    FROM Users
    LIMIT 1
)
SELECT
    (SELECT count FROM AdminCount) AS admin_count,
    (SELECT count FROM SchoolCount) AS school_count,
    (SELECT count FROM StudentCount) AS student_count,
    (SELECT count FROM VolunteerCount) AS volunteer_count,
    (SELECT count FROM ApprovedCount) AS approved_count,
    (SELECT count FROM NewRegistrationsCount) AS new_registrations_count,
    (SELECT count FROM LastMonthNewUsersCount) AS last_month_new_users_count,
    (SELECT yesterday_approved_count FROM YesterdayApprovedCount) AS yesterday_approved_count;

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
WHERE UserRole IN ('volunteer', 'admin')
  AND Status = 'approved'
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalVolunteersAndAdminsCount :one
SELECT COUNT(*) FROM Users
WHERE UserRole IN ('volunteer', 'admin')
  AND Status = 'approved'
  AND deleted_at IS NULL;

-- name: UpdatePasswordAndClearResetCode :exec
WITH updated_users AS (
    UPDATE Users
    SET Password = $2, reset_token = NULL, reset_token_expires = NULL
    WHERE Users.UserID = $1
    RETURNING UserID, UserRole
)
UPDATE UserProfiles
SET Password = $2
WHERE UserProfiles.UserID = (SELECT UserID FROM updated_users);

UPDATE Students
SET Password = $2
WHERE Students.UserID = (SELECT UserID FROM updated_users)
  AND (SELECT UserRole FROM updated_users) = 'student';

UPDATE Volunteers
SET Password = $2
WHERE Volunteers.UserID = (SELECT UserID FROM updated_users)
  AND (SELECT UserRole FROM updated_users) = 'volunteer';

-- name: SetPasswordResetCodeAndGetUser :one
UPDATE Users
SET reset_token = $2, reset_token_expires = $3
WHERE UserID = $1
RETURNING UserID, Email, Name;

-- name: ValidateResetCodeAndGetUser :one
SELECT UserID, Email, Name, reset_token, reset_token_expires
FROM Users
WHERE UserID = $1 AND reset_token IS NOT NULL AND reset_token_expires > NOW();