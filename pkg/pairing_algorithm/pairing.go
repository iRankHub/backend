package pairing_algorithm

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Team struct {
	ID              int
	Name            string
	Wins            int
	TotalPoints     float64
	AverageRank     float64
	LastSide        int // 1 for affirmative, -1 for negative, 0 for no previous side
	Opponents       map[int]bool
	SpeakerIDs      []int
	EliminationRank int
}

type Judge struct {
	ID          int
	Name        string
	IsHeadJudge bool
}

type Debate struct {
	Team1  *Team
	Team2  *Team
	Judges []*Judge
	Room   int
}

type TournamentSpecs struct {
	PreliminaryRounds     int
	EliminationRounds     int
	JudgesPerDebate       int
	TeamsAdvancingToElims int
}

var ErrUnableToPair = errors.New("unable to generate valid pairings")

func GeneratePairings(teams []*Team, judges []*Judge, rooms []int, specs TournamentSpecs, roundNumber int, isElimination bool) ([]*Debate, error) {
	if isElimination {
		return generateEliminationPairings(teams, judges, rooms, specs, roundNumber)
	}
	return generatePreliminaryPairings(teams, judges, rooms, specs)
}


func generateEliminationPairings(teams []*Team, judges []*Judge, rooms []int, specs TournamentSpecs, roundNumber int) ([]*Debate, error) {
	teamsNeeded := int(math.Pow(2, float64(specs.EliminationRounds-roundNumber+1)))
	if len(teams) < teamsNeeded {
		return nil, fmt.Errorf("not enough teams for elimination round %d: have %d, need %d", roundNumber, len(teams), teamsNeeded)
	}

	if roundNumber == 1 {
		// First elimination round: sort based on preliminary performance
		sort.Slice(teams, func(i, j int) bool {
			if teams[i].Wins != teams[j].Wins {
				return teams[i].Wins > teams[j].Wins
			}
			if teams[i].TotalPoints != teams[j].TotalPoints {
				return teams[i].TotalPoints > teams[j].TotalPoints
			}
			return teams[i].AverageRank < teams[j].AverageRank
		})
	}

	// Select top teams for this elimination round
	selectedTeams := teams[:teamsNeeded]

	// Create pairings
	debates := make([]*Debate, teamsNeeded/2)
	for i := 0; i < teamsNeeded/2; i++ {
		debates[i] = &Debate{
			Team1: selectedTeams[i],
			Team2: selectedTeams[teamsNeeded-1-i],
		}
		// Set elimination ranks
		debates[i].Team1.EliminationRank = i + 1
		debates[i].Team2.EliminationRank = teamsNeeded - i
	}

	// Assign rooms and judges
	debates = assignJudgesAndRooms(debates, judges, rooms, specs.JudgesPerDebate)

	return debates, nil
}


func GeneratePreliminaryPairingIDs(teamIDs []int, rounds int) ([][][]int, error) {
	originalLength := len(teamIDs)
	if originalLength%2 != 0 {
		teamIDs = append(teamIDs, -1) // Add a "Public Speaking" team with ID -1
	}

	n := len(teamIDs)
	pairings := make([][][]int, rounds)

	for round := 0; round < rounds; round++ {
		roundPairings := make([][]int, n/2)

		for i := 0; i < n/2; i++ {
			team1 := teamIDs[i]
			team2 := teamIDs[n-1-i]
			if team1 == -1 || team2 == -1 {
				// Handle bye
				if team1 == -1 {
					roundPairings[i] = []int{team2, -1}
				} else {
					roundPairings[i] = []int{team1, -1}
				}
			} else {
				roundPairings[i] = []int{team1, team2}
			}
		}

		pairings[round] = roundPairings

		// Rotate teams for the next round
		if round < rounds-1 {
			teamIDs = append(teamIDs[1:n-1], teamIDs[0], teamIDs[n-1])
		}
	}

	return pairings, nil
}

func assignJudgesAndRooms(debates []*Debate, judges []*Judge, rooms []int, judgesPerDebate int) []*Debate {
	availableJudges := make([]*Judge, len(judges))
	copy(availableJudges, judges)
	availableRooms := make([]int, len(rooms))
	copy(availableRooms, rooms)

	for _, debate := range debates {
		// Assign judges
		debate.Judges = make([]*Judge, 0, judgesPerDebate)
		for i := 0; i < judgesPerDebate; i++ {
			if len(availableJudges) > 0 {
				judgeIndex := rand.Intn(len(availableJudges))
				debate.Judges = append(debate.Judges, availableJudges[judgeIndex])
				availableJudges = append(availableJudges[:judgeIndex], availableJudges[judgeIndex+1:]...)
			} else {
				// If we run out of judges, break the loop
				break
			}
		}

		// Assign head judge
		if len(debate.Judges) > 0 {
			headJudgeIndex := 0
			debate.Judges[headJudgeIndex].IsHeadJudge = true
			for i := 1; i < len(debate.Judges); i++ {
				debate.Judges[i].IsHeadJudge = false
			}
		}

		// Assign room
		if len(availableRooms) > 0 {
			roomIndex := rand.Intn(len(availableRooms))
			debate.Room = availableRooms[roomIndex]
			availableRooms = append(availableRooms[:roomIndex], availableRooms[roomIndex+1:]...)
		} else {
			// If we run out of rooms, assign a placeholder value
			debate.Room = -1
		}
	}

	return debates
}

func generatePreliminaryPairings(teams []*Team, judges []*Judge, rooms []int, specs TournamentSpecs) ([]*Debate, error) {
	teamIDs := make([]int, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}

	prelimPairings, err := GeneratePreliminaryPairingIDs(teamIDs, specs.PreliminaryRounds)
	if err != nil {
		return nil, err
	}

	if len(judges) < specs.JudgesPerDebate {
		return nil, fmt.Errorf("not enough judges: have %d, need at least %d", len(judges), specs.JudgesPerDebate)
	}

	if len(rooms) < len(teams)/2 {
		return nil, fmt.Errorf("not enough rooms: have %d, need at least %d", len(rooms), len(teams)/2)
	}

	// We'll use the pairings for the current round (index 0)
	currentRoundPairings := prelimPairings[0]
	debates := make([]*Debate, len(currentRoundPairings))

	for i, pair := range currentRoundPairings {
		team1 := findTeamByID(teams, pair[0])
		team2 := findTeamByID(teams, pair[1])
		if team1 == nil || team2 == nil {
			return nil, errors.New("invalid team ID in pairings")
		}
		debates[i] = &Debate{
			Team1: team1,
			Team2: team2,
		}
	}

	debates = assignJudgesAndRooms(debates, judges, rooms, specs.JudgesPerDebate)

	return debates, nil
}

func findTeamByID(teams []*Team, id int) *Team {
	for _, team := range teams {
		if team.ID == id {
			return team
		}
	}
	if id == -1 {
		return &Team{ID: -1, Name: "Public Speaking"}
	}
	return nil
}