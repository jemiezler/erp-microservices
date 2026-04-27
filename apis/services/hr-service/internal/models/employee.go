package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	EmployeeID  string    `json:"employee_id" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	Position    string    `json:"position"`
	Department  string    `json:"department"`
	Status      string    `json:"status" gorm:"default:active"`
	ManagerID   *uint     `json:"manager_id"`
	Manager     *Employee `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Role        string    `json:"role" gorm:"default:employee"`
}
