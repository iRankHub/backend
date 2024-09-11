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

### GetAllUsers

Endpoint: `UserManagementService.GetAllUsers`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "page": 1,
  "pageSize": 10
}
```

### GetUserStatistics

Endpoint: `UserManagementService.GetUserStatistics`
Authorization: Admin Only

Request:
```json
{
  "token": "your_auth_token_here"
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

### GetUserProfile

Endpoint: `UserManagementService.GetUserProfile`
Authorization: User can retrieve their own profile, Admin can retrieve any profile

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### UpdateAdminProfile

Endpoint: `UserManagementService.UpdateAdminProfile`
Authorization: Admin only (can only update their own profile)

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123,
  "name": "John Doe",
  "gender": "male",
  "address": "123 Admin St",
  "bio": "Experienced administrator",
  "phone": "555-1234",
  "profilePicture": "base64encodedimage"
}
```

Note: All fields required. Whether you updated it or not send it back in the request.

### UpdateSchoolProfile

Endpoint: `UserManagementService.UpdateSchoolProfile`
Authorization: School contact person only (can only update their own school's profile)

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123,
  "contactPersonName": "Jane Smith",
  "gender": "female",
  "address": "456 School Ave",
  "schoolName": "New Example High School",
  "schoolEmail": "contact@newexample.edu",
  "schoolType": "Private",
  "contactEmail": "jane.smith@newexample.edu",
  "contactPersonNationalId": "ID12345678",
  "phone": "555-5678",
  "profilePicture": "base64_encoded_image",
  "bio": "Dedicated school administrator"
}
```

Note: All fields are required. Whether you updated it or not. It must be included in the request body. Also Only `Approved Users` can update their profile because they are the only one who have their userprofile row filled.

### UpdateStudentProfile

Endpoint: `UserManagementService.UpdateStudentProfile`
Authorization: Student only (can only update their own profile)

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123,
  "firstName": "John",
  "lastName": "Doe",
  "gender": "male",
  "email": "john.doe@example.com",
  "grade": "10th",
  "dateOfBirth": "2005-01-15",
  "address": "123 Student St",
  "bio": "Dedicated student",
  "profilePicture": "base64_encoded_image",
  "phone": "555-1234"
}
```

Note: All fields are required. Whether you updated it or not, it must be included in the request body. Also, only `Approved Users` can update their profile because they are the only ones who have their userprofile row filled.

### UpdateVolunteerProfile

Endpoint: `UserManagementService.UpdateVolunteerProfile`
Authorization: Volunteer only (can only update their own profile)

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123,
  "firstName": "Jane",
  "lastName": "Smith",
  "gender": "female",
  "email": "jane.smith@example.com",
  "nationalID": "ID87654321",
  "graduateYear": 2022,
  "isEnrolledInUniversity": true,
  "hasInternship": false,
  "address": "456 Volunteer Ave",
  "bio": "Passionate volunteer",
  "profilePicture": "base64_encoded_image",
  "role": "Mentor",
  "phone": "555-5678"
}
```

Note: All fields are required. Whether you updated it or not, it must be included in the request body. Also, only `Approved Users` can update their profile because they are the only ones who have their userprofile row filled.

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

Note: This operation performs a soft delete, marking the user's account as deleted but retaining the data for potential future recovery or auditing purposes.

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

### GetVolunteersAndAdmins

Endpoint: `UserManagementService.GetVolunteersAndAdmins`
Authorization: Admin only

Request:
```json
{
  "token": "your_auth_token_here",
  "page": 1,
  "pageSize": 10
}
```

### GetSchoolsNoAuth

Endpoint: `UserManagementService.GetSchoolsNoAuth`
Authorization: No authentication required

Request:
```json
{
  "page": 1,
  "pageSize": 1000
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

### InitiatePasswordUpdate

Endpoint: `UserManagementService.InitiatePasswordUpdate`
Authorization: Authenticated user (can only initiate for their own account)

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123
}
```

### VerifyAndUpdatePassword

Endpoint: `UserManagementService.VerifyAndUpdatePassword`
Authorization: Authenticated user (can only update their own password)

Request:
```json
{
  "token": "your_auth_token_here",
  "userID": 123,
  "verificationCode": "123456",
  "newPassword": "new_secure_password"
}
```

### GetSchoolIDsByNames

Endpoint: `UserManagementService.GetSchoolIDsByNames`
Authorization: Any authenticated user

Request:
```json
{
  "token": "your_auth_token_here",
  "school_names": ["School A", "School B", "School C"]
}
```


Note: This endpoint allows you to retrieve school IDs for multiple school names in a single request. If a school name is not found, it will not be included in the response. This is useful for efficiently mapping school names to their corresponding IDs.

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

   l. User Management:
   - Use `GetAllUsers` to retrieve a list of users.

   m. Password Update:
   - Use `InitiatePasswordUpdate` to start the password update process for a user.
   - Verify that a verification email is sent to the user's email address.
   - Use `VerifyAndUpdatePassword` with the correct verification code and new password.
   - Attempt to log in with the old password (should fail).
   - Verify successful login with the new password.
   - Test with incorrect verification codes and expired codes (should fail).
   - Attempt to update password for a different user ID (should fail).

   n. Admin Profile Update:
   - Use `UpdateAdminProfile` to modify an admin's profile information.
   - Verify that only admins can use this endpoint.
   - Attempt to update another admin's profile (should fail).
   - Use `GetUserProfile` to verify the changes.

   o. School Profile Update:
   - Use `UpdateSchoolProfile` to modify a school's profile information.
   - Verify that only school contact persons can use this endpoint.
   - Attempt to update another school's profile (should fail).
   - Use `GetUserProfile` and `GetSchools` to verify the changes.
   p. School ID Lookup:
   - Use `GetSchoolIDsByNames` to retrieve IDs for multiple school names.
   - Test with a mix of existing and non-existing school names.
   - Verify that only existing schools are returned in the response.
   - Test with an empty list of school names.
   - Test with a large number of school names to ensure performance.



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

11. Password Update Testing:
    - Test initiating password update for non-existent users.
    - Test verifying with incorrect or expired verification codes.
    - Test password update with weak passwords (if password strength rules are implemented).
    - Verify that password update emails are sent and contain the correct information.
    - Test the expiration of verification codes (e.g., try using a code after 15 minutes).