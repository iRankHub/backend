package pairing_algorithm

import (
	"errors"
	"math/rand"
)

type Team struct {
	ID         int
	Name       string
	Wins       int
	LastSide   int // 1 for affirmative, -1 for negative, 0 for no previous side
	Opponents  map[int]bool
	SpeakerIDs []int
}

type Judge struct {
	ID   int
	Name string
	IsHeadJudge bool
}

type Debate struct {
	Team1  *Team
	Team2  *Team
	Judges []*Judge
	Room   int
}

type TournamentSpecs struct {
	PreliminaryRounds int
	EliminationRounds int
	JudgesPerDebate   int
}

var ErrUnableToPair = errors.New("unable to generate valid pairings")

func GeneratePairings(teams []*Team, judges []*Judge, rooms []int, specs TournamentSpecs) ([][]*Debate, error) {
	allPairings := make([][]*Debate, specs.PreliminaryRounds)

	teamIDs := make([]int, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}

	prelimPairings, err := GeneratePreliminaryPairings(teamIDs, specs.PreliminaryRounds)
	if err != nil {
		return nil, err
	}

	for round, roundPairings := range prelimPairings {
		debates := make([]*Debate, len(roundPairings))
		for i, pair := range roundPairings {
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
		allPairings[round] = debates

		// Update team statistics
		for _, debate := range debates {
			debate.Team1.Opponents[debate.Team2.ID] = true
			debate.Team2.Opponents[debate.Team1.ID] = true
			debate.Team1.LastSide = 1
			debate.Team2.LastSide = -1
			debate.Team1.Wins++
			debate.Team2.Wins++
		}
	}

	return allPairings, nil
}

func assignJudgesAndRooms(pairings []*Debate, judges []*Judge, rooms []int, judgesPerDebate int) []*Debate {
	availableJudges := make([]*Judge, len(judges))
	copy(availableJudges, judges)
	availableRooms := make([]int, len(rooms))
	copy(availableRooms, rooms)

	for _, debate := range pairings {
		// Assign judges
		debate.Judges = make([]*Judge, judgesPerDebate)
		for i := 0; i < judgesPerDebate; i++ {
			if len(availableJudges) > 0 {
				judgeIndex := rand.Intn(len(availableJudges))
				debate.Judges[i] = availableJudges[judgeIndex]
				availableJudges = append(availableJudges[:judgeIndex], availableJudges[judgeIndex+1:]...)
			}
		}

		// Assign head judge
		if len(debate.Judges) > 0 {
			headJudgeIndex := 0
			if len(debate.Judges) > 1 {
				headJudgeIndex = rand.Intn(len(debate.Judges))
			}
			debate.Judges[headJudgeIndex].IsHeadJudge = true
		}

		// Assign room
		if len(availableRooms) > 0 {
			roomIndex := rand.Intn(len(availableRooms))
			debate.Room = availableRooms[roomIndex]
			availableRooms = append(availableRooms[:roomIndex], availableRooms[roomIndex+1:]...)
		}
	}

	return pairings
}

func GeneratePreliminaryPairings(teamIDs []int, rounds int) ([][][]int, error) {
	originalLength := len(teamIDs)
	if originalLength%2 != 0 {
		teamIDs = append(teamIDs, -1) // Add a "Public Speaking" team with ID -1
	}

	n := len(teamIDs) / 2
	proposition := make([]int, n)
	opposition := make([]int, n)
	copy(proposition, teamIDs[:n])
	copy(opposition, teamIDs[n:])

	pairings := make([][][]int, rounds)

	for round := 0; round < rounds; round++ {
		roundPairings := make([][]int, n)

		if round == rounds-1 { // Last round
			// Switch sides first
			proposition, opposition = opposition, proposition

			// Combine arrays
			combined := append(proposition, opposition...)

			// Create new proposition and opposition
			newProp := make([]int, 0, n)
			newOpp := make([]int, 0, n)
			for i := 0; i < len(combined); i++ {
				if i%2 == 0 {
					newProp = append(newProp, combined[i])
				} else {
					newOpp = append(newOpp, combined[i])
				}
			}

			// Pair across
			for i := 0; i < n; i++ {
				if i < len(newProp) && i < len(newOpp) {
					roundPairings[i] = []int{newProp[i], newOpp[i]}
				} else if i < len(newProp) {
					roundPairings[i] = []int{newProp[i], -1} // Handle odd number of teams
				}
			}
		} else {
			switch round % 4 {
			case 0: // Pair across
				for i := 0; i < n; i++ {
					roundPairings[i] = []int{proposition[i], opposition[i]}
				}
			case 1: // Pair diagonal first-last
				for i := 0; i < n; i++ {
					roundPairings[i] = []int{opposition[i], proposition[n-1-i]}
				}
			case 2: // Pair diagonal two by two
				for i := 0; i < n; i += 2 {
					if i+1 < n {
						roundPairings[i] = []int{proposition[i], opposition[i+1]}
						roundPairings[i+1] = []int{proposition[i+1], opposition[i]}
					} else {
						roundPairings[i] = []int{proposition[i], -1} // Handle odd number of teams
					}
				}
			case 3: // Pair diagonal first and second last
				for i := 0; i < n; i++ {
					roundPairings[i] = []int{opposition[i], proposition[(i+n/2)%n]}
				}
			}

			// Switch sides for next round, except after the first round
			if round != 0 {
				proposition, opposition = opposition, proposition
			}
		}

		pairings[round] = roundPairings
	}

	return pairings, nil
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