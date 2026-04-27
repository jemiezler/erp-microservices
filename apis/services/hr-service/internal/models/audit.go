package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	EmployeeID uint   `json:"employee_id"`
	Action     string `json:"action"` // e.g., "PROFILE_UPDATE", "SALARY_CHANGE"
	Field      string `json:"field"`
	OldValue   string `json:"old_value"`
	NewValue   string `json:"new_value"`
	ChangedBy  uint   `json:"changed_by"`
}
