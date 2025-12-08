# Using RestySched Without n8n

RestySched can be used as a standalone schedule management system without n8n integration. This guide shows you how to use all features without setting up n8n.

## Quick Start (No n8n Required)

### 1. Setup

```bash
# Start MongoDB
docker-compose up -d

# Copy environment file
cp .env.example .env
```

### 2. Configure `.env`

```env
SERVER_PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=restysched
N8N_WEBHOOK_URL=  # Leave empty - no n8n needed!
ENABLE_SCHEDULER=true
```

### 3. Run the Application

```bash
go run cmd/server/main.go
```

### 4. Access the Web UI

Open your browser: http://localhost:8080

## Features Available Without n8n

✅ **Full Employee Management**
- Add, edit, and delete employees
- Set roles and descriptions
- Configure monthly hours
- Active/inactive status

✅ **Schedule Generation**
- Manual schedule generation
- Automatic biweekly generation
- View all generated schedules
- See employee details per schedule

✅ **Web Interface**
- Modern, responsive UI
- HTMX-powered interactions
- Real-time updates
- No page reloads needed

## What Happens Without n8n?

### When You Click "Send to n8n"

If n8n is not configured, you'll see an error message:

```
n8n webhook URL not configured - please set N8N_WEBHOOK_URL in .env file
```

The schedule is **still saved** in your database, you just won't send it to n8n.

### Automated Scheduler

If `ENABLE_SCHEDULER=true`, the system will:

1. ✅ Generate schedules every 2 weeks
2. ✅ Save them to MongoDB
3. ⚠️ Log a warning about n8n not being configured
4. ✅ Continue running normally

You'll see logs like:

```
Schedule generated successfully: abc123...
WARNING: Failed to send schedule to n8n: n8n webhook URL not configured
Schedule saved but not sent to n8n. You can manually send it later from the UI.
```

## Exporting Schedule Data

Without n8n, you can still export your data:

### 1. View in MongoDB

```bash
mongosh

use restysched

# View all schedules
db.schedules.find().pretty()

# View all employees
db.employees.find().pretty()
```

### 2. Export to JSON

```bash
# Export schedules
mongoexport --uri="mongodb://localhost:27017" \
  --db=restysched \
  --collection=schedules \
  --out=schedules.json

# Export employees
mongoexport --uri="mongodb://localhost:27017" \
  --db=restysched \
  --collection=employees \
  --out=employees.json
```

### 3. Use the Data Yourself

The exported JSON can be:
- Imported into Excel/Google Sheets
- Processed with custom scripts
- Sent to other systems
- Backed up for archival

## Alternative Integrations

Instead of n8n, you can integrate with other systems:

### 1. Direct API Access

The schedule data is stored in MongoDB and accessible via the web UI. You can:

- Read from MongoDB directly
- Build custom scripts to process schedules
- Create your own webhook endpoint

### 2. MongoDB Change Streams

Monitor for new schedules in real-time:

```javascript
const changeStream = db.schedules.watch();

changeStream.on('change', (change) => {
  if (change.operationType === 'insert') {
    console.log('New schedule created:', change.fullDocument);
    // Send to your system
  }
});
```

### 3. Custom Export Script

Create a Go script to export schedules:

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "os"

    "github.com/isak/restySched/internal/repository/mongodb"
)

func main() {
    db, _ := mongodb.InitDB("mongodb://localhost:27017", "restysched")
    scheduleRepo := mongodb.NewScheduleRepository(db)

    schedules, _ := scheduleRepo.GetAll(context.Background())

    data, _ := json.MarshalIndent(schedules, "", "  ")
    os.WriteFile("export.json", data, 0644)

    log.Printf("Exported %d schedules", len(schedules))
}
```

## Manual Workflow

Here's how to use RestySched in a manual workflow:

### Weekly Process

1. **Monday Morning:**
   - Open http://localhost:8080/schedules
   - Click "Generate Biweekly Schedule"
   - Review the employee list

2. **Review:**
   - Check all employees are included
   - Verify roles and hours are correct
   - Note the schedule ID

3. **Distribution:**
   - Export to JSON or view in MongoDB
   - Copy data to Excel/Google Sheets
   - Email to team manually
   - Post to Slack/Teams

4. **Tracking:**
   - Mark schedule as "completed" in your tracking system
   - Keep record in MongoDB for history

## Adding n8n Later

When you're ready to add n8n automation:

1. **Set up n8n:**
   - Install n8n locally or use n8n Cloud
   - Create a webhook workflow
   - Copy the webhook URL

2. **Update `.env`:**
   ```env
   N8N_WEBHOOK_URL=https://your-n8n.com/webhook/abc123
   ```

3. **Restart the application:**
   ```bash
   go run cmd/server/main.go
   ```

4. **Test:**
   - Generate a new schedule
   - Click "Send to n8n"
   - Should work immediately!

No code changes needed - just configuration.

## Benefits of Running Without n8n

### 1. Simpler Setup
- One less system to configure
- Faster initial deployment
- Easier troubleshooting

### 2. Lower Complexity
- Direct database access
- Clear data ownership
- No external dependencies

### 3. Full Control
- Choose your own integration
- Custom export formats
- Manual review process

### 4. Gradual Migration
- Start simple, add automation later
- Test workflow before automating
- Validate data before integration

## Common Use Cases

### Case 1: Small Team

**Scenario:** 5-10 employees, manual scheduling

**Workflow:**
1. Add employees once
2. Generate schedules biweekly
3. Export to Excel
4. Email to team

**No n8n needed** - manual process is fast enough

### Case 2: Testing Phase

**Scenario:** Evaluating the system

**Workflow:**
1. Set up without n8n
2. Test employee management
3. Generate test schedules
4. Review data in MongoDB
5. Add n8n when ready to automate

**Perfect for proof-of-concept**

### Case 3: Custom Integration

**Scenario:** Integrating with existing HR system

**Workflow:**
1. Use RestySched for schedule generation
2. Read from MongoDB directly
3. Custom script pushes to HR system
4. Skip n8n entirely

**Direct integration is simpler**

### Case 4: Offline Environment

**Scenario:** Air-gapped or restricted network

**Workflow:**
1. Run everything locally
2. No external webhooks
3. Export data manually
4. Transfer via approved methods

**n8n cloud not an option**

## Support

For questions about running without n8n:
- Check [README.md](README.md) for general setup
- See [MONGODB_SETUP.md](MONGODB_SETUP.md) for database help
- Review [QUICKSTART.md](QUICKSTART.md) for basic usage

For n8n integration when ready:
- See [N8N_SETUP.md](N8N_SETUP.md) for detailed instructions

## Conclusion

RestySched is **fully functional without n8n**. Use it as:
- A standalone schedule management tool
- A database for employee scheduling
- A foundation for custom integrations
- A manual workflow system

Add n8n automation whenever you're ready - it's completely optional! 🎉
