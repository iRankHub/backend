-- name: CreateUserProfile :one
INSERT INTO UserProfiles (UserID, Name, UserRole, Email, Password, VerificationStatus, Address, Phone, Bio, ProfilePicture, Gender)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetUserProfile :one
SELECT
    u.UserID, u.Name, u.Email, u.UserRole, u.password, u.Gender, u.VerificationStatus,
    u.two_factor_enabled, u.created_at AS SignUpDate,
    up.Address, up.Phone, up.Bio, up.ProfilePicture,
    s.Grade, s.DateOfBirth, s.SchoolID,
    sch.SchoolName, sch.Address AS SchoolAddress, sch.Country, sch.Province, sch.District, sch.SchoolType,
    v.Role AS VolunteerRole, v.GraduateYear, v.SafeGuardCertificate, v.HasInternship, v.IsEnrolledInUniversity,
    CASE WHEN EXISTS (
        SELECT 1 FROM WebAuthnCredentials wac WHERE wac.UserID = u.UserID
    ) THEN true ELSE false END AS biometric_auth_enabled
FROM Users u
LEFT JOIN UserProfiles up ON u.UserID = up.UserID
LEFT JOIN Students s ON u.UserID = s.UserID
LEFT JOIN Schools sch ON u.UserID = sch.ContactPersonID
LEFT JOIN Volunteers v ON u.UserID = v.UserID
WHERE u.UserID = $1 AND u.deleted_at IS NULL;

-- name: UpdateUserProfile :one
WITH updated_user AS (
    UPDATE Users
    SET Name = $2, Email = $3, Gender = $4, Password = $5, updated_at = CURRENT_TIMESTAMP
    WHERE Users.UserID = $1
    RETURNING UserID
)
UPDATE UserProfiles
SET Name = $2, Email = $3, Gender = $4, Address = $6, Phone = $7, Bio = $8, ProfilePicture = $9
WHERE UserID = (SELECT UserID FROM updated_user)
RETURNING
    UserProfiles.*,
    (SELECT Password FROM Users WHERE UserID = UserProfiles.UserID) AS Password,
    (SELECT updated_at FROM Users WHERE UserID = UserProfiles.UserID) AS updated_at;

-- name: UpdateStudentProfile :exec
UPDATE Students
SET Grade = $2, DateOfBirth = $3, SchoolID = $4
WHERE UserID = $1;

-- name: UpdateSchoolProfile :exec
UPDATE Schools
SET SchoolName = $2, Address = $3, Country = $4, Province = $5, District = $6, SchoolType = $7
WHERE ContactPersonID = $1;

-- name: UpdateVolunteerProfile :exec
UPDATE Volunteers
SET Role = $2, GraduateYear = $3, SafeGuardCertificate = $4, HasInternship = $5, IsEnrolledInUniversity = $6
WHERE UserID = $1;

-- name: SoftDeleteUserProfile :exec
UPDATE Users
SET deleted_at = CURRENT_TIMESTAMP
WHERE UserID = $1;