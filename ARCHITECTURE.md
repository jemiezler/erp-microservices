# Enterprise HRMS Architecture - System Design Document

**Version**: 1.0  
**Date**: April 2026  
**Target Users**: Mid-to-large enterprises (1,000 - 100,000+ employees)  
**Comparable Systems**: Workday, SuccessFactors, Bamboo HR

---

## Table of Contents
1. [System Overview](#system-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Core Principles](#core-principles)
4. [Microservices Architecture](#microservices-architecture)
5. [Module Specifications](#module-specifications)
6. [Database Design](#database-design)
7. [API Design](#api-design)
8. [Event-Driven Architecture](#event-driven-architecture)
9. [Microfrontend Strategy](#microfrontend-strategy)
10. [Security & RBAC](#security--rbac)
11. [Deployment Strategy](#deployment-strategy)
12. [Development Phases](#development-phases)
13. [Common Pitfalls](#common-pitfalls)

---

## System Overview

### High-Level Vision
Enterprise HRMS supporting:
- **Multi-tenant SaaS** with complete data isolation
- **11 core HR modules** with deep feature sets
- **Global compliance** (GDPR, local payroll laws)
- **Millions of users** with sub-100ms API responses
- **Role-based access** with granular permissions
- **Workflow automation** with complex approval chains
- **Real-time analytics** and reporting

### Key Metrics
- **Expected Scale**: 5,000 - 100,000 employees per tenant
- **Daily Active Users**: 20-30% of total employees
- **Peak Requests**: 10,000 req/s during payroll processing
- **Data Volume**: ~500GB - 5TB per enterprise tenant
- **Availability Target**: 99.95% SLA

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     CLIENT LAYER (Web/Mobile)                   │
├─────────────────────────────────────────────────────────────────┤
│  Host App (Shell) │ Dashboard │ ESS │ Analytics │ Admin Console │
└─────────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────────┐
│                  API GATEWAY (Kong/Envoy)                       │
│  - Rate limiting, Auth token validation, Request routing        │
│  - Circuit breaker, Caching (Redis)                             │
└─────────────────────────────────────────────────────────────────┘
                           ↓
┌──────────────────────────────────────────────────────────────────┐
│                    MICROSERVICES TIER                            │
├──────────────────────────────────────────────────────────────────┤
│ Core-HR │ Recruitment │ Attendance │ Leave │ Payroll │ Performance
│ Learning │ Compensation │ ESS │ Admin │ Workflow │ Analytics    │
├──────────────────────────────────────────────────────────────────┤
│                    Service Mesh (Istio)                          │
│              - Distributed tracing, Circuit breaking             │
│              - Mutual TLS (mTLS), Service-to-service auth        │
└──────────────────────────────────────────────────────────────────┘
                           ↓
┌──────────────────────────────────────────────────────────────────┐
│                   DATA & MESSAGING TIER                          │
├──────────────────────────────────────────────────────────────────┤
│ PostgreSQL (OLTP) │ TimescaleDB (Analytics) │ Redis (Cache)      │
│ Elasticsearch (Search) │ NATS/Kafka (Event Bus)                  │
└──────────────────────────────────────────────────────────────────┘
                           ↓
┌──────────────────────────────────────────────────────────────────┐
│              SUPPORTING SERVICES                                 │
├──────────────────────────────────────────────────────────────────┤
│ Auth Service │ Notification │ File Storage │ Audit │ Search      │
└──────────────────────────────────────────────────────────────────┘

TENANT ISOLATION:
- Row-Level Security (RLS) in PostgreSQL per tenant
- Separate Redis namespaces per tenant
- Tenant context in every request header
```

---

## Core Principles

### 1. **Microservices Boundaries**
- **By Business Capability**: Each service owns one business domain
- **Loose Coupling**: Services communicate via async events (preferred) or gRPC
- **High Cohesion**: Related features live in same service
- **Database per Service**: No direct DB access between services

### 2. **Data Consistency**
- **Eventual Consistency** for cross-service operations
- **Saga Pattern** for distributed transactions (e.g., hire → create email → add to org)
- **Event Sourcing** optional for audit trail in finance modules
- **ACID** within single service database

### 3. **Scalability**
- **Horizontal scaling**: Stateless services behind load balancers
- **Database partitioning**: Tenant-based sharding for very large deployments
- **Caching strategy**: Multi-level (HTTP, Redis, CDN)
- **Batch processing**: Async jobs for bulk operations (payroll, reports)

### 4. **Security**
- **JWT tokens** for API authentication (short-lived, refresh tokens)
- **RBAC + ABAC**: Role-based + Attribute-based access control
- **mTLS**: Service-to-service authentication
- **Encryption**: At-rest (DB) and in-transit (TLS)
- **Audit everything**: All state changes logged

---

## Microservices Architecture

### Service Topology

```
TIER 1: CORE SERVICES (Always needed)
├── Auth Service
│   └── Handles JWT, OAuth2, SAML, MFA
├── API Gateway
│   └── Request routing, rate limiting, caching
└── Tenant Service
    └── Multi-tenant context, onboarding

TIER 2: DOMAIN SERVICES (Business logic)
├── Core-HR Service
│   └── Employees, org structure, master data
├── Recruitment Service
│   └── Job requisitions, candidates, hiring workflows
├── Attendance Service
│   └── Check-in/out, shifts, timesheets
├── Leave Service
│   └── Leave policies, applications, approvals
├── Payroll Service
│   └── Salary structure, processing, tax compliance
├── Performance Service
│   └── Goals, reviews, feedback, calibration
├── Learning Service
│   └── Courses, certifications, tracking
├── Compensation Service
│   └── Bonus, equity, variable pay, budgets
├── ESS Service
│   └── Employee self-service portal (read-heavy)
└── Workflow Service
    └── Approval engine, state machines

TIER 3: CROSS-CUTTING SERVICES
├── Notification Service
│   └── Email, SMS, push notifications
├── File Service
│   └── Document storage, version control
├── Analytics Service
│   └── Data warehouse, reporting
├── Audit Service
│   └── Change logs, compliance, forensics
└── Search Service
    └── Employee directory, document search
```

---

## Module Specifications

### MODULE 1: CORE HR

**Purpose**: Central employee repository and org structure management

#### Key Features
- Employee master data (personal, contact, emergency)
- Organizational hierarchy with reporting lines
- Job classification and career leveling
- Competency management
- Employment contracts and terms
- Separation/exit management
- Bulk employee import/sync

#### Database Schema

```sql
-- Multi-tenant employee table
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    personal_email VARCHAR(255),
    
    -- Employment details
    employment_status ENUM('Active', 'OnLeave', 'Terminated', 'OnNotice'),
    employment_type ENUM('FullTime', 'PartTime', 'Contract', 'Intern'),
    hire_date DATE NOT NULL,
    termination_date DATE,
    termination_reason VARCHAR(500),
    
    -- Organization
    department_id UUID NOT NULL REFERENCES departments(id),
    job_id UUID NOT NULL REFERENCES jobs(id),
    manager_id BIGINT REFERENCES employees(id),
    cost_center_id UUID,
    
    -- Identity
    date_of_birth DATE,
    gender ENUM('Male', 'Female', 'Other', 'NotSpecified'),
    nationality VARCHAR(100),
    personal_id_number VARCHAR(50), -- SSN, Aadhaar, etc.
    
    -- Contact
    phone_number VARCHAR(20),
    alternate_phone VARCHAR(20),
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    
    CONSTRAINT unique_employee_per_tenant UNIQUE(tenant_id, employee_id),
    INDEX idx_tenant_status (tenant_id, employment_status),
    INDEX idx_department (department_id),
    INDEX idx_manager (manager_id)
);

-- Organizational structure
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    parent_department_id UUID REFERENCES departments(id),
    department_head_id BIGINT REFERENCES employees(id),
    description TEXT,
    cost_center_id VARCHAR(50),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_parent (tenant_id, parent_department_id)
);

CREATE TABLE org_hierarchy (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    manager_id BIGINT REFERENCES employees(id),
    level INTEGER, -- 0=CEO, 1=C-level, 2=VP, etc.
    hierarchy_path TEXT, -- e.g., "1/12/456" for ancestry queries
    effective_from DATE,
    effective_to DATE,
    CONSTRAINT no_self_manager CHECK (employee_id != manager_id),
    INDEX idx_hierarchy (tenant_id, manager_id)
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    job_title VARCHAR(255) NOT NULL,
    job_code VARCHAR(50),
    job_family VARCHAR(100),
    job_level VARCHAR(50),
    description TEXT,
    min_salary DECIMAL(15,2),
    max_salary DECIMAL(15,2),
    reports_to_job_id UUID REFERENCES jobs(id),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE employment_contracts (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    contract_type ENUM('Permanent', 'Fixed-Term', 'Probation'),
    start_date DATE NOT NULL,
    end_date DATE,
    notice_period_days INTEGER,
    terms_conditions TEXT,
    document_url VARCHAR(500),
    created_at TIMESTAMP,
    INDEX idx_employee_contract (tenant_id, employee_id)
);
```

#### API Design (gRPC + REST)

```protobuf
// proto/core_hr/v1/employees.proto
service EmployeeService {
    // Query operations
    rpc GetEmployee(GetEmployeeRequest) returns (Employee);
    rpc ListEmployees(ListEmployeesRequest) returns (ListEmployeesResponse);
    rpc GetOrgChart(GetOrgChartRequest) returns (OrgChartResponse);
    
    // Write operations
    rpc CreateEmployee(CreateEmployeeRequest) returns (Employee);
    rpc UpdateEmployee(UpdateEmployeeRequest) returns (Employee);
    rpc TerminateEmployee(TerminateEmployeeRequest) returns (Employee);
    
    // Bulk operations
    rpc BulkImportEmployees(stream BulkImportRequest) returns (BulkImportResponse);
}

message Employee {
    string employee_id = 1;
    string first_name = 2;
    string last_name = 3;
    string email = 4;
    EmploymentStatus employment_status = 5;
    string department_id = 6;
    string job_id = 7;
    string manager_id = 8;
    google.protobuf.Timestamp hire_date = 9;
    google.protobuf.Timestamp created_at = 10;
}
```

#### REST Endpoints

```
GET    /api/v1/employees                           - List with filters
GET    /api/v1/employees/{id}                      - Get employee details
POST   /api/v1/employees                           - Create employee
PATCH  /api/v1/employees/{id}                      - Update employee
DELETE /api/v1/employees/{id}                      - Terminate (soft delete)
GET    /api/v1/org-chart?department_id={id}        - Get org hierarchy
GET    /api/v1/employees/{id}/reports              - Direct reports
POST   /api/v1/employees/bulk-import               - Bulk upload
GET    /api/v1/search/employees?q={query}          - Search employees
```

#### Events Published

```
✓ EmployeeCreated
  {employee_id, email, department_id, manager_id, hire_date}
  → Consumed by: ESS, Payroll, Attendance, Notification

✓ EmployeeUpdated
  {employee_id, changed_fields, timestamp}

✓ EmployeeTerminated
  {employee_id, termination_date, reason}
  → Consumed by: Payroll, ESS, Notification, Learning

✓ DepartmentChanged
  {department_id, change_type, affected_employees_count}

✓ OrganizationHierarchyChanged
  {changes: [{employee_id, old_manager, new_manager}]}
```

---

### MODULE 2: RECRUITMENT (ATS)

**Purpose**: End-to-end hiring process management

#### Key Features
- Job requisitions and publishing
- Candidate pipeline management
- Interview scheduling and feedback
- Offer management
- Background check integration
- Recruitment analytics
- Career portal (public-facing)

#### Database Schema (Abbreviated)

```sql
CREATE TABLE job_requisitions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    job_id UUID NOT NULL REFERENCES jobs(id),
    department_id UUID NOT NULL,
    requisition_code VARCHAR(50) UNIQUE,
    headcount_needed INTEGER,
    status ENUM('Draft', 'Approved', 'Open', 'Filled', 'Cancelled'),
    approval_chain_id UUID, -- linked to workflow
    created_by UUID,
    created_at TIMESTAMP
);

CREATE TABLE candidates (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    source ENUM('LinkedIn', 'Indeed', 'Referral', 'Career_Site', 'Direct_Apply'),
    resume_url VARCHAR(500),
    created_at TIMESTAMP,
    INDEX idx_tenant_email (tenant_id, email)
);

CREATE TABLE candidate_applications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    requisition_id UUID NOT NULL REFERENCES job_requisitions(id),
    stage ENUM('Applied', 'Screening', 'Interview', 'Offer', 'Hired', 'Rejected'),
    stage_entered_at TIMESTAMP,
    score_by_recruiter DECIMAL(3,1),
    score_by_hiring_manager DECIMAL(3,1),
    applied_at TIMESTAMP,
    INDEX idx_stage (stage, tenant_id)
);

CREATE TABLE interview_rounds (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    application_id UUID NOT NULL REFERENCES candidate_applications(id),
    round_number INTEGER,
    interview_type ENUM('Phone', 'Video', 'OnSite', 'Panel'),
    interviewer_id BIGINT REFERENCES employees(id),
    scheduled_at TIMESTAMP,
    feedback TEXT,
    rating ENUM('Strong_Yes', 'Yes', 'Maybe', 'No'),
    duration_minutes INTEGER
);

CREATE TABLE offers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    application_id UUID NOT NULL,
    candidate_id UUID NOT NULL,
    job_id UUID NOT NULL,
    base_salary DECIMAL(15,2),
    signing_bonus DECIMAL(15,2),
    stock_grants DECIMAL(15,2),
    benefits_package TEXT,
    offer_letter_url VARCHAR(500),
    start_date DATE,
    status ENUM('Draft', 'Sent', 'Accepted', 'Rejected', 'Expired'),
    sent_at TIMESTAMP,
    expires_at TIMESTAMP,
    accepted_at TIMESTAMP
);
```

#### API Endpoints

```
POST   /api/v1/requisitions                        - Create job req
GET    /api/v1/requisitions                        - List reqs
POST   /api/v1/candidates                          - Add candidate
GET    /api/v1/candidates/{id}/applications        - Get applications
POST   /api/v1/applications/{id}/advance           - Move to next stage
POST   /api/v1/interviews                          - Schedule interview
POST   /api/v1/offers                              - Create offer
PATCH  /api/v1/offers/{id}/send                    - Send offer letter
GET    /api/v1/analytics/hiring-funnel             - Recruitment metrics
```

#### Events

```
✓ RequisitionCreated → Notification (post to career portal)
✓ CandidateApplied → Recruiter notification
✓ ApplicationStageChanged → Candidate notification, Hiring manager alert
✓ InterviewScheduled → Calendar sync, Notification
✓ OfferAccepted → Core-HR triggers EmployeeCreated event
✓ ApplicationRejected → CRM cleanup
```

---

### MODULE 3: ATTENDANCE & TIME TRACKING

**Purpose**: Shift management, biometric tracking, timesheet management

#### Key Features
- Multiple shift types (9-5, rotating, on-call)
- Biometric device integration
- Real-time geolocation tracking
- Timesheet approval workflows
- Attendance rules and exceptions
- Auto-overtime calculation
- Attendance analytics (leaderboards, patterns)

#### Database Schema (Abbreviated)

```sql
CREATE TABLE shifts (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100),
    start_time TIME,
    end_time TIME,
    break_duration_minutes INTEGER,
    work_hours_per_day DECIMAL(4,2),
    created_at TIMESTAMP
);

CREATE TABLE employee_shifts (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    shift_id UUID NOT NULL REFERENCES shifts(id),
    effective_from DATE,
    effective_to DATE,
    INDEX idx_employee_shifts (tenant_id, employee_id, effective_from)
);

CREATE TABLE attendance_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL,
    check_in_time TIMESTAMP NOT NULL,
    check_out_time TIMESTAMP,
    device_id VARCHAR(100), -- biometric device ID
    location_lat DECIMAL(10,8),
    location_lon DECIMAL(11,8),
    check_in_method ENUM('Biometric', 'Mobile', 'Web', 'Manual'),
    status ENUM('Present', 'Late', 'EarlyLeave', 'Absent', 'WFH'),
    remarks TEXT,
    approved_by BIGINT REFERENCES employees(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP,
    INDEX idx_employee_date (tenant_id, employee_id, DATE(check_in_time))
);

CREATE TABLE timesheets (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL,
    week_start_date DATE,
    week_end_date DATE,
    total_hours DECIMAL(6,2),
    overtime_hours DECIMAL(6,2),
    status ENUM('Draft', 'Submitted', 'Approved', 'Rejected'),
    submitted_at TIMESTAMP,
    approved_by BIGINT,
    approved_at TIMESTAMP,
    approver_comments TEXT,
    created_at TIMESTAMP,
    UNIQUE(tenant_id, employee_id, week_start_date)
);

CREATE TABLE timesheet_entries (
    id UUID PRIMARY KEY,
    timesheet_id UUID NOT NULL REFERENCES timesheets(id),
    day_of_week SMALLINT, -- 0=Monday, 6=Sunday
    hours_worked DECIMAL(4,2),
    project_id VARCHAR(50), -- billable project
    task_description TEXT,
    notes TEXT
);
```

#### REST Endpoints

```
POST   /api/v1/attendance/check-in                 - Clock in
POST   /api/v1/attendance/check-out                - Clock out
GET    /api/v1/attendance/{employee_id}            - Daily attendance
GET    /api/v1/attendance/report                   - Attendance report
POST   /api/v1/timesheets                          - Submit timesheet
GET    /api/v1/timesheets/{id}                     - Get timesheet
PATCH  /api/v1/timesheets/{id}/approve             - Approve timesheet
GET    /api/v1/shifts                              - Get available shifts
POST   /api/v1/employee-shifts                     - Assign shift to employee
```

#### Events

```
✓ CheckInRecorded → Update employee status in real-time
✓ CheckOutRecorded → Calculate work hours, detect late checkout
✓ TimesheetSubmitted → Route to manager approval
✓ AttendanceAnomaly → Late arrival, early leave (for analytics)
✓ TimesheetApproved → Payroll consumption
```

---

### MODULE 4: LEAVE MANAGEMENT

**Purpose**: Leave policy enforcement, applications, and approvals

#### Key Features
- Multiple leave types (PTO, Sick, Sabbatical, Unpaid)
- Carry-over policies
- Company holidays
- Approval workflows
- Leave balance tracking
- Team view for managers
- Compliance with local laws (e.g., France 30 days minimum)

#### Database Schema

```sql
CREATE TABLE leave_types (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100), -- "Annual Leave", "Sick Leave", etc.
    code VARCHAR(20),
    description TEXT,
    max_accrual_per_year DECIMAL(6,2),
    carry_forward_limit DECIMAL(6,2),
    gender_based BOOLEAN, -- e.g., maternity leave only for women
    approval_required BOOLEAN,
    country_applicable VARCHAR(100), -- e.g., "IN", "US"
    created_at TIMESTAMP
);

CREATE TABLE leave_balances (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    leave_type_id UUID NOT NULL REFERENCES leave_types(id),
    financial_year_start DATE,
    financial_year_end DATE,
    opening_balance DECIMAL(6,2),
    accrued DECIMAL(6,2),
    utilized DECIMAL(6,2),
    carried_forward DECIMAL(6,2),
    lapsed DECIMAL(6,2),
    closing_balance DECIMAL(6,2),
    last_updated TIMESTAMP,
    UNIQUE(tenant_id, employee_id, leave_type_id, financial_year_start)
);

CREATE TABLE leave_applications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    leave_type_id UUID NOT NULL REFERENCES leave_types(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    duration_days DECIMAL(4,1),
    reason TEXT,
    attachment_url VARCHAR(500),
    status ENUM('Draft', 'Submitted', 'Approved', 'Rejected', 'Cancelled'),
    approver_id BIGINT REFERENCES employees(id),
    approval_date TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP,
    INDEX idx_employee_dates (tenant_id, employee_id, start_date)
);

CREATE TABLE company_holidays (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    holiday_date DATE,
    holiday_name VARCHAR(100),
    country_code VARCHAR(10), -- For multi-country
    applicable_departments TEXT[], -- NULL = all departments
    created_at TIMESTAMP
);

CREATE TABLE leave_accrual_rules (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    leave_type_id UUID NOT NULL REFERENCES leave_types(id),
    employment_type ENUM('FullTime', 'PartTime', 'Contract'),
    accrual_frequency ENUM('Monthly', 'Quarterly', 'Annually'),
    accrual_amount DECIMAL(6,2),
    vesting_period_months INTEGER,
    effective_from DATE,
    effective_to DATE
);
```

#### API Endpoints

```
GET    /api/v1/leave-types                         - Get leave types
GET    /api/v1/leave-balances/{employee_id}        - Check balance
POST   /api/v1/leave-applications                  - Apply for leave
GET    /api/v1/leave-applications                  - Pending approvals
PATCH  /api/v1/leave-applications/{id}/approve     - Approve leave
GET    /api/v1/team-calendar?manager_id={id}       - Team leave calendar
GET    /api/v1/holidays                            - Company holidays
POST   /api/v1/holidays                            - Add holiday (Admin)
```

#### Events

```
✓ LeaveApplicationSubmitted → Manager notification
✓ LeaveApplicationApproved → ESS update, Calendar event
✓ LeaveApplicationRejected → Employee notification
✓ LeaveUtilized → Deduct from balance, Attendance marking
✓ LeavesExpiring → Reminder notification (90 days before year end)
✓ AccrualProcessed → Annual leave accrual batch job
```

---

### MODULE 5: PAYROLL (Advanced)

**Purpose**: Salary structure, processing, compliance, tax calculations

#### Key Features
- Salary structure with fixed/variable components
- Earnings (basic, HRA, DA, allowances)
- Deductions (PF, ESI, IT, loans)
- Tax compliance (multiple countries)
- Statutory benefits (provident fund, gratuity)
- Mid-month salary advances
- Payroll runs and approval
- Salary slips and disbursement
- Reimbursement tracking
- Stock option vesting

#### Database Schema (Critical)

```sql
CREATE TABLE salary_structures (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100), -- "Standard", "Senior Manager", etc.
    employee_id BIGINT REFERENCES employees(id), -- NULL for template
    employment_type ENUM('FullTime', 'PartTime', 'Contract'),
    country VARCHAR(10), -- "IN", "US", "GB"
    currency VARCHAR(3),
    base_salary DECIMAL(15,2),
    effective_from DATE,
    effective_to DATE,
    is_template BOOLEAN,
    created_at TIMESTAMP,
    INDEX idx_employee_salary (tenant_id, employee_id, effective_from)
);

CREATE TABLE salary_components (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100), -- "Basic", "HRA", "Bonus", "PF"
    component_type ENUM('Earning', 'Deduction', 'Statutory'),
    category ENUM('Fixed', 'Variable', 'Arrear', 'Advance'),
    is_taxable BOOLEAN,
    is_pf_applicable BOOLEAN,
    calculation_method ENUM('Fixed', 'Percentage', 'Formula'), -- formula = "Basic * 0.5"
    formula_expression TEXT, -- For complex calculations
    max_limit DECIMAL(15,2), -- e.g., HRA max limit
    min_limit DECIMAL(15,2),
    rounding_method ENUM('None', 'Round', 'RoundDown', 'RoundUp'),
    created_at TIMESTAMP
);

CREATE TABLE salary_structure_assignments (
    id UUID PRIMARY KEY,
    salary_structure_id UUID NOT NULL REFERENCES salary_structures(id),
    component_id UUID NOT NULL REFERENCES salary_components(id),
    component_amount DECIMAL(15,2),
    percentage_of_base DECIMAL(5,2), -- NULL if fixed amount
    order_index INTEGER, -- For display order
    created_at TIMESTAMP
);

CREATE TABLE payroll_runs (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    payroll_period_id UUID,
    month_year DATE, -- e.g., 2026-03-01
    payroll_cycle ENUM('Monthly', 'Bi-Weekly', 'Weekly'),
    status ENUM('Draft', 'Submitted', 'Reviewed', 'Approved', 'Processed', 'Disbursed'),
    total_employees INTEGER,
    total_earnings DECIMAL(18,2),
    total_deductions DECIMAL(18,2),
    net_payroll DECIMAL(18,2),
    processing_date DATE,
    disbursement_date DATE,
    submitted_by UUID,
    submitted_at TIMESTAMP,
    approved_by UUID,
    approved_at TIMESTAMP,
    created_at TIMESTAMP,
    INDEX idx_payroll_status (tenant_id, status, month_year)
);

CREATE TABLE employee_salary_details (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    payroll_run_id UUID NOT NULL REFERENCES payroll_runs(id),
    salary_structure_id UUID NOT NULL REFERENCES salary_structures(id),
    
    -- Components calculated
    component_id UUID NOT NULL REFERENCES salary_components(id),
    calculated_amount DECIMAL(15,2),
    remarks TEXT,
    
    -- Leave deduction
    leave_without_pay_days INTEGER,
    leave_deduction DECIMAL(15,2),
    
    -- Overtime
    overtime_hours DECIMAL(6,2),
    overtime_rate DECIMAL(15,2),
    overtime_amount DECIMAL(15,2),
    
    created_at TIMESTAMP,
    INDEX idx_employee_payroll (tenant_id, employee_id, payroll_run_id)
);

CREATE TABLE salary_slips (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    payroll_run_id UUID NOT NULL REFERENCES payroll_runs(id),
    salary_month DATE,
    gross_earnings DECIMAL(15,2),
    gross_deductions DECIMAL(15,2),
    net_pay DECIMAL(15,2),
    working_days INTEGER,
    present_days INTEGER,
    paid_days DECIMAL(5,2),
    
    -- Tax details
    income_tax_deducted DECIMAL(15,2),
    income_tax_regime ENUM('Old', 'New'),
    pan_number VARCHAR(20),
    
    -- PF details
    employee_pf_contribution DECIMAL(15,2),
    employer_pf_contribution DECIMAL(15,2),
    pf_balance DECIMAL(15,2),
    
    -- Statutory
    esi_deducted DECIMAL(15,2),
    
    -- Payment
    payment_method ENUM('Bank', 'Check', 'Cash'),
    bank_account_id UUID,
    payment_date DATE,
    payment_reference VARCHAR(100),
    
    created_at TIMESTAMP,
    INDEX idx_employee_slip (tenant_id, employee_id, salary_month)
);

CREATE TABLE tax_regimes (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    country VARCHAR(10),
    regime_name VARCHAR(100), -- "Old Regime", "New Regime"
    financial_year_start DATE,
    financial_year_end DATE,
    slabs JSONB, -- [{threshold: 250000, rate: 0.05}, ...]
    deductions JSONB, -- {standard: 50000, hra_limit: 1500000}
    surcharge_applicable BOOLEAN,
    cess_rate DECIMAL(5,2),
    created_at TIMESTAMP
);

CREATE TABLE payroll_taxes (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL,
    financial_year DATE,
    gross_income DECIMAL(15,2),
    income_tax DECIMAL(15,2),
    surcharge DECIMAL(15,2),
    cess DECIMAL(15,2),
    total_tax DECIMAL(15,2),
    tax_regime_id UUID REFERENCES tax_regimes(id),
    filed_status ENUM('Not_Filed', 'Filed', 'Approved', 'Amended'),
    created_at TIMESTAMP
);

CREATE TABLE statutory_benefits (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    benefit_type ENUM('PF', 'Gratuity', 'ESI', 'LWOP'),
    opening_balance DECIMAL(15,2),
    current_year_contribution DECIMAL(15,2),
    withdrawals DECIMAL(15,2),
    closing_balance DECIMAL(15,2),
    financial_year DATE,
    account_number VARCHAR(50),
    created_at TIMESTAMP,
    INDEX idx_employee_benefit (tenant_id, employee_id, benefit_type)
);

CREATE TABLE reimbursements (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    expense_category VARCHAR(50), -- Travel, Meals, Supplies
    expense_date DATE,
    amount DECIMAL(15,2),
    description TEXT,
    receipt_url VARCHAR(500),
    status ENUM('Draft', 'Submitted', 'Approved', 'Rejected', 'Reimbursed'),
    approved_by BIGINT REFERENCES employees(id),
    approved_at TIMESTAMP,
    reimbursement_date DATE,
    created_at TIMESTAMP,
    INDEX idx_employee_reimbursement (tenant_id, employee_id, status)
);
```

#### REST Endpoints

```
GET    /api/v1/salary-structures/{employee_id}     - Get salary structure
POST   /api/v1/payroll-runs                        - Create payroll run
GET    /api/v1/payroll-runs                        - List payroll runs
POST   /api/v1/payroll-runs/{id}/calculate         - Calculate salaries
GET    /api/v1/payroll-runs/{id}/details           - Payroll details
PATCH  /api/v1/payroll-runs/{id}/approve           - Approve payroll
POST   /api/v1/payroll-runs/{id}/disburse          - Disburse salaries
GET    /api/v1/salary-slips/{employee_id}          - Get salary slip
GET    /api/v1/tax-calculations/{employee_id}      - Tax summary
POST   /api/v1/reimbursements                      - Submit expense
PATCH  /api/v1/reimbursements/{id}/approve         - Approve reimbursement
```

#### Key Business Logic

```go
// Payroll Calculation Engine (pseudocode)
CalculateSalary(employee, payrollRun):
  1. Get salary structure for employee
  2. For each earning component:
     - Calculate amount (fixed/percentage/formula)
     - Apply limits (min/max)
  3. For each deduction:
     - Calculate (can reference earnings)
  4. Deduct leave without pay
  5. Add overtime
  6. Calculate taxes (PF, IT, ESI, etc.)
  7. Calculate net pay
  8. Create salary slip entry
  9. Emit SalaryCalculated event
  
TaxCalculation(employee, financialYear):
  1. Sum all taxable earnings
  2. Apply deductions (standard, HRA, etc.)
  3. Look up tax regime
  4. Calculate IT using progressive slabs
  5. Calculate surcharge, cess
  6. Handle TDS adjustments
  7. Update payroll_taxes table
```

#### Events

```
✓ PayrollRunCreated → Notification to CFO/HR
✓ SalaryCalculated → Validation engine checks anomalies
✓ PayrollApproved → Finance notification
✓ PaymentProcessed → Bank integration trigger
✓ SalarySlipGenerated → ESS notification, PDF generation
✓ TaxFilingRequired → Notification (year-end)
✓ ReimbursementApproved → Accounting integration
```

---

### MODULE 6: PERFORMANCE MANAGEMENT

**Purpose**: Goal setting, reviews, feedback, calibration

#### Key Features
- Goal management with OKRs
- 360-degree feedback
- Performance reviews (annual, mid-year, ad-hoc)
- Calibration sessions
- Performance ratings and distributions
- Development plans
- Competency assessments
- Performance trends and analytics

#### Database Schema (Abbreviated)

```sql
CREATE TABLE goals (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    goal_title VARCHAR(255),
    goal_description TEXT,
    goal_type ENUM('Business', 'Development', 'Personal'),
    alignment_parent_goal_id UUID REFERENCES goals(id), -- For alignment to higher-level goals
    start_date DATE,
    end_date DATE,
    target_value DECIMAL(15,2),
    achieved_value DECIMAL(15,2),
    progress_percentage DECIMAL(5,2),
    status ENUM('Draft', 'Active', 'Completed', 'Abandoned'),
    owner_id BIGINT NOT NULL REFERENCES employees(id), -- Usually employee's manager
    created_at TIMESTAMP,
    INDEX idx_employee_goals (tenant_id, employee_id, status)
);

CREATE TABLE goal_check_ins (
    id UUID PRIMARY KEY,
    goal_id UUID NOT NULL REFERENCES goals(id),
    check_in_date DATE,
    progress_percentage DECIMAL(5,2),
    comments TEXT,
    updated_by BIGINT NOT NULL REFERENCES employees(id),
    created_at TIMESTAMP
);

CREATE TABLE performance_reviews (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    review_cycle_id UUID, -- Links to a specific review cycle (e.g., "2025 Annual")
    review_period_start DATE,
    review_period_end DATE,
    manager_id BIGINT NOT NULL REFERENCES employees(id),
    review_status ENUM('Draft', 'Submitted', 'In_Progress', 'Completed'),
    overall_rating ENUM('Exceeds_Expectations', 'Meets_Expectations', 'Needs_Improvement', 'Unsatisfactory'),
    manager_comments TEXT,
    submitted_at TIMESTAMP,
    completed_at TIMESTAMP,
    INDEX idx_review_status (tenant_id, employee_id, review_status)
);

CREATE TABLE review_competencies (
    id UUID PRIMARY KEY,
    review_id UUID NOT NULL REFERENCES performance_reviews(id),
    competency_id UUID NOT NULL,
    rating ENUM('1_Below', '2_Developing', '3_Proficient', '4_Advanced', '5_Expert'),
    comments TEXT,
    rater_id BIGINT REFERENCES employees(id),
    rater_type ENUM('Manager', 'Peer', 'Direct_Report', 'Self'),
    created_at TIMESTAMP
);

CREATE TABLE feedback (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    feedback_from_id BIGINT NOT NULL REFERENCES employees(id),
    feedback_to_id BIGINT NOT NULL REFERENCES employees(id),
    feedback_type ENUM('Appreciation', 'Coaching', 'Request', 'General'),
    message TEXT,
    is_anonymous BOOLEAN,
    sentiment ENUM('Positive', 'Neutral', 'Negative'),
    created_at TIMESTAMP,
    INDEX idx_feedback_to (tenant_id, feedback_to_id, created_at)
);

CREATE TABLE calibration_sessions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    review_cycle_id UUID,
    session_date DATE,
    participants_count INTEGER,
    department_id UUID,
    normalization_curve JSONB, -- {exceeds: 10%, meets: 70%, needs_improvement: 20%}
    finalized BOOLEAN,
    created_at TIMESTAMP
);
```

#### API Endpoints

```
POST   /api/v1/goals                               - Create goal
GET    /api/v1/employees/{id}/goals                - Get employee goals
PATCH  /api/v1/goals/{id}                          - Update goal
POST   /api/v1/goals/{id}/check-in                 - Add goal check-in
POST   /api/v1/reviews                             - Create performance review
GET    /api/v1/reviews/{id}                        - Get review details
POST   /api/v1/reviews/{id}/submit                 - Submit review
POST   /api/v1/feedback                            - Give feedback
GET    /api/v1/calibration-sessions                - List calibration sessions
POST   /api/v1/calibration-sessions/{id}/finalize  - Finalize ratings
GET    /api/v1/analytics/performance-trends        - Performance data
```

#### Events

```
✓ GoalCreated → Notification to employee & manager
✓ GoalProgress Updated → Real-time updates
✓ ReviewCycleStarted → All employees notified
✓ ReviewSubmitted → Manager receives notification
✓ FeedbackReceived → Recipient notification
✓ CalibrationCompleted → Rating distribution locked
✓ RatingsFinalized → Linked to compensation decisions
```

---

### MODULE 7: LEARNING & DEVELOPMENT (LMS)

**Purpose**: Training, certifications, skill development

#### Key Features
- Course catalog with internal/external courses
- Learning paths for roles
- Enrollment and tracking
- Certification management
- Learning outcomes and assessments
- Compliance training tracking
- Learning analytics

#### Database Schema (Abbreviated)

```sql
CREATE TABLE courses (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    course_code VARCHAR(50),
    course_name VARCHAR(255),
    description TEXT,
    course_type ENUM('Online', 'Classroom', 'Blended', 'Self-Paced'),
    duration_hours DECIMAL(6,2),
    provider VARCHAR(100),
    external_url VARCHAR(500),
    audience VARCHAR(500), -- JSON array of roles
    is_mandatory BOOLEAN,
    completion_required_for_role BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE enrollments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    course_id UUID NOT NULL REFERENCES courses(id),
    enrollment_date DATE,
    completion_date DATE,
    status ENUM('Not_Started', 'In_Progress', 'Completed', 'Failed', 'Dropped'),
    progress_percentage DECIMAL(5,2),
    score DECIMAL(5,2),
    certificate_url VARCHAR(500),
    created_at TIMESTAMP,
    INDEX idx_employee_enrollment (tenant_id, employee_id, status)
);

CREATE TABLE learning_paths (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    job_id UUID NOT NULL REFERENCES jobs(id),
    path_name VARCHAR(255),
    description TEXT,
    target_completion_months INTEGER,
    created_at TIMESTAMP
);

CREATE TABLE learning_path_courses (
    id UUID PRIMARY KEY,
    learning_path_id UUID NOT NULL REFERENCES learning_paths(id),
    course_id UUID NOT NULL REFERENCES courses(id),
    sequence_order INTEGER,
    is_mandatory BOOLEAN
);

CREATE TABLE certifications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    certification_name VARCHAR(255),
    issuing_body VARCHAR(255),
    validity_months INTEGER, -- NULL for lifetime
    created_at TIMESTAMP
);

CREATE TABLE employee_certifications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    certification_id UUID NOT NULL REFERENCES certifications(id),
    certificate_number VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    document_url VARCHAR(500),
    created_at TIMESTAMP,
    INDEX idx_employee_cert (tenant_id, employee_id)
);

CREATE TABLE compliance_trainings (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    training_name VARCHAR(255),
    description TEXT,
    applicable_roles TEXT[], -- JSON array
    renewal_frequency_months INTEGER,
    created_at TIMESTAMP
);

CREATE TABLE compliance_training_records (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL,
    compliance_training_id UUID NOT NULL,
    completion_date DATE,
    expiry_date DATE,
    certificate_url VARCHAR(500),
    is_compliant BOOLEAN,
    created_at TIMESTAMP,
    INDEX idx_compliance_status (tenant_id, employee_id, is_compliant)
);
```

#### API Endpoints

```
GET    /api/v1/courses                             - Course catalog
POST   /api/v1/enrollments                         - Enroll in course
GET    /api/v1/enrollments/{employee_id}           - Learning progress
GET    /api/v1/learning-paths/{job_id}             - Get learning path
GET    /api/v1/compliance-trainings                - Compliance status
GET    /api/v1/analytics/learning-metrics          - Learning insights
```

#### Events

```
✓ CourseCompleted → Certificate generation, Notification
✓ ComplianceTrainingExpiring → Reminder notification
✓ CertificationExpiring → Manager alert
✓ LearningPathStarted → Tracking begins
```

---

### MODULE 8: COMPENSATION & BENEFITS

**Purpose**: Bonus, equity, variable pay, benefits administration

#### Key Features
- Performance-based bonuses
- Stock options and equity vesting
- Benefits enrollment (medical, dental, retirement)
- Benefits utilization tracking
- Flexible benefits/FSA
- Dependent management
- Claims processing

#### Database Schema (Abbreviated)

```sql
CREATE TABLE bonus_plans (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    plan_name VARCHAR(100),
    fiscal_year DATE,
    bonus_type ENUM('Performance', 'Referral', 'Milestone'),
    base_percentage DECIMAL(5,2), -- % of base salary
    performance_metrics JSONB, -- {company_target: 100M, metric: 'revenue'}
    payout_date DATE,
    created_at TIMESTAMP
);

CREATE TABLE employee_bonus_allocations (
    id UUID PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    bonus_plan_id UUID NOT NULL REFERENCES bonus_plans(id),
    allocated_percentage DECIMAL(5,2),
    performance_score DECIMAL(5,2), -- 0-100
    bonus_amount DECIMAL(15,2),
    payout_status ENUM('Pending', 'Approved', 'Paid'),
    paid_date DATE,
    created_at TIMESTAMP
);

CREATE TABLE equity_grants (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    grant_type ENUM('Stock_Options', 'RSU', 'ESPP'),
    shares_granted INTEGER,
    grant_price DECIMAL(15,2),
    vesting_schedule_id UUID,
    grant_date DATE,
    expiry_date DATE,
    created_at TIMESTAMP,
    INDEX idx_employee_equity (tenant_id, employee_id)
);

CREATE TABLE vesting_schedules (
    id UUID PRIMARY KEY,
    schedule_name VARCHAR(100), -- "4-year cliff 1-year"
    total_vesting_years INTEGER,
    cliff_months INTEGER,
    vesting_frequency ENUM('Monthly', 'Quarterly', 'Annually'),
    percentage_per_period DECIMAL(5,2),
    created_at TIMESTAMP
);

CREATE TABLE equity_vesting_transactions (
    id UUID PRIMARY KEY,
    equity_grant_id UUID NOT NULL REFERENCES equity_grants(id),
    vesting_date DATE,
    shares_vested INTEGER,
    vesting_price DECIMAL(15,2),
    vested_value DECIMAL(15,2),
    created_at TIMESTAMP
);

CREATE TABLE benefits_plans (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    plan_name VARCHAR(100), -- "Health Plus", "Dental Basic"
    benefit_type ENUM('Health', 'Dental', 'Vision', 'Life', 'Retirement', 'Wellness'),
    coverage_type ENUM('Individual', 'Family', 'Family+Dependents'),
    employer_contribution DECIMAL(15,2),
    employee_contribution DECIMAL(15,2),
    plan_year_start DATE,
    plan_year_end DATE,
    provider_name VARCHAR(100),
    created_at TIMESTAMP
);

CREATE TABLE employee_benefits (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    benefits_plan_id UUID NOT NULL REFERENCES benefits_plans(id),
    enrollment_date DATE,
    coverage_type ENUM('Individual', 'Family', 'Family+Dependents'),
    dependents_count INTEGER,
    status ENUM('Active', 'Inactive', 'Terminated'),
    created_at TIMESTAMP,
    INDEX idx_employee_benefits (tenant_id, employee_id, status)
);

CREATE TABLE dependents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    dependent_name VARCHAR(100),
    relationship ENUM('Spouse', 'Child', 'Parent'),
    date_of_birth DATE,
    is_covered_under_benefits BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE benefits_claims (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL,
    plan_id UUID NOT NULL REFERENCES benefits_plans(id),
    claim_date DATE,
    claim_amount DECIMAL(15,2),
    claim_status ENUM('Submitted', 'Under_Review', 'Approved', 'Denied'),
    claim_type VARCHAR(100), -- Medical, Dental, etc.
    supporting_documents TEXT[], -- URLs
    approved_amount DECIMAL(15,2),
    approved_date DATE,
    created_at TIMESTAMP
);
```

#### Events

```
✓ BonusCalculated → Finance approval required
✓ EquityGranted → Notification, Vesting begins
✓ EquityVested → Tax implications, ESL tracking
✓ BenefitsEnrollmentOpened → Notification to employees
✓ CoverageEffective → Insurance provider integration
✓ ClaimSubmitted → Insurance claim processing
```

---

### MODULE 9: EMPLOYEE SELF-SERVICE (ESS)

**Purpose**: Employee-facing portal for personal data, applications, documents

#### Key Features
- Profile management (personal, contact, emergency)
- Leave/time-off requests from employee perspective
- Expense claims
- Salary slip download
- Benefit enrollment
- Document requests
- Grievance management
- Career progression view

#### Database Schema (Mostly references other modules)

```sql
CREATE TABLE grievances (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    grievance_type ENUM('Harassment', 'Discrimination', 'Safety', 'Pay', 'Benefits', 'Work_Condition', 'Other'),
    title VARCHAR(255),
    description TEXT,
    attachment_urls TEXT[],
    is_anonymous BOOLEAN,
    status ENUM('Submitted', 'Acknowledged', 'Investigation', 'Resolved', 'Closed'),
    severity ENUM('Low', 'Medium', 'High', 'Critical'),
    assigned_to BIGINT REFERENCES employees(id), -- HR representative
    resolution_date DATE,
    resolution_comments TEXT,
    submitted_at TIMESTAMP,
    INDEX idx_grievance_status (tenant_id, employee_id, status)
);

CREATE TABLE employee_documents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    document_type ENUM('Resume', 'CertificateOfWork', 'ApprovalLetter', 'BankDetails', 'IdentityProof', 'AddressProof', 'TaxDocument'),
    document_url VARCHAR(500),
    upload_date DATE,
    expiry_date DATE,
    verified_by BIGINT REFERENCES employees(id),
    verified_at TIMESTAMP,
    created_at TIMESTAMP
);

CREATE TABLE document_requests (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    request_type VARCHAR(100), -- "Salary Certificate", "Experience Letter"
    request_date DATE,
    delivery_date DATE,
    status ENUM('Pending', 'Generated', 'Delivered'),
    document_url VARCHAR(500),
    created_at TIMESTAMP
);
```

#### REST Endpoints (All employee-scoped)

```
GET    /api/v1/me/profile                          - Get my profile
PATCH  /api/v1/me/profile                          - Update profile
GET    /api/v1/me/balance/leaves                   - My leave balance
GET    /api/v1/me/balance/salary                   - My salary info
GET    /api/v1/me/salary-slips                     - Download salary slips
GET    /api/v1/me/grievances                       - My grievances
POST   /api/v1/me/grievances                       - File grievance
GET    /api/v1/me/benefits                         - My benefits
POST   /api/v1/me/benefit-enrollment               - Enroll benefits
POST   /api/v1/me/document-requests                - Request document
GET    /api/v1/me/reports                          - My reports chain
```

#### Events

```
✓ ProfileUpdated → Sync to payroll, HR
✓ GrievanceSubmitted → HR notification
✓ DocumentRequested → HR task creation
```

---

### MODULE 10: ADMIN & RBAC

**Purpose**: System configuration, security, role management

#### Key Features
- User management and onboarding
- Role and permission definitions
- Custom workflows
- Audit logs
- System configurations
- Tenant management (for SaaS)

#### Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    employee_id BIGINT REFERENCES employees(id), -- NULL for non-employee users
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(500),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    mfa_enabled BOOLEAN DEFAULT false,
    mfa_method ENUM('Totp', 'SMS', 'Email'),
    last_login TIMESTAMP,
    last_password_change TIMESTAMP,
    password_expires_at TIMESTAMP,
    created_at TIMESTAMP,
    INDEX idx_tenant_users (tenant_id, is_active)
);

CREATE TABLE roles (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    role_name VARCHAR(100),
    role_code VARCHAR(50), -- "HR_Manager", "Employee"
    description TEXT,
    is_system_role BOOLEAN, -- Cannot be modified
    created_at TIMESTAMP,
    UNIQUE(tenant_id, role_code)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    permission_code VARCHAR(100), -- "employee.create", "salary.view", "payroll.approve"
    module VARCHAR(50), -- "core-hr", "payroll", "recruitment"
    action VARCHAR(50), -- "create", "read", "update", "delete", "approve"
    description TEXT,
    created_at TIMESTAMP,
    UNIQUE(permission_code)
);

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY,
    role_id UUID NOT NULL REFERENCES roles(id),
    permission_id UUID NOT NULL REFERENCES permissions(id),
    UNIQUE(role_id, permission_id)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    effective_from DATE,
    effective_to DATE,
    assigned_by UUID REFERENCES users(id),
    created_at TIMESTAMP,
    INDEX idx_user_roles (user_id, effective_from)
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID,
    entity_type VARCHAR(50), -- "Employee", "Salary", "Leave"
    entity_id VARCHAR(100),
    action VARCHAR(50), -- "Create", "Update", "Delete", "Approve"
    old_values JSONB,
    new_values JSONB,
    changed_fields TEXT[], -- ["salary", "department"]
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP,
    INDEX idx_audit_entity (tenant_id, entity_type, entity_id),
    INDEX idx_audit_user (tenant_id, user_id, created_at)
);

CREATE TABLE system_configurations (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    config_key VARCHAR(255),
    config_value TEXT,
    value_type ENUM('String', 'Number', 'Boolean', 'JSON'),
    description TEXT,
    is_secret BOOLEAN,
    created_at TIMESTAMP,
    UNIQUE(tenant_id, config_key)
);

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_name VARCHAR(255) NOT NULL,
    tenant_slug VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(255),
    
    -- Subscription
    subscription_tier ENUM('Free', 'Starter', 'Professional', 'Enterprise'),
    subscription_started_at DATE,
    subscription_expires_at DATE,
    max_employees INTEGER,
    
    -- Configuration
    primary_country VARCHAR(10), -- For payroll/compliance
    currency VARCHAR(3),
    timezone VARCHAR(50),
    date_format VARCHAR(10),
    
    -- Organization
    logo_url VARCHAR(500),
    company_name VARCHAR(255),
    company_website VARCHAR(255),
    industry VARCHAR(100),
    employees_count INTEGER,
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    INDEX idx_tenant_slug (tenant_slug)
);

CREATE TABLE data_retention_policies (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    entity_type VARCHAR(50), -- "audit_logs", "documents"
    retention_days INTEGER,
    archive_location VARCHAR(255),
    delete_after_archive BOOLEAN DEFAULT false,
    created_at TIMESTAMP
);
```

#### API Endpoints

```
POST   /api/v1/admin/users                         - Create user
GET    /api/v1/admin/users                         - List users
PATCH  /api/v1/admin/users/{id}                    - Update user
DELETE /api/v1/admin/users/{id}                    - Delete user
POST   /api/v1/admin/roles                         - Create role
GET    /api/v1/admin/roles                         - List roles
POST   /api/v1/admin/roles/{id}/permissions        - Assign permissions
GET    /api/v1/admin/audit-logs                    - View audit logs
GET    /api/v1/admin/configurations                - Get configurations
PATCH  /api/v1/admin/configurations/{key}          - Update configuration
```

#### Events

```
✓ UserCreated → Email verification sent
✓ RoleAssigned → Permissions cached
✓ AuditLogged → Searchable in analytics
✓ PermissionDenied → Logged as security event
```

---

### MODULE 11: WORKFLOW & APPROVAL ENGINE

**Purpose**: Generic workflow orchestration for multi-level approvals

#### Key Features
- Configurable approval chains
- Multi-level approvals
- Conditional routing (approval depends on amount, department, etc.)
- Escalation rules
- SLA tracking
- Workflow history and audit

#### Database Schema

```sql
CREATE TABLE workflow_definitions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    workflow_name VARCHAR(255), -- "Leave Approval", "Requisition Approval"
    workflow_type VARCHAR(50),
    description TEXT,
    is_active BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE workflow_steps (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES workflow_definitions(id),
    step_number INTEGER,
    step_name VARCHAR(100),
    approval_required BOOLEAN,
    approver_role_id UUID REFERENCES roles(id), -- OR approver_user_id for specific user
    approver_user_id UUID REFERENCES users(id),
    approver_hierarchy_level VARCHAR(50), -- "DirectManager", "DepartmentHead"
    approval_logic ENUM('Any', 'All', 'FirstReject'),
    sla_hours INTEGER,
    escalation_on_delay BOOLEAN,
    escalate_to_user_id UUID,
    auto_approve_if_timeout BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    workflow_id UUID NOT NULL REFERENCES workflow_definitions(id),
    request_id VARCHAR(100), -- e.g., leave application ID
    request_type VARCHAR(50),
    current_step_number INTEGER,
    status ENUM('Pending', 'InProgress', 'Approved', 'Rejected', 'Withdrawn'),
    initiated_by UUID,
    initiated_at TIMESTAMP,
    completed_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP,
    INDEX idx_workflow_status (tenant_id, status, created_at)
);

CREATE TABLE workflow_approvals (
    id UUID PRIMARY KEY,
    workflow_instance_id UUID NOT NULL REFERENCES workflow_instances(id),
    step_id UUID NOT NULL REFERENCES workflow_steps(id),
    approver_id UUID NOT NULL,
    action ENUM('Pending', 'Approved', 'Rejected', 'Escalated'),
    comments TEXT,
    action_taken_at TIMESTAMP,
    sla_breached BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE workflow_triggers (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES workflow_definitions(id),
    trigger_event VARCHAR(100), -- "LeaveApplicationSubmitted", "RequisitionCreated"
    condition JSONB, -- {amount: {gt: 100000}, department: {in: [1, 2, 3]}}
    actions_json JSONB, -- [{"type": "StartWorkflow", "workflow_id": "xyz"}]
    created_at TIMESTAMP
);
```

#### REST Endpoints

```
POST   /api/v1/workflows                           - Create workflow
GET    /api/v1/workflows                           - List workflows
POST   /api/v1/workflows/{id}/instances             - Start workflow
GET    /api/v1/workflows/{id}/instances             - Get instances
GET    /api/v1/workflows/instances/{id}             - Get instance details
PATCH  /api/v1/workflows/instances/{id}/approve     - Approve request
PATCH  /api/v1/workflows/instances/{id}/reject      - Reject request
GET    /api/v1/workflows/my-approvals               - My pending approvals
```

#### Events

```
✓ WorkflowInstanceStarted → First approver notified
✓ WorkflowStepApproved → Next step begins
✓ ApprovalSLABreach → Escalation triggered
✓ WorkflowCompleted → Original requester notified
```

---

### MODULE 12: ANALYTICS & REPORTING

**Purpose**: Business intelligence, dashboards, custom reports

#### Key Features
- Pre-built dashboards (HR metrics, payroll, attendance)
- Custom report builder
- Data export (Excel, CSV, PDF)
- Real-time KPIs
- Predictive analytics (attrition, hiring trends)
- Compliance reports
- Data warehouse (columnar storage)

#### Architecture

```
Application DBs (PostgreSQL OLTP)
         ↓
Kafka Topic: analytics_events
         ↓
ETL Pipeline (Apache Spark / dbt)
         ↓
TimescaleDB (Time-series for metrics)
Elasticsearch (Full-text search)
ClickHouse (Columnar OLAP)
         ↓
Analytics Service (Query aggregation)
         ↓
Dashboard / BI Tools (Grafana, Metabase)
```

#### Key Metrics

```json
{
  "HR Metrics": {
    "Total Employees": "COUNT(employees WHERE employment_status='Active')",
    "Headcount Growth": "YoY employee count change",
    "Attrition Rate": "Terminated / Average headcount",
    "Hiring Rate": "New hires / time period",
    "Employee Tenure": "AVG(hire_date)",
    "Diversity Metrics": "Gender, age, ethnicity distribution",
    "Department Distribution": "Employees per department"
  },
  
  "Payroll Metrics": {
    "Total Payroll Cost": "SUM(net_pay) per month",
    "Cost Per Employee": "Payroll / Headcount",
    "Tax Efficiency": "Total tax / Gross income",
    "Bonus Payout Rate": "Total bonuses / Payroll %"
  },
  
  "Attendance Metrics": {
    "Average Attendance": "AVG(attendance %)",
    "Overtime Hours": "SUM(overtime)",
    "Late Arrivals": "COUNT(late check-ins)",
    "Work From Home %": "WFH days / Total days"
  },
  
  "Recruitment Metrics": {
    "Time to Hire": "AVG(offer_accepted - application_date)",
    "Offer Acceptance Rate": "Accepted / Offers sent %",
    "Cost Per Hire": "Total recruitment spend / Hires",
    "Hiring Funnel": "Stage-wise conversion rates"
  },
  
  "Performance Metrics": {
    "High Performer %": "Exceeds_Expectations / Total",
    "Average Rating": "AVG(overall_rating)",
    "Goal Achievement": "AVG(goal completion %)",
    "360 Feedback Score": "AVG(feedback ratings)"
  }
}
```

---

## Database Design Principles

### Multi-Tenancy Pattern

```
STRATEGY: Row-Level Security (RLS)
- Every table has tenant_id
- PostgreSQL RLS policies enforce tenant isolation
- Indexes on (tenant_id, other_columns)

ISOLATION LEVEL: Logical isolation (same database, different tenants)
- Pro: Cost-efficient, operational simplicity
- Con: Less isolation than separate DBs

For very large deployments:
- Tenant-based sharding (shard_key = tenant_id)
- Separate database servers per shard
```

### Partitioning Strategy

```sql
-- Large tables should be partitioned by date/range
CREATE TABLE salary_slips (
    -- ... columns
) PARTITION BY RANGE (DATE_TRUNC('month', salary_month));

CREATE TABLE salary_slips_2026_01 PARTITION OF salary_slips
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
```

### Indexing Strategy

```
Rule 1: Index on tenant_id (first column)
  ├─ Multi-column indexes: (tenant_id, foreign_key)
  ├─ Unique constraint: UNIQUE(tenant_id, business_key)
  └─ Partial indexes for status: WHERE status = 'Active'

Rule 2: Hot query indexes
  ├─ (tenant_id, employee_id, created_at DESC) for recent changes
  ├─ (tenant_id, manager_id) for org hierarchy
  └─ Full-text indexes on TEXT fields

Rule 3: Avoid over-indexing
  ├─ Each index = ~10% storage + maintenance cost
  ├─ Monitor unused indexes
  └─ Analyze query plans
```

---

## API Design Principles

### REST Conventions

```
Resource-Based URLs:
  /api/v1/{module}/{resource}
  /api/v1/{module}/{resource}/{id}
  /api/v1/{module}/{resource}/{id}/{subresource}

Methods:
  POST   /api/v1/employees              - Create
  GET    /api/v1/employees              - List (with pagination, filters)
  GET    /api/v1/employees/{id}         - Get single
  PATCH  /api/v1/employees/{id}         - Update
  DELETE /api/v1/employees/{id}         - Delete (soft delete usually)
  
Actions:
  POST   /api/v1/leaves/{id}/approve    - Custom action
  POST   /api/v1/payroll-runs/{id}/disburse

Query Parameters:
  ?page=1&limit=20                      - Pagination
  ?sort=-created_at&sort=name           - Sorting
  ?filter=status:Active&filter=dept:HR  - Filtering
  ?include=manager,department           - Eager loading relationships
```

### Error Responses

```json
{
  "status": 400,
  "code": "INVALID_REQUEST",
  "message": "Validation error",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ],
  "request_id": "req_12345",
  "timestamp": "2026-04-27T20:00:00Z"
}
```

### Pagination

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1250,
    "pages": 63
  },
  "meta": {
    "request_id": "req_123",
    "response_time_ms": 145
  }
}
```

---

## Event-Driven Architecture

### Event Broker: NATS or Kafka

```
Advantages of NATS (chosen for this design):
✓ Low-latency (<100ms)
✓ Lightweight (Go native)
✓ Built-in request-reply pattern
✓ Streaming (for replay)

Advantages of Kafka:
✓ Event replay from beginning
✓ High throughput
✓ Consumer group management
✓ Exactly-once guarantees
```

### Event Schema

```protobuf
// proto/events/domain_events.proto
message DomainEvent {
    string event_id = 1;           // UUID
    string event_type = 2;          // "EmployeeCreated"
    string aggregate_id = 3;        // employee_id
    string aggregate_type = 4;      // "Employee"
    google.protobuf.Timestamp occurred_at = 5;
    google.protobuf.Any data = 6;   // Event-specific data
    string tenant_id = 7;
    string source_service = 8;      // "core-hr-service"
    int32 version = 9;              // For schema evolution
}

message EmployeeCreatedEvent {
    string employee_id = 1;
    string email = 2;
    string first_name = 3;
    string last_name = 4;
    string department_id = 5;
    string job_id = 6;
    string manager_id = 7;
    google.protobuf.Timestamp hire_date = 8;
    string created_by = 9;
}
```

### Event Consumption Example

```go
// Internal: Service A listening to other services' events
func (svc *PayrollService) handleEmployeeCreated(event *events.EmployeeCreatedEvent) {
    // Create corresponding salary structure record
    // Initialize payroll account
    // Emit: SalaryStructureCreated
}

// Idempotency: Use event_id as deduplication key
func (svc *PayrollService) ProcessEvent(eventID string, data []byte) error {
    // Check: Is this event already processed?
    if exists := cache.Get("event:" + eventID); exists {
        return nil // Already processed
    }
    
    // Process...
    
    // Store processed event ID
    cache.Set("event:" + eventID, true, 30*24*time.Hour)
    return nil
}
```

### Saga Pattern for Distributed Transactions

```
Example: Hire an employee (spans multiple services)

1. Core-HR Service
   └─ Create employee record
   └─ Emit: EmployeeCreated

2. Payroll Service (listening to EmployeeCreated)
   └─ Create salary structure
   └─ Emit: PayrollSetupComplete

3. ESS Service (listening to EmployeeCreated)
   └─ Create ESS account
   └─ Emit: ESSAccountCreated

4. Notification Service (listening to all)
   └─ Send welcome email
   └─ Send IT setup email

All parallel, eventual consistency, retry logic on failures
```

---

## Microfrontend Strategy

### Shell Architecture

```
┌─────────────────────────────────────────┐
│  HOST APP (Shell) - Main Navigation      │
│  ├─ Top bar (logo, user menu)            │
│  ├─ Left sidebar (module navigation)     │
│  └─ Main content area (outlet)           │
└─────────────────────────────────────────┘
       ↓
┌─────────────────────────────────────────┐
│  Module Microfrontends (Independent)     │
│  ├─ Core HR MFE (employee mgmt)          │
│  ├─ Recruitment MFE (ATS)                │
│  ├─ Payroll MFE (salary mgmt)            │
│  ├─ Attendance MFE (check-in/out)        │
│  ├─ Leave MFE (leave apps)               │
│  ├─ Performance MFE (reviews)            │
│  ├─ Learning MFE (courses)               │
│  ├─ Compensation MFE (bonus, equity)     │
│  ├─ ESS MFE (employee self-service)      │
│  ├─ Admin MFE (users, config)            │
│  ├─ Analytics MFE (dashboards)           │
│  └─ Workflow MFE (approvals)             │
└─────────────────────────────────────────┘
```

### Module Federation (Next.js)

```typescript
// Host App (apps/host-app/next.config.ts)
import { NextFederationPlugin } from '@module-federation/nextjs-mf';

export default {
  webpack: (config, options) => {
    config.plugins.push(
      new NextFederationPlugin({
        name: 'host-app',
        filename: 'static/chunks/remoteEntry.js',
        remotes: {
          'core-hr': 'core_hr@http://localhost:3011/remoteEntry.js',
          'recruitment': 'recruitment@http://localhost:3012/remoteEntry.js',
          'payroll': 'payroll@http://localhost:3013/remoteEntry.js',
          'attendance': 'attendance@http://localhost:3014/remoteEntry.js',
          'leave': 'leave@http://localhost:3015/remoteEntry.js',
          'performance': 'performance@http://localhost:3016/remoteEntry.js',
          // ... others
        },
        exposes: {
          './Layout': './components/Layout',
          './useAuth': './hooks/useAuth',
        },
        shared: {
          react: { singleton: true },
          'react-dom': { singleton: true },
          '@erp/logger': { singleton: true },
        },
      }),
    );
    return config;
  },
};
```

### State Management Across MFEs

```typescript
// Shared stores (Redis cache + local state)
// Global auth context via provider
// Event-based communication for cross-module state updates

// Example: Leave MFE -> Attendance MFE
// When leave is approved, Attendance MFE marks employee as "On Leave"
// Via event: LeaveApplicationApproved → update cache → UI refresh
```

---

## Security & RBAC

### Authentication Flow

```
1. User logs in (email + password)
2. Auth Service validates
3. Generate JWT token (1 hour)
   └─ Payload: {user_id, email, tenant_id, roles: [role_ids], permissions: [perm_codes]}
4. Client stores JWT in httpOnly cookie
5. API Gateway validates JWT on every request
6. Service decodes JWT, checks permissions

REFRESH TOKEN:
- Store in secure storage
- 30-day expiry
- Rotated on each use
- Used to get new access token
```

### Permission Model

```
Role-Based Access Control (RBAC):
  User
    ├─ Admin
    │   └─ All permissions
    ├─ HR Manager
    │   ├─ employee.create, employee.read, employee.update
    │   ├─ leave.approve
    │   └─ recruitment.manage
    ├─ Manager
    │   ├─ employee.read (direct reports only)
    │   ├─ leave.approve (direct reports)
    │   ├─ attendance.view (team)
    │   └─ performance.create (self + team)
    └─ Employee
        ├─ employee.read (self)
        ├─ leave.apply
        ├─ attendance.checkin, attendance.checkout
        ├─ performance.read (self)
        └─ salary_slip.read (self)

Attribute-Based Access Control (ABAC):
  - Extend RBAC for complex rules
  - "HR Manager can approve leaves for salary < $100k"
  - "Finance can view payroll after approval"
  - Rules in workflow engine
```

### Data Encryption

```
At Rest:
  ├─ PostgreSQL: Column-level encryption for PII
  │  └─ ssn, pan, passport_number, bank_account
  └─ Disk: Full-disk encryption (EBS encryption in AWS)

In Transit:
  ├─ TLS 1.3 for all APIs
  ├─ Service-to-service: mTLS (certificate-based)
  └─ Database connections: SSL

Key Management:
  ├─ AWS KMS for key encryption
  ├─ Rotate keys annually
  └─ Separate keys per tenant (optional, for very sensitive)
```

### Audit & Compliance

```
AUDIT TRAIL:
- Every entity create/update/delete logged to audit_logs table
- Include: user_id, old_values, new_values, timestamp, ip_address
- Immutable: Never update audit logs
- Retention: 7+ years for payroll/tax records

COMPLIANCE:
  ├─ GDPR: Right to be forgotten, data portability
  ├─ CCPA: Opt-out, data transparency
  ├─ SOC 2: Access controls, encryption, monitoring
  ├─ Local Payroll: Tax compliance, statutory deductions
  └─ Industry: Healthcare (HIPAA), Finance (PCI)

Reporting:
  ├─ Who changed what and when
  ├─ Approval chains for sensitive operations
  ├─ Export audit logs for auditors
  └─ Real-time alerts on suspicious activities
```

---

## Deployment Strategy

### Docker & Kubernetes

```dockerfile
# Dockerfile (Go service)
FROM golang:1.26-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
```

```yaml
# kubernetes/core-hr-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: core-hr-service
  labels:
    app: core-hr-service
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: core-hr-service
  template:
    metadata:
      labels:
        app: core-hr-service
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - core-hr-service
              topologyKey: kubernetes.io/hostname
      
      containers:
      - name: core-hr-service
        image: myregistry.azurecr.io/core-hr-service:1.2.3
        imagePullPolicy: IfNotPresent
        
        ports:
        - name: http
          containerPort: 8080
        - name: grpc
          containerPort: 9090
        
        env:
        - name: PORT
          value: "8080"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
        
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        
        resources:
          requests:
            cpu: 250m
            memory: 512Mi
          limits:
            cpu: 500m
            memory: 1Gi

---
apiVersion: v1
kind: Service
metadata:
  name: core-hr-service
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    name: http
  - port: 9090
    targetPort: 9090
    name: grpc
  selector:
    app: core-hr-service

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: core-hr-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: core-hr-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Database Deployment

```yaml
# kubernetes/postgresql.yaml (using Helm / Operator)
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: hrms-postgres
spec:
  instances: 3
  primaryUpdateStrategy: unsupervised
  postgresql:
    parameters:
      max_connections: "1000"
      shared_buffers: "256MB"
      effective_cache_size: "1GB"
      work_mem: "4MB"
      maintenance_work_mem: "64MB"
      random_page_cost: "1.1"
  bootstrap:
    initdb:
      database: hrms
      owner: hrms_user
      secret:
        name: db-secret
  storage:
    size: 500Gi
    storageClass: fast-ssd
  monitoring:
    enabled: true
  backup:
    retentionPolicy: "30d"
    barmanObjectStore:
      destinationPath: s3://backup-bucket/hrms/
      s3Credentials:
        accessKeyId:
          name: aws-credentials
          key: accesskey
        secretAccessKey:
          name: aws-credentials
          key: secretkey
```

### Observability Stack

```yaml
# kubernetes/monitoring-stack.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: 'core-hr-service'
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
          - production
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: core-hr-service

---
# Loki (Log aggregation)
apiVersion: v1
kind: ConfigMap
metadata:
  name: loki-config
data:
  loki-config.yaml: |
    ingester:
      chunk_idle_period: 3m
      max_chunk_age: 1h
    storage_config:
      boltdb_shipper:
        active_index_directory: /loki/index
        cache_location: /loki/boltdb-cache
    limits_config:
      enforce_metric_name: false

---
# Jaeger (Distributed tracing)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jaeger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jaeger
  template:
    metadata:
      labels:
        app: jaeger
    spec:
      containers:
      - name: jaeger
        image: jaegertracing/all-in-one:latest
        ports:
        - containerPort: 16686  # UI
        - containerPort: 14250  # gRPC receiver
```

---

## Folder Structure (TurboRepo Monorepo)

```
erp-microservices/
├── package.json (root, workspace definition)
├── turbo.json (build pipeline)
├── pnpm-workspace.yaml
├── docker-compose.yml (local dev)
├── .github/workflows/ (CI/CD)
│
├── apis/                          # Backend services
│   ├── go.work (Go workspace)
│   ├── shared/
│   │   ├── auth/
│   │   │   ├── go.mod
│   │   │   ├── jwt.go
│   │   │   └── rbac.go
│   │   ├── logger/
│   │   │   ├── go.mod
│   │   │   └── zap.go
│   │   ├── middleware/
│   │   │   ├── cors.go
│   │   │   ├── tenant.go
│   │   │   └── logging.go
│   │   └── database/
│   │       ├── go.mod
│   │       ├── postgres.go
│   │       └── migrations.go
│   │
│   ├── api-gateway/
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── config/
│   │   ├── routes/
│   │   └── middleware/
│   │
│   └── services/
│       ├── core-hr/
│       │   ├── go.mod
│       │   ├── cmd/api/main.go
│       │   ├── internal/
│       │   │   ├── domain/
│       │   │   │   └── employee.go
│       │   │   ├── repository/
│       │   │   │   └── employee_repo.go
│       │   │   ├── service/
│       │   │   │   └── employee_service.go
│       │   │   ├── handler/
│       │   │   │   └── employee_handler.go
│       │   │   └── events/
│       │   │       └── publisher.go
│       │   ├── proto/
│       │   │   └── core_hr/v1/
│       │   │       ├── employees.proto
│       │   │       └── departments.proto
│       │   └── migrations/
│       │       ├── 001_create_employees_table.sql
│       │       └── 002_create_departments_table.sql
│       │
│       ├── recruitment/
│       ├── attendance/
│       ├── leave/
│       ├── payroll/
│       ├── performance/
│       ├── learning/
│       ├── compensation/
│       ├── ess/
│       ├── admin/
│       ├── workflow/
│       └── analytics/
│
├── frontend/                      # Frontend monorepo
│   ├── package.json
│   ├── pnpm-lock.yaml
│   ├── turbo.json
│   ├── tsconfig.json
│   │
│   ├── apps/
│   │   ├── host-app/              # Shell application
│   │   │   ├── package.json
│   │   │   ├── next.config.ts
│   │   │   ├── tsconfig.json
│   │   │   ├── app/
│   │   │   │   ├── layout.tsx
│   │   │   │   └── page.tsx
│   │   │   ├── components/
│   │   │   │   ├── Layout.tsx
│   │   │   │   ├── Navbar.tsx
│   │   │   │   └── Sidebar.tsx
│   │   │   ├── hooks/
│   │   │   ├── lib/
│   │   │   └── public/
│   │   │
│   │   ├── core-hr-mfe/            # Core HR Microfrontend
│   │   │   ├── package.json
│   │   │   ├── next.config.ts
│   │   │   ├── app/
│   │   │   │   ├── employees/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── [id]/page.tsx
│   │   │   │   └── org-chart/
│   │   │   └── components/
│   │   │
│   │   ├── recruitment-mfe/
│   │   ├── payroll-mfe/
│   │   ├── attendance-mfe/
│   │   ├── leave-mfe/
│   │   ├── performance-mfe/
│   │   ├── learning-mfe/
│   │   ├── compensation-mfe/
│   │   ├── ess-mfe/
│   │   ├── admin-mfe/
│   │   ├── analytics-mfe/
│   │   └── workflow-mfe/
│   │
│   └── packages/
│       ├── shared-ui/              # Shared UI components
│       │   ├── package.json
│       │   ├── src/
│       │   │   ├── Button.tsx
│       │   │   ├── Modal.tsx
│       │   │   ├── Table.tsx
│       │   │   ├── Form.tsx
│       │   │   └── index.ts
│       │   └── tsconfig.json
│       │
│       ├── shared-hooks/            # Shared React hooks
│       │   ├── package.json
│       │   ├── src/
│       │   │   ├── useAuth.ts
│       │   │   ├── useTenant.ts
│       │   │   ├── useApi.ts
│       │   │   ├── usePagination.ts
│       │   │   └── index.ts
│       │   └── tsconfig.json
│       │
│       ├── shared-types/            # Shared TypeScript types
│       │   ├── package.json
│       │   ├── src/
│       │   │   ├── employee.ts
│       │   │   ├── leave.ts
│       │   │   ├── payroll.ts
│       │   │   ├── common.ts
│       │   │   └── index.ts
│       │   └── tsconfig.json
│       │
│       ├── logger/                  # @erp/logger
│       │   ├── package.json
│       │   └── src/
│       │       └── index.ts
│       │
│       └── api-client/              # API client library
│           ├── package.json
│           ├── src/
│           │   ├── client.ts
│           │   ├── interceptors.ts
│           │   └── endpoints.ts
│           └── tsconfig.json
│
├── infrastructure/                # IaC and deployment configs
│   ├── kubernetes/
│   │   ├── namespaces.yaml
│   │   ├── configmaps.yaml
│   │   ├── secrets.yaml
│   │   ├── database/
│   │   │   └── postgresql-operator.yaml
│   │   ├── services/
│   │   │   ├── core-hr-service.yaml
│   │   │   ├── recruitment-service.yaml
│   │   │   └── ... (all services)
│   │   ├── ingress.yaml
│   │   ├── istio-config.yaml
│   │   └── monitoring/
│   │       ├── prometheus.yaml
│   │       ├── loki.yaml
│   │       ├── grafana.yaml
│   │       └── jaeger.yaml
│   │
│   ├── terraform/
│   │   ├── main.tf
│   │   ├── vpc.tf
│   │   ├── eks.tf (AWS EKS)
│   │   ├── rds.tf (PostgreSQL)
│   │   ├── elasticache.tf (Redis)
│   │   └── variables.tf
│   │
│   └── docker-compose.yml
│
├── proto/                         # Protocol Buffers
│   ├── common/
│   │   └── v1/
│   │       └── common.proto
│   ├── core_hr/
│   │   └── v1/
│   │       ├── employees.proto
│   │       ├── departments.proto
│   │       └── org_hierarchy.proto
│   ├── recruitment/
│   ├── payroll/
│   └── ... (all services)
│
├── migrations/                    # Database migrations
│   ├── postgres/
│   │   ├── 001_init_schema.sql
│   │   ├── 002_core_hr_tables.sql
│   │   ├── 003_recruitment_tables.sql
│   │   ├── 004_attendance_tables.sql
│   │   ├── 005_leave_tables.sql
│   │   ├── 006_payroll_tables.sql
│   │   ├── 007_performance_tables.sql
│   │   ├── 008_learning_tables.sql
│   │   ├── 009_compensation_tables.sql
│   │   ├── 010_ess_tables.sql
│   │   ├── 011_admin_tables.sql
│   │   ├── 012_workflow_tables.sql
│   │   ├── 013_analytics_tables.sql
│   │   ├── 014_audit_tables.sql
│   │   └── 015_add_indexes.sql
│   │
│   └── scripts/
│       ├── run-migrations.sh
│       └── rollback.sh
│
├── docs/
│   ├── README.md
│   ├── ARCHITECTURE.md
│   ├── API.md
│   ├── DEVELOPMENT.md
│   ├── DEPLOYMENT.md
│   └── TROUBLESHOOTING.md
│
└── .github/
    └── workflows/
        ├── ci.yml (tests, linting)
        ├── cd-api.yml (deploy APIs)
        └── cd-frontend.yml (deploy frontend)
```

---

## Development Phases

### PHASE 1: MVP (3-4 months)
**Scope**: Core HR + ESS
**Modules**: Core HR, ESS, Admin, Auth
**Capacity**: 5 engineers

```
Core Features:
✓ Employee master data management
✓ Org hierarchy
✓ Basic RBAC
✓ Employee self-service (view profile, salary slip)
✓ Audit logging
✓ Multi-tenancy foundation

Output:
  - Core HR Service + DB
  - ESS Frontend
  - Admin Panel (user management)
  - Auth Service
  - Basic API Gateway
```

### PHASE 2: Talent Lifecycle (2-3 months)
**Scope**: Recruitment + Attendance + Leave
**Modules**: Recruitment, Attendance, Leave, Workflow Engine
**Capacity**: 8 engineers

```
New Features:
✓ Job requisitions
✓ Candidate tracking
✓ Offer management
✓ Check-in/check-out
✓ Timesheet approval
✓ Leave applications & approvals
✓ Approval workflows

Infrastructure:
  - Event streaming (NATS/Kafka)
  - Workflow engine
  - Notification service
```

### PHASE 3: Compensation (2-3 months)
**Scope**: Payroll + Performance
**Modules**: Payroll, Performance, Compensation
**Capacity**: 10 engineers

```
Complex Features:
✓ Salary structure & components
✓ Tax calculations
✓ Payroll processing
✓ Performance reviews
✓ Bonuses & equity
✓ Benefits administration

Integrations:
  - Bank integrations for disbursement
  - External payroll API (if needed)
  - Tax compliance (country-specific)
```

### PHASE 4: Analytics & Advanced (2 months)
**Scope**: Learning, Analytics, Workflow Advanced
**Modules**: Learning, Analytics, Advanced Workflow
**Capacity**: 6 engineers

```
Final Features:
✓ Course management
✓ BI dashboards
✓ Predictive analytics (attrition, hiring trends)
✓ Advanced workflow rules
✓ Compliance reporting
✓ Data warehouse

Infrastructure:
  - TimescaleDB / ClickHouse
  - ETL pipeline
  - Analytics service
  - Dashboarding tool integration
```

### PHASE 5: Enterprise Hardening (Ongoing)
**Scope**: Security, Performance, Scale
**Capacity**: 4-6 engineers

```
Non-functional:
✓ Load testing
✓ Security audits (GDPR, SOC 2)
✓ Performance optimization
✓ High availability setup
✓ Disaster recovery
✓ Advanced monitoring
```

---

## Common Pitfalls in HR Systems

### 1. **Underestimating Payroll Complexity**
❌ **Pitfall**: "Payroll is just salary / 12"
✅ **Reality**: 
- Earnings with cascading deductions (HRA, DA, etc.)
- Statutory contributions (PF, ESI, IT)
- Complex tax regimes (progressive slabs, surcharge, cess)
- Country-specific calculations (France min 30 days, India gratuity)
- Variable components, bonuses, overtime
- Tax implications of different components

**Mitigation**: Build tax engine early. Hire domain expert (payroll consultant).

### 2. **Ignoring Data Quality from Day 1**
❌ **Pitfall**: "We'll clean data later"
✅ **Reality**: Bad data in production ruins everything

**Mitigation**:
- Strict validation on employee creation
- Duplicate detection (same SSN, email)
- PII data governance (who can see what)
- Regular data quality audits

### 3. **Tight Coupling Between Services**
❌ **Pitfall**: Payroll Service directly queries Core-HR DB
✅ **Better**: Event-driven → Payroll consumes EmployeeCreated events

**Mitigation**: Define contracts (protocols/events), enforce via code review.

### 4. **Not Designing for Multi-Tenancy from Start**
❌ **Pitfall**: "We'll add tenants later"
✅ **Cost**: Massive refactoring (row-level security, data isolation, billing)

**Mitigation**: Include tenant_id in all queries. Test with multiple tenants in development.

### 5. **Underestimating Compliance Complexity**
❌ **Pitfall**: "We'll handle GDPR/compliance later"
✅ **Reality**: 
- Right to be forgotten (hard with audit logs)
- Data portability (export all employee data)
- Consent management (marketing emails)
- Tax compliance varies by country

**Mitigation**: Involve legal/compliance team early. Document privacy policies.

### 6. **Poor Workflow Engine Design**
❌ **Pitfall**: Hardcoding approval chains (if manager_id → approve → send email)
✅ **Better**: Configurable workflow engine with SLA tracking, escalations

**Mitigation**: Build workflow engine as standalone service. Support complex rules.

### 7. **Ignoring Performance Under Load**
❌ **Pitfall**: Works fine with 100 employees, crashes at 10,000
✅ **Reality**: 
- Payroll run: Process 5,000 employees' salaries in minutes
- Reporting: Aggregate analytics across millions of records
- Search: Employee directory with faceted filters

**Mitigation**: Load test from Month 3. Use TimescaleDB for time-series. Index aggressively.

### 8. **Weak Audit Trail**
❌ **Pitfall**: "User updated salary → can't track what changed"
✅ **Better**: Store old_values & new_values, IP address, timestamp

**Mitigation**: Immutable audit logs. Retention: 7+ years.

### 9. **Not Handling Leave Edge Cases**
❌ **Pitfall**: "Employee applied for leave, leave approved, but employee left company"
✅ **Reality**: 
- Carry-over policies
- Country-specific minimums (France 30 days/year)
- Gender-based leave (maternity, paternity)
- Probation period restrictions
- Pro-rata leave for mid-year joins

**Mitigation**: Build comprehensive leave policy engine. Test extensively.

### 10. **Ignoring Scalability of Notifications**
❌ **Pitfall**: Synchronous email sending blocking API
✅ **Better**: Async queue → Notification service → Email provider

**Mitigation**: Use message queues (NATS, RabbitMQ). Retry logic. Fallback channels.

### 11. **Weak Search/Discovery**
❌ **Pitfall**: No easy way to find employees by skills, location, projects
✅ **Better**: Elasticsearch for full-text search + faceted navigation

**Mitigation**: Build search index early. Keep it in sync with primary DB.

### 12. **Not Planning for Legacy Data**
❌ **Pitfall**: "We'll migrate from legacy HRMS manually"
✅ **Reality**: Millions of historical records, inconsistent data formats

**Mitigation**: Build data migration tool. Plan for 2-3 months of parallel running.

---

## Key Metrics & KPIs

### System Health
- API Response Time: p50 < 200ms, p99 < 1000ms
- Service Availability: 99.95% (< 22 minutes downtime/month)
- Database Query Performance: p95 < 500ms
- Error Rate: < 0.1%

### Business Metrics
- Payroll Processing Time: Complete 5,000 salaries in < 30 minutes
- Leave Application TTL: Approve within 2 business days (SLA)
- Recruitment Time to Hire: < 45 days average
- Data Accuracy: > 99.9%

### User Adoption
- Monthly Active Users: > 80% of assigned users
- Feature Adoption: > 70% using core features
- Support Tickets: < 5 per 1,000 employees
- Training Completion: > 90%

---

## Conclusion

This HRMS architecture is:
- **Modular**: 11 independent services
- **Scalable**: Handles millions of users
- **Compliant**: GDPR, local payroll laws, SOC 2
- **Extensible**: New modules can be added without disrupting existing ones
- **Future-Proof**: Event-driven foundation allows for AI/ML additions

Next: Detailed implementation guides for each module, with Go code examples and Next.js component examples.
