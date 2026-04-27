​# HRMS Quick Reference Guide

## Project Structure Summary

### Created Models (12 HR Modules)

#### 1. **organization.go** - Organization & Structure
- Organization
- Department  
- JobTitle
- Location
- EmployeeHierarchy

#### 2. **employee.go** - Core HR Master Data
- Employee (100+ fields including personal, employment, compensation)
- Dependent
- EmergencyContact
- Qualification
- WorkExperience

#### 3. **recruitment.go** - ATS & Onboarding
- JobPosting
- Candidate
- Interview
- JobOffer
- OnboardingTask

#### 4. **attendance.go** - Time & Attendance
- Attendance (check-in/out with geo-location)
- TimeSheet & TimeSheetEntry
- Shift & ShiftAssignment
- BiometricData
- WorkFromHome

#### 5. **leave.go** - Leave Management
- LeaveType
- LeavePolicy
- LeaveAllocation
- Leave (with multi-level approval)
- LeaveApproval
- HolidayCalendar & Holiday
- LeaveEncashment

#### 6. **performance.go** - Performance Management
- PerformanceGoal & GoalProgress
- PerformanceRating
- CompetencyRating
- Competency
- EmployeeReview & ReviewQuestion
- TrainingRequest

#### 7. **compensation.go** - Payroll & Benefits
- SalaryStructure & SalaryComponent
- EmployeeSalary & SalaryComponentValue
- Payroll & PayrollLine
- BenefitPlan
- EmployeeBenefit & BenefitClaim

#### 8. **workflow.go** - Approval & RBAC
- ApprovalWorkflow & ApprovalLevel
- ApprovalRequest & Approval
- RBACRole & Permission
- AuditLog
- SystemSetting

#### 9. **lms.go** - Learning Management
- LMSCourse & CourseModule
- CourseEnrollment & LessonProgress
- Quiz & QuizQuestion
- QuizAttempt & QuizAnswer
- Certificate

### Created Repository Layer

#### Data Access Layer Pattern
```
BaseRepository (common CRUD)
    ├── EmployeeRepository (employee operations)
    ├── LeaveRepository (leave & holiday operations)
    ├── AttendanceRepository (attendance & WFH operations)
    └── PayrollRepository (payroll & benefits operations)
```

### Database Features

✅ **Multi-Tenancy**: All tables have tenant_id index
✅ **Audit Trail**: AuditLog table for tracking changes
✅ **Soft Deletes**: is_active flag for logical deletion
✅ **Relationships**: Proper foreign keys and eager loading
✅ **Indexing**: Strategic indexes for performance
✅ **JSON Support**: JSONB fields for flexible data

## Key API Patterns

### Resource URL Pattern
```
/api/v1/{module}/{resource}
/api/v1/employees
/api/v1/leaves
/api/v1/attendance
/api/v1/payroll
/api/v1/recruitment
/api/v1/performance
/api/v1/learning
```

### Standard Operations (CRUD)
```
GET    /{resource}           # List with filters
POST   /{resource}           # Create
GET    /{resource}/:id       # Get one
PATCH  /{resource}/:id       # Update
DELETE /{resource}/:id       # Soft delete
```

### Special Operations
```
POST   /{resource}/:id/approve        # Approve request
POST   /{resource}/:id/reject         # Reject request
GET    /{resource}/pending            # Get pending approvals
PATCH  /{resource}/:id/delegate       # Delegate approval
```

## Repository Method Naming Convention

- `GetAll{Entity}`: Fetch multiple records with filters
- `Get{Entity}ByID`: Fetch by primary key
- `Create{Entity}`: Insert new record
- `Update{Entity}`: Update existing record
- `Delete{Entity}`: Soft delete (set is_active=false)
- `Get{Entity}Stats`: Aggregate statistics
- `Get{Entity}History`: Get historical records

## Database Connection

### Environment Variables
```bash
DATABASE_URL=postgres://user:password@localhost:5432/hr_db
NATS_URL=nats://localhost:4222
REDIS_URL=redis://localhost:6379
LOG_LEVEL=info
```

### Connection Pool Settings (GORM)
```go
sqlDB, _ := DB.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Multi-Tenant Architecture

### Tenant Isolation Pattern
```go
// Always filter by tenant_id
query := db.Where("tenant_id = ?", userTenantID)

// Example: Get employee
db.Where("tenant_id = ? AND id = ?", tenantID, empID).First(&employee)
```

### Tenant ID in Context
```go
// Extract from JWT token
tenantID := c.Locals("tenant_id").(string)

// Or from request header
tenantID := c.Get("X-Tenant-ID")
```

## Event Publishing Pattern

### Event Types to Publish
```
employee.{created|updated|terminated}
leave.{created|approved|rejected|cancelled}
payroll.{generated|approved|posted|paid}
attendance.{checkin|checkout}
performance.review.{created|submitted|approved}
recruitment.offer.{created|accepted|rejected}
```

### Event Payload Structure
```json
{
  "event_type": "employee.created",
  "event_id": "uuid",
  "timestamp": "ISO8601",
  "tenant_id": "acme_corp",
  "user_id": "user_123",
  "data": { /* entity data */ }
}
```

## RBAC Permissions Matrix

### Employee Resource
- `employees.view` - View employee details
- `employees.create` - Create new employee
- `employees.edit` - Edit employee
- `employees.delete` - Deactivate employee
- `employees.export` - Export employee data

### Leave Resource
- `leaves.view` - View leaves
- `leaves.create` - Apply for leave
- `leaves.approve` - Approve leave
- `leaves.reject` - Reject leave
- `leaves.edit` - Edit leave settings

### Payroll Resource
- `payroll.view` - View payroll
- `payroll.generate` - Generate payroll
- `payroll.approve` - Approve payroll
- `payroll.post` - Post to finance
- `payroll.edit` - Edit payroll

## Common Response Codes

### Success (2xx)
- `200 OK` - Successful operation
- `201 Created` - Resource created
- `204 No Content` - No data returned

### Client Error (4xx)
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication failed
- `403 Forbidden` - No permission
- `404 Not Found` - Resource not found
- `409 Conflict` - Business rule violation

### Server Error (5xx)
- `500 Internal Server Error` - Unexpected error
- `503 Service Unavailable` - Service down

## Testing Checklist

### Unit Tests
- [ ] Model validation
- [ ] Repository CRUD operations
- [ ] Business logic (leave calculations, payroll computations)
- [ ] Permission checks

### Integration Tests
- [ ] API endpoints
- [ ] Multi-tenant isolation
- [ ] Approval workflows
- [ ] Event publishing

### E2E Tests
- [ ] Complete employee lifecycle
- [ ] Leave request to approval
- [ ] Payroll generation and posting
- [ ] User roles and permissions

## Performance Optimization Tips

1. **Eager Loading**: Preload relationships to avoid N+1 queries
   ```go
   db.Preload("Employee").Preload("Department").Find(&leaves)
   ```

2. **Pagination**: Always paginate large result sets
   ```go
   db.Limit(20).Offset((page-1)*20).Find(&employees)
   ```

3. **Indexes**: Create indexes on frequently filtered columns
   - tenant_id (all tables)
   - employee_id (most tables)
   - status columns
   - date ranges

4. **Caching**: Cache master data (leave types, departments, etc.)
   ```go
   cache.Set("leave_types:"+tenantID, leaveTypes, 24*time.Hour)
   ```

5. **Connection Pooling**: Configure appropriate pool size based on load

## Migration Strategy

### New Environment Setup
```bash
# Run migrations
./migrate -path ./migrations -database $DATABASE_URL up

# Seed data
./seed-data.sh
```

### Backward Compatibility
- Always create nullable columns for new fields
- Use migration scripts for schema changes
- Test migrations on staging first

## Monitoring & Debugging

### Key Logs
```log
[hr-service] Employee created: EMP-001 by user_123
[hr-service] Leave approved: Leave ID 42 by manager_5
[hr-service] Payroll posted: Payroll 2024-04 by admin
[hr-service] RBAC check failed: user_3 lacks employees.edit
```

### Health Check Endpoint
```
GET /health
Response: {"status": "healthy", "uptime": "1h23m"}
```

### Metrics to Track
- Average response time per endpoint
- Error rate by endpoint
- Database query latency
- Payroll processing duration
- Leave approval SLA compliance

## Common Issues & Solutions

### Issue: Slow leave approval queries
**Solution**: Add index on (tenant_id, status, approver_id)

### Issue: Concurrent payroll generation conflicts
**Solution**: Use row-level locking with SELECT FOR UPDATE

### Issue: Employee hierarchy queries timeout
**Solution**: Cache hierarchy or use ltree PostgreSQL extension

### Issue: Multi-tenant data leakage
**Solution**: Always add tenant_id filter in WHERE clause (use middleware)

## Next Development Steps

1. **Service Layer**: Implement business logic (validations, calculations)
2. **API Handlers**: Create controllers for all endpoints
3. **Middleware**: RBAC, auth, logging, error handling
4. **Event Publishing**: Implement NATS/Kafka integration
5. **Frontend**: Build microfrontends for each module
6. **Testing**: Unit and integration tests
7. **Documentation**: Swagger/OpenAPI specs
8. **Deployment**: Docker and K8s configs

---

**Status**: Architecture & Models Complete ✓ 
**Repositories**: 50% Complete
**Services**: Pending
**Handlers**: Pending
**Frontend**: In Progress
