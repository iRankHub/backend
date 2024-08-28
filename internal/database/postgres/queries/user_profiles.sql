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

-- name: UpdateUserBasicInfo :exec
UPDATE Users
SET
    Name = COALESCE($2, Name),
    Email = COALESCE($3, Email),
    Gender = COALESCE($4, Gender),
    updated_at = CURRENT_TIMESTAMP
WHERE UserID = $1;

-- name: UpdateUserProfile :exec
UPDATE UserProfiles
SET
    Name = COALESCE($2, Name),
    Email = COALESCE($3, Email),
    Gender = COALESCE($4, Gender),
    Address = COALESCE($5, Address),
    Phone = COALESCE($6, Phone),
    Bio = COALESCE($7, Bio),
    ProfilePicture = COALESCE($8, ProfilePicture)
WHERE UserID = $1;

-- name: UpdateStudentProfile :exec
UPDATE Students
SET
    Grade = COALESCE($2, Grade),
    DateOfBirth = COALESCE($3, DateOfBirth),
    SchoolID = COALESCE($4, SchoolID)
WHERE UserID = $1;

-- name: UpdateSchoolProfile :exec
UPDATE Schools
SET
    SchoolName = COALESCE($2, SchoolName),
    Address = COALESCE($3, Address),
    Country = COALESCE($4, Country),
    Province = COALESCE($5, Province),
    District = COALESCE($6, District),
    SchoolType = COALESCE($7, SchoolType)
WHERE ContactPersonID = $1;

-- name: UpdateVolunteerProfile :exec
UPDATE Volunteers
SET
    Role = COALESCE($2, Role),
    GraduateYear = COALESCE($3, GraduateYear),
    SafeGuardCertificate = COALESCE($4, SafeGuardCertificate),
    HasInternship = COALESCE($5, HasInternship),
    IsEnrolledInUniversity = COALESCE($6, IsEnrolledInUniversity)
WHERE UserID = $1;

-- name: SoftDeleteUserProfile :exec
UPDATE Users
SET deleted_at = CURRENT_TIMESTAMP
WHERE UserID = $1;