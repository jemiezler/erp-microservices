# API Gateway Architecture

## Overview

All services in this microservices architecture **must pass through the API Gateway** before being accessed by clients (frontend, mobile apps, external integrations).

```
┌─────────────────┐
│   Frontend      │
│ (localhost:3000)│
└────────┬────────┘
         │
         │ HTTP Request
         │ (through gateway)
         ↓
┌─────────────────────────────────────┐
│   API Gateway (localhost:8080)      │
│  - CORS Headers                     │
│  - Request ID Tracking              │
│  - Request Logging                  │
│  - Error Handling                   │
│  - Routing                          │
└────────┬────────────────────────────┘
         │
    ┌────┴────┬──────────┐
    ↓         ↓          ↓
┌────────┐ ┌────────┐ ┌────────┐
│   HR   │ │Finance │ │  Auth  │
│Service │ │Service │ │Service │
│:8081   │ │:8082   │ │:8083   │
└────────┘ └────────┘ └────────┘
    ↓         ↓          ↓
    └────────→ Database (localhost:5432)
```

## Gateway Features

### 1. **Request Routing**
- All requests start at `/api/v1/*`
- Routes are dynamically configured via environment variables
- Service URLs are internal (localhost) - not exposed

### 2. **Request/Response Handling**
- Request ID injection (`X-Request-ID`)
- CORS middleware for frontend origins
- Structured logging for all requests
- Error handling with meaningful responses

### 3. **Security**
- CORS configuration restricts origins
- Request validation
- Rate limiting ready (can be added)
- Services are internal-only

### 4. **Monitoring**
- Health check endpoint: `GET /health`
- Info endpoint: `GET /info` (lists all services)
- Request logging with timestamps
- Error tracking with request IDs

## Service Routes

### HR Service
```
GET    /api/v1/hr/employees
POST   /api/v1/hr/employees
GET    /api/v1/hr/employees/:id
PATCH  /api/v1/hr/employees/:id
DELETE /api/v1/hr/employees/:id
```

### Finance Service (Ready for expansion)
```
/api/v1/finance/*
```

### Auth Service (Ready for expansion)
```
/api/v1/auth/*
```

## Running the Gateway

### 1. **Start the Gateway**

```bash
cd apis/api-gateway
go run main.go
```

Output:
```
Starting API-GATEWAY on port 8080
HR Service: http://localhost:8081
Finance Service: http://localhost:8082
```

### 2. **Start Backend Services**

Each service should be running on its configured internal port:

```bash
# Terminal 1: HR Service
cd apis/services/hr-service
go run cmd/api/main.go

# Terminal 2: Finance Service
cd apis/services/finance-service
go run cmd/api/main.go
```

### 3. **Test Gateway**

```bash
# Health check
curl http://localhost:8080/health

# Gateway info
curl http://localhost:8080/info

# Access HR Service through gateway
curl -H "Authorization: Bearer token" \
     -H "X-Tenant-ID: default-tenant" \
     http://localhost:8080/api/v1/hr/employees
```

## Environment Configuration

### .env File

Create `.env` in the `apis/` directory:

```env
GATEWAY_PORT=8080
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

HR_SERVICE_URL=http://localhost:8081
FINANCE_SERVICE_URL=http://localhost:8082
AUTH_SERVICE_URL=http://localhost:8083

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=erp_db
```

### Frontend Configuration

Update `frontend/apps/hr-mfe/.env.local`:

```env
# Point to gateway, not direct service
NEXT_PUBLIC_HR_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_TENANT_ID=default-tenant
```

## Request Flow

### 1. Frontend Makes Request
```
GET http://localhost:8080/api/v1/hr/employees
Headers:
  - Authorization: Bearer <token>
  - X-Tenant-ID: default-tenant
  - Content-Type: application/json
```

### 2. Gateway Processes Request
```
1. Generates X-Request-ID (e.g., "abc123")
2. Applies CORS headers
3. Logs request with ID
4. Routes to HR Service: http://localhost:8081/api/v1/hr/employees
```

### 3. HR Service Responds
```
Response from http://localhost:8081
```

### 4. Gateway Returns to Frontend
```
Response with:
  - Original CORS headers
  - X-Request-ID header
  - Response body
```

## Debugging

### Check Gateway Health
```bash
curl http://localhost:8080/health
# Response: {"status":"healthy","timestamp":1234567890}
```

### View Service Info
```bash
curl http://localhost:8080/info
# Response shows all configured services
```

### Monitor Logs
Gateway logs show request flow:
```
[request-id] HR Service: GET /api/v1/hr/employees
[request-id] Response Status: 200
```

### Service Not Responding?
If you see 502 Bad Gateway:
1. Check that the service is running on the configured port
2. Verify the service URL in gateway configuration
3. Check firewall settings for localhost access

## Adding New Services

### 1. Start New Service on Internal Port
```go
// your-service/cmd/api/main.go
app.Listen("127.0.0.1:8084")  // Internal only
```

### 2. Add Gateway Route

Update `apis/api-gateway/main.go`:

```go
// Add to ServiceConfig struct
type ServiceConfig struct {
    HR      string
    Finance string
    Auth    string
    NewService string  // Add this
}

// Add to getServiceConfig
NewService: getEnv("NEW_SERVICE_URL", "http://localhost:8084"),

// Add route in main()
app.All("/api/v1/new-service/*", func(c fiber.Ctx) error {
    path := c.Path()
    target := config.NewService + path
    log.Printf("[%s] New Service: %s %s", c.Get("X-Request-ID"), c.Method(), path)
    
    if err := proxy.Do(c, target); err != nil {
        return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
            "error": "Failed to reach New Service",
            "request_id": c.Get("X-Request-ID"),
        })
    }
    return nil
})
```

### 3. Update Environment
```env
NEW_SERVICE_URL=http://localhost:8084
```

### 4. Update Frontend URL (if needed)
```bash
NEXT_PUBLIC_NEW_SERVICE_URL=http://localhost:8080/api/v1/new-service
```

## Security Considerations

### Production Deployment

1. **Internal Network Only**
   - Services communicate on private network
   - Only gateway is public-facing

2. **Authentication**
   - Add JWT verification in gateway middleware
   - Validate token before routing

3. **HTTPS**
   - Use HTTPS for all external communication
   - Internal services can use HTTP

4. **Rate Limiting**n   - Add rate limiting middleware to gateway
   - Protect backend services from overload

5. **Service Discovery**
   - Consider Consul, Kubernetes, or similar for dynamic routing

### Example JWT Middleware (Future)

```go
app.Use(func(c fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Missing Authorization header",
        })
    }
    
    // Validate JWT token here
    // token := strings.TrimPrefix(authHeader, "Bearer ")
    
    return c.Next()
})
```

## Monitoring & Metrics

### Metrics to Track
- Request count per service
- Response times per service
- Error rates
- Gateway uptime

### Log Aggregation
All requests are logged with request ID for easy tracing:
```
[abc123] GET /api/v1/hr/employees - 200 - 45ms
[abc124] POST /api/v1/finance/invoices - 201 - 123ms
[abc125] GET /api/v1/auth/verify - 401 - 5ms
```

## Troubleshooting

### Issue: 502 Bad Gateway

**Cause:** Backend service not responding

**Solution:**
```bash
# Check if service is running
curl http://localhost:8081/health

# Check logs
# Look for connection refused errors
```

### Issue: CORS Errors in Frontend

**Cause:** CORS_ORIGINS not configured properly

**Solution:**
Update `.env`:
```env
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
```

### Issue: Request Timeout

**Cause:** Backend service slow or hanging

**Solution:**
1. Check service logs
2. Increase timeout (if needed)
3. Check database connection

## Next Steps

1. ✅ Gateway routing for HR Service - DONE
2. ⏳ Add JWT middleware for authentication
3. ⏳ Add rate limiting
4. ⏳ Add response caching
5. ⏳ Add service health checks
6. ⏳ Add metrics/monitoring
7. ⏳ Deploy to production infrastructure

---

**Architecture Pattern:** API Gateway + Microservices
**Status:** Ready for Development & Testing
