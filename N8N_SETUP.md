# n8n Workflow Setup Guide

This guide will help you set up an n8n workflow to receive schedule data from RestySched.

## Quick Start

1. Create a new workflow in n8n
2. Add a Webhook trigger node
3. Configure the webhook
4. Add your processing nodes
5. Copy the webhook URL to your `.env` file

## Detailed Setup

### Step 1: Create Webhook Trigger

1. In n8n, create a new workflow
2. Add a **Webhook** node as the trigger
3. Configure the webhook:
   - **HTTP Method**: POST
   - **Path**: `/schedule-webhook` (or your preferred path)
   - **Authentication**: None (or configure as needed)
   - **Response Mode**: "Respond Immediately"

### Step 2: Webhook URL

After saving, n8n will provide a webhook URL like:
```
https://your-n8n-instance.com/webhook/abc123def456
```

Copy this URL to your `.env` file:
```env
N8N_WEBHOOK_URL=https://your-n8n-instance.com/webhook/abc123def456
```

### Step 3: Understanding the Payload

RestySched sends the following JSON structure:

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
    },
    {
      "id": "employee-uuid-2",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "role": "Designer",
      "role_description": "UI/UX designer creating user interfaces",
      "monthly_hours": 120
    }
  ],
  "generated_at": "2024-01-01T00:00:00Z"
}
```

### Step 4: Processing the Data

Here are some example workflow nodes you might add:

#### Example 1: Send Email Notifications

Add a **Split In Batches** node to process employees:
```
Webhook → Split In Batches → Email (Gmail/Outlook)
```

Configure Email node:
- **To**: `{{ $json.email }}`
- **Subject**: `Schedule for {{ $json.period_start }} - {{ $json.period_end }}`
- **Message**:
```
Hi {{ $json.name }},

Your schedule for the period {{ $json.period_start }} to {{ $json.period_end }}:

Role: {{ $json.role }}
Description: {{ $json.role_description }}
Monthly Hours: {{ $json.monthly_hours }}

Best regards,
Schedule System
```

#### Example 2: Store in Google Sheets

Add a **Google Sheets** node:
```
Webhook → Split In Batches → Google Sheets
```

Configure Google Sheets:
- **Operation**: Append Row
- **Sheet**: Your schedule sheet
- **Columns**: Map each field to columns

#### Example 3: Create Calendar Events

Add a **Google Calendar** node:
```
Webhook → Split In Batches → Google Calendar
```

#### Example 4: Send to Slack

Add a **Slack** node:
```
Webhook → Slack
```

Message template:
```
🗓️ New Schedule Generated!

Period: {{ $json.period_start }} - {{ $json.period_end }}
Employees: {{ $json.employees.length }}

Employees:
{{ $json.employees.map(e => `• ${e.name} (${e.role}) - ${e.monthly_hours}h/month`).join('\n') }}
```

#### Example 5: Complex Workflow

```
                    ┌──────────────┐
                    │   Webhook    │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │ Set Variables│
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──────┐ ┌──▼─────┐ ┌───▼────┐
       │Split Batches│ │ Slack  │ │ Log DB │
       └──────┬──────┘ └────────┘ └────────┘
              │
       ┌──────▼──────┐
       │ Send Emails │
       └─────────────┘
```

### Step 5: Testing the Integration

#### Manual Test from n8n

1. In n8n, click "Execute Workflow"
2. Select "Listen for test webhook"
3. From RestySched, generate a schedule and click "Send to n8n"
4. Check n8n for the received data

#### Test Payload

You can also manually test with this curl command:

```bash
curl -X POST https://your-n8n-instance.com/webhook/abc123def456 \
  -H "Content-Type: application/json" \
  -d '{
    "schedule_id": "test-123",
    "period_start": "2024-01-01T00:00:00Z",
    "period_end": "2024-01-15T00:00:00Z",
    "employees": [
      {
        "id": "emp-1",
        "name": "Test Employee",
        "email": "test@example.com",
        "role": "Developer",
        "role_description": "Test description",
        "monthly_hours": 160
      }
    ],
    "generated_at": "2024-01-01T00:00:00Z"
  }'
```

## Common Use Cases

### Use Case 1: Weekly Schedule Emails

**Workflow**:
1. Webhook receives schedule
2. Store in database for historical tracking
3. For each employee:
   - Generate personalized schedule PDF
   - Send email with PDF attachment

**n8n Nodes**:
- Webhook
- MySQL/PostgreSQL
- Split In Batches
- Function (generate PDF)
- Gmail/SendGrid

### Use Case 2: Multi-Channel Notifications

**Workflow**:
1. Webhook receives schedule
2. Send summary to Slack
3. For each employee:
   - Send email
   - Send SMS (if urgent)
   - Create calendar event

**n8n Nodes**:
- Webhook
- Slack
- Split In Batches
- Gmail
- Twilio
- Google Calendar

### Use Case 3: Approval Workflow

**Workflow**:
1. Webhook receives schedule
2. Send to manager for approval via Slack
3. Wait for approval
4. If approved:
   - Send to all employees
   - Create calendar events
5. If rejected:
   - Notify scheduler
   - Log rejection reason

**n8n Nodes**:
- Webhook
- Slack (with buttons)
- Wait
- IF condition
- Split In Batches
- Gmail
- Google Calendar

## Accessing Data in n8n

Use these expressions to access data in your nodes:

```javascript
// Schedule ID
{{ $json.schedule_id }}

// Period dates
{{ $json.period_start }}
{{ $json.period_end }}

// Number of employees
{{ $json.employees.length }}

// Current employee (in Split In Batches)
{{ $json.name }}
{{ $json.email }}
{{ $json.role }}
{{ $json.role_description }}
{{ $json.monthly_hours }}

// All employee names
{{ $json.employees.map(e => e.name).join(', ') }}

// Filter employees by criteria
{{ $json.employees.filter(e => e.monthly_hours > 100) }}
```

## Error Handling

Add error handling to your workflow:

1. **Webhook Errors**: Add an "Error Trigger" node
2. **Retry Logic**: Configure retry settings on critical nodes
3. **Notifications**: Send alerts on failures via Slack/Email

Example error workflow:
```
Error Trigger → Format Error Message → Slack/Email
```

## Security Best Practices

1. **Use Authentication**: Enable basic auth or header auth on webhook
2. **Validate Payload**: Add a Function node to validate the payload structure
3. **Rate Limiting**: Configure rate limits in n8n
4. **HTTPS Only**: Always use HTTPS webhooks
5. **Secret Tokens**: Add a secret token header for verification

Example validation function:
```javascript
// Validate required fields
const required = ['schedule_id', 'period_start', 'period_end', 'employees'];
const missing = required.filter(field => !items[0].json[field]);

if (missing.length > 0) {
  throw new Error(`Missing required fields: ${missing.join(', ')}`);
}

// Validate employees array
if (!Array.isArray(items[0].json.employees) || items[0].json.employees.length === 0) {
  throw new Error('Employees must be a non-empty array');
}

return items;
```

## Monitoring and Logging

1. **Execution History**: Review in n8n's execution log
2. **Custom Logging**: Add HTTP Request nodes to log to external services
3. **Metrics**: Track schedule counts, processing times, etc.

## Troubleshooting

### Webhook Not Receiving Data

1. Check webhook URL in `.env` file
2. Verify n8n workflow is activated
3. Check n8n logs for errors
4. Test with curl command

### Payload Issues

1. Check JSON format in n8n webhook test
2. Verify all required fields are present
3. Check data types match expectations

### Processing Errors

1. Review execution log in n8n
2. Check individual node configurations
3. Add logging nodes to debug

## Additional Resources

- [n8n Documentation](https://docs.n8n.io/)
- [n8n Webhook Node](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.webhook/)
- [n8n Community](https://community.n8n.io/)

## Support

If you encounter issues:
1. Check n8n execution logs
2. Verify webhook URL and configuration
3. Test with the provided curl command
4. Review RestySched logs for errors
