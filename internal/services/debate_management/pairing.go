package services

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *PairingService) GetPairings(ctx context.Context, req *debate_management.GetPairingsRequest) ([]*debate_management.Pairing, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	pairings, err := queries.GetPairingsByTournamentAndRound(ctx, models.GetPairingsByTournamentAndRoundParams{
		Tournamentid:       req.GetTournamentId(),
		Roundnumber:        req.GetRoundNumber(),
		Iseliminationround: req.GetIsElimination(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pairings: %v", err)
	}

	return convertPairings(pairings), nil
}

func (s *PairingService) GetPairing(ctx context.Context, req *debate_management.GetPairingRequest) (*debate_management.Pairing, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	pairing, err := queries.GetPairingByID(ctx, req.GetPairingId())
	if err != nil {
		return nil, fmt.Errorf("failed to get pairing: %v", err)
	}

	return convertSinglePairing(pairing), nil
}
func (s *PairingService) UpdatePairing(ctx context.Context, req *debate_management.UpdatePairingRequest) (*debate_management.Pairing, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	err = queries.UpdatePairing(ctx, models.UpdatePairingParams{
		Debateid: req.GetPairing().GetPairingId(),
		Team1id:  req.GetPairing().GetTeam1().GetTeamId(),
		Team2id:  req.GetPairing().GetTeam2().GetTeamId(),
		Roomid:   req.GetPairing().GetRoomId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update pairing: %v", err)
	}

	// Fetch the updated pairing
	updatedPairing, err := queries.GetPairingByID(ctx, req.GetPairing().GetPairingId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated pairing: %v", err)
	}

	return convertSinglePairing(updatedPairing), nil
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

    // Get teams for the tournament
    teams, err := queries.GetTeamsByTournament(ctx, req.GetTournamentId())
    if err != nil {
        return nil, fmt.Errorf("failed to get teams: %v", err)
    }

    // Get previous pairings
    prevPairings, err := queries.GetPreviousPairings(ctx, models.GetPreviousPairingsParams{
        Tournamentid: req.GetTournamentId(),
        Roundnumber:  req.GetRoundNumber(),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get previous pairings: %v", err)
    }

    // Convert teams and previous pairings to the format expected by the pairing algorithm
algorithmTeams := make([]*pairing_algorithm.Team, len(teams))

for i, team := range teams {
    algorithmTeams[i] = &pairing_algorithm.Team{
        ID:        int(team.Teamid),
        Name:      team.Name,
        Opponents: make(map[int]bool),
    }
}

    // Use prevPairings to set up conflicts in algorithmTeams
    for _, prevPairing := range prevPairings {
        team1Index := findTeamIndex(algorithmTeams, int(prevPairing.Team1id))
        team2Index := findTeamIndex(algorithmTeams, int(prevPairing.Team2id))
        if team1Index != -1 && team2Index != -1 {
            algorithmTeams[team1Index].Opponents[int(prevPairing.Team2id)] = true
            algorithmTeams[team2Index].Opponents[int(prevPairing.Team1id)] = true
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

    // Create rooms
    numDebates := len(teams) / 2
    rooms := make([]int, numDebates)
    for i := 0; i < numDebates; i++ {
        room, err := queries.CreateRoom(ctx, models.CreateRoomParams{
            Roomname:  fmt.Sprintf("Room %d", i+1),
            Location:  "TBD",
            Capacity:  int32(tournament.Judgesperdebatepreliminary + 2), // Judges + 2 teams
        })
        if err != nil {
            return nil, fmt.Errorf("failed to create room: %v", err)
        }
        rooms[i] = int(room.Roomid)
    }

    // Generate pairings using the pairing algorithm
    specs := pairing_algorithm.TournamentSpecs{
        PreliminaryRounds: int(tournament.Numberofpreliminaryrounds),
        EliminationRounds: int(tournament.Numberofeliminationrounds),
        JudgesPerDebate:   int(tournament.Judgesperdebatepreliminary),
    }

    allPairings, err := pairing_algorithm.GeneratePairings(algorithmTeams, algorithmJudges, rooms, specs)
    if err != nil {
        return nil, fmt.Errorf("failed to generate pairings: %v", err)
    }

    // Save new pairings to the database
 dbPairings := make([]*debate_management.Pairing, 0)
    for roundNumber, roundPairings := range allPairings {
        for _, pair := range roundPairings {
            startTime := time.Now().Add(time.Duration(roundNumber) * time.Hour) // Example: each round starts 1 hour after the previous

            debate, err := queries.CreateDebate(ctx, models.CreateDebateParams{
                Tournamentid:       req.GetTournamentId(),
                Roundnumber:        int32(roundNumber + 1),
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
            for _, judge := range pair.Judges {
                err := queries.AssignJudgeToDebate(ctx, models.AssignJudgeToDebateParams{
                    Tournamentid:  req.GetTournamentId(),
                    Judgeid:       int32(judge.ID),
                    Debateid:      debate,
                    Roundnumber:   int32(roundNumber + 1),
                    Iselimination: false,
                    Isheadjudge:   false, // You might want to implement logic to determine head judge
                })
                if err != nil {
                    return nil, fmt.Errorf("failed to assign judge to debate: %v", err)
                }
            }

            dbPairings = append(dbPairings, &debate_management.Pairing{
                PairingId:          debate,
                RoundNumber:        int32(roundNumber + 1),
                IsEliminationRound: false,
                RoomId:             int32(pair.Room),
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
                Iselimination: false,
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

	// Delete existing pairings for the tournament
	err = queries.DeletePairingsForTournament(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing pairings: %v", err)
	}

	// Generate new pairings
	newPairings, err := s.GeneratePairings(ctx, &debate_management.GeneratePairingsRequest{
		TournamentId: req.GetTournamentId(),
		Token:        req.GetToken(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate new pairings: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return newPairings, nil
}

func convertPairings(dbPairings []models.GetPairingsByTournamentAndRoundRow) []*debate_management.Pairing {
	pairings := make([]*debate_management.Pairing, len(dbPairings))
	for i, dbPairing := range dbPairings {
		pairings[i] = &debate_management.Pairing{
			PairingId:          dbPairing.Debateid,
			RoundNumber:        dbPairing.Roundnumber,
			IsEliminationRound: dbPairing.Iseliminationround,
			RoomId:             dbPairing.Roomid,
			Team1: &debate_management.Team{
				TeamId: dbPairing.Team1id,
				Name:   dbPairing.Team1name,
			},
			Team2: &debate_management.Team{
				TeamId: dbPairing.Team2id,
				Name:   dbPairing.Team2name,
			},
		}
	}
	return pairings
}

func convertPairingFromRow(dbPairing models.GetPairingsByTournamentAndRoundRow) *debate_management.Pairing {
	return &debate_management.Pairing{
		PairingId:          dbPairing.Debateid,
		RoundNumber:        dbPairing.Roundnumber,
		IsEliminationRound: dbPairing.Iseliminationround,
		RoomId:             dbPairing.Roomid,
		Team1: &debate_management.Team{
			TeamId: dbPairing.Team1id,
			Name:   dbPairing.Team1name,
		},
		Team2: &debate_management.Team{
			TeamId: dbPairing.Team2id,
			Name:   dbPairing.Team2name,
		},
	}
}

func convertSinglePairing(dbPairing models.GetPairingByIDRow) *debate_management.Pairing {
	return &debate_management.Pairing{
		PairingId:          dbPairing.Debateid,
		RoundNumber:        dbPairing.Roundnumber,
		IsEliminationRound: dbPairing.Iseliminationround,
		RoomId:             dbPairing.Roomid,
		Team1: &debate_management.Team{
			TeamId: dbPairing.Team1id,
			Name:   dbPairing.Team1name,
		},
		Team2: &debate_management.Team{
			TeamId: dbPairing.Team2id,
			Name:   dbPairing.Team2name,
		},
	}
}

func getTeamName(teams []pairing_algorithm.Team, id int) string {
	for _, team := range teams {
		if team.ID == id {
			return team.Name
		}
	}
	return ""
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

func findTeamIndex(teams []*pairing_algorithm.Team, teamID int) int {
    for i, team := range teams {
        if team.ID == teamID {
            return i
        }
    }
    return -1
}
