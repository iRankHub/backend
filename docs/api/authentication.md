# Authentication API

### SignUp

Endpoint: `AuthService.SignUp`

Demo Data:

1. Admin Sign Up:
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@admin.com",
  "password": "admin",
  "userRole": "admin",
  "gender": "male",
}
```

1. Student Sign Up:
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@student.com",
  "password": "studentPass123!",
  "userRole": "student",
  "gender": "male",
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
  "gender": "female",
  "schoolName": "Springfield High",
  "address": "KK 123 St",
  "country": "United States of America",
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
  "gender": "female",
  "dateOfBirth": "1990-08-20",
  "roleInterestedIn": "Mentor",
  "graduationYear": 2012,
  "hasInternship": true,
  "isEnrolledInUniversity": true,
  "safeguardingCertificate": true
}
```

### Batch Import Users

Endpoint: `AuthService.BatchImportUsers`

Description: Import multiple users at once. This endpoint is typically used by admins to bulk import user data.

Demo Data:
```json
{
  "users": [
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "userRole": "student",
      "gender": "female",
      "dateOfBirth": "2005-05-15",
      "grade": "Grade 4",
      "schoolID": 2
    },
    {
      "firstName": "Jane",
      "lastName": "Smith",
      "email": "jane.smith@example.com",
      "userRole": "school",
      "schoolName": "Springfield High",
      "address": "KK 123 St",
      "country": "United States of America",
      "province": "Illinois",
      "district": "Springfield",
      "contactEmail": "contact@springfieldhigh.edu",
      "schoolType": "Public"
    },
    {
      "firstName": "Alex",
      "lastName": "Johnson",
      "email": "alex.johnson@example.com",
      "userRole": "volunteer",
      "gender": "female",
      "dateOfBirth": "1990-08-20",
      "roleInterestedIn": "Mentor",
      "graduationYear": 2012,
      "hasInternship": true,
      "isEnrolledInUniversity": true,
      "safeguardingCertificate": true
    }
  ]
}
```

Note:
- Passwords for imported users are automatically generated and sent to their email addresses.
- Users are notified via email about their account creation and temporary password.
- If any imports fail, their email addresses will be listed in the `failedEmails` array.

### Login

Endpoint: `AuthService.Login`

Demo Data:

```json
{
  "email_or_id": "user@example.com",
  "password": "password123"
}
```
Note: The login process now has two steps when 2FA is enabled:
1. Initial login attempt with email/ID and password
2. If 2FA is required, use the `VerifyTwoFactor` endpoint to complete the authentication

### Generate Two-Factor Authentication OTP

Endpoint: `AuthService.GenerateTwoFactorOTP`

Description: Generate and send a new 2FA OTP to the user's email. This can be used when setting up 2FA or when the user needs a new OTP.

Demo Data:
```json
{
  "email": "user@example.com"
}
```

Response:
```json
{
  "success": true,
  "message": "Two-factor authentication OTP sent. Please check your email."
}
```

### Verify Two-Factor Authentication

Endpoint: `AuthService.VerifyTwoFactor`

Description: Verify the 2FA code provided by the user. This is used during the login process when 2FA is required.

Demo Data:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

Response:
```json
{
  "success": true,
  "message": "Two-factor authentication verified successfully."
}
```

## Testing Two-Factor Authentication

To test the two-factor authentication flow:

1. Log in using the `Login` endpoint to obtain an authentication token.
2. Use the `EnableTwoFactor` endpoint with the obtained token to enable 2FA for the user.
3. The response will include a secret and a QR code URL. Use an authenticator app to scan the QR code or manually enter the secret.
4. Generate a 2FA code using the authenticator app.
5. Use the `VerifyTwoFactor` endpoint with the generated code and the token to verify and complete the 2FA setup.
6. For subsequent logins, the `Login` endpoint will return `requireTwoFactor: true` if 2FA is enabled.
7. Provide the 2FA code and the token using the `VerifyTwoFactor` endpoint to complete the login process.
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

### Begin WebAuthn Registration

Endpoint: `AuthService.BeginWebAuthnRegistration`

Description: Begin the WebAuthn registration process for a user account. Requires authentication.

Demo Data:
```json
{
  "userID": 1,
  "token": "your_auth_token_here"
}
```

### Finish WebAuthn Registration

Endpoint: `AuthService.FinishWebAuthnRegistration`

Description: Complete the WebAuthn registration process for a user account. Requires authentication.

Demo Data:
```json
{
  "userID": 1,
  "token": "your_auth_token_here",
  "credential": "base64_encoded_credential_data"
}
```

### Begin WebAuthn Login

Endpoint: `AuthService.BeginWebAuthnLogin`

Description: Begin the WebAuthn login process for a user account.

Demo Data:
```json
{
  "email": "user@example.com"
}
```

### Finish WebAuthn Login

Endpoint: `AuthService.FinishWebAuthnLogin`

Description: Complete the WebAuthn login process for a user account.

Demo Data:
```json
{
  "email": "user@example.com",
  "credential": "base64_encoded_credential_data"
}
```


Endpoint: `AuthService.BiometricLogin`

Description: Authenticate a user using a biometric token.

Demo Data:
```json
{
  "biometricToken": "token_here"
}
```