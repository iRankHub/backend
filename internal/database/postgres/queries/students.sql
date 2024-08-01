-- name: GetStudentByID :one
SELECT * FROM Students
WHERE StudentID = $1;

-- name: GetStudentByEmail :one
SELECT * FROM Students
WHERE Email = $1;

-- name: GetAllStudents :many
SELECT * FROM Students;

-- name: GetStudentsPaginated :many
SELECT s.*, sch.SchoolName
FROM Students s
JOIN Schools sch ON s.SchoolID = sch.SchoolID
ORDER BY s.StudentID
LIMIT $1 OFFSET $2;

-- name: GetTotalStudentCount :one
SELECT COUNT(*) FROM Students;

-- name: CreateStudent :one
INSERT INTO Students (FirstName, LastName, Grade, DateOfBirth, Email, Password, SchoolID, UserID)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateStudent :one
UPDATE Students
SET FirstName = $2, LastName = $3, Grade = $4, DateOfBirth = $5, Email = $6, Password = $7, SchoolID = $8
WHERE StudentID = $1
RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM Students
WHERE StudentID = $1;