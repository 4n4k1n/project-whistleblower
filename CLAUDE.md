# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based web application for reporting suspected academic dishonesty at 42 school. The system provides secure reporting for students and review tools for staff, with built-in abuse prevention mechanisms.

## Architecture

### Core Components
- **main.go**: Entry point with Gin router setup and middleware
- **handlers/**: HTTP request handlers for web and API endpoints
- **auth/**: 42 OAuth integration for authentication
- **database/**: SQLite database layer with schema management
- **models/**: Data structures for users, reports, and API requests
- **templates/**: HTML templates for the web interface

### Key Features
- 42 OAuth authentication with session management
- Report threshold system (staff notified after 3+ reports per project)
- False report tracking and user abuse prevention
- Staff-only admin interface for report review
- RESTful API with protected endpoints

## Common Commands

### Development
```bash
make dev              # Run with race detection
make run              # Standard run
make build            # Build binary
make deps             # Download and tidy dependencies
```

### Database Operations
```bash
make init-db          # Initialize SQLite database with schema
make clean            # Remove binary and database files
```

### Testing
```bash
make test             # Run Go tests
./test.sh             # Integration test suite with server startup
./test_*.sh           # Specific functionality tests
```

### Docker Operations
```bash
make docker-build     # Build Docker image
make docker-run       # Run in container with environment file
```

## Database Architecture

The system uses SQLite with the following core tables:
- **users**: 42 OAuth user accounts with staff privileges
- **reports**: Academic dishonesty reports with status tracking
- **staff_notifications**: Threshold-based notifications to staff
- **user_report_stats**: False report ratio tracking for abuse prevention
- **report_reasons**: Predefined report categories

Database initialization automatically populates report reasons and creates all required tables.

## Environment Configuration

Required environment variables:
- `OAUTH_42_CLIENT_ID`: 42 API application client ID
- `OAUTH_42_CLIENT_SECRET`: 42 API application secret
- `OAUTH_42_REDIRECT_URL`: OAuth callback URL (typically `http://localhost:8080/callback`)
- `PORT`: Server port (defaults to 8080)

## API Structure

### Authentication Flow
1. `/login` - Initiates 42 OAuth with state verification
2. `/callback` - Handles OAuth callback and creates/updates user session
3. Session stored in cookies for subsequent requests

### Public Endpoints
- `GET /` - Landing page
- `GET /dashboard` - Main authenticated dashboard
- `GET /api/report-reasons` - Available report categories

### Authenticated API
- `GET /api/students/search` - Search 42 students by login
- `GET /api/students/:login/projects` - Fetch student's project list
- `POST /api/reports` - Submit new report
- `GET /api/stats` - User's report statistics

### Staff-Only Endpoints
- `GET /admin` - Admin interface (staff permission required)
- `GET /api/staff/reports` - Pending reports for review
- `PUT /api/staff/reports/:id` - Approve/reject reports

## Security Features

### Abuse Prevention
- Report threshold system: Staff only notified after 3+ reports per project
- False report tracking: Users with >70% rejection rate are flagged
- Session-based authentication with 42 OAuth only
- CSRF protection via state parameter in OAuth flow

### Access Control
- Staff privileges stored in database `is_staff` field
- Protected endpoints check user authentication status
- Admin page restricted to staff users only

## Development Notes

- Uses Gin framework with HTML template rendering
- SQLite database with automatic schema initialization
- 42 API integration for student/project data fetching
- Thread-safe request handling with Gin's built-in features
- Comprehensive test suite covering endpoints and database operations

## Testing Strategy

The application includes both unit tests (`go test`) and integration tests via shell scripts:
- Server startup and endpoint availability
- Database schema and data integrity
- OAuth flow simulation
- API response validation
- Admin access control verification