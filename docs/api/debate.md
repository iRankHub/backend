# Debate Management API

## Pairing Management

### GeneratePairings

Endpoint: `DebateService.GeneratePairings`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
  "round_number": 1,
  "is_elimination_round": false,
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
  "is_elimination_round": false,
  "token": "your_auth_token_here"
}
```

### UpdatePairing

Endpoint: `DebateService.UpdatePairing`
Authorization: Admin only

Request:
```json
{
  "pairing_id": 1,
  "team1_id": 2,
  "team2_id": 3,
  "room_id": 4,
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
  "is_elimination_round": false,
  "token": "your_auth_token_here"
}
```

### UpdateRoom

Endpoint: `DebateService.UpdateRoom`
Authorization: Admin only

Request:
```json
{
  "room_id": 1,
  "room_name": "Updated Room Name",
  "token": "your_auth_token_here"
}
```

### AssignRoomsToDebates

Endpoint: `DebateService.AssignRoomsToDebates`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
  "round_number": 1,
  "is_elimination_round": false,
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
  "is_elimination_round": false,
  "token": "your_auth_token_here"
}
```

### AssignJudgeToDebate

Endpoint: `DebateService.AssignJudgeToDebate`
Authorization: Admin only

Request:
```json
{
  "tournament_id": 1,
  "judge_id": 2,
  "debate_id": 3,
  "round_number": 1,
  "is_elimination_round": false,
  "is_head_judge": true,
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
  "is_elimination_round": false,
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
  "ballot_id": 1,
  "team1_total_score": 75.5,
  "team2_total_score": 78.0,
  "recording_status": "completed",
  "verdict": "team2_win",
  "speakers": [
    {
      "score_id": 1,
      "rank": 1,
      "points": 37.5,
      "feedback": "Excellent argumentation"
    },
    {
      "score_id": 2,
      "rank": 2,
      "points": 36.0,
      "feedback": "Strong rebuttal"
    }
  ],
  "token": "your_auth_token_here"
}
```

## Testing Debate Management Features

To test the debate management features:

1. Use the `Login` endpoint to authenticate as an admin and receive a token.
2. Include the token in the request body for subsequent authenticated requests.
3. Test the following scenarios:

   a. Pairing Generation and Management:
   - Use `GeneratePairings` to create pairings for a tournament round.
   - Use `GetPairings` to retrieve the generated pairings.
   - Use `UpdatePairing` to modify a specific pairing.

   b. Room Management:
   - Use `GetRooms` to retrieve available rooms for a tournament round.
   - Use `UpdateRoom` to modify room details.
   - Use `AssignRoomsToDebates` to automatically assign rooms to debates.

   c. Judge Management:
   - Use `GetJudges` to retrieve available judges for a tournament round.
   - Use `AssignJudgeToDebate` to assign a judge to a specific debate.

   d. Ballot Management:
   - Use `GetBallots` to retrieve all ballots for a tournament round.
   - Use `GetBallot` to retrieve a specific ballot.
   - Use `UpdateBallot` to submit or modify ballot results and speaker scores.

4. For each test, verify that the appropriate actions are taken:
   - Pairings are generated correctly and can be modified.
   - Rooms are assigned to debates efficiently.
   - Judges are assigned to debates according to tournament rules.
   - Ballots can be submitted and updated with accurate scoring.

5. Test error scenarios:
   - Attempt to generate pairings for a non-existent tournament or round.
   - Try to assign more judges than available to a debate.
   - Attempt to submit a ballot with invalid scores or for a non-existent debate.

Remember to include the authentication token in the request body for each request:
- Key: `token`
- Value: `<token_received_from_login>`

This testing process will help ensure that all aspects of debate management within a tournament are functioning correctly.