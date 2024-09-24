-- name: GetStudentByID :one
SELECT * FROM Students
WHERE StudentID = $1;

-- name: GetStudentByIDebateID :one
SELECT * FROM Students
WHERE iDebateStudentID = $1;

-- name: GetStudentByUserID :one
SELECT
    s.StudentID,
    s.FirstName,
    s.LastName,
    s.Email,
    s.SchoolID
FROM
    Users u
JOIN
    Students s ON u.UserID = s.UserID
WHERE
    u.UserID = $1
    AND u.UserRole = 'student';

-- name: GetStudentByEmail :one
SELECT * FROM Students
WHERE Email = $1;

-- name: GetAllStudents :many
SELECT s.*
FROM Students s
JOIN Users u ON s.UserID = u.UserID
WHERE u.UserRole = 'student'
  AND u.Status = 'approved'
  AND u.deleted_at IS NULL;

-- name: GetStudentsPaginated :many
SELECT s.*, sch.SchoolName
FROM Students s
JOIN Schools sch ON s.SchoolID = sch.SchoolID
ORDER BY s.StudentID
LIMIT $1 OFFSET $2;

-- name: GetTotalStudentCount :one
SELECT COUNT(*) FROM Students;

-- name: CreateStudent :one
INSERT INTO Students (FirstName, LastName, Gender, Grade, DateOfBirth, Email, Password, SchoolID, UserID)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)
RETURNING *;

-- name: UpdateStudent :one
UPDATE Students
SET FirstName = $2, LastName = $3, Grade = $4, DateOfBirth = $5, Email = $6, Password = $7, SchoolID = $8
WHERE StudentID = $1
RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM Students
WHERE StudentID = $1;

-- name: GetSchoolByContactPersonID :one
SELECT * FROM Schools
WHERE ContactPersonID = $1;

-- name: GetStudentsBySchoolID :many
SELECT * FROM Students
WHERE SchoolID = $1
ORDER BY StudentID
LIMIT $2 OFFSET $3;

-- name: GetStudentCountBySchoolID :one
SELECT COUNT(*) FROM Students
WHERE SchoolID = $1;