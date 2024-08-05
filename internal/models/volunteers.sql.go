// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: volunteers.sql

package models

import (
	"context"
	"database/sql"
)

const createVolunteer = `-- name: CreateVolunteer :one
INSERT INTO Volunteers (
  FirstName, LastName, DateOfBirth, Role, GraduateYear,
  Password, SafeGuardCertificate, HasInternship, UserID, IsEnrolledInUniversity
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING volunteerid, idebatevolunteerid, firstname, lastname, dateofbirth, role, graduateyear, password, safeguardcertificate, hasinternship, isenrolledinuniversity, userid
`

type CreateVolunteerParams struct {
	Firstname              string        `json:"firstname"`
	Lastname               string        `json:"lastname"`
	Dateofbirth            sql.NullTime  `json:"dateofbirth"`
	Role                   string        `json:"role"`
	Graduateyear           sql.NullInt32 `json:"graduateyear"`
	Password               string        `json:"password"`
	Safeguardcertificate   sql.NullBool  `json:"safeguardcertificate"`
	Hasinternship          sql.NullBool  `json:"hasinternship"`
	Userid                 int32         `json:"userid"`
	Isenrolledinuniversity sql.NullBool  `json:"isenrolledinuniversity"`
}

func (q *Queries) CreateVolunteer(ctx context.Context, arg CreateVolunteerParams) (Volunteer, error) {
	row := q.db.QueryRowContext(ctx, createVolunteer,
		arg.Firstname,
		arg.Lastname,
		arg.Dateofbirth,
		arg.Role,
		arg.Graduateyear,
		arg.Password,
		arg.Safeguardcertificate,
		arg.Hasinternship,
		arg.Userid,
		arg.Isenrolledinuniversity,
	)
	var i Volunteer
	err := row.Scan(
		&i.Volunteerid,
		&i.Idebatevolunteerid,
		&i.Firstname,
		&i.Lastname,
		&i.Dateofbirth,
		&i.Role,
		&i.Graduateyear,
		&i.Password,
		&i.Safeguardcertificate,
		&i.Hasinternship,
		&i.Isenrolledinuniversity,
		&i.Userid,
	)
	return i, err
}

const deleteVolunteer = `-- name: DeleteVolunteer :exec
DELETE FROM Volunteers
WHERE VolunteerID = $1
`

func (q *Queries) DeleteVolunteer(ctx context.Context, volunteerid int32) error {
	_, err := q.db.ExecContext(ctx, deleteVolunteer, volunteerid)
	return err
}

const getAllVolunteers = `-- name: GetAllVolunteers :many
SELECT volunteerid, idebatevolunteerid, firstname, lastname, dateofbirth, role, graduateyear, password, safeguardcertificate, hasinternship, isenrolledinuniversity, userid FROM Volunteers
`

func (q *Queries) GetAllVolunteers(ctx context.Context) ([]Volunteer, error) {
	rows, err := q.db.QueryContext(ctx, getAllVolunteers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Volunteer{}
	for rows.Next() {
		var i Volunteer
		if err := rows.Scan(
			&i.Volunteerid,
			&i.Idebatevolunteerid,
			&i.Firstname,
			&i.Lastname,
			&i.Dateofbirth,
			&i.Role,
			&i.Graduateyear,
			&i.Password,
			&i.Safeguardcertificate,
			&i.Hasinternship,
			&i.Isenrolledinuniversity,
			&i.Userid,
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

const getTotalVolunteerCount = `-- name: GetTotalVolunteerCount :one
SELECT COUNT(*) FROM Volunteers
`

func (q *Queries) GetTotalVolunteerCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTotalVolunteerCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getVolunteerByID = `-- name: GetVolunteerByID :one
SELECT volunteerid, idebatevolunteerid, firstname, lastname, dateofbirth, role, graduateyear, password, safeguardcertificate, hasinternship, isenrolledinuniversity, userid FROM Volunteers
WHERE VolunteerID = $1
`

func (q *Queries) GetVolunteerByID(ctx context.Context, volunteerid int32) (Volunteer, error) {
	row := q.db.QueryRowContext(ctx, getVolunteerByID, volunteerid)
	var i Volunteer
	err := row.Scan(
		&i.Volunteerid,
		&i.Idebatevolunteerid,
		&i.Firstname,
		&i.Lastname,
		&i.Dateofbirth,
		&i.Role,
		&i.Graduateyear,
		&i.Password,
		&i.Safeguardcertificate,
		&i.Hasinternship,
		&i.Isenrolledinuniversity,
		&i.Userid,
	)
	return i, err
}

const getVolunteerByUserID = `-- name: GetVolunteerByUserID :one
SELECT volunteerid, idebatevolunteerid, firstname, lastname, dateofbirth, role, graduateyear, password, safeguardcertificate, hasinternship, isenrolledinuniversity, userid FROM volunteers
WHERE UserID = $1 LIMIT 1
`

func (q *Queries) GetVolunteerByUserID(ctx context.Context, userid int32) (Volunteer, error) {
	row := q.db.QueryRowContext(ctx, getVolunteerByUserID, userid)
	var i Volunteer
	err := row.Scan(
		&i.Volunteerid,
		&i.Idebatevolunteerid,
		&i.Firstname,
		&i.Lastname,
		&i.Dateofbirth,
		&i.Role,
		&i.Graduateyear,
		&i.Password,
		&i.Safeguardcertificate,
		&i.Hasinternship,
		&i.Isenrolledinuniversity,
		&i.Userid,
	)
	return i, err
}

const getVolunteersPaginated = `-- name: GetVolunteersPaginated :many
SELECT volunteerid, idebatevolunteerid, firstname, lastname, dateofbirth, role, graduateyear, password, safeguardcertificate, hasinternship, isenrolledinuniversity, userid
FROM Volunteers
ORDER BY VolunteerID
LIMIT $1 OFFSET $2
`

type GetVolunteersPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetVolunteersPaginated(ctx context.Context, arg GetVolunteersPaginatedParams) ([]Volunteer, error) {
	rows, err := q.db.QueryContext(ctx, getVolunteersPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Volunteer{}
	for rows.Next() {
		var i Volunteer
		if err := rows.Scan(
			&i.Volunteerid,
			&i.Idebatevolunteerid,
			&i.Firstname,
			&i.Lastname,
			&i.Dateofbirth,
			&i.Role,
			&i.Graduateyear,
			&i.Password,
			&i.Safeguardcertificate,
			&i.Hasinternship,
			&i.Isenrolledinuniversity,
			&i.Userid,
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

const updateVolunteer = `-- name: UpdateVolunteer :one
UPDATE Volunteers
SET FirstName = $2, LastName = $3, DateOfBirth = $4, Role = $5, GraduateYear = $6,
    Password = $7, SafeGuardCertificate = $8, hasinternship = $9, IsEnrolledInUniversity = $10
WHERE VolunteerID = $1
RETURNING volunteerid, idebatevolunteerid, firstname, lastname, dateofbirth, role, graduateyear, password, safeguardcertificate, hasinternship, isenrolledinuniversity, userid
`

type UpdateVolunteerParams struct {
	Volunteerid            int32         `json:"volunteerid"`
	Firstname              string        `json:"firstname"`
	Lastname               string        `json:"lastname"`
	Dateofbirth            sql.NullTime  `json:"dateofbirth"`
	Role                   string        `json:"role"`
	Graduateyear           sql.NullInt32 `json:"graduateyear"`
	Password               string        `json:"password"`
	Safeguardcertificate   sql.NullBool  `json:"safeguardcertificate"`
	Hasinternship          sql.NullBool  `json:"hasinternship"`
	Isenrolledinuniversity sql.NullBool  `json:"isenrolledinuniversity"`
}

func (q *Queries) UpdateVolunteer(ctx context.Context, arg UpdateVolunteerParams) (Volunteer, error) {
	row := q.db.QueryRowContext(ctx, updateVolunteer,
		arg.Volunteerid,
		arg.Firstname,
		arg.Lastname,
		arg.Dateofbirth,
		arg.Role,
		arg.Graduateyear,
		arg.Password,
		arg.Safeguardcertificate,
		arg.Hasinternship,
		arg.Isenrolledinuniversity,
	)
	var i Volunteer
	err := row.Scan(
		&i.Volunteerid,
		&i.Idebatevolunteerid,
		&i.Firstname,
		&i.Lastname,
		&i.Dateofbirth,
		&i.Role,
		&i.Graduateyear,
		&i.Password,
		&i.Safeguardcertificate,
		&i.Hasinternship,
		&i.Isenrolledinuniversity,
		&i.Userid,
	)
	return i, err
}
