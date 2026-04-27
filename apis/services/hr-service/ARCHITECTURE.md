# Enterprise HRMS Architecture Documentation

## System Overview

This is a comprehensive Human Resource Management System (HRMS) designed for mid-to-large enterprises with millions of users. The system implements microservices architecture with a dedicated HR service handling all HR operations.

## Technology Stack

- **Backend**: Go 1.26.1 with Fiber framework
- **Database**: PostgreSQL with GORM ORM
- **API Style**: REST with JSON
- **Authentication**: JWT-based RBAC
- **Frontend**: Next.js with microfrontends
- **Monorepo**: Turbo-based workspace
- **Message Queue**: NATS/Kafka (for event-driven architecture)
- **Deployment**: Docker + Kubernetes

## Core HR Modules (12 Modules)

### 1. Core HR Module
**Features:**
- Employee master data management (100+ fields)
- Organization structure (hierarchical)
- Department, Location, Job Title management
- Employee hierarchy tracking
- Dependents and emergency contacts
- Education and work experience tracking
- Multi-tenant support

**Key Entities:**
```
Organization → Department → Employee
          ↓
      Location
       ↓
   JobTitle
```

**API Endpoints:**
```
GET    /api/v1/employees           # List all employees
POST   /api/v1/employees           # Create employee
GET    /api/v1/employees/:id       # Get employee details
PATCH  /api/v1/employees/:id       # Update employee
DELETE /api/v1/employees/:id       # Deactivate employee
GET    /api/v1/employees/:id/hierarchy # Get reporting structure
GET    /api/v1/departments         # List departments
POST   /api/v1/departments         # Create department
GET    /api/v1/locations           # List locations
```

### 2. Recruitment & ATS
**Features:**
- Job posting management
- Candidate management with resume parsing
- Multi-round interview scheduling
- Interview feedback and ratings
- Job offer generation and tracking
- Offer approval workflow
- Onboarding task management

**Key Entities:**
```
JobPosting → Candidate → Interview → JobOffer
                    ↓
            OnboardingTask
```

**API Endpoints:**
```
GET    /api/v1/recruitment/job-postings
POST   /api/v1/recruitment/job-postings
GET    /api/v1/recruitment/job-postings/:id/candidates
POST   /api/v1/recruitment/candidates/:id/interviews
POST   /api/v1/recruitment/candidates/:id/offers
GET    /api/v1/recruitment/onboarding-tasks
```

### 3. Attendance & Time Tracking
**Features:**
- Daily attendance with check-in/check-out
- Geolocation tracking
- Shift management
- Biometric integration
- Work from home (WFH) approvals
- Timesheet management with project allocation
- Monthly attendance reports

**Key Entities:**
```
Attendance ← Shift, BiometricData
         ↓
      TimeSheet → TimeSheetEntry
```

**API Endpoints:**
```
POST   /api/v1/attendance/check-in       # Record check-in
POST   /api/v1/attendance/check-out      # Record check-out
GET    /api/v1/attendance/today          # Today's attendance
GET    /api/v1/attendance/history        # Monthly history
GET    /api/v1/attendance/stats          # Dashboard stats
POST   /api/v1/timesheets                # Submit timesheet
GET    /api/v1/timesheets                # Get timesheets
POST   /api/v1/wfh-requests              # Request WFH
```

### 4. Leave Management
**Features:**
- Multiple leave types (Annual, Sick, Casual, Maternity, etc.)
- Leave policies per department/role
- Leave allocation and balance tracking
- Leave application with multi-level approval
- Holiday calendar management
- Leave encashment
- Leave carryforward rules

**Key Entities:**
```
LeaveType → LeavePolicy → LeaveAllocation → Leave → LeaveApproval
     ↓
HolidayCalendar
```

**API Endpoints:**
```
GET    /api/v1/leaves/types              # Available leave types
GET    /api/v1/leaves/balance            # Leave balance
POST   /api/v1/leaves                    # Apply for leave
GET    /api/v1/leaves/pending            # Pending approvals
PATCH  /api/v1/leaves/:id/approve        # Approve leave
GET    /api/v1/leaves/holidays           # Holiday calendar
POST   /api/v1/leaves/encashment         # Encash leave
```

### 5. Payroll (Basic to Advanced)
**Features:**
- Salary structure configuration
- Flexible salary components (Basic, HRA, DA, Allowances, Deductions)
- Formula-based calculations
- Monthly payroll generation
- Payroll approval workflow
- Tax calculations (TDS, IT, surcharge)
- Compliance for multiple countries (India: PAN, PF, ESIC)
- Payroll posting to finance module
- Payment reconciliation

**Key Entities:**
```
SalaryStructure → SalaryComponent
          ↓
EmployeeSalary → SalaryComponentValue
       ↓
    Payroll → PayrollLine
```

**API Endpoints:**
```
GET    /api/v1/payroll/salary-structures
POST   /api/v1/payroll/salary-structures
GET    /api/v1/payroll/employee-salary/:empId
POST   /api/v1/payroll/generate              # Generate monthly payroll
GET    /api/v1/payroll/pending-approval
PATCH  /api/v1/payroll/:id/approve
PATCH  /api/v1/payroll/:id/post
```

### 6. Compensation & Benefits
**Features:**
- Salary reviews and adjustments
- Benefits enrollment
- Insurance claims tracking
- Benefit reconciliation
- Cost to Company (CTC) breakdown
- Compensation statements

**Key Entities:**
```
BenefitPlan → EmployeeBenefit → BenefitClaim
```

**API Endpoints:**
```
GET    /api/v1/benefits/plans
POST   /api/v1/benefits/enrollment
GET    /api/v1/benefits/my-benefits
POST   /api/v1/benefits/claims
GET    /api/v1/benefits/claims/pending
```

### 7. Performance Management
**Features:**
- Goal setting (OKRs, KPIs)
- Goal tracking and progress
- 360-degree feedback
- Performance ratings (Annual/Half-yearly)
- Competency assessments
- Training and development plans
- Performance review workflows

**Key Entities:**
```
PerformanceGoal → GoalProgress
         ↓
PerformanceRating → CompetencyRating
         ↓
EmployeeReview → ReviewQuestion
```

**API Endpoints:**
```
POST   /api/v1/performance/goals
GET    /api/v1/performance/goals
PATCH  /api/v1/performance/goals/:id/progress
POST   /api/v1/performance/reviews
GET    /api/v1/performance/reviews
POST   /api/v1/performance/ratings
```

### 8. Learning Management System (LMS)
**Features:**
- Course creation and management
- Course enrollment
- Progress tracking
- Quizzes and assessments
- Certificate generation
- Learning records store (LRS) integration
- Mandatory vs optional courses

**Key Entities:**
```
LMSCourse → CourseModule
     ↓
CourseEnrollment → LessonProgress
     ↓
Quiz → QuizQuestion → QuizAttempt → Certificate
```

**API Endpoints:**
```
GET    /api/v1/learning/courses
POST   /api/v1/learning/enrollment
GET    /api/v1/learning/my-courses
POST   /api/v1/learning/quizzes/:id/submit
GET    /api/v1/learning/certificates
```

### 9. Employee Self-Service Portal (ESS)
**Features:**
- Personal information updates
- Leave request submission
- WFH request submission
- Expense submission
- Training requests
- Document upload (resume, certs)
- View payslips
- Download tax documents (IT, PF)
- View benefits and insurance

**Base URL:** `/api/v1/ess/`

### 10. Workflow & Approval Engine
**Features:**
- Dynamic approval workflows
- Multi-level approvals
- Conditional approval logic
- Approval delegation
- Timeout and escalation
- Workflow status tracking
- Audit trail

**Key Entities:**
```
ApprovalWorkflow → ApprovalLevel
         ↓
ApprovalRequest → Approval
```

**API Endpoints:**
```
GET    /api/v1/workflows
POST   /api/v1/workflows
GET    /api/v1/approvals/pending
POST   /api/v1/approvals/:id/approve
POST   /api/v1/approvals/:id/delegate
```

### 11. Admin & RBAC
**Features:**
- Role-based access control (Admin, HR, Manager, Employee)
- Granular permissions
- Tenant isolation
- System settings management
- Module configuration
- User role assignment

**Key Entities:**
```
RBACRole ←→ Permission (many-to-many)
    ↓
SystemSetting
```

**API Endpoints:**
```
GET    /api/v1/admin/roles
POST   /api/v1/admin/roles
GET    /api/v1/admin/permissions
POST   /api/v1/admin/role-permissions/:roleId/assign
GET    /api/v1/admin/settings
PATCH  /api/v1/admin/settings
```

### 12. Reporting & Analytics
**Features:**
- Pre-built reports (headcount, attrition, payroll, etc.)
- Custom report builder
- Dashboard with KPIs
- Export to PDF/Excel
- Scheduled report delivery
- Data visualization

**Base URL:** `/api/v1/reports/`

## Database Schema

### Multi-Tenant Architecture

All tables include `tenant_id` field for multi-tenancy:

```sql
-- Example: Employees table
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    -- ... more fields
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, employee_id),
    UNIQUE(tenant_id, email),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id)
);
```

## API Response Format

### Success Response (200 OK)
```json
{
  "success": true,
  "code": "SUCCESS",
  "message": "Operation completed successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP-001",
    "name": "John Doe",
    ...
  },
  "timestamp": "2026-04-27T20:30:00Z"
}
```

### Paginated Response
```json
{
  "success": true,
  "code": "SUCCESS",
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### Error Response (400, 401, 403, 500)
```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ],
  "timestamp": "2026-04-27T20:30:00Z"
}
```

## Event-Driven Architecture

### Events Published to Message Queue (NATS/Kafka)

```
employee.created
employee.updated
employee.terminated

leave.created
leave.approved
leave.rejected
leave.encashed

payroll.generated
payroll.approved
payroll.posted

attendance.checkin
attendance.checkout

performance.review.completed

recruitment.offer.created
recruitment.offer.accepted
```

### Event Payload Example
```json
{
  "event_type": "employee.created",
  "event_id": "evt_12345",
  "timestamp": "2026-04-27T20:30:00Z",
  "tenant_id": "acme_corp",
  "data": {
    "employee_id": "EMP-001",
    "name": "John Doe",
    "email": "john@acme.com",
    "department": "Engineering"
  }
}
```

## Folder Structure

```
erp-microservices/
├── apis/
│   ├── hr-service/
│   │   ├── cmd/
│   │   │   └── api/
│   │   │       └── main.go
│   │   ├── internal/
│   │   │   ├── models/              # 12 modules: employee, organization, recruitment, etc.
│   │   │   ├── repository/          # Data access layer
│   │   │   ├── service/             # Business logic layer
│   │   │   ├── handlers/            # HTTP handlers (controllers)
│   │   │   ├── middleware/          # RBAC, Auth, Logging
│   │   │   ├── database/
│   │   │   ├── config/
│   │   │   └── utils/
│   │   ├── migrations/              # Database migrations
│   │   ├── go.mod
│   │   └── go.sum
│   ├── finance-service/             # For payroll posting
│   ├── api-gateway/
│   └── shared/
│       ├── auth/
│       └── logger/
├── frontend/
│   ├── apps/
│   │   ├── host-app/               # Main dashboard
│   │   ├── hr-mfe/                 # HR module microfrontend
│   │   ├── recruitment-mfe/        # Recruitment module
│   │   ├── payroll-mfe/            # Payroll module
│   │   ├── performance-mfe/        # Performance management
│   │   └── learning-mfe/           # LMS
│   ├── packages/
│   │   ├── ui-components/          # Shared UI components
│   │   ├── api-client/             # API client library
│   │   └── hooks/                  # Custom React hooks
│   ├── pnpm-workspace.yaml
│   └── turbo.json
├── migrations/
│   └── postgres/
│       ├── 001_init_hrms_schema.sql
│       ├── 002_recruitment_tables.sql
│       ├── 003_payroll_tables.sql
│       └── ...
└── docker-compose.yml
```

## Deployment Strategy

### Docker Containers

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: erp_admin
      POSTGRES_PASSWORD: secure_password
      POSTGRES_DB: hr_db
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    ports:
      - "6379:6379"

  nats:
    image: nats:latest
    ports:
      - "4222:4222"

  hr-service:
    build:
      context: ./apis/hr-service
      dockerfile: Dockerfile
    depends_on:
      - postgres
      - nats
    environment:
      DATABASE_URL: postgres://erp_admin:password@postgres:5432/hr_db
      NATS_URL: nats://nats:4222
    ports:
      - "8081:8081"

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
```

### Kubernetes Deployment

```yaml
# k8s/hr-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hr-service
  namespace: hrms
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hr-service
  template:
    metadata:
      labels:
        app: hr-service
    spec:
      containers:
      - name: hr-service
        image: hrms/hr-service:latest
        ports:
        - containerPort: 8081
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: hr-secrets
              key: db-url
        - name: NATS_URL
          value: nats://nats-service:4222
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: hr-service
  namespace: hrms
spec:
  selector:
    app: hr-service
  ports:
  - protocol: TCP
    port: 8081
    targetPort: 8081
  type: ClusterIP
```

## Development Phases

### Phase 1 (MVP) - 2-3 Months
- Core HR (Employee Management)
- Basic Leave Management
- Attendance Tracking
- Simple Payroll

### Phase 2 - 2-3 Months
- Recruitment & ATS
- Performance Management (Goals)
- Employee Self-Service Portal
- RBAC & Workflow Engine

### Phase 3 - 2-3 Months
- Advanced Payroll (Tax, Compliance)
- Compensation & Benefits
- Learning Management System
- Analytics & Reporting

## Common Pitfalls & Solutions

### 1. Payroll Complexity
**Pitfall:** Underestimating salary calculation complexity (different countries, tax rules, statutory compliance)
**Solution:** Use formula-based components, support multiple localization rules, maintain audit trail

### 2. Approval Workflows
**Pitfall:** Hardcoding approval logic instead of making it configurable
**Solution:** Implement workflow engine with rule-based logic, support delegation and timeout

### 3. Data Privacy
**Pitfall:** Not implementing proper RBAC and audit trails
**Solution:** Granular permissions, encrypt sensitive data, audit every change

### 4. Scalability
**Pitfall:** Not designing for multi-tenancy from the start
**Solution:** All tables have tenant_id, use connection pooling, shard at tenant level if needed

### 5. Integration
**Pitfall:** Tightly coupling HR service with other services
**Solution:** Use event-driven architecture, implement service-to-service contracts

## Authentication & Security

### JWT Token Structure
```json
{
  "sub": "user_id",
  "tenant_id": "acme_corp",
  "email": "user@acme.com",
  "roles": ["employee", "manager"],
  "permissions": ["employees.view", "leaves.approve"],
  "iat": 1715100600,
  "exp": 1715187000
}
```

### RBAC Matrix

| Role | Employees | Leave | Payroll | Performance | Admin |
|------|-----------|-------|---------|-------------|-------|
| Admin | CRUD | CRUD | CRUD | CRUD | CRUD |
| HR Manager | R | CRUD | RU | RU | R |
| Manager | R | Approve | R | R | - |
| Employee | R (self) | Create | R (self) | R | - |

## Monitoring & Logging

### Key Metrics
- API response times (p50, p95, p99)
- Error rates by endpoint
- Payroll processing time
- Leave approval SLA
- Attendance check-in success rate
- Database query performance

### Logs to Track
- Employee CRUD operations
- Leave approvals/rejections
- Payroll processing
- Authentication failures
- Approval workflow status changes

## Next Steps

1. **Database Migrations**: Create SQL migration scripts for all tables
2. **Service Layer**: Implement business logic (validations, calculations, workflows)
3. **API Handlers**: Create HTTP handlers for all endpoints
4. **Event Publishing**: Integrate NATS for event-driven architecture
5. **Frontend**: Build microfrontends for each module
6. **Testing**: Unit, integration, and E2E tests
7. **Documentation**: API docs (Swagger/OpenAPI)
8. **Deployment**: Docker and Kubernetes setup

---

*Architecture Document - Enterprise HRMS System*
*Last Updated: April 27, 2026*
