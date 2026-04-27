package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Attendance tracks daily attendance records
type Attendance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	AttendanceDate time.Time `json:"attendance_date" gorm:"index"`

	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`

	WorkingHours float64 `json:"working_hours"` // Calculated

	Status    string  `json:"status"`     // Present, Absent, Late, HalfDay, OnLeave, WFH
	LeaveType *string `json:"leave_type"` // If status is OnLeave
	LeaveID   *uint   `json:"leave_id"`

	Location          string  `json:"location"` // Office/Remote
	CheckInLatitude   float64 `json:"check_in_latitude"`
	CheckInLongitude  float64 `json:"check_in_longitude"`
	CheckOutLatitude  float64 `json:"check_out_latitude"`
	CheckOutLongitude float64 `json:"check_out_longitude"`

	CheckInDeviceID  string `json:"check_in_device_id"`
	CheckOutDeviceID string `json:"check_out_device_id"`

	Notes        string     `json:"notes"`
	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedDate *time.Time `json:"approved_date"`

	gorm.Model
}

// TimeSheet tracks work hours and project allocation
type TimeSheet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	WeekStartDate time.Time `json:"week_start_date" gorm:"index"`
	WeekEndDate   time.Time `json:"week_end_date"`

	TotalHours       float64 `json:"total_hours"`
	BillableHours    float64 `json:"billable_hours"`
	NonBillableHours float64 `json:"non_billable_hours"`

	Status        string     `json:"status"` // Draft, Submitted, Approved, Rejected
	SubmittedDate *time.Time `json:"submitted_date"`
	ApprovedByID  *uint      `json:"approved_by_id"`
	ApprovedDate  *time.Time `json:"approved_date"`

	RejectionReason string `json:"rejection_reason"`

	Entries []TimeSheetEntry `json:"entries,omitempty"`

	gorm.Model
}

// TimeSheetEntry daily entry in timesheet
type TimeSheetEntry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TimeSheetID uint       `json:"time_sheet_id" gorm:"index"`
	TimeSheet   *TimeSheet `json:"time_sheet,omitempty"`

	EntryDate time.Time `json:"entry_date"`

	ProjectID       string `json:"project_id"` // Optional project allocation
	ProjectName     string `json:"project_name"`
	TaskDescription string `json:"task_description"`

	Hours      float64 `json:"hours"`
	IsBillable bool    `json:"is_billable"`

	Notes string `json:"notes"`

	gorm.Model
}

// ShiftManagement defines work shifts
type Shift struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name string `json:"name"` // Morning, Afternoon, Night
	Code string `json:"code" gorm:"uniqueIndex"`

	StartTime     string `json:"start_time"` // HH:MM format
	EndTime       string `json:"end_time"`
	BreakDuration int    `json:"break_duration"` // Minutes

	WeekDays datatypes.JSONSlice[int] `json:"week_days" gorm:"type:jsonb"` // [1,2,3,4,5] = Mon-Fri

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// ShiftAssignment assigns shifts to employees
type ShiftAssignment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	ShiftID uint   `json:"shift_id" gorm:"index"`
	Shift   *Shift `json:"shift,omitempty"`

	EffectiveDate time.Time  `json:"effective_date"`
	EndDate       *time.Time `json:"end_date"` // NULL = Ongoing

	Status string `json:"status"` // Active, Inactive

	gorm.Model
}

// BiometricData for fingerprint/iris authentication
type BiometricData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	BiometricType string `json:"biometric_type"` // Fingerprint, Iris, Face, RFID
	BiometricID   string `json:"biometric_id" gorm:"index"`

	DeviceID string `json:"device_id"`
	IsActive bool   `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// WorkFromHome tracks WFH approvals
type WorkFromHome struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	Date time.Time `json:"date" gorm:"index"`

	Reason string `json:"reason"`
	Status string `json:"status"` // Pending, Approved, Rejected

	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedDate *time.Time `json:"approved_date"`

	RejectionReason string `json:"rejection_reason"`

	ContactNumber string `json:"contact_number"`

	gorm.Model
}
