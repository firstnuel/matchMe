# Match-Me

A full-stack dating/matching application built with React TypeScript frontend and Go backend, featuring real-time chat, user profiles, and location-based matching.

## Tech Stack

**Frontend:**
- React 19 with TypeScript
- Vite for build tooling
- React Router for navigation
- Bootstrap for styling
- Zustand for state management
- React Query for API data management
- WebSocket for real-time features

**Backend:**
- Go 1.24+ with Gin web framework
- Ent ORM for database operations
- PostgreSQL database with PostGIS for location features
- JWT authentication
- WebSocket for real-time chat
- Cloudinary for image storage

## Prerequisites

- Node.js 18+ and npm
- Go 1.24+
- PostgreSQL with PostGIS extension

## Setup

### Database Setup

**Install PostgreSQL with PostGIS:**
```bash
# macOS
brew install postgresql postgis

# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib postgis

# Create database
createdb matchme
psql -d matchme -c "CREATE EXTENSION postgis;"
```

### Application Setup

1. **Clone and setup environment:**
   ```bash
   git clone <repository-url>
   cd match-me
   ```

2. **Configure environment variables:**
   Create `server/.env` file:
   ```env
   APP_ENV=development
   PORT=8080
   DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/matchme?sslmode=disable
   CLIENT_ADDR=http://someclientdomain #for CORS allowed host config defaults to all if not set
   DATABASE_NAME=matchme
   JWT_SECRET=your-secret-key
   CLOUDINARY_URL=your-cloudinary-url
   ```

3. **Install dependencies and run:**
   ```bash
   # Install client dependencies
   make install-client

   # Run development servers (in separate terminals)
   make dev-client    # Frontend on http://localhost:5173
   make run-server    # Backend on http://localhost:8080
   ```

## Available Commands

```bash
make help              # Show all available commands
make dev-setup         # Setup development environment
make dev-client        # Start frontend development server
make run-server        # Start backend server
make build-all         # Build both client and server for production
make clean             # Clean build artifacts
```

## Database Management

The server includes built-in database management commands:

```bash
cd server
go run ./cmd/server -h         # Show help
go run ./cmd/server -r         # Reset database
go run ./cmd/server -p 50      # Add 50 test users
go run ./cmd/server -rp 25     # Reset and add 25 test users
```

## Project Structure

```
match-me/
├── client/           # React TypeScript frontend
│   ├── src/
│   │   ├── features/ # Feature-based organization
│   │   └── shared/   # Shared components and utilities
│   └── package.json
├── server/           # Go backend
│   ├── cmd/server/   # Main application entry
│   ├── api/          # HTTP routes and handlers
│   ├── ent/          # Database models and migrations
│   └── internal/     # Business logic and services
└── Makefile         # Build and development commands
```

## Features

- User authentication and profiles
- Photo upload with Cloudinary
- Location-based matching
- Real-time chat with WebSocket
- Connection requests and management
- Mobile-responsive design

