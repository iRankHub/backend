package pairing_algorithm

import (
	"testing"

)

func TestPreliminaryPairings(t *testing.T) {
	testCases := []struct {
		name   string
		teams  int
		rounds int
	}{
		{"8 teams, 3 rounds", 8, 3},
		{"8 teams, 4 rounds", 8, 4},
		{"8 teams, 5 rounds", 8, 5},
		{"16 teams, 3 rounds", 16, 3},
		{"16 teams, 4 rounds", 16, 4},
		{"16 teams, 5 rounds", 16, 5},
		{"20 teams, 3 rounds", 20, 3},
		{"20 teams, 4 rounds", 20, 4},
		{"20 teams, 5 rounds", 20, 5},
		{"30 teams, 3 rounds", 30, 3},
		{"30 teams, 4 rounds", 30, 4},
		{"30 teams, 5 rounds", 30, 5},
		{"40 teams, 3 rounds", 30, 3},
		{"40 teams, 4 rounds", 30, 4},
		{"40 teams, 5 rounds", 30, 5},
		{"50 teams, 3 rounds", 30, 3},
		{"50 teams, 4 rounds", 30, 4},
		{"50 teams, 5 rounds", 30, 5},
		{"60 teams, 3 rounds", 30, 3},
		{"60 teams, 4 rounds", 30, 4},
		{"60 teams, 5 rounds", 30, 5},
		{"70 teams, 3 rounds", 30, 3},
		{"70 teams, 4 rounds", 30, 4},
		{"70 teams, 5 rounds", 30, 5},
		{"80 teams, 3 rounds", 30, 3},
		{"80 teams, 4 rounds", 30, 4},
		{"80 teams, 5 rounds", 30, 5},
		{"90 teams, 3 rounds", 30, 3},
		{"90 teams, 4 rounds", 30, 4},
		{"90 teams, 5 rounds", 30, 5},
		{"100 teams, 3 rounds", 30, 3},
		{"100 teams, 4 rounds", 30, 4},
		{"100 teams, 5 rounds", 30, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teams := make([]int, tc.teams)
			for i := range teams {
				teams[i] = i + 1
			}

			pairings, err := GeneratePreliminaryPairings(teams, tc.rounds)
			if err != nil {
				t.Fatalf("Error generating pairings: %v", err)
			}

			// Check if the number of rounds is correct
			if len(pairings) != tc.rounds {
				t.Errorf("Expected %d rounds, got %d", tc.rounds, len(pairings))
			}

			// Check if each round has the correct number of pairings
			expectedPairingsPerRound := tc.teams / 2
			for i, round := range pairings {
				if len(round) != expectedPairingsPerRound {
					t.Errorf("Round %d: Expected %d pairings, got %d", i+1, expectedPairingsPerRound, len(round))
				}
			}

			// Explicitly check that no team meets twice
			meetingCount := make(map[int]map[int]int)
			for roundIndex, round := range pairings {
				for pairingIndex, pairing := range round {
					if meetingCount[pairing[0]] == nil {
						meetingCount[pairing[0]] = make(map[int]int)
					}
					if meetingCount[pairing[1]] == nil {
						meetingCount[pairing[1]] = make(map[int]int)
					}
					meetingCount[pairing[0]][pairing[1]]++
					meetingCount[pairing[1]][pairing[0]]++
					if meetingCount[pairing[0]][pairing[1]] > 1 {
						t.Errorf("Teams %d and %d meet more than once (Round %d, Pairing %d)", pairing[0], pairing[1], roundIndex+1, pairingIndex+1)
					}
				}
			}
			t.Log("No team meets twice in the preliminaries")

			// Check if side alternation is mostly respected
			sides := make(map[int][]int)
			for roundIndex, round := range pairings {
				for _, pairing := range round {
					sides[pairing[0]] = append(sides[pairing[0]], roundIndex)
					sides[pairing[1]] = append(sides[pairing[1]], -roundIndex)
				}
			}

			sideViolations := 0
			for _, teamSides := range sides {
				for i := 1; i < len(teamSides)-1; i++ {
					if (teamSides[i] > 0 && teamSides[i-1] > 0) || (teamSides[i] < 0 && teamSides[i-1] < 0) {
						sideViolations++
					}
				}
			}

			// Allow for some flexibility in side alternation
			maxAllowedViolations := tc.teams * tc.rounds / 10 // 10% tolerance
			if sideViolations > maxAllowedViolations {
				t.Errorf("Too many side alternation violations: %d (max allowed: %d)", sideViolations, maxAllowedViolations)
			}
		})
	}
}

func TestAssignRoomsAndJudges(t *testing.T) {
	testCases := []struct {
		name     string
		pairings [][][]int
		rooms    []int
		judges   []int
	}{
		{
			name: "Single round",
			pairings: [][][]int{
				{{1, 2}, {3, 4}, {5, 6}},
			},
			rooms:  []int{101, 102, 103},
			judges: []int{1, 2, 3},
		},
		{
			name: "Multiple rounds",
			pairings: [][][]int{
				{{1, 2}, {3, 4}, {5, 6}},
				{{1, 3}, {2, 5}, {4, 6}},
				{{1, 4}, {2, 6}, {3, 5}},
			},
			rooms:  []int{101, 102, 103},
			judges: []int{1, 2, 3},
		},
		{
			name: "Larger tournament",
			pairings: [][][]int{
				{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}},
				{{1, 3}, {2, 5}, {4, 7}, {6, 9}, {8, 10}},
			},
			rooms:  []int{101, 102, 103, 104, 105},
			judges: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			debates, err := AssignRoomsAndJudges(tc.pairings, tc.rooms, tc.judges)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(debates) != len(tc.pairings) {
				t.Errorf("Expected %d rounds of debates, got %d", len(tc.pairings), len(debates))
			}

			for round, roundDebates := range debates {
				if len(roundDebates) != len(tc.pairings[round]) {
					t.Errorf("Round %d: Expected %d debates, got %d", round, len(tc.pairings[round]), len(roundDebates))
				}

				usedJudges := make(map[int]bool)

				for _, debate := range roundDebates {
					// Check if room is within the provided range
					if debate.Room < tc.rooms[0] || debate.Room > tc.rooms[len(tc.rooms)-1] {
						t.Errorf("Invalid room assigned: %d", debate.Room)
					}

					// Check if judge is within the provided range
					if debate.Judge < tc.judges[0] || debate.Judge > tc.judges[len(tc.judges)-1] {
						t.Errorf("Invalid judge assigned: %d", debate.Judge)
					}

					// Check if judge is already assigned in this round
					if usedJudges[debate.Judge] {
						t.Errorf("Judge %d assigned more than once in round %d", debate.Judge, round)
					}
					usedJudges[debate.Judge] = true

					// Check if teams are correctly assigned
					team1Found, team2Found := false, false
					for _, pairing := range tc.pairings[round] {
						if pairing[0] == debate.Team1 && pairing[1] == debate.Team2 {
							team1Found, team2Found = true, true
							break
						}
					}
					if !team1Found || !team2Found {
						t.Errorf("Incorrect team pairing in round %d: %d vs %d", round, debate.Team1, debate.Team2)
					}
				}

				// Check if all debates have been assigned a room and judge
				if len(usedJudges) != len(roundDebates) {
					t.Errorf("Round %d: Not all debates were assigned a judge", round)
				}
			}
		})
	}
}