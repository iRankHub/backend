# Architecture

## Overview

The iRankHub backend is built using a microservices architecture with gRPC for communication.

## Components

1. gRPC Server
   - Handles all business logic and database operations
   - Implemented in Go
   - Listens on port 8080

2. Envoy Proxy
   - Acts as a reverse proxy for the gRPC server
   - Handles protocol translation (HTTP/1.1 to gRPC)
   - Listens on port 10000

3. PostgreSQL Database
   - Stores all application data

4. Docker
   - Used for containerizing the Envoy proxy

## Flow

1. Client sends requests to Envoy proxy (port 10000)
2. Envoy forwards requests to gRPC server (port 8080)
3. gRPC server processes requests and interacts with the database
4. Responses follow the reverse path back to the client

## Security

- Authentication is handled using PASETO tokens
- Passwords are hashed using bcrypt before storage
-
## Network Protocol

The system prioritizes HTTP/2 for improved security and performance:

- Envoy is configured to use HTTP/2 by default for incoming connections.
- If a client doesn't support HTTP/2, the system will automatically fall back to HTTP/1.1.
- Communication between Envoy and the gRPC backend always uses HTTP/2.

This configuration ensures optimal performance for modern clients while maintaining compatibility with older systems.

## Future Considerations

- Implement service discovery for scalability
- Add caching layer for improved performance
- Implement logging and monitoring solutions