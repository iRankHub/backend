-- name: GetJudgesByTournamentAndRound :many
SELECT u.UserID, u.Name, u.Email, ja.IsHeadJudge
FROM Users u
JOIN JudgeAssignments ja ON u.UserID = ja.JudgeID
WHERE ja.TournamentID = $1 AND ja.RoundNumber = $2 AND ja.IsElimination = $3;

-- name: GetJudgeByID :one
SELECT u.UserID, u.Name, u.Email
FROM Users u
WHERE u.UserID = $1;

-- name: AssignJudgeToDebate :exec
INSERT INTO JudgeAssignments (TournamentID, JudgeID, DebateID, RoundNumber, IsElimination, IsHeadJudge)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetAvailableJudges :many
SELECT DISTINCT u.UserID, u.Name, u.Email, v.iDebateVolunteerID
FROM Users u
         JOIN Volunteers v ON u.UserID = v.UserID
         JOIN TournamentInvitations ti ON ti.InviteeID = v.iDebateVolunteerID
WHERE ti.TournamentID = $1
  AND ti.Status = 'accepted'
  AND ti.InviteeRole = 'volunteer';

-- name: DeletePairingsForTournament :exec
DELETE FROM Debates
WHERE TournamentID = $1;

-- name: UpdatePairing :exec
UPDATE Debates
SET Team1ID = $2, Team2ID = $3
WHERE DebateID = $1;

-- name: GetBallotsByTournamentAndRound :many
SELECT b.BallotID, d.RoundNumber, d.IsEliminationRound, r.RoomName,
       u.Name AS HeadJudgeName, b.RecordingStatus, b.Verdict,
       t1.TeamID as Team1ID, t1.Name as Team1Name,
       t2.TeamID as Team2ID, t2.Name as Team2Name
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Users u ON b.JudgeID = u.UserID
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
WHERE d.TournamentID = $1 AND d.RoundNumber = $2 AND d.IsEliminationRound = $3;

-- name: GetBallotByID :one
SELECT b.BallotID, d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.RoomID, r.roomname AS RoomName, b.JudgeID, u.Name AS JudgeName,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       b.Team1TotalScore, b.Team2TotalScore, b.RecordingStatus, b.Verdict,
       b.Team1Feedback, b.Team2Feedback, b.last_updated_by, b.last_updated_at,
       b.head_judge_submitted
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Users u ON b.JudgeID = u.UserID
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
WHERE b.BallotID = $1;

-- name: GetBallotByJudgeID :one
SELECT b.BallotID, d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.RoomID, r.roomname AS RoomName, b.JudgeID, u.Name AS JudgeName,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       b.Team1TotalScore, b.Team2TotalScore, b.RecordingStatus, b.Verdict,
       b.Team1Feedback, b.Team2Feedback, b.last_updated_by, b.last_updated_at,
       b.head_judge_submitted
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Users u ON b.JudgeID = u.UserID
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
WHERE b.JudgeID = $1 AND d.TournamentID = $2
ORDER BY d.RoundNumber DESC
LIMIT 1;

-- name: UpdateBallot :exec
UPDATE Ballots
SET Team1TotalScore = $2, Team2TotalScore = $3, RecordingStatus = $4, Verdict = $5,
    Team1Feedback = $6, Team2Feedback = $7, last_updated_by = $8,
    last_updated_at = CURRENT_TIMESTAMP, head_judge_submitted = $9
WHERE BallotID = $1;

-- name: CreateInitialSpeakerScores :exec
WITH ballot_info AS (
    SELECT b.BallotID, d.Team1ID, d.Team2ID
    FROM Ballots b
    JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.DebateID = $1
    ORDER BY b.BallotID  -- Added explicit ordering
),
team_speakers AS (
    SELECT tm.StudentID as SpeakerID, t.TeamID,
           CASE
               WHEN t.TeamID = bi.Team1ID THEN 1
               WHEN t.TeamID = bi.Team2ID THEN 2
           END as TeamNumber
    FROM TeamMembers tm
    JOIN Teams t ON tm.TeamID = t.TeamID
    JOIN ballot_info bi ON t.TeamID IN (bi.Team1ID, bi.Team2ID)
)
INSERT INTO SpeakerScores (BallotID, SpeakerID, SpeakerRank, SpeakerPoints)
SELECT bi.BallotID, ts.SpeakerID,
       ROW_NUMBER() OVER (PARTITION BY bi.BallotID, ts.TeamNumber ORDER BY ts.SpeakerID) as SpeakerRank,
       0 as SpeakerPoints
FROM ballot_info bi
JOIN team_speakers ts ON (ts.TeamNumber = 1 AND bi.Team1ID = ts.TeamID)
                      OR (ts.TeamNumber = 2 AND bi.Team2ID = ts.TeamID);

-- name: GetSpeakerScoresByBallot :many
SELECT ss.ScoreID, ss.SpeakerID, s.FirstName, s.LastName,
       ss.SpeakerRank, ss.SpeakerPoints, ss.Feedback,
       t.TeamID, t.Name AS TeamName
FROM SpeakerScores ss
         JOIN Students s ON ss.SpeakerID = s.StudentID
         JOIN Ballots b ON ss.BallotID = b.BallotID
         JOIN Debates d ON b.DebateID = d.DebateID
         JOIN TeamMembers tm ON s.StudentID = tm.StudentID
         JOIN Teams t ON tm.TeamID = t.TeamID
WHERE ss.BallotID = $1
  AND (t.TeamID = d.Team1ID OR t.TeamID = d.Team2ID)
ORDER BY t.TeamID, ss.SpeakerRank;


-- name: UpdateSpeakerScore :exec
UPDATE SpeakerScores
SET SpeakerRank = $3, SpeakerPoints = $4, Feedback = $5
WHERE BallotID = $1 AND SpeakerID = $2;

-- name: IsHeadJudgeForBallot :one
SELECT COUNT(*) > 0 as is_head_judge
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
JOIN JudgeAssignments ja ON d.DebateID = ja.DebateID
WHERE b.BallotID = $1 AND ja.JudgeID = $2 AND ja.IsHeadJudge = true;

-- name: GetTeamWins :one
SELECT COUNT(*) as wins
FROM Debates d
JOIN Ballots b ON d.DebateID = b.DebateID
WHERE (d.Team1ID = $1 AND b.Team1TotalScore > b.Team2TotalScore)
   OR (d.Team2ID = $1 AND b.Team2TotalScore > b.Team1TotalScore)
   AND d.TournamentID = $2;

-- name: CreateDebate :one
INSERT INTO Debates (TournamentID, RoundID, RoundNumber, IsEliminationRound, Team1ID, Team2ID, RoomID, StartTime)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING DebateID;

-- name: CreateBallot :one
INSERT INTO ballots (
    debateid,
    judgeid,
    recordingstatus,
    verdict
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetTeamsByTournament :many
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
GROUP BY t.TeamID, t.Name, t.TournamentID, l.Name;


-- name: GetPreviousPairings :many
SELECT Team1ID, Team2ID
FROM Debates
WHERE TournamentID = $1 AND RoundNumber < $2;

-- name: CreatePairingHistory :exec
INSERT INTO PairingHistory (TournamentID, Team1ID, Team2ID, RoundNumber, IsElimination)
VALUES ($1, $2, $3, $4, $5);

-- name: CreateTeam :one
INSERT INTO Teams (Name, TournamentID)
VALUES ($1, $2)
RETURNING TeamID, Name, TournamentID;

-- name: AddTeamMember :one
INSERT INTO TeamMembers (TeamID, StudentID)
VALUES ($1, $2)
RETURNING TeamID, StudentID;

-- name: CheckExistingTeamMembership :one
SELECT COUNT(*) > 0 AS has_team
FROM TeamMembers tm
JOIN Teams t ON tm.TeamID = t.TeamID
WHERE t.TournamentID = $1 AND tm.StudentID = $2;

-- name: GetTeamByID :one
SELECT t.TeamID, t.Name, t.TournamentID,
       array_agg(tm.StudentID) as SpeakerIDs
FROM Teams t
LEFT JOIN TeamMembers tm ON t.TeamID = tm.TeamID
WHERE t.TeamID = $1
GROUP BY t.TeamID, t.Name, t.TournamentID;

-- name: UpdateTeam :exec
UPDATE Teams
SET Name = $2
WHERE TeamID = $1;

-- name: RemoveTeamMembers :exec
DELETE FROM TeamMembers
WHERE TeamID = $1;

-- name: GetTeamMembers :many
SELECT tm.TeamID, tm.StudentID, s.FirstName, s.LastName
FROM TeamMembers tm
JOIN Students s ON tm.StudentID = s.StudentID
WHERE tm.TeamID = $1;

-- name: DeleteTeam :exec
WITH debate_check AS (
    SELECT 1
    FROM Debates
    WHERE Team1ID = $1 OR Team2ID = $1
    LIMIT 1
)
DELETE FROM Teams
WHERE TeamID = $1 AND NOT EXISTS (SELECT 1 FROM debate_check);

-- name: DeleteTeamMembers :exec
DELETE FROM TeamMembers
WHERE TeamID = $1;


-- name: GetRoundByTournamentAndNumber :one
SELECT * FROM Rounds
WHERE TournamentID = $1 AND RoundNumber = $2 AND IsEliminationRound = $3
LIMIT 1;

-- name: DeleteJudgeAssignmentsForTournament :exec
DELETE FROM JudgeAssignments
WHERE TournamentID = $1;

-- name: DeleteDebatesForTournament :exec
DELETE FROM Debates
WHERE TournamentID = $1;

-- name: DeleteRoundsForTournament :exec
DELETE FROM Rounds
WHERE TournamentID = $1;

-- name: DeletePairingHistoryForTournament :exec
DELETE FROM PairingHistory
WHERE TournamentID = $1;


-- name: GetRoomByID :one
SELECT RoomID, RoomName, TournamentID, Location, Capacity
FROM Rooms
WHERE RoomID = $1;

-- name: CreateRoom :one
INSERT INTO Rooms (RoomName, Location, Capacity, TournamentID)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateRoom :one
UPDATE Rooms
SET RoomName = $2
WHERE RoomID = $1
RETURNING *;


-- name: AssignRoomToDebate :exec
UPDATE Debates
SET RoomID = $2
WHERE DebateID = $1;

-- name: DeleteRoomsForTournament :exec
DELETE FROM Rooms
WHERE TournamentID = $1;

-- name: GetPairingsByTournamentAndRound :many
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
         l1.Name, l2.Name, t1_points.TotalPoints, t2_points.TotalPoints;

-- name: GetPairings :many
SELECT
    d.debateid,
    d.roundnumber,
    d.iseliminationround,
    d.roomid,
    r.roomname,
    t1.teamid AS team1id,
    t1.name AS team1name,
    t2.teamid AS team2id,
    t2.name AS team2name,
    j.name AS headjudgename
FROM
    Debates d
    JOIN Rooms r ON d.roomid = r.roomid
    JOIN Teams t1 ON d.team1id = t1.teamid
    JOIN Teams t2 ON d.team2id = t2.teamid
    LEFT JOIN JudgeAssignments ja ON d.debateid = ja.debateid AND ja.isheadjudge = true
    LEFT JOIN Users j ON ja.judgeid = j.userid
WHERE
    d.tournamentid = $1
    AND d.roundnumber = $2
    AND d.iseliminationround = $3
ORDER BY
    d.debateid;

-- name: GetPairingByID :one
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
         l1.Name, l2.Name, t1_points.TotalPoints, t2_points.TotalPoints;

-- name: GetSinglePairing :one
SELECT
    d.debateid,
    d.roundnumber,
    d.iseliminationround,
    d.roomid,
    r.roomname,
    t1.teamid AS team1id,
    t1.name AS team1name,
    t2.teamid AS team2id,
    t2.name AS team2name,
    j.name AS headjudgename
FROM
    Debates d
    JOIN Rooms r ON d.roomid = r.roomid
    JOIN Teams t1 ON d.team1id = t1.teamid
    JOIN Teams t2 ON d.team2id = t2.teamid
    LEFT JOIN JudgeAssignments ja ON d.debateid = ja.debateid AND ja.isheadjudge = true
    LEFT JOIN Users j ON ja.judgeid = j.userid
WHERE
    d.debateid = $1;


-- name: GetRoomsByTournament :many
SELECT * FROM Rooms
WHERE TournamentID = $1;

-- name: GetDebateByRoomAndRound :one
SELECT *
FROM Debates
WHERE TournamentID = $1 AND RoomID = $2 AND RoundNumber = $3 AND IsEliminationRound = $4
LIMIT 1;

-- name: GetDebatesByRoomAndTournament :many
SELECT *
FROM Debates
WHERE TournamentID = $1 AND RoomID = $2 AND IsEliminationRound = $3;

-- name: GetJudgesForDebate :many
SELECT ja.JudgeID, u.Name
FROM JudgeAssignments ja
JOIN Users u ON ja.JudgeID = u.UserID
WHERE ja.DebateID = $1;

-- name: GetJudgesForTournament :many
SELECT
    u.UserID as JudgeID,
    u.Name,
    v.iDebateVolunteerID,
    COALESCE(r.RoomName, 'Unassigned') AS RoomName,
    COALESCE(ja.IsHeadJudge, false) AS IsHeadJudge,
    CASE WHEN ja.JudgeID IS NULL THEN false ELSE true END AS IsAssigned
FROM
    Users u
        JOIN
    Volunteers v ON u.UserID = v.UserID
        JOIN
    TournamentInvitations ti ON ti.InviteeID = v.iDebateVolunteerID
        LEFT JOIN
    JudgeAssignments ja ON u.UserID = ja.JudgeID AND ja.TournamentID = $1
        LEFT JOIN
    Debates d ON ja.DebateID = d.DebateID
        LEFT JOIN
    Rooms r ON d.RoomID = r.RoomID
WHERE
    ti.TournamentID = $1
  AND ti.Status = 'accepted'
  AND ti.InviteeRole = 'volunteer'
ORDER BY
    IsAssigned DESC, u.Name;

-- name: CountJudgeDebates :one
SELECT
    COUNT(DISTINCT d.DebateID) as DebateCount
FROM
    JudgeAssignments ja
JOIN
    Debates d ON ja.DebateID = d.DebateID
WHERE
    ja.JudgeID = $1 AND
    d.TournamentID = $2 AND
    d.IsEliminationRound = $3;

-- name: GetJudgeDetails :one
SELECT
    u.UserID as JudgeID,
    u.Name,
    v.iDebateVolunteerID
FROM
    Users u
JOIN
    Volunteers v ON u.UserID = v.UserID
WHERE
    u.UserID = $1;

-- name: GetJudgeRooms :many
SELECT
    d.RoundNumber,
    d.RoomID,
    r.RoomName,
    ja.IsHeadJudge
FROM
    JudgeAssignments ja
        JOIN
    Debates d ON ja.DebateID = d.DebateID
        JOIN
    Rooms r ON d.RoomID = r.RoomID
WHERE
    ja.JudgeID = $1 AND
    d.TournamentID = $2 AND
    d.IsEliminationRound = $3
ORDER BY
    d.RoundNumber;

-- name: UpdateJudgeRoom :exec
WITH new_debate AS (
    SELECT d.DebateID
    FROM Debates d
    WHERE d.TournamentID = $2
    AND d.RoundNumber = $3
    AND d.RoomID = $4
    AND d.IsEliminationRound = $5
),
old_assignment AS (
    SELECT ja.DebateID, ja.IsHeadJudge
    FROM JudgeAssignments ja
    JOIN Debates d ON ja.DebateID = d.DebateID
    WHERE ja.JudgeID = $1
    AND d.TournamentID = $2
    AND d.RoundNumber = $3
    AND d.IsEliminationRound = $5
)
UPDATE JudgeAssignments ja
SET DebateID = nd.DebateID
FROM new_debate nd, old_assignment oa
WHERE ja.JudgeID = $1
AND ja.TournamentID = $2
AND ja.DebateID = oa.DebateID
AND nd.DebateID IS NOT NULL;

-- name: GetEliminationRoundTeams :many
SELECT
    CASE
        WHEN b.Verdict = t1.Name THEN d.Team1ID
        WHEN b.Verdict = t2.Name THEN d.Team2ID
    END AS TeamID,
    CASE
        WHEN b.Verdict = t1.Name THEN t1.Name
        WHEN b.Verdict = t2.Name THEN t2.Name
    END AS TeamName,
    d.TournamentID,
    CASE
        WHEN b.Verdict = t1.Name THEN b.Team1TotalScore
        WHEN b.Verdict = t2.Name THEN b.Team2TotalScore
    END AS TotalScore
FROM
    Debates d
JOIN
    Ballots b ON d.DebateID = b.DebateID
JOIN
    Teams t1 ON d.Team1ID = t1.TeamID
JOIN
    Teams t2 ON d.Team2ID = t2.TeamID
WHERE
    d.TournamentID = $1
    AND d.RoundNumber = $2
    AND d.IsEliminationRound = true
    AND b.RecordingStatus = 'Recorded'
    AND (b.Verdict = t1.Name OR b.Verdict = t2.Name)
ORDER BY
    CASE
        WHEN b.Verdict = t1.Name THEN b.Team1TotalScore
        WHEN b.Verdict = t2.Name THEN b.Team2TotalScore
    END DESC
LIMIT $3;

-- name: GetTopPerformingTeams :many
WITH ballot_check AS (
    SELECT COUNT(*) = 0 AS all_recorded
    FROM Ballots b
             JOIN Debates d ON b.DebateID = d.DebateID
    WHERE d.TournamentID = $1 AND d.IsEliminationRound = false AND b.RecordingStatus != 'Recorded'
),
     team_performance AS (
         SELECT
             t.TeamID,
             t.Name,
             t.TournamentID,
             COALESCE(t.TotalWins, 0) as Wins,
             COALESCE(t.TotalSpeakerPoints, 0) as TotalSpeakerPoints,
             COALESCE(t.AverageRank, 99) as AverageRank  -- Use 99 as default for ranking purposes
         FROM
             Teams t
         WHERE
             t.TournamentID = $1
           AND EXISTS (
             -- Only include teams that have participated in debates
             SELECT 1 FROM Debates d
             WHERE (d.Team1ID = t.TeamID OR d.Team2ID = t.TeamID)
               AND d.IsEliminationRound = false
         )
     )
SELECT
    tp.TeamID, tp.Name, tp.TournamentID, tp.Wins, tp.TotalSpeakerPoints, tp.AverageRank
FROM
    team_performance tp, ballot_check
WHERE
    ballot_check.all_recorded = true
ORDER BY
    tp.Wins DESC, tp.TotalSpeakerPoints DESC, tp.AverageRank ASC
LIMIT $2;

-- name: GetDebateByBallotID :one
SELECT d.DebateID, d.Team1ID, d.Team2ID, d.IsEliminationRound, d.TournamentID
FROM Debates d
JOIN Ballots b ON d.DebateID = b.DebateID
WHERE b.BallotID = $1;

-- name: GetDebateByID :one
SELECT * FROM Debates
WHERE DebateID = $1;

-- name: GetBallotsByDebateID :many
SELECT * FROM Ballots
WHERE DebateID = $1;

-- name: DeleteSpeakerScoresByBallot :exec
DELETE FROM SpeakerScores
WHERE BallotID = $1;

-- name: ResetBallotAfterTeamChange :exec
UPDATE Ballots
SET Team1TotalScore = $2,
    Team2TotalScore = $3,
    RecordingStatus = $4,
    Verdict = $5,
    Team1Feedback = $6,
    Team2Feedback = $7,
    head_judge_submitted = $8,
    last_updated_at = CURRENT_TIMESTAMP
WHERE BallotID = $1;

-- name: UpdateTeamScore :exec
UPDATE TeamScores
SET TotalScore = $3, IsElimination = $4
WHERE TeamID = $1 AND DebateID = $2;

-- name: InsertTeamScore :exec
INSERT INTO TeamScores (TeamID, DebateID, TotalScore, IsElimination)
SELECT $1, $2, $3, $4
WHERE NOT EXISTS (
    SELECT 1 FROM TeamScores
    WHERE TeamID = $1 AND DebateID = $2
);

-- name: GetTeamAverageRank :one
WITH speaker_ranks AS (
    SELECT ss.SpeakerRank
    FROM SpeakerScores ss
    JOIN TeamMembers tm ON ss.SpeakerID = tm.StudentID
    WHERE tm.TeamID = $1 AND ss.BallotID = $2
)
SELECT
    AVG(SpeakerRank)::FLOAT as AvgRank,
    COUNT(*) as SpeakerCount,
    array_agg(SpeakerRank) as AllRanks
FROM speaker_ranks;

-- name: UpdateTeamScoreRank :exec
UPDATE TeamScores
SET Rank = $3
WHERE TeamID = $1 AND DebateID = $2;

-- name: UpdateTeamStats :exec
WITH team_stats AS (
    SELECT
        t.TeamID,
        t.TournamentID,
        COUNT(CASE WHEN b.Verdict = t.Name THEN 1 ELSE NULL END) AS TotalWins,
        AVG(ts.Rank) AS AvgRank,
        SUM(ts.TotalScore::NUMERIC) AS TotalSpeakerPoints
    FROM
        Teams t
    JOIN
        Debates d ON (t.TeamID = d.Team1ID OR t.TeamID = d.Team2ID) AND t.TournamentID = d.TournamentID
    JOIN
        Ballots b ON d.DebateID = b.DebateID
    JOIN
        TeamScores ts ON t.TeamID = ts.TeamID AND d.DebateID = ts.DebateID
    WHERE
        t.TeamID = $1 AND t.TournamentID = $2
    GROUP BY
        t.TeamID, t.TournamentID
)
UPDATE Teams
SET
    TotalWins = team_stats.TotalWins,
    AverageRank = team_stats.AvgRank,
    TotalSpeakerPoints = team_stats.TotalSpeakerPoints
FROM
    team_stats
WHERE
    Teams.TeamID = team_stats.TeamID AND Teams.TournamentID = team_stats.TournamentID;

-- name: GetTournamentStudentRanking :many
WITH RankedStudents AS (
    SELECT
        s.StudentID,
        s.FirstName || ' ' || s.LastName AS StudentName,
        sch.SchoolName,
        t.TotalWins AS TotalWins,
        CAST(SUM(ss.SpeakerPoints) AS DECIMAL(10,2)) AS TotalPoints,
        AVG(ss.SpeakerRank) AS AverageRank,
        RANK() OVER (ORDER BY SUM(ss.SpeakerPoints) DESC,
            t.TotalWins DESC,
            AVG(ss.SpeakerRank) ASC) as place
    FROM Students s
             JOIN TeamMembers tm ON s.StudentID = tm.StudentID
             JOIN Teams t ON tm.TeamID = t.TeamID
             JOIN Debates d ON (t.TeamID = d.Team1ID OR t.TeamID = d.Team2ID)
             JOIN Ballots b ON d.DebateID = b.DebateID
             JOIN SpeakerScores ss ON s.StudentID = ss.SpeakerID AND b.BallotID = ss.BallotID
             JOIN Schools sch ON s.SchoolID = sch.SchoolID
    WHERE d.TournamentID = $1 AND d.IsEliminationRound = false
    GROUP BY s.StudentID, StudentName, sch.SchoolName, t.TotalWins, t.TeamID
),
     TopThree AS (
         SELECT *
         FROM RankedStudents
         WHERE place <= 3
     ),
     SearchResults AS (
         SELECT *
         FROM RankedStudents
         WHERE $4::text IS NULL
            OR StudentName ILIKE '%' || $4 || '%'
            OR SchoolName ILIKE '%' || $4 || '%'
     )
SELECT DISTINCT ON (place) *
FROM (
         SELECT * FROM TopThree
         UNION
         SELECT * FROM SearchResults WHERE place > 3
     ) combined
ORDER BY place
LIMIT $2 OFFSET $3;

-- name: GetOverallStudentRanking :many
WITH student_ranking AS (
    SELECT
        s.StudentID,
        s.FirstName || ' ' || s.LastName AS StudentName,
        CAST(SUM(ss.SpeakerPoints) AS DECIMAL(10,2)) AS TotalPoints,
        AVG(ss.SpeakerRank) AS AverageRank,
        COUNT(DISTINCT d.TournamentID) AS TournamentsParticipated,
        RANK() OVER (ORDER BY SUM(ss.SpeakerPoints) DESC, AVG(ss.SpeakerRank) ASC) AS CurrentRank,
        COUNT(*) OVER () AS TotalStudents,
        MAX(t.StartDate) AS LastTournamentDate
    FROM
        Students s
    JOIN TeamMembers tm ON s.StudentID = tm.StudentID
    JOIN Teams te ON tm.TeamID = te.TeamID
    JOIN Debates d ON (te.TeamID = d.Team1ID OR te.TeamID = d.Team2ID)
    JOIN Ballots b ON d.DebateID = b.DebateID
    JOIN SpeakerScores ss ON s.StudentID = ss.SpeakerID AND b.BallotID = ss.BallotID
    JOIN Tournaments t ON d.TournamentID = t.TournamentID
    GROUP BY
        s.StudentID, s.FirstName, s.LastName
)
SELECT *
FROM student_ranking
ORDER BY CurrentRank;

-- name: GetStudentOverallPerformance :many
WITH tournament_performance AS (
    SELECT
        d.TournamentID,
        t.StartDate,
        s.StudentID,
        CAST(SUM(ss.SpeakerPoints) AS NUMERIC(10,2)) AS StudentTotalPoints,
        CAST(AVG(ss.SpeakerPoints) AS NUMERIC(10,2)) AS StudentAveragePoints,
        CAST(AVG(SUM(ss.SpeakerPoints)) OVER (PARTITION BY d.TournamentID) AS NUMERIC(10,2)) AS OverallAverageTotalPoints,
        CAST(AVG(AVG(ss.SpeakerPoints)) OVER (PARTITION BY d.TournamentID) AS NUMERIC(10,2)) AS OverallAveragePoints,
        RANK() OVER (PARTITION BY d.TournamentID ORDER BY SUM(ss.SpeakerPoints) DESC) AS TournamentRank
    FROM
        Students s
    JOIN TeamMembers tm ON s.StudentID = tm.StudentID
    JOIN Teams te ON tm.TeamID = te.TeamID
    JOIN Debates d ON (te.TeamID = d.Team1ID OR te.TeamID = d.Team2ID)
    JOIN Ballots b ON d.DebateID = b.DebateID
    JOIN SpeakerScores ss ON s.StudentID = ss.SpeakerID AND b.BallotID = ss.BallotID
    JOIN Tournaments t ON d.TournamentID = t.TournamentID
    WHERE
        s.StudentID = $1 AND t.StartDate BETWEEN $2 AND $3
    GROUP BY
        d.TournamentID, t.StartDate, s.StudentID
)
SELECT
    StartDate,
    StudentTotalPoints,
    StudentAveragePoints,
    OverallAverageTotalPoints,
    OverallAveragePoints,
    TournamentRank
FROM
    tournament_performance
ORDER BY
    StartDate;

-- name: GetTournamentTeamsRanking :many
WITH TeamData AS (
    SELECT
        t.TeamID,
        t.Name AS TeamName,
        t.TotalWins AS Wins,
        t.TotalSpeakerPoints AS TotalPoints,
        t.AverageRank,
        t.TournamentID
    FROM
        Teams t
    WHERE
        t.TournamentID = $1
),
     TeamSchools AS (
         SELECT
             td.TeamID,
             td.TeamName,
             td.Wins,
             td.TotalPoints,
             td.AverageRank,
             ARRAY_AGG(DISTINCT s.SchoolName) FILTER (WHERE s.SchoolName IS NOT NULL) AS SchoolNames
         FROM
             TeamData td
                 LEFT JOIN TeamMembers tm ON td.TeamID = tm.TeamID
                 LEFT JOIN Students stu ON tm.StudentID = stu.StudentID
                 LEFT JOIN Schools s ON stu.SchoolID = s.SchoolID
         GROUP BY
             td.TeamID, td.TeamName, td.Wins, td.TotalPoints, td.AverageRank
     ),
     RankedTeams AS (
         SELECT
             TeamID,
             TeamName,
             SchoolNames,
             Wins,
             TotalPoints,
             AverageRank,
             RANK() OVER (ORDER BY Wins DESC, TotalPoints DESC, AverageRank ASC) as place
         FROM
             TeamSchools
     ),
     TopThree AS (
         SELECT * FROM RankedTeams WHERE place <= 3
     ),
     FilteredTeams AS (
         SELECT * FROM RankedTeams
         WHERE $4::text IS NULL
            OR TeamName ILIKE '%' || $4 || '%'
            OR EXISTS (
             SELECT 1
             FROM unnest(SchoolNames) AS school
             WHERE school ILIKE '%' || $4 || '%'
         )
     )
SELECT * FROM (
                  SELECT * FROM TopThree
                  UNION ALL
                  SELECT * FROM FilteredTeams WHERE place > 3
              ) combined
ORDER BY place
LIMIT $2 OFFSET $3;

-- name: GetTournamentTeamsRankingCount :one
SELECT COUNT(DISTINCT t.TeamID)
FROM Teams t
WHERE t.TournamentID = $1;

-- name: GetTournamentSchoolRanking :many
WITH SchoolTeams AS (
    -- First, get distinct teams by school
    SELECT DISTINCT
        s.SchoolID,
        s.SchoolName,
        t.TeamID,
        t.TotalWins,
        t.TotalSpeakerPoints,
        t.AverageRank
    FROM Schools s
             JOIN Students stu ON s.SchoolID = stu.SchoolID
             JOIN TeamMembers tm ON stu.StudentID = tm.StudentID
             JOIN Teams t ON tm.TeamID = t.TeamID
             JOIN Tournaments tour ON t.TournamentID = tour.TournamentID
             LEFT JOIN Leagues l ON tour.LeagueID = l.LeagueID
    WHERE t.TournamentID = $1
      AND (l.Name != 'DAC' OR l.Name IS NULL)
),
     SchoolData AS (
         -- Aggregate school data using the distinct teams
         SELECT
             SchoolID,
             SchoolName,
             COUNT(TeamID) AS TeamCount,
             SUM(TotalWins) AS TotalWins,
             CAST(SUM(TotalSpeakerPoints) AS numeric(10,2)) AS TotalPoints,
             CAST(AVG(AverageRank) AS numeric(5,2)) AS AverageRank
         FROM SchoolTeams
         GROUP BY SchoolID, SchoolName
         HAVING COUNT(TeamID) > 0
     ),
     RankedSchools AS (
         SELECT
             SchoolName,
             TeamCount,
             TotalWins,
             CAST(AverageRank AS text) AS AverageRank,
             CAST(TotalPoints AS text) AS TotalPoints,
             RANK() OVER (
                 ORDER BY TotalWins DESC,
                     TotalPoints DESC,
                     AverageRank ASC
                 ) as place
         FROM SchoolData
     ),
     TopThree AS (
         SELECT *
         FROM RankedSchools
         WHERE place <= 3
     ),
     SearchResults AS (
         SELECT *
         FROM RankedSchools
         WHERE $4::text IS NULL
            OR SchoolName ILIKE '%' || $4 || '%'
     )
SELECT DISTINCT ON (place) *
FROM (
         SELECT * FROM TopThree
         UNION
         SELECT * FROM SearchResults WHERE place > 3
     ) combined
ORDER BY place
LIMIT $2 OFFSET $3;

-- name: GetTournamentSchoolRankingCount :one
WITH SchoolTeams AS (
    -- Get distinct teams per school to avoid duplicates
    SELECT DISTINCT
        s.SchoolID
    FROM Schools s
             JOIN Students stu ON s.SchoolID = stu.SchoolID
             JOIN TeamMembers tm ON stu.StudentID = tm.StudentID
             JOIN Teams t ON tm.TeamID = t.TeamID
             JOIN Tournaments tour ON t.TournamentID = tour.TournamentID
             LEFT JOIN Leagues l ON tour.LeagueID = l.LeagueID
    WHERE t.TournamentID = $1
      AND (l.Name != 'DAC' OR l.Name IS NULL)
)
SELECT COUNT(DISTINCT SchoolID)
FROM SchoolTeams;

-- name: GetOverallSchoolRanking :many
WITH school_ranking AS (
    SELECT
        s.SchoolID,
        s.SchoolName,
        CAST(COALESCE(SUM(ts.TotalScore), 0) AS DECIMAL(10,2)) AS TotalPoints,
        COALESCE(AVG(ts.Rank), 0) AS AverageRank,
        COUNT(DISTINCT tour.TournamentID) as TournamentsParticipated,
        RANK() OVER (ORDER BY COALESCE(SUM(ts.TotalScore), 0) DESC, COALESCE(AVG(ts.Rank), 0) ASC) as CurrentRank,
        COUNT(*) OVER () as TotalSchools,
        MAX(tour.StartDate) as LastTournamentDate
    FROM
        Schools s
        LEFT JOIN Students stu ON s.SchoolID = stu.SchoolID
        LEFT JOIN TeamMembers tm ON stu.StudentID = tm.StudentID
        LEFT JOIN Teams te ON tm.TeamID = te.TeamID
        LEFT JOIN Debates d ON (te.TeamID = d.Team1ID OR te.TeamID = d.Team2ID)
        LEFT JOIN TeamScores ts ON te.TeamID = ts.TeamID AND d.DebateID = ts.DebateID
        LEFT JOIN Tournaments tour ON d.TournamentID = tour.TournamentID
        LEFT JOIN Leagues l ON tour.LeagueID = l.LeagueID
    WHERE
        s.SchoolID = $1 OR s.SchoolID IN (
            SELECT s2.SchoolID
            FROM Schools s2
            LEFT JOIN Students stu2 ON s2.SchoolID = stu2.SchoolID
            LEFT JOIN TeamMembers tm2 ON stu2.StudentID = tm2.StudentID
            LEFT JOIN Teams te2 ON tm2.TeamID = te2.TeamID
            LEFT JOIN Debates d2 ON (te2.TeamID = d2.Team1ID OR te2.TeamID = d2.Team2ID)
            LEFT JOIN TeamScores ts2 ON te2.TeamID = ts2.TeamID AND d2.DebateID = ts2.DebateID
            GROUP BY s2.SchoolID
            ORDER BY COALESCE(SUM(ts2.TotalScore), 0) DESC
            LIMIT 3
        )
    GROUP BY
        s.SchoolID, s.SchoolName
)
SELECT *
FROM school_ranking
ORDER BY CurrentRank;

-- name: GetSchoolOverallPerformance :many
WITH tournament_performance AS (
  SELECT
    d.TournamentID,
    t.StartDate,
    s.SchoolID,
    CAST(SUM(ts.TotalScore) AS DECIMAL(10,2)) AS SchoolTotalPoints,
    CAST(AVG(ts.TotalScore) AS DECIMAL(10,2)) AS SchoolAveragePoints,
    CAST(AVG(SUM(ts.TotalScore)) OVER (PARTITION BY d.TournamentID) AS DECIMAL(10,2)) AS OverallAverageTotalPoints,
    CAST(AVG(AVG(ts.TotalScore)) OVER (PARTITION BY d.TournamentID) AS DECIMAL(10,2)) AS OverallAveragePoints,
    RANK() OVER (PARTITION BY d.TournamentID ORDER BY SUM(ts.TotalScore) DESC) AS TournamentRank
  FROM
    Schools s
    JOIN Students stu ON s.SchoolID = stu.SchoolID
    JOIN TeamMembers tm ON stu.StudentID = tm.StudentID
    JOIN Teams te ON tm.TeamID = te.TeamID
    JOIN Debates d ON (te.TeamID = d.Team1ID OR te.TeamID = d.Team2ID)
    JOIN TeamScores ts ON te.TeamID = ts.TeamID AND d.DebateID = ts.DebateID
    JOIN Tournaments t ON d.TournamentID = t.TournamentID
    LEFT JOIN Leagues l ON t.LeagueID = l.LeagueID
  WHERE
    s.SchoolID = $1 AND t.StartDate BETWEEN $2 AND $3 AND l.Name != 'DAC'
  GROUP BY
    d.TournamentID, t.StartDate, s.SchoolID
)
SELECT
  StartDate,
  SchoolTotalPoints,
  SchoolAveragePoints,
  OverallAverageTotalPoints,
  OverallAveragePoints,
  TournamentRank
FROM
  tournament_performance
ORDER BY
  StartDate;

-- name: GetStudentTournamentStats :one
WITH student_stats AS (
    SELECT
        COUNT(DISTINCT t.TournamentID) AS attended_tournaments,
        (SELECT COUNT(*) FROM Tournaments WHERE StartDate >= CURRENT_DATE - INTERVAL '365 days') AS total_tournaments_last_year
    FROM
        Students s
    JOIN TeamMembers tm ON s.StudentID = tm.StudentID
    JOIN Teams te ON tm.TeamID = te.TeamID
    JOIN Tournaments t ON te.TournamentID = t.TournamentID
    WHERE
        s.StudentID = $1 AND t.StartDate >= CURRENT_DATE - INTERVAL '365 days'
),
current_stats AS (
    SELECT
        (SELECT COUNT(*) FROM Tournaments WHERE deleted_at IS NULL) AS total_tournaments,
        (SELECT COUNT(*) FROM Tournaments WHERE StartDate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days' AND deleted_at IS NULL) AS upcoming_tournaments
),
yesterday_stats AS (
    SELECT yesterday_total_count, yesterday_upcoming_count
    FROM Tournaments
    WHERE TournamentID = (SELECT MIN(TournamentID) FROM Tournaments)
)
SELECT
    cs.total_tournaments,
    ys.yesterday_total_count,
    cs.upcoming_tournaments,
    ys.yesterday_upcoming_count,
    ss.attended_tournaments,
    ss.total_tournaments_last_year
FROM
    current_stats cs, yesterday_stats ys, student_stats ss;

-- name: GetVolunteerTournamentStats :one
WITH volunteer_stats AS (
    SELECT
        COUNT(DISTINCT CASE WHEN d.DebateID IS NOT NULL THEN d.DebateID END) AS total_rounds_judged,
        COUNT(DISTINCT CASE WHEN t.TournamentID IS NOT NULL THEN t.TournamentID END) AS attended_tournaments,
        (
            SELECT COUNT(*)
            FROM TournamentInvitations ti
            JOIN Tournaments t ON ti.TournamentID = t.TournamentID
            JOIN Volunteers v2 ON ti.InviteeID = v2.iDebateVolunteerID
            WHERE v2.UserID = $1
            AND ti.Status = 'accepted'
            AND t.StartDate > CURRENT_DATE
        ) AS upcoming_tournaments
    FROM
        Volunteers v
        LEFT JOIN JudgeAssignments ja ON v.UserID = ja.JudgeID
        LEFT JOIN Debates d ON ja.DebateID = d.DebateID
        LEFT JOIN Tournaments t ON d.TournamentID = t.TournamentID
    WHERE
        v.UserID = $1
    GROUP BY v.UserID
)
SELECT
    vs.*,
    COALESCE(v.yesterday_rounds_judged, 0) as yesterday_rounds_judged,
    COALESCE(v.yesterday_tournaments_attended, 0) as yesterday_tournaments_attended,
    COALESCE(v.yesterday_upcoming_tournaments, 0) as yesterday_upcoming_tournaments
FROM volunteer_stats vs
JOIN Volunteers v ON v.UserID = $1;

-- name: GetStudentFeedback :many
WITH debate_judges AS (
    SELECT
        ja.DebateID,
        json_agg(
            json_build_object(
                'judge_id', ja.JudgeID,
                'judge_name', u.Name,
                'is_head_judge', ja.IsHeadJudge
            )
        ) as judges
    FROM JudgeAssignments ja
    JOIN Users u ON ja.JudgeID = u.UserID
    GROUP BY ja.DebateID
)
SELECT
    d.RoundNumber,
    d.IsEliminationRound,
    d.DebateID,
    b.BallotID,
    ss.SpeakerPoints,
    ss.Feedback,
    ss.IsRead,
    u.Name as HeadJudgeName,
    r.RoomName,
    t1.Name as OpponentTeamName,
    t2.Name as StudentTeamName,
    dj.judges::text as JudgesInfo
FROM SpeakerScores ss
JOIN Ballots b ON ss.BallotID = b.BallotID
JOIN Debates d ON b.DebateID = d.DebateID
JOIN JudgeAssignments ja ON d.DebateID = ja.DebateID AND ja.IsHeadJudge = true
JOIN Users u ON ja.JudgeID = u.UserID
JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Teams t1 ON (d.Team1ID = t1.TeamID)
JOIN Teams t2 ON (d.Team2ID = t2.TeamID AND t2.TeamID IN (
    SELECT TeamID FROM TeamMembers WHERE StudentID = $1
))
LEFT JOIN debate_judges dj ON d.DebateID = dj.DebateID
WHERE ss.SpeakerID = $1
  AND d.TournamentID = $2
ORDER BY d.RoundNumber
LIMIT $3 OFFSET $4;

-- name: GetStudentFeedbackCount :one
SELECT COUNT(*)
FROM SpeakerScores ss
JOIN Ballots b ON ss.BallotID = b.BallotID
JOIN Debates d ON b.DebateID = d.DebateID
WHERE ss.SpeakerID = $1 AND d.TournamentID = $2;

-- name: CreateJudgeFeedback :one
INSERT INTO JudgeFeedback (
    JudgeID, StudentID, DebateID,
    ClarityRating, ConstructivenessRating, TimelinessRating,
    FairnessRating, EngagementRating, TextFeedback
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetJudgeFeedbackList :many
SELECT
    f.*,
    d.RoundNumber,
    d.IsEliminationRound,
    t.StartDate as TournamentDate,
    b.BallotID as BallotID
FROM JudgeFeedback f
JOIN Debates d ON f.DebateID = d.DebateID
JOIN Tournaments t ON d.TournamentID = t.TournamentID
JOIN Ballots b ON d.DebateID = b.DebateID
WHERE f.JudgeID = $1
ORDER BY f.CreatedAt DESC
LIMIT $2 OFFSET $3;

-- name: GetJudgeFeedbackCount :one
SELECT COUNT(*)
FROM JudgeFeedback
WHERE JudgeID = $1;

-- name: GetOverallVolunteerRanking :many
WITH volunteer_ranking AS (
    SELECT
        v.VolunteerID,
        v.FirstName || ' ' || v.LastName AS VolunteerName,
        AVG(jf.AverageRating) as AverageRating,
        COUNT(DISTINCT d.TournamentID) as TournamentsJudged,
        RANK() OVER (ORDER BY AVG(jf.AverageRating) DESC) as CurrentRank,
        COUNT(*) OVER () as TotalVolunteers,
        MAX(t.StartDate) as LastTournamentDate
    FROM
        Volunteers v
        LEFT JOIN JudgeAssignments ja ON v.UserID = ja.JudgeID
        LEFT JOIN Debates d ON ja.DebateID = d.DebateID
        LEFT JOIN JudgeFeedback jf ON v.UserID = jf.JudgeID
        LEFT JOIN Tournaments t ON d.TournamentID = t.TournamentID
    GROUP BY
        v.VolunteerID, v.FirstName, v.LastName
)
SELECT *
FROM volunteer_ranking
ORDER BY CurrentRank;

-- name: GetVolunteerPerformance :many
WITH tournament_performance AS (
    SELECT
        t.StartDate,
        v.UserID,
        AVG(jf.AverageRating) as VolunteerAverageRating,
        AVG(AVG(jf.AverageRating)) OVER (PARTITION BY d.TournamentID) as OverallAverageRating,
        RANK() OVER (PARTITION BY d.TournamentID ORDER BY AVG(jf.AverageRating) DESC) as TournamentRank
    FROM
        Volunteers v
        JOIN JudgeAssignments ja ON v.UserID = ja.JudgeID
        JOIN Debates d ON ja.DebateID = d.DebateID
        JOIN JudgeFeedback jf ON v.UserID = jf.JudgeID
        JOIN Tournaments t ON d.TournamentID = t.TournamentID
    WHERE
        v.UserID = $1 AND t.StartDate BETWEEN $2 AND $3
    GROUP BY
        t.StartDate, d.TournamentID, v.UserID
)
SELECT *
FROM tournament_performance
ORDER BY StartDate;

-- name: MarkStudentFeedbackAsRead :exec
UPDATE SpeakerScores
SET IsRead = true
WHERE SpeakerID = $1 AND BallotID = $2;

-- name: MarkJudgeFeedbackAsRead :exec
UPDATE JudgeFeedback
SET IsRead = true
WHERE FeedbackID = $1 AND JudgeID = $2;

-- name: GetTournamentVolunteerRanking :many
WITH RankedVolunteers AS (
    SELECT
        v.VolunteerID,
        v.FirstName || ' ' || v.LastName AS VolunteerName,
        COUNT(DISTINCT CASE WHEN d.IsEliminationRound = false THEN d.DebateID END) as PreliminaryRounds,
        COUNT(DISTINCT CASE WHEN d.IsEliminationRound = true THEN d.DebateID END) as EliminationRounds,
        COALESCE(AVG(jf.AverageRating), 0) as AverageRating,
        RANK() OVER (
            ORDER BY
                COALESCE(AVG(jf.AverageRating), 0) DESC,
                COUNT(DISTINCT CASE WHEN d.IsEliminationRound = true THEN d.DebateID END) DESC,
                COUNT(DISTINCT CASE WHEN d.IsEliminationRound = false THEN d.DebateID END) DESC
            ) as place
    FROM Volunteers v
             JOIN JudgeAssignments ja ON v.UserID = ja.JudgeID
             JOIN Debates d ON ja.DebateID = d.DebateID
             LEFT JOIN JudgeFeedback jf ON v.UserID = jf.JudgeID AND jf.DebateID = d.DebateID
    WHERE d.TournamentID = $1
    GROUP BY v.VolunteerID, v.FirstName, v.LastName
),
     TopThree AS (
         SELECT *
         FROM RankedVolunteers
         WHERE place <= 3
     ),
     SearchResults AS (
         SELECT *
         FROM RankedVolunteers
         WHERE $4::text IS NULL
            OR VolunteerName ILIKE '%' || $4 || '%'
     )
SELECT DISTINCT ON (place) *
FROM (
         SELECT * FROM TopThree
         UNION
         SELECT * FROM SearchResults WHERE place > 3
     ) combined
ORDER BY place
LIMIT $2 OFFSET $3;

-- name: GetTournamentVolunteerRankingCount :one
SELECT COUNT(DISTINCT v.VolunteerID)
FROM Volunteers v
JOIN JudgeAssignments ja ON v.UserID = ja.JudgeID
JOIN Debates d ON ja.DebateID = d.DebateID
WHERE d.TournamentID = $1;

-- name: GetRankingVisibility :one
SELECT IsVisible
FROM RankingVisibility
WHERE TournamentID = $1 AND RankingType = $2 AND VisibleTo = $3;

-- name: SetRankingVisibility :exec
INSERT INTO RankingVisibility (TournamentID, RankingType, VisibleTo, IsVisible)
VALUES ($1, $2, $3, $4)
ON CONFLICT (TournamentID, RankingType, VisibleTo)
    DO UPDATE SET IsVisible = $4, UpdatedAt = CURRENT_TIMESTAMP;

-- name: GetSpeakerScoresByBallotAndTeam :many
SELECT ss.ScoreID, ss.SpeakerID, ss.BallotID, ss.SpeakerRank,
       ss.SpeakerPoints, ss.Feedback, tm.TeamID
FROM SpeakerScores ss
         JOIN TeamMembers tm ON ss.SpeakerID = tm.StudentID
WHERE ss.BallotID = $1 AND tm.TeamID = $2
ORDER BY ss.SpeakerRank;

-- name: CreateSpeakerScore :exec
INSERT INTO SpeakerScores (BallotID, SpeakerID, SpeakerRank, SpeakerPoints, Feedback)
VALUES ($1, $2, $3, $4, $5);

-- name: TeamHasDebates :one
SELECT EXISTS (
    SELECT 1 FROM Debates
    WHERE Team1ID = $1 OR Team2ID = $1
) AS has_debates;

-- name: GetDebatesByTeam :many
SELECT * FROM Debates
WHERE Team1ID = $1 OR Team2ID = $1;

-- name: DeleteSpeakerScoresByBallotAndTeam :exec
DELETE FROM SpeakerScores
WHERE BallotID = $1
  AND SpeakerID IN (
    SELECT tm.StudentID
    FROM TeamMembers tm
    WHERE tm.TeamID = $2
);