# QnA - Question and Answer API

A RESTful API service built with Go for managing questions and answers. This application provides endpoints to create, read, and delete questions, as well as manage answers associated with those questions.

## Features

- ✅ Create, read, and delete questions
- ✅ Create, read, and delete answers
- ✅ PostgreSQL database with migrations
- ✅ Docker containerization
- ✅ Clean architecture with separation of concerns
- ✅ UUID-based entity identification
- ✅ Timestamp tracking for all entities

## Tech Stack

- **Language**: Go 1.24
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **HTTP Server**: Standard library `net/http`
- **Migrations**: Goose
- **Containerization**: Docker & Docker Compose

## Project Structure

```
QnA/ 
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/
│   │   ├── migrations/          # Database migrations
│   │   └── postgres/            # Database connection
│   ├── domain/                  # Domain models
│   ├── repository/              # Data access layer
│   ├── usecase/                 # Business logic
│   ├── transport/
│   │   └── my_http/             # HTTP handlers
│   ├── server/                  # HTTP server setup
│   └── shared/                  # Shared utilities and errors
└── docker/
    ├── Dockerfile               # Application Docker image
    └── docker-compose.yaml      # Docker Compose configuration

```

## Prerequisites

- Go 1.24 or later
- Docker and Docker Compose
- PostgreSQL 15 (if running locally without Docker)

## Installation

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd QnA
```

2. Create a `.env` file in the project root (optional, defaults are provided):
```env
POSTGRES_DB=qa
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
DATABASE_URL=postgres://user:pass@db:5432/qa?sslmode=disable
HTTP_PORT=8080
```

3. Build and start the services:
```bash
cd docker
docker compose up --build
```

The application will:
- Start PostgreSQL database
- Run database migrations automatically
- Start the API server on port 8080

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Set up PostgreSQL database and create a database named `qa` (or your preferred name).

3. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/qa?sslmode=disable"
export HTTP_PORT=8080
```

Or create a `.env` file in the project root:
```env
DATABASE_URL=postgres://user:password@localhost:5432/qa?sslmode=disable
HTTP_PORT=8080
```

4. Run migrations:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./internal/database/migrations postgres "$DATABASE_URL" up
```

5. Run the application:
```bash
go run cmd/app/main.go
```

## API Endpoints

### Questions

#### Create Question
```http
POST /questions
Content-Type: application/json

{
  "text": "What is Go?"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "text": "What is Go?",
  "created_at": "2025-12-04T18:20:31Z"
}
```

#### Get All Questions
```http
GET /questions
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "text": "What is Go?",
    "created_at": "2025-12-04T18:20:31Z"
  }
]
```

#### Get Question by ID
```http
GET /questions/{id}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "text": "What is Go?",
  "created_at": "2025-12-04T18:20:31Z"
}
```

#### Delete Question
```http
DELETE /questions/{id}
```

**Response:** `204 No Content`

### Answers

#### Create Answer
```http
POST /questions/{question_id}/answers
Content-Type: application/json

{
  "user_id": "user123",
  "text": "Go is a programming language developed by Google."
}
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "question_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user123",
  "text": "Go is a programming language developed by Google.",
  "created_at": "2025-12-04T18:20:31Z"
}
```

#### Get Answer by ID
```http
GET /answers/{id}
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "question_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user123",
  "text": "Go is a programming language developed by Google.",
  "created_at": "2025-12-04T18:20:31Z"
}
```

#### Delete Answer
```http
DELETE /answers/{id}
```

**Response:** `204 No Content`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `HTTP_PORT` | Port for HTTP server | `8080` |
| `POSTGRES_DB` | PostgreSQL database name | `qa` |
| `POSTGRES_USER` | PostgreSQL username | `user` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `pass` |

## Docker Services

The `docker-compose.yaml` defines three services:

1. **db**: PostgreSQL 15 database
   - Port: `5433:5432` (host:container)
   - Persistent volume: `postgres_data`

2. **migrate**: Runs database migrations
   - Automatically runs on startup
   - Exits after completion

3. **app**: Main application server
   - Port: `8080:8080`
   - Depends on `db` and `migrate`

## Database Migrations

Migrations are managed using [Goose](https://github.com/pressly/goose). Migration files are located in `internal/database/migrations/`.

### Running Migrations Manually

```bash
# Up
goose -dir ./internal/database/migrations postgres "$DATABASE_URL" up

# Down
goose -dir ./internal/database/migrations postgres "$DATABASE_URL" down

# Status
goose -dir ./internal/database/migrations postgres "$DATABASE_URL" status
```

## Development

### Project Architecture

The project follows clean architecture principles:

- **Domain**: Core business entities (Question, Answer)
- **Repository**: Data access layer (database operations)
- **Usecase**: Business logic layer
- **Transport**: HTTP handlers and request/response handling
- **Server**: HTTP server configuration and routing

### Adding New Features

1. Define domain models in `internal/domain/`
2. Create repository interfaces and implementations in `internal/repository/`
3. Implement business logic in `internal/usecase/`
4. Add HTTP handlers in `internal/transport/my_http/`
5. Register routes in `internal/server/router.go`
6. Create database migrations if needed

## Testing

Run tests (when implemented):
```bash
go test ./...
```

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]

