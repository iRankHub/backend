# Debate Management API

## Pairing Management

### GeneratePairings

Endpoint: `DebateService.GeneratePairings`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
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

### GetPairing

Endpoint: `DebateService.GetPairing`

Request:
```json
{
  "pairing_id": 1,
  "token": "your_auth_token_here"
}
```

### UpdatePairing

Endpoint: `DebateService.UpdatePairing`
Authorization: Admin only

Request:
```json
{
  "pairing": {
    "pairing_id": 1,
    "team1": {
      "team_id": 2
    },
    "team2": {
      "team_id": 3
    },
    "room_id": 4
  },
  "token": "your_auth_token_here"
}
```

### RegeneratePairings

Endpoint: `DebateService.RegeneratePairings`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
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
  "round_number": 1,
  "is_elimination": false,
  "token": "your_auth_token_here"
}
```

### GetRoom

Endpoint: `DebateService.GetRoom`

Request:
```json
{
  "room_id": 1,
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
  "round_number": 1,
  "is_elimination": false,
  "token": "your_auth_token_here"
}
```

### GetJudge

Endpoint: `DebateService.GetJudge`

Request:
```json
{
  "judge_id": 1,
  "token": "your_auth_token_here"
}
```

### AssignJudges

Endpoint: `DebateService.AssignJudges`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
  "round_number": 1,
  "is_elimination": false,
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

### UpdateBallot

Endpoint: `DebateService.UpdateBallot`
Authorization: Judge or Admin

Request:
```json
{
  "ballot": {
    "ballot_id": 1,
    "team1": {
      "total_points": 75.5,
      "speakers": [
        {
          "score_id": 1,
          "rank": 1,
          "points": 37.5,
          "feedback": "Excellent argumentation"
        }
      ]
    },
    "team2": {
      "total_points": 78.0,
      "speakers": [
        {
          "score_id": 2,
          "rank": 2,
          "points": 36.0,
          "feedback": "Strong rebuttal"
        }
      ]
    },
    "recording_status": "completed",
    "verdict": "team2_win"
  },
  "token": "your_auth_token_here"
}
```

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
   - Use `AssignJudges` to automatically assign judges to debates in a round.

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