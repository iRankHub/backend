// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: judges.sql

package models

import (
	"context"
	"database/sql"
)

const assignJudgeToPosition = `-- name: AssignJudgeToPosition :exec
WITH judge_b_info AS (
    SELECT ja.debateid, ja.isheadjudge
    FROM judgeassignments ja
    WHERE ja.judgeid = $2
      AND ja.tournamentid = $3
      AND ja.roundnumber = $4
      AND ja.iselimination = $5
)
INSERT INTO judgeassignments (
    judgeid,
    tournamentid,
    debateid,
    roundnumber,
    iselimination,
    isheadjudge
)
SELECT
    $1,
    $3,
    debateid,
    $4,
    $5,
    isheadjudge
FROM judge_b_info
`

type AssignJudgeToPositionParams struct {
	Judgeid       int32 `json:"judgeid"`
	Judgeid_2     int32 `json:"judgeid_2"`
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) AssignJudgeToPosition(ctx context.Context, arg AssignJudgeToPositionParams) error {
	_, err := q.db.ExecContext(ctx, assignJudgeToPosition,
		arg.Judgeid,
		arg.Judgeid_2,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Iselimination,
	)
	return err
}

const checkHeadJudgeExists = `-- name: CheckHeadJudgeExists :one
SELECT EXISTS (
    SELECT 1
    FROM judgeassignments ja
    JOIN debates d ON ja.debateid = d.debateid
    WHERE ja.tournamentid = $1
    AND d.roomid = $2
    AND ja.roundnumber = $3
    AND ja.iselimination = $4
    AND ja.isheadjudge = true
)
`

type CheckHeadJudgeExistsParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Roomid        int32 `json:"roomid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) CheckHeadJudgeExists(ctx context.Context, arg CheckHeadJudgeExistsParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkHeadJudgeExists,
		arg.Tournamentid,
		arg.Roomid,
		arg.Roundnumber,
		arg.Iselimination,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const demoteCurrentHeadJudge = `-- name: DemoteCurrentHeadJudge :exec
UPDATE judgeassignments ja
SET isheadjudge = false
WHERE ja.tournamentid = $1
AND ja.roundnumber = $2
AND ja.debateid IN (
    SELECT d.debateid
    FROM debates d
    WHERE d.roomid = $3
    AND d.tournamentid = $1
)
AND ja.iselimination = $4
AND ja.isheadjudge = true
`

type DemoteCurrentHeadJudgeParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Roomid        int32 `json:"roomid"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) DemoteCurrentHeadJudge(ctx context.Context, arg DemoteCurrentHeadJudgeParams) error {
	_, err := q.db.ExecContext(ctx, demoteCurrentHeadJudge,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Roomid,
		arg.Iselimination,
	)
	return err
}

const ensureUnassignedRoomExists = `-- name: EnsureUnassignedRoomExists :one
WITH existing_room AS (
    SELECT r.roomid
    FROM rooms r
    WHERE r.roomname = 'Unassigned'
      AND r.tournamentid = $1
),
     inserted_room AS (
         INSERT INTO rooms (roomname, location, capacity, tournamentid)
             SELECT 'Unassigned', 'N/A', 100, $1
             WHERE NOT EXISTS (SELECT 1 FROM existing_room)
             RETURNING roomid
     )
SELECT er.roomid FROM existing_room er
UNION ALL
SELECT ir.roomid FROM inserted_room ir
`

func (q *Queries) EnsureUnassignedRoomExists(ctx context.Context, tournamentid sql.NullInt32) (int32, error) {
	row := q.db.QueryRowContext(ctx, ensureUnassignedRoomExists, tournamentid)
	var roomid int32
	err := row.Scan(&roomid)
	return roomid, err
}

const getEligibleHeadJudge = `-- name: GetEligibleHeadJudge :one
SELECT ja.judgeid
FROM judgeassignments ja
JOIN debates d ON ja.debateid = d.debateid
WHERE ja.tournamentid = $1
AND ja.roundnumber = $2
AND d.roomid = $3
AND ja.iselimination = $4
AND ja.judgeid != $5
AND NOT ja.isheadjudge
LIMIT 1
`

type GetEligibleHeadJudgeParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Roomid        int32 `json:"roomid"`
	Iselimination bool  `json:"iselimination"`
	Judgeid       int32 `json:"judgeid"`
}

func (q *Queries) GetEligibleHeadJudge(ctx context.Context, arg GetEligibleHeadJudgeParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, getEligibleHeadJudge,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Roomid,
		arg.Iselimination,
		arg.Judgeid,
	)
	var judgeid int32
	err := row.Scan(&judgeid)
	return judgeid, err
}

const getJudgeAssignment = `-- name: GetJudgeAssignment :one
SELECT ja.assignmentid, ja.tournamentid, ja.judgeid, ja.debateid, ja.roundnumber, ja.iselimination, ja.isheadjudge, d.roomid as roomid, d.debateid
FROM judgeassignments ja
JOIN debates d ON ja.debateid = d.debateid
WHERE ja.tournamentid = $1
AND ja.judgeid = $2
AND ja.roundnumber = $3
AND ja.iselimination = $4
`

type GetJudgeAssignmentParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Judgeid       int32 `json:"judgeid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

type GetJudgeAssignmentRow struct {
	Assignmentid  int32 `json:"assignmentid"`
	Tournamentid  int32 `json:"tournamentid"`
	Judgeid       int32 `json:"judgeid"`
	Debateid      int32 `json:"debateid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
	Isheadjudge   bool  `json:"isheadjudge"`
	Roomid        int32 `json:"roomid"`
	Debateid_2    int32 `json:"debateid_2"`
}

func (q *Queries) GetJudgeAssignment(ctx context.Context, arg GetJudgeAssignmentParams) (GetJudgeAssignmentRow, error) {
	row := q.db.QueryRowContext(ctx, getJudgeAssignment,
		arg.Tournamentid,
		arg.Judgeid,
		arg.Roundnumber,
		arg.Iselimination,
	)
	var i GetJudgeAssignmentRow
	err := row.Scan(
		&i.Assignmentid,
		&i.Tournamentid,
		&i.Judgeid,
		&i.Debateid,
		&i.Roundnumber,
		&i.Iselimination,
		&i.Isheadjudge,
		&i.Roomid,
		&i.Debateid_2,
	)
	return i, err
}

const getJudgeInRoom = `-- name: GetJudgeInRoom :one
SELECT ja.assignmentid, ja.tournamentid, ja.judgeid, ja.debateid, ja.roundnumber, ja.iselimination, ja.isheadjudge, d.roomid as roomid, d.debateid
FROM judgeassignments ja
         JOIN debates d ON ja.debateid = d.debateid
WHERE ja.tournamentid = $1
  AND d.roomid = $2
  AND ja.roundnumber = $3
  AND ja.iselimination = $4
LIMIT 1
`

type GetJudgeInRoomParams struct {
	Tournamentid  int32 `json:"tournamentid"`
	Roomid        int32 `json:"roomid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

type GetJudgeInRoomRow struct {
	Assignmentid  int32 `json:"assignmentid"`
	Tournamentid  int32 `json:"tournamentid"`
	Judgeid       int32 `json:"judgeid"`
	Debateid      int32 `json:"debateid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
	Isheadjudge   bool  `json:"isheadjudge"`
	Roomid        int32 `json:"roomid"`
	Debateid_2    int32 `json:"debateid_2"`
}

func (q *Queries) GetJudgeInRoom(ctx context.Context, arg GetJudgeInRoomParams) (GetJudgeInRoomRow, error) {
	row := q.db.QueryRowContext(ctx, getJudgeInRoom,
		arg.Tournamentid,
		arg.Roomid,
		arg.Roundnumber,
		arg.Iselimination,
	)
	var i GetJudgeInRoomRow
	err := row.Scan(
		&i.Assignmentid,
		&i.Tournamentid,
		&i.Judgeid,
		&i.Debateid,
		&i.Roundnumber,
		&i.Iselimination,
		&i.Isheadjudge,
		&i.Roomid,
		&i.Debateid_2,
	)
	return i, err
}

const swapJudges = `-- name: SwapJudges :exec
WITH judge_a_info AS (
    SELECT ja.debateid as a_debateid, ja.isheadjudge as a_isheadjudge
    FROM judgeassignments ja
    WHERE ja.judgeid = $1
      AND ja.tournamentid = $3
      AND ja.roundnumber = $4
      AND ja.iselimination = $5
),
     judge_b_info AS (
         SELECT ja.debateid as b_debateid, ja.isheadjudge as b_isheadjudge
         FROM judgeassignments ja
         WHERE ja.judgeid = $2
           AND ja.tournamentid = $3
           AND ja.roundnumber = $4
           AND ja.iselimination = $5
     )
UPDATE judgeassignments ja
SET debateid = CASE
                   WHEN ja.judgeid = $1 THEN (SELECT b_debateid FROM judge_b_info)
                   WHEN ja.judgeid = $2 THEN (SELECT a_debateid FROM judge_a_info)
    END,
    isheadjudge = CASE
                      WHEN ja.judgeid = $1 THEN (SELECT b_isheadjudge FROM judge_b_info)
                      WHEN ja.judgeid = $2 THEN (SELECT a_isheadjudge FROM judge_a_info)
        END
WHERE ja.judgeid IN ($1, $2)
  AND ja.tournamentid = $3
  AND ja.roundnumber = $4
  AND ja.iselimination = $5
`

type SwapJudgesParams struct {
	Judgeid       int32 `json:"judgeid"`
	Judgeid_2     int32 `json:"judgeid_2"`
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) SwapJudges(ctx context.Context, arg SwapJudgesParams) error {
	_, err := q.db.ExecContext(ctx, swapJudges,
		arg.Judgeid,
		arg.Judgeid_2,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Iselimination,
	)
	return err
}

const transferBallotOwnership = `-- name: TransferBallotOwnership :exec
UPDATE ballots b
SET judgeid = $2,
    last_updated_by = $2,
    last_updated_at = CURRENT_TIMESTAMP
WHERE b.judgeid = $1
  AND b.debateid IN (
    SELECT debateid
    FROM judgeassignments ja
    WHERE ja.judgeid = $2
)
`

type TransferBallotOwnershipParams struct {
	Judgeid   int32 `json:"judgeid"`
	Judgeid_2 int32 `json:"judgeid_2"`
}

func (q *Queries) TransferBallotOwnership(ctx context.Context, arg TransferBallotOwnershipParams) error {
	_, err := q.db.ExecContext(ctx, transferBallotOwnership, arg.Judgeid, arg.Judgeid_2)
	return err
}

const unassignJudge = `-- name: UnassignJudge :exec
DELETE FROM judgeassignments
WHERE judgeid = $1
  AND tournamentid = $2
  AND roundnumber = $3
  AND iselimination = $4
`

type UnassignJudgeParams struct {
	Judgeid       int32 `json:"judgeid"`
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) UnassignJudge(ctx context.Context, arg UnassignJudgeParams) error {
	_, err := q.db.ExecContext(ctx, unassignJudge,
		arg.Judgeid,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Iselimination,
	)
	return err
}

const updateJudgeAssignment = `-- name: UpdateJudgeAssignment :exec
UPDATE judgeassignments
SET debateid = (
    SELECT d.debateid
    FROM debates d
    WHERE d.roomid = $4
      AND d.tournamentid = $2
      AND d.roundnumber = $3
      AND d.iseliminationround = $6
    LIMIT 1
),
    isheadjudge = $5
WHERE judgeid = $1
  AND tournamentid = $2
  AND roundnumber = $3
`

type UpdateJudgeAssignmentParams struct {
	Judgeid            int32 `json:"judgeid"`
	Tournamentid       int32 `json:"tournamentid"`
	Roundnumber        int32 `json:"roundnumber"`
	Roomid             int32 `json:"roomid"`
	Isheadjudge        bool  `json:"isheadjudge"`
	Iseliminationround bool  `json:"iseliminationround"`
}

func (q *Queries) UpdateJudgeAssignment(ctx context.Context, arg UpdateJudgeAssignmentParams) error {
	_, err := q.db.ExecContext(ctx, updateJudgeAssignment,
		arg.Judgeid,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Roomid,
		arg.Isheadjudge,
		arg.Iseliminationround,
	)
	return err
}

const updateJudgeToHeadJudge = `-- name: UpdateJudgeToHeadJudge :exec
UPDATE judgeassignments ja
SET isheadjudge = true
WHERE ja.judgeid = $1
AND ja.tournamentid = $2
AND ja.roundnumber = $3
AND ja.debateid IN (
    SELECT d.debateid
    FROM debates d
    WHERE d.roomid = $4
    AND d.tournamentid = $2
)
AND ja.iselimination = $5
`

type UpdateJudgeToHeadJudgeParams struct {
	Judgeid       int32 `json:"judgeid"`
	Tournamentid  int32 `json:"tournamentid"`
	Roundnumber   int32 `json:"roundnumber"`
	Roomid        int32 `json:"roomid"`
	Iselimination bool  `json:"iselimination"`
}

func (q *Queries) UpdateJudgeToHeadJudge(ctx context.Context, arg UpdateJudgeToHeadJudgeParams) error {
	_, err := q.db.ExecContext(ctx, updateJudgeToHeadJudge,
		arg.Judgeid,
		arg.Tournamentid,
		arg.Roundnumber,
		arg.Roomid,
		arg.Iselimination,
	)
	return err
}
