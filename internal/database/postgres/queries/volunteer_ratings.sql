-- name: GetVolunteerRating :one
SELECT * FROM VolunteerRatings WHERE RatingID = $1;

-- name: GetVolunteerRatings :many
SELECT * FROM VolunteerRatings;

-- name: CreateVolunteerRating :one
INSERT INTO VolunteerRatings (VolunteerID, RatingTypeID, RatingScore, RatingComments, CumulativeRating)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateVolunteerRating :one
UPDATE VolunteerRatings
SET VolunteerID = $2, RatingTypeID = $3, RatingScore = $4, RatingComments = $5, CumulativeRating = $6
WHERE RatingID = $1
RETURNING *;

-- name: DeleteVolunteerRating :exec
DELETE FROM VolunteerRatings WHERE RatingID = $1;