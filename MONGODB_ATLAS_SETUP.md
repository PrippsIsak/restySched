# MongoDB Atlas Setup Guide

Complete guide for using RestySched with MongoDB Atlas (cloud-hosted MongoDB).

## Why MongoDB Atlas?

- ✅ **Free tier available** - Perfect for getting started
- ✅ **Automatic backups** - Built-in data protection
- ✅ **Global deployment** - Low latency worldwide
- ✅ **Auto-scaling** - Grows with your needs
- ✅ **Monitoring** - Built-in performance tracking
- ✅ **High availability** - 99.995% uptime SLA

## Quick Start

### Step 1: Create MongoDB Atlas Account

1. Go to [MongoDB Atlas](https://www.mongodb.com/cloud/atlas/register)
2. Sign up with email or Google/GitHub
3. Verify your email address

### Step 2: Create a Free Cluster

1. Click **"Build a Database"**
2. Choose **"M0 FREE"** tier
   - 512 MB storage
   - Shared RAM
   - Perfect for development and small teams

3. **Select Cloud Provider & Region:**
   - Provider: AWS, Google Cloud, or Azure
   - Region: Choose closest to your location
   - Example: `AWS - us-east-1 (N. Virginia)`

4. **Cluster Name:**
   - Default: `Cluster0`
   - Or custom: `restysched-prod`

5. Click **"Create Cluster"**
   - Wait 3-5 minutes for deployment

### Step 3: Configure Database Access

1. **Create Database User:**
   - Go to **"Database Access"** (left sidebar)
   - Click **"Add New Database User"**
   - Choose **"Password"** authentication
   - Username: `restysched`
   - Password: Click **"Autogenerate Secure Password"** (copy it!)
   - Or set custom password
   - Database User Privileges: **"Read and write to any database"**
   - Click **"Add User"**

2. **Save Your Credentials:**
   ```
   Username: restysched
   Password: [your-generated-password]
   ```
   ⚠️ **IMPORTANT:** Save these credentials securely!

### Step 4: Configure Network Access

1. **Add IP Address:**
   - Go to **"Network Access"** (left sidebar)
   - Click **"Add IP Address"**

2. **Options:**

   **Option A: Development (Allow from Anywhere)**
   - Click **"Allow Access from Anywhere"**
   - IP: `0.0.0.0/0`
   - ⚠️ Not recommended for production
   - ✅ Good for testing

   **Option B: Production (Specific IP)**
   - Click **"Add Current IP Address"**
   - Or manually add your server's IP
   - More secure

   **Option C: Multiple IPs**
   - Add each developer's IP
   - Add your server's IP
   - Add CI/CD IP if applicable

3. Click **"Confirm"**

### Step 5: Get Connection String

1. Go to **"Database"** (left sidebar)
2. Click **"Connect"** on your cluster
3. Choose **"Connect your application"**
4. **Driver:** Go / Version: 1.7 or later
5. **Copy the connection string:**

   ```
   mongodb+srv://restysched:<password>@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```

6. **Replace `<password>`** with your actual password:

   ```
   mongodb+srv://restysched:YourActualPassword@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```

### Step 6: Configure RestySched

1. **Update `.env` file:**

   ```env
   # MongoDB Atlas Configuration
   MONGO_URI=mongodb+srv://restysched:YourActualPassword@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
   MONGO_DATABASE=restysched

   # Other settings
   SERVER_PORT=8080
   N8N_WEBHOOK_URL=
   ENABLE_SCHEDULER=true
   ```

2. **Important Notes:**
   - No spaces in the connection string
   - Password must be URL-encoded if it contains special characters
   - Database name is set separately in `MONGO_DATABASE`

### Step 7: Test Connection

1. **Start RestySched:**
   ```bash
   go run cmd/server/main.go
   ```

2. **Look for success message:**
   ```
   Connected to MongoDB database: restysched
   Server starting on port 8080
   ```

3. **If connection fails**, check:
   - Password is correct
   - IP address is whitelisted
   - Connection string format is correct

### Step 8: Verify in Atlas

1. Go to **"Database"** → **"Browse Collections"**
2. After adding first employee:
   - Database: `restysched`
   - Collections: `employees`, `schedules`
3. Click collection to view data

## Connection String Formats

### Standard Format

```
mongodb+srv://username:password@cluster.xxxxx.mongodb.net/database?options
```

### With Specific Database

```
mongodb+srv://restysched:password@cluster0.xxxxx.mongodb.net/restysched?retryWrites=true&w=majority
```

### With Additional Options

```
mongodb+srv://restysched:password@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority&maxPoolSize=50&connectTimeoutMS=10000
```

### URL Encoding Passwords

If password contains special characters:

| Character | Encoded |
|-----------|---------|
| @ | %40 |
| : | %3A |
| / | %2F |
| ? | %3F |
| # | %23 |
| [ | %5B |
| ] | %5D |

Example:
- Password: `P@ss:word/123`
- Encoded: `P%40ss%3Aword%2F123`

## Configuration Options

### Basic Configuration

```env
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/
MONGO_DATABASE=restysched
```

### With Connection Options

```env
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/?retryWrites=true&w=majority&maxPoolSize=20&serverSelectionTimeoutMS=5000
```

### Common Options

| Option | Description | Default | Recommended |
|--------|-------------|---------|-------------|
| `retryWrites` | Retry failed writes | true | true |
| `w` | Write concern | majority | majority |
| `maxPoolSize` | Max connections | 100 | 20-50 |
| `minPoolSize` | Min connections | 0 | 5 |
| `serverSelectionTimeoutMS` | Timeout for server selection | 30000 | 10000 |
| `connectTimeoutMS` | Connection timeout | 10000 | 10000 |

## Security Best Practices

### 1. Strong Passwords

```bash
# Generate secure password
openssl rand -base64 32
```

Use this for your database user password.

### 2. IP Whitelisting

**Development:**
```
Allow from anywhere: 0.0.0.0/0
```

**Production:**
```
Specific IPs only:
- 203.0.113.0 (production server)
- 198.51.100.0 (backup server)
```

### 3. Database User Privileges

**Least Privilege Principle:**
- App user: Read/Write to `restysched` database only
- Admin user: Separate account for maintenance

**Create Limited User in Atlas:**
1. Database Access → Add New User
2. Custom Role → Select specific database
3. Privileges: `readWrite` on `restysched` only

### 4. Environment Variables

**Never commit `.env` to Git:**

```bash
# .gitignore
.env
.env.local
.env.*.local
```

**Use secrets management in production:**
- AWS Secrets Manager
- HashiCorp Vault
- Kubernetes Secrets

### 5. TLS/SSL

MongoDB Atlas enforces TLS by default.

Verify with:
```
mongodb+srv://... (uses TLS automatically)
```

## Monitoring & Performance

### Atlas Built-in Monitoring

1. **Real-time Performance:**
   - Go to **"Metrics"** tab
   - View operations/second
   - Monitor connections
   - Check latency

2. **Query Performance:**
   - **"Performance Advisor"** tab
   - Suggests indexes
   - Identifies slow queries

3. **Alerts:**
   - **"Alerts"** tab → **"Configure Alert"**
   - Alert on high connections
   - Alert on low disk space
   - Email/SMS notifications

### Viewing Data

1. **Atlas UI:**
   - **"Browse Collections"**
   - View/edit documents
   - Run queries

2. **MongoDB Compass:**
   - Download [MongoDB Compass](https://www.mongodb.com/products/compass)
   - Connect with same connection string
   - Visual query builder
   - Schema analysis

3. **Command Line:**
   ```bash
   # Install mongosh
   brew install mongosh  # macOS
   choco install mongosh # Windows

   # Connect
   mongosh "mongodb+srv://cluster.mongodb.net/" --username restysched
   ```

## Scaling

### Free Tier (M0)
- Storage: 512 MB
- RAM: Shared
- Connections: 500
- **Good for:** Development, small teams (< 20 employees)

### Paid Tiers

**M10 ($0.08/hour = ~$57/month):**
- Storage: 10 GB
- RAM: 2 GB
- Connections: Unlimited
- **Good for:** Small production (20-100 employees)

**M20 ($0.20/hour = ~$145/month):**
- Storage: 20 GB
- RAM: 4 GB
- Auto-scaling available
- **Good for:** Growing teams (100-500 employees)

**M30+:** Enterprise scale

### When to Upgrade

Monitor these metrics in Atlas:
- ⚠️ **Storage > 80%** - Upgrade soon
- ⚠️ **Connections > 400** - Upgrade soon
- ⚠️ **CPU > 75%** - Performance degradation
- ⚠️ **Slow queries** - Add indexes or upgrade

## Backup & Restore

### Automatic Backups (Paid Tiers)

1. **Enable Backups:**
   - Cluster settings → Enable
   - Continuous backups
   - Point-in-time recovery

2. **Restore:**
   - Backups tab → Select snapshot
   - Restore to new cluster
   - Or download data

### Manual Backups (Free Tier)

```bash
# Export collections
mongoexport --uri="mongodb+srv://user:pass@cluster.mongodb.net/restysched" \
  --collection=employees \
  --out=employees.json

mongoexport --uri="mongodb+srv://user:pass@cluster.mongodb.net/restysched" \
  --collection=schedules \
  --out=schedules.json

# Import collections
mongoimport --uri="mongodb+srv://user:pass@cluster.mongodb.net/restysched" \
  --collection=employees \
  --file=employees.json

mongoimport --uri="mongodb+srv://user:pass@cluster.mongodb.net/restysched" \
  --collection=schedules \
  --file=schedules.json
```

## Multi-Region Deployment

### Global Clusters (M30+)

1. **Create Global Cluster:**
   - Database → Create → Global Cluster
   - Select multiple regions

2. **Configure Zones:**
   - US East: Primary
   - EU West: Read replica
   - Asia Pacific: Read replica

3. **Update Connection String:**
   - Automatically routes to nearest region
   - Low latency worldwide

## Troubleshooting

### Connection Timeout

**Error:**
```
Failed to initialize MongoDB: connection timeout
```

**Solutions:**
1. Check IP whitelist
2. Verify connection string
3. Check network/firewall
4. Test with mongosh

### Authentication Failed

**Error:**
```
authentication failed
```

**Solutions:**
1. Verify username/password
2. Check URL encoding
3. Ensure user has correct permissions
4. Re-create database user

### Connection String Issues

**Error:**
```
error parsing uri
```

**Solutions:**
1. Check for spaces in string
2. Ensure special characters are encoded
3. Verify format: `mongodb+srv://user:pass@host/`
4. No database name in URI (use MONGO_DATABASE env var)

### Slow Queries

**Solutions:**
1. Check indexes in Atlas Performance Advisor
2. Create suggested indexes
3. Upgrade cluster tier
4. Optimize queries

### Out of Connections

**Error:**
```
too many connections
```

**Solutions:**
1. Check for connection leaks
2. Reduce maxPoolSize
3. Upgrade to paid tier
4. Monitor active connections

## Migration from Local MongoDB

### Export from Local

```bash
mongodump --uri="mongodb://localhost:27017" --db=restysched --out=./backup
```

### Import to Atlas

```bash
mongorestore --uri="mongodb+srv://user:pass@cluster.mongodb.net" \
  --db=restysched \
  ./backup/restysched
```

### Update Configuration

```env
# Before (local)
MONGO_URI=mongodb://localhost:27017

# After (Atlas)
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/
```

Restart RestySched - done!

## Cost Optimization

### Free Tier Limits

- Storage: 512 MB
- Operations: No limit
- Data transfer: 10 GB/month

**Tips to stay in free tier:**
- Clean old schedules regularly
- Delete inactive employees
- Monitor storage usage

### Paid Tier Optimization

1. **Right-size cluster:**
   - Start with M10
   - Scale up if needed
   - Don't over-provision

2. **Use auto-scaling:**
   - Saves money during low usage
   - Scales up when needed

3. **Pause development clusters:**
   - Pause when not in use
   - Resume when needed

## Best Practices for Production

### 1. Separate Clusters

```
restysched-dev  (Free M0)
restysched-staging (M10)
restysched-prod (M20+)
```

### 2. Connection Pooling

```env
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/?maxPoolSize=20&minPoolSize=5
```

### 3. Monitoring

- Enable alerts for 80% storage
- Monitor slow queries
- Track connection count
- Set up status page

### 4. Indexes

Atlas auto-creates indexes, verify:
```javascript
db.employees.getIndexes()
db.schedules.getIndexes()
```

### 5. Regular Backups

- Enable continuous backups (paid tier)
- Or manual exports weekly (free tier)

## Support Resources

- **Atlas Documentation:** https://docs.atlas.mongodb.com/
- **Community Forums:** https://www.mongodb.com/community/forums/
- **Atlas Support:** Available in Atlas dashboard
- **RestySched Issues:** [GitHub Repository]

## Quick Reference

### Connection String Template

```env
MONGO_URI=mongodb+srv://USERNAME:PASSWORD@CLUSTER.mongodb.net/?retryWrites=true&w=majority
MONGO_DATABASE=restysched
```

### Common Commands

```bash
# Test connection
mongosh "mongodb+srv://cluster.mongodb.net/" -u restysched -p

# Export data
mongoexport --uri="CONNECTION_STRING" --collection=employees --out=backup.json

# Import data
mongoimport --uri="CONNECTION_STRING" --collection=employees --file=backup.json

# View indexes
mongosh> db.employees.getIndexes()
```

### Atlas Checklist

- [ ] Create cluster
- [ ] Create database user
- [ ] Whitelist IP addresses
- [ ] Copy connection string
- [ ] Update .env file
- [ ] Test connection
- [ ] Enable monitoring alerts
- [ ] Set up backups (if paid tier)

---

**You're now ready to use MongoDB Atlas with RestySched!** 🚀

For local MongoDB setup, see [MONGODB_SETUP.md](MONGODB_SETUP.md)
