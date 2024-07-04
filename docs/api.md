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

## Testing with Postman

1. Set up a gRPC request in Postman.
2. Use `localhost:10000` as the server URL (Envoy proxy).
3. Import the `.proto` file and select the desired method.
4. Input the demo data in the "Message" tab.
5. Click "Invoke" to send the request.

For authenticated requests, add the token to the metadata:
- Key: `authorization`
- Value: `Bearer <token_received_from_login>`