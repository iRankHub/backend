-- name: GetSchoolByID :one
SELECT s.* FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE s.SchoolID = $1 AND u.deleted_at IS NULL;

-- name: GetSchoolByIDebateID :one
SELECT s.* FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE s.iDebateSchoolID = $1 AND u.deleted_at IS NULL;

-- name: GetSchoolByUserID :one
SELECT s.* FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE s.ContactPersonID = $1 AND u.deleted_at IS NULL;

-- name: GetSchoolIDByUserID :one
SELECT s.SchoolID
FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE u.UserID = $1;

-- name: GetSchoolByContactEmail :one
SELECT s.* FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE s.ContactEmail = $1 AND u.deleted_at IS NULL;

-- name: GetSchoolsPaginated :many
SELECT *
FROM Schools
ORDER BY SchoolID
LIMIT $1 OFFSET $2;

-- name: GetTotalSchoolCount :one
SELECT COUNT(*) FROM Schools;

-- name: GetSchoolsByDistrict :many
SELECT s.* FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE s.District = $1 AND u.deleted_at IS NULL;

-- name: GetSchoolsByCountry :many
SELECT s.* FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
WHERE s.Country = $1 AND u.deleted_at IS NULL;

-- name: GetSchoolsByLeague :many
SELECT s.*
FROM Schools s
JOIN Users u ON s.ContactPersonID = u.UserID
JOIN Leagues l ON l.LeagueID = $1
WHERE u.deleted_at IS NULL
  AND (
    (l.LeagueType = 'local' AND s.District = ANY(SELECT jsonb_array_elements_text(l.Details->'districts')))
    OR
    (l.LeagueType = 'international' AND s.Country = ANY(SELECT jsonb_array_elements_text(l.Details->'countries')))
  );


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