# API Documentation

## Overview

The iRankHub backend API is built using gRPC and is accessible through an Envoy proxy.

- gRPC Server Port: 8080
- Envoy Proxy Port: 10000

## Testing with Postman

To test the API endpoints using Postman:

1. Set up a new gRPC request in Postman.
2. Use `localhost:10000` as the server URL (Envoy proxy address).
3. Import the `.proto` file into Postman and select the desired method.
4. Input the appropriate demo data in the "Message" tab, including the `token` field for authenticated requests.
5. Click "Invoke" to send the request.

### Testing Flow

1. Start by using the SignUp endpoint to create a new user.
2. Use the Login endpoint to authenticate and receive a token.
3. Include the token in the request body for subsequent authenticated requests.
4. Test other endpoints as needed, ensuring to use the correct user ID and token.

## Error Handling

The API uses standard gRPC error codes. Common errors include:

- `INVALID_ARGUMENT`: Missing or invalid input data
- `UNAUTHENTICATED`: Invalid or missing authentication token
- `NOT_FOUND`: Requested resource not found
- `INTERNAL`: Server-side error
- `PERMISSION_DENIED`: User doesn't have the required permissions for the action

Detailed error messages are provided in the response for easier debugging and user feedback.