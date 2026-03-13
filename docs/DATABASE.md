# Database Management Guide

Lumine provides comprehensive database management with support for multiple database engines and admin panels.

## Supported Databases

### Relational Databases
- **MySQL 8.0** - Most popular open-source database
- **PostgreSQL 16** - Advanced open-source database
- **MariaDB 11.2** - MySQL fork with additional features

### NoSQL Databases
- **MongoDB 7.0** - Document-oriented database
- **Redis 7.2** - In-memory data structure store
- **Elasticsearch 8.11** - Search and analytics engine

## Quick Setup

### Start All Databases

```bash
make db-setup
```

This will start:
- MySQL on port 3306
- PostgreSQL on port 5432
- MariaDB on port 3307
- MongoDB on port 27017
- Redis on port 6379
- Elasticsearch on port 9200

Plus admin panels:
- phpMyAdmin on port 8080
- Adminer on port 8081
- Mongo Express on port 8082
- Redis Commander on port 8083
- pgAdmin on port 8084

### Stop Databases

```bash
make db-stop
```

### Restart Databases

```bash
make db-restart
```

### View Logs

```bash
make db-logs
```

## Admin Panels

### phpMyAdmin (MySQL/MariaDB)
- **URL**: http://localhost:8080
- **Features**:
  - Visual database management
  - SQL query editor
  - Import/Export data
  - User management
  - Table designer

### Adminer (Universal)
- **URL**: http://localhost:8081
- **Features**:
  - Supports MySQL, PostgreSQL, SQLite, MS SQL, Oracle
  - Lightweight (single PHP file)
  - Clean interface
  - Export in multiple formats

### Mongo Express (MongoDB)
- **URL**: http://localhost:8082
- **Login**: admin / admin
- **Features**:
  - Database/collection management
  - Document CRUD operations
  - JSON editor
  - Import/Export

### Redis Commander (Redis)
- **URL**: http://localhost:8083
- **Features**:
  - Key management
  - Value inspection
  - TTL management
  - CLI interface

### pgAdmin (PostgreSQL)
- **URL**: http://localhost:8084
- **Login**: admin@lumine.local / admin
- **Features**:
  - Full PostgreSQL management
  - Query tool
  - Schema designer
  - Backup/Restore

## Connection Details

### MySQL

```bash
# CLI
mysql -h localhost -P 3306 -u root -proot

# Connection String
mysql://root:root@localhost:3306/lumine

# DSN (Go)
root:root@tcp(localhost:3306)/lumine?charset=utf8mb4&parseTime=True
```

**Default Credentials:**
- Host: localhost
- Port: 3306
- User: root
- Password: root
- Database: lumine

### PostgreSQL

```bash
# CLI
psql -h localhost -p 5432 -U postgres -d lumine

# Connection String
postgresql://postgres:postgres@localhost:5432/lumine

# DSN
host=localhost port=5432 user=postgres password=postgres dbname=lumine sslmode=disable
```

**Default Credentials:**
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: lumine

### MariaDB

```bash
# CLI
mysql -h localhost -P 3307 -u root -proot

# Connection String
mysql://root:root@localhost:3307/lumine
```

**Default Credentials:**
- Host: localhost
- Port: 3307
- User: root
- Password: root
- Database: lumine

### MongoDB

```bash
# CLI
mongosh mongodb://root:root@localhost:27017/lumine

# Connection String
mongodb://root:root@localhost:27017/lumine
```

**Default Credentials:**
- Host: localhost
- Port: 27017
- User: root
- Password: root
- Database: lumine

### Redis

```bash
# CLI
redis-cli -h localhost -p 6379

# Connection String
redis://localhost:6379
```

**Default Credentials:**
- Host: localhost
- Port: 6379
- No password (default)

### Elasticsearch

```bash
# CLI
curl http://localhost:9200

# Connection String
http://localhost:9200
```

**Default Credentials:**
- Host: localhost
- Port: 9200
- No authentication (development mode)

## Database Initialization

### MySQL/MariaDB

Custom initialization scripts in `scripts/mysql-init/`:

```sql
-- 01-init.sql
CREATE DATABASE IF NOT EXISTS myapp;
CREATE USER IF NOT EXISTS 'dev'@'%' IDENTIFIED BY 'dev';
GRANT ALL PRIVILEGES ON myapp.* TO 'dev'@'%';
FLUSH PRIVILEGES;
```

### PostgreSQL

Custom initialization scripts in `scripts/postgres-init/`:

```sql
-- 01-init.sql
CREATE DATABASE myapp;
CREATE USER dev WITH PASSWORD 'dev';
GRANT ALL PRIVILEGES ON DATABASE myapp TO dev;
```

### MongoDB

Custom initialization scripts in `scripts/mongo-init/`:

```javascript
// 01-init.js
db = db.getSiblingDB('myapp');
db.createCollection('users');
db.users.createIndex({ "email": 1 }, { unique: true });
```

## Using with Projects

### Laravel (MySQL)

```env
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=lumine
DB_USERNAME=root
DB_PASSWORD=root
```

### Django (PostgreSQL)

```python
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': 'lumine',
        'USER': 'postgres',
        'PASSWORD': 'postgres',
        'HOST': 'localhost',
        'PORT': '5432',
    }
}
```

### Next.js (MongoDB)

```env
MONGODB_URI=mongodb://root:root@localhost:27017/lumine
```

### Express (Redis)

```javascript
const redis = require('redis');
const client = redis.createClient({
  url: 'redis://localhost:6379'
});
```

## Backup & Restore

### MySQL Backup

```bash
# Backup
docker exec lumine-mysql mysqldump -u root -proot lumine > backup.sql

# Restore
docker exec -i lumine-mysql mysql -u root -proot lumine < backup.sql
```

### PostgreSQL Backup

```bash
# Backup
docker exec lumine-postgres pg_dump -U postgres lumine > backup.sql

# Restore
docker exec -i lumine-postgres psql -U postgres lumine < backup.sql
```

### MongoDB Backup

```bash
# Backup
docker exec lumine-mongodb mongodump --uri="mongodb://root:root@localhost:27017/lumine" --out=/backup

# Restore
docker exec lumine-mongodb mongorestore --uri="mongodb://root:root@localhost:27017/lumine" /backup/lumine
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :3306

# Kill the process
kill -9 <PID>

# Or change port in docker-compose.db.yml
```

### Container Won't Start

```bash
# Check logs
docker logs lumine-mysql

# Remove and recreate
docker rm -f lumine-mysql
make db-setup
```

### Data Persistence

All database data is stored in Docker volumes:
- `lumine_mysql_data`
- `lumine_postgres_data`
- `lumine_mariadb_data`
- `lumine_mongodb_data`
- `lumine_redis_data`
- `lumine_elasticsearch_data`

To completely remove data:
```bash
make db-clean
```

## Performance Tuning

### MySQL

Edit `docker-compose.db.yml`:
```yaml
mysql:
  command: --max_connections=200 --innodb_buffer_pool_size=1G
```

### PostgreSQL

```yaml
postgres:
  command: -c shared_buffers=256MB -c max_connections=200
```

### Redis

```yaml
redis:
  command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
```

## Security Notes

⚠️ **Development Only**: The default credentials are for development only. Never use these in production!

For production:
1. Change all default passwords
2. Enable SSL/TLS
3. Configure firewalls
4. Use environment variables
5. Enable authentication
6. Regular backups
7. Monitor access logs


## 🗑️ Cleanup & Removal

### Safe Cleanup (Containers Only)

```bash
# Stop all containers
make containers-stop

# Remove stopped containers only
make containers-clean

# Remove all containers (keeps data)
make containers-remove
```

### Complete Cleanup (Containers + Data)

```bash
# Remove containers and volumes
make containers-prune

# Reset databases (remove and recreate)
make db-reset
```

### Nuclear Option

```bash
# Remove EVERYTHING (requires typing 'DESTROY')
make clean-everything
```

This will remove:
- All Lumine containers
- All Lumine volumes (database data)
- Lumine network
- Build artifacts
- Docker cache

### Selective Removal

#### Remove Specific Container

```bash
# Using Docker directly
docker rm -f lumine-mysql

# Or stop first
docker stop lumine-mysql
docker rm lumine-mysql
```

#### Remove Specific Volume

```bash
# List volumes
make volumes-list

# Remove specific volume
docker volume rm lumine_mysql_data
```

#### Remove Specific Database

```bash
# Stop and remove MySQL only
docker compose -f docker-compose.db.yml stop mysql
docker compose -f docker-compose.db.yml rm -f mysql

# Remove MySQL volume
docker volume rm lumine_mysql_data
```

### Cleanup from TUI

Inside Lumine TUI:

1. Navigate to service/database
2. Press `delete` or `backspace`
3. Choose cleanup option:
   - **Remove Container** - Remove container only
   - **Remove with Volume** - Remove container + data
   - **Remove All Containers** - Remove all Lumine containers
   - **Nuclear Cleanup** - Remove everything

4. Confirm action

### Cleanup Workflow

#### Before Removing

```bash
# 1. Check what's running
make containers-list
make volumes-list

# 2. Backup important data
docker exec lumine-mysql mysqldump -u root -proot lumine > backup.sql

# 3. Stop containers
make containers-stop
```

#### After Removing

```bash
# 1. Verify removal
docker ps -a | grep lumine
docker volume ls | grep lumine

# 2. Recreate if needed
make db-setup

# 3. Restore data
docker exec -i lumine-mysql mysql -u root -proot lumine < backup.sql
```

### Troubleshooting Cleanup

#### "Volume is in use"

```bash
# Stop all containers first
make containers-stop

# Then remove volumes
make volumes-remove
```

#### "Network has active endpoints"

```bash
# Remove all containers first
make containers-remove

# Then remove network
make network-remove
```

#### "Permission denied"

```bash
# Use sudo for Docker commands
sudo make containers-prune

# Or add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### Cleanup Checklist

Before major cleanup:

- [ ] Backup important databases
- [ ] Export configurations
- [ ] Note custom environment variables
- [ ] Save project files
- [ ] Document custom setups

After cleanup:

- [ ] Verify all containers removed
- [ ] Check volumes removed
- [ ] Confirm network removed
- [ ] Test fresh installation
- [ ] Restore backups if needed

### Automated Cleanup Scripts

Create a cleanup script:

```bash
#!/bin/bash
# cleanup-lumine.sh

echo "Stopping containers..."
make containers-stop

echo "Backing up databases..."
docker exec lumine-mysql mysqldump -u root -proot --all-databases > backup-$(date +%Y%m%d).sql

echo "Removing containers..."
make containers-remove

echo "Removing volumes..."
make volumes-remove

echo "Cleanup complete!"
```

Make it executable:
```bash
chmod +x cleanup-lumine.sh
./cleanup-lumine.sh
```

### Recovery After Cleanup

If you accidentally removed everything:

```bash
# 1. Reinstall Lumine
make install

# 2. Setup databases
make db-setup

# 3. Restore from backup
docker exec -i lumine-mysql mysql -u root -proot < backup.sql

# 4. Verify
make containers-list
```
