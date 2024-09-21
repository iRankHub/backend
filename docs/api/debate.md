# Debate Management API

## Pairing Management

### GeneratePreliminaryPairings

Endpoint: `DebateService.GeneratePreliminaryPairings`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### GenerateEliminationPairings

Endpoint: `DebateService.GenerateEliminationPairings`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
  "round_number": 1,
  "token": "your_auth_token_here"
}
```

### GetPairings

Endpoint: `DebateService.GetPairings`

Request:
```json
{
  "tournament_id": 1,
  "round_number": 1,
  "is_elimination": false,
  "token": "your_auth_token_here"
}
```

### UpdatePairings

Endpoint: `DebateService.UpdatePairings`
Authorization: Admin only

Request:
```json
{
  "pairings": [
    {
      "judges": [],
      "pairing_id": 70,
      "round_number": 1,
      "is_elimination_round": false,
      "room_id": 66,
      "room_name": "Room 5",
      "team1": {
        "speakers": [],
        "speaker_names": [],
        "team_id": 5,
        "name": "Team A",
        "total_points": 0,
        "league_name": "",
        "feedback": ""
      },
      "team2": {
        "speakers": [],
        "speaker_names": [],
        "team_id": 5,
        "name": "Team 1",
        "total_points": 0,
        "league_name": "",
        "feedback": ""
      },
      "head_judge_name": "Al John"
    }
    // other pairings...
  ],
  "token": "your_auth_token_here"
}
```

## Room Management

### GetRooms

Endpoint: `DebateService.GetRooms`

Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### GetRoom

Endpoint: `DebateService.GetRoom`

Request:
```json
{
  "room_id": 1,
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### UpdateRoom

Endpoint: `DebateService.UpdateRoom`
Authorization: Admin only

Request:
```json
{
  "room": {
    "room_id": 1,
    "room_name": "Updated Room Name"
  },
  "token": "your_auth_token_here"
}
```

## Judge Management

### GetJudges

Endpoint: `DebateService.GetJudges`

Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### GetJudge

Endpoint: `DebateService.GetJudge`

Request:
```json
{
  "judge_id": 1,
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### UpdateJudge
Endpoint: `DebateService.UpdateJudge`
Authorization: Admin only

Request:
```json
{
  "judge_id": 1,
  "tournament_id": 1,
   "preliminary": {
        "1": {
            "room_id": 64
        },
        "2": {
            "room_id": 66
        }
    },
       "elimination": {
        "1": {
            "room_id": 64
        },
        "2": {
            "room_id": 66
        }
    },
  "token": "your_auth_token_here"
}
```

## Ballot Management

### GetBallots

Endpoint: `DebateService.GetBallots`

Request:
```json
{
  "tournament_id": 1,
  "round_number": 1,
  "is_elimination": false,
  "token": "your_auth_token_here"
}
```

### GetBallot

Endpoint: `DebateService.GetBallot`

Request:
```json
{
  "ballot_id": 1,
  "token": "your_auth_token_here"
}
```

### GetBallotByJudgeID

Endpoint: `DebateService.GetBallotByJudgeID`

Request:
```json
{
  "judge_id": 1,
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### UpdateBallot

Endpoint: `DebateService.UpdateBallot`
Authorization: Head Judge or Admin

Request:
```json
{
  "ballot": {
    "ballot_id": 1,
    "team1": {
      "team_id": 1,
      "total_points": 75.5,
      "feedback": "Excellent performance overall",
      "speakers": [
        {
          "score_id": 1,
          "speaker_id": 101,
          "rank": 1,
          "points": 37.5,
          "feedback": "Strong argumentation and rebuttal"
        },
        {
          "score_id": 2,
          "speaker_id": 102,
          "rank": 2,
          "points": 38.0,
          "feedback": "Clear presentation and good teamwork"
        }
      ]
    },
    "team2": {
      "team_id": 2,
      "total_points": 78.0,
      "feedback": "Very persuasive arguments and excellent teamwork",
      "speakers": [
        {
          "score_id": 3,
          "speaker_id": 201,
          "rank": 1,
          "points": 39.0,
          "feedback": "Outstanding speaker, very convincing"
        },
        {
          "score_id": 4,
          "speaker_id": 202,
          "rank": 2,
          "points": 39.0,
          "feedback": "Excellent rebuttal and time management"
        }
      ]
    },
    "verdict": "Team B"
  },
  "token": "your_auth_token_here"
}
```

Notes for UpdateBallot:
- Only the head judge assigned to the debate or an admin can update the ballot.
- The `recording_status` will automatically be set to "Recorded" upon successful update.
- `last_updated_by` will be set to the user ID of the person making the update (derived from the token).
- `last_updated_at` will be automatically set to the current timestamp.
- If the update is made by the head judge, `head_judge_submitted` will be set to true.
- Once a head judge has submitted a ballot (i.e., `head_judge_submitted` is true), further updates are not allowed unless made by an admin.
- Ensure that all speaker scores and team total points are included in the update request.
- The `verdict` should be one of: "team1 name", "team2 name".
- Any existing speaker scores will be overwritten by the new scores provided in the update request.

## Ranking Management

### GetTournamentStudentRanking

Endpoint: `DebateService.GetTournamentStudentRanking`

Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### GetOverallStudentRanking

Endpoint: `DebateService.GetOverallStudentRanking`

Request:
```json
{
  "user_id": 1,
  "token": "your_auth_token_here"
}
```

### GetStudentOverallPerformance

Endpoint: `DebateService.GetStudentOverallPerformance`

Request:
```json
{
  "user_id": 1,
  "start_date": "2023-01-01",
  "end_date": "2023-12-31",
  "token": "your_auth_token_here"
}
```

### GetTournamentTeamsRanking

Endpoint: `DebateService.GetTournamentTeamsRanking`
Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### GetTournamentSchoolRanking

Endpoint: `DebateService.GetTournamentSchoolRanking`

Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### GetOverallSchoolRanking

Endpoint: `DebateService.GetOverallSchoolRanking`

Request:
```json
{
  "user_id": 1,
  "token": "your_auth_token_here"
}
```

### GetSchoolOverallPerformance

Endpoint: `DebateService.GetSchoolOverallPerformance`

Request:
```json
{
  "user_id": 1,
  "start_date": "2023-01-01",
  "end_date": "2023-12-31",
  "token": "your_auth_token_here"
}
```

Notes for Ranking Endpoints:
- All ranking endpoints require authentication.
- Dates in requests and responses should be in ISO 8601 format.
- The `rank_change` field in responses indicates improvement (positive value) or decline (negative value) in ranking.
- For `GetOverallStudentRanking`, the response includes the top 3 students' information along with the requested student's ranking.
- `GetStudentOverallPerformance` allows for querying performance data within a specific date range.

## Team Management

### CreateTeam

Endpoint: `DebateService.CreateTeam`
Authorization: Admin only

Request:
```json
{
  "name": "New Team",
  "tournament_id": 1,
  "speakers": [
    {
      "speaker_id": 1
    },
    {
      "speaker_id": 2
    }
  ],
  "token": "your_auth_token_here"
}
```

### GetTeam

Endpoint: `DebateService.GetTeam`

Request:
```json
{
  "team_id": 1,
  "token": "your_auth_token_here"
}
```

### UpdateTeam

Endpoint: `DebateService.UpdateTeam`
Authorization: Admin only

Request:
```json
{
  "team": {
    "team_id": 1,
    "name": "Updated Team Name",
    "speakers": [
      {
        "speaker_id": 1
      },
      {
        "speaker_id": 3
      }
    ]
  },
  "token": "your_auth_token_here"
}
```

### GetTeamsByTournament

Endpoint: `DebateService.GetTeamsByTournament`

Request:
```json
{
  "tournament_id": 1,
  "token": "your_auth_token_here"
}
```

### Delete Team

Endpoint: `DebateService.DeleteTeam`

Request:
```json
{
  "team_id": 1,
  "token": "your_auth_token_here"
}
```

## Testing Debate Management Features

To test the debate management features:

1. Use the `Login` endpoint to authenticate as an admin and receive a token.
2. Include the token in the request body for subsequent authenticated requests.
3. Test the following scenarios:

   a. Team Management:
   - Use `CreateTeam` to create a new team for a tournament.
   - Use `GetTeam` to retrieve team details.
   - Use `UpdateTeam` to modify team information.
   - Use `GetTeamsByTournament` to list all teams in a tournament.
   - Use `DeleteTeam` to delete the team in that tournament

   b. Pairing Generation and Management:
   - Use `GeneratePairings` to create pairings for a tournament.
   - Use `GetPairings` to retrieve the generated pairings.
   - Use `GetPairing` to retrieve a specific pairing.
   - Use `UpdatePairing` to modify a specific pairing.
   - Use `RegeneratePairings` to recreate pairings for a tournament.

   c. Room Management:
   - Use `GetRooms` to retrieve available rooms for a tournament round.
   - Use `GetRoom` to retrieve details of a specific room.
   - Use `UpdateRoom` to modify room details.

   d. Judge Management:
   - Use `GetJudges` to retrieve available judges for a tournament round.
   - Use `GetJudge` to retrieve details of a specific judge.
   - Use `UpdateJudge` to automatically update the room which the judge is assigned.

   e. Ballot Management:
   - Use `GetBallots` to retrieve all ballots for a tournament round.
   - Use `GetBallot` to retrieve a specific ballot.
   - Use `UpdateBallot` to submit or modify ballot results and speaker scores.

4. For each test, verify that the appropriate actions are taken and data is correctly managed.

5. Test error scenarios and edge cases for each endpoint.

Remember to include the authentication token in the request body for each request:
- Key: `token`
- Value: `<token_received_from_login>`

This testing process will help ensure that all aspects of debate management within a tournament are functioning correctly.