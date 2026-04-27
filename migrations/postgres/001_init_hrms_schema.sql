-- HRMS Database Schema - Phase 1 (Core HR + ESS)
-- PostgreSQL 14+
-- Tenant-isolated using Row-Level Security

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================================================
-- MULTI-TENANCY & CORE INFRASTRUCTURE
-- =============================================================================

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
    primary_country VARCHAR(10), -- 'IN', 'US', 'GB', etc.
    currency VARCHAR(3),
    timezone VARCHAR(50),
    date_format VARCHAR(10),
    
    -- Organization
    logo_url VARCHAR(500),
    company_name VARCHAR(255),
    company_website VARCHAR(255),
    industry VARCHAR(100),
    employees_count INTEGER DEFAULT 0,
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_tenant_slug UNIQUE(tenant_slug),
    INDEX idx_tenant_active (is_active),
    INDEX idx_tenant_created (created_at)
);

-- =============================================================================
-- USERS & AUTHENTICATION
-- =============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id BIGINT, -- NULL for non-employee users (e.g., auditors)
    
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(500),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    
    is_active BOOLEAN DEFAULT true,
    mfa_enabled BOOLEAN DEFAULT false,
    mfa_method VARCHAR(50), -- 'totp', 'sms', 'email'
    
    last_login TIMESTAMP,
    last_password_change TIMESTAMP,
    password_expires_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_user_per_tenant UNIQUE(tenant_id, email),
    INDEX idx_tenant_users (tenant_id, is_active),
    INDEX idx_user_email (email)
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    role_name VARCHAR(100) NOT NULL,
    role_code VARCHAR(50) NOT NULL,
    description TEXT,
    
    is_system_role BOOLEAN DEFAULT false, -- Cannot be deleted/modified
    is_active BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_role_per_tenant UNIQUE(tenant_id, role_code),
    INDEX idx_tenant_roles (tenant_id, is_active)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    permission_code VARCHAR(100) NOT NULL UNIQUE,
    module VARCHAR(50), -- 'core_hr', 'recruitment', 'payroll'
    action VARCHAR(50), -- 'create', 'read', 'update', 'delete', 'approve'
    description TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    
    CONSTRAINT unique_role_permission UNIQUE(role_id, permission_id),
    INDEX idx_role_permissions (role_id)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    
    effective_from DATE,
    effective_to DATE,
    assigned_by UUID REFERENCES users(id),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_roles (user_id, effective_from),
    INDEX idx_active_roles (user_id, effective_to)
);

-- =============================================================================
-- CORE HR: EMPLOYEES & ORGANIZATION
-- =============================================================================

CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    parent_department_id UUID REFERENCES departments(id),
    
    department_head_id BIGINT, -- Will reference employees after creation
    description TEXT,
    cost_center_id VARCHAR(50),
    
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_dept_per_tenant UNIQUE(tenant_id, code),
    INDEX idx_tenant_dept (tenant_id, active),
    INDEX idx_dept_parent (parent_department_id)
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    job_title VARCHAR(255) NOT NULL,
    job_code VARCHAR(50),
    job_family VARCHAR(100), -- e.g., "Engineering", "Finance"
    job_level VARCHAR(50), -- e.g., "Junior", "Senior", "Lead"
    
    description TEXT,
    min_salary DECIMAL(15,2),
    max_salary DECIMAL(15,2),
    
    reports_to_job_id UUID REFERENCES jobs(id),
    
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_job_per_tenant UNIQUE(tenant_id, job_code),
    INDEX idx_tenant_jobs (tenant_id, active)
);

CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Employment identification
    employee_id VARCHAR(50) NOT NULL, -- Human-readable: EMP001, EMP002, etc.
    
    -- Personal information
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    personal_email VARCHAR(255),
    
    -- Employment details
    employment_status ENUM('Active', 'OnLeave', 'Terminated', 'OnNotice') DEFAULT 'Active',
    employment_type ENUM('FullTime', 'PartTime', 'Contract', 'Intern') DEFAULT 'FullTime',
    hire_date DATE NOT NULL,
    termination_date DATE,
    termination_reason VARCHAR(500),
    
    -- Organization
    department_id UUID NOT NULL REFERENCES departments(id),
    job_id UUID NOT NULL REFERENCES jobs(id),
    manager_id BIGINT REFERENCES employees(id),
    cost_center_id UUID,
    
    -- Personal
    date_of_birth DATE,
    gender ENUM('Male', 'Female', 'Other', 'PreferNotToSay') DEFAULT 'PreferNotToSay',
    nationality VARCHAR(100),
    personal_id_number VARCHAR(50), -- SSN, PAN, Aadhaar (ENCRYPTED in production)
    
    -- Contact
    phone_number VARCHAR(20),
    alternate_phone VARCHAR(20),
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    
    -- Emergency contact
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relationship VARCHAR(50),
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    CONSTRAINT unique_employee_per_tenant UNIQUE(tenant_id, employee_id),
    CONSTRAINT unique_email_per_tenant UNIQUE(tenant_id, email),
    CONSTRAINT no_self_manager CHECK (manager_id != id),
    
    INDEX idx_tenant_employee (tenant_id, employment_status),
    INDEX idx_employee_manager (manager_id),
    INDEX idx_employee_department (department_id),
    INDEX idx_employee_job (job_id),
    INDEX idx_employee_email (email)
);

-- Organizational hierarchy with path-based queries
CREATE TABLE org_hierarchy (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    manager_id BIGINT REFERENCES employees(id) ON DELETE SET NULL,
    
    level INTEGER, -- 0=CEO, 1=C-level, 2=VP, etc.
    hierarchy_path TEXT, -- e.g., "1/12/456" for ancestor queries
    
    effective_from DATE,
    effective_to DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT no_self_manager CHECK (employee_id != manager_id),
    INDEX idx_hierarchy (tenant_id, manager_id),
    INDEX idx_hierarchy_path (hierarchy_path)
);

-- Employment history
CREATE TABLE employment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    
    previous_employee_id VARCHAR(50),
    previous_department_id UUID REFERENCES departments(id),
    previous_job_id UUID REFERENCES jobs(id),
    previous_manager_id BIGINT REFERENCES employees(id),
    previous_salary DECIMAL(15,2),
    
    change_date DATE,
    change_reason VARCHAR(200),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_employee_history (tenant_id, employee_id, change_date DESC)
);

-- Employment contracts
CREATE TABLE employment_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    
    contract_type ENUM('Permanent', 'FixedTerm', 'Probation', 'Apprentice'),
    start_date DATE NOT NULL,
    end_date DATE,
    notice_period_days INTEGER,
    
    terms_conditions TEXT,
    document_url VARCHAR(500),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_employee_contract (tenant_id, employee_id)
);

-- =============================================================================
-- AUDIT & COMPLIANCE
-- =============================================================================

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    user_id UUID,
    entity_type VARCHAR(50), -- 'Employee', 'Department', 'Salary'
    entity_id VARCHAR(100),
    
    action VARCHAR(50), -- 'Create', 'Update', 'Delete', 'Approve'
    
    old_values JSONB,
    new_values JSONB,
    changed_fields TEXT[], -- Array of field names
    
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_audit_entity (tenant_id, entity_type, entity_id),
    INDEX idx_audit_user (tenant_id, user_id, created_at),
    INDEX idx_audit_time (tenant_id, created_at DESC)
);

-- System configurations
CREATE TABLE system_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    config_key VARCHAR(255) NOT NULL,
    config_value TEXT,
    value_type ENUM('String', 'Number', 'Boolean', 'JSON'),
    
    description TEXT,
    is_secret BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_config_per_tenant UNIQUE(tenant_id, config_key),
    INDEX idx_tenant_config (tenant_id)
);

-- =============================================================================
-- ROW-LEVEL SECURITY (RLS) POLICIES
-- =============================================================================

-- Enable RLS on all tables
ALTER TABLE employees ENABLE ROW LEVEL SECURITY;
ALTER TABLE departments ENABLE ROW LEVEL SECURITY;
ALTER TABLE jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE employment_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE employment_contracts ENABLE ROW LEVEL SECURITY;
ALTER TABLE org_hierarchy ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Example RLS policy: Employees can see their own data + their reports
CREATE POLICY employees_tenant_isolation ON employees
    FOR SELECT
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY employees_tenant_isolation_modify ON employees
    FOR UPDATE
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Similar policies for other tables...

-- =============================================================================
-- INDEXES FOR PERFORMANCE
-- =============================================================================

CREATE INDEX idx_employees_last_name ON employees(tenant_id, last_name);
CREATE INDEX idx_employees_department ON employees(department_id, employment_status);
CREATE INDEX idx_org_hierarchy_ancestry ON org_hierarchy USING GIST(hierarchy_path);

-- =============================================================================
-- VIEWS FOR COMMON QUERIES
-- =============================================================================

-- View: Active employees with full hierarchy info
CREATE VIEW active_employees_view AS
SELECT 
    e.id,
    e.tenant_id,
    e.employee_id,
    e.first_name,
    e.last_name,
    e.email,
    d.name as department_name,
    j.job_title,
    m.first_name as manager_first_name,
    m.last_name as manager_last_name,
    m.email as manager_email,
    e.hire_date,
    e.employment_status
FROM employees e
LEFT JOIN departments d ON e.department_id = d.id
LEFT JOIN jobs j ON e.job_id = j.id
LEFT JOIN employees m ON e.manager_id = m.id
WHERE e.employment_status = 'Active';

-- View: Org hierarchy tree
CREATE VIEW org_hierarchy_view AS
SELECT 
    e.id,
    e.employee_id,
    e.first_name || ' ' || e.last_name as name,
    j.job_title,
    oh.level,
    oh.hierarchy_path,
    e.manager_id
FROM employees e
LEFT JOIN jobs j ON e.job_id = j.id
LEFT JOIN org_hierarchy oh ON e.id = oh.employee_id
ORDER BY oh.level, oh.hierarchy_path;

-- =============================================================================
-- INITIAL DATA (System Roles & Permissions)
-- =============================================================================

-- Roles
INSERT INTO roles (id, role_name, role_code, description, is_system_role) 
VALUES 
    (gen_random_uuid(), 'System Administrator', 'admin', 'Full system access', true),
    (gen_random_uuid(), 'HR Manager', 'hr_manager', 'HR team member', true),
    (gen_random_uuid(), 'Department Manager', 'manager', 'Line manager', true),
    (gen_random_uuid(), 'Employee', 'employee', 'Regular employee', true),
    (gen_random_uuid(), 'Finance Manager', 'finance_manager', 'Finance team member', true),
    (gen_random_uuid(), 'Auditor', 'auditor', 'Compliance auditor', true);

-- Permissions (sample)
INSERT INTO permissions (id, permission_code, module, action, description) VALUES
    (gen_random_uuid(), 'employee.create', 'core_hr', 'create', 'Create new employee'),
    (gen_random_uuid(), 'employee.read', 'core_hr', 'read', 'View employee details'),
    (gen_random_uuid(), 'employee.update', 'core_hr', 'update', 'Edit employee information'),
    (gen_random_uuid(), 'employee.delete', 'core_hr', 'delete', 'Terminate employee'),
    (gen_random_uuid(), 'org.view', 'core_hr', 'read', 'View organization structure'),
    (gen_random_uuid(), 'audit.view', 'admin', 'read', 'View audit logs');
