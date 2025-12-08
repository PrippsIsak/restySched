# n8n Workflow Diagrams

Visual representations of the RestySched n8n workflows.

## Simple Starter Workflow

```
┌──────────────────┐
│                  │
│  RestySched App  │
│                  │
└────────┬─────────┘
         │
         │ HTTP POST
         │ /webhook/restysched
         │
         ▼
┌────────────────────┐
│                    │
│  Webhook Trigger   │
│  Receives schedule │
│                    │
└────┬──────────┬────┘
     │          │
     │          └──────────────────┐
     │                             │
     ▼                             ▼
┌─────────────────┐      ┌─────────────────────┐
│                 │      │                     │
│  Send Response  │      │  Send to Slack      │
│  {success:true} │      │  (Optional)         │
│                 │      │                     │
└─────────────────┘      └─────────────────────┘
     │
     │
     ▼
┌──────────────────┐
│                  │
│ Split Employees  │
│ Into Items       │
│                  │
└────────┬─────────┘
         │
         │ For Each Employee
         │
         ▼
┌──────────────────┐
│                  │
│  Send Email      │
│  (Optional)      │
│                  │
└──────────────────┘
```

## Full Automation Workflow

```
┌──────────────────┐
│  RestySched App  │
└────────┬─────────┘
         │
         │ HTTP POST
         │
         ▼
┌────────────────────┐
│  Webhook Trigger   │
└────┬──────────┬────┘
     │          │
     ▼          └──────────────────┐
┌──────────────┐                   │
│   Response   │                   │
└──────────────┘                   │
                                   ▼
                          ┌─────────────────┐
                          │ Process Data    │
                          │ Format dates    │
                          │ Create summary  │
                          └────┬───────┬────┘
                               │       │
              ┌────────────────┘       └────────────────┐
              │                                         │
              ▼                                         ▼
    ┌─────────────────┐                    ┌──────────────────┐
    │ Create Summary  │                    │ Split Employees  │
    │ Build message   │                    │ Into Items       │
    └────────┬────────┘                    └────────┬─────────┘
             │                                       │
             ▼                                       │ For Each
    ┌─────────────────┐                            │ Employee
    │ Send to Slack   │                            │
    │ Team summary    │                            ▼
    └─────────────────┘              ┌──────────────────────────┐
                                     │                          │
                         ┌───────────┼──────────┬───────────────┼──────────┐
                         │           │          │               │          │
                         ▼           ▼          ▼               ▼          ▼
                   ┌──────────┐ ┌────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐
                   │  Email   │ │ Sheets │ │ Calendar │ │  Check   │ │  ...   │
                   │  Send    │ │  Log   │ │  Event   │ │  Hours   │ │        │
                   └──────────┘ └────────┘ └──────────┘ └────┬─────┘ └────────┘
                                                              │
                                                              │ If > 120h
                                                              │
                                                              ▼
                                                        ┌──────────┐
                                                        │  Alert   │
                                                        │  Slack   │
                                                        └──────────┘
```

## Data Flow

### Input (from RestySched)

```json
{
  "schedule_id": "abc-123",
  "period_start": "2024-01-01T00:00:00Z",
  "period_end": "2024-01-15T00:00:00Z",
  "employees": [
    {
      "id": "emp-1",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "Developer",
      "role_description": "Full-stack developer",
      "monthly_hours": 160
    },
    {
      "id": "emp-2",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "role": "Designer",
      "role_description": "UI/UX Designer",
      "monthly_hours": 120
    }
  ],
  "generated_at": "2024-01-01T10:00:00Z"
}
```

### Processing Steps

```
1. Webhook Receives
   ↓
2. Immediate Response
   {success: true, schedule_id: "abc-123"}
   ↓
3. Process Data
   - Format dates
   - Calculate totals
   - Create summaries
   ↓
4. Split Path:
   ├─ Send team summary to Slack
   └─ Process each employee
      ↓
5. For Each Employee:
   ├─ Send personalized email
   ├─ Log to Google Sheets
   ├─ Create calendar event
   └─ Check if hours > threshold
      └─ Send alert if needed
```

## Notification Examples

### Slack Summary Message

```
📅 New Schedule Generated

• Schedule ID: `abc-123`
• Period: January 1, 2024 to January 15, 2024
• Total Employees: 2

Employees:
• John Doe (Developer) - 160h/month
• Jane Smith (Designer) - 120h/month
```

### Individual Email

```
Subject: Your Work Schedule - January 1, 2024

Hi John Doe,

Your schedule is ready!

📋 Details:
• Role: Developer
• Description: Full-stack developer
• Monthly Hours: 160
• Period: January 1, 2024 to January 15, 2024

Best regards,
HR Team
```

### Google Sheets Entry

| Schedule ID | Period Start | Period End | Employee Name | Email | Role | Monthly Hours | Generated At |
|-------------|--------------|------------|---------------|-------|------|---------------|--------------|
| abc-123 | 2024-01-01 | 2024-01-15 | John Doe | john@example.com | Developer | 160 | 2024-01-01 10:00 |
| abc-123 | 2024-01-01 | 2024-01-15 | Jane Smith | jane@example.com | Designer | 120 | 2024-01-01 10:00 |

### Calendar Event

```
Title: Schedule: John Doe
Date: Jan 1, 2024 - Jan 15, 2024
Description:
  Role: Developer
  Description: Full-stack developer
  Monthly Hours: 160
Attendees: john@example.com
```

## Integration Patterns

### Pattern 1: Notification Only

```
RestySched → n8n Webhook → Slack/Teams/Email
```

**Use:** Simple notifications
**Pros:** Fast, easy setup
**Cons:** No data logging

### Pattern 2: Notification + Logging

```
RestySched → n8n Webhook → Notifications
                         → Database/Sheets
```

**Use:** Track history
**Pros:** Audit trail, reporting
**Cons:** More setup

### Pattern 3: Full Automation

```
RestySched → n8n Webhook → Notifications
                         → Logging
                         → Calendar
                         → Approval flow
```

**Use:** Complete workflow
**Pros:** Fully automated
**Cons:** Complex, more maintenance

### Pattern 4: Conditional Routing

```
RestySched → n8n → Check conditions
                   ├─ Managers → Slack
                   ├─ Full-time → Email + Calendar
                   └─ Part-time → Email only
```

**Use:** Different handling per role
**Pros:** Customized per employee
**Cons:** Complex logic

## Error Handling Flow

```
┌──────────────┐
│  Any Node    │
│  Fails       │
└──────┬───────┘
       │
       │ On Error
       │
       ▼
┌──────────────┐
│ Error        │
│ Trigger      │
└──────┬───────┘
       │
       ├─────────────────────┐
       │                     │
       ▼                     ▼
┌──────────────┐    ┌───────────────┐
│ Log Error to │    │ Send Alert to │
│ Database     │    │ Slack/Email   │
└──────────────┘    └───────────────┘
```

## Approval Workflow (Advanced)

```
┌──────────────┐
│  RestySched  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  n8n Webhook │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│ Send Approval    │
│ Request to Slack │
│ [Approve/Reject] │
└──────┬───────────┘
       │
       │ Wait for response
       │
       ▼
┌──────────────┐
│ Slack Button │
│ Clicked      │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ IF Approved? │
└──┬────────┬──┘
   │        │
   │ Yes    │ No
   │        │
   ▼        ▼
┌──────┐ ┌──────┐
│ Send │ │ Notify│
│ Emails│ │ Reject│
└──────┘ └──────┘
```

## Timeline Example

### Automated Biweekly Schedule

```
Week 1, Monday 9:00 AM
│
├─ RestySched Scheduler runs
│  └─ Generate schedule
│     └─ Send to n8n webhook
│
├─ n8n receives (9:00:05 AM)
│  └─ Respond to RestySched
│  └─ Process data
│
├─ Send Slack summary (9:00:06 AM)
│  └─ Team sees notification
│
├─ Split employees (9:00:07 AM)
│  └─ Process 50 employees
│
├─ Send emails (9:00:08 - 9:00:30 AM)
│  └─ Rate limited: 2/second
│  └─ All employees notified
│
├─ Log to Sheets (9:00:08 AM)
│  └─ 50 rows added
│
├─ Create calendar events (9:00:08 - 9:00:45 AM)
│  └─ Rate limited by Google
│  └─ All events created
│
└─ Check high hours (9:00:10 AM)
   └─ 3 employees > 120h
   └─ Alert sent to #alerts
```

## Monitoring Dashboard

```
┌─────────────────────────────────────────┐
│        n8n Execution History            │
├─────────────────────────────────────────┤
│                                         │
│  Today's Schedules: 1                   │
│  Emails Sent: 50                        │
│  Success Rate: 100%                     │
│  Avg Duration: 42s                      │
│                                         │
│  Recent Executions:                     │
│  ✅ 09:00 - Schedule abc-123 (42s)      │
│  ✅ 08:45 - Test execution (2s)         │
│  ❌ 08:30 - Failed - SMTP error         │
│                                         │
│  Active Workflows: 2                    │
│  • RestySched - Full Automation         │
│  • RestySched - Simple Starter          │
│                                         │
└─────────────────────────────────────────┘
```

## Scaling Considerations

### Small Team (1-20 employees)

```
Simple workflow:
Webhook → Slack → Email

Execution time: < 5 seconds
```

### Medium Team (20-100 employees)

```
Standard workflow:
Webhook → Slack → Split → Email → Sheets

Execution time: 30-60 seconds
Consider: Rate limiting, batching
```

### Large Team (100+ employees)

```
Optimized workflow:
Webhook → Slack
        → Queue to database
        → Background processor
        → Batch emails (10 at a time)
        → Rate limit: 2/second

Execution time: 5-10 minutes
Consider: Dedicated email service, queue system
```

## Best Practices Diagram

```
✅ DO                           ❌ DON'T
┌──────────────────┐           ┌──────────────────┐
│ Respond to       │           │ Block webhook    │
│ webhook first    │           │ waiting for      │
│                  │           │ all processing   │
└──────────────────┘           └──────────────────┘

┌──────────────────┐           ┌──────────────────┐
│ Use error        │           │ Ignore errors    │
│ handling         │           │                  │
└──────────────────┘           └──────────────────┘

┌──────────────────┐           ┌──────────────────┐
│ Log executions   │           │ No audit trail   │
│ to database      │           │                  │
└──────────────────┘           └──────────────────┘

┌──────────────────┐           ┌──────────────────┐
│ Rate limit API   │           │ Spam APIs with   │
│ calls            │           │ requests         │
└──────────────────┘           └──────────────────┘

┌──────────────────┐           ┌──────────────────┐
│ Test with small  │           │ Deploy to full   │
│ dataset first    │           │ team untested    │
└──────────────────┘           └──────────────────┘
```

## Quick Reference

### Common Expressions

```javascript
// Format date
{{ new Date($json.period_start).toLocaleDateString() }}

// Count employees
{{ $json.employees.length }}

// Map employee names
{{ $json.employees.map(e => e.name).join(', ') }}

// Access previous node data
{{ $('Webhook Trigger').item.json.schedule_id }}

// Conditional text
{{ $json.monthly_hours > 120 ? 'High' : 'Normal' }}
```

### Common Filters

```javascript
// Only full-time
{{ $json.employees.filter(e => e.monthly_hours >= 160) }}

// Only managers
{{ $json.employees.filter(e => e.role.includes('Manager')) }}

// Sort by name
{{ $json.employees.sort((a, b) => a.name.localeCompare(b.name)) }}
```

---

For more details, see [README.md](README.md) in this directory.
