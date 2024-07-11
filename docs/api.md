# API Documentation

## Overview

The iRankHub backend API is built using gRPC and is accessible through an Envoy proxy.

- gRPC Server Port: 8080
- Envoy Proxy Port: 10000

## Authentication API

### SignUp

Endpoint: `AuthService.SignUp`

Demo Data:

1. Student Sign Up:
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@student.com",
  "password": "studentPass123!",
  "userRole": "student",
  "dateOfBirth": "2005-05-15",
  "grade": "Grade 4",
  "schoolID": 2
}
```

2. School Sign Up:
```json
{
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "jane.smith@school.com",
  "password": "schoolAdmin456!",
  "userRole": "school",
  "schoolName": "Springfield High",
  "address": "KK 123 St",
  "country": "USA",
  "province": "Illinois",
  "district": "Springfield",
  "contactEmail": "contact@springfieldhigh.edu",
  "schoolType": "Public"
}
```

3. Volunteer Sign Up:
```json
{
  "firstName": "Alex",
  "lastName": "Johnson",
  "email": "alex.johnson@volunteer.com",
  "password": "volunteer789!",
  "userRole": "volunteer",
  "dateOfBirth": "1990-08-20",
  "roleInterestedIn": "Mentor",
  "graduationYear": 2012,
  "safeguardingCertificate": true
}
```

### Login

Endpoint: `AuthService.Login`

Demo Data:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Enable Two-Factor Authentication

Endpoint: `AuthService.EnableTwoFactor`

Description: Enable 2FA for a user account. Requires authentication.

Demo Data:
```json
{
  "userID": 1,
  "token": "your_auth_token_here"
}
```

### Verify Two-Factor Authentication

Endpoint: `AuthService.VerifyTwoFactor`

Description: Verify the 2FA code provided by the user.

Demo Data:
```json
{
  "userID": 1,
  "code": "123456"
}
```

### Disable Two-Factor Authentication

Endpoint: `AuthService.DisableTwoFactor`

Description: Disable 2FA for a user account. Requires authentication.

Demo Data:
```json
{
  "userID": 1,
  "token": "your_auth_token_here"
}
```
## User Management API

### GetPendingUsers

Endpoint: `UserManagementService.GetPendingUsers`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here"
}
```

### GetUserDetails

Endpoint: `UserManagementService.GetUserDetails`
Authorization: User can access their own details, Admin can access any user's details

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### ApproveUser

Endpoint: `UserManagementService.ApproveUser`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### RejectUser

Endpoint: `UserManagementService.RejectUser`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### UpdateUserProfile

Endpoint: `UserManagementService.UpdateUserProfile`
Authorization: User can update their own profile

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123,
  "name": "John Doe",
  "email": "john.doe@example.com",
  "address": "123 Main St",
  "phone": "555-1234",
  "bio": "A brief bio",
  "profilePicture": ""
}
```

### DeleteUserProfile

Endpoint: `UserManagementService.DeleteUserProfile`
Authorization: User can delete their own profile, Admin can delete any profile

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### DeactivateAccount

Endpoint: `UserManagementService.DeactivateAccount`
Authorization: User can deactivate their own account, Admin can deactivate any account

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### ReactivateAccount

Endpoint: `UserManagementService.ReactivateAccount`
Authorization: User can reactivate their own account, Admin can reactivate any account

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### GetAccountStatus

Endpoint: `UserManagementService.GetAccountStatus`
Authorization: User can get their own account status, Admin can get any account's status

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

## Testing User Management Features

To test the user management features:

1. Start by creating a new user using the `SignUp` endpoint.
2. Use the `Login` endpoint to authenticate and receive a token.
3. Use the token in the metadata for subsequent authenticated requests.
4. Test the following scenarios:

   a. Pending User Approval:
   - Use `GetPendingUsers` to retrieve a list of pending users.
   - Use `ApproveUser` or `RejectUser` to approve or reject a pending user.
   - Verify the user's status has changed using `GetUserDetails`.

   b. User Profile Management:
   - Use `UpdateUserProfile` to modify a user's profile information.
   - Use `GetUserDetails` to verify the changes.
   - Use `DeleteUserProfile` to remove a user's profile.

   c. Account Deactivation and Reactivation:
   - Use `DeactivateAccount` to deactivate a user's account.
   - Attempt to log in with the deactivated account (should fail).
   - Use `ReactivateAccount` to reactivate the account.
   - Verify successful login after reactivation.

   d. Account Status:
   - Use `GetAccountStatus` to check the current status of a user's account at any point.

5. For each test, verify that the appropriate email notifications are sent (approval, rejection, deactivation, reactivation).

## Testing Two-Factor Authentication

To test the two-factor authentication flow:

1. Log in using the `Login` endpoint to obtain an authentication token.
2. Use the `EnableTwoFactor` endpoint with the obtained token to enable 2FA for the user.
3. The response will include a secret and a QR code URL. Use an authenticator app to scan the QR code or manually enter the secret.
4. Generate a 2FA code using the authenticator app.
5. Use the `VerifyTwoFactor` endpoint with the generated code to verify and complete the 2FA setup.
6. For subsequent logins, the `Login` endpoint will return `requireTwoFactor: true` if 2FA is enabled.
7. Provide the 2FA code using the `VerifyTwoFactor` endpoint to complete the login process.
8. To disable 2FA, use the `DisableTwoFactor` endpoint with a valid authentication token.

### Request Password Reset

Endpoint: `AuthService.RequestPasswordReset`

Description: Request a password reset for a user account.

Demo Data:
```json
{
  "email": "user@example.com"
}
```

### Reset Password

Endpoint: `AuthService.ResetPassword`

Description: Reset a user's password using the provided token.

Demo Data:
```json
{
  "token": "reset_token_here",
  "newPassword": "newPassword123!"
}
```

### Enable Biometric Login

Endpoint: `AuthService.EnableBiometricLogin`

Description: Enable biometric login for a user account.

Demo Data:
```json
{
  "userID": 1
}
```

### Biometric Login

Endpoint: `AuthService.BiometricLogin`

Description: Authenticate a user using a biometric token.

Demo Data:
```json
{
  "biometricToken": "token_here"
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
  "startDate": "2023-07-15T09:00:00Z",
  "endDate": "2023-07-17T18:00:00Z",
  "location": "City Convention Center",
  "formatId": 1,
  "leagueId": 2,
  "numberOfPreliminaryRounds": 4,
  "numberOfEliminationRounds": 2,
  "judgesPerDebatePreliminary": 3,
  "judgesPerDebateElimination": 5,
  "tournamentFee": 100.00
}
```

### GetTournament

Endpoint: `TournamentService.GetTournament`

Request:
```json
{
  "tournamentId": 1
}
```

### ListTournaments

Endpoint: `TournamentService.ListTournaments`

Request:
```json
{
  "pageSize": 10,
  "pageToken": 0
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
  "tournamentFee": 120.00
}
```

### DeleteTournament

Endpoint: `TournamentService.DeleteTournament`
Authorization: Admin only

Request:
```json
{
  "tournamentId": 1
}
```

## Testing Tournament Management Features

To test the tournament management features:

1. Use the `Login` endpoint to authenticate as an admin and receive a token.
2. Use the token in the metadata for subsequent authenticated requests.
3. Test the following scenarios:

   a. Tournament Creation and Invitation:
   - Use `CreateTournament` to create a new tournament.
   - Verify that invitation emails are sent to relevant schools (check your email service or logs).
   - Use `GetTournament` to retrieve the created tournament details.

   b. Tournament Listing and Updates:
   - Use `ListTournaments` to get a list of tournaments.
   - Use `UpdateTournament` to modify a tournament's details.
   - Use `GetTournament` again to verify the changes.

   c. Tournament Deletion:
   - Use `DeleteTournament` to remove a tournament.
   - Attempt to `GetTournament` for the deleted tournament (should fail).

4. For each test, verify that the appropriate email notifications are sent (tournament creation confirmation, invitations).

Remember to include the authentication token in the metadata for each request:
- Key: `authorization`
- Value: `Bearer <token_received_from_login>`
  
## Testing with Postman

To test the API endpoints using Postman:

1. Set up a new gRPC request in Postman.
2. Use `localhost:10000` as the server URL (Envoy proxy address).
3. Import the `.proto` file into Postman and select the desired method.
4. Input the appropriate demo data in the "Message" tab.
5. For authenticated requests, add the token to the metadata:
   - Key: `authorization`
   - Value: `Bearer <token_received_from_login>`
6. Click "Invoke" to send the request.

### Testing Flow

1. Start by using the SignUp endpoint to create a new user.
2. Use the Login endpoint to authenticate and receive a token.
3. Use the token in the metadata for subsequent authenticated requests.
4. Test other endpoints as needed, ensuring to use the correct user ID and token.

## Error Handling

The API uses standard gRPC error codes. Common errors include:

- `INVALID_ARGUMENT`: Missing or invalid input data
- `UNAUTHENTICATED`: Invalid or missing authentication token
- `NOT_FOUND`: Requested resource not found
- `INTERNAL`: Server-side error
- `PERMISSION_DENIED`: User doesn't have the required permissions for the action

Detailed error messages are provided in the response for easier debugging and user feedback.