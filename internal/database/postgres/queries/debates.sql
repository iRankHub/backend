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
SELECT u.UserID, u.Name, u.Email
FROM Users u
JOIN Volunteers v ON u.UserID = v.UserID
LEFT JOIN JudgeAssignments ja ON u.UserID = ja.JudgeID AND ja.TournamentID = $1
WHERE v.Role = 'Judge' AND ja.JudgeID IS NULL;

-- name: DeletePairingsForTournament :exec
DELETE FROM Debates
WHERE TournamentID = $1;

-- name: UpdatePairing :exec
UPDATE Debates
SET Team1ID = $2, Team2ID = $3, RoomID = $4
WHERE DebateID = $1;

-- name: GetBallotsByTournamentAndRound :many
SELECT b.BallotID, d.RoundNumber, d.IsEliminationRound, r.RoomName,
       u.Name AS HeadJudgeName, b.RecordingStatus, b.Verdict
FROM Ballots b
JOIN Debates d ON b.DebateID = d.DebateID
JOIN Rooms r ON d.RoomID = r.RoomID
JOIN Users u ON b.JudgeID = u.UserID
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

-- name: UpdateBallot :exec
UPDATE Ballots
SET Team1TotalScore = $2, Team2TotalScore = $3, RecordingStatus = $4, Verdict = $5,
    Team1Feedback = $6, Team2Feedback = $7, last_updated_by = $8,
    last_updated_at = CURRENT_TIMESTAMP, head_judge_submitted = $9
WHERE BallotID = $1;

-- name: GetSpeakerScoresByBallot :many
SELECT ss.ScoreID, ss.SpeakerID, s.FirstName, s.LastName,
       ss.SpeakerRank, ss.SpeakerPoints, ss.Feedback,
       t.TeamID, t.Name AS TeamName
FROM SpeakerScores ss
JOIN Students s ON ss.SpeakerID = s.StudentID
JOIN TeamMembers tm ON s.StudentID = tm.StudentID
JOIN Teams t ON tm.TeamID = t.TeamID
JOIN Debates d ON t.TournamentID = d.TournamentID
JOIN Ballots b ON d.DebateID = b.DebateID
WHERE b.BallotID = $1;


-- name: UpdateSpeakerScore :exec
UPDATE SpeakerScores
SET SpeakerRank = $2, SpeakerPoints = $3, Feedback = $4
WHERE ScoreID = $1;

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
INSERT INTO Ballots (DebateID, JudgeID, RecordingStatus, Verdict)
VALUES ($1, $2, $3, $4)
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
    v.iDebateVolunteerID
FROM
    Users u
JOIN
    Volunteers v ON u.UserID = v.UserID
JOIN
    JudgeAssignments ja ON u.UserID = ja.JudgeID
WHERE
    ja.TournamentID = $1
GROUP BY
    u.UserID, v.iDebateVolunteerID;

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
    r.RoomName
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
UPDATE JudgeAssignments ja
SET DebateID = (
    SELECT d.DebateID
    FROM Debates d
    WHERE d.TournamentID = $2 AND d.RoundNumber = $3 AND d.RoomID = $4
)
WHERE ja.JudgeID = $1
  AND ja.TournamentID = $2
  AND ja.DebateID IN (
    SELECT d.DebateID
    FROM Debates d
    WHERE d.TournamentID = $2 AND d.RoundNumber = $3
  );