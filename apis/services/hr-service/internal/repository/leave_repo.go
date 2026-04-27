package repository

import (
	"time"

	"erp/hr-service/internal/models"

	"gorm.io/gorm"
)

// LeaveRepository handles leave data access
type LeaveRepository struct {
	*BaseRepository
}

// NewLeaveRepository creates new leave repository
func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{
		BaseRepository: &BaseRepository{DB: db},
	}
}

// CreateLeave creates new leave application
func (r *LeaveRepository) CreateLeave(leave *models.Leave) error {
	return r.DB.Create(leave).Error
}

// GetLeaveByID retrieves leave by ID
func (r *LeaveRepository) GetLeaveByID(tenantID string, leaveID uint) (*models.Leave, error) {
	var leave models.Leave

	err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, leaveID).
		Preload("Employee").
		Preload("LeaveType").
		Preload("ApprovalChain").
		First(&leave).Error

	if err != nil {
		return nil, err
	}
	return &leave, nil
}

// GetEmployeeLeaves retrieves all leaves for an employee
func (r *LeaveRepository) GetEmployeeLeaves(tenantID string, empID uint, year int) ([]models.Leave, error) {
	var leaves []models.Leave

	err := r.DB.Where("tenant_id = ? AND employee_id = ? AND EXTRACT(YEAR FROM start_date) = ?",
		tenantID, empID, year).
		Preload("LeaveType").
		Preload("ApprovalChain").
		Order("start_date DESC").
		Find(&leaves).Error

	if err != nil {
		return nil, err
	}
	return leaves, nil
}

// GetPendingLeaveApprovals retrieves pending approvals for a manager
func (r *LeaveRepository) GetPendingLeaveApprovals(tenantID string, managerID uint) ([]models.Leave, error) {
	var leaves []models.Leave

	err := r.DB.Where("tenant_id = ? AND status = ?", tenantID, "Submitted").
		Joins("JOIN leave_approvals ON leave_approvals.leave_id = leaves.id").
		Where("leave_approvals.approver_id = ? AND leave_approvals.status = ?", managerID, "Pending").
		Preload("Employee").
		Preload("LeaveType").
		Preload("ApprovalChain").
		Find(&leaves).Error

	if err != nil {
		return nil, err
	}
	return leaves, nil
}

// UpdateLeaveStatus updates leave status
func (r *LeaveRepository) UpdateLeaveStatus(tenantID string, leaveID uint, status string) error {
	return r.DB.Model(&models.Leave{}).
		Where("tenant_id = ? AND id = ?", tenantID, leaveID).
		Update("status", status).Error
}

// ApproveLeave approves leave and creates approval record
func (r *LeaveRepository) ApproveLeave(tenantID string, leaveID uint, approverID uint, comments string) error {
	var leave models.Leave

	// Get leave
	if err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, leaveID).First(&leave).Error; err != nil {
		return err
	}

	// Create approval
	approval := models.Approval{
		RequestID:    leaveID,
		ApproverID:   approverID,
		Status:       "Approved",
		ApprovalDate: &time.Time{},
		Comments:     comments,
	}

	*approval.ApprovalDate = time.Now()

	if err := r.DB.Create(&approval).Error; err != nil {
		return err
	}

	// Update leave status
	return r.DB.Model(&models.Leave{}).
		Where("tenant_id = ? AND id = ?", tenantID, leaveID).
		Update("status", "Approved").Error
}

// RejectLeave rejects leave
func (r *LeaveRepository) RejectLeave(tenantID string, leaveID uint, approverID uint, reason string) error {
	var leave models.Leave

	if err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, leaveID).First(&leave).Error; err != nil {
		return err
	}

	return r.DB.Model(&models.Leave{}).
		Where("tenant_id = ? AND id = ?", tenantID, leaveID).
		Updates(map[string]interface{}{
			"status":           "Rejected",
			"rejection_reason": reason,
		}).Error
}

// GetLeaveBalance retrieves leave balance for employee
func (r *LeaveRepository) GetLeaveBalance(tenantID string, empID uint, leaveTypeID uint) (*models.LeaveAllocation, error) {
	var allocation models.LeaveAllocation
	currentYear := time.Now().Format("2006")

	err := r.DB.Where("tenant_id = ? AND employee_id = ? AND leave_type_id = ? AND financial_year = ?",
		tenantID, empID, leaveTypeID, currentYear+"-"+string(rune(time.Now().Year()+1))).
		First(&allocation).Error

	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

// GetOverlappingLeaves checks for overlapping leave applications
func (r *LeaveRepository) GetOverlappingLeaves(tenantID string, empID uint, startDate, endDate time.Time) ([]models.Leave, error) {
	var leaves []models.Leave

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND status IN (?, ?) AND ((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?))",
		tenantID, empID, "Approved", "Submitted",
		endDate, startDate, endDate, startDate).
		Find(&leaves).Error

	if err != nil {
		return nil, err
	}
	return leaves, nil
}

// GetLeaveTypes retrieves all available leave types
func (r *LeaveRepository) GetLeaveTypes(tenantID string) ([]models.LeaveType, error) {
	var leaveTypes []models.LeaveType

	err := r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Find(&leaveTypes).Error

	if err != nil {
		return nil, err
	}
	return leaveTypes, nil
}

// GetHolidays retrieves holidays for current year
func (r *LeaveRepository) GetHolidays(tenantID string) ([]models.Holiday, error) {
	var holidays []models.Holiday
	year := time.Now().Year()

	err := r.DB.Where("tenant_id = ?", tenantID).
		Joins("JOIN holiday_calendars ON holiday_calendars.id = holidays.holiday_calendar_id").
		Where("holiday_calendars.year = ? AND holiday_calendars.is_active = ?", year, true).
		Find(&holidays).Error

	if err != nil {
		return nil, err
	}
	return holidays, nil
}

// GetLeaveStats retrieves leave statistics for dashboard
func (r *LeaveRepository) GetLeaveStats(tenantID string, empID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalApplied int64
	var totalApproved int64
	var totalRejected int64
	var totalPending int64

	r.DB.Where("tenant_id = ? AND employee_id = ?", tenantID, empID).
		Model(&models.Leave{}).
		Count(&totalApplied)

	r.DB.Where("tenant_id = ? AND employee_id = ? AND status = ?", tenantID, empID, "Approved").
		Model(&models.Leave{}).
		Count(&totalApproved)

	r.DB.Where("tenant_id = ? AND employee_id = ? AND status = ?", tenantID, empID, "Rejected").
		Model(&models.Leave{}).
		Count(&totalRejected)

	r.DB.Where("tenant_id = ? AND employee_id = ? AND status IN (?, ?)", tenantID, empID, "Draft", "Submitted").
		Model(&models.Leave{}).
		Count(&totalPending)

	stats["total_applied"] = totalApplied
	stats["total_approved"] = totalApproved
	stats["total_rejected"] = totalRejected
	stats["total_pending"] = totalPending

	return stats, nil
}
