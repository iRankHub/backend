package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
    dbPairings, err := queries.GetPairingsByTournamentAndRound(ctx, models.GetPairingsByTournamentAndRoundParams{
        Tournamentid:       req.GetTournamentId(),
        Roundnumber:        req.GetRoundNumber(),
        Iseliminationround: req.GetIsElimination(),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get pairings: %v", err)
    }

    pairings := make([]*debate_management.Pairing, len(dbPairings))
    for i, dbPairing := range dbPairings {
        pairings[i] = convertSinglePairingFromRow(dbPairing)
        judges, err := s.getJudgesForPairing(ctx, dbPairing.Debateid)
        if err != nil {
            return nil, err
        }
        pairings[i].Judges = judges
    }

    return &debate_management.GetPairingsResponse{Pairings: pairings}, nil
}

func (s *PairingService) GetPairing(ctx context.Context, req *debate_management.GetPairingRequest) (*debate_management.GetPairingResponse, error) {
    if err := s.validateAuthentication(req.GetToken()); err != nil {
        return nil, err
    }

    queries := models.New(s.db)
    dbPairing, err := queries.GetPairingByID(ctx, req.GetPairingId())
    if err != nil {
        return nil, fmt.Errorf("failed to get pairing: %v", err)
    }

    pairing := convertSinglePairing(dbPairing)
    judges, err := s.getJudgesForPairing(ctx, dbPairing.Debateid)
    if err != nil {
        return nil, err
    }
    pairing.Judges = judges

    return &debate_management.GetPairingResponse{Pairing: pairing}, nil
}

func (s *PairingService) UpdatePairings(ctx context.Context, req *debate_management.UpdatePairingsRequest) ([]*debate_management.Pairing, error) {
    _, err := s.validateAdminRole(req.GetToken())
    if err != nil {
        return nil, err
    }

    queries := models.New(s.db)
    updatedPairings := make([]*debate_management.Pairing, 0, len(req.GetPairings()))

    for _, pairing := range req.GetPairings() {
        err = queries.UpdatePairing(ctx, models.UpdatePairingParams{
            Debateid: pairing.GetPairingId(),
            Team1id:  pairing.GetTeam1().GetTeamId(),
            Team2id:  pairing.GetTeam2().GetTeamId(),
            Roomid:   pairing.GetRoomId(),
        })
        if err != nil {
            return nil, fmt.Errorf("failed to update pairing: %v", err)
        }

        // Fetch the updated pairing
        updatedPairing, err := queries.GetPairingByID(ctx, pairing.GetPairingId())
        if err != nil {
            return nil, fmt.Errorf("failed to fetch updated pairing: %v", err)
        }

        updatedPairings = append(updatedPairings, convertSinglePairing(updatedPairing))
    }

    return updatedPairings, nil
}

func (s *PairingService) GeneratePairings(ctx context.Context, req *debate_management.GeneratePairingsRequest) ([]*debate_management.Pairing, error) {
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

    // Create rounds if they don't exist
    roundIDs := make(map[int32]int32)
    roundCount := tournament.Numberofpreliminaryrounds
    if req.GetIsEliminationRound() {
        roundCount = tournament.Numberofeliminationrounds
    }
    for roundNumber := 1; roundNumber <= int(roundCount); roundNumber++ {
        round, err := queries.CreateRound(ctx, models.CreateRoundParams{
            Tournamentid:       req.GetTournamentId(),
            Roundnumber:        int32(roundNumber),
            Iseliminationround: req.GetIsEliminationRound(),
        })
        if err != nil {
            if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
                existingRound, err := queries.GetRoundByTournamentAndNumber(ctx, models.GetRoundByTournamentAndNumberParams{
                    Tournamentid:       req.GetTournamentId(),
                    Roundnumber:        int32(roundNumber),
                    Iseliminationround: req.GetIsEliminationRound(),
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
    teams, err := queries.GetTeamsByTournament(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to get teams: %v", err)
    }

    // If the number of teams is odd, create a "Public Speaking" team
    if len(teams)%2 != 0 {
        publicSpeakingTeam, err := queries.CreateTeam(ctx, models.CreateTeamParams{
            Name:         "Public Speaking",
            Tournamentid: req.GetTournamentId(),
        })
        if err != nil {
            return nil, fmt.Errorf("failed to create Public Speaking team: %v", err)
        }

        // Fetch additional information for the new team
        newTeamInfo, err := queries.GetTeamsByTournament(ctx, req.GetTournamentId())
        if err != nil {
            return nil, fmt.Errorf("failed to fetch new team information: %v", err)
        }

        var publicSpeakingTeamRow models.GetTeamsByTournamentRow
        for _, team := range newTeamInfo {
            if team.Teamid == publicSpeakingTeam.Teamid {
                publicSpeakingTeamRow = team
                break
            }
        }

        if publicSpeakingTeamRow.Teamid == 0 {
            return nil, fmt.Errorf("failed to find newly created team in tournament teams")
        }

        teams = append(teams, publicSpeakingTeamRow)
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
    judges, err := queries.GetAvailableJudges(ctx, req.GetTournamentId())
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
    rooms, err := queries.GetRoomsByTournament(ctx, sql.NullInt32{Int32: req.GetTournamentId(), Valid: true})
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
                Tournamentid: sql.NullInt32{Int32: req.GetTournamentId(), Valid: true},
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
        PreliminaryRounds: int(roundCount),
        JudgesPerDebate:   int(tournament.Judgesperdebatepreliminary),
    }

    allPairings, err := pairing_algorithm.GeneratePairings(algorithmTeams, algorithmJudges, roomIDs, specs)
    if err != nil {
        return nil, fmt.Errorf("failed to generate pairings: %v", err)
    }

    // Save new pairings to the database
  dbPairings := make([]*debate_management.Pairing, 0)
    for roundNumber, roundPairings := range allPairings {
        for _, pair := range roundPairings {
            startTime := time.Now().Add(time.Duration(roundNumber) * time.Hour)

            roundID, ok := roundIDs[int32(roundNumber+1)]
            if !ok {
                return nil, fmt.Errorf("round ID not found for round number %d", roundNumber+1)
            }

            debate, err := queries.CreateDebate(ctx, models.CreateDebateParams{
                Roundid:            roundID,
                Tournamentid:       req.GetTournamentId(),
                Roundnumber:        int32(roundNumber + 1),
                Iseliminationround: req.GetIsEliminationRound(),
                Team1id:            int32(pair.Team1.ID),
                Team2id:            int32(pair.Team2.ID),
                Roomid:             int32(pair.Room),
                Starttime:          startTime,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create debate: %v", err)
            }

            // Assign judges
            for i, judge := range pair.Judges {
                isHeadJudge := (i == 0) || (len(pair.Judges) == 1)
                err := queries.AssignJudgeToDebate(ctx, models.AssignJudgeToDebateParams{
                    Tournamentid:  req.GetTournamentId(),
                    Judgeid:       int32(judge.ID),
                    Debateid:      debate,
                    Roundnumber:   int32(roundNumber + 1),
                    Iselimination: req.GetIsEliminationRound(),
                    Isheadjudge:   isHeadJudge,
                })
                if err != nil {
                    return nil, fmt.Errorf("failed to assign judge to debate: %v", err)
                }
            }

            // Fetch room name
            room, err := queries.GetRoomByID(ctx, int32(pair.Room))
            if err != nil {
                return nil, fmt.Errorf("failed to get room: %v", err)
            }

            dbPairings = append(dbPairings, &debate_management.Pairing{
                PairingId:          debate,
                RoundNumber:        int32(roundNumber + 1),
                IsEliminationRound: req.GetIsEliminationRound(),
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
                Tournamentid:  req.GetTournamentId(),
                Team1id:       int32(pair.Team1.ID),
                Team2id:       int32(pair.Team2.ID),
                Roundnumber:   int32(roundNumber + 1),
                Iselimination: req.GetIsEliminationRound(),
            })
            if err != nil {
                return nil, fmt.Errorf("failed to record pairing history: %v", err)
            }
        }
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return dbPairings, nil
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

func (s *PairingService) RegeneratePairings(ctx context.Context, req *debate_management.RegeneratePairingsRequest) ([]*debate_management.Pairing, error) {
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

    // Delete existing data in the specified order
    err = queries.DeleteJudgeAssignmentsForTournament(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to delete judge assignments: %v", err)
    }

    err = queries.DeleteDebatesForTournament(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to delete debates: %v", err)
    }

    err = queries.DeleteRoomsForTournament(ctx, sql.NullInt32{Int32: req.GetTournamentId(), Valid: true})
    if err != nil {
        return nil, fmt.Errorf("failed to delete rooms: %v", err)
    }


    err = queries.DeleteRoundsForTournament(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to delete rounds: %v", err)
    }

    err = queries.DeletePairingHistoryForTournament(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to delete pairing history: %v", err)
    }

    // Generate new pairings
    newPairings, err := s.GeneratePairings(ctx, &debate_management.GeneratePairingsRequest{
        TournamentId:  req.GetTournamentId(),
        Token:         req.GetToken(),
        IsEliminationRound: req.GetIsEliminationRound(),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to generate new pairings: %v", err)
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return newPairings, nil
}

func convertSinglePairing(dbPairing models.GetPairingByIDRow) *debate_management.Pairing {
	return convertSinglePairingFromRow(models.GetPairingsByTournamentAndRoundRow(dbPairing))
}

func convertSinglePairingFromRow(dbPairing models.GetPairingsByTournamentAndRoundRow) *debate_management.Pairing {
    return &debate_management.Pairing{
        PairingId:          dbPairing.Debateid,
        RoundNumber:        dbPairing.Roundnumber,
        IsEliminationRound: dbPairing.Iseliminationround,
        RoomId:             dbPairing.Roomid,
        RoomName:           dbPairing.Roomname.String,
        Team1:              convertTeamPairing(dbPairing.Team1id, dbPairing.Team1name, parseSpeakerNames(dbPairing.Team1speakernames), dbPairing.Team1leaguename.String, float64(dbPairing.Team1totalpoints)),
        Team2:              convertTeamPairing(dbPairing.Team2id, dbPairing.Team2name, parseSpeakerNames(dbPairing.Team2speakernames), dbPairing.Team2leaguename.String, float64(dbPairing.Team2totalpoints)),
        HeadJudgeName:      dbPairing.Headjudgename,
    }
}

func convertTeamPairing(teamID int32, teamName string, speakerNames []string, leagueName string, totalPoints float64) *debate_management.Team {
    return &debate_management.Team{
        TeamId:       teamID,
        Name:         teamName,
        SpeakerNames: speakerNames,
        TotalPoints:  totalPoints,
        LeagueName:   leagueName,
    }
}

func parseSpeakerNames(speakerNamesData interface{}) []string {
    var speakerNames []string

    switch v := speakerNamesData.(type) {
    case []string:
        speakerNames = v
    case string:
        err := json.Unmarshal([]byte(v), &speakerNames)
        if err != nil {
            speakerNames = strings.Split(v, ",")
            for i, name := range speakerNames {
                speakerNames[i] = strings.TrimSpace(name)
            }
        }
    }

    return speakerNames
}

func (s *PairingService) getJudgesForPairing(ctx context.Context, debateID int32) ([]*debate_management.Judge, error) {
    queries := models.New(s.db)
    dbJudges, err := queries.GetJudgesForDebate(ctx, debateID)
    if err != nil {
        return nil, fmt.Errorf("failed to get judges for debate: %v", err)
    }

    judges := make([]*debate_management.Judge, len(dbJudges))
    for i, dbJudge := range dbJudges {
        judges[i] = &debate_management.Judge{
            JudgeId: dbJudge.Judgeid,
            Name:    dbJudge.Name,
        }
    }
    return judges, nil
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
