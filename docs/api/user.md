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
### ApproveUsers

Endpoint: `UserManagementService.ApproveUsers`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "userIDs": [123, 124, 125]
}
```

### RejectUsers

Endpoint: `UserManagementService.RejectUsers`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "userIDs": [123, 124, 125]
}
```

### DeleteUsers

Endpoint: `UserManagementService.DeleteUsers`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "userIDs": [123, 124, 125]
}
```

... [Rest of the API documentation remains unchanged] ...

## Testing User Management Features

To test the user management features:

1. Start by creating multiple new users using the `SignUp` endpoint.
2. Use the `Login` endpoint to authenticate as an admin and receive a token.
3. Include the token in the request body for subsequent authenticated requests.
4. Test the following scenarios:

   a. Pending User Approval (Single and Batch):
   - Use `GetPendingUsers` to retrieve a list of pending users.
   - Use `ApproveUser` to approve a single pending user.
   - Use `ApproveUsers` to approve multiple pending users at once.
   - Verify the users' statuses have changed using `GetUserDetails`.

   b. User Rejection (Single and Batch):
   - Use `RejectUser` to reject a single pending user.
   - Use `RejectUsers` to reject multiple pending users at once.
   - Verify the users' statuses have changed or that they have been removed using `GetUserDetails`.

   c. User Deletion (Single and Batch):
   - Use `DeleteUserProfile` to delete a single user's profile.
   - Use `DeleteUsers` to delete multiple users at once.
   - Attempt to retrieve the deleted users' details (should fail).

   d. Batch Operation Error Handling:
   - Attempt to approve, reject, or delete a mix of valid and invalid user IDs.
   - Check the response for `failedUserIDs` to ensure proper error handling.

   e. User Profile Management:
   - Use `UpdateUserProfile` to modify a user's profile information.
   - Use `GetUserDetails` to verify the changes.

   f. Account Deactivation and Reactivation:
   - Use `DeactivateAccount` to deactivate a user's account.
   - Attempt to log in with the deactivated account (should fail).
   - Use `ReactivateAccount` to reactivate the account.
   - Verify successful login after reactivation.

   g. Account Status:
   - Use `GetAccountStatus` to check the current status of a user's account at any point.

   h. Student Management:
   - Use `GetStudents` to retrieve a paginated list of students.
   - Verify the student information is correct and complete.

   i. Volunteer Management:
   - Use `GetVolunteers` to retrieve a paginated list of volunteers.
   - Verify the volunteer information is correct and complete.

   j. School Management:
   - Use `GetSchools` to retrieve a paginated list of schools.
   - Verify the school information is correct and complete.

   k. Country Information:
   - Use `GetCountries` to retrieve a list of countries and their codes.
   - Verify the country information is correct and complete.

5. For each test, verify that the appropriate email notifications are sent:
   - Single user operations: approval, rejection, deletion, deactivation, reactivation
   - Batch operations: verify that emails are sent for each successfully processed user

6. Test pagination for endpoints that support it (GetStudents, GetVolunteers, GetSchools) by varying the page and pageSize parameters.

7. Performance Testing for Batch Operations:
   - Test `ApproveUsers`, `RejectUsers`, and `DeleteUsers` with varying numbers of user IDs (e.g., 10, 100, 1000) to ensure the system can handle large batches efficiently.
   - Monitor response times and system resources during these tests.

8. Concurrency Testing:
   - Simulate multiple admins performing batch operations simultaneously to ensure data consistency and proper handling of concurrent requests.

9. Edge Cases for Batch Operations:
   - Test with an empty list of user IDs.
   - Test with a list containing only invalid user IDs.
   - Test with a very large list of user IDs to verify any upper limits on batch size.

10. Authorization Testing:
    - Attempt to use batch operations with non-admin user tokens to ensure proper access control.