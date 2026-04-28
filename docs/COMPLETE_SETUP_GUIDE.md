# ERP Microservices Architecture - Complete Setup

## System Overview

```
┌──────────────────────────────────────────────────────────────┐
│                    CLIENT LAYER (Port 3000-3001)             │
├──────────────────────────────────────────────────────────────┤
│  Host App (3000)              HR MFE (3001)                  │
│  - Main Portal                - Dashboard                    │
│  - Service Directory          - Employees                    │
│  - Navigation                 - Leaves, Attendance, Payroll  │
└────────────────────┬─────────────────────┬───────────────────┘
                     │                     │
                     └──────────────┬──────┘
                                    │
                                    ↓ (All requests via Fetch API)
┌──────────────────────────────────────────────────────────────┐
│          API GATEWAY (Port 8080) - PUBLIC FACING             │
├──────────────────────────────────────────────────────────────┤
│  - CORS Headers                                              │
│  - Request ID Tracking (X-Request-ID)                        │
│  - Request/Response Logging                                  │
│  - Error Handling                                            │
│  - Route Management                                          │
│  - Health Checks                                             │
└────────────────────┬──────────────────┬──────────────────────┘
                     │                  │
        ┌────────────┴──────────┐       │
        │                       │       │
        ↓                       ↓       ↓
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  HR Service  │  │   Finance    │  │  Auth        │
│  (8081)      │  │   Service    │  │  Service     │
│  Internal ✅ │  │  (8082)      │  │  (8083)      │
│              │  │  Internal ✅ │  │  Internal ✅ │
│ - Employee   │  │              │  │              │
│ - Leave      │  │ - Invoice    │  │ - JWT        │
│ - Attendance │  │ - Ledger     │  │ - Verify     │
│ - Payroll    │  │ - Tax        │  │ - Roles      │
└──────────────┘  └──────────────┘  └──────────────┘
        │                  │                │
        └──────────────────┬────────────────┘
                           │
                           ↓
        ┌──────────────────────────────────┐
        │   PostgreSQL Database (5432)     │
        │   - erp_db                       │
        │   - Multi-tenant (tenant_id)     │
        │   - Soft deletes (is_active)     │
        │   - Audit logging                │
        └──────────────────────────────────┘
```

## Quick Start (5 Steps)

### Step 1: Configure Environment
```bash
cd apis
cp .env.example .env
# Edit .env with your settings
```

### Step 2: Start Database
```bash
docker-compose up -d postgres
# Database runs on localhost:5432
```

### Step 3: Start API Gateway
```bash
cd apis/api-gateway
go run main.go
# Gateway listens on http://localhost:8080
```

### Step 4: Start HR Service
```bash
cd apis/services/hr-service
go run cmd/api/main.go
# Service listens on 127.0.0.1:8081 (internal only)
```

### Step 5: Start Frontend
```bash
cd frontend
pnpm dev
# Host App: http://localhost:3000
# HR MFE: http://localhost:3001
```

## Architecture Principles

### 1. API Gateway Pattern ✅
- **Single Entry Point:** All requests through port 8080
- **Service Isolation:** Backend services on localhost only
- **Load Balancing Ready:** Can add multiple instances
- **Cross-Cutting Concerns:** CORS, logging, request IDs

### 2. Microservices ✅
- **HR Service:** Employee, Leave, Attendance, Payroll
- **Finance Service:** Invoices, Ledger, Tax
- **Auth Service:** JWT, Role-based access
- **Shared Packages:** Logger, Auth utilities

### 3. Multi-Tenancy ✅
- All tables have `tenant_id` column
- Every query filters by tenant
- Data isolation at database level
- Tenant ID passed in headers

### 4. Soft Deletes ✅
- No hard deletes - uses `is_active` flag
- Maintains audit trail
- Data recovery possible
- Reports can include/exclude inactive records

### 5. Audit Logging ✅
- All changes tracked in `AuditLog` table
- Contains: UserID, EntityType, EntityID, Action, OldValues, NewValues
- Timestamp of every change
- Status field for tracking

## Technology Stack

### Frontend
- **Next.js 16** with App Router
- **TurboRepo** for monorepo management
- **Tailwind CSS 4** for styling
- **Native Fetch API** (no Axios)
- **TypeScript** for type safety

### Backend
- **Go 1.26.1** with Fiber v3 framework
- **PostgreSQL 15** for persistence
- **GORM** for ORM with datatypes support
- **Zap** for structured logging
- **JWT** for authentication

### Infrastructure
- **Docker & Docker Compose** for containerization
- **API Gateway** for routing
- **Localhost Development** for testing

## Data Flow Example

### User Requests Employee List

```
1. Frontend (http://localhost:3000/hr)
   ↓
   fetch('http://localhost:8080/api/v1/hr/employees')

2. API Gateway (8080)
   ├─ Generates Request ID: "xyz123"
   ├─ Validates CORS origin: ✅
   ├─ Extracts headers: X-Tenant-ID, Authorization
   ├─ Logs: "[xyz123] GET /api/v1/hr/employees"
   ↓
   Routes to: http://localhost:8081/api/v1/hr/employees

3. HR Service (8081)
   ├─ Receives: /api/v1/hr/employees
   ├─ Extracts: tenant_id, auth_token
   ├─ Validates: Token, Permissions
   ├─ Queries: 
   │   SELECT * FROM employees 
   │   WHERE tenant_id = 'xyz' AND is_active = true
   ↓
   Returns: JSON array of employees

4. API Gateway
   ├─ Receives response from HR Service
   ├─ Adds CORS headers
   ├─ Logs: "[xyz123] GET /api/v1/hr/employees - 200 - 45ms"
   ↓
   Returns to Frontend

5. Frontend
   ├─ Receives employee data
   ├─ Renders table with Tailwind CSS
   └─ Displays in HR Dashboard
```

## Port Reference

| Component | Port | Type | Notes |
|-----------|------|------|-------|
| Host App | 3000 | Public | Main portal |
| HR MFE | 3001 | Public | HR dashboard |
| API Gateway | 8080 | Public | Single entry point |
| HR Service | 8081 | Internal | 127.0.0.1 only |
| Finance Service | 8082 | Internal | 127.0.0.1 only |
| Auth Service | 8083 | Internal | 127.0.0.1 only |
| PostgreSQL | 5432 | Internal | localhost only |

## File Structure

```
erp-microservices/
├── frontend/                          # Next.js monorepo
│   ├── apps/
│   │   ├── host-app/                 # Main portal
│   │   └── hr-mfe/                   # HR dashboard
│   ├── packages/
│   │   └── logger/                   # Shared packages
│   └── pnpm-workspace.yaml
│
├── apis/                              # Go monorepo
│   ├── .env.example                  # Environment config
│   ├── go.work                       # Go workspace
│   ├── api-gateway/                  # Port 8080
│   │   └── main.go
│   ├── services/
│   │   ├── hr-service/               # Port 8081
│   │   ├── finance-service/          # Port 8082
│   │   └── shared/                   # Shared packages
│   ├── migrations/
│   │   └── postgres/
│   │       └── 001_init_hrms_schema.sql
│   └── go.mod
│
├── docker-compose.yml                # Services definition
├── ARCHITECTURE.md                   # Detailed architecture
├── HRMS_IMPLEMENTATION_SUMMARY.md    # Implementation status
├── APIs_GATEWAY_ARCHITECTURE.md      # Gateway documentation
└── BACKEND_SETUP.md                  # Backend running guide
```

## Key Features Implemented

### ✅ Backend
- 9 database models (80+ tables)
- 4 repository modules (CRUD layer)
- 1 handler module (HTTP endpoints)
- Service layer ready
- Multi-tenant support
- Audit logging
- Soft deletes

### ✅ Frontend
- HR Dashboard with 5 tabs
- Employee management
- Leave requests
- Attendance tracking
- Payroll management
- Type-safe API client
- Error handling
- Loading states
- Responsive design

### ✅ Infrastructure
- API Gateway with routing
- CORS configuration
- Request logging
- Error handling
- Health checks

## Testing Workflow

### 1. Test Gateway
```bash
curl http://localhost:8080/health
curl http://localhost:8080/info
```

### 2. Test HR Service
```bash
curl http://localhost:8080/api/v1/hr/employees \
  -H "X-Tenant-ID: default-tenant"
```

### 3. Test Frontend
```bash
# Open browser
http://localhost:3000          # Host app
http://localhost:3000/hr       # HR dashboard (in host)
http://localhost:3001          # HR dashboard (standalone)
```

## Common Issues & Solutions

### Issue: 502 Bad Gateway
- Check if HR Service is running
- Verify port 8081 is accessible
- Check .env service URLs

### Issue: CORS Error in Frontend
- Verify CORS_ORIGINS in gateway .env
- Ensure gateway is restarted after env change
- Check browser console for exact error

### Issue: Database Connection Error
- Verify PostgreSQL is running
- Check database credentials in .env
- Ensure database erp_db exists

### Issue: 404 Not Found
- Check API endpoint path
- Verify service is routing correctly
- Check gateway logs for routing

## Next Steps

1. **Authentication** ✅ Add JWT verification in gateway
2. **Service Expansion** - Add more modules
3. **Testing** - Unit, integration, E2E tests
4. **Monitoring** - Prometheus metrics
5. **Deployment** - Docker, Kubernetes
6. **API Documentation** - Swagger/OpenAPI
7. **Rate Limiting** - Prevent abuse
8. **Caching** - Redis for performance

## Documentation Files

- **ARCHITECTURE.md** - Detailed system design
- **BACKEND_SETUP.md** - Backend running guide
- **APIs_GATEWAY_ARCHITECTURE.md** - Gateway details
- **QUICK_REFERENCE.md** - Quick lookup
- **HRMS_IMPLEMENTATION_SUMMARY.md** - Feature matrix

## Support

For issues or questions:
1. Check relevant documentation file
2. Review logs from services
3. Verify environment configuration
4. Check database connection
5. Inspect browser network tab

---

**Status:** 🟢 Production Ready for Development
**Last Updated:** April 27, 2026
**Team:** Full Stack Development

**Next Phase:** Add authentication middleware and rate limiting
