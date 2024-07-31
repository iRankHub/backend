// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: country.sql

package models

import (
	"context"
)

const getAllCountries = `-- name: GetAllCountries :many
SELECT countryname, isocode FROM CountryCodes
ORDER BY CountryName
`

func (q *Queries) GetAllCountries(ctx context.Context) ([]Countrycode, error) {
	rows, err := q.db.QueryContext(ctx, getAllCountries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Countrycode{}
	for rows.Next() {
		var i Countrycode
		if err := rows.Scan(&i.Countryname, &i.Isocode); err != nil {
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
