# echo-basics

> A simple REST API built with Go + Echo + GORM for storing and managing logs.
The project demonstrates a clean backend architecture with separated layers for handlers, models, utilities, and routes.
> This project is based on and forked from the original tutorial repository, this particular repository is my interpretation built on it. The repository mentioned:  
https://github.com/smsk-dev/go-basics/tree/main/echo-basics

---

## Project Structure

```
cmd/
  api/
    main.go        -> application entry point
  migrate/
    main.go        -> migration runner

internal/

  app/
    context.go     -> shared application context
    health.go      -> health check endpoint

  config/
    config.go      -> environment configuration

  db/
    db.go          -> database connection

  handlers/        -> HTTP handlers
    create.go
    delete.go
    fetch.go

  models/          -> database models
    log.go

  middleware/      -> API key authentication
    auth.go

  routes/          -> route registration
    routes.go

  utils/           -> helper utilities
    enums.go
    helpers.go
    levels.go

  migrations/      -> database migrations
    migrations.go
```

---

##Features
- Create logs
- Fetch logs
- Filter logs by:
- ID
- flag
- timestamp
- Delete logs
- API key authentication
- Rate limiting
- Database migrations
- Structured logging
- Modular project structure

## Stack

- [Go 1.25+](https://go.dev/)
- [Echo v5](https://github.com/labstack/echo) — HTTP framework
- [GORM v2](https://gorm.io/) — ORM
- [gormigrate](https://github.com/go-gormigrate/gormigrate) — versioned migrations
- [PostgreSQL](https://www.postgresql.org/) — database

---

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL running locally (or remotely, doesn't matter)
- `jq` (optional, but makes test output much nicer)


## Running The Project

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Configure environment variables

Create a .env file in the project root.
Example:

```
PORT=6070
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=echo_logs
API_KEY=secret
LIMIT_RATE=10
ENV=dev
```

### 3. Run the API

```
go run ./cmd/api
```
The server will start at:
```
http://localhost:6070
```

---

## API Endpoints

### Health Check
Checks if the API is running.

**Request**
```
GET /api/health
```

**Response **
```json
{
  "status": "ok",
  "message": "Yeppers, seems good."
}
```

---

### Create Log
Creates a new log entry.

***Request***
```
POST /api/logs
```


**Request headers**
```
X-API-Key: your_api_key
Content-Type: application/json
```

**Request body**
```json
{
  "flag": "info",
  "message": "something happened"
}
```

***Response***
```json
{
  "id": "uuid",
  "flag": "info",
  "message": "test message",
  "timestamp": "2026-03-15T15:00:00Z"
}
```
---
### Fetch All Logs

**Request**
```
GET /api/logs
```

**Headers**
```
X-API-Key: your_api_key
```

***Response***
```json
{
  "id": "uuid",
  "flag": "info",
  "message": "test message",
  "timestamp": "2026-03-15T15:00:00Z"
}
```

---

### Fetch by ID

***Request***
```
GET /api/logs/id/:id
```
***Example:***
```
GET /api/logs/id/7dff1c48-4a71-4a6c-9a21-0cfa1e5d0e45
````

***Headers***
```
X-API-Key: your_api_key
```

---


### Fetch by Timestamp
Returns logs filtered by timestamp.

***Request***
```
GET /api/logs/timestamp/:timestamp
```

***Example:***
```
GET /api/logs/timestamp/2026-03-15T15:00:00Z
```
***Headers***
```
X-API-Key: your_api_key
```


Returns the latest log at or before the given timestamp. Timestamp must be RFC3339 format (e.g. `2026-02-23T00:00:00Z`).

**Response 200** — single log object  
**Response 400** — bad timestamp  
**Response 404** — nothing found before that timestamp

---

### Fetch by Flag
Returns logs filtered by severity level.

***Request***
```
GET /api/logs/flag/:flag
```

***Example:***
```
GET /api/logs/flag/info`
```

***Headers***
```
X-API-Key: your_api_key
```

---

### Delete Log
Deletes a log entry.

***Request***
```
DELETE /api/logs/:id
```

***Example***
```
DELETE /api/logs/7dff1c48-4a71-4a6c-9a21-0cfa1e5d0e45
```

***Headers***
```
X-API-Key: your_api_key
```

***Note***
Logs with high severity levels (for example error or trace) may be protected from deletion depending on application logic.

---


## Migrations
```
go run ./cmd/migrate
```
This migration will:

- create the uuid-ossp extension
- create the log_flag enum type
- create the logs table

---
##Architecture Notes

The project follows a modular structure:
- handlers handle HTTP requests
- models define database entities
- routes register endpoints
-middleware handles authentication and request filtering
- utils contains shared helper functions
- migrations manage database schema changes
  
This separation improves maintainability and scalability.

---

## License

[MIT](https://devsimsek.mit-license.org) — Metin Şimşek (devsimsek)
