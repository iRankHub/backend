// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: communications.sql

package models

import (
	"context"
	"time"
)

const createCommunication = `-- name: CreateCommunication :one
INSERT INTO Communications (UserID, SchoolID, Type, Content, Timestamp)
VALUES ($1, $2, $3, $4, $5)
RETURNING communicationid, userid, schoolid, type, content, timestamp
`

type CreateCommunicationParams struct {
	Userid    int32     `json:"userid"`
	Schoolid  int32     `json:"schoolid"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func (q *Queries) CreateCommunication(ctx context.Context, arg CreateCommunicationParams) (Communication, error) {
	row := q.db.QueryRowContext(ctx, createCommunication,
		arg.Userid,
		arg.Schoolid,
		arg.Type,
		arg.Content,
		arg.Timestamp,
	)
	var i Communication
	err := row.Scan(
		&i.Communicationid,
		&i.Userid,
		&i.Schoolid,
		&i.Type,
		&i.Content,
		&i.Timestamp,
	)
	return i, err
}

const deleteCommunication = `-- name: DeleteCommunication :exec
DELETE FROM Communications WHERE CommunicationID = $1
`

func (q *Queries) DeleteCommunication(ctx context.Context, communicationid int32) error {
	_, err := q.db.ExecContext(ctx, deleteCommunication, communicationid)
	return err
}

const getCommunication = `-- name: GetCommunication :one
SELECT communicationid, userid, schoolid, type, content, timestamp FROM Communications WHERE CommunicationID = $1
`

func (q *Queries) GetCommunication(ctx context.Context, communicationid int32) (Communication, error) {
	row := q.db.QueryRowContext(ctx, getCommunication, communicationid)
	var i Communication
	err := row.Scan(
		&i.Communicationid,
		&i.Userid,
		&i.Schoolid,
		&i.Type,
		&i.Content,
		&i.Timestamp,
	)
	return i, err
}

const getCommunications = `-- name: GetCommunications :many
SELECT communicationid, userid, schoolid, type, content, timestamp FROM Communications
`

func (q *Queries) GetCommunications(ctx context.Context) ([]Communication, error) {
	rows, err := q.db.QueryContext(ctx, getCommunications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Communication{}
	for rows.Next() {
		var i Communication
		if err := rows.Scan(
			&i.Communicationid,
			&i.Userid,
			&i.Schoolid,
			&i.Type,
			&i.Content,
			&i.Timestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCommunication = `-- name: UpdateCommunication :one
UPDATE Communications
SET UserID = $2, SchoolID = $3, Type = $4, Content = $5, Timestamp = $6
WHERE CommunicationID = $1
RETURNING communicationid, userid, schoolid, type, content, timestamp
`

type UpdateCommunicationParams struct {
	Communicationid int32     `json:"communicationid"`
	Userid          int32     `json:"userid"`
	Schoolid        int32     `json:"schoolid"`
	Type            string    `json:"type"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
}

func (q *Queries) UpdateCommunication(ctx context.Context, arg UpdateCommunicationParams) (Communication, error) {
	row := q.db.QueryRowContext(ctx, updateCommunication,
		arg.Communicationid,
		arg.Userid,
		arg.Schoolid,
		arg.Type,
		arg.Content,
		arg.Timestamp,
	)
	var i Communication
	err := row.Scan(
		&i.Communicationid,
		&i.Userid,
		&i.Schoolid,
		&i.Type,
		&i.Content,
		&i.Timestamp,
	)
	return i, err
}
