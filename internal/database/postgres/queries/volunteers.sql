-- name: GetVolunteerByID :one
SELECT * FROM Volunteers
WHERE VolunteerID = $1;

-- name: GetVolunteerByIDebateID :one
SELECT * FROM Volunteers
WHERE iDebateVolunteerID = $1;

-- name: GetVolunteerByUserID :one
SELECT * FROM volunteers
WHERE UserID = $1 LIMIT 1;

-- name: GetAllVolunteers :many
SELECT * FROM Volunteers;

-- name: GetVolunteersPaginated :many
SELECT *
FROM Volunteers
ORDER BY VolunteerID
LIMIT $1 OFFSET $2;

-- name: GetTotalVolunteerCount :one
SELECT COUNT(*) FROM Volunteers;

-- name: CreateVolunteer :one
INSERT INTO Volunteers (
  FirstName, LastName, DateOfBirth, Role, GraduateYear,
  Password, Gender, SafeGuardCertificate, HasInternship, UserID,
  IsEnrolledInUniversity, NationalID
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateVolunteer :one
UPDATE Volunteers
SET FirstName = $2, LastName = $3, DateOfBirth = $4, Role = $5, GraduateYear = $6,
    Password = $7, SafeGuardCertificate = $8, HasInternship = $9,
    IsEnrolledInUniversity = $10, NationalID = $11
WHERE VolunteerID = $1
RETURNING *;

-- name: DeleteVolunteer :exec
DELETE FROM Volunteers
WHERE VolunteerID = $1;