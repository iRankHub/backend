-- name: GetStudent :one
SELECT studentid, userid, dateofbirth, schoolid, uniquestudentid FROM Students WHERE StudentID = $1;

-- name: ListStudents :many
SELECT studentid, userid, dateofbirth, schoolid, uniquestudentid FROM Students;

-- name: CreateStudent :one
INSERT INTO Students (UserID, DateOfBirth, SchoolID, UniqueStudentID) VALUES ($1, $2, $3, $4) RETURNING studentid, userid, dateofbirth, schoolid, uniquestudentid;

-- name: UpdateStudent :one
UPDATE Students SET UserID = $2, DateOfBirth = $3, SchoolID = $4, UniqueStudentID = $5 WHERE StudentID = $1 RETURNING studentid, userid, dateofbirth, schoolid, uniquestudentid;

-- name: DeleteStudent :exec
DELETE FROM Students WHERE StudentID = $1;