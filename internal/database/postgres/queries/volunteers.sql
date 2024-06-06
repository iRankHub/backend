-- name: GetVolunteer :one
SELECT * FROM Volunteers WHERE VolunteerID = $1;

-- name: GetVolunteers :many
SELECT * FROM Volunteers;

-- name: CreateVolunteer :one
INSERT INTO Volunteers (Name, Role, UserID)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateVolunteer :one
UPDATE Volunteers
SET Name = $2, Role = $3, UserID = $4
WHERE VolunteerID = $1
RETURNING *;

-- name: DeleteVolunteer :exec
DELETE FROM Volunteers WHERE VolunteerID = $1;