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
		{"40 teams, 3 rounds", 40, 3},
		{"40 teams, 4 rounds", 40, 4},
		{"40 teams, 5 rounds", 40, 5},
		{"50 teams, 3 rounds", 50, 3},
		{"50 teams, 4 rounds", 50, 4},
		{"50 teams, 5 rounds", 50, 5},
		{"60 teams, 3 rounds", 60, 3},
		{"60 teams, 4 rounds", 60, 4},
		{"60 teams, 5 rounds", 60, 5},
		{"70 teams, 3 rounds", 70, 3},
		{"70 teams, 4 rounds", 70, 4},
		{"70 teams, 5 rounds", 70, 5},
		{"80 teams, 3 rounds", 80, 3},
		{"80 teams, 4 rounds", 80, 4},
		{"80 teams, 5 rounds", 80, 5},
		{"90 teams, 3 rounds", 90, 3},
		{"90 teams, 4 rounds", 90, 4},
		{"90 teams, 5 rounds", 90, 5},
		{"100 teams, 3 rounds", 100, 3},
		{"100 teams, 4 rounds", 100, 4},
		{"100 teams, 5 rounds", 100, 5},
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
			expectedPairingsPerRound := (tc.teams + 1) / 2 // Account for odd number of teams
			for i, round := range pairings {
				if len(round) != expectedPairingsPerRound {
					t.Errorf("Round %d: Expected %d pairings, got %d", i+1, expectedPairingsPerRound, len(round))
				}
			}

			// Check that no team meets twice
			meetingCount := make(map[int]map[int]bool)
			for roundIndex, round := range pairings {
				for pairingIndex, pairing := range round {
					if pairing[0] == -1 || pairing[1] == -1 {
						continue // Skip checking for byes
					}
					if meetingCount[pairing[0]] == nil {
						meetingCount[pairing[0]] = make(map[int]bool)
					}
					if meetingCount[pairing[1]] == nil {
						meetingCount[pairing[1]] = make(map[int]bool)
					}
					if meetingCount[pairing[0]][pairing[1]] || meetingCount[pairing[1]][pairing[0]] {
						t.Errorf("Teams %d and %d meet more than once (Round %d, Pairing %d)", pairing[0], pairing[1], roundIndex+1, pairingIndex+1)
					}
					meetingCount[pairing[0]][pairing[1]] = true
					meetingCount[pairing[1]][pairing[0]] = true
				}
			}

			// Check that each team participates in each round (except for bye in odd-numbered tournaments)
			for roundIndex, round := range pairings {
				participatingTeams := make(map[int]bool)
				for _, pairing := range round {
					participatingTeams[pairing[0]] = true
					if pairing[1] != -1 {
						participatingTeams[pairing[1]] = true
					}
				}
				if len(participatingTeams) != tc.teams {
					t.Errorf("Round %d: Not all teams are participating. Expected %d, got %d", roundIndex+1, tc.teams, len(participatingTeams))
				}
			}
		})
	}
}


func TestAssignRoomsAndJudges(t *testing.T) {
    testCases := []struct {
        name     string
        pairings [][]*Debate
        rooms    []int
        judges   []*Judge
    }{
        {
            name: "Single round",
            pairings: [][]*Debate{
                {
                    {Team1: &Team{ID: 1}, Team2: &Team{ID: 2}},
                    {Team1: &Team{ID: 3}, Team2: &Team{ID: 4}},
                    {Team1: &Team{ID: 5}, Team2: &Team{ID: 6}},
                },
            },
            rooms:  []int{101, 102, 103},
            judges: []*Judge{{ID: 1}, {ID: 2}, {ID: 3}},
        },
        {
            name: "Multiple rounds",
            pairings: [][]*Debate{
                {
                    {Team1: &Team{ID: 1}, Team2: &Team{ID: 2}},
                    {Team1: &Team{ID: 3}, Team2: &Team{ID: 4}},
                    {Team1: &Team{ID: 5}, Team2: &Team{ID: 6}},
                },
                {
                    {Team1: &Team{ID: 1}, Team2: &Team{ID: 3}},
                    {Team1: &Team{ID: 2}, Team2: &Team{ID: 5}},
                    {Team1: &Team{ID: 4}, Team2: &Team{ID: 6}},
                },
                {
                    {Team1: &Team{ID: 1}, Team2: &Team{ID: 4}},
                    {Team1: &Team{ID: 2}, Team2: &Team{ID: 6}},
                    {Team1: &Team{ID: 3}, Team2: &Team{ID: 5}},
                },
            },
            rooms:  []int{101, 102, 103},
            judges: []*Judge{{ID: 1}, {ID: 2}, {ID: 3}},
        },
        {
            name: "Larger tournament",
            pairings: [][]*Debate{
                {
                    {Team1: &Team{ID: 1}, Team2: &Team{ID: 2}},
                    {Team1: &Team{ID: 3}, Team2: &Team{ID: 4}},
                    {Team1: &Team{ID: 5}, Team2: &Team{ID: 6}},
                    {Team1: &Team{ID: 7}, Team2: &Team{ID: 8}},
                    {Team1: &Team{ID: 9}, Team2: &Team{ID: 10}},
                },
                {
                    {Team1: &Team{ID: 1}, Team2: &Team{ID: 3}},
                    {Team1: &Team{ID: 2}, Team2: &Team{ID: 5}},
                    {Team1: &Team{ID: 4}, Team2: &Team{ID: 7}},
                    {Team1: &Team{ID: 6}, Team2: &Team{ID: 9}},
                    {Team1: &Team{ID: 8}, Team2: &Team{ID: 10}},
                },
            },
            rooms:  []int{101, 102, 103, 104, 105},
            judges: []*Judge{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            for _, roundPairings := range tc.pairings {
                assignedDebates := assignJudgesAndRooms(roundPairings, tc.judges, tc.rooms, 1)

                if len(assignedDebates) != len(roundPairings) {
                    t.Errorf("Expected %d debates, got %d", len(roundPairings), len(assignedDebates))
                }

                usedRooms := make(map[int]bool)
                usedJudges := make(map[int]bool)

                for _, debate := range assignedDebates {
                    if debate.Room == 0 {
                        t.Errorf("Debate not assigned a room")
                    }
                    if usedRooms[debate.Room] {
                        t.Errorf("Room %d assigned more than once", debate.Room)
                    }
                    usedRooms[debate.Room] = true

                    if len(debate.Judges) == 0 {
                        t.Errorf("Debate not assigned a judge")
                    }
                    for _, judge := range debate.Judges {
                        if usedJudges[judge.ID] {
                            t.Errorf("Judge %d assigned more than once", judge.ID)
                        }
                        usedJudges[judge.ID] = true
                    }
                }
                // Check if exactly one head judge is assigned per debate
                for _, debate := range assignedDebates {
                    headJudgeCount := 0
                    for _, judge := range debate.Judges {
                        if judge.IsHeadJudge {
                            headJudgeCount++
                        }
                    }
                    if headJudgeCount != 1 {
                        t.Errorf("Debate does not have exactly one head judge: %d", headJudgeCount)
                    }
                }
            }
        })
    }
}
