# User Management API

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
### GetStudents

Endpoint: `UserManagementService.GetStudents`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "page": 1,
  "pageSize": 10
}
```

### GetVolunteers

Endpoint: `UserManagementService.GetVolunteers`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "page": 1,
  "pageSize": 10
}
```

### GetSchools

Endpoint: `UserManagementService.GetSchools`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "page": 1,
  "pageSize": 10
}
```


### GetCountries

Endpoint: `UserManagementService.GetCountries`
Authorization: Any authenticated user

Request:
```json
{
  "token": "your_auth_token_here"
}
```

## Testing User Management Features

To test the user management features:

1. Start by creating a new user using the `SignUp` endpoint.
2. Use the `Login` endpoint to authenticate and receive a token.
3. Include the token in the request body for subsequent authenticated requests.
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

   e. Student Management:
   - Use `GetStudents` to retrieve a paginated list of students.
   - Verify the student information is correct and complete.

   f. Volunteer Management:
   - Use `GetVolunteers` to retrieve a paginated list of volunteers.
   - Verify the volunteer information is correct and complete.

   g. School Management:
   - Use `GetSchools` to retrieve a paginated list of schools.
   - Verify the school information is correct and complete.

   h. Country Information:
   - Use `GetCountries` to retrieve a list of countries and their codes.
   - Verify the country information is correct and complete.

5. For each test, verify that the appropriate email notifications are sent (approval, rejection, deactivation, reactivation).
6. Test pagination for endpoints that support it (GetStudents, GetVolunteers, GetSchools) by varying the page and pageSize parameters.