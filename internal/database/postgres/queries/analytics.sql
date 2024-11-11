-- name: GetTournamentIncomeOverview :many
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
SELECT * FROM TournamentIncome
ORDER BY startdate DESC;

-- name: GetSchoolPerformanceByCategory :many
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
GROUP BY s.schooltype;

-- name: GetSchoolPerformanceByLocation :many
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
    END;

-- name: GetExpensesByTournament :many
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
ORDER BY t.startdate DESC;

-- name: GetExpensesSummary :one
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
AND ($3::INTEGER IS NULL OR t.tournamentid = $3);
