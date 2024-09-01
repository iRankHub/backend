-- name: GetSchoolByID :one
SELECT * FROM Schools
WHERE SchoolID = $1;

-- name: GetSchoolByIDebateID :one
SELECT * FROM Schools
WHERE iDebateSchoolID = $1;

-- name: GetSchoolByUserID :one
SELECT * FROM Schools WHERE ContactPersonID = $1;

-- name: GetSchoolByContactEmail :one
SELECT * FROM Schools
WHERE ContactEmail = $1;

-- name: GetSchoolsPaginated :many
SELECT *
FROM Schools
ORDER BY SchoolID
LIMIT $1 OFFSET $2;

-- name: GetTotalSchoolCount :one
SELECT COUNT(*) FROM Schools;

-- name: GetSchoolsByDistrict :many
SELECT * FROM Schools
WHERE District = $1;

-- name: GetSchoolsByCountry :many
SELECT * FROM Schools
WHERE Country = $1;

-- name: GetSchoolsByLeague :many
SELECT s.*
FROM Schools s
JOIN Leagues l ON l.LeagueID = $1
WHERE
    (l.LeagueType = 'local' AND s.District = ANY(SELECT jsonb_array_elements_text(l.Details->'districts')))
    OR
    (l.LeagueType = 'international' AND s.Country = ANY(SELECT jsonb_array_elements_text(l.Details->'countries')));

-- name: CreateSchool :one
INSERT INTO Schools (
  SchoolName, Address, Country, Province, District, ContactPersonID,
  ContactEmail, SchoolEmail, SchoolType, ContactPersonNationalID
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateSchool :one
UPDATE Schools
SET SchoolName = $2, Address = $3, Country = $4, Province = $5, District = $6,
    ContactPersonID = $7, ContactEmail = $8, SchoolEmail = $9, SchoolType = $10,
    ContactPersonNationalID = $11
WHERE SchoolID = $1
RETURNING *;

-- name: DeleteSchool :exec
DELETE FROM Schools
WHERE SchoolID = $1;

-- name: UpdateSchoolAddress :one
UPDATE Schools
SET Address = $2
WHERE ContactPersonID = $1
RETURNING *;

-- name: GetSchoolIDByName :one
SELECT SchoolID FROM Schools WHERE SchoolName = $1;

-- name: GetSchoolIDsByNames :many
SELECT SchoolID, SchoolName FROM Schools WHERE LOWER(SchoolName) = ANY(ARRAY(SELECT LOWER(unnest($1::text[]))));