package pairing_algorithm

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

)

type Team struct {
	ID        int
	Name	  string
	Wins      int
	LastSide  int // 1 for affirmative, -1 for negative, 0 for no previous side
	Opponents map[int]bool
}

type Debate struct {
	Team1 int
	Team2 int
	Room  int
	Judge int
}

const maxAttempts = 1000

func AssignRoomsAndJudges(pairings [][][]int, rooms, judges []int) ([][]*Debate, error) {
	if len(rooms) < len(pairings[0]) || len(judges) < len(pairings[0]) {
		return nil, errors.New("not enough rooms or judges for all debates")
	}

	debates := make([][]*Debate, len(pairings))

	for round, roundPairings := range pairings {
		debates[round] = make([]*Debate, len(roundPairings))
		availableRooms := make([]int, len(rooms))
		copy(availableRooms, rooms)
		availableJudges := make([]int, len(judges))
		copy(availableJudges, judges)

		for i, pairing := range roundPairings {
			roomIndex := rand.Intn(len(availableRooms))
			judgeIndex := rand.Intn(len(availableJudges))

			debates[round][i] = &Debate{
				Team1: pairing[0],
				Team2: pairing[1],
				Room:  availableRooms[roomIndex],
				Judge: availableJudges[judgeIndex],
			}

			// Remove assigned room and judge from available lists
			availableRooms = append(availableRooms[:roomIndex], availableRooms[roomIndex+1:]...)
			availableJudges = append(availableJudges[:judgeIndex], availableJudges[judgeIndex+1:]...)
		}
	}

	return debates, nil
}

func GeneratePreliminaryPairings(teamIDs []int, rounds int) ([][][]int, error) {
	if len(teamIDs)%2 != 0 {
		return nil, errors.New("number of teams must be even")
	}

	teams := make([]Team, len(teamIDs))
	for i, id := range teamIDs {
		teams[i] = Team{ID: id, Opponents: make(map[int]bool)}
	}

	pairings := make([][][]int, rounds)
	totalAttempts := 0

	for round := 0; round < rounds; {
		attempts := 0
		for attempts < maxAttempts {
			roundPairings, err := generateRoundPairings(teams, round == rounds-1)
			if err == nil {
				pairings[round] = roundPairings

				// Update team statistics
				for _, pair := range roundPairings {
					team1 := findTeam(&teams, pair[0])
					team2 := findTeam(&teams, pair[1])

					if team1 == nil || team2 == nil {
						return nil, errors.New("team not found")
					}

					team1.Opponents[team2.ID] = true
					team2.Opponents[team1.ID] = true
					team1.LastSide = 1
					team2.LastSide = -1
					team1.Wins++
					team2.Wins++
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

func generateRoundPairings(teams []Team, isLastRound bool) ([][]int, error) {
	shuffledTeams := make([]Team, len(teams))
	copy(shuffledTeams, teams)
	rand.Shuffle(len(shuffledTeams), func(i, j int) {
		shuffledTeams[i], shuffledTeams[j] = shuffledTeams[j], shuffledTeams[i]
	})

	sort.Slice(shuffledTeams, func(i, j int) bool {
		return shuffledTeams[i].Wins > shuffledTeams[j].Wins
	})

	pairings := [][]int{}
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
				pairings = append(pairings, []int{shuffledTeams[i].ID, shuffledTeams[j].ID})
				paired[shuffledTeams[i].ID] = true
				paired[shuffledTeams[j].ID] = true
				break
			}
		}
	}

	if len(pairings) < len(teams)/2 {
		return nil, errors.New("unable to generate valid pairings")
	}

	return pairings, nil
}

func findTeam(teams *[]Team, id int) *Team {
	for i := range *teams {
		if (*teams)[i].ID == id {
			return &(*teams)[i]
		}
	}
	return nil
}