# n8n Workflows for RestySched

This directory contains ready-to-use n8n workflows for automating schedule distribution and processing.

## Available Workflows

### 1. Simple Starter (`simple-starter.json`)
**Best for:** Getting started quickly

**Features:**
- ✅ Webhook trigger
- ✅ Slack notifications (optional)
- ✅ Email to each employee (optional)
- ✅ Minimal setup required

**Use case:** Testing the integration or simple notifications

### 2. Full Automation (`schedule-automation.json`)
**Best for:** Production use with multiple integrations

**Features:**
- ✅ Webhook trigger
- ✅ Slack team notifications
- ✅ Individual employee emails
- ✅ Google Sheets logging
- ✅ Google Calendar events
- ✅ High hours alerts
- ✅ Data processing and formatting

**Use case:** Complete automated workflow with logging and calendar

## Quick Start

### Step 1: Import Workflow

1. Open n8n (locally or cloud)
2. Click **"Add workflow"** → **"Import from File"**
3. Select `simple-starter.json` or `schedule-automation.json`
4. Click **"Import"**

### Step 2: Get Webhook URL

1. Click on the **"Webhook Trigger"** node
2. Click **"Test step"** or **"Listen for test event"**
3. Copy the **Production URL** (looks like: `https://your-n8n.com/webhook/restysched`)
4. Save this URL for later

### Step 3: Configure RestySched

Update your `.env` file:

```env
N8N_WEBHOOK_URL=https://your-n8n.com/webhook/restysched
```

Restart RestySched:
```bash
go run cmd/server/main.go
```

### Step 4: Test the Integration

1. In RestySched, go to http://localhost:8080/schedules
2. Click **"Generate Biweekly Schedule"**
3. Click **"Send to n8n"**
4. Check n8n - you should see the execution!

## Workflow Details

### Simple Starter Workflow

```
Webhook → Send Response → Slack Notification
                       → Split Employees → Send Emails
```

**Nodes:**
1. **Webhook Trigger**: Receives schedule data from RestySched
2. **Send Response**: Confirms receipt to RestySched
3. **Send to Slack**: Posts summary to Slack channel (disabled by default)
4. **Split Employees**: Separates each employee into individual items
5. **Send Email**: Sends personalized email to each employee (disabled by default)

**To enable Slack:**
1. Click the **"Send to Slack"** node
2. Add Slack credentials
3. Select your channel
4. Click the node → **Settings** → Uncheck **"Disabled"**

**To enable Email:**
1. Click the **"Send Email"** node
2. Add SMTP credentials
3. Update the from email address
4. Click the node → **Settings** → Uncheck **"Disabled"**

### Full Automation Workflow

```
                     → Slack Summary
                     → Split Employees → Email
Webhook → Response                    → Google Sheets
         → Process Data               → Calendar
                                      → High Hours Alert
```

**Additional Nodes:**
- **Process Schedule Data**: Formats dates and creates summaries
- **Create Summary**: Builds a formatted message
- **Save to Google Sheets**: Logs all schedules
- **Create Calendar Event**: Adds events for each employee
- **High Hours Check**: Alerts if hours > 120
- **Alert High Hours**: Sends warning to Slack

## Configuration Guides

### Slack Integration

1. **Create Slack App:**
   - Go to https://api.slack.com/apps
   - Click **"Create New App"** → **"From scratch"**
   - Name: "RestySched"
   - Select your workspace

2. **Add Permissions:**
   - Go to **"OAuth & Permissions"**
   - Add scopes:
     - `chat:write`
     - `chat:write.public`
   - Click **"Install to Workspace"**
   - Copy the **Bot User OAuth Token**

3. **Configure in n8n:**
   - In n8n, click the Slack node
   - Click **"Create New Credential"**
   - Paste your token
   - Save

4. **Select Channel:**
   - Click the Slack node
   - Choose your channel (e.g., `#schedules`)
   - Test the node

### Gmail Integration

1. **Enable Gmail API:**
   - Go to https://console.cloud.google.com
   - Create a new project or select existing
   - Enable **Gmail API**

2. **Create OAuth Credentials:**
   - Go to **"Credentials"** → **"Create Credentials"** → **"OAuth client ID"**
   - Application type: **Web application**
   - Add authorized redirect URI: `https://your-n8n.com/rest/oauth2-credential/callback`
   - Copy **Client ID** and **Client Secret**

3. **Configure in n8n:**
   - Click the Gmail node
   - Create new credential
   - Enter Client ID and Secret
   - Click **"Connect my account"**
   - Authorize access

### Google Sheets Integration

1. **Create Spreadsheet:**
   - Go to Google Sheets
   - Create new spreadsheet named **"RestySched Logs"**
   - Add headers: `Schedule ID`, `Period Start`, `Period End`, `Employee Name`, `Email`, `Role`, `Monthly Hours`, `Generated At`
   - Copy the spreadsheet ID from URL

2. **Enable Sheets API:**
   - Follow same steps as Gmail
   - Enable **Google Sheets API**

3. **Configure in n8n:**
   - Use same OAuth credentials as Gmail
   - Select your spreadsheet
   - Map columns to data

### Google Calendar Integration

1. **Enable Calendar API:**
   - Follow same steps as Gmail/Sheets
   - Enable **Google Calendar API**

2. **Configure in n8n:**
   - Use same OAuth credentials
   - Select calendar (usually "Primary")
   - Configure event details

### SMTP Email (Alternative to Gmail)

1. **Get SMTP Credentials:**
   - From your email provider
   - Or use services like SendGrid, Mailgun

2. **Common SMTP Settings:**
   ```
   Gmail:
   - Host: smtp.gmail.com
   - Port: 587
   - Secure: TLS

   Outlook:
   - Host: smtp.office365.com
   - Port: 587
   - Secure: STARTTLS

   SendGrid:
   - Host: smtp.sendgrid.net
   - Port: 587
   - User: apikey
   - Password: <your-api-key>
   ```

3. **Configure in n8n:**
   - Click Email node
   - Add SMTP credentials
   - Enter host, port, username, password
   - Test

## Customization

### Modify Slack Message

Edit the **"Send to Slack"** node:

```javascript
📅 **New Schedule Generated**

• Schedule ID: `{{ $json.schedule_id }}`
• Period: {{ new Date($json.period_start).toLocaleDateString() }} to {{ new Date($json.period_end).toLocaleDateString() }}
• Total Employees: {{ $json.employees.length }}

**Employees:**
{{ $json.employees.map(e => `• ${e.name} (${e.role}) - ${e.monthly_hours}h/month`).join('\n') }}
```

### Modify Email Template

Edit the **"Send Email"** node:

```
Subject: Your Work Schedule - {{ new Date($('Webhook Trigger').item.json.period_start).toLocaleDateString() }}

Body:
Hi {{ $json.name }},

Your schedule is ready!

📋 Details:
• Role: {{ $json.role }}
• Description: {{ $json.role_description }}
• Monthly Hours: {{ $json.monthly_hours }}
• Period: {{ ... }}

Best regards,
HR Team
```

### Add Teams Integration

1. **Add Microsoft Teams node**
2. **Create Teams webhook:**
   - In Teams channel: **"..."** → **"Connectors"** → **"Incoming Webhook"**
   - Name: "RestySched"
   - Copy webhook URL

3. **Configure node:**
   ```json
   {
     "text": "New schedule generated!",
     "sections": [{
       "activityTitle": "Schedule Details",
       "facts": [{
         "name": "Period",
         "value": "{{ $json.period_start }} to {{ $json.period_end }}"
       }, {
         "name": "Employees",
         "value": "{{ $json.employees.length }}"
       }]
     }]
   }
   ```

### Add Database Logging

1. **Add Postgres/MySQL node**
2. **Create table:**
   ```sql
   CREATE TABLE schedules (
     id SERIAL PRIMARY KEY,
     schedule_id VARCHAR(255),
     period_start TIMESTAMP,
     period_end TIMESTAMP,
     employee_name VARCHAR(255),
     employee_email VARCHAR(255),
     role VARCHAR(255),
     monthly_hours INT,
     created_at TIMESTAMP DEFAULT NOW()
   );
   ```

3. **Configure node:**
   - Operation: Insert
   - Map fields from schedule data

### Add Approval Step

1. **Add Slack node** with buttons:
   ```javascript
   {
     "text": "New schedule ready. Approve?",
     "attachments": [{
       "text": "Schedule for {{ $json.period_start }}",
       "callback_id": "schedule_approval",
       "actions": [{
         "name": "approve",
         "text": "Approve",
         "type": "button",
         "value": "approved"
       }, {
         "name": "reject",
         "text": "Reject",
         "type": "button",
         "value": "rejected"
       }]
     }]
   }
   ```

2. **Add Slack Webhook node** to receive button clicks
3. **Add IF node** to check approval
4. **Route to email nodes** only if approved

## Payload Reference

RestySched sends this JSON structure:

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
      "role_description": "Full-stack developer",
      "monthly_hours": 160,
      "active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "generated_at": "2024-01-01T00:00:00Z"
}
```

### Accessing Data in Nodes

```javascript
// Schedule information
$json.schedule_id
$json.period_start
$json.period_end
$json.generated_at

// Employee array
$json.employees.length
$json.employees[0].name

// Inside Split Employees node
$json.name
$json.email
$json.role
$json.monthly_hours

// Reference previous nodes
$('Webhook Trigger').item.json.schedule_id
$('Process Schedule Data').item.json.period_start_formatted
```

## Troubleshooting

### Webhook not receiving data

1. **Check webhook URL:**
   - Ensure URL in `.env` matches n8n webhook URL
   - Should be Production URL, not Test URL

2. **Activate workflow:**
   - Click **"Active"** toggle in top-right
   - Workflow must be active to receive webhooks

3. **Test manually:**
   ```bash
   curl -X POST https://your-n8n.com/webhook/restysched \
     -H "Content-Type: application/json" \
     -d '{"schedule_id":"test","period_start":"2024-01-01T00:00:00Z","period_end":"2024-01-15T00:00:00Z","employees":[],"generated_at":"2024-01-01T00:00:00Z"}'
   ```

### Slack messages not sending

1. **Check credentials:**
   - Verify token is correct
   - Check bot has permission to post

2. **Check channel:**
   - Verify channel name is correct
   - Channel should be public or bot should be invited

3. **Test manually:**
   - Click node → **"Test step"**

### Emails not sending

1. **Check SMTP settings:**
   - Verify host, port, username, password
   - Check TLS/SSL settings

2. **Check from address:**
   - Some providers require specific from addresses

3. **Check rate limits:**
   - Gmail: 500 emails/day
   - SendGrid: varies by plan

### Google Sheets errors

1. **Check permissions:**
   - Verify OAuth has Sheets access
   - Check spreadsheet is accessible

2. **Check spreadsheet ID:**
   - Copy from URL
   - Don't include `/edit#gid=0`

3. **Check column mapping:**
   - Verify column names match exactly

## Performance Tips

### 1. Batch Processing

Instead of sending individual emails, batch them:

```javascript
// In a Code node after splitting
const batch = [];
for (const item of $input.all()) {
  batch.push(item.json);
}

// Send one email with all employees
return [{json: {employees: batch}}];
```

### 2. Error Handling

Add error handling nodes:

1. Add **"Error Trigger"** node
2. Connect to Slack for error notifications
3. Log errors to database

### 3. Rate Limiting

For large employee lists:

1. Use **"Split In Batches"** with batch size: 10
2. Add **"Wait"** node between batches: 1 second
3. Prevents API rate limits

## Examples

### Example 1: Simple Slack Only

```
Webhook → Response → Slack
```

**Use case:** Just want notifications, no emails

### Example 2: Email Only

```
Webhook → Response → Split → Email
```

**Use case:** Send to employees, no Slack

### Example 3: Conditional Routing

```
Webhook → Split → IF (check role)
                    ├─ Managers → Slack
                    └─ Staff → Email
```

**Use case:** Different notifications per role

### Example 4: Multi-Channel

```
Webhook → Split → Email
                → Slack
                → Teams
                → Sheets
```

**Use case:** Notify everywhere

## Best Practices

1. **Always respond to webhook** - Use "Respond to Webhook" node
2. **Enable error notifications** - Add Error Trigger node
3. **Log executions** - Save to database or sheets
4. **Test with one employee first** - Verify before full rollout
5. **Set execution timeout** - Prevent long-running workflows
6. **Use environment variables** - For credentials and config
7. **Monitor execution history** - Check for failures
8. **Version your workflows** - Export and save changes

## Support

- **n8n Documentation**: https://docs.n8n.io
- **n8n Community**: https://community.n8n.io
- **RestySched Issues**: See main README.md

## Next Steps

1. Import a workflow
2. Configure credentials
3. Test with RestySched
4. Enable disabled nodes gradually
5. Customize for your needs
6. Add more integrations

Happy automating! 🚀
