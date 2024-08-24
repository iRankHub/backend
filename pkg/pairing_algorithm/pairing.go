package pairing_algorithm

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
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

const maxAttempts = 1000
var ErrUnableToPair = errors.New("unable to generate valid pairings")

func GeneratePairings(teams []*Team, judges []*Judge, rooms []int, specs TournamentSpecs) ([][]*Debate, error) {
	allPairings := make([][]*Debate, specs.PreliminaryRounds)

	for round := 0; round < specs.PreliminaryRounds; round++ {
		pairings, err := generateRoundPairings(teams, round == specs.PreliminaryRounds-1)
		if err != nil {
			return nil, err
		}

		debates := assignJudgesAndRooms(pairings, judges, rooms, specs.JudgesPerDebate)
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
	if len(teamIDs)%2 != 0 {
		return nil, errors.New("number of teams must be even")
	}

    teams := make([]*Team, len(teamIDs))
    for i, id := range teamIDs {
        teams[i] = &Team{ID: id, Opponents: make(map[int]bool)}
    }


	pairings := make([][][]int, rounds)
	totalAttempts := 0

    for round := 0; round < rounds; {
        attempts := 0
        for attempts < maxAttempts {
            roundPairings, err := generateRoundPairings(teams, round == rounds-1)
            if err == nil {
                pairingInts := make([][]int, len(roundPairings))
                for i, debate := range roundPairings {
                    pairingInts[i] = []int{debate.Team1.ID, debate.Team2.ID}
                }
                pairings[round] = pairingInts

                // Update team statistics
                for _, debate := range roundPairings {
                    debate.Team1.Opponents[debate.Team2.ID] = true
                    debate.Team2.Opponents[debate.Team1.ID] = true
                    debate.Team1.LastSide = 1
                    debate.Team2.LastSide = -1
                    debate.Team1.Wins++
                    debate.Team2.Wins++
                }
                round++
                break
            }
            attempts++
            totalAttempts++
        }

		if attempts == maxAttempts {
			// If we can't generate valid pairings, reset and try again
			for i := range teams {
				teams[i].LastSide = 0
				teams[i].Opponents = make(map[int]bool)
				teams[i].Wins = 0
			}
			round = 0
			pairings = make([][][]int, rounds)
		}

		if totalAttempts >= maxAttempts*rounds {
			return nil, fmt.Errorf("unable to generate valid pairings after %d total attempts", totalAttempts)
		}
	}

	return pairings, nil
}

func generateRoundPairings(teams []*Team, isLastRound bool) ([]*Debate, error) {
	shuffledTeams := make([]*Team, len(teams))
	copy(shuffledTeams, teams)
	rand.Shuffle(len(shuffledTeams), func(i, j int) {
		shuffledTeams[i], shuffledTeams[j] = shuffledTeams[j], shuffledTeams[i]
	})

	sort.Slice(shuffledTeams, func(i, j int) bool {
		return shuffledTeams[i].Wins > shuffledTeams[j].Wins
	})

	pairings := []*Debate{}
	paired := make(map[int]bool)

	for i := 0; i < len(shuffledTeams); i++ {
		if paired[shuffledTeams[i].ID] {
			continue
		}

		for j := i + 1; j < len(shuffledTeams); j++ {
			if paired[shuffledTeams[j].ID] {
				continue
			}

			if !shuffledTeams[i].Opponents[shuffledTeams[j].ID] &&
				(isLastRound || (shuffledTeams[i].LastSide != 1 || shuffledTeams[j].LastSide != -1)) {
				pairings = append(pairings, &Debate{
					Team1: shuffledTeams[i],
					Team2: shuffledTeams[j],
				})
				paired[shuffledTeams[i].ID] = true
				paired[shuffledTeams[j].ID] = true
				break
			}
		}
	}

	if len(pairings) < len(teams)/2 {
		return nil, ErrUnableToPair
	}

	return pairings, nil
}
