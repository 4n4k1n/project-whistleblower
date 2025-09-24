# 42 Whistleblower System

A web application for reporting suspected academic dishonesty at 42 school. This system allows students to report suspicious projects while providing staff tools to review and manage reports.

## Features

### For Students
- **42 OAuth Authentication**: Secure login using 42 Intra credentials
- **Student Search**: Find students by login name
- **Project Selection**: View and select from a student's projects
- **Report Submission**: Submit reports with predefined reasons and detailed explanations
- **Abuse Prevention**: System tracks false report patterns

### For Staff
- **Report Review**: View and approve/reject pending reports
- **Automatic Notifications**: Staff alerted when report threshold reached (3+ reports per project)
- **Abuse Monitoring**: Track users with high false report ratios

### Security Features
- **Threshold System**: Staff notified only after multiple reports (default: 3)
- **False Report Tracking**: Users with high rejection rates are flagged
- **Authentication Required**: All actions require 42 OAuth authentication
- **Audit Trail**: All actions are logged with timestamps and user IDs

## Setup

### Prerequisites
- Go 1.21 or higher
- SQLite3
- 42 API Application (for OAuth)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd whistleblower
```

2. Install dependencies:
```bash
make deps
```

3. Create environment configuration:
```bash
cp .env.example .env
```

4. Configure your 42 API application:
   - Go to https://profile.intra.42.fr/oauth/applications
   - Create a new application
   - Set redirect URI to: `http://localhost:8080/callback`
   - Update `.env` with your client credentials

5. Initialize the database:
```bash
make init-db
```

### Running the Application

Development mode:
```bash
make dev
```

Production mode:
```bash
make build
./whistleblower
```

The application will be available at `http://localhost:8080`

## Usage

### Student Workflow
1. Visit the application and click "Login with 42 Intra"
2. Authenticate with your 42 credentials
3. Search for the student you want to report
4. Select the specific project
5. Choose a reason from the dropdown
6. Provide a detailed explanation
7. Submit the report

### Staff Workflow
1. Access `/api/staff/reports` to view pending reports
2. Review report details and evidence
3. Approve or reject reports via `/api/staff/reports/:id`

## API Endpoints

### Public Endpoints
- `GET /` - Landing page
- `GET /login` - Initiate 42 OAuth flow
- `GET /callback` - OAuth callback
- `GET /dashboard` - Main dashboard

### Authenticated Endpoints
- `GET /api/students/search?q=<query>` - Search students
- `GET /api/students/:login/projects` - Get student's projects
- `POST /api/reports` - Submit a report
- `GET /api/report-reasons` - Get available report reasons

### Staff-Only Endpoints
- `GET /api/staff/reports` - Get pending reports
- `PUT /api/staff/reports/:id` - Review a report

## Database Schema

The system uses SQLite with the following main tables:
- `users` - User accounts from 42 OAuth
- `reports` - Submitted reports with status tracking
- `staff_notifications` - Notifications sent to staff
- `user_report_stats` - False report tracking
- `report_reasons` - Predefined report categories

## Environment Variables

- `OAUTH_42_CLIENT_ID` - Your 42 application client ID
- `OAUTH_42_CLIENT_SECRET` - Your 42 application client secret  
- `OAUTH_42_REDIRECT_URL` - OAuth callback URL
- `PORT` - Server port (default: 8080)

## Abuse Prevention

The system includes several mechanisms to prevent abuse:

1. **Report Threshold**: Staff only notified after 3+ reports per project
2. **False Report Tracking**: Users with >70% rejection rate are flagged
3. **Authentication Required**: All actions require valid 42 credentials
4. **Detailed Logging**: All reports and reviews are logged with user IDs
5. **Staff Review**: All reports require human review before action

## Development

Run tests:
```bash
make test
```

Clean build artifacts:
```bash
make clean
```

## License

This project is for educational and administrative use at 42 school only.