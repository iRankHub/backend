// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: schools.sql

package models

import (
	"context"
	"database/sql"
)

const createSchool = `-- name: CreateSchool :one
INSERT INTO Schools (SchoolName, Address, Country, Province, District, ContactPersonID, ContactEmail, SchoolEmail, SchoolType)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype
`

type CreateSchoolParams struct {
	Schoolname      string         `json:"schoolname"`
	Address         string         `json:"address"`
	Country         sql.NullString `json:"country"`
	Province        sql.NullString `json:"province"`
	District        sql.NullString `json:"district"`
	Contactpersonid int32          `json:"contactpersonid"`
	Contactemail    string         `json:"contactemail"`
	Schoolemail     string         `json:"schoolemail"`
	Schooltype      string         `json:"schooltype"`
}

func (q *Queries) CreateSchool(ctx context.Context, arg CreateSchoolParams) (School, error) {
	row := q.db.QueryRowContext(ctx, createSchool,
		arg.Schoolname,
		arg.Address,
		arg.Country,
		arg.Province,
		arg.District,
		arg.Contactpersonid,
		arg.Contactemail,
		arg.Schoolemail,
		arg.Schooltype,
	)
	var i School
	err := row.Scan(
		&i.Schoolid,
		&i.Idebateschoolid,
		&i.Schoolname,
		&i.Address,
		&i.Country,
		&i.Province,
		&i.District,
		&i.Contactpersonid,
		&i.Contactemail,
		&i.Schoolemail,
		&i.Schooltype,
	)
	return i, err
}

const deleteSchool = `-- name: DeleteSchool :exec
DELETE FROM Schools
WHERE SchoolID = $1
`

func (q *Queries) DeleteSchool(ctx context.Context, schoolid int32) error {
	_, err := q.db.ExecContext(ctx, deleteSchool, schoolid)
	return err
}

const getSchoolByContactEmail = `-- name: GetSchoolByContactEmail :one
SELECT schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype FROM Schools
WHERE ContactEmail = $1
`

func (q *Queries) GetSchoolByContactEmail(ctx context.Context, contactemail string) (School, error) {
	row := q.db.QueryRowContext(ctx, getSchoolByContactEmail, contactemail)
	var i School
	err := row.Scan(
		&i.Schoolid,
		&i.Idebateschoolid,
		&i.Schoolname,
		&i.Address,
		&i.Country,
		&i.Province,
		&i.District,
		&i.Contactpersonid,
		&i.Contactemail,
		&i.Schoolemail,
		&i.Schooltype,
	)
	return i, err
}

const getSchoolByID = `-- name: GetSchoolByID :one
SELECT schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype FROM Schools
WHERE SchoolID = $1
`

func (q *Queries) GetSchoolByID(ctx context.Context, schoolid int32) (School, error) {
	row := q.db.QueryRowContext(ctx, getSchoolByID, schoolid)
	var i School
	err := row.Scan(
		&i.Schoolid,
		&i.Idebateschoolid,
		&i.Schoolname,
		&i.Address,
		&i.Country,
		&i.Province,
		&i.District,
		&i.Contactpersonid,
		&i.Contactemail,
		&i.Schoolemail,
		&i.Schooltype,
	)
	return i, err
}

const getSchoolByUserID = `-- name: GetSchoolByUserID :one
SELECT schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype FROM Schools WHERE ContactPersonID = $1
`

func (q *Queries) GetSchoolByUserID(ctx context.Context, contactpersonid int32) (School, error) {
	row := q.db.QueryRowContext(ctx, getSchoolByUserID, contactpersonid)
	var i School
	err := row.Scan(
		&i.Schoolid,
		&i.Idebateschoolid,
		&i.Schoolname,
		&i.Address,
		&i.Country,
		&i.Province,
		&i.District,
		&i.Contactpersonid,
		&i.Contactemail,
		&i.Schoolemail,
		&i.Schooltype,
	)
	return i, err
}

const getSchoolsByCountry = `-- name: GetSchoolsByCountry :many
SELECT schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype FROM Schools
WHERE Country = $1
`

func (q *Queries) GetSchoolsByCountry(ctx context.Context, country sql.NullString) ([]School, error) {
	rows, err := q.db.QueryContext(ctx, getSchoolsByCountry, country)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []School{}
	for rows.Next() {
		var i School
		if err := rows.Scan(
			&i.Schoolid,
			&i.Idebateschoolid,
			&i.Schoolname,
			&i.Address,
			&i.Country,
			&i.Province,
			&i.District,
			&i.Contactpersonid,
			&i.Contactemail,
			&i.Schoolemail,
			&i.Schooltype,
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

const getSchoolsByDistrict = `-- name: GetSchoolsByDistrict :many
SELECT schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype FROM Schools
WHERE District = $1
`

func (q *Queries) GetSchoolsByDistrict(ctx context.Context, district sql.NullString) ([]School, error) {
	rows, err := q.db.QueryContext(ctx, getSchoolsByDistrict, district)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []School{}
	for rows.Next() {
		var i School
		if err := rows.Scan(
			&i.Schoolid,
			&i.Idebateschoolid,
			&i.Schoolname,
			&i.Address,
			&i.Country,
			&i.Province,
			&i.District,
			&i.Contactpersonid,
			&i.Contactemail,
			&i.Schoolemail,
			&i.Schooltype,
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

const getSchoolsByLeague = `-- name: GetSchoolsByLeague :many
SELECT s.schoolid, s.idebateschoolid, s.schoolname, s.address, s.country, s.province, s.district, s.contactpersonid, s.contactemail, s.schoolemail, s.schooltype
FROM Schools s
JOIN Leagues l ON l.LeagueID = $1
WHERE
    (l.LeagueType = 'local' AND s.District = ANY(SELECT jsonb_array_elements_text(l.Details->'districts')))
    OR
    (l.LeagueType = 'international' AND s.Country = ANY(SELECT jsonb_array_elements_text(l.Details->'countries')))
`

func (q *Queries) GetSchoolsByLeague(ctx context.Context, leagueid int32) ([]School, error) {
	rows, err := q.db.QueryContext(ctx, getSchoolsByLeague, leagueid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []School{}
	for rows.Next() {
		var i School
		if err := rows.Scan(
			&i.Schoolid,
			&i.Idebateschoolid,
			&i.Schoolname,
			&i.Address,
			&i.Country,
			&i.Province,
			&i.District,
			&i.Contactpersonid,
			&i.Contactemail,
			&i.Schoolemail,
			&i.Schooltype,
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

const getSchoolsPaginated = `-- name: GetSchoolsPaginated :many
SELECT schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype
FROM Schools
ORDER BY SchoolID
LIMIT $1 OFFSET $2
`

type GetSchoolsPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetSchoolsPaginated(ctx context.Context, arg GetSchoolsPaginatedParams) ([]School, error) {
	rows, err := q.db.QueryContext(ctx, getSchoolsPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []School{}
	for rows.Next() {
		var i School
		if err := rows.Scan(
			&i.Schoolid,
			&i.Idebateschoolid,
			&i.Schoolname,
			&i.Address,
			&i.Country,
			&i.Province,
			&i.District,
			&i.Contactpersonid,
			&i.Contactemail,
			&i.Schoolemail,
			&i.Schooltype,
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

const getTotalSchoolCount = `-- name: GetTotalSchoolCount :one
SELECT COUNT(*) FROM Schools
`

func (q *Queries) GetTotalSchoolCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTotalSchoolCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateSchool = `-- name: UpdateSchool :one
UPDATE Schools
SET SchoolName = $2, Address = $3, Country = $4, Province = $5, District = $6, ContactPersonID = $7, ContactEmail = $8, SchoolEmail = $9, SchoolType = $10
WHERE SchoolID = $1
RETURNING schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype
`

type UpdateSchoolParams struct {
	Schoolid        int32          `json:"schoolid"`
	Schoolname      string         `json:"schoolname"`
	Address         string         `json:"address"`
	Country         sql.NullString `json:"country"`
	Province        sql.NullString `json:"province"`
	District        sql.NullString `json:"district"`
	Contactpersonid int32          `json:"contactpersonid"`
	Contactemail    string         `json:"contactemail"`
	Schoolemail     string         `json:"schoolemail"`
	Schooltype      string         `json:"schooltype"`
}

func (q *Queries) UpdateSchool(ctx context.Context, arg UpdateSchoolParams) (School, error) {
	row := q.db.QueryRowContext(ctx, updateSchool,
		arg.Schoolid,
		arg.Schoolname,
		arg.Address,
		arg.Country,
		arg.Province,
		arg.District,
		arg.Contactpersonid,
		arg.Contactemail,
		arg.Schoolemail,
		arg.Schooltype,
	)
	var i School
	err := row.Scan(
		&i.Schoolid,
		&i.Idebateschoolid,
		&i.Schoolname,
		&i.Address,
		&i.Country,
		&i.Province,
		&i.District,
		&i.Contactpersonid,
		&i.Contactemail,
		&i.Schoolemail,
		&i.Schooltype,
	)
	return i, err
}

const updateSchoolAddress = `-- name: UpdateSchoolAddress :one
UPDATE Schools
SET Address = $2
WHERE ContactPersonID = $1
RETURNING schoolid, idebateschoolid, schoolname, address, country, province, district, contactpersonid, contactemail, schoolemail, schooltype
`

type UpdateSchoolAddressParams struct {
	Contactpersonid int32  `json:"contactpersonid"`
	Address         string `json:"address"`
}

func (q *Queries) UpdateSchoolAddress(ctx context.Context, arg UpdateSchoolAddressParams) (School, error) {
	row := q.db.QueryRowContext(ctx, updateSchoolAddress, arg.Contactpersonid, arg.Address)
	var i School
	err := row.Scan(
		&i.Schoolid,
		&i.Idebateschoolid,
		&i.Schoolname,
		&i.Address,
		&i.Country,
		&i.Province,
		&i.District,
		&i.Contactpersonid,
		&i.Contactemail,
		&i.Schoolemail,
		&i.Schooltype,
	)
	return i, err
}
