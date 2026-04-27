package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// LeaveType defines types of leaves available (Annual, Sick, Casual, etc.)
type LeaveType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name        string `json:"name"` // Annual, Sick, Casual, Maternity, Paternity
	Code        string `json:"code" gorm:"uniqueIndex"`
	Description string `json:"description"`

	DaysPerYear     float64 `json:"days_per_year"`
	IsCarryForward  bool    `json:"is_carry_forward"` // Can unused days be carried to next year
	MaxCarryForward float64 `json:"max_carry_forward"`

	ApprovalRequired bool `json:"approval_required"`
	ApprovalLevel    int  `json:"approval_level"` // 1 = Manager, 2 = HR Manager, 3 = Director

	IsGenderSpecific bool   `json:"is_gender_specific"`
	ApplicableGender string `json:"applicable_gender"` // If gender specific: M, F, All

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// LeavePolicy defines leave policies for departments/roles
type LeavePolicy struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name        string `json:"name"`
	Description string `json:"description"`

	DepartmentID uint        `json:"department_id" gorm:"index"`
	Department   *Department `json:"department,omitempty"`

	EmploymentType string `json:"employment_type"` // Full-time, Part-time

	LeaveAllocations []LeaveAllocation `json:"leave_allocations,omitempty"`

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// LeaveAllocation allocates leave days to employees
type LeaveAllocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	LeaveTypeID uint       `json:"leave_type_id" gorm:"index"`
	LeaveType   *LeaveType `json:"leave_type,omitempty"`

	FinancialYear  string    `json:"financial_year"` // 2024-2025
	AllocationDate time.Time `json:"allocation_date"`

	AllocatedDays  float64 `json:"allocated_days"`
	UsedDays       float64 `json:"used_days"`       // Calculated from Leave records
	PendingDays    float64 `json:"pending_days"`    // In approval
	AvailableDays  float64 `json:"available_days"`  // Calculated: allocated - used
	CarriedForward float64 `json:"carried_forward"` // From previous year

	Status string `json:"status"` // Active, Expired

	Notes string `json:"notes"`

	gorm.Model
}

// Leave tracks leave applications
type Leave struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	LeaveNumber string    `json:"leave_number" gorm:"uniqueIndex"` // AUTO-2024-001
	EmployeeID  uint      `json:"employee_id" gorm:"index"`
	Employee    *Employee `json:"employee,omitempty"`

	LeaveTypeID uint       `json:"leave_type_id" gorm:"index"`
	LeaveType   *LeaveType `json:"leave_type,omitempty"`

	StartDate time.Time `json:"start_date" gorm:"index"`
	EndDate   time.Time `json:"end_date"`

	// Half-day support
	IsFullDay   bool   `json:"is_full_day" gorm:"default:true"`
	HalfDayType string `json:"half_day_type"` // First, Second (if half-day)

	NumberOfDays float64 `json:"number_of_days"` // Calculated including weekends
	WorkingDays  float64 `json:"working_days"`   // Excluding weekends

	Reason             string `json:"reason" gorm:"type:text"`
	ContactDuringLeave string `json:"contact_during_leave"`

	Status        string     `json:"status"` // Draft, Submitted, Approved, Rejected, Cancelled
	SubmittedDate *time.Time `json:"submitted_date"`

	ApprovalChain []LeaveApproval `json:"approval_chain,omitempty"`

	RejectionReason string `json:"rejection_reason"`

	IsBackdated        bool `json:"is_backdated"`
	IsCarryForwardUsed bool `json:"is_carry_forward_used"`

	gorm.Model
}

// LeaveApproval tracks multi-level leave approvals
type LeaveApproval struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LeaveID uint   `json:"leave_id" gorm:"index"`
	Leave   *Leave `json:"leave,omitempty"`

	ApprovalLevel int       `json:"approval_level"` // 1, 2, 3
	ApproverID    uint      `json:"approver_id" gorm:"index"`
	Approver      *Employee `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`

	Status       string     `json:"status"` // Pending, Approved, Rejected
	ApprovalDate *time.Time `json:"approval_date"`

	Comments string `json:"comments"`

	gorm.Model
}

// HolidayCalendar defines national and company holidays
type HolidayCalendar struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name    string `json:"name"` // India Holidays 2024
	Year    int    `json:"year"`
	Country string `json:"country"`

	Holidays []Holiday `json:"holidays,omitempty"`

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// Holiday individual holiday entry
type Holiday struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	HolidayCalendarID uint             `json:"holiday_calendar_id" gorm:"index"`
	HolidayCalendar   *HolidayCalendar `json:"holiday_calendar,omitempty"`

	Name        string    `json:"name"` // New Year Day
	HolidayDate time.Time `json:"holiday_date" gorm:"index"`
	IsOptional  bool      `json:"is_optional"`
	Description string    `json:"description"`

	ApplicableStates datatypes.JSONSlice[string] `json:"applicable_states" gorm:"type:jsonb"` // For state-specific holidays

	gorm.Model
}

// LeaveEncashment tracks leave encashment (conversion to pay)
type LeaveEncashment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	LeaveTypeID uint       `json:"leave_type_id"`
	LeaveType   *LeaveType `json:"leave_type,omitempty"`

	FinancialYear string `json:"financial_year"`

	DaysEncashed float64 `json:"days_encashed"`
	RatePerDay   float64 `json:"rate_per_day"`
	TotalAmount  float64 `json:"total_amount"` // Calculated

	RequestDate   time.Time  `json:"request_date"`
	ProcessedDate *time.Time `json:"processed_date"`

	Status  string `json:"status"` // Pending, Approved, Processed
	Remarks string `json:"remarks"`

	gorm.Model
}
