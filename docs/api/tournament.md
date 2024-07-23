## League Management API

### CreateLeague

Endpoint: `TournamentService.CreateLeague`
Authorization: Admin only

Request:
```json
{
  "name": "National Debate League",
  "league_type": "international",
  "international_details": {
    "continents": ["North America"],
    "countries": ["United States of America"]
  },
  "token": "your_auth_token_here"
}
```

### GetLeague

Endpoint: `TournamentService.GetLeague`

Request:
```json
{
  "league_id": 1,
  "token": "your_auth_token_here"
}
```

### ListLeagues

Endpoint: `TournamentService.ListLeagues`

Request:
```json
{
  "page_size": 10,
  "page_token": 0,
  "token": "your_auth_token_here"
}
```

### UpdateLeague

Endpoint: `TournamentService.UpdateLeague`
Authorization: Admin only

Request:
```json
{
  "league_id": 1,
  "name": "Updated National Debate League",
  "league_type": "local",
  "local_details": {
    "provinces": ["East"],
    "districts": ["Kigali"]
  },
  "token": "your_auth_token_here"
}
```

### DeleteLeague

Endpoint: `TournamentService.DeleteLeague`
Authorization: Admin only

Request:
```json
{
  "league_id": 1,
  "token": "your_auth_token_here"
}
```

## Tournament Format Management API

### CreateTournamentFormat

Endpoint: `TournamentService.CreateTournamentFormat`
Authorization: Admin only

Request:
```json
{
  "format_name": "British Parliamentary",
  "description": "A globally recognized debate format",
  "speakers_per_team": 2,
  "token": "your_auth_token_here"
}
```

### GetTournamentFormat

Endpoint: `TournamentService.GetTournamentFormat`

Request:
```json
{
  "format_id": 1,
  "token": "your_auth_token_here"
}
```

### ListTournamentFormats

Endpoint: `TournamentService.ListTournamentFormats`

Request:
```json
{
  "page_size": 10,
  "page_token": 0,
  "token": "your_auth_token_here"
}
```

### UpdateTournamentFormat

Endpoint: `TournamentService.UpdateTournamentFormat`
Authorization: Admin only

Request:
```json
{
  "format_id": 1,
  "format_name": "Updated British Parliamentary",
  "description": "An updated globally recognized debate format",
  "speakers_per_team": 2,
  "token": "your_auth_token_here"
}
```

### DeleteTournamentFormat

Endpoint: `TournamentService.DeleteTournamentFormat`
Authorization: Admin only

Request:
```json
{
  "format_id": 1,
  "token": "your_auth_token_here"
}
```

## Tournament Management API
### CreateTournament

Endpoint: `TournamentService.CreateTournament`
Authorization: Admin only

Request:
```json
{
  "name": "Summer Debate Championship",
  "start_date": "2023-07-15 09:00",
  "end_date": "2023-07-17 18:00",
  "location": "City Convention Center",
  "format_id": 1,
  "league_id": 1,
  "number_of_preliminary_rounds": 4,
  "number_of_elimination_rounds": 2,
  "judges_per_debate_preliminary": 3,
  "judges_per_debate_elimination": 5,
  "tournament_fee": 100.00,
  "token": "your_auth_token_here"
}
```

### GetTournament

Endpoint: `TournamentService.GetTournament`

Request:
```json
{
  "tournamentId": 1,
  "token": "your_auth_token_here"
}
```

### ListTournaments

Endpoint: `TournamentService.ListTournaments`

Request:
```json
{
  "pageSize": 10,
  "pageToken": 0,
  "token": "your_auth_token_here"
}
```

### UpdateTournament

Endpoint: `TournamentService.UpdateTournament`
Authorization: Admin only

Request:
```json
{
  "tournamentId": 1,
  "name": "Updated Summer Debate Championship",
  "startDate": "2023-07-16T09:00:00Z",
  "endDate": "2023-07-18T18:00:00Z",
  "location": "Updated City Convention Center",
  "formatId": 2,
  "leagueId": 3,
  "numberOfPreliminaryRounds": 5,
  "numberOfEliminationRounds": 3,
  "judgesPerDebatePreliminary": 4,
  "judgesPerDebateElimination": 6,
  "tournamentFee": 120.00,
  "token": "your_auth_token_here"
}
```

### DeleteTournament

Endpoint: `TournamentService.DeleteTournament`
Authorization: Admin only

Request:
```json
{
  "tournamentId": 1,
  "token": "your_auth_token_here"
}
```

## Testing Tournament Management Features

To test the tournament management features, including leagues and formats:

1. Use the `Login` endpoint to authenticate as an admin and receive a token.
2. Include the token in the request body for subsequent authenticated requests.
3. Test the following scenarios:

   a. League Management:
   - Use `CreateLeague` to create a new league.
   - Use `GetLeague` to retrieve the created league details.
   - Use `ListLeagues` to get a list of leagues.
   - Use `UpdateLeague` to modify a league's details.
   - Use `DeleteLeague` to remove a league (ensure it's not associated with any tournaments first).

   b. Tournament Format Management:
   - Use `CreateTournamentFormat` to create a new tournament format.
   - Use `GetTournamentFormat` to retrieve the created format details.
   - Use `ListTournamentFormats` to get a list of formats.
   - Use `UpdateTournamentFormat` to modify a format's details.
   - Use `DeleteTournamentFormat` to remove a format (ensure it's not associated with any tournaments first).

   c. Tournament Creation and Invitation:
   - Use `CreateTournament` to create a new tournament, using the IDs of the created league and format.
   - Verify that invitation emails are sent to relevant schools (check your email service or logs).
   - Use `GetTournament` to retrieve the created tournament details.

   d. Tournament Listing and Updates:
   - Use `ListTournaments` to get a list of tournaments.
   - Use `UpdateTournament` to modify a tournament's details.
   - Use `GetTournament` again to verify the changes.

   e. Tournament Deletion:
   - Use `DeleteTournament` to remove a tournament.
   - Attempt to `GetTournament` for the deleted tournament (should fail).

4. For each test, verify that the appropriate email notifications are sent (tournament creation confirmation, invitations).

Remember to include the authentication token in the request body for each request:
- Key: `token`
- Value: `<token_received_from_login>`

Note: When testing deletion operations, ensure that you're not deleting entities that are still referenced by others (e.g., don't delete a league or format that's still used by an active tournament).
