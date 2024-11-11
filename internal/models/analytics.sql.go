// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: analytics.sql

package models

import (
	"context"
	"database/sql"
	"time"
)

const getExpensesByTournament = `-- name: GetExpensesByTournament :many
SELECT
    t.tournamentid,
    t.name as tournament_name,
    te.foodexpense,
    te.transportexpense,
    te.perdiemexpense,
    te.awardingexpense,
    te.stationaryexpense,
    te.otherexpenses,
    te.totalexpense
FROM tournaments t
JOIN tournamentexpenses te ON t.tournamentid = te.tournamentid
WHERE t.startdate BETWEEN $1 AND $2
AND ($3::INTEGER IS NULL OR t.tournamentid = $3)
ORDER BY t.startdate DESC
`

type GetExpensesByTournamentParams struct {
	Startdate   time.Time `json:"startdate"`
	Startdate_2 time.Time `json:"startdate_2"`
	Column3     int32     `json:"column_3"`
}

type GetExpensesByTournamentRow struct {
	Tournamentid      int32          `json:"tournamentid"`
	TournamentName    string         `json:"tournament_name"`
	Foodexpense       string         `json:"foodexpense"`
	Transportexpense  string         `json:"transportexpense"`
	Perdiemexpense    string         `json:"perdiemexpense"`
	Awardingexpense   string         `json:"awardingexpense"`
	Stationaryexpense string         `json:"stationaryexpense"`
	Otherexpenses     string         `json:"otherexpenses"`
	Totalexpense      sql.NullString `json:"totalexpense"`
}

func (q *Queries) GetExpensesByTournament(ctx context.Context, arg GetExpensesByTournamentParams) ([]GetExpensesByTournamentRow, error) {
	rows, err := q.db.QueryContext(ctx, getExpensesByTournament, arg.Startdate, arg.Startdate_2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetExpensesByTournamentRow{}
	for rows.Next() {
		var i GetExpensesByTournamentRow
		if err := rows.Scan(
			&i.Tournamentid,
			&i.TournamentName,
			&i.Foodexpense,
			&i.Transportexpense,
			&i.Perdiemexpense,
			&i.Awardingexpense,
			&i.Stationaryexpense,
			&i.Otherexpenses,
			&i.Totalexpense,
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

const getExpensesSummary = `-- name: GetExpensesSummary :one
SELECT
    COALESCE(SUM(foodexpense), 0) as food_expense,
    COALESCE(SUM(transportexpense), 0) as transport_expense,
    COALESCE(SUM(perdiemexpense), 0) as per_diem_expense,
    COALESCE(SUM(awardingexpense), 0) as awarding_expense,
    COALESCE(SUM(stationaryexpense), 0) as stationary_expense,
    COALESCE(SUM(otherexpenses), 0) as other_expenses,
    COALESCE(SUM(totalexpense), 0) as total_expense
FROM tournamentexpenses te
JOIN tournaments t ON te.tournamentid = t.tournamentid
WHERE t.startdate BETWEEN $1 AND $2
AND ($3::INTEGER IS NULL OR t.tournamentid = $3)
`

type GetExpensesSummaryParams struct {
	Startdate   time.Time `json:"startdate"`
	Startdate_2 time.Time `json:"startdate_2"`
	Column3     int32     `json:"column_3"`
}

type GetExpensesSummaryRow struct {
	FoodExpense       interface{} `json:"food_expense"`
	TransportExpense  interface{} `json:"transport_expense"`
	PerDiemExpense    interface{} `json:"per_diem_expense"`
	AwardingExpense   interface{} `json:"awarding_expense"`
	StationaryExpense interface{} `json:"stationary_expense"`
	OtherExpenses     interface{} `json:"other_expenses"`
	TotalExpense      interface{} `json:"total_expense"`
}

func (q *Queries) GetExpensesSummary(ctx context.Context, arg GetExpensesSummaryParams) (GetExpensesSummaryRow, error) {
	row := q.db.QueryRowContext(ctx, getExpensesSummary, arg.Startdate, arg.Startdate_2, arg.Column3)
	var i GetExpensesSummaryRow
	err := row.Scan(
		&i.FoodExpense,
		&i.TransportExpense,
		&i.PerDiemExpense,
		&i.AwardingExpense,
		&i.StationaryExpense,
		&i.OtherExpenses,
		&i.TotalExpense,
	)
	return i, err
}

const getSchoolPerformanceByCategory = `-- name: GetSchoolPerformanceByCategory :many
SELECT
    s.schooltype as group_name,
    COALESCE(SUM(str.actualpaidamount), 0) as total_amount,
    COUNT(DISTINCT s.schoolid) as school_count
FROM schools s
JOIN schooltournamentregistrations str ON s.schoolid = str.schoolid
JOIN tournaments t ON str.tournamentid = t.tournamentid
WHERE t.startdate BETWEEN $1 AND $2
AND ($3::INTEGER IS NULL OR t.tournamentid = $3)
AND str.paymentstatus = 'paid'
GROUP BY s.schooltype
`

type GetSchoolPerformanceByCategoryParams struct {
	Startdate   time.Time `json:"startdate"`
	Startdate_2 time.Time `json:"startdate_2"`
	Column3     int32     `json:"column_3"`
}

type GetSchoolPerformanceByCategoryRow struct {
	GroupName   string      `json:"group_name"`
	TotalAmount interface{} `json:"total_amount"`
	SchoolCount int64       `json:"school_count"`
}

func (q *Queries) GetSchoolPerformanceByCategory(ctx context.Context, arg GetSchoolPerformanceByCategoryParams) ([]GetSchoolPerformanceByCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getSchoolPerformanceByCategory, arg.Startdate, arg.Startdate_2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetSchoolPerformanceByCategoryRow{}
	for rows.Next() {
		var i GetSchoolPerformanceByCategoryRow
		if err := rows.Scan(&i.GroupName, &i.TotalAmount, &i.SchoolCount); err != nil {
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

const getSchoolPerformanceByLocation = `-- name: GetSchoolPerformanceByLocation :many
SELECT
    CASE
        WHEN s.country = 'Rwanda' THEN s.province
        ELSE s.country
    END as group_name,
    COALESCE(SUM(str.actualpaidamount), 0) as total_amount,
    COUNT(DISTINCT s.schoolid) as school_count
FROM schools s
JOIN schooltournamentregistrations str ON s.schoolid = str.schoolid
JOIN tournaments t ON str.tournamentid = t.tournamentid
WHERE t.startdate BETWEEN $1 AND $2
AND ($3::INTEGER IS NULL OR t.tournamentid = $3)
AND str.paymentstatus = 'paid'
GROUP BY
    CASE
        WHEN s.country = 'Rwanda' THEN s.province
        ELSE s.country
    END
`

type GetSchoolPerformanceByLocationParams struct {
	Startdate   time.Time `json:"startdate"`
	Startdate_2 time.Time `json:"startdate_2"`
	Column3     int32     `json:"column_3"`
}

type GetSchoolPerformanceByLocationRow struct {
	GroupName   interface{} `json:"group_name"`
	TotalAmount interface{} `json:"total_amount"`
	SchoolCount int64       `json:"school_count"`
}

func (q *Queries) GetSchoolPerformanceByLocation(ctx context.Context, arg GetSchoolPerformanceByLocationParams) ([]GetSchoolPerformanceByLocationRow, error) {
	rows, err := q.db.QueryContext(ctx, getSchoolPerformanceByLocation, arg.Startdate, arg.Startdate_2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetSchoolPerformanceByLocationRow{}
	for rows.Next() {
		var i GetSchoolPerformanceByLocationRow
		if err := rows.Scan(&i.GroupName, &i.TotalAmount, &i.SchoolCount); err != nil {
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

const getTournamentIncomeOverview = `-- name: GetTournamentIncomeOverview :many
WITH TournamentIncome AS (
    SELECT
        t.tournamentid,
        t.name as tournament_name,
        t.leagueid,
        l.name as league_name,
        t.startdate,
        COALESCE(SUM(str.actualpaidamount), 0) as total_income,
        COALESCE(SUM(str.actualpaidamount - COALESCE(str.discountamount, 0)), 0) as net_revenue,
        COALESCE(SUM(str.actualpaidamount - COALESCE(str.discountamount, 0)), 0) -
        COALESCE((SELECT te.totalexpense
                  FROM tournamentexpenses te
                  WHERE te.tournamentid = t.tournamentid), 0) as net_profit
    FROM tournaments t
    LEFT JOIN leagues l ON t.leagueid = l.leagueid
    LEFT JOIN schooltournamentregistrations str ON t.tournamentid = str.tournamentid
    WHERE t.deleted_at IS NULL
    AND t.startdate BETWEEN $1 AND $2
    AND ($3::INTEGER IS NULL OR t.tournamentid = $3)
    AND str.paymentstatus = 'paid'
    GROUP BY t.tournamentid, t.name, t.leagueid, l.name, t.startdate
)
SELECT tournamentid, tournament_name, leagueid, league_name, startdate, total_income, net_revenue, net_profit FROM TournamentIncome
ORDER BY startdate DESC
`

type GetTournamentIncomeOverviewParams struct {
	Startdate   time.Time `json:"startdate"`
	Startdate_2 time.Time `json:"startdate_2"`
	Column3     int32     `json:"column_3"`
}

type GetTournamentIncomeOverviewRow struct {
	Tournamentid   int32          `json:"tournamentid"`
	TournamentName string         `json:"tournament_name"`
	Leagueid       sql.NullInt32  `json:"leagueid"`
	LeagueName     sql.NullString `json:"league_name"`
	Startdate      time.Time      `json:"startdate"`
	TotalIncome    interface{}    `json:"total_income"`
	NetRevenue     interface{}    `json:"net_revenue"`
	NetProfit      int32          `json:"net_profit"`
}

func (q *Queries) GetTournamentIncomeOverview(ctx context.Context, arg GetTournamentIncomeOverviewParams) ([]GetTournamentIncomeOverviewRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentIncomeOverview, arg.Startdate, arg.Startdate_2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTournamentIncomeOverviewRow{}
	for rows.Next() {
		var i GetTournamentIncomeOverviewRow
		if err := rows.Scan(
			&i.Tournamentid,
			&i.TournamentName,
			&i.Leagueid,
			&i.LeagueName,
			&i.Startdate,
			&i.TotalIncome,
			&i.NetRevenue,
			&i.NetProfit,
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
