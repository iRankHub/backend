-- name: CheckHeadJudgeExists :one
SELECT EXISTS (
    SELECT 1
    FROM judgeassignments ja
    JOIN debates d ON ja.debateid = d.debateid
    WHERE ja.tournamentid = $1
    AND d.roomid = $2
    AND ja.roundnumber = $3
    AND ja.iselimination = $4
    AND ja.isheadjudge = true
);

-- name: GetJudgeAssignment :one
SELECT ja.*, d.roomid as roomid, d.debateid
FROM judgeassignments ja
JOIN debates d ON ja.debateid = d.debateid
WHERE ja.tournamentid = $1
AND ja.judgeid = $2
AND ja.roundnumber = $3
AND ja.iselimination = $4;

-- name: GetEligibleHeadJudge :one
SELECT ja.judgeid
FROM judgeassignments ja
JOIN debates d ON ja.debateid = d.debateid
WHERE ja.tournamentid = $1
AND ja.roundnumber = $2
AND d.roomid = $3
AND ja.iselimination = $4
AND ja.judgeid != $5
AND NOT ja.isheadjudge
LIMIT 1;

-- name: UpdateJudgeToHeadJudge :exec
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
AND ja.iselimination = $5;

-- name: TransferBallotOwnership :exec
UPDATE ballots b
SET judgeid = $2,
    last_updated_by = $2,
    last_updated_at = CURRENT_TIMESTAMP
WHERE b.judgeid = $1
  AND b.debateid IN (
    SELECT debateid
    FROM judgeassignments ja
    WHERE ja.judgeid = $2
);

-- name: DemoteCurrentHeadJudge :exec
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
AND ja.isheadjudge = true;

-- name: UpdateJudgeAssignment :exec
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
  AND roundnumber = $3;

-- name: GetJudgeInRoom :one
SELECT ja.*, d.roomid as roomid, d.debateid
FROM judgeassignments ja
         JOIN debates d ON ja.debateid = d.debateid
WHERE ja.tournamentid = $1
  AND d.roomid = $2
  AND ja.roundnumber = $3
  AND ja.iselimination = $4
LIMIT 1;

-- name: SwapJudges :exec
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
  AND ja.iselimination = $5;

-- name: UnassignJudge :exec
DELETE FROM judgeassignments
WHERE judgeid = $1
  AND tournamentid = $2
  AND roundnumber = $3
  AND iselimination = $4;

-- name: AssignJudgeToPosition :exec
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
FROM judge_b_info;

-- name: EnsureUnassignedRoomExists :one
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
SELECT ir.roomid FROM inserted_room ir;