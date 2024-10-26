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
UPDATE ballots
SET judgeid = $2,
    last_updated_by = $2,
    last_updated_at = CURRENT_TIMESTAMP
WHERE judgeid = $1
AND debateid = $3;

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
UPDATE judgeassignments ja
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
WHERE ja.judgeid = $1
AND ja.tournamentid = $2
AND ja.roundnumber = $3;