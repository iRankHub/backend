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
   - Rename the `.env.example` to `.env` file in the root directory and add your database configuration:

4. Install Air for live reloading:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

5. Start the server:
   ```bash
   air
   ```

   This command will start the gRPC server on port 8080, run database migrations, and start the Envoy proxy on port 10000.

## Project Structure

The project follows a standard Go project structure:

- `cmd/`: Contains the main entry points for the application.
- `docs/`: Contains the project documentation.
- `internal/`: Contains the internal packages used by the application.
- `pkg/`: Contains reusable packages used across the application.
- `scripts/`: Contains shell scripts for building, deploying, and testing the application.
- `tests/`: Contains the integration and unit tests for the backend components.

## API Documentation

For detailed API documentation, refer to the `docs/api.md` file.

## Architecture

For information about the system architecture, refer to the `docs/architecture.md` file.

## Database

For database-related information, refer to the `docs/database.md` file.

## Deployment

For deployment instructions, refer to the `docs/deployment.md` file.

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