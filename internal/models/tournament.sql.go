// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tournament.sql

package models

import (
	"context"
	"database/sql"
	"time"
)

const createInternationalLeagueDetails = `-- name: CreateInternationalLeagueDetails :exec
INSERT INTO InternationalLeagueDetails (LeagueID, Continent, Country)
VALUES ($1, $2, $3)
`

type CreateInternationalLeagueDetailsParams struct {
	Leagueid  int32          `json:"leagueid"`
	Continent sql.NullString `json:"continent"`
	Country   sql.NullString `json:"country"`
}

func (q *Queries) CreateInternationalLeagueDetails(ctx context.Context, arg CreateInternationalLeagueDetailsParams) error {
	_, err := q.db.ExecContext(ctx, createInternationalLeagueDetails, arg.Leagueid, arg.Continent, arg.Country)
	return err
}

const createLeague = `-- name: CreateLeague :one
INSERT INTO Leagues (Name, LeagueType)
VALUES ($1, $2)
RETURNING leagueid, name, leaguetype, deleted_at
`

type CreateLeagueParams struct {
	Name       string `json:"name"`
	Leaguetype string `json:"leaguetype"`
}

// League Queries
func (q *Queries) CreateLeague(ctx context.Context, arg CreateLeagueParams) (League, error) {
	row := q.db.QueryRowContext(ctx, createLeague, arg.Name, arg.Leaguetype)
	var i League
	err := row.Scan(
		&i.Leagueid,
		&i.Name,
		&i.Leaguetype,
		&i.DeletedAt,
	)
	return i, err
}

const createLocalLeagueDetails = `-- name: CreateLocalLeagueDetails :exec
INSERT INTO LocalLeagueDetails (LeagueID, Province, District)
VALUES ($1, $2, $3)
`

type CreateLocalLeagueDetailsParams struct {
	Leagueid int32          `json:"leagueid"`
	Province sql.NullString `json:"province"`
	District sql.NullString `json:"district"`
}

func (q *Queries) CreateLocalLeagueDetails(ctx context.Context, arg CreateLocalLeagueDetailsParams) error {
	_, err := q.db.ExecContext(ctx, createLocalLeagueDetails, arg.Leagueid, arg.Province, arg.District)
	return err
}

const createTournamentEntry = `-- name: CreateTournamentEntry :one
INSERT INTO Tournaments (Name, StartDate, EndDate, Location, FormatID, LeagueID, NumberOfPreliminaryRounds, NumberOfEliminationRounds, JudgesPerDebatePreliminary, JudgesPerDebateElimination, TournamentFee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING tournamentid, name, startdate, enddate, location, formatid, leagueid, numberofpreliminaryrounds, numberofeliminationrounds, judgesperdebatepreliminary, judgesperdebateelimination, tournamentfee, deleted_at
`

type CreateTournamentEntryParams struct {
	Name                       string        `json:"name"`
	Startdate                  time.Time     `json:"startdate"`
	Enddate                    time.Time     `json:"enddate"`
	Location                   string        `json:"location"`
	Formatid                   int32         `json:"formatid"`
	Leagueid                   sql.NullInt32 `json:"leagueid"`
	Numberofpreliminaryrounds  int32         `json:"numberofpreliminaryrounds"`
	Numberofeliminationrounds  int32         `json:"numberofeliminationrounds"`
	Judgesperdebatepreliminary int32         `json:"judgesperdebatepreliminary"`
	Judgesperdebateelimination int32         `json:"judgesperdebateelimination"`
	Tournamentfee              string        `json:"tournamentfee"`
}

// Tournament Queries
func (q *Queries) CreateTournamentEntry(ctx context.Context, arg CreateTournamentEntryParams) (Tournament, error) {
	row := q.db.QueryRowContext(ctx, createTournamentEntry,
		arg.Name,
		arg.Startdate,
		arg.Enddate,
		arg.Location,
		arg.Formatid,
		arg.Leagueid,
		arg.Numberofpreliminaryrounds,
		arg.Numberofeliminationrounds,
		arg.Judgesperdebatepreliminary,
		arg.Judgesperdebateelimination,
		arg.Tournamentfee,
	)
	var i Tournament
	err := row.Scan(
		&i.Tournamentid,
		&i.Name,
		&i.Startdate,
		&i.Enddate,
		&i.Location,
		&i.Formatid,
		&i.Leagueid,
		&i.Numberofpreliminaryrounds,
		&i.Numberofeliminationrounds,
		&i.Judgesperdebatepreliminary,
		&i.Judgesperdebateelimination,
		&i.Tournamentfee,
		&i.DeletedAt,
	)
	return i, err
}

const createTournamentFormat = `-- name: CreateTournamentFormat :one
INSERT INTO TournamentFormats (FormatName, Description, SpeakersPerTeam)
VALUES ($1, $2, $3)
RETURNING formatid, formatname, description, speakersperteam, deleted_at
`

type CreateTournamentFormatParams struct {
	Formatname      string         `json:"formatname"`
	Description     sql.NullString `json:"description"`
	Speakersperteam int32          `json:"speakersperteam"`
}

// Tournament Format Queries
func (q *Queries) CreateTournamentFormat(ctx context.Context, arg CreateTournamentFormatParams) (Tournamentformat, error) {
	row := q.db.QueryRowContext(ctx, createTournamentFormat, arg.Formatname, arg.Description, arg.Speakersperteam)
	var i Tournamentformat
	err := row.Scan(
		&i.Formatid,
		&i.Formatname,
		&i.Description,
		&i.Speakersperteam,
		&i.DeletedAt,
	)
	return i, err
}

const deleteLeagueByID = `-- name: DeleteLeagueByID :exec
UPDATE Leagues
SET deleted_at = CURRENT_TIMESTAMP
WHERE LeagueID = $1
`

func (q *Queries) DeleteLeagueByID(ctx context.Context, leagueid int32) error {
	_, err := q.db.ExecContext(ctx, deleteLeagueByID, leagueid)
	return err
}

const deleteTournamentByID = `-- name: DeleteTournamentByID :exec
UPDATE Tournaments
SET deleted_at = CURRENT_TIMESTAMP
WHERE TournamentID = $1
`

func (q *Queries) DeleteTournamentByID(ctx context.Context, tournamentid int32) error {
	_, err := q.db.ExecContext(ctx, deleteTournamentByID, tournamentid)
	return err
}

const deleteTournamentFormatByID = `-- name: DeleteTournamentFormatByID :exec
UPDATE TournamentFormats
SET deleted_at = CURRENT_TIMESTAMP
WHERE FormatID = $1
`

func (q *Queries) DeleteTournamentFormatByID(ctx context.Context, formatid int32) error {
	_, err := q.db.ExecContext(ctx, deleteTournamentFormatByID, formatid)
	return err
}

const getInternationalLeagueDetails = `-- name: GetInternationalLeagueDetails :one
SELECT leagueid, continent, country FROM InternationalLeagueDetails WHERE LeagueID = $1
`

func (q *Queries) GetInternationalLeagueDetails(ctx context.Context, leagueid int32) (Internationalleaguedetail, error) {
	row := q.db.QueryRowContext(ctx, getInternationalLeagueDetails, leagueid)
	var i Internationalleaguedetail
	err := row.Scan(&i.Leagueid, &i.Continent, &i.Country)
	return i, err
}

const getLeagueByID = `-- name: GetLeagueByID :one
SELECT l.leagueid, l.name, l.leaguetype, l.deleted_at,
       COALESCE(lld.Province, ild.Continent) AS detail1,
       COALESCE(lld.District, ild.Country) AS detail2
FROM Leagues l
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
WHERE l.LeagueID = $1 AND l.deleted_at IS NULL
`

type GetLeagueByIDRow struct {
	Leagueid   int32          `json:"leagueid"`
	Name       string         `json:"name"`
	Leaguetype string         `json:"leaguetype"`
	DeletedAt  sql.NullTime   `json:"deleted_at"`
	Detail1    sql.NullString `json:"detail1"`
	Detail2    sql.NullString `json:"detail2"`
}

func (q *Queries) GetLeagueByID(ctx context.Context, leagueid int32) (GetLeagueByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getLeagueByID, leagueid)
	var i GetLeagueByIDRow
	err := row.Scan(
		&i.Leagueid,
		&i.Name,
		&i.Leaguetype,
		&i.DeletedAt,
		&i.Detail1,
		&i.Detail2,
	)
	return i, err
}

const getLocalLeagueDetails = `-- name: GetLocalLeagueDetails :one
SELECT leagueid, province, district FROM LocalLeagueDetails WHERE LeagueID = $1
`

func (q *Queries) GetLocalLeagueDetails(ctx context.Context, leagueid int32) (Localleaguedetail, error) {
	row := q.db.QueryRowContext(ctx, getLocalLeagueDetails, leagueid)
	var i Localleaguedetail
	err := row.Scan(&i.Leagueid, &i.Province, &i.District)
	return i, err
}

const getTournamentByID = `-- name: GetTournamentByID :one
SELECT t.tournamentid, t.name, t.startdate, t.enddate, t.location, t.formatid, t.leagueid, t.numberofpreliminaryrounds, t.numberofeliminationrounds, t.judgesperdebatepreliminary, t.judgesperdebateelimination, t.tournamentfee, t.deleted_at, tf.formatid, tf.formatname, tf.description, tf.speakersperteam, tf.deleted_at, l.leagueid, l.name, l.leaguetype, l.deleted_at,
       COALESCE(lld.Province, ild.Continent) AS league_detail1,
       COALESCE(lld.District, ild.Country) AS league_detail2,
       tc.VolunteerID as CoordinatorID
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
LEFT JOIN TournamentCoordinators tc ON t.TournamentID = tc.TournamentID
WHERE t.TournamentID = $1 AND t.deleted_at IS NULL
`

type GetTournamentByIDRow struct {
	Tournamentid               int32          `json:"tournamentid"`
	Name                       string         `json:"name"`
	Startdate                  time.Time      `json:"startdate"`
	Enddate                    time.Time      `json:"enddate"`
	Location                   string         `json:"location"`
	Formatid                   int32          `json:"formatid"`
	Leagueid                   sql.NullInt32  `json:"leagueid"`
	Numberofpreliminaryrounds  int32          `json:"numberofpreliminaryrounds"`
	Numberofeliminationrounds  int32          `json:"numberofeliminationrounds"`
	Judgesperdebatepreliminary int32          `json:"judgesperdebatepreliminary"`
	Judgesperdebateelimination int32          `json:"judgesperdebateelimination"`
	Tournamentfee              string         `json:"tournamentfee"`
	DeletedAt                  sql.NullTime   `json:"deleted_at"`
	Formatid_2                 int32          `json:"formatid_2"`
	Formatname                 string         `json:"formatname"`
	Description                sql.NullString `json:"description"`
	Speakersperteam            int32          `json:"speakersperteam"`
	DeletedAt_2                sql.NullTime   `json:"deleted_at_2"`
	Leagueid_2                 int32          `json:"leagueid_2"`
	Name_2                     string         `json:"name_2"`
	Leaguetype                 string         `json:"leaguetype"`
	DeletedAt_3                sql.NullTime   `json:"deleted_at_3"`
	LeagueDetail1              sql.NullString `json:"league_detail1"`
	LeagueDetail2              sql.NullString `json:"league_detail2"`
	Coordinatorid              sql.NullInt32  `json:"coordinatorid"`
}

func (q *Queries) GetTournamentByID(ctx context.Context, tournamentid int32) (GetTournamentByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getTournamentByID, tournamentid)
	var i GetTournamentByIDRow
	err := row.Scan(
		&i.Tournamentid,
		&i.Name,
		&i.Startdate,
		&i.Enddate,
		&i.Location,
		&i.Formatid,
		&i.Leagueid,
		&i.Numberofpreliminaryrounds,
		&i.Numberofeliminationrounds,
		&i.Judgesperdebatepreliminary,
		&i.Judgesperdebateelimination,
		&i.Tournamentfee,
		&i.DeletedAt,
		&i.Formatid_2,
		&i.Formatname,
		&i.Description,
		&i.Speakersperteam,
		&i.DeletedAt_2,
		&i.Leagueid_2,
		&i.Name_2,
		&i.Leaguetype,
		&i.DeletedAt_3,
		&i.LeagueDetail1,
		&i.LeagueDetail2,
		&i.Coordinatorid,
	)
	return i, err
}

const getTournamentFormatByID = `-- name: GetTournamentFormatByID :one
SELECT formatid, formatname, description, speakersperteam, deleted_at FROM TournamentFormats
WHERE FormatID = $1 AND deleted_at IS NULL
`

func (q *Queries) GetTournamentFormatByID(ctx context.Context, formatid int32) (Tournamentformat, error) {
	row := q.db.QueryRowContext(ctx, getTournamentFormatByID, formatid)
	var i Tournamentformat
	err := row.Scan(
		&i.Formatid,
		&i.Formatname,
		&i.Description,
		&i.Speakersperteam,
		&i.DeletedAt,
	)
	return i, err
}

const listLeaguesPaginated = `-- name: ListLeaguesPaginated :many
SELECT l.leagueid, l.name, l.leaguetype, l.deleted_at,
       COALESCE(lld.Province, ild.Continent) AS detail1,
       COALESCE(lld.District, ild.Country) AS detail2
FROM Leagues l
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
WHERE l.deleted_at IS NULL
ORDER BY l.LeagueID
LIMIT $1 OFFSET $2
`

type ListLeaguesPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListLeaguesPaginatedRow struct {
	Leagueid   int32          `json:"leagueid"`
	Name       string         `json:"name"`
	Leaguetype string         `json:"leaguetype"`
	DeletedAt  sql.NullTime   `json:"deleted_at"`
	Detail1    sql.NullString `json:"detail1"`
	Detail2    sql.NullString `json:"detail2"`
}

func (q *Queries) ListLeaguesPaginated(ctx context.Context, arg ListLeaguesPaginatedParams) ([]ListLeaguesPaginatedRow, error) {
	rows, err := q.db.QueryContext(ctx, listLeaguesPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListLeaguesPaginatedRow{}
	for rows.Next() {
		var i ListLeaguesPaginatedRow
		if err := rows.Scan(
			&i.Leagueid,
			&i.Name,
			&i.Leaguetype,
			&i.DeletedAt,
			&i.Detail1,
			&i.Detail2,
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

const listTournamentFormatsPaginated = `-- name: ListTournamentFormatsPaginated :many
SELECT formatid, formatname, description, speakersperteam, deleted_at FROM TournamentFormats
WHERE deleted_at IS NULL
ORDER BY FormatID
LIMIT $1 OFFSET $2
`

type ListTournamentFormatsPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTournamentFormatsPaginated(ctx context.Context, arg ListTournamentFormatsPaginatedParams) ([]Tournamentformat, error) {
	rows, err := q.db.QueryContext(ctx, listTournamentFormatsPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Tournamentformat{}
	for rows.Next() {
		var i Tournamentformat
		if err := rows.Scan(
			&i.Formatid,
			&i.Formatname,
			&i.Description,
			&i.Speakersperteam,
			&i.DeletedAt,
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

const listTournamentsPaginated = `-- name: ListTournamentsPaginated :many
SELECT t.tournamentid, t.name, t.startdate, t.enddate, t.location, t.formatid, t.leagueid, t.numberofpreliminaryrounds, t.numberofeliminationrounds, t.judgesperdebatepreliminary, t.judgesperdebateelimination, t.tournamentfee, t.deleted_at, tf.formatid, tf.formatname, tf.description, tf.speakersperteam, tf.deleted_at, l.leagueid, l.name, l.leaguetype, l.deleted_at,
       COALESCE(lld.Province, ild.Continent) AS league_detail1,
       COALESCE(lld.District, ild.Country) AS league_detail2,
       tc.VolunteerID as CoordinatorID
FROM Tournaments t
JOIN TournamentFormats tf ON t.FormatID = tf.FormatID
JOIN Leagues l ON t.LeagueID = l.LeagueID
LEFT JOIN LocalLeagueDetails lld ON l.LeagueID = lld.LeagueID
LEFT JOIN InternationalLeagueDetails ild ON l.LeagueID = ild.LeagueID
LEFT JOIN TournamentCoordinators tc ON t.TournamentID = tc.TournamentID
WHERE t.deleted_at IS NULL
ORDER BY t.TournamentID
LIMIT $1 OFFSET $2
`

type ListTournamentsPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListTournamentsPaginatedRow struct {
	Tournamentid               int32          `json:"tournamentid"`
	Name                       string         `json:"name"`
	Startdate                  time.Time      `json:"startdate"`
	Enddate                    time.Time      `json:"enddate"`
	Location                   string         `json:"location"`
	Formatid                   int32          `json:"formatid"`
	Leagueid                   sql.NullInt32  `json:"leagueid"`
	Numberofpreliminaryrounds  int32          `json:"numberofpreliminaryrounds"`
	Numberofeliminationrounds  int32          `json:"numberofeliminationrounds"`
	Judgesperdebatepreliminary int32          `json:"judgesperdebatepreliminary"`
	Judgesperdebateelimination int32          `json:"judgesperdebateelimination"`
	Tournamentfee              string         `json:"tournamentfee"`
	DeletedAt                  sql.NullTime   `json:"deleted_at"`
	Formatid_2                 int32          `json:"formatid_2"`
	Formatname                 string         `json:"formatname"`
	Description                sql.NullString `json:"description"`
	Speakersperteam            int32          `json:"speakersperteam"`
	DeletedAt_2                sql.NullTime   `json:"deleted_at_2"`
	Leagueid_2                 int32          `json:"leagueid_2"`
	Name_2                     string         `json:"name_2"`
	Leaguetype                 string         `json:"leaguetype"`
	DeletedAt_3                sql.NullTime   `json:"deleted_at_3"`
	LeagueDetail1              sql.NullString `json:"league_detail1"`
	LeagueDetail2              sql.NullString `json:"league_detail2"`
	Coordinatorid              sql.NullInt32  `json:"coordinatorid"`
}

func (q *Queries) ListTournamentsPaginated(ctx context.Context, arg ListTournamentsPaginatedParams) ([]ListTournamentsPaginatedRow, error) {
	rows, err := q.db.QueryContext(ctx, listTournamentsPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListTournamentsPaginatedRow{}
	for rows.Next() {
		var i ListTournamentsPaginatedRow
		if err := rows.Scan(
			&i.Tournamentid,
			&i.Name,
			&i.Startdate,
			&i.Enddate,
			&i.Location,
			&i.Formatid,
			&i.Leagueid,
			&i.Numberofpreliminaryrounds,
			&i.Numberofeliminationrounds,
			&i.Judgesperdebatepreliminary,
			&i.Judgesperdebateelimination,
			&i.Tournamentfee,
			&i.DeletedAt,
			&i.Formatid_2,
			&i.Formatname,
			&i.Description,
			&i.Speakersperteam,
			&i.DeletedAt_2,
			&i.Leagueid_2,
			&i.Name_2,
			&i.Leaguetype,
			&i.DeletedAt_3,
			&i.LeagueDetail1,
			&i.LeagueDetail2,
			&i.Coordinatorid,
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

const updateInternationalLeagueDetailsInfo = `-- name: UpdateInternationalLeagueDetailsInfo :exec
UPDATE InternationalLeagueDetails
SET Continent = $2, Country = $3
WHERE LeagueID = $1
`

type UpdateInternationalLeagueDetailsInfoParams struct {
	Leagueid  int32          `json:"leagueid"`
	Continent sql.NullString `json:"continent"`
	Country   sql.NullString `json:"country"`
}

func (q *Queries) UpdateInternationalLeagueDetailsInfo(ctx context.Context, arg UpdateInternationalLeagueDetailsInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateInternationalLeagueDetailsInfo, arg.Leagueid, arg.Continent, arg.Country)
	return err
}

const updateLeagueDetails = `-- name: UpdateLeagueDetails :one
UPDATE Leagues
SET Name = $2, LeagueType = $3
WHERE LeagueID = $1
RETURNING leagueid, name, leaguetype, deleted_at
`

type UpdateLeagueDetailsParams struct {
	Leagueid   int32  `json:"leagueid"`
	Name       string `json:"name"`
	Leaguetype string `json:"leaguetype"`
}

func (q *Queries) UpdateLeagueDetails(ctx context.Context, arg UpdateLeagueDetailsParams) (League, error) {
	row := q.db.QueryRowContext(ctx, updateLeagueDetails, arg.Leagueid, arg.Name, arg.Leaguetype)
	var i League
	err := row.Scan(
		&i.Leagueid,
		&i.Name,
		&i.Leaguetype,
		&i.DeletedAt,
	)
	return i, err
}

const updateLocalLeagueDetailsInfo = `-- name: UpdateLocalLeagueDetailsInfo :exec
UPDATE LocalLeagueDetails
SET Province = $2, District = $3
WHERE LeagueID = $1
`

type UpdateLocalLeagueDetailsInfoParams struct {
	Leagueid int32          `json:"leagueid"`
	Province sql.NullString `json:"province"`
	District sql.NullString `json:"district"`
}

func (q *Queries) UpdateLocalLeagueDetailsInfo(ctx context.Context, arg UpdateLocalLeagueDetailsInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateLocalLeagueDetailsInfo, arg.Leagueid, arg.Province, arg.District)
	return err
}

const updateTournamentDetails = `-- name: UpdateTournamentDetails :one
UPDATE Tournaments
SET Name = $2, StartDate = $3, EndDate = $4, Location = $5, FormatID = $6, LeagueID = $7, NumberOfPreliminaryRounds = $8, NumberOfEliminationRounds = $9, JudgesPerDebatePreliminary = $10, JudgesPerDebateElimination = $11, TournamentFee = $12
WHERE TournamentID = $1
RETURNING tournamentid, name, startdate, enddate, location, formatid, leagueid, numberofpreliminaryrounds, numberofeliminationrounds, judgesperdebatepreliminary, judgesperdebateelimination, tournamentfee, deleted_at
`

type UpdateTournamentDetailsParams struct {
	Tournamentid               int32         `json:"tournamentid"`
	Name                       string        `json:"name"`
	Startdate                  time.Time     `json:"startdate"`
	Enddate                    time.Time     `json:"enddate"`
	Location                   string        `json:"location"`
	Formatid                   int32         `json:"formatid"`
	Leagueid                   sql.NullInt32 `json:"leagueid"`
	Numberofpreliminaryrounds  int32         `json:"numberofpreliminaryrounds"`
	Numberofeliminationrounds  int32         `json:"numberofeliminationrounds"`
	Judgesperdebatepreliminary int32         `json:"judgesperdebatepreliminary"`
	Judgesperdebateelimination int32         `json:"judgesperdebateelimination"`
	Tournamentfee              string        `json:"tournamentfee"`
}

func (q *Queries) UpdateTournamentDetails(ctx context.Context, arg UpdateTournamentDetailsParams) (Tournament, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentDetails,
		arg.Tournamentid,
		arg.Name,
		arg.Startdate,
		arg.Enddate,
		arg.Location,
		arg.Formatid,
		arg.Leagueid,
		arg.Numberofpreliminaryrounds,
		arg.Numberofeliminationrounds,
		arg.Judgesperdebatepreliminary,
		arg.Judgesperdebateelimination,
		arg.Tournamentfee,
	)
	var i Tournament
	err := row.Scan(
		&i.Tournamentid,
		&i.Name,
		&i.Startdate,
		&i.Enddate,
		&i.Location,
		&i.Formatid,
		&i.Leagueid,
		&i.Numberofpreliminaryrounds,
		&i.Numberofeliminationrounds,
		&i.Judgesperdebatepreliminary,
		&i.Judgesperdebateelimination,
		&i.Tournamentfee,
		&i.DeletedAt,
	)
	return i, err
}

const updateTournamentFormatDetails = `-- name: UpdateTournamentFormatDetails :one
UPDATE TournamentFormats
SET FormatName = $2, Description = $3, SpeakersPerTeam = $4
WHERE FormatID = $1
RETURNING formatid, formatname, description, speakersperteam, deleted_at
`

type UpdateTournamentFormatDetailsParams struct {
	Formatid        int32          `json:"formatid"`
	Formatname      string         `json:"formatname"`
	Description     sql.NullString `json:"description"`
	Speakersperteam int32          `json:"speakersperteam"`
}

func (q *Queries) UpdateTournamentFormatDetails(ctx context.Context, arg UpdateTournamentFormatDetailsParams) (Tournamentformat, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentFormatDetails,
		arg.Formatid,
		arg.Formatname,
		arg.Description,
		arg.Speakersperteam,
	)
	var i Tournamentformat
	err := row.Scan(
		&i.Formatid,
		&i.Formatname,
		&i.Description,
		&i.Speakersperteam,
		&i.DeletedAt,
	)
	return i, err
}
