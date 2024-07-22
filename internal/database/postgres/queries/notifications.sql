-- name: CreateNotification :one
INSERT INTO Notifications (UserID, Type, Message)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUnreadNotifications :many
SELECT * FROM Notifications
WHERE UserID = $1 AND IsRead = FALSE
ORDER BY CreatedAt DESC;

-- name: MarkNotificationsAsRead :exec
UPDATE Notifications
SET IsRead = TRUE
WHERE UserID = $1 AND IsRead = FALSE;