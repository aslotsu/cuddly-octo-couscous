# Form API Project

## Project Overview

This is a Go-based REST API for managing form data, built with the Gin web framework and PostgreSQL database. The application provides CRUD (Create, Read, Update, Delete) operations for forms stored in a PostgreSQL database with JSONB data fields. The API includes CORS configuration for cross-origin requests and uses environment variables for configuration.

**Main Technologies:**
- Go 1.25.4
- Gin web framework
- PostgreSQL with pgx driver
- GORM for database interactions
- CORS middleware

**Key Features:**
- RESTful API endpoints for form management
- JSONB data storage for flexible form data
- Environment-based configuration
- CORS support for web applications
- PostgreSQL database integration

## Project Structure

```
/Users/alfredlotsu/backend/form-api/
├── .env.example          # Example environment variables
├── .gitignore           # Git ignore rules
├── CLAUDE.md            # Additional documentation
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Main application entry point
├── config/              # Configuration files
│   ├── cors.go          # CORS configuration
│   └── database.go      # Database connection and setup
├── handlers/            # API handlers
│   └── form.go          # Form-related HTTP handlers
└── models/              # Data models
    └── form.go          # Form data structures
```

## Building and Running

### Prerequisites
- Go 1.25.4 or later
- PostgreSQL database

### Setup Instructions

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up the database:**
   - Install and start PostgreSQL
   - Create a database for the application (e.g., `forms_db`)

3. **Configure environment variables:**
   - Copy `.env.example` to `.env`:
     ```bash
     cp .env.example .env
     ```
   - Update the `DATABASE_URL` in the `.env` file to match your PostgreSQL setup

4. **Run the application:**
   ```bash
   go run main.go
   ```

The API will start on `http://localhost:8080`.

### API Endpoints

| Method | Endpoint               | Description                    |
|--------|------------------------|--------------------------------|
| GET    | `/api/forms`           | Retrieve all forms             |
| GET    | `/api/forms/:id`       | Retrieve a specific form       |
| POST   | `/api/forms`           | Create a new form              |
| PUT    | `/api/forms/:id`       | Update an existing form        |
| DELETE | `/api/forms/:id`       | Delete a form                  |

### Request/Response Examples

**Create Form (POST /api/forms):**
```json
{
  "title": "Contact Form",
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "message": "Hello world"
  }
}
```

**Response:**
```json
{
  "id": 1
}
```

## Development Conventions

### Code Structure
- Handlers are organized by functionality in the `handlers/` directory
- Data models are defined in the `models/` directory
- Configuration is handled in the `config/` directory
- The main entry point is in `main.go`

### Database Schema
The application creates and uses a single table called `forms` with the following structure:
- `id`: SERIAL PRIMARY KEY
- `title`: VARCHAR(255) NOT NULL
- `data`: JSONB NOT NULL (stores form data as JSON)
- `created_at`: TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- `updated_at`: TIMESTAMP DEFAULT CURRENT_TIMESTAMP

### Error Handling
- Standard HTTP status codes are used
- JSON error responses follow the format: `{"error": "error message"}`
- Validation errors are handled by Gin's binding functionality

### CORS Configuration
The API supports cross-origin requests from:
- `https://monkreflections.com`
- `http://localhost:3000`

## Testing

There are no explicit test files in the current structure, but tests can be added following Go's testing conventions in files with the `_test.go` suffix.

## Security Considerations

- Input validation is performed using Gin's binding features
- SQL queries use parameterized statements to prevent injection
- CORS is configured for specific origins only

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string (required)

## Deployment

For production deployment:
1. Ensure the database is properly configured and secured
2. Set the `DATABASE_URL` environment variable
3. Consider using a process manager to keep the application running
4. Set up proper logging and monitoring