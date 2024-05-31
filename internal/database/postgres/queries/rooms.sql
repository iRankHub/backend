-- name: GetRoom :one
SELECT * FROM Rooms WHERE RoomID = $1;

-- name: GetRooms :many
SELECT * FROM Rooms;

-- name: CreateRoom :one
INSERT INTO Rooms (RoomName, Location, Capacity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRoom :one
UPDATE Rooms
SET RoomName = $2, Location = $3, Capacity = $4
WHERE RoomID = $1
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM Rooms WHERE RoomID = $1;