# RestySched Architecture

## Overview

RestySched is a schedule planning automation system built with Go, following clean architecture principles and using the repository pattern for dependency injection and testability.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                      Web Layer                          │
│  (Templ Templates + HTMX + Tailwind CSS)               │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                  Handler Layer                          │
│   (HTTP Handlers - Employee, Schedule, Home)           │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                  Service Layer                          │
│   (Business Logic - Employee, Schedule Services)       │
└─────────┬──────────────────────────────────┬───────────┘
          │                                  │
          │                                  │
┌─────────▼────────────────┐    ┌───────────▼────────────┐
│   Repository Layer       │    │   External Services    │
│   (Data Access)          │    │   (n8n Client)         │
├──────────────────────────┤    └────────────────────────┘
│ - Employee Repository    │
│ - Schedule Repository    │
│                          │
│ Implementation:          │
│ - SQLite                 │
└──────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                  Domain Layer                           │
│  (Entities, Value Objects, Domain Errors)              │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                  Scheduler                              │
│  (Automated Biweekly Schedule Generation)              │
└─────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Core business entities and rules

**Components**:
- `employee.go`: Employee entity with validation
- `schedule.go`: Schedule entity and n8n payload structures
- `errors.go`: Domain-specific errors

**Rules**:
- No dependencies on other layers
- Contains pure business logic
- Defines interfaces for repositories (consumed by other layers)

### 2. Repository Layer (`internal/repository/`)

**Purpose**: Data access abstraction

**Components**:
- `employee_repository.go`: Employee data access interface
- `schedule_repository.go`: Schedule data access interface
- `sqlite/`: SQLite implementation
  - `employee_repository.go`: SQLite employee implementation
  - `schedule_repository.go`: SQLite schedule implementation
  - `db.go`: Database initialization and migrations

**Key Pattern**: Repository Pattern
- Interfaces defined in repository package
- Implementations in subdirectories (e.g., `sqlite/`)
- Easy to swap implementations or add new ones
- Enables dependency injection for testing

### 3. Service Layer (`internal/service/`)

**Purpose**: Business logic orchestration

**Components**:
- `employee_service.go`: Employee business logic
- `schedule_service.go`: Schedule generation and n8n integration
- `employee_service_test.go`: Test examples with mocks

**Responsibilities**:
- Coordinate between repositories
- Enforce business rules
- Handle transaction boundaries
- Transform data between layers

### 4. Handler Layer (`internal/handler/`)

**Purpose**: HTTP request handling

**Components**:
- `home_handler.go`: Home page
- `employee_handler.go`: Employee CRUD operations
- `schedule_handler.go`: Schedule generation and sending

**Responsibilities**:
- Parse HTTP requests
- Validate input
- Call service layer
- Render responses using Templ templates

### 5. Web Layer (`web/templates/`)

**Purpose**: User interface

**Components**:
- `layout.templ`: Base layout with navigation
- `home.templ`: Home page
- `employees.templ`: Employee list and forms
- `schedules.templ`: Schedule list and cards

**Technologies**:
- Templ: Type-safe Go templating
- HTMX: Dynamic interactions
- Tailwind CSS: Styling

### 6. External Services (`internal/n8n/`)

**Purpose**: Integration with external systems

**Components**:
- `client.go`: n8n webhook client

**Responsibilities**:
- Send schedule data to n8n webhook
- Handle HTTP communication
- Error handling and retries

### 7. Scheduler (`internal/scheduler/`)

**Purpose**: Automated task execution

**Components**:
- `scheduler.go`: Biweekly schedule automation

**Responsibilities**:
- Run tasks on schedule (every 2 weeks)
- Generate schedules automatically
- Send schedules to n8n
- Logging and error handling

### 8. Configuration (`internal/config/`)

**Purpose**: Application configuration

**Components**:
- `config.go`: Environment-based configuration

**Responsibilities**:
- Load environment variables
- Validate configuration
- Provide configuration to application

## Dependency Injection

The application uses constructor injection throughout:

```go
// Repositories
employeeRepo := sqlite.NewEmployeeRepository(db)
scheduleRepo := sqlite.NewScheduleRepository(db)

// External services
n8nClient := n8n.NewClient(webhookURL)

// Services (injecting repositories and clients)
employeeService := service.NewEmployeeService(employeeRepo)
scheduleService := service.NewScheduleService(
    scheduleRepo,
    employeeRepo,
    n8nClient,
)

// Handlers (injecting services)
employeeHandler := handler.NewEmployeeHandler(employeeService)
scheduleHandler := handler.NewScheduleHandler(scheduleService)

// Scheduler (injecting services)
scheduler := scheduler.NewScheduler(scheduleService)
```

## Testing Strategy

### Unit Testing with Mocks

Example from `employee_service_test.go`:

```go
// 1. Create mock repository
mockRepo := NewMockEmployeeRepository()

// 2. Inject into service
service := NewEmployeeService(mockRepo)

// 3. Test service logic
employee, err := service.CreateEmployee(ctx, input)
```

**Benefits**:
- Test business logic in isolation
- No database required
- Fast test execution
- Easy to simulate edge cases

### Integration Testing

For integration tests:
1. Use in-memory SQLite (`:memory:`)
2. Test full stack except HTTP layer
3. Verify database interactions

## Data Flow

### Creating an Employee

```
User Form → Handler → Service → Repository → Database
                ↓         ↓
            Validation  Business
                       Logic
```

### Automated Schedule Generation

```
Scheduler → Service → Get Active Employees → Repository
             ↓
        Create Schedule
             ↓
        Send to n8n → n8n Client → Webhook
             ↓
        Mark as Sent → Repository
```

## Key Design Patterns

### 1. Repository Pattern

**Purpose**: Abstract data access

```go
type EmployeeRepository interface {
    Create(ctx context.Context, employee *Employee) error
    GetByID(ctx context.Context, id string) (*Employee, error)
    // ...
}
```

### 2. Service Layer Pattern

**Purpose**: Encapsulate business logic

```go
type EmployeeService struct {
    repo repository.EmployeeRepository
}

func (s *EmployeeService) CreateEmployee(
    ctx context.Context,
    input EmployeeCreateInput,
) (*Employee, error) {
    // Business logic here
}
```

### 3. Dependency Injection

**Purpose**: Loose coupling, testability

All dependencies injected via constructors, never created internally.

### 4. Interface Segregation

Each repository interface focuses on a single entity with specific operations.

## Configuration Management

Environment variables:
- `SERVER_PORT`: HTTP server port
- `DATABASE_PATH`: SQLite database file
- `N8N_WEBHOOK_URL`: n8n webhook endpoint
- `ENABLE_SCHEDULER`: Enable/disable automation

## Error Handling

### Domain Errors

Defined in `internal/domain/errors.go`:
- `ErrEmployeeNotFound`
- `ErrScheduleNotFound`
- `ErrInvalidEmployeeEmail`
- etc.

### Error Flow

```
Repository → Service → Handler → HTTP Response
    ↓           ↓          ↓
Domain Err  Check Err   Status Code
```

## Database Schema

### Employees Table

```sql
CREATE TABLE employees (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    role_description TEXT NOT NULL,
    monthly_hours INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### Schedules Table

```sql
CREATE TABLE schedules (
    id TEXT PRIMARY KEY,
    period_start DATETIME NOT NULL,
    period_end DATETIME NOT NULL,
    employees TEXT NOT NULL,  -- JSON array
    status TEXT NOT NULL,
    sent_to_n8n BOOLEAN NOT NULL DEFAULT 0,
    sent_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

## Future Enhancements

Potential improvements maintaining current architecture:

1. **Add PostgreSQL Repository**
   - Implement `postgresql.EmployeeRepository`
   - No changes to service or handler layers

2. **Add API Layer**
   - Create REST API handlers
   - Reuse existing service layer

3. **Add Caching**
   - Implement caching repository decorator
   - Wrap existing repositories

4. **Add Authentication**
   - Add middleware layer
   - No changes to core business logic

5. **Multiple Schedule Types**
   - Extend domain models
   - Add new service methods
   - Update templates

The clean architecture makes these enhancements straightforward without major refactoring.
