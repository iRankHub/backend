package utils

import (
	"fmt"
	"time"
)

// ScheduleCalculator holds tournament schedule parameters
type ScheduleCalculator struct {
	StartTime         time.Time
	PreliminaryRounds int
	EliminationRounds int
	CheckInBuffer     time.Duration // 2 hours before start
	RoundDuration     time.Duration // 30 minutes for regular, 45 for finals
	SwitchBuffer      time.Duration // 10 minutes between rounds
	LunchDuration     time.Duration // 2 hours
	ImpromptuPrepTime time.Duration // 15 minutes
}

// NewScheduleCalculator creates a new schedule calculator with default timings
func NewScheduleCalculator(startTime time.Time, prelimRounds, elimRounds int) *ScheduleCalculator {
	return &ScheduleCalculator{
		StartTime:         startTime,
		PreliminaryRounds: prelimRounds,
		EliminationRounds: elimRounds,
		CheckInBuffer:     2 * time.Hour,
		RoundDuration:     30 * time.Minute,
		SwitchBuffer:      10 * time.Minute,
		LunchDuration:     2 * time.Hour,
		ImpromptuPrepTime: 15 * time.Minute,
	}
}

func (sc *ScheduleCalculator) GetCheckInTime() time.Time {
	return sc.StartTime.Add(-sc.CheckInBuffer)
}

func (sc *ScheduleCalculator) GetFirstDebateTime() time.Time {
	return sc.StartTime.Add(30 * time.Minute) // 30 minutes after official start for opening
}

func (sc *ScheduleCalculator) GetLunchTime() time.Time {
	firstDebate := sc.GetFirstDebateTime()
	roundTime := sc.RoundDuration + sc.SwitchBuffer
	return firstDebate.Add(time.Duration(sc.PreliminaryRounds) * roundTime)
}

func (sc *ScheduleCalculator) GetEliminationStartTime() time.Time {
	return sc.GetLunchTime().Add(sc.LunchDuration)
}

func (sc *ScheduleCalculator) GetFinalRoundTime() time.Time {
	elimStart := sc.GetEliminationStartTime()
	regularElimRounds := sc.EliminationRounds - 1 // exclude finals
	roundTime := sc.RoundDuration + sc.SwitchBuffer
	return elimStart.Add(time.Duration(regularElimRounds) * roundTime)
}

func (sc *ScheduleCalculator) GetAwardsCeremonyTime() time.Time {
	finalStart := sc.GetFinalRoundTime()
	return finalStart.Add(45 * time.Minute) // Finals last 45 minutes
}

// CalculateTeamFees calculates tournament fees based on team count
func CalculateTeamFees(baseAmount float64, maxTeams int) map[int]float64 {
	fees := make(map[int]float64)
	for i := 1; i <= maxTeams; i++ {
		fees[i] = baseAmount + float64(i-1)*15000
	}
	return fees
}

// FormatSchedule generates a formatted schedule string
func (sc *ScheduleCalculator) FormatSchedule() string {
	schedule := fmt.Sprintf("Check-in and Registration: %s\n", formatTime(sc.GetCheckInTime()))
	schedule += fmt.Sprintf("Opening Ceremony: %s\n", formatTime(sc.StartTime))
	schedule += fmt.Sprintf("First Debate: %s\n", formatTime(sc.GetFirstDebateTime()))
	schedule += fmt.Sprintf("Lunch Break: %s - %s\n",
		formatTime(sc.GetLunchTime()),
		formatTime(sc.GetEliminationStartTime()))
	schedule += fmt.Sprintf("Elimination Rounds: %s\n", formatTime(sc.GetEliminationStartTime()))
	schedule += fmt.Sprintf("Finals: %s\n", formatTime(sc.GetFinalRoundTime()))
	schedule += fmt.Sprintf("Awards Ceremony: %s\n", formatTime(sc.GetAwardsCeremonyTime()))
	return schedule
}

// FormatRoundSchedule formats the schedule for a specific round
func (sc *ScheduleCalculator) FormatRoundSchedule(roundNum int, isElimination bool) string {
	var roundTime time.Time
	if isElimination {
		elimStart := sc.GetEliminationStartTime()
		roundTime = elimStart.Add(time.Duration(roundNum-1) * (sc.RoundDuration + sc.SwitchBuffer))
	} else {
		firstDebate := sc.GetFirstDebateTime()
		roundTime = firstDebate.Add(time.Duration(roundNum-1) * (sc.RoundDuration + sc.SwitchBuffer))
	}

	duration := sc.RoundDuration
	if isElimination && roundNum == sc.EliminationRounds {
		duration = 45 * time.Minute // Finals are 45 minutes
	}

	return fmt.Sprintf("%s - %s",
		formatTime(roundTime),
		formatTime(roundTime.Add(duration)))
}

func formatTime(t time.Time) string {
	return t.Format("3:04 PM")
}
