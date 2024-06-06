-- name: GetJudgeReview :one
SELECT * FROM JudgeReviews WHERE ReviewID = $1;

-- name: GetJudgeReviews :many
SELECT * FROM JudgeReviews;

-- name: CreateJudgeReview :one
INSERT INTO JudgeReviews (StudentID, JudgeID, Rating, Comments)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateJudgeReview :one
UPDATE JudgeReviews
SET StudentID = $2, JudgeID = $3, Rating = $4, Comments = $5
WHERE ReviewID = $1
RETURNING *;

-- name: DeleteJudgeReview :exec
DELETE FROM JudgeReviews WHERE ReviewID = $1;