-- name: GetRoomsByTournamentAndRound :many
SELECT r.RoomID, r.roomname, rb.RoundNumber, rb.IsElimination, rb.IsOccupied
FROM Rooms r
JOIN RoomBookings rb ON r.RoomID = rb.RoomID
WHERE rb.TournamentID = $1 AND rb.RoundNumber = $2 AND rb.IsElimination = $3;

-- name: GetRoomByID :many
SELECT r.RoomID, r.RoomName, rb.RoundNumber, rb.IsElimination, rb.IsOccupied
FROM Rooms r
LEFT JOIN RoomBookings rb ON r.RoomID = rb.RoomID
WHERE r.RoomID = $1;

-- name: CreateRoom :one
INSERT INTO Rooms (RoomName, Location, Capacity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRoom :one
UPDATE Rooms
SET RoomName = $2
WHERE RoomID = $1
RETURNING *;

-- name: GetAvailableRooms :many
SELECT r.*
FROM Rooms r
LEFT JOIN RoomBookings rb ON r.RoomID = rb.RoomID
  AND rb.TournamentID = $1
  AND rb.RoundNumber = $2
  AND rb.IsElimination = $3
WHERE rb.RoomID IS NULL OR rb.IsOccupied = FALSE;

-- name: GetDebatesWithoutRooms :many
SELECT *
FROM Debates
WHERE TournamentID = $1 AND RoundNumber = $2 AND IsEliminationRound = $3 AND RoomID IS NULL;

-- name: AssignRoomToDebate :exec
UPDATE Debates
SET RoomID = $2
WHERE DebateID = $1;

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

-- name: GetPairingsByTournamentAndRound :many
SELECT d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       d.RoomID, r.roomname AS RoomName
FROM Debates d
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
WHERE d.TournamentID = $1 AND d.RoundNumber = $2 AND d.IsEliminationRound = $3;

-- name: GetPairingByID :one
SELECT d.DebateID, d.RoundNumber, d.IsEliminationRound,
       d.Team1ID, t1.Name AS Team1Name, d.Team2ID, t2.Name AS Team2Name,
       d.RoomID, r.roomname AS RoomName
FROM Debates d
JOIN Teams t1 ON d.Team1ID = t1.TeamID
JOIN Teams t2 ON d.Team2ID = t2.TeamID
LEFT JOIN Rooms r ON d.RoomID = r.RoomID
WHERE d.DebateID = $1;

-- name: UpdatePairing :exec
UPDATE Debates
SET Team1ID = $2, Team2ID = $3, RoomID = $4
WHERE DebateID = $1;

-- name: GetBallotsByTournamentAndRound :many
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
WHERE d.TournamentID = $1 AND d.RoundNumber = $2 AND d.IsEliminationRound = $3;

-- name: GetBallotByID :one
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
WHERE b.BallotID = $1;

-- name: UpdateBallot :exec
UPDATE Ballots
SET Team1TotalScore = $2, Team2TotalScore = $3, RecordingStatus = $4, Verdict = $5
WHERE BallotID = $1;

-- name: GetSpeakerScoresByBallot :many
SELECT ss.ScoreID, ss.SpeakerID, s.FirstName, s.LastName, ss.SpeakerRank, ss.SpeakerPoints, ss.Feedback
FROM SpeakerScores ss
JOIN Students s ON ss.SpeakerID = s.StudentID
WHERE ss.BallotID = $1;

-- name: UpdateSpeakerScore :exec
UPDATE SpeakerScores
SET SpeakerRank = $2, SpeakerPoints = $3, Feedback = $4
WHERE ScoreID = $1;

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

-- name: GetRoomsByTournament :many
SELECT r.RoomID, r.RoomName, r.Location, r.Capacity
FROM Rooms r
JOIN RoomBookings rb ON r.RoomID = rb.RoomID
WHERE rb.TournamentID = $1;

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

-- name: DeleteRoomBookingsForTournament :exec
DELETE FROM Rooms;

-- name: DeleteRoundsForTournament :exec
DELETE FROM Rounds
WHERE TournamentID = $1;

-- name: DeletePairingHistoryForTournament :exec
DELETE FROM PairingHistory
WHERE TournamentID = $1;