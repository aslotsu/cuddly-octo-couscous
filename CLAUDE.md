# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go REST API built with Gin-gonic framework for managing forms. Uses PostgreSQL with pgx/v5 connection pooling for database operations. The API stores form data as JSONB in PostgreSQL.

Module name: `github.com/aslotsu/monkreflections-form-api`

## Development Commands

### Running the Application
```bash
go run main.go
```
Server starts on port 8080.

### Building
```bash
go build -o form-api
```

### Managing Dependencies
```bash
go mod tidy        # Clean up dependencies
go mod download    # Download dependencies
```

## Architecture

### Application Flow
1. **main.go**: Entry point - loads environment variables, establishes database connection, sets up Gin router with CORS, registers routes, starts server
2. **config/database.go**: Database connection using pgxpool, auto-creates tables on startup
3. **config/cors.go**: CORS configuration allowing monkreflections.com and localhost:3000
4. **handlers/form.go**: HTTP handlers with dependency injection pattern (FormHandler struct receives db pool)
5. **models/form.go**: Data models and request/response structs

### Database Connection Pattern
- Uses `pgxpool.Pool` (connection pooling) from jackc/pgx/v5
- Connection established via `DATABASE_URL` environment variable
- Tables auto-created on startup via `CreateTables()` function in config/database.go:68
- Single forms table with JSONB data column

### Handler Pattern
- Handlers use dependency injection: `FormHandler` struct holds db connection pool
- Create handlers with `NewFormHandler(db)` constructor
- All handlers use `context.Background()` for database operations
- Error responses use Gin's `c.JSON()` with appropriate HTTP status codes

### Data Storage
- Forms have flexible JSONB `data` field stored as JSON string in PostgreSQL
- Request models use `map[string]any` for incoming JSON
- JSON marshaling/unmarshaling handled in handlers before database operations

## API Routes

All routes under `/api/forms`:
- `GET /api/forms` - List all forms
- `POST /api/forms` - Create form (requires title and data)
- `GET /api/forms/:id` - Get single form
- `PUT /api/forms/:id` - Update form (partial updates supported)
- `DELETE /api/forms/:id` - Delete form

## Environment Configuration

Required environment variable:
- `DATABASE_URL`: PostgreSQL connection string (format: `postgres://user:password@host:port/dbname?sslmode=mode`)

See `.env.example` for reference configuration.

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS forms (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

The `updated_at` field is managed via SQL query in handlers/form.go:135 (UPDATE queries set `updated_at = CURRENT_TIMESTAMP`).
