-- name: CreateNotification :one
INSERT INTO Notifications (UserID, Type, Message, RecipientEmail, Subject, IsRead)
VALUES ($1, $2, $3, $4, $5, FALSE)
RETURNING *;

-- name: GetUnreadNotifications :many
SELECT * FROM Notifications
WHERE UserID = $1 AND IsRead = FALSE
ORDER BY CreatedAt DESC;

-- name: MarkNotificationsAsRead :exec
UPDATE Notifications
SET IsRead = TRUE
WHERE UserID = $1 AND IsRead = FALSE;