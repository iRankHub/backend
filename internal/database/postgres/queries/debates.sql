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

-- name: CreateDebate :one
INSERT INTO Debates (TournamentID, RoundNumber, IsEliminationRound, Team1ID, Team2ID)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTeamsByTournament :many
SELECT t.TeamID, t.Name
FROM Teams t
JOIN TournamentInvitations ti ON t.InvitationID = ti.InvitationID
WHERE ti.TournamentID = $1 AND ti.Status = 'accepted';

-- name: GetPreviousPairings :many
SELECT Team1ID, Team2ID
FROM Debates
WHERE TournamentID = $1 AND RoundNumber < $2;

-- name: CreatePairingHistory :exec
INSERT INTO PairingHistory (TournamentID, Team1ID, Team2ID, RoundNumber, IsElimination)
VALUES ($1, $2, $3, $4, $5);