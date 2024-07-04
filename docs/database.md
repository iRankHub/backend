# Database

## Overview

iRankHub uses PostgreSQL as its primary database.

## Configuration

Database configuration is stored in the `.env` file and loaded using Viper. The connection string is constructed in the `main.go` file.

## Migrations

Database migrations are handled using the `golang-migrate` library. Migration files are located in `internal/database/postgres/migrations`.

To run migrations:
```bash
go run cmd/server/main.go
```
This command will automatically run all pending migrations before starting the server.

## Models

Database models are defined in the `internal/models` directory. These models correspond to the database tables and are used throughout the application for data operations.

## Queries

Database queries are generated using `sqlc`. The query definitions are located in `internal/database/postgres/queries`.

To regenerate queries after modifying the SQL files:
```bash
sqlc generate
```

## Connections

Database connections are managed using a connection pool to ensure efficient use of resources.

## Backup and Recovery

[To be implemented]

## Performance Tuning

[To be implemented]