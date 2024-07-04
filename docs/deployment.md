# Deployment

## Development Environment

### Prerequisites

- Go (version 1.22.2 or later)
- PostgreSQL
- Docker
- Air (for live reloading)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/iRankHub/backend.git
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the `.env` file with your database credentials and other configuration.

4. Install Air for live reloading:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

5. Run the application:
   ```bash
   air
   ```

   This will start the server, run migrations, and start the Envoy proxy.

### Troubleshooting

- Ensure PostgreSQL is running and accessible.
- Check that all required environment variables are set in the `.env` file.
- If you encounter any dependency issues, try cleaning the Go module cache:
  ```bash
  go clean -modcache
  go mod download
  ```

## Production Environment

[To be implemented]