# iRankHub Backend

This repository contains the backend application for the iRankHub project. It provides the necessary APIs and services for managing debate tournaments, users, schools, teams, and more.

## Getting Started

To get started with the iRankHub backend, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/iRankHub/backend.git
   ```

2. Install the dependencies:
   ```bash
   cd backend
   go mod download
   ```

3. Set up the database:
   - Create a PostgreSQL database for the project.
   - Update the database configuration in `internal/config/database.go`.
   - Run the database migrations:
     ```bash
     go run scripts/migrate.go
     ```

4. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

   The server will start running on `http://localhost:8080`.

## Project Structure

The project follows a standard Go project structure:

- `cmd/`: Contains the main entry points for the application.
  - `server/`: Contains the main server application.
- `docs/`: Contains the project documentation.
- `internal/`: Contains the internal packages used by the application.
  - `config/`: Contains the configuration files and utilities.
  - `database/`: Contains the database-related files and migrations.
  - `grpc/`: Contains the gRPC server and protocol buffer files.
  - `handlers/`: Contains the gRPC handler functions.
  - `middleware/`: Contains the middleware functions.
  - `models/`: Contains the data models.
  - `repositories/`: Contains the data access layer and repository interfaces.
  - `services/`: Contains the business logic and service-level operations.
  - `utils/`: Contains utility functions.
- `pkg/`: Contains reusable packages used across the application.
- `scripts/`: Contains shell scripts for building, deploying, and testing the application.
- `tests/`: Contains the integration and unit tests for the backend components.

## Testing

To run the tests, use the following command:

```bash
go test ./...
```

This command will run all the tests in the project.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

N/A