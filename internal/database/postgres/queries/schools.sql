-- name: GetSchoolByID :one
SELECT * FROM Schools
WHERE SchoolID = $1;

-- name: GetSchoolByUserID :one
SELECT * FROM Schools WHERE ContactPersonID = $1;

-- name: GetSchoolByContactEmail :one
SELECT * FROM Schools
WHERE ContactEmail = $1;

-- name: GetSchoolsByDistrict :many
SELECT * FROM Schools
WHERE District = $1;

-- name: GetSchoolsByCountry :many
SELECT * FROM Schools
WHERE Country = $1;

-- name: GetSchoolsByProvinceOrCountry :many
SELECT * FROM Schools WHERE Province = $1 OR Country = $1;

-- name: GetSchoolAddressByUserID :one
SELECT Address FROM Schools WHERE ContactPersonID = $1;

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