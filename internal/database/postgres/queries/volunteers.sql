-- name: GetVolunteerByID :one
SELECT * FROM Volunteers
WHERE VolunteerID = $1;

-- name: GetAllVolunteers :many
SELECT * FROM Volunteers;

-- name: CreateVolunteer :one
INSERT INTO Volunteers (FirstName, LastName, DateOfBirth, Role, GraduateYear, Password, SafeGuardCertificate, UserID)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateVolunteer :one
UPDATE Volunteers
SET FirstName = $2, LastName = $3, DateOfBirth = $4, Role = $5, GraduateYear = $6, Password = $7, SafeGuardCertificate = $8
WHERE VolunteerID = $1
RETURNING *;

-- name: DeleteVolunteer :exec
DELETE FROM Volunteers
WHERE VolunteerID = $1;