package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
	"github.com/iRankHub/backend/pkg/pairing_algorithm"
)

type PairingService struct {
	db *sql.DB
}

func NewPairingService(db *sql.DB) *PairingService {
	return &PairingService{db: db}
}

func (s *PairingService) GetPairings(ctx context.Context, req *debate_management.GetPairingsRequest) (*debate_management.GetPairingsResponse, error) {
    if err := s.validateAuthentication(req.GetToken()); err != nil {
        return nil, err
    }

    queries := models.New(s.db)
    dbPairings, err := queries.GetPairings(ctx, models.GetPairingsParams{
        Tournamentid:       req.GetTournamentId(),
        Roundnumber:        req.GetRoundNumber(),
        Iseliminationround: req.GetIsElimination(),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get pairings: %v", err)
    }

    pairings := make([]*debate_management.Pairing, len(dbPairings))
    for i, dbPairing := range dbPairings {
        pairings[i] = &debate_management.Pairing{
            PairingId:          dbPairing.Debateid,
            RoundNumber:        dbPairing.Roundnumber,
            IsEliminationRound: dbPairing.Iseliminationround,
            RoomId:             dbPairing.Roomid,
            RoomName:           dbPairing.Roomname,
            Team1: &debate_management.Team{
                TeamId: dbPairing.Team1id,
                Name:   dbPairing.Team1name,
            },
            Team2: &debate_management.Team{
                TeamId: dbPairing.Team2id,
                Name:   dbPairing.Team2name,
            },
            HeadJudgeName: dbPairing.Headjudgename.String,
        }
    }

    return &debate_management.GetPairingsResponse{Pairings: pairings}, nil
}


func (s *PairingService) UpdatePairings(ctx context.Context, req *debate_management.UpdatePairingsRequest) (*debate_management.UpdatePairingsResponse, error) {
    _, err := s.validateAdminRole(req.GetToken())
    if err != nil {
        return nil, err
    }

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(s.db).WithTx(tx)

    updatedPairings := make([]*debate_management.Pairing, 0, len(req.GetPairings()))

    for _, pairing := range req.GetPairings() {
        err = queries.UpdatePairing(ctx, models.UpdatePairingParams{
            Debateid: pairing.GetPairingId(),
            Team1id:  pairing.GetTeam1().GetTeamId(),
            Team2id:  pairing.GetTeam2().GetTeamId(),
        })
        if err != nil {
            return nil, fmt.Errorf("failed to update pairing: %v", err)
        }

        // Fetch the updated pairing
        updatedPairing, err := queries.GetSinglePairing(ctx, pairing.GetPairingId())
        if err != nil {
            return nil, fmt.Errorf("failed to fetch updated pairing: %v", err)
        }

        updatedPairings = append(updatedPairings, &debate_management.Pairing{
            PairingId:          updatedPairing.Debateid,
            RoundNumber:        updatedPairing.Roundnumber,
            IsEliminationRound: updatedPairing.Iseliminationround,
            RoomId:             updatedPairing.Roomid,
            RoomName:           updatedPairing.Roomname,
            Team1: &debate_management.Team{
                TeamId: updatedPairing.Team1id,
                Name:   updatedPairing.Team1name,
            },
            Team2: &debate_management.Team{
                TeamId: updatedPairing.Team2id,
                Name:   updatedPairing.Team2name,
            },
            HeadJudgeName: updatedPairing.Headjudgename.String,
        })
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &debate_management.UpdatePairingsResponse{Pairings: updatedPairings}, nil
}

func (s *PairingService) GeneratePreliminaryPairings(ctx context.Context, req *debate_management.GeneratePreliminaryPairingsRequest) (*debate_management.GeneratePairingsResponse, error) {
    _, err := s.validateAdminRole(req.GetToken())
    if err != nil {
        return nil, err
    }

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(s.db).WithTx(tx)

    // Get tournament details
    tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to get tournament: %v", err)
    }

    pairings, err := s.generatePreliminaryPairings(ctx, queries, tournament, req.GetTournamentId())
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &debate_management.GeneratePairingsResponse{Pairings: pairings}, nil
}

func (s *PairingService) GenerateEliminationPairings(ctx context.Context, req *debate_management.GenerateEliminationPairingsRequest) (*debate_management.GeneratePairingsResponse, error) {
    _, err := s.validateAdminRole(req.GetToken())
    if err != nil {
        return nil, err
    }

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(s.db).WithTx(tx)

    // Get tournament details
    tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to get tournament: %v", err)
    }

    pairings, err := s.generateEliminationPairings(ctx, queries, tournament, req.GetTournamentId(), req.GetRoundNumber())
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &debate_management.GeneratePairingsResponse{Pairings: pairings}, nil
}

func (s *PairingService) generatePreliminaryPairings(ctx context.Context, queries *models.Queries, tournament models.GetTournamentByIDRow, tournamentID int32) ([]*debate_management.Pairing, error) {
    // Create rounds if they don't exist
    roundIDs := make(map[int32]int32)
    for roundNumber := 1; roundNumber <= int(tournament.Numberofpreliminaryrounds); roundNumber++ {
        round, err := queries.CreateRound(ctx, models.CreateRoundParams{
            Tournamentid:       tournamentID,
            Roundnumber:        int32(roundNumber),
            Iseliminationround: false,
        })
        if err != nil {
            if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
                existingRound, err := queries.GetRoundByTournamentAndNumber(ctx, models.GetRoundByTournamentAndNumberParams{
                    Tournamentid:       tournamentID,
                    Roundnumber:        int32(roundNumber),
                    Iseliminationround: false,
                })
                if err != nil {
                    return nil, fmt.Errorf("failed to get existing round: %v", err)
                }
                roundIDs[int32(roundNumber)] = existingRound.Roundid
            } else {
                return nil, fmt.Errorf("failed to create round: %v", err)
            }
        } else {
            roundIDs[int32(roundNumber)] = round.Roundid
        }
    }

    // Get teams for the tournament
    teams, err := queries.GetTeamsByTournament(ctx, tournamentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get teams: %v", err)
    }

    // If the number of teams is odd, create a "Public Speaking" team
    if len(teams)%2 != 0 {
        publicSpeakingTeam, err := queries.CreateTeam(ctx, models.CreateTeamParams{
            Name:         "Public Speaking",
            Tournamentid: tournamentID,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to create Public Speaking team: %v", err)
        }
        teams = append(teams, models.GetTeamsByTournamentRow{
            Teamid: publicSpeakingTeam.Teamid,
            Name:   publicSpeakingTeam.Name,
        })
    }

    // Convert teams to the format expected by the pairing algorithm
    algorithmTeams := make([]*pairing_algorithm.Team, len(teams))
    for i, team := range teams {
        algorithmTeams[i] = &pairing_algorithm.Team{
            ID:        int(team.Teamid),
            Name:      team.Name,
            Opponents: make(map[int]bool),
        }
    }

    // Get available judges
    judges, err := queries.GetAvailableJudges(ctx, tournamentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get available judges: %v", err)
    }

    algorithmJudges := make([]*pairing_algorithm.Judge, len(judges))
    for i, judge := range judges {
        algorithmJudges[i] = &pairing_algorithm.Judge{
            ID:   int(judge.Userid),
            Name: judge.Name,
        }
    }

    // Create or get rooms
    rooms, err := queries.GetRoomsByTournament(ctx, sql.NullInt32{Int32: tournamentID, Valid: true})
    if err != nil {
        return nil, fmt.Errorf("failed to get rooms: %v", err)
    }

    if len(rooms) < len(teams)/2 {
        neededRooms := len(teams)/2 - len(rooms)
        for i := 0; i < neededRooms; i++ {
            roomName := fmt.Sprintf("Room %d", len(rooms)+i+1)
            room, err := queries.CreateRoom(ctx, models.CreateRoomParams{
                Roomname:     roomName,
                Location:     "TBD",
                Capacity:     int32(tournament.Judgesperdebatepreliminary + 12),
                Tournamentid: sql.NullInt32{Int32: tournamentID, Valid: true},
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create room: %v", err)
            }
            rooms = append(rooms, room)
        }
    }
    roomIDs := make([]int, len(rooms))
    for i, room := range rooms {
        roomIDs[i] = int(room.Roomid)
    }

    // Generate pairings using the pairing algorithm
    specs := pairing_algorithm.TournamentSpecs{
        PreliminaryRounds: int(tournament.Numberofpreliminaryrounds),
        JudgesPerDebate:   int(tournament.Judgesperdebatepreliminary),
    }

    debates, err := pairing_algorithm.GeneratePairings(algorithmTeams, algorithmJudges, roomIDs, specs, 0, false)
    if err != nil {
        return nil, fmt.Errorf("failed to generate pairings: %v", err)
    }

    // Save new pairings to the database
    dbPairings := make([]*debate_management.Pairing, 0, len(debates))
    debateIndex := 0
    for roundNumber := 1; roundNumber <= int(tournament.Numberofpreliminaryrounds); roundNumber++ {
        for i := 0; i < len(teams)/2; i++ {
            pair := debates[debateIndex]
            startTime := time.Now().Add(time.Duration(roundNumber) * time.Hour)

            roundID, ok := roundIDs[int32(roundNumber)]
            if !ok {
                return nil, fmt.Errorf("round ID not found for round number %d", roundNumber)
            }

            debate, err := queries.CreateDebate(ctx, models.CreateDebateParams{
                Roundid:            roundID,
                Tournamentid:       tournamentID,
                Roundnumber:        int32(roundNumber),
                Iseliminationround: false,
                Team1id:            int32(pair.Team1.ID),
                Team2id:            int32(pair.Team2.ID),
                Roomid:             int32(pair.Room),
                Starttime:          startTime,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create debate: %v", err)
            }

            // Assign judges
            var headJudgeID int32
            for i, judge := range pair.Judges {
                isHeadJudge := (i == 0) || (len(pair.Judges) == 1)
                err := queries.AssignJudgeToDebate(ctx, models.AssignJudgeToDebateParams{
                    Tournamentid:  tournamentID,
                    Judgeid:       int32(judge.ID),
                    Debateid:      debate,
                    Roundnumber:   int32(roundNumber),
                    Iselimination: false,
                    Isheadjudge:   isHeadJudge,
                })
                if err != nil {
                    return nil, fmt.Errorf("failed to assign judge to debate: %v", err)
                }
                if isHeadJudge {
                    headJudgeID = int32(judge.ID)
                }
            }

            // Create a ballot for the debate
            _, err = queries.CreateBallot(ctx, models.CreateBallotParams{
                Debateid:        debate,
                Judgeid:         headJudgeID,
                Recordingstatus: "not yet",
                Verdict:         "pending",
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create ballot: %v", err)
            }

            // Create initial speaker scores
            err = queries.CreateInitialSpeakerScores(ctx, debate)
            if err != nil {
                return nil, fmt.Errorf("failed to create initial speaker scores: %v", err)
            }

            // Fetch room name
            room, err := queries.GetRoomByID(ctx, int32(pair.Room))
            if err != nil {
                return nil, fmt.Errorf("failed to get room: %v", err)
            }

            dbPairings = append(dbPairings, &debate_management.Pairing{
                PairingId:          debate,
                RoundNumber:        int32(roundNumber),
                IsEliminationRound: false,
                RoomId:             int32(pair.Room),
                RoomName:           room.Roomname,
                Team1: &debate_management.Team{
                    TeamId: int32(pair.Team1.ID),
                    Name:   pair.Team1.Name,
                },
                Team2: &debate_management.Team{
                    TeamId: int32(pair.Team2.ID),
                    Name:   pair.Team2.Name,
                },
                Judges: convertJudgesToProto(pair.Judges),
            })

            // Record pairing history
            err = queries.CreatePairingHistory(ctx, models.CreatePairingHistoryParams{
                Tournamentid:  tournamentID,
                Team1id:       int32(pair.Team1.ID),
                Team2id:       int32(pair.Team2.ID),
                Roundnumber:   int32(roundNumber),
                Iselimination: false,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to record pairing history: %v", err)
            }

            debateIndex++
        }
    }

    return dbPairings, nil
}

func (s *PairingService) generateEliminationPairings(ctx context.Context, queries *models.Queries, tournament models.GetTournamentByIDRow, tournamentID int32, roundNumber int32) ([]*debate_management.Pairing, error) {

    // Get tournament details
    tournament, err := queries.GetTournamentByID(ctx, tournamentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get tournament: %v", err)
    }

    // Create or get the elimination round
    round, err := queries.CreateRound(ctx, models.CreateRoundParams{
        Tournamentid:       tournamentID,
        Roundnumber:        roundNumber,
        Iseliminationround: true,
    })
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
            existingRound, err := queries.GetRoundByTournamentAndNumber(ctx, models.GetRoundByTournamentAndNumberParams{
                Tournamentid:       tournamentID,
                Roundnumber:        roundNumber,
                Iseliminationround: true,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to get existing round: %v", err)
            }
            round = existingRound
        } else {
            return nil, fmt.Errorf("failed to create round: %v", err)
        }
    }

    // Get teams for the elimination round
    var teams []models.GetTeamsByTournamentRow
    teamsNeededForElims := int32(math.Pow(2, float64(tournament.Numberofeliminationrounds-roundNumber+1)))

    if roundNumber == 1 {
        // For the first elimination round, get top performing teams from preliminaries
        topTeams, err := queries.GetTopPerformingTeams(ctx, models.GetTopPerformingTeamsParams{
            Tournamentid: tournamentID,
            Limit:        teamsNeededForElims,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to get top performing teams: %v", err)
        }
        teams = convertTopPerformingTeamsToTeamsByTournament(topTeams)
    } else {
        // For subsequent elimination rounds, get winning teams from the previous round
        winningTeams, err := queries.GetEliminationRoundTeams(ctx, models.GetEliminationRoundTeamsParams{
            Tournamentid: tournamentID,
            Roundnumber:  roundNumber - 1,
            Limit:        teamsNeededForElims,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to get winning teams from previous elimination round: %v", err)
        }
        teams = convertEliminationRoundTeamsToTeamsByTournament(winningTeams)
    }

    // Convert teams to the format expected by the pairing algorithm
    algorithmTeams := make([]*pairing_algorithm.Team, len(teams))
    for i, team := range teams {
        algorithmTeams[i] = &pairing_algorithm.Team{
            ID:        int(team.Teamid),
            Name:      team.Name,
            Opponents: make(map[int]bool),
        }
    }

    // Get available judges
    judges, err := queries.GetAvailableJudges(ctx, tournamentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get available judges: %v", err)
    }

    algorithmJudges := make([]*pairing_algorithm.Judge, len(judges))
    for i, judge := range judges {
        algorithmJudges[i] = &pairing_algorithm.Judge{
            ID:   int(judge.Userid),
            Name: judge.Name,
        }
    }

    // Create or get rooms
    rooms, err := queries.GetRoomsByTournament(ctx, sql.NullInt32{Int32: tournamentID, Valid: true})
    if err != nil {
        return nil, fmt.Errorf("failed to get rooms: %v", err)
    }

    if len(rooms) < len(teams)/2 {
        neededRooms := len(teams)/2 - len(rooms)
        for i := 0; i < neededRooms; i++ {
            roomName := fmt.Sprintf("Room %d", len(rooms)+i+1)
            room, err := queries.CreateRoom(ctx, models.CreateRoomParams{
                Roomname:     roomName,
                Location:     "TBD",
                Capacity:     int32(tournament.Judgesperdebateelimination + 12),
                Tournamentid: sql.NullInt32{Int32: tournamentID, Valid: true},
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create room: %v", err)
            }
            rooms = append(rooms, room)
        }
    }
    roomIDs := make([]int, len(rooms))
    for i, room := range rooms {
        roomIDs[i] = int(room.Roomid)
    }

    // Generate pairings using the pairing algorithm
    specs := pairing_algorithm.TournamentSpecs{
        PreliminaryRounds:     int(tournament.Numberofpreliminaryrounds),
        EliminationRounds:     int(tournament.Numberofeliminationrounds),
        JudgesPerDebate:       int(tournament.Judgesperdebateelimination),
        TeamsAdvancingToElims: int(math.Pow(2, float64(tournament.Numberofeliminationrounds))),
    }

    debates, err := pairing_algorithm.GeneratePairings(algorithmTeams, algorithmJudges, roomIDs, specs, int(roundNumber), true)
    if err != nil {
        return nil, fmt.Errorf("failed to generate elimination pairings: %v", err)
    }

    dbPairings, err := s.saveDebatesToDatabase(ctx, queries, debates, tournamentID, roundNumber, true, round.Roundid)
    if err != nil {
        return nil, fmt.Errorf("failed to save debates to database: %v", err)
    }

    return dbPairings, nil
}


func (s *PairingService) saveDebatesToDatabase(ctx context.Context, queries *models.Queries, debates []*pairing_algorithm.Debate, tournamentID int32, roundNumber int32, isElimination bool, roundID int32) ([]*debate_management.Pairing, error) {
    dbPairings := make([]*debate_management.Pairing, 0, len(debates))
    for _, debate := range debates {
        startTime := time.Now().Add(time.Hour) // Adjust as needed

        dbDebate, err := queries.CreateDebate(ctx, models.CreateDebateParams{
            Roundid:            roundID,
            Tournamentid:       tournamentID,
            Roundnumber:        roundNumber,
            Iseliminationround: isElimination,
            Team1id:            int32(debate.Team1.ID),
            Team2id:            int32(debate.Team2.ID),
            Roomid:             int32(debate.Room),
            Starttime:          startTime,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to create debate: %v", err)
        }

        // Assign judges
        var headJudgeID int32
        for i, judge := range debate.Judges {
            isHeadJudge := (i == 0) || (len(debate.Judges) == 1)
            err := queries.AssignJudgeToDebate(ctx, models.AssignJudgeToDebateParams{
                Tournamentid:  tournamentID,
                Judgeid:       int32(judge.ID),
                Debateid:      dbDebate,
                Roundnumber:   roundNumber,
                Iselimination: isElimination,
                Isheadjudge:   isHeadJudge,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to assign judge to debate: %v", err)
            }
            if isHeadJudge {
                headJudgeID = int32(judge.ID)
            }
        }

        // Create a ballot for the debate
        _, err = queries.CreateBallot(ctx, models.CreateBallotParams{
            Debateid:        dbDebate,
            Judgeid:         headJudgeID,
            Recordingstatus: "not yet",
            Verdict:         "pending",
        })
        if err != nil {
            // If creating the ballot fails, log the error but continue with the process
            log.Printf("Warning: Failed to create ballot for debate %d: %v", dbDebate, err)
        }

        // Create initial speaker scores
        err = queries.CreateInitialSpeakerScores(ctx, dbDebate)
        if err != nil {
            return nil, fmt.Errorf("failed to create initial speaker scores: %v", err)
        }

        // Fetch room name
        room, err := queries.GetRoomByID(ctx, int32(debate.Room))
        if err != nil {
            return nil, fmt.Errorf("failed to get room: %v", err)
        }

        dbPairings = append(dbPairings, &debate_management.Pairing{
            PairingId:          dbDebate,
            RoundNumber:        roundNumber,
            IsEliminationRound: isElimination,
            RoomId:             int32(debate.Room),
            RoomName:           room.Roomname,
            Team1: &debate_management.Team{
                TeamId: int32(debate.Team1.ID),
                Name:   debate.Team1.Name,
            },
            Team2: &debate_management.Team{
                TeamId: int32(debate.Team2.ID),
                Name:   debate.Team2.Name,
            },
            Judges: convertJudgesToProto(debate.Judges),
        })

        // Record pairing history
        err = queries.CreatePairingHistory(ctx, models.CreatePairingHistoryParams{
            Tournamentid:  tournamentID,
            Team1id:       int32(debate.Team1.ID),
            Team2id:       int32(debate.Team2.ID),
            Roundnumber:   roundNumber,
            Iselimination: isElimination,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to record pairing history: %v", err)
        }
    }

    return dbPairings, nil
}

func convertTopPerformingTeamsToTeamsByTournament(topTeams []models.GetTopPerformingTeamsRow) []models.GetTeamsByTournamentRow {
    result := make([]models.GetTeamsByTournamentRow, len(topTeams))
    for i, team := range topTeams {
        result[i] = models.GetTeamsByTournamentRow{
            Teamid:       team.Teamid,
            Name:         team.Name,
            Tournamentid: team.Tournamentid,
        }
    }
    return result
}

func convertEliminationRoundTeamsToTeamsByTournament(elimTeams []models.GetEliminationRoundTeamsRow) []models.GetTeamsByTournamentRow {
    result := make([]models.GetTeamsByTournamentRow, len(elimTeams))
    for i, team := range elimTeams {
        result[i] = models.GetTeamsByTournamentRow{
            Teamid:       team.Teamid,
            Name:         team.Name,
            Tournamentid: team.Tournamentid,
        }
    }
    return result
}

func convertJudgesToProto(judges []*pairing_algorithm.Judge) []*debate_management.Judge {
    protoJudges := make([]*debate_management.Judge, len(judges))
    for i, judge := range judges {
        protoJudges[i] = &debate_management.Judge{
            JudgeId: int32(judge.ID),
            Name:    judge.Name,
        }
    }
    return protoJudges
}

func (s *PairingService) validateAuthentication(token string) error {
    _, err := utils.ValidateToken(token)
    if err != nil {
        return fmt.Errorf("authentication failed: %v", err)
    }
    return nil
}

func (s *PairingService) validateAdminRole(token string) (map[string]interface{}, error) {
    claims, err := utils.ValidateToken(token)
    if err != nil {
        return nil, fmt.Errorf("authentication failed: %v", err)
    }

    userRole, ok := claims["user_role"].(string)
    if !ok || userRole != "admin" {
        return nil, fmt.Errorf("unauthorized: only admins can perform this action")
    }

    return claims, nil
}