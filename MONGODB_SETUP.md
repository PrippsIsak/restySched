# MongoDB Setup Guide

This guide will help you set up MongoDB for RestySched.

## Quick Start with Docker (Recommended)

The easiest way to run MongoDB is using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Start MongoDB on port 27017
- Create a database named `restysched`
- Persist data in a Docker volume

### Stop MongoDB

```bash
docker-compose down
```

### Stop and Remove Data

```bash
docker-compose down -v
```

## Manual MongoDB Installation

### Windows

1. **Download MongoDB:**
   - Visit [MongoDB Download Center](https://www.mongodb.com/try/download/community)
   - Download MongoDB Community Server for Windows
   - Run the installer

2. **Start MongoDB:**
   ```powershell
   # Start as a service
   net start MongoDB

   # Or run manually
   "C:\Program Files\MongoDB\Server\7.0\bin\mongod.exe" --dbpath="C:\data\db"
   ```

3. **Verify Installation:**
   ```powershell
   "C:\Program Files\MongoDB\Server\7.0\bin\mongosh.exe"
   ```

### macOS

1. **Install with Homebrew:**
   ```bash
   brew tap mongodb/brew
   brew install mongodb-community
   ```

2. **Start MongoDB:**
   ```bash
   brew services start mongodb-community
   ```

3. **Verify Installation:**
   ```bash
   mongosh
   ```

### Linux (Ubuntu/Debian)

1. **Import MongoDB GPG Key:**
   ```bash
   curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | \
      sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor
   ```

2. **Add MongoDB Repository:**
   ```bash
   echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | \
      sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
   ```

3. **Install MongoDB:**
   ```bash
   sudo apt-get update
   sudo apt-get install -y mongodb-org
   ```

4. **Start MongoDB:**
   ```bash
   sudo systemctl start mongod
   sudo systemctl enable mongod
   ```

5. **Verify Installation:**
   ```bash
   mongosh
   ```

## MongoDB Atlas (Cloud)

For production or remote access, consider MongoDB Atlas:

1. **Create Account:**
   - Visit [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
   - Sign up for free

2. **Create Cluster:**
   - Click "Build a Cluster"
   - Choose Free Tier (M0)
   - Select your region
   - Click "Create Cluster"

3. **Configure Access:**
   - Add IP Address: Click "Network Access" → "Add IP Address"
   - Create Database User: Click "Database Access" → "Add New Database User"

4. **Get Connection String:**
   - Click "Connect" on your cluster
   - Choose "Connect your application"
   - Copy the connection string

5. **Update `.env`:**
   ```env
   MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
   MONGO_DATABASE=restysched
   ```

## Configuration

Update your `.env` file with MongoDB settings:

```env
# Local MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=restysched

# Or MongoDB Atlas
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DATABASE=restysched
```

### Connection String Formats

**Local MongoDB:**
```
mongodb://localhost:27017
```

**MongoDB with Authentication:**
```
mongodb://username:password@localhost:27017
```

**MongoDB Replica Set:**
```
mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myReplicaSet
```

**MongoDB Atlas:**
```
mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
```

## Verifying Your Setup

### 1. Check MongoDB is Running

```bash
# Using mongosh
mongosh

# Should see:
# > Connected to MongoDB
```

### 2. Test RestySched Connection

Run the application:

```bash
go run cmd/server/main.go
```

Look for this log message:
```
Connected to MongoDB database: restysched
```

### 3. View Database in MongoDB Compass

1. Download [MongoDB Compass](https://www.mongodb.com/products/compass)
2. Connect to `mongodb://localhost:27017`
3. You should see the `restysched` database after creating employees

## Collections and Indexes

RestySched automatically creates these collections and indexes:

### employees Collection

**Indexes:**
- `email` (unique) - Ensures no duplicate emails
- `active` - Optimizes queries for active employees

### schedules Collection

**Indexes:**
- `period_start`, `period_end` (compound) - Optimizes period queries
- `status` - Optimizes status filtering

## Common Operations

### View All Employees

```javascript
use restysched
db.employees.find().pretty()
```

### View All Schedules

```javascript
db.schedules.find().pretty()
```

### Count Active Employees

```javascript
db.employees.countDocuments({ active: true })
```

### Find Employee by Email

```javascript
db.employees.findOne({ email: "john@example.com" })
```

### View Indexes

```javascript
db.employees.getIndexes()
db.schedules.getIndexes()
```

### Drop Database (Careful!)

```javascript
use restysched
db.dropDatabase()
```

## Backup and Restore

### Backup Database

```bash
mongodump --uri="mongodb://localhost:27017" --db=restysched --out=./backup
```

### Restore Database

```bash
mongorestore --uri="mongodb://localhost:27017" --db=restysched ./backup/restysched
```

### Export Collection to JSON

```bash
mongoexport --uri="mongodb://localhost:27017" --db=restysched --collection=employees --out=employees.json
```

### Import Collection from JSON

```bash
mongoimport --uri="mongodb://localhost:27017" --db=restysched --collection=employees --file=employees.json
```

## Performance Tips

### 1. Use Indexes

RestySched creates indexes automatically, but if you add custom queries, add appropriate indexes.

### 2. Connection Pooling

The MongoDB driver automatically manages connection pooling. Default settings work well for most use cases.

### 3. Monitor Performance

Use MongoDB Atlas monitoring or install Mongo Express:

```bash
docker run -d -p 8081:8081 \
  -e ME_CONFIG_MONGODB_URL=mongodb://mongodb:27017 \
  --link mongodb:mongo \
  mongo-express
```

Access at: http://localhost:8081

## Troubleshooting

### Connection Refused

**Problem:** `connection refused`

**Solution:**
1. Verify MongoDB is running: `mongosh`
2. Check port 27017 is available: `netstat -an | grep 27017`
3. Review MongoDB logs

### Authentication Failed

**Problem:** `authentication failed`

**Solution:**
1. Verify username/password in connection string
2. Check user has appropriate permissions
3. Use admin database for authentication

### Database Not Created

**Problem:** Database doesn't appear

**Solution:**
- MongoDB creates databases on first write operation
- Add an employee through the UI to trigger database creation

### Slow Queries

**Problem:** Queries are slow

**Solution:**
1. Check indexes exist: `db.employees.getIndexes()`
2. Use query profiler: `db.setProfilingLevel(2)`
3. View slow queries: `db.system.profile.find().pretty()`

### Docker Volume Issues

**Problem:** Data not persisting

**Solution:**
```bash
# Remove and recreate volume
docker-compose down -v
docker-compose up -d
```

## Security Best Practices

### 1. Enable Authentication

```bash
mongod --auth --bind_ip localhost
```

### 2. Create Admin User

```javascript
use admin
db.createUser({
  user: "admin",
  pwd: "securepassword",
  roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
})
```

### 3. Create Application User

```javascript
use restysched
db.createUser({
  user: "restysched",
  pwd: "apppassword",
  roles: [ { role: "readWrite", db: "restysched" } ]
})
```

### 4. Update Connection String

```env
MONGO_URI=mongodb://restysched:apppassword@localhost:27017/restysched
```

### 5. Network Security

- Bind to localhost only for local development
- Use firewalls to restrict access
- Enable SSL/TLS for production

## Additional Resources

- [MongoDB Documentation](https://docs.mongodb.com/)
- [MongoDB University (Free Courses)](https://university.mongodb.com/)
- [MongoDB Compass](https://www.mongodb.com/products/compass)
- [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
- [Go MongoDB Driver](https://www.mongodb.com/docs/drivers/go/current/)

## Support

For MongoDB-specific issues:
1. Check [MongoDB Community Forums](https://www.mongodb.com/community/forums/)
2. Review [Go Driver Documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
3. Open an issue in the RestySched repository for integration questions
