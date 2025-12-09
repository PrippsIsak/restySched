# RestySched - Quick Start Guide

Get up and running with RestySched in 5 minutes!

## Prerequisites

- Go 1.23 or higher
- MongoDB 4.4 or higher
- An n8n instance with webhook access

## Installation

### Option 1: Automated Setup (Recommended)

**Windows:**
```bash
setup.bat
```

**Linux/Mac:**
```bash
chmod +x setup.sh
./setup.sh
```

### Option 2: Manual Setup

1. **Install Templ CLI:**
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

2. **Download dependencies:**
```bash
go mod download
```

3. **Generate templates:**
```bash
templ generate
```

4. **Set up MongoDB (choose one):**

**Option A: MongoDB Atlas (Cloud)**
- [Follow Atlas setup guide](MONGODB_ATLAS_SETUP.md)
- Free tier available
- Recommended for production

**Option B: Local with Docker**
```bash
docker-compose up -d
```

**Option C: Local Installation**
```bash
mongod
```

5. **Create configuration:**
```bash
cp .env.example .env
```

6. **Update `.env` with your configuration:**

**For MongoDB Atlas:**
```env
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DATABASE=restysched
N8N_WEBHOOK_URL=  # Optional
```

**For Local MongoDB:**
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=restysched
N8N_WEBHOOK_URL=  # Optional
```

**Note:** n8n is optional and can be added later.

## Running the Application

### Using Make (Recommended)

```bash
# Run the application
make run

# Build for production
make build

# Run tests
make test

# Generate templates
make generate

# View all commands
make help
```

### Using Go Commands

```bash
# Run the application
go run cmd/server/main.go

# Build the application
go build -o restysched cmd/server/main.go

# Run tests
go test ./...
```

## First Steps

### 1. Access the Application

Open your browser and navigate to:
```
http://localhost:8080
```

### 2. Add Your First Employee

1. Click "Go to Employees"
2. Click "Add Employee"
3. Fill in the form:
   - **Name**: John Doe
   - **Email**: john@example.com
   - **Role**: Developer
   - **Role Description**: Full-stack developer working on web applications
   - **Monthly Hours**: 160
4. Click "Create"

### 3. Add More Employees

Repeat step 2 to add more employees. Example:

- **Name**: Jane Smith
- **Email**: jane@example.com
- **Role**: Designer
- **Role Description**: UI/UX designer creating beautiful interfaces
- **Monthly Hours**: 120

### 4. Generate Your First Schedule

1. Click "Go to Schedules"
2. Click "Generate Biweekly Schedule"
3. Review the generated schedule with all active employees
4. Click "Send to n8n" to trigger your workflow

### 5. Set Up n8n Workflow

See [N8N_SETUP.md](N8N_SETUP.md) for detailed instructions on creating your n8n workflow.

Basic steps:
1. Create new workflow in n8n
2. Add Webhook trigger (POST)
3. Add your processing nodes (Email, Slack, etc.)
4. Copy webhook URL to `.env` file
5. Activate workflow

## Configuration Options

Edit `.env` to customize:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DATABASE_PATH=restysched.db

# n8n Webhook Configuration
N8N_WEBHOOK_URL=https://your-n8n-instance.com/webhook/your-id

# Scheduler Configuration (true/false)
ENABLE_SCHEDULER=true
```

### Disable Automatic Scheduling

If you want to generate schedules manually only:
```env
ENABLE_SCHEDULER=false
```

## Common Tasks

### Add an Employee
1. Navigate to `/employees`
2. Click "Add Employee"
3. Fill in the form
4. Submit

### Edit an Employee
1. Navigate to `/employees`
2. Click "Edit" next to the employee
3. Update fields
4. Submit

### Delete an Employee
1. Navigate to `/employees`
2. Click "Delete" next to the employee
3. Confirm deletion
4. Employee is soft-deleted (marked as inactive)

### Generate a Schedule Manually
1. Navigate to `/schedules`
2. Click "Generate Biweekly Schedule"
3. Review the schedule

### Send a Schedule to n8n
1. Navigate to `/schedules`
2. Find the schedule you want to send
3. Click "Send to n8n"
4. Schedule is marked as "Sent"

### View Schedule History
1. Navigate to `/schedules`
2. All generated schedules are listed
3. View employee details for each schedule

## Automated Schedule Generation

When `ENABLE_SCHEDULER=true`, the system will:
- **Automatically** generate a new schedule every 2 weeks
- **Include** all active employees
- **Send** the schedule to your n8n webhook
- **Log** all operations to console

You'll see log messages like:
```
Scheduler started - will generate schedules every 2 weeks
Starting biweekly schedule generation...
Schedule generated successfully: abc123...
Schedule sent to n8n successfully: abc123...
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Service Tests with Details
```bash
go test ./internal/service -v
```

### Run with Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Project Structure

```
restySched/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Business entities
│   ├── handler/         # HTTP handlers
│   ├── n8n/            # n8n integration
│   ├── repository/      # Data access layer
│   │   └── mongodb/    # MongoDB implementation
│   ├── scheduler/       # Automated scheduling
│   └── service/         # Business logic
├── web/templates/       # UI templates
├── .env                # Configuration (create from .env.example)
├── go.mod              # Go dependencies
├── Makefile            # Build commands
└── README.md           # Full documentation
```

## Troubleshooting

### Port Already in Use

If port 8080 is already in use, change it in `.env`:
```env
SERVER_PORT=3000
```

### n8n Webhook Not Working

1. Check webhook URL in `.env`
2. Verify n8n workflow is activated
3. Test webhook with curl:
```bash
curl -X POST YOUR_WEBHOOK_URL \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```

### Templates Not Found Error

Generate templates:
```bash
templ generate
```

### MongoDB Connection Error

1. Verify MongoDB is running:
```bash
docker ps  # Check if MongoDB container is running
```

2. Check connection string in `.env`
3. For Atlas: Verify IP whitelist and credentials

### Build Errors

1. Update dependencies:
```bash
go mod tidy
```

2. Clean and rebuild:
```bash
make clean
make build
```

## Development Workflow

### Adding New Features

1. **Define domain models** in `internal/domain/`
2. **Create repository interface** in `internal/repository/`
3. **Implement repository** in `internal/repository/mongodb/`
4. **Add business logic** in `internal/service/`
5. **Create handlers** in `internal/handler/`
6. **Build templates** in `web/templates/`
7. **Generate templates**: `templ generate`
8. **Write tests** in `*_test.go` files
9. **Update routes** in `cmd/server/main.go`

### Testing Strategy

1. **Mock repositories** for service tests
2. **Test business logic** in isolation
3. **Use dependency injection** for flexibility

Example:
```go
// Create mock
mockRepo := NewMockEmployeeRepository()

// Inject into service
service := service.NewEmployeeService(mockRepo)

// Test
employee, err := service.CreateEmployee(ctx, input)
```

## Next Steps

1. ✅ Add employees
2. ✅ Generate schedules
3. ✅ Set up n8n workflow
4. ✅ Test the integration
5. 📚 Read [ARCHITECTURE.md](ARCHITECTURE.md) for deeper understanding
6. 🔧 Customize for your needs

## Production Deployment

### Build for Production

```bash
make build
```

This creates `bin/restysched` executable.

### Run in Production

```bash
# Set environment variables
export N8N_WEBHOOK_URL=https://your-n8n.com/webhook/id
export ENABLE_SCHEDULER=true

# Run the binary
./bin/restysched
```

### Using Docker (Optional)

Create `Dockerfile`:
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go build -o restysched cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/restysched .
COPY .env .
CMD ["./restysched"]
```

Build and run:
```bash
docker build -t restysched .
docker run -p 8080:8080 --env-file .env restysched
```

## Getting Help

- 📖 Read the [README.md](README.md) for full documentation
- 🏗️ Check [ARCHITECTURE.md](ARCHITECTURE.md) for design details
- 🔗 See [N8N_SETUP.md](N8N_SETUP.md) for n8n integration help
- 💬 Open an issue on GitHub for support

## What's Included

✅ Employee management CRUD
✅ Schedule generation
✅ n8n webhook integration (optional)
✅ Automated biweekly scheduling
✅ MongoDB database (local or Atlas)
✅ Repository pattern with DI
✅ Comprehensive tests
✅ Clean architecture
✅ HTMX-powered UI
✅ Responsive design

## License

MIT License - See LICENSE file for details

---

**Happy Scheduling! 🎉**
