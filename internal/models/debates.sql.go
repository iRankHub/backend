// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: debates.sql

package models

import (
	"context"
	"database/sql"
	"time"
)

const addTeamMember = `-- name: AddTeamMember :one
INSERT INTO TeamMembers (TeamID, StudentID)
VALUES ($1, $2)
RETURNING TeamID, StudentID
`

type AddTeamMemberParams struct {
	Teamid    int32 `json:"teamid"`
	Studentid int32 `json:"studentid"`
}

func (q *Queries) AddTeamMember(ctx context.Context, arg AddTeamMemberParams) (Teammember, error) {
	row := q.db.QueryRowContext(ctx, addTeamMember, arg.Teamid, arg.Studentid)
	var i Teammember
	err := row.Scan(&i.Teamid, &i.Studentid)
	return i, err
}

const assignJudgeToDebate = `-- name: AssignJudgeToDebate :exec
INSERT INTO JudgeAssignments (TournamentID, JudgeID, DebateID, RoundNumber, IsElimination, IsHeadJudge)
VALUES ($1, $2, $3, $4, $5, $6)
`

type AssignJudgeToDebateParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Judgeid       int32 `json:"judgeid"`
	Debateid      int32 `json:"debateid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
	Isheadjudge   bool  `json:"isheadjudge"`
}

func (q *Queries) AssignJudgeToDebate(ctx context.Context, arg AssignJudgeToDebateParams) error {
	_, err := q.db.ExecContext(ctx, assignJudgeToDebate,
		arg.Tournamentid,
		arg.Judgeid,
		arg.Debateid,
		arg.Roundnumber,
		arg.Iselimination,
		arg.Isheadjudge,
	)
	return err
}

const assignRoomToDebate = `-- name: AssignRoomToDebate :exec
UPDATE Debates
SET RoomID = $2
WHERE DebateID = $1
`

type AssignRoomToDebateParams struct {
	Debateid int32 `json:"debateid"`
	Roomid   int32 `json:"roomid"`
}

func (q *Queries) AssignRoomToDebate(ctx context.Context, arg AssignRoomToDebateParams) error {
	_, err := q.db.ExecContext(ctx, assignRoomToDebate, arg.Debateid, arg.Roomid)
	return err
}

const checkExistingTeamMembership = `-- name: CheckExistingTeamMembership :one
SELECT COUNT(*) > 0 AS has_team
FROM TeamMembers tm
JOIN Teams t ON tm.TeamID = t.TeamID
WHERE t.TournamentID = $1 AND tm.StudentID = $2
`

type CheckExistingTeamMembershipParams struct {
	Tournamentid int32 `json:"tournamentid"`
	Studentid    int32 `json:"studentid"`
}

func (q *Queries) CheckExistingTeamMembership(ctx context.Context, arg CheckExistingTeamMembershipParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkExistingTeamMembership, arg.Tournamentid, arg.Studentid)
	var has_team bool
	err := row.Scan(&has_team)
	return has_team, err
}

const createDebate = `-- name: CreateDebate :one
INSERT INTO Debates (TournamentID, RoundID, RoundNumber, IsEliminationRound, Team1ID, Team2ID, RoomID, StartTime)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING DebateID
`

type CreateDebateParams struct {
	Tournamentid       int32     `json:"tournamentid"`
	Roundid            int32     `json:"roundid"`
	Roundnumber        int32     `json:"roundnumber"`
	Iseliminationround bool      `json:"iseliminationround"`
	Team1id            int32     `json:"team1id"`
	Team2id            int32     `json:"team2id"`
	Roomid             int32     `json:"roomid"`
	Starttime          time.Time `json:"starttime"`
}

func (q *Queries) CreateDebate(ctx context.Context, arg CreateDebateParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createDebate,
		arg.Tournamentid,
		arg.Roundid,
		arg.Roundnumber,
		arg.Iseliminationround,
		arg.Team1id,
		arg.Team2id,
		arg.Roomid,
		arg.Starttime,
	)
	var debateid int32
	err := row.Scan(&debateid)
	return debateid, err
}

const createPairingHistory = `-- name: CreatePairingHistory :exec
INSERT INTO PairingHistory (TournamentID, Team1ID, Team2ID, RoundNumber, IsElimination)
VALUES ($1, $2, $3, $4, $5)
`

type CreatePairingHistoryParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Team1id       int32 `json:"team1id"`
	Team2id       int32 `json:"team2id"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) CreatePairingHistory(ctx context.Context, arg CreatePairingHistoryParams) error {
	_, err := q.db.ExecContext(ctx, createPairingHistory,
		arg.Tournamentid,
		arg.Team1id,
		arg.Team2id,
		arg.Roundnumber,
		arg.Iselimination,
	)
	return err
}

const createRoom = `-- name: CreateRoom :one
INSERT INTO Rooms (RoomName, Location, Capacity, TournamentID)
VALUES ($1, $2, $3, $4)
RETURNING roomid, roomname, location, capacity, tournamentid
`

type CreateRoomParams struct {
	Roomname     string        `json:"roomname"`
	Location     string        `json:"location"`
	Capacity     int32         `json:"capacity"`
	Tournamentid sql.NullInt32 `json:"tournamentid"`
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom,
		arg.Roomname,
		arg.Location,
		arg.Capacity,
		arg.Tournamentid,
	)
	var i Room
	err := row.Scan(
		&i.Roomid,
		&i.Roomname,
		&i.Location,
		&i.Capacity,
		&i.Tournamentid,
	)
	return i, err
}

const createTeam = `-- name: CreateTeam :one
INSERT INTO Teams (Name, TournamentID)
VALUES ($1, $2)
RETURNING TeamID, Name, TournamentID
`

type CreateTeamParams struct {
	Name         string `json:"name"`
	Tournamentid int32  `json:"tournamentid"`
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, createTeam, arg.Name, arg.Tournamentid)
	var i Team
	err := row.Scan(&i.Teamid, &i.Name, &i.Tournamentid)
	return i, err
}

const deleteDebatesForTournament = `-- name: DeleteDebatesForTournament :exec
DELETE FROM Debates
WHERE TournamentID = $1
`

func (q *Queries) DeleteDebatesForTournament(ctx context.Context, tournamentid int32) error {
	_, err := q.db.ExecContext(ctx, deleteDebatesForTournament, tournamentid)
	return err
}

const deleteJudgeAssignmentsForTournament = `-- name: DeleteJudgeAssignmentsForTournament :exec
DELETE FROM JudgeAssignments
WHERE TournamentID = $1
`

func (q *Queries) DeleteJudgeAssignmentsForTournament(ctx context.Context, tournamentid int32) error {
	_, err := q.db.ExecContext(ctx, deleteJudgeAssignmentsForTournament, tournamentid)
	return err
}

const deletePairingHistoryForTournament = `-- name: DeletePairingHistoryForTournament :exec
DELETE FROM PairingHistory
WHERE TournamentID = $1
`

func (q *Queries) DeletePairingHistoryForTournament(ctx context.Context, tournamentid int32) error {
	_, err := q.db.ExecContext(ctx, deletePairingHistoryForTournament, tournamentid)
	return err
}

const deletePairingsForTournament = `-- name: DeletePairingsForTournament :exec
DELETE FROM Debates
WHERE TournamentID = $1
`

func (q *Queries) DeletePairingsForTournament(ctx context.Context, tournamentid int32) error {
	_, err := q.db.ExecContext(ctx, deletePairingsForTournament, tournamentid)
	return err
}

const deleteRoomsForTournament = `-- name: DeleteRoomsForTournament :exec
DELETE FROM Rooms
WHERE TournamentID = $1
`

func (q *Queries) DeleteRoomsForTournament(ctx context.Context, tournamentid sql.NullInt32) error {
	_, err := q.db.ExecContext(ctx, deleteRoomsForTournament, tournamentid)
	return err
}

const deleteRoundsForTournament = `-- name: DeleteRoundsForTournament :exec
DELETE FROM Rounds
WHERE TournamentID = $1
`

func (q *Queries) DeleteRoundsForTournament(ctx context.Context, tournamentid int32) error {
	_, err := q.db.ExecContext(ctx, deleteRoundsForTournament, tournamentid)
	return err
}

const deleteTeam = `-- name: DeleteTeam :exec
WITH debate_check AS (
    SELECT 1
    FROM Debates
    WHERE Team1ID = $1 OR Team2ID = $1
    LIMIT 1
)
DELETE FROM Teams
WHERE TeamID = $1 AND NOT EXISTS (SELECT 1 FROM debate_check)
`

func (q *Queries) DeleteTeam(ctx context.Context, teamid int32) error {
	_, err := q.db.ExecContext(ctx, deleteTeam, teamid)
	return err
}

const deleteTeamMembers = `-- name: DeleteTeamMembers :exec
DELETE FROM TeamMembers
WHERE TeamID = $1
`

func (q *Queries) DeleteTeamMembers(ctx context.Context, teamid int32) error {
	_, err := q.db.ExecContext(ctx, deleteTeamMembers, teamid)
	return err
}

const getAvailableJudges = `-- name: GetAvailableJudges :many
SELECT u.UserID, u.Name, u.Email
FROM Users u
JOIN Volunteers v ON u.UserID = v.UserID
LEFT JOIN JudgeAssignments ja ON u.UserID = ja.JudgeID AND ja.TournamentID = $1
WHERE v.Role = 'Judge' AND ja.JudgeID IS NULL
`

type GetAvailableJudgesRow struct {
	Userid int32  `json:"userid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (q *Queries) GetAvailableJudges(ctx context.Context, tournamentid int32) ([]GetAvailableJudgesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAvailableJudges, tournamentid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAvailableJudgesRow{}
	for rows.Next() {
		var i GetAvailableJudgesRow
		if err := rows.Scan(&i.Userid, &i.Name, &i.Email); err != nil {
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

const getBallotByID = `-- name: GetBallotByID :one
SELECT b.BallotID, d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.RoomID, r.roomname AS RoomName, b.JudgeID, u.Name AS JudgeName,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       b.Team1TotalScore, b.Team2TotalScore, b.RecordingStatus, b.Verdict
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Users u ON b.JudgeID = u.UserID
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
WHERE b.BallotID = $1
`

type GetBallotByIDRow struct {
	Ballotid           int32          `json:"ballotid"`
	Debateid           int32          `json:"debateid"`
	Roundnumber        int32          `json:"roundnumber"`
	Iseliminationround bool           `json:"iseliminationround"`
	Roomid             int32          `json:"roomid"`
	Roomname           sql.NullString `json:"roomname"`
	Judgeid            int32          `json:"judgeid"`
	Judgename          string         `json:"judgename"`
	Team1id            int32          `json:"team1id"`
	Team1name          string         `json:"team1name"`
	Team2id            int32          `json:"team2id"`
	Team2name          string         `json:"team2name"`
	Team1totalscore    sql.NullString `json:"team1totalscore"`
	Team2totalscore    sql.NullString `json:"team2totalscore"`
	Recordingstatus    string         `json:"recordingstatus"`
	Verdict            string         `json:"verdict"`
}

func (q *Queries) GetBallotByID(ctx context.Context, ballotid int32) (GetBallotByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getBallotByID, ballotid)
	var i GetBallotByIDRow
	err := row.Scan(
		&i.Ballotid,
		&i.Debateid,
		&i.Roundnumber,
		&i.Iseliminationround,
		&i.Roomid,
		&i.Roomname,
		&i.Judgeid,
		&i.Judgename,
		&i.Team1id,
		&i.Team1name,
		&i.Team2id,
		&i.Team2name,
		&i.Team1totalscore,
		&i.Team2totalscore,
		&i.Recordingstatus,
		&i.Verdict,
	)
	return i, err
}

const getBallotsByTournamentAndRound = `-- name: GetBallotsByTournamentAndRound :many
SELECT b.BallotID, d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.RoomID, r.roomname AS RoomName, b.JudgeID, u.Name AS JudgeName,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       b.Team1TotalScore, b.Team2TotalScore, b.RecordingStatus, b.Verdict
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Users u ON b.JudgeID = u.UserID
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
WHERE d.TournamentID = $1 AND d.RoundNumber = $2 AND d.IsEliminationRound = $3
`

type GetBallotsByTournamentAndRoundParams struct {
	Tournamentid       int32 `json:"tournamentid"`
	Roundnumber        int32 `json:"roundnumber"`
	Iseliminationround bool  `json:"iseliminationround"`
}

type GetBallotsByTournamentAndRoundRow struct {
	Ballotid           int32          `json:"ballotid"`
	Debateid           int32          `json:"debateid"`
	Roundnumber        int32          `json:"roundnumber"`
	Iseliminationround bool           `json:"iseliminationround"`
	Roomid             int32          `json:"roomid"`
	Roomname           sql.NullString `json:"roomname"`
	Judgeid            int32          `json:"judgeid"`
	Judgename          string         `json:"judgename"`
	Team1id            int32          `json:"team1id"`
	Team1name          string         `json:"team1name"`
	Team2id            int32          `json:"team2id"`
	Team2name          string         `json:"team2name"`
	Team1totalscore    sql.NullString `json:"team1totalscore"`
	Team2totalscore    sql.NullString `json:"team2totalscore"`
	Recordingstatus    string         `json:"recordingstatus"`
	Verdict            string         `json:"verdict"`
}

func (q *Queries) GetBallotsByTournamentAndRound(ctx context.Context, arg GetBallotsByTournamentAndRoundParams) ([]GetBallotsByTournamentAndRoundRow, error) {
	rows, err := q.db.QueryContext(ctx, getBallotsByTournamentAndRound, arg.Tournamentid, arg.Roundnumber, arg.Iseliminationround)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetBallotsByTournamentAndRoundRow{}
	for rows.Next() {
		var i GetBallotsByTournamentAndRoundRow
		if err := rows.Scan(
			&i.Ballotid,
			&i.Debateid,
			&i.Roundnumber,
			&i.Iseliminationround,
			&i.Roomid,
			&i.Roomname,
			&i.Judgeid,
			&i.Judgename,
			&i.Team1id,
			&i.Team1name,
			&i.Team2id,
			&i.Team2name,
			&i.Team1totalscore,
			&i.Team2totalscore,
			&i.Recordingstatus,
			&i.Verdict,
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

const getDebateByRoomAndRound = `-- name: GetDebateByRoomAndRound :one
SELECT debateid, roundid, roundnumber, iseliminationround, tournamentid, team1id, team2id, starttime, endtime, roomid, status
FROM Debates
WHERE TournamentID = $1 AND RoomID = $2 AND RoundNumber = $3 AND IsEliminationRound = $4
LIMIT 1
`

type GetDebateByRoomAndRoundParams struct {
	Tournamentid       int32 `json:"tournamentid"`
	Roomid             int32 `json:"roomid"`
	Roundnumber        int32 `json:"roundnumber"`
	Iseliminationround bool  `json:"iseliminationround"`
}

func (q *Queries) GetDebateByRoomAndRound(ctx context.Context, arg GetDebateByRoomAndRoundParams) (Debate, error) {
	row := q.db.QueryRowContext(ctx, getDebateByRoomAndRound,
		arg.Tournamentid,
		arg.Roomid,
		arg.Roundnumber,
		arg.Iseliminationround,
	)
	var i Debate
	err := row.Scan(
		&i.Debateid,
		&i.Roundid,
		&i.Roundnumber,
		&i.Iseliminationround,
		&i.Tournamentid,
		&i.Team1id,
		&i.Team2id,
		&i.Starttime,
		&i.Endtime,
		&i.Roomid,
		&i.Status,
	)
	return i, err
}

const getDebatesByRoomAndTournament = `-- name: GetDebatesByRoomAndTournament :many
SELECT debateid, roundid, roundnumber, iseliminationround, tournamentid, team1id, team2id, starttime, endtime, roomid, status
FROM Debates
WHERE TournamentID = $1 AND RoomID = $2 AND IsEliminationRound = $3
`

type GetDebatesByRoomAndTournamentParams struct {
	Tournamentid       int32 `json:"tournamentid"`
	Roomid             int32 `json:"roomid"`
	Iseliminationround bool  `json:"iseliminationround"`
}

func (q *Queries) GetDebatesByRoomAndTournament(ctx context.Context, arg GetDebatesByRoomAndTournamentParams) ([]Debate, error) {
	rows, err := q.db.QueryContext(ctx, getDebatesByRoomAndTournament, arg.Tournamentid, arg.Roomid, arg.Iseliminationround)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Debate{}
	for rows.Next() {
		var i Debate
		if err := rows.Scan(
			&i.Debateid,
			&i.Roundid,
			&i.Roundnumber,
			&i.Iseliminationround,
			&i.Tournamentid,
			&i.Team1id,
			&i.Team2id,
			&i.Starttime,
			&i.Endtime,
			&i.Roomid,
			&i.Status,
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

const getJudgeByID = `-- name: GetJudgeByID :one
SELECT u.UserID, u.Name, u.Email
FROM Users u
WHERE u.UserID = $1
`

type GetJudgeByIDRow struct {
	Userid int32  `json:"userid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (q *Queries) GetJudgeByID(ctx context.Context, userid int32) (GetJudgeByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getJudgeByID, userid)
	var i GetJudgeByIDRow
	err := row.Scan(&i.Userid, &i.Name, &i.Email)
	return i, err
}

const getJudgesByTournamentAndRound = `-- name: GetJudgesByTournamentAndRound :many
SELECT u.UserID, u.Name, u.Email, ja.IsHeadJudge
FROM Users u
JOIN JudgeAssignments ja ON u.UserID = ja.JudgeID
WHERE ja.TournamentID = $1 AND ja.RoundNumber = $2 AND ja.IsElimination = $3
`

type GetJudgesByTournamentAndRoundParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

type GetJudgesByTournamentAndRoundRow struct {
	Userid      int32  `json:"userid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Isheadjudge bool   `json:"isheadjudge"`
}

func (q *Queries) GetJudgesByTournamentAndRound(ctx context.Context, arg GetJudgesByTournamentAndRoundParams) ([]GetJudgesByTournamentAndRoundRow, error) {
	rows, err := q.db.QueryContext(ctx, getJudgesByTournamentAndRound, arg.Tournamentid, arg.Roundnumber, arg.Iselimination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetJudgesByTournamentAndRoundRow{}
	for rows.Next() {
		var i GetJudgesByTournamentAndRoundRow
		if err := rows.Scan(
			&i.Userid,
			&i.Name,
			&i.Email,
			&i.Isheadjudge,
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

const getPairingByID = `-- name: GetPairingByID :one
SELECT d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       d.RoomID, r.roomname AS RoomName,
       array_agg(DISTINCT s1.FirstName || ' ' || s1.LastName) AS Team1SpeakerNames,
       array_agg(DISTINCT s2.FirstName || ' ' || s2.LastName) AS Team2SpeakerNames,
       l1.Name AS Team1LeagueName, l2.Name AS Team2LeagueName,
       COALESCE(t1_points.TotalPoints, 0) AS Team1TotalPoints,
       COALESCE(t2_points.TotalPoints, 0) AS Team2TotalPoints,
       (SELECT u.Name FROM JudgeAssignments ja
        JOIN Users u ON ja.JudgeID = u.UserID
        WHERE ja.DebateID = d.DebateID AND ja.IsHeadJudge = true
        LIMIT 1) AS HeadJudgeName
FROM Debates d
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
LEFT JOIN TeamMembers tm1 ON t1.TeamID = tm1.TeamID
LEFT JOIN TeamMembers tm2 ON t2.TeamID = tm2.TeamID
LEFT JOIN Students s1 ON tm1.StudentID = s1.StudentID
LEFT JOIN Students s2 ON tm2.StudentID = s2.StudentID
LEFT JOIN Tournaments tour ON d.TournamentID = tour.TournamentID
LEFT JOIN Leagues l1 ON tour.LeagueID = l1.LeagueID
LEFT JOIN Leagues l2 ON tour.LeagueID = l2.LeagueID
LEFT JOIN (
    SELECT Team1ID AS TeamID, SUM(Team1TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = (SELECT TournamentID FROM Debates WHERE d.DebateID = $1)
    GROUP BY Team1ID
    UNION ALL
    SELECT Team2ID AS TeamID, SUM(Team2TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = (SELECT TournamentID FROM Debates WHERE d.DebateID = $1)
    GROUP BY Team2ID
) t1_points ON t1.TeamID = t1_points.TeamID
LEFT JOIN (
    SELECT Team1ID AS TeamID, SUM(Team1TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = (SELECT TournamentID FROM Debates WHERE d.DebateID = $1)
    GROUP BY Team1ID
    UNION ALL
    SELECT Team2ID AS TeamID, SUM(Team2TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = (SELECT TournamentID FROM Debates WHERE d.DebateID = $1)
    GROUP BY Team2ID
) t2_points ON t2.TeamID = t2_points.TeamID
WHERE d.DebateID = $1
GROUP BY d.DebateID, d.RoundNumber, d.IsEliminationRound, d.Team1ID, t1.Name, d.Team2ID, t2.Name, d.RoomID, r.RoomName,
         l1.Name, l2.Name, t1_points.TotalPoints, t2_points.TotalPoints
`

type GetPairingByIDRow struct {
	Debateid           int32          `json:"debateid"`
	Roundnumber        int32          `json:"roundnumber"`
	Iseliminationround bool           `json:"iseliminationround"`
	Team1id            int32          `json:"team1id"`
	Team1name          string         `json:"team1name"`
	Team2id            int32          `json:"team2id"`
	Team2name          string         `json:"team2name"`
	Roomid             int32          `json:"roomid"`
	Roomname           sql.NullString `json:"roomname"`
	Team1speakernames  interface{}    `json:"team1speakernames"`
	Team2speakernames  interface{}    `json:"team2speakernames"`
	Team1leaguename    sql.NullString `json:"team1leaguename"`
	Team2leaguename    sql.NullString `json:"team2leaguename"`
	Team1totalpoints   int64          `json:"team1totalpoints"`
	Team2totalpoints   int64          `json:"team2totalpoints"`
	Headjudgename      string         `json:"headjudgename"`
}

func (q *Queries) GetPairingByID(ctx context.Context, debateid int32) (GetPairingByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getPairingByID, debateid)
	var i GetPairingByIDRow
	err := row.Scan(
		&i.Debateid,
		&i.Roundnumber,
		&i.Iseliminationround,
		&i.Team1id,
		&i.Team1name,
		&i.Team2id,
		&i.Team2name,
		&i.Roomid,
		&i.Roomname,
		&i.Team1speakernames,
		&i.Team2speakernames,
		&i.Team1leaguename,
		&i.Team2leaguename,
		&i.Team1totalpoints,
		&i.Team2totalpoints,
		&i.Headjudgename,
	)
	return i, err
}

const getPairingsByTournamentAndRound = `-- name: GetPairingsByTournamentAndRound :many
SELECT d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       d.RoomID, r.roomname AS RoomName,
       array_agg(DISTINCT s1.FirstName || ' ' || s1.LastName) AS Team1SpeakerNames,
       array_agg(DISTINCT s2.FirstName || ' ' || s2.LastName) AS Team2SpeakerNames,
       l1.Name AS Team1LeagueName, l2.Name AS Team2LeagueName,
       COALESCE(t1_points.TotalPoints, 0) AS Team1TotalPoints,
       COALESCE(t2_points.TotalPoints, 0) AS Team2TotalPoints,
       (SELECT u.Name FROM JudgeAssignments ja
        JOIN Users u ON ja.JudgeID = u.UserID
        WHERE ja.DebateID = d.DebateID AND ja.IsHeadJudge = true
        LIMIT 1) AS HeadJudgeName
FROM Debates d
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
LEFT JOIN TeamMembers tm1 ON t1.TeamID = tm1.TeamID
LEFT JOIN TeamMembers tm2 ON t2.TeamID = tm2.TeamID
LEFT JOIN Students s1 ON tm1.StudentID = s1.StudentID
LEFT JOIN Students s2 ON tm2.StudentID = s2.StudentID
LEFT JOIN Tournaments tour ON d.TournamentID = tour.TournamentID
LEFT JOIN Leagues l1 ON tour.LeagueID = l1.LeagueID
LEFT JOIN Leagues l2 ON tour.LeagueID = l2.LeagueID
LEFT JOIN (
    SELECT Team1ID AS TeamID, SUM(Team1TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = $1
    GROUP BY Team1ID
    UNION ALL
    SELECT Team2ID AS TeamID, SUM(Team2TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = $1
    GROUP BY Team2ID
) t1_points ON t1.TeamID = t1_points.TeamID
LEFT JOIN (
    SELECT Team1ID AS TeamID, SUM(Team1TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = $1
    GROUP BY Team1ID
    UNION ALL
    SELECT Team2ID AS TeamID, SUM(Team2TotalScore) AS TotalPoints
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = $1
    GROUP BY Team2ID
) t2_points ON t2.TeamID = t2_points.TeamID
WHERE d.TournamentID = $1 AND d.RoundNumber = $2 AND d.IsEliminationRound = $3
GROUP BY d.DebateID, d.RoundNumber, d.IsEliminationRound, d.Team1ID, t1.Name, d.Team2ID, t2.Name, d.RoomID, r.RoomName,
         l1.Name, l2.Name, t1_points.TotalPoints, t2_points.TotalPoints
`

type GetPairingsByTournamentAndRoundParams struct {
	Tournamentid       int32 `json:"tournamentid"`
	Roundnumber        int32 `json:"roundnumber"`
	Iseliminationround bool  `json:"iseliminationround"`
}

type GetPairingsByTournamentAndRoundRow struct {
	Debateid           int32          `json:"debateid"`
	Roundnumber        int32          `json:"roundnumber"`
	Iseliminationround bool           `json:"iseliminationround"`
	Team1id            int32          `json:"team1id"`
	Team1name          string         `json:"team1name"`
	Team2id            int32          `json:"team2id"`
	Team2name          string         `json:"team2name"`
	Roomid             int32          `json:"roomid"`
	Roomname           sql.NullString `json:"roomname"`
	Team1speakernames  interface{}    `json:"team1speakernames"`
	Team2speakernames  interface{}    `json:"team2speakernames"`
	Team1leaguename    sql.NullString `json:"team1leaguename"`
	Team2leaguename    sql.NullString `json:"team2leaguename"`
	Team1totalpoints   int64          `json:"team1totalpoints"`
	Team2totalpoints   int64          `json:"team2totalpoints"`
	Headjudgename      string         `json:"headjudgename"`
}

func (q *Queries) GetPairingsByTournamentAndRound(ctx context.Context, arg GetPairingsByTournamentAndRoundParams) ([]GetPairingsByTournamentAndRoundRow, error) {
	rows, err := q.db.QueryContext(ctx, getPairingsByTournamentAndRound, arg.Tournamentid, arg.Roundnumber, arg.Iseliminationround)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPairingsByTournamentAndRoundRow{}
	for rows.Next() {
		var i GetPairingsByTournamentAndRoundRow
		if err := rows.Scan(
			&i.Debateid,
			&i.Roundnumber,
			&i.Iseliminationround,
			&i.Team1id,
			&i.Team1name,
			&i.Team2id,
			&i.Team2name,
			&i.Roomid,
			&i.Roomname,
			&i.Team1speakernames,
			&i.Team2speakernames,
			&i.Team1leaguename,
			&i.Team2leaguename,
			&i.Team1totalpoints,
			&i.Team2totalpoints,
			&i.Headjudgename,
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

const getPreviousPairings = `-- name: GetPreviousPairings :many
SELECT Team1ID, Team2ID
FROM Debates
WHERE TournamentID = $1 AND RoundNumber < $2
`

type GetPreviousPairingsParams struct {
	Tournamentid int32 `json:"tournamentid"`
	Roundnumber  int32 `json:"roundnumber"`
}

type GetPreviousPairingsRow struct {
	Team1id int32 `json:"team1id"`
	Team2id int32 `json:"team2id"`
}

func (q *Queries) GetPreviousPairings(ctx context.Context, arg GetPreviousPairingsParams) ([]GetPreviousPairingsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPreviousPairings, arg.Tournamentid, arg.Roundnumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPreviousPairingsRow{}
	for rows.Next() {
		var i GetPreviousPairingsRow
		if err := rows.Scan(&i.Team1id, &i.Team2id); err != nil {
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

const getRoomByID = `-- name: GetRoomByID :one
SELECT RoomID, RoomName, TournamentID, Location, Capacity
FROM Rooms
WHERE RoomID = $1
`

type GetRoomByIDRow struct {
	Roomid       int32         `json:"roomid"`
	Roomname     string        `json:"roomname"`
	Tournamentid sql.NullInt32 `json:"tournamentid"`
	Location     string        `json:"location"`
	Capacity     int32         `json:"capacity"`
}

func (q *Queries) GetRoomByID(ctx context.Context, roomid int32) (GetRoomByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getRoomByID, roomid)
	var i GetRoomByIDRow
	err := row.Scan(
		&i.Roomid,
		&i.Roomname,
		&i.Tournamentid,
		&i.Location,
		&i.Capacity,
	)
	return i, err
}

const getRoomsByTournament = `-- name: GetRoomsByTournament :many
SELECT roomid, roomname, location, capacity, tournamentid FROM Rooms
WHERE TournamentID = $1
`

func (q *Queries) GetRoomsByTournament(ctx context.Context, tournamentid sql.NullInt32) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getRoomsByTournament, tournamentid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Room{}
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.Roomid,
			&i.Roomname,
			&i.Location,
			&i.Capacity,
			&i.Tournamentid,
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

const getRoundByTournamentAndNumber = `-- name: GetRoundByTournamentAndNumber :one
SELECT roundid, tournamentid, roundnumber, iseliminationround FROM Rounds
WHERE TournamentID = $1 AND RoundNumber = $2 AND IsEliminationRound = $3
LIMIT 1
`

type GetRoundByTournamentAndNumberParams struct {
	Tournamentid       int32 `json:"tournamentid"`
	Roundnumber        int32 `json:"roundnumber"`
	Iseliminationround bool  `json:"iseliminationround"`
}

func (q *Queries) GetRoundByTournamentAndNumber(ctx context.Context, arg GetRoundByTournamentAndNumberParams) (Round, error) {
	row := q.db.QueryRowContext(ctx, getRoundByTournamentAndNumber, arg.Tournamentid, arg.Roundnumber, arg.Iseliminationround)
	var i Round
	err := row.Scan(
		&i.Roundid,
		&i.Tournamentid,
		&i.Roundnumber,
		&i.Iseliminationround,
	)
	return i, err
}

const getSpeakerScoresByBallot = `-- name: GetSpeakerScoresByBallot :many
SELECT ss.ScoreID, ss.SpeakerID, s.FirstName, s.LastName, ss.SpeakerRank, ss.SpeakerPoints, ss.Feedback
FROM SpeakerScores ss
JOIN Students s ON ss.SpeakerID = s.StudentID
WHERE ss.BallotID = $1
`

type GetSpeakerScoresByBallotRow struct {
	Scoreid       int32          `json:"scoreid"`
	Speakerid     int32          `json:"speakerid"`
	Firstname     string         `json:"firstname"`
	Lastname      string         `json:"lastname"`
	Speakerrank   int32          `json:"speakerrank"`
	Speakerpoints string         `json:"speakerpoints"`
	Feedback      sql.NullString `json:"feedback"`
}

func (q *Queries) GetSpeakerScoresByBallot(ctx context.Context, ballotid int32) ([]GetSpeakerScoresByBallotRow, error) {
	rows, err := q.db.QueryContext(ctx, getSpeakerScoresByBallot, ballotid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetSpeakerScoresByBallotRow{}
	for rows.Next() {
		var i GetSpeakerScoresByBallotRow
		if err := rows.Scan(
			&i.Scoreid,
			&i.Speakerid,
			&i.Firstname,
			&i.Lastname,
			&i.Speakerrank,
			&i.Speakerpoints,
			&i.Feedback,
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

const getTeamByID = `-- name: GetTeamByID :one
SELECT t.TeamID, t.Name, t.TournamentID,
       array_agg(tm.StudentID) as SpeakerIDs
FROM Teams t
LEFT JOIN TeamMembers tm ON t.TeamID = tm.TeamID
WHERE t.TeamID = $1
GROUP BY t.TeamID, t.Name, t.TournamentID
`

type GetTeamByIDRow struct {
	Teamid       int32       `json:"teamid"`
	Name         string      `json:"name"`
	Tournamentid int32       `json:"tournamentid"`
	Speakerids   interface{} `json:"speakerids"`
}

func (q *Queries) GetTeamByID(ctx context.Context, teamid int32) (GetTeamByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getTeamByID, teamid)
	var i GetTeamByIDRow
	err := row.Scan(
		&i.Teamid,
		&i.Name,
		&i.Tournamentid,
		&i.Speakerids,
	)
	return i, err
}

const getTeamMembers = `-- name: GetTeamMembers :many
SELECT tm.TeamID, tm.StudentID, s.FirstName, s.LastName
FROM TeamMembers tm
JOIN Students s ON tm.StudentID = s.StudentID
WHERE tm.TeamID = $1
`

type GetTeamMembersRow struct {
	Teamid    int32  `json:"teamid"`
	Studentid int32  `json:"studentid"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (q *Queries) GetTeamMembers(ctx context.Context, teamid int32) ([]GetTeamMembersRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeamMembers, teamid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTeamMembersRow{}
	for rows.Next() {
		var i GetTeamMembersRow
		if err := rows.Scan(
			&i.Teamid,
			&i.Studentid,
			&i.Firstname,
			&i.Lastname,
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

const getTeamWins = `-- name: GetTeamWins :one
SELECT COUNT(*) as wins
FROM Debates d
JOIN Ballots b ON d.DebateID = b.DebateID
WHERE (d.Team1ID = $1 AND b.Team1TotalScore > b.Team2TotalScore)
   OR (d.Team2ID = $1 AND b.Team2TotalScore > b.Team1TotalScore)
   AND d.TournamentID = $2
`

type GetTeamWinsParams struct {
	Team1id      int32 `json:"team1id"`
	Tournamentid int32 `json:"tournamentid"`
}

func (q *Queries) GetTeamWins(ctx context.Context, arg GetTeamWinsParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTeamWins, arg.Team1id, arg.Tournamentid)
	var wins int64
	err := row.Scan(&wins)
	return wins, err
}

const getTeamsByTournament = `-- name: GetTeamsByTournament :many
SELECT t.TeamID, t.Name, t.TournamentID,
       array_agg(tm.StudentID) as SpeakerIDs,
       (SELECT COUNT(*)
        FROM Debates d
        JOIN Ballots b ON d.DebateID = b.DebateID
        WHERE ((d.Team1ID = t.TeamID AND b.Team1TotalScore > b.Team2TotalScore)
           OR (d.Team2ID = t.TeamID AND b.Team2TotalScore > b.Team1TotalScore))
           AND d.TournamentID = $1) as Wins,
       l.Name as LeagueName
FROM Teams t
LEFT JOIN TeamMembers tm ON t.TeamID = tm.TeamID
JOIN Tournaments tour ON t.TournamentID = tour.TournamentID
JOIN Leagues l ON tour.LeagueID = l.LeagueID
WHERE t.TournamentID = $1
GROUP BY t.TeamID, t.Name, t.TournamentID, l.Name
`

type GetTeamsByTournamentRow struct {
	Teamid       int32       `json:"teamid"`
	Name         string      `json:"name"`
	Tournamentid int32       `json:"tournamentid"`
	Speakerids   interface{} `json:"speakerids"`
	Wins         int64       `json:"wins"`
	Leaguename   string      `json:"leaguename"`
}

func (q *Queries) GetTeamsByTournament(ctx context.Context, tournamentid int32) ([]GetTeamsByTournamentRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeamsByTournament, tournamentid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTeamsByTournamentRow{}
	for rows.Next() {
		var i GetTeamsByTournamentRow
		if err := rows.Scan(
			&i.Teamid,
			&i.Name,
			&i.Tournamentid,
			&i.Speakerids,
			&i.Wins,
			&i.Leaguename,
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

const removeTeamMembers = `-- name: RemoveTeamMembers :exec
DELETE FROM TeamMembers
WHERE TeamID = $1
`

func (q *Queries) RemoveTeamMembers(ctx context.Context, teamid int32) error {
	_, err := q.db.ExecContext(ctx, removeTeamMembers, teamid)
	return err
}

const updateBallot = `-- name: UpdateBallot :exec
UPDATE Ballots
SET Team1TotalScore = $2, Team2TotalScore = $3, RecordingStatus = $4, Verdict = $5
WHERE BallotID = $1
`

type UpdateBallotParams struct {
	Ballotid        int32          `json:"ballotid"`
	Team1totalscore sql.NullString `json:"team1totalscore"`
	Team2totalscore sql.NullString `json:"team2totalscore"`
	Recordingstatus string         `json:"recordingstatus"`
	Verdict         string         `json:"verdict"`
}

func (q *Queries) UpdateBallot(ctx context.Context, arg UpdateBallotParams) error {
	_, err := q.db.ExecContext(ctx, updateBallot,
		arg.Ballotid,
		arg.Team1totalscore,
		arg.Team2totalscore,
		arg.Recordingstatus,
		arg.Verdict,
	)
	return err
}

const updatePairing = `-- name: UpdatePairing :exec
UPDATE Debates
SET Team1ID = $2, Team2ID = $3, RoomID = $4
WHERE DebateID = $1
`

type UpdatePairingParams struct {
	Debateid int32 `json:"debateid"`
	Team1id  int32 `json:"team1id"`
	Team2id  int32 `json:"team2id"`
	Roomid   int32 `json:"roomid"`
}

func (q *Queries) UpdatePairing(ctx context.Context, arg UpdatePairingParams) error {
	_, err := q.db.ExecContext(ctx, updatePairing,
		arg.Debateid,
		arg.Team1id,
		arg.Team2id,
		arg.Roomid,
	)
	return err
}

const updateRoom = `-- name: UpdateRoom :one
UPDATE Rooms
SET RoomName = $2
WHERE RoomID = $1
RETURNING roomid, roomname, location, capacity, tournamentid
`

type UpdateRoomParams struct {
	Roomid   int32  `json:"roomid"`
	Roomname string `json:"roomname"`
}

func (q *Queries) UpdateRoom(ctx context.Context, arg UpdateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, updateRoom, arg.Roomid, arg.Roomname)
	var i Room
	err := row.Scan(
		&i.Roomid,
		&i.Roomname,
		&i.Location,
		&i.Capacity,
		&i.Tournamentid,
	)
	return i, err
}

const updateSpeakerScore = `-- name: UpdateSpeakerScore :exec
UPDATE SpeakerScores
SET SpeakerRank = $2, SpeakerPoints = $3, Feedback = $4
WHERE ScoreID = $1
`

type UpdateSpeakerScoreParams struct {
	Scoreid       int32          `json:"scoreid"`
	Speakerrank   int32          `json:"speakerrank"`
	Speakerpoints string         `json:"speakerpoints"`
	Feedback      sql.NullString `json:"feedback"`
}

func (q *Queries) UpdateSpeakerScore(ctx context.Context, arg UpdateSpeakerScoreParams) error {
	_, err := q.db.ExecContext(ctx, updateSpeakerScore,
		arg.Scoreid,
		arg.Speakerrank,
		arg.Speakerpoints,
		arg.Feedback,
	)
	return err
}

const updateTeam = `-- name: UpdateTeam :exec
UPDATE Teams
SET Name = $2
WHERE TeamID = $1
`

type UpdateTeamParams struct {
	Teamid int32  `json:"teamid"`
	Name   string `json:"name"`
}

func (q *Queries) UpdateTeam(ctx context.Context, arg UpdateTeamParams) error {
	_, err := q.db.ExecContext(ctx, updateTeam, arg.Teamid, arg.Name)
	return err
}
