# Backend Services Setup & Running Guide

## Architecture Overview

```
Frontend (port 3000-3001)
         ↓
API Gateway (port 8080) - Public facing
         ↓
├─ HR Service (port 8081) - Internal only
├─ Finance Service (port 8082) - Internal only
└─ Auth Service (port 8083) - Internal only
         ↓
PostgreSQL Database (port 5432)
```

**All frontend requests go through the API Gateway at port 8080**

## Prerequisites

1. **Go 1.26.1** installed
2. **PostgreSQL 15** running
3. **.env file** configured in `apis/` directory

## Quick Start

### 1. Configure Environment

```bash
cd apis
cp .env.example .env
```

Edit `.env`:
```env
GATEWAY_PORT=8080
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

HR_SERVICE_URL=http://localhost:8081
FINANCE_SERVICE_URL=http://localhost:8082

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=erp_db
DB_SSL_MODE=disable
```

### 2. Start Database

```bash
# Using Docker Compose
cd ../
docker-compose up -d
```

### 3. Start All Services

**Terminal 1: API Gateway**
```bash
cd apis/api-gateway
go run main.go
# Output: Starting API-GATEWAY on port 8080
```

**Terminal 2: HR Service**
```bash
cd apis/services/hr-service
go run cmd/api/main.go
# Output: HR-SERVICE listening on 127.0.0.1:8081
```

**Terminal 3: Frontend**
```bash
cd frontend
pnpm dev
# Outputs running on http://localhost:3000 and http://localhost:3001
```

## Service Ports

| Service | Port | Access | Purpose |
|---------|------|--------|---------|
| API Gateway | 8080 | ✅ Public | Frontend → Gateway |
| HR Service | 8081 | 🔒 Internal | Gateway only |
| Finance Service | 8082 | 🔒 Internal | Gateway only |
| Auth Service | 8083 | 🔒 Internal | Gateway only |
| PostgreSQL | 5432 | 🔒 Internal | All services |
| Host App | 3000 | ✅ Public | Main portal |
| HR MFE | 3001 | ✅ Public | HR dashboard |

**🔒 Internal** = Localhost only (not exposed)
**✅ Public** = Accessible from browser

## Testing the Setup

### 1. Check Gateway Health

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": 1234567890
}
```

### 2. Check Gateway Info

```bash
curl http://localhost:8080/info
```

Response:
```bash
{
  "service": "API-GATEWAY",
  "version": "1.0.0",
  "services": {
    "hr": "http://localhost:8081",
    "finance": "http://localhost:8082",
    "auth": "http://localhost:8083"
  }
}
```

### 3. Test HR Service Through Gateway

```bash
curl -H "X-Tenant-ID: default-tenant" \
     http://localhost:8080/api/v1/hr/employees
```

## How Frontend Calls Backend

### Old Way (DEPRECATED) ❌
```
Frontend → HR Service (8081)
```

### New Way (CORRECT) ✅
```
Frontend → API Gateway (8080) → HR Service (8081)
```

### Frontend Configuration

File: `frontend/apps/hr-mfe/.env.local`
```env
# Point to gateway, not direct service
NEXT_PUBLIC_HR_API_URL=http://localhost:8080/api/v1
```

## Troubleshooting

### Problem: "Connection refused" error

**Check if service is running:**
```bash
# Check if gateway is running
curl http://localhost:8080/health

# Check if hr-service is running
curl -v http://localhost:8081/health

# Check database connection
psql -h localhost -U postgres -d erp_db
```

### Problem: HR app shows "Failed to load dashboard"

**Check the flow:**
1. ✅ Gateway running? → `curl http://localhost:8080/health`
2. ✅ HR Service running? → `curl http://localhost:8081/health`
3. ✅ Frontend .env.local points to gateway? → Check `NEXT_PUBLIC_HR_API_URL`

### Problem: CORS errors in browser

**Frontend error:** "Access to fetch blocked by CORS policy"

**Solution:** 
Check gateway `.env`:
```env
# Must include frontend origins
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
```

### Problem: Database connection error

**Check database is running:**
```bash
# Windows
netstat -ano | findstr :5432

# macOS/Linux
lsof -i :5432
```

**Start database:**
```bash
docker-compose up -d
```

## Development Workflow

### 1. Update HR Service Code

```bash
# Edit file
vim apis/services/hr-service/cmd/api/main.go

# Restart service (in Terminal 2)
# Ctrl+C to stop
go run cmd/api/main.go
```

### 2. Hot Reload (Optional)

Use `air` for auto-reload:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### 3. Database Changes

```bash
# Database migrations are in:
migrations/postgres/001_init_hrms_schema.sql

# After changes:
# 1. Update file
# 2. Run migrations: psql -U postgres -d erp_db -f migrations/postgres/001_init_hrms_schema.sql
# 3. Restart services
```

## Directory Structure

```
apis/
├── .env.example              # Environment template
├── go.work                   # Workspace config
├── api-gateway/              # Main gateway (port 8080)
│   └── main.go
├── services/
│   ├── hr-service/           # HR service (port 8081, internal)
│   │   └── cmd/api/main.go
│   ├── finance-service/      # Finance service (port 8082, internal)
│   └── shared/               # Shared packages
│       ├── logger/
│       └── auth/
└── migrations/
    └── postgres/
        └── 001_init_hrms_schema.sql
```

## API Endpoints (Through Gateway)

### Employee Management
```
GET    /api/v1/hr/employees                  # List employees
POST   /api/v1/hr/employees                  # Create employee
GET    /api/v1/hr/employees/:id              # Get employee
PATCH  /api/v1/hr/employees/:id              # Update employee
DELETE /api/v1/hr/employees/:id              # Delete employee
```

All requests must include:
```
X-Tenant-ID: default-tenant
Authorization: Bearer <token> (if required)
```

## Common Commands

```bash
# Start all backend services
cd apis/api-gateway && go run main.go &
cd apis/services/hr-service && go run cmd/api/main.go &

# Check service status
curl http://localhost:8080/info

# View gateway logs
tail -f <gateway-output-file>

# Test an endpoint
curl http://localhost:8080/api/v1/hr/employees

# Stop all services
pkill -f "go run"
```

## Next Steps

1. ✅ Gateway setup complete
2. ⏳ Add more services (Finance, Auth, etc.)
3. ⏳ Implement authentication middleware
4. ⏳ Add database migrations
5. ⏳ Deploy to production

## Performance Tips

1. **Database Pooling** - Keep persistent connections
2. **Response Caching** - Cache frequently accessed data
3. **Rate Limiting** - Add limits to prevent abuse
4. **Monitoring** - Track response times and errors

---

**Status:** Ready for Development
**Last Updated:** April 27, 2026
