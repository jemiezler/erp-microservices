package repository

import (
	"time"

	"erp/hr-service/internal/models"

	"gorm.io/gorm"
)

// AttendanceRepository handles attendance data access
type AttendanceRepository struct {
	*BaseRepository
}

// NewAttendanceRepository creates new attendance repository
func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{
		BaseRepository: &BaseRepository{DB: db},
	}
}

// RecordCheckIn records employee check-in
func (r *AttendanceRepository) RecordCheckIn(tenantID string, empID uint, latitude, longitude float64, deviceID string) (*models.Attendance, error) {
	now := time.Now()
	attendance := &models.Attendance{
		TenantID:         tenantID,
		EmployeeID:       empID,
		AttendanceDate:   now,
		CheckInTime:      &now,
		CheckInLatitude:  latitude,
		CheckInLongitude: longitude,
		CheckInDeviceID:  deviceID,
		Status:           "Present",
	}

	if err := r.DB.Create(attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

// RecordCheckOut records employee check-out and calculates working hours
func (r *AttendanceRepository) RecordCheckOut(tenantID string, empID uint, latitude, longitude float64, deviceID string) (*models.Attendance, error) {
	now := time.Now()

	// Get today's attendance
	var attendance models.Attendance
	if err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND DATE(attendance_date) = DATE(?)",
		tenantID, empID, now).
		First(&attendance).Error; err != nil {
		return nil, err
	}

	// Update checkout
	attendance.CheckOutTime = &now
	attendance.CheckOutLatitude = latitude
	attendance.CheckOutLongitude = longitude
	attendance.CheckOutDeviceID = deviceID

	// Calculate working hours
	if attendance.CheckInTime != nil && attendance.CheckOutTime != nil {
		duration := attendance.CheckOutTime.Sub(*attendance.CheckInTime)
		attendance.WorkingHours = duration.Hours()
	}

	if err := r.DB.Save(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

// GetTodayAttendance gets today's attendance for employee
func (r *AttendanceRepository) GetTodayAttendance(tenantID string, empID uint) (*models.Attendance, error) {
	var attendance models.Attendance
	today := time.Now().Format("2006-01-02")

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND DATE(attendance_date) = ?",
		tenantID, empID, today).
		First(&attendance).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &attendance, nil
}

// GetAttendanceHistory retrieves attendance history for date range
func (r *AttendanceRepository) GetAttendanceHistory(tenantID string, empID uint, startDate, endDate time.Time) ([]models.Attendance, error) {
	var records []models.Attendance

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND attendance_date BETWEEN ? AND ?",
		tenantID, empID, startDate, endDate).
		Order("attendance_date DESC").
		Find(&records).Error

	if err != nil {
		return nil, err
	}
	return records, nil
}

// GetMonthlyAttendanceStats gets monthly attendance statistics
func (r *AttendanceRepository) GetMonthlyAttendanceStats(tenantID string, empID uint, year, month int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalDays int64
	var presentDays int64
	var absentDays int64
	var lateDays int64
	var halfDays int64
	var wfhDays int64

	query := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND EXTRACT(YEAR FROM attendance_date) = ? AND EXTRACT(MONTH FROM attendance_date) = ?",
		tenantID, empID, year, month)

	query.Model(&models.Attendance{}).Count(&totalDays)
	query.Where("status = ?", "Present").Model(&models.Attendance{}).Count(&presentDays)
	query.Where("status = ?", "Absent").Model(&models.Attendance{}).Count(&absentDays)
	query.Where("status = ?", "Late").Model(&models.Attendance{}).Count(&lateDays)
	query.Where("status = ?", "HalfDay").Model(&models.Attendance{}).Count(&halfDays)
	query.Where("status = ?", "WFH").Model(&models.Attendance{}).Count(&wfhDays)

	var totalWorkingHours float64
	query.Model(&models.Attendance{}).Select("COALESCE(SUM(working_hours), 0)").Row().Scan(&totalWorkingHours)

	stats["total_days"] = totalDays
	stats["present_days"] = presentDays
	stats["absent_days"] = absentDays
	stats["late_days"] = lateDays
	stats["half_days"] = halfDays
	stats["wfh_days"] = wfhDays
	stats["total_working_hours"] = totalWorkingHours

	return stats, nil
}

// GetDepartmentAttendance gets attendance summary for department
func (r *AttendanceRepository) GetDepartmentAttendance(tenantID string, deptID uint, date time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalEmp int64
	var present int64
	var absent int64
	var leave int64

	dateStr := date.Format("2006-01-02")

	// Get total employees in department
	r.DB.Where("tenant_id = ? AND department_id = ? AND is_active = ?", tenantID, deptID, true).
		Model(&models.Employee{}).Count(&totalEmp)

	// Get attendance summary
	r.DB.Joins("JOIN employees ON employees.id = attendances.employee_id").
		Where("attendances.tenant_id = ? AND employees.department_id = ? AND DATE(attendance_date) = ? AND status = ?",
			tenantID, deptID, dateStr, "Present").
		Model(&models.Attendance{}).Count(&present)

	r.DB.Joins("JOIN employees ON employees.id = attendances.employee_id").
		Where("attendances.tenant_id = ? AND employees.department_id = ? AND DATE(attendance_date) = ? AND status = ?",
			tenantID, deptID, dateStr, "Absent").
		Model(&models.Attendance{}).Count(&absent)

	r.DB.Joins("JOIN employees ON employees.id = attendances.employee_id").
		Where("attendances.tenant_id = ? AND employees.department_id = ? AND DATE(attendance_date) = ? AND status = ?",
			tenantID, deptID, dateStr, "OnLeave").
		Model(&models.Attendance{}).Count(&leave)

	stats["total_employees"] = totalEmp
	stats["present"] = present
	stats["absent"] = absent
	stats["on_leave"] = leave
	if totalEmp > 0 {
		stats["attendance_percentage"] = float64(present) / float64(totalEmp) * 100
	}

	return stats, nil
}

// UpdateAttendanceStatus updates attendance status manually (for corrections)
func (r *AttendanceRepository) UpdateAttendanceStatus(tenantID string, attendanceID uint, status string, notes string) error {
	return r.DB.Model(&models.Attendance{}).
		Where("tenant_id = ? AND id = ?", tenantID, attendanceID).
		Updates(map[string]interface{}{
			"status": status,
			"notes":  notes,
		}).Error
}

// GetWFHRequests gets work from home requests
func (r *AttendanceRepository) GetWFHRequests(tenantID string, empID uint, status string) ([]models.WorkFromHome, error) {
	var wfhRequests []models.WorkFromHome

	query := r.DB.Where("tenant_id = ? AND employee_id = ?", tenantID, empID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("date DESC").Find(&wfhRequests).Error
	if err != nil {
		return nil, err
	}
	return wfhRequests, nil
}

// ApproveWFH approves work from home request
func (r *AttendanceRepository) ApproveWFH(tenantID string, wfhID uint, approverID uint) error {
	now := time.Now()
	return r.DB.Model(&models.WorkFromHome{}).
		Where("tenant_id = ? AND id = ?", tenantID, wfhID).
		Updates(map[string]interface{}{
			"status":         "Approved",
			"approved_by_id": approverID,
			"approved_date":  now,
		}).Error
}
