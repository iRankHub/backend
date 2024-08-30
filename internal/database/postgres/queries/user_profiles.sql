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

-- name: UpdateAdminProfile :exec
WITH updated_admin AS (
    UPDATE Users
    SET Name = COALESCE($2, Name),
        Gender = COALESCE($3, Gender),
        Email = COALESCE($4, Email)
    WHERE Users.UserID = $1
    RETURNING Users.UserID
)
UPDATE UserProfiles
SET Name = COALESCE($2, Name),
    Gender = COALESCE($3, Gender),
    Email = COALESCE($4, Email),
    Address = COALESCE($5, Address),
    Bio = COALESCE($6, Bio),
    Phone = COALESCE($7, Phone),
    ProfilePicture = COALESCE($8, ProfilePicture)
WHERE UserProfiles.UserID = (SELECT UserID FROM updated_admin);

-- name: UpdateSchoolUser :exec
UPDATE Users
SET Name = COALESCE($2, Name),
    Gender = COALESCE($3, Gender),
    Email = COALESCE($4, Email)
WHERE UserID = $1;

-- name: UpdateSchoolUserProfile :exec
UPDATE UserProfiles
SET Name = COALESCE($2, Name),
    Email = COALESCE($3, Email),
    Gender = COALESCE($4, Gender),
    Address = COALESCE($5, Address),
    Phone = COALESCE($6, Phone),
    Bio = COALESCE($7, Bio),
    ProfilePicture = COALESCE($8, ProfilePicture)
WHERE UserProfiles.UserID = $1;

-- name: UpdateSchoolDetails :exec
UPDATE Schools
SET ContactPersonNationalID = COALESCE($2, ContactPersonNationalID),
    SchoolName = COALESCE($3, SchoolName),
    Address = COALESCE($4, Address),
    SchoolEmail = COALESCE($5, SchoolEmail),
    SchoolType = COALESCE($6, SchoolType),
    ContactEmail = COALESCE($7, ContactEmail)
WHERE ContactPersonID = $1;

-- name: SoftDeleteUserProfile :exec
UPDATE Users
SET deleted_at = CURRENT_TIMESTAMP
WHERE UserID = $1;