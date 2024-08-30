// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: teams.sql

package models

import (
	"context"
)

const deleteTeam = `-- name: DeleteTeam :exec
DELETE FROM Teams WHERE TeamID = $1
`

func (q *Queries) DeleteTeam(ctx context.Context, teamid int32) error {
	_, err := q.db.ExecContext(ctx, deleteTeam, teamid)
	return err
}

const getTeam = `-- name: GetTeam :one
SELECT teamid, name, schoolid, tournamentid FROM Teams WHERE TeamID = $1
`

func (q *Queries) GetTeam(ctx context.Context, teamid int32) (Team, error) {
	row := q.db.QueryRowContext(ctx, getTeam, teamid)
	var i Team
	err := row.Scan(
		&i.Teamid,
		&i.Name,
		&i.Tournamentid,
	)
	return i, err
}
