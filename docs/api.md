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
  "schoolID": "1234"
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
  "safeguardingCertificate": "cert123"
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

Description: Enable 2FA for a user account.

Demo Data:
```json
{
  "userID": 1
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

Detailed error messages are provided in the response for easier debugging and user feedback.