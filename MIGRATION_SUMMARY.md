# MongoDB Migration Summary

RestySched has been successfully migrated from SQLite to MongoDB!

## What Changed

### 1. Database Layer

**Before (SQLite):**
- File-based database (`restysched.db`)
- SQL queries
- Limited scalability

**After (MongoDB):**
- NoSQL database
- Document-based storage
- Highly scalable
- Cloud-ready (MongoDB Atlas)

### 2. Repository Implementation

**New Files:**
- `internal/repository/mongodb/db.go` - MongoDB connection and indexes
- `internal/repository/mongodb/employee_repository.go` - Employee data access
- `internal/repository/mongodb/schedule_repository.go` - Schedule data access

**Removed Files:**
- `internal/repository/sqlite/` - No longer needed

### 3. Configuration

**New Environment Variables:**
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=restysched
```

**Removed:**
```env
DATABASE_PATH=restysched.db
```

### 4. Domain Models

**Updated with BSON tags:**
```go
type Employee struct {
    ID    string `json:"id" bson:"id"`
    Name  string `json:"name" bson:"name"`
    // ... other fields
}
```

### 5. Dependencies

**Added:**
- `go.mongodb.org/mongo-driver v1.17.6`

**Removed:**
- `github.com/mattn/go-sqlite3`

## Data Structure

### employees Collection

```json
{
  "_id": ObjectId("..."),
  "id": "uuid-string",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "Developer",
  "role_description": "Full-stack developer",
  "monthly_hours": 160,
  "active": true,
  "created_at": ISODate("2024-01-01T00:00:00Z"),
  "updated_at": ISODate("2024-01-01T00:00:00Z")
}
```

**Indexes:**
- `email` (unique)
- `active`

### schedules Collection

```json
{
  "_id": ObjectId("..."),
  "id": "uuid-string",
  "period_start": ISODate("2024-01-01T00:00:00Z"),
  "period_end": ISODate("2024-01-15T00:00:00Z"),
  "employees": [
    {
      "id": "employee-uuid",
      "name": "John Doe",
      // ... full employee object
    }
  ],
  "status": "draft",
  "sent_to_n8n": false,
  "sent_at": null,
  "created_at": ISODate("2024-01-01T00:00:00Z"),
  "updated_at": ISODate("2024-01-01T00:00:00Z")
}
```

**Indexes:**
- `period_start`, `period_end` (compound)
- `status`

## Benefits of MongoDB

### 1. Scalability
- Horizontal scaling with sharding
- Replica sets for high availability
- Handle millions of documents

### 2. Flexibility
- Schema-less design
- Easy to add new fields
- Nested documents (employees in schedules)

### 3. Performance
- Automatic indexing
- Fast queries with proper indexes
- Connection pooling built-in

### 4. Cloud-Ready
- MongoDB Atlas integration
- Automatic backups
- Global distribution

### 5. Developer Experience
- JSON-like documents
- Rich query language
- Built-in aggregation framework

## Migration Steps (if you have existing data)

### From SQLite to MongoDB

If you have existing SQLite data, here's how to migrate:

1. **Export SQLite Data:**

```bash
# Export employees
sqlite3 restysched.db << 'EOF'
.mode json
.output employees.json
SELECT * FROM employees;
.quit
EOF

# Export schedules
sqlite3 restysched.db << 'EOF'
.mode json
.output schedules.json
SELECT * FROM schedules;
.quit
EOF
```

2. **Transform Data:**

Create a migration script (`migrate.go`):

```go
package main

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "log"

    "github.com/isak/restySched/internal/repository/mongodb"
)

func main() {
    // Connect to MongoDB
    db, err := mongodb.InitDB("mongodb://localhost:27017", "restysched")
    if err != nil {
        log.Fatal(err)
    }

    // Read SQLite JSON exports
    employeesData, _ := ioutil.ReadFile("employees.json")
    schedulesData, _ := ioutil.ReadFile("schedules.json")

    // Insert into MongoDB
    // ... migration logic here
}
```

3. **Run Migration:**

```bash
go run migrate.go
```

## Setup Instructions

### Quick Start

1. **Start MongoDB:**
   ```bash
   docker-compose up -d
   ```

2. **Update Configuration:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Run Application:**
   ```bash
   go run cmd/server/main.go
   ```

### Verify Setup

You should see:
```
Connected to MongoDB database: restysched
Server starting on port 8080
```

## Testing

All existing tests still pass:

```bash
go test ./...
```

The repository pattern ensures tests work exactly as before with mock implementations.

## Backward Compatibility

### Repository Pattern Preserved

The repository interfaces remain unchanged:

```go
type EmployeeRepository interface {
    Create(ctx context.Context, employee *Employee) error
    GetByID(ctx context.Context, id string) (*Employee, error)
    // ... same as before
}
```

Only the implementation changed from SQLite to MongoDB.

### Service Layer Unchanged

All business logic remains the same:

```go
// Still works exactly as before
employeeService := service.NewEmployeeService(employeeRepo)
employee, err := employeeService.CreateEmployee(ctx, input)
```

### Handlers Unchanged

HTTP handlers require no changes:

```go
// No changes needed
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
    // Same logic
}
```

## Performance Comparison

### SQLite
- Single-file database
- Good for small-scale
- Limited concurrent writes
- No horizontal scaling

### MongoDB
- Distributed database
- Excellent for production
- High concurrent writes
- Horizontal scaling with sharding

## What Stayed the Same

✅ Repository interfaces
✅ Service layer logic
✅ HTTP handlers
✅ Templates and UI
✅ n8n integration
✅ Scheduler logic
✅ Domain models (just added BSON tags)
✅ Testing approach
✅ API endpoints
✅ Configuration pattern

## What's New

🆕 MongoDB driver
🆕 BSON tags on models
🆕 MongoDB repository implementation
🆕 docker-compose.yml for easy setup
🆕 MongoDB-specific documentation
🆕 Automatic index creation

## Deployment Options

### Development
```bash
# Local MongoDB
docker-compose up -d
```

### Production Options

**Option 1: Self-Hosted**
- Install MongoDB on your server
- Configure replica sets
- Set up backups

**Option 2: MongoDB Atlas (Recommended)**
- Free tier available
- Automatic backups
- Global distribution
- Monitoring included

**Option 3: Cloud Providers**
- AWS DocumentDB
- Azure Cosmos DB (MongoDB API)
- Google Cloud MongoDB

## Next Steps

1. ✅ MongoDB is set up and working
2. ✅ All tests pass
3. ✅ Documentation updated
4. 🔜 Add employees through the UI
5. 🔜 Generate schedules
6. 🔜 Send to n8n

## Rollback Plan

If you need to go back to SQLite:

1. The old SQLite code is preserved in git history
2. Restore `internal/repository/sqlite/` files
3. Update `cmd/server/main.go` to use SQLite
4. Revert dependencies in `go.mod`

## Support

For MongoDB questions:
- See [MONGODB_SETUP.md](MONGODB_SETUP.md)
- Check [MongoDB Documentation](https://docs.mongodb.com/)

For RestySched questions:
- See [README.md](README.md)
- Check [QUICKSTART.md](QUICKSTART.md)

---

**Migration completed successfully! 🎉**

RestySched is now running on MongoDB with improved scalability and flexibility.
