package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ApprovalWorkflow defines approval processes
type ApprovalWorkflow struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name       string `json:"name"`                     // Leave Approval, Salary Approval, Expense Approval
	EntityType string `json:"entity_type" gorm:"index"` // Leave, Expense, AssetRequest, etc

	Description string `json:"description"`

	Levels []ApprovalLevel `json:"levels,omitempty"`

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// ApprovalLevel individual level in approval workflow
type ApprovalLevel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	WorkflowID uint              `json:"workflow_id" gorm:"index"`
	Workflow   *ApprovalWorkflow `json:"workflow,omitempty"`

	Level int    `json:"level"` // 1, 2, 3...
	Name  string `json:"name"`  // Manager Approval, HR Approval

	ApprovesFrom string `json:"approves_from"` // Manager, Role, Static
	// If ApprovesFrom = Manager: next manager in hierarchy
	// If ApprovesFrom = Role: users with specific role
	// If ApprovesFrom = Static: specific users defined in ApproversID

	ApproversID       datatypes.JSONSlice[int] `json:"approvers_id" gorm:"type:jsonb"` // IDs if static
	RequiredApprovals int                      `json:"required_approvals"`             // If multiple approvers, how many needed

	CanDelegate bool `json:"can_delegate"`

	TimeoutDays int `json:"timeout_days"` // Auto-escalate after days

	gorm.Model
}

// ApprovalRequest individual approval instance
type ApprovalRequest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	RequestNumber string `json:"request_number" gorm:"uniqueIndex"`

	WorkflowID uint              `json:"workflow_id" gorm:"index"`
	Workflow   *ApprovalWorkflow `json:"workflow,omitempty"`

	EntityType string `json:"entity_type"` // Leave, Expense, etc
	EntityID   uint   `json:"entity_id" gorm:"index"`

	RequestedByID uint      `json:"requested_by_id" gorm:"index"`
	RequestedBy   *Employee `json:"requested_by,omitempty" gorm:"foreignKey:RequestedByID"`

	RequestedDate time.Time `json:"requested_date"`

	CurrentLevel int    `json:"current_level"`
	Status       string `json:"status"` // Pending, Approved, Rejected, Cancelled

	RejectionReason string `json:"rejection_reason"`

	Approvals []Approval `json:"approvals,omitempty"`

	gorm.Model
}

// Approval individual approval decision
type Approval struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	RequestID uint             `json:"request_id" gorm:"index"`
	Request   *ApprovalRequest `json:"request,omitempty"`

	Level int `json:"level"`

	ApprovingLevel uint `json:"approving_level_id"`

	ApproverID uint      `json:"approver_id" gorm:"index"`
	Approver   *Employee `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`

	DelegatedFromID *uint     `json:"delegated_from_id"` // If delegated
	DelegatedFrom   *Employee `json:"delegated_from,omitempty" gorm:"foreignKey:DelegatedFromID"`

	Status string `json:"status"` // Pending, Approved, Rejected, Delegated

	ApprovalDate *time.Time `json:"approval_date"`
	Comments     string     `json:"comments" gorm:"type:text"`

	SequentialApproval bool `json:"sequential_approval"` // Must wait for previous level

	gorm.Model
}

// RBACRole role-based access control roles
type RBACRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name        string `json:"name" gorm:"uniqueIndex:idx_role_name_tenant"` // Admin, Manager, HR, Employee
	Description string `json:"description"`

	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`

	IsSystemRole bool `json:"is_system_role"` // Cannot be deleted
	IsActive     bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// Permission granular permissions
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Resource    string `json:"resource"`                // employees, leaves, payroll
	Action      string `json:"action"`                  // view, create, edit, delete, approve
	Code        string `json:"code" gorm:"uniqueIndex"` // employees.view, leaves.approve
	Description string `json:"description"`

	Roles []RBACRole `json:"roles,omitempty" gorm:"many2many:role_permissions;"`

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// AuditLog tracks all system changes
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	UserID uint      `json:"user_id" gorm:"index"`
	User   *Employee `json:"user,omitempty" gorm:"foreignKey:UserID"`

	EntityType string `json:"entity_type" gorm:"index"` // Employee, Leave, Payroll
	EntityID   uint   `json:"entity_id" gorm:"index"`

	Action string `json:"action"` // Create, Update, Delete, View, Export

	OldValues datatypes.JSON `json:"old_values" gorm:"type:jsonb"`
	NewValues datatypes.JSON `json:"new_values" gorm:"type:jsonb"`

	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	Status       string `json:"status"` // Success, Failed
	ErrorMessage string `json:"error_message"`

	gorm.Model
}

// SystemSetting global system configuration
type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	SettingKey   string `json:"setting_key" gorm:"uniqueIndex:idx_setting_key_tenant"`
	SettingValue string `json:"setting_value" gorm:"type:jsonb"`
	DataType     string `json:"data_type"` // string, integer, boolean, json

	Description string `json:"description"`
	Module      string `json:"module"` // HR, Payroll, Leave, etc

	gorm.Model
}
