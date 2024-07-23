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

-- name: GetSchoolsByLeague :many
SELECT s.*
FROM Schools s
JOIN Leagues l ON (
    (l.LeagueType = 'local' AND s.District = ANY(l.Details->>'districts'))
    OR
    (l.LeagueType = 'international' AND s.Country = ANY(l.Details->>'countries'))
)
WHERE l.LeagueID = $1;

-- name: CreateSchool :one
INSERT INTO Schools (SchoolName, Address, Country, Province, District, ContactPersonID, ContactEmail, SchoolEmail, SchoolType)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateSchool :one
UPDATE Schools
SET SchoolName = $2, Address = $3, Country = $4, Province = $5, District = $6, ContactPersonID = $7, ContactEmail = $8, SchoolEmail = $9, SchoolType = $10
WHERE SchoolID = $1
RETURNING *;

-- name: DeleteSchool :exec
DELETE FROM Schools
WHERE SchoolID = $1;