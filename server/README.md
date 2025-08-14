# Match-Me Server

A Go-based backend server for the Match-Me application, built with Gin web framework and Ent ORM for database operations.

## Architecture Overview

This server implements a clean architecture pattern with distinct layers for API routing, business logic, data models, and database operations. The application uses PostgreSQL as its primary database with automatic schema migration via Ent ORM.

## Directory Structure

```
server/
  api/                    # HTTP API layer
  cmd/                    # Application entry points
  config/                 # Configuration management
  ent/                    # Ent ORM generated code and schemas
  internal/               # Private application code
  go.mod                  # Go module definition
  go.sum                  # Go module checksums
  README.md               # This documentation
```

## Core Components

### API Layer (`api/`)

**`server.go`**
- HTTP server initialization and configuration
- Gin router setup with middleware integration
- Server timeout and connection limits configuration
- Environment-based mode switching (development/production)

**`routes.go`**
- HTTP route definitions and handlers
- Currently minimal implementation (placeholder for future endpoints)

**`middleware/middlewares.go`**
- Custom middleware implementations
- **Ping Middleware**: Health check endpoint responding to `/ping` with "pong"

### Application Entry Point (`cmd/server/`)

**`main.go`**
- Application bootstrap and initialization
- Configuration loading via environment variables
- Database client setup with automatic migration
- HTTP server lifecycle management
- Graceful shutdown handling with configurable timeouts

**`lifecycle.go`**
- Server lifecycle management utilities
- Signal handling for graceful shutdown (SIGTERM, SIGINT)
- Background server monitoring and error handling
- Clean shutdown procedures with timeout controls

### Configuration (`config/`)

**`config.go`**
- Environment-based configuration loading using godotenv
- Singleton pattern implementation for configuration management
- Database connection and server address configuration

**`structs.go`**
- Configuration structure definitions
- Environment variable helper functions with validation
- Support for string, integer, and boolean environment variables
- Required vs optional configuration parameters

**Configuration Parameters:**
- `APP_ENV`: Application environment (development/production)
- `PORT`: HTTP server port (default: 8080)
- `HOST`: Server host address
- `DATABASE_URL`: PostgreSQL connection string
- `DATABASE_NAME`: Database name/driver
- `JWT_SECRET`: JWT signing secret
- `SERVER_ADDR`: Server address configuration
- `CLIENT_ADDR`: Client address configuration

### Database Schema (`ent/schema/`)

**`user.go`**
- User entity schema definition using Ent ORM
- **Fields:**
  - `id`: UUID primary key with auto-generation
  - `email`: Unique, validated email address
  - `password_hash`: Sensitive password storage
  - `first_name`: User's first name (max 50 chars)
  - `username`: Unique username (3-30 chars, alphanumeric)
  - `created_at`/`updated_at`: Timestamp tracking
  - `is_online`: Real-time presence indicator
  - `age`: User age (18-100 range validation)
  - `gender`: Gender specification
  - `looking_for`: JSON array of relationship types sought
  - `interests`: JSON array of user interests
  - `music_preferences`: JSON array of music preferences
  - `food_preferences`: JSON array of food preferences
  - `communication_style`: Communication preference
  - `prompts`: JSON array of profile prompts and responses
- **Relationships:** One-to-many relationship with user photos

**`user_photo.go`**
- User photo entity schema for profile images
- **Fields:**
  - `id`: UUID primary key
  - `photo_url`: Image URL storage
  - `order`: Photo display order (minimum 1)
  - `user_id`: Foreign key to user entity
- **Relationships:** Many-to-one relationship with users (cascade delete)

### Internal Models (`internal/models/`)
Contains core business models and validation logic for different Match-Me features.
Each file is self-descriptive by name (`user.go`, `chat.go`, `event.go`, etc.) and defines the corresponding domain model along with any necessary validation rules.
Validation helpers and shared types are also kept here to support entity integrity across the application.

### Repository Layer (`internal/repositories/`)

**`entclient.go`**
- Database client initialization and connection management
- Automatic schema migration on application startup
- PostgreSQL driver integration with connection string configuration
- Error handling for database connection failures

### Dependencies

The application uses the following key dependencies:

- **Web Framework**: Gin (github.com/gin-gonic/gin) - HTTP web framework
- **ORM**: Ent (entgo.io/ent) - Entity framework for Go
- **Database**: PostgreSQL (github.com/lib/pq) - PostgreSQL driver
- **Validation**: go-playground/validator - Struct validation
- **Configuration**: godotenv (github.com/joho/godotenv) - Environment variable loading
- **UUID**: Google UUID (github.com/google/uuid) - UUID generation


## Getting Started

1. Set up required environment variables in `.env` file
2. Ensure PostgreSQL database is running and accessible
3. Run `go mod tidy` to install dependencies
4. Execute `go run cmd/server/main.go` to start the server
5. Server will automatically create database schema on first run
6. Health check available at `GET /ping`

The server implements graceful shutdown and will handle SIGTERM/SIGINT signals properly, allowing for clean database connection closure and ongoing request completion.