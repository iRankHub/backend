-- name: GetSchool :one
SELECT * FROM Schools WHERE SchoolID = $1;

-- name: GetSchools :many
SELECT * FROM Schools;

-- name: CreateSchool :one
INSERT INTO Schools (Name, Address, ContactPersonID, ContactEmail, Category)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSchool :one
UPDATE Schools
SET Name = $2, Address = $3, ContactPersonID = $4, ContactEmail = $5, Category = $6
WHERE SchoolID = $1
RETURNING *;

-- name: DeleteSchool :exec
DELETE FROM Schools WHERE SchoolID = $1;