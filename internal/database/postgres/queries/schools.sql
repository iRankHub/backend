-- name: GetSchoolByID :one
SELECT * FROM Schools
WHERE SchoolID = $1;

-- name: GetSchoolByContactEmail :one
SELECT * FROM Schools
WHERE ContactEmail = $1;

-- name: CreateSchool :one
INSERT INTO Schools (SchoolName, Address, Country, Province, District, ContactPersonID, ContactEmail, SchoolType)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateSchool :one
UPDATE Schools
SET SchoolName = $2, Address = $3, Country = $4, Province = $5, District = $6, ContactPersonID = $7, ContactEmail = $8, SchoolType = $9
WHERE SchoolID = $1
RETURNING *;

-- name: DeleteSchool :exec
DELETE FROM Schools
WHERE SchoolID = $1;