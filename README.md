# RestySched - Automated Schedule Planning System

A Go-based schedule planning automation system that manages employee schedules and integrates with n8n workflows for automated processing.

## Features

- **Employee Management**: Create, update, and manage employees with roles and descriptions
- **Automated Schedule Generation**: Automatically generates biweekly schedules
- **n8n Integration** (Optional): Sends schedule data to n8n webhooks for workflow automation
- **Repository Pattern**: Clean architecture with dependency injection for easy testing
- **Templ Templates**: Modern Go templating with HTMX for dynamic UI
- **MongoDB Database**: Scalable NoSQL database for data persistence
- **Biweekly Scheduler**: Automated schedule generation every 2 weeks

## Architecture

The project follows clean architecture principles:

```
restySched/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and errors
│   ├── handler/         # HTTP handlers
│   ├── n8n/             # n8n webhook client
│   ├── repository/      # Repository interfaces and implementations
│   │   └── mongodb/     # MongoDB implementation
│   ├── scheduler/       # Biweekly schedule automation
│   └── service/         # Business logic layer
└── web/
    └── templates/       # Templ templates
```

## Prerequisites

- Go 1.23 or higher
- MongoDB 4.4 or higher (running locally or remote)
- Templ CLI (for template generation)
- n8n instance with webhook configured (optional - can be added later)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd restySched
```

2. Install dependencies:
```bash
go mod download
```

3. Install Templ CLI:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

4. Generate Templ templates:
```bash
templ generate
```

5. Create a `.env` file from the example:
```bash
cp .env.example .env
```

6. Set up MongoDB (choose one):

**Option A: MongoDB Atlas (Cloud - Recommended)**
- Free tier available
- See [MONGODB_ATLAS_SETUP.md](MONGODB_ATLAS_SETUP.md) for complete guide
- Connection string: `mongodb+srv://user:pass@cluster.mongodb.net/`

**Option B: Local with Docker (Easiest for local dev)**
```bash
docker-compose up -d
```

**Option C: Local MongoDB Installation**
```bash
mongod
```

7. Update the `.env` file with your configuration:

**For MongoDB Atlas:**
```env
SERVER_PORT=8080
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DATABASE=restysched
N8N_WEBHOOK_URL=  # Optional - leave empty to run without n8n
ENABLE_SCHEDULER=true
```

**For Local MongoDB:**
```env
SERVER_PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=restysched
N8N_WEBHOOK_URL=  # Optional - leave empty to run without n8n
ENABLE_SCHEDULER=true
```

**Note:** n8n integration is optional. You can start using RestySched without configuring n8n and add it later when ready.

## Running the Application

1. Generate templates (required after any template changes):
```bash
templ generate
```

2. Run the server:
```bash
go run cmd/server/main.go
```

3. Access the application at `http://localhost:8080`

## Building for Production

```bash
# Generate templates
templ generate

# Build the binary
go build -o restysched cmd/server/main.go

# Run the binary
./restysched
```

## Testing

The project includes example tests demonstrating the repository pattern with dependency injection:

```bash
go test ./...
```

Example test with mock repository:
```bash
go test ./internal/service -v
```

## Usage

### Managing Employees

1. Navigate to `/employees`
2. Click "Add Employee" to create a new employee
3. Fill in the form:
   - **Name**: Employee's full name
   - **Email**: Contact email
   - **Role**: Job title (e.g., Developer, Designer)
   - **Role Description**: Detailed description of responsibilities
   - **Monthly Hours**: Required hours per month

### Generating Schedules

1. Navigate to `/schedules`
2. Click "Generate Biweekly Schedule" to create a new schedule
3. Review the generated schedule with all active employees
4. Click "Send to n8n" to trigger the workflow

### Automated Schedule Generation

When `ENABLE_SCHEDULER=true`, the system automatically:
- Generates a new schedule every 2 weeks
- Includes all active employees
- Sends the schedule to n8n webhook
- Logs all operations

## n8n Webhook Integration

### Payload Format

The system sends the following JSON structure to your n8n webhook:

```json
{
  "schedule_id": "uuid-string",
  "period_start": "2024-01-01T00:00:00Z",
  "period_end": "2024-01-15T00:00:00Z",
  "employees": [
    {
      "id": "employee-uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "Developer",
      "role_description": "Full-stack developer working on web applications",
      "monthly_hours": 160
    }
  ],
  "generated_at": "2024-01-01T00:00:00Z"
}
```

### Setting up n8n Webhook

1. In n8n, create a new workflow
2. Add a "Webhook" trigger node
3. Set the HTTP method to POST
4. Copy the webhook URL
5. Add it to your `.env` file as `N8N_WEBHOOK_URL`

## API Endpoints

### Web UI
- `GET /` - Home page
- `GET /employees` - Employee list
- `GET /schedules` - Schedule list

### Employee API
- `POST /employees` - Create employee
- `PUT /employees/{id}` - Update employee
- `DELETE /employees/{id}` - Delete employee (soft delete)
- `GET /employees/new` - New employee form
- `GET /employees/{id}/edit` - Edit employee form

### Schedule API
- `POST /schedules/generate` - Generate new biweekly schedule
- `POST /schedules/{id}/send` - Send schedule to n8n
- `DELETE /schedules/{id}` - Delete schedule

## Dependency Injection Example

The repository pattern allows easy testing with mock implementations:

```go
// Production code
employeeRepo := mongodb.NewEmployeeRepository(db)
service := service.NewEmployeeService(employeeRepo)

// Test code
mockRepo := NewMockEmployeeRepository()
service := service.NewEmployeeService(mockRepo)
```

## Configuration

All configuration is managed through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | 8080 |
| `MONGO_URI` | MongoDB connection URI (required) | mongodb://localhost:27017 |
| `MONGO_DATABASE` | MongoDB database name (required) | restysched |
| `N8N_WEBHOOK_URL` | n8n webhook URL (optional) | empty |
| `ENABLE_SCHEDULER` | Enable automated scheduling | true |

## MongoDB Collections

### employees Collection

```json
{
  "id": "uuid-string",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "Developer",
  "role_description": "Full-stack developer",
  "monthly_hours": 160,
  "active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Indexes:**
- `email` (unique)
- `active`

### schedules Collection

```json
{
  "id": "uuid-string",
  "period_start": "2024-01-01T00:00:00Z",
  "period_end": "2024-01-15T00:00:00Z",
  "employees": [
    {
      "id": "employee-uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "Developer",
      "role_description": "Full-stack developer",
      "monthly_hours": 160,
      "active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "status": "draft",
  "sent_to_n8n": false,
  "sent_at": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Indexes:**
- `period_start`, `period_end` (compound)
- `status`

## Tech Stack

- **Go 1.23**: Programming language
- **Templ**: Type-safe Go templating
- **HTMX**: Dynamic UI interactions
- **Tailwind CSS**: Styling
- **MongoDB**: NoSQL database
- **gocron**: Scheduled task execution
- **godotenv**: Environment configuration

## Development

### Adding New Features

1. Define domain models in `internal/domain/`
2. Create repository interface in `internal/repository/`
3. Implement repository in `internal/repository/mongodb/`
4. Add business logic in `internal/service/`
5. Create handlers in `internal/handler/`
6. Build templates in `web/templates/`
7. Generate templates with `templ generate`

### Testing Strategy

Use the repository pattern for dependency injection:

1. Define interfaces in repository layer
2. Create mock implementations for testing
3. Inject dependencies through constructors
4. Test services with mock repositories

## License

MIT

## Contributing

Contributions are welcome! Please ensure:
- Tests pass: `go test ./...`
- Code is formatted: `go fmt ./...`
- Templates are generated: `templ generate`
