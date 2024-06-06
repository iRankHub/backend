-- name: GetVolunteerRatingType :one
SELECT * FROM VolunteerRatingTypes WHERE RatingTypeID = $1;

-- name: GetVolunteerRatingTypes :many
SELECT * FROM VolunteerRatingTypes;

-- name: CreateVolunteerRatingType :one
INSERT INTO VolunteerRatingTypes (RatingTypeName)
VALUES ($1)
RETURNING *;

-- name: UpdateVolunteerRatingType :one
UPDATE VolunteerRatingTypes
SET RatingTypeName = $2
WHERE RatingTypeID = $1
RETURNING *;

-- name: DeleteVolunteerRatingType :exec
DELETE FROM VolunteerRatingTypes WHERE RatingTypeID = $1;