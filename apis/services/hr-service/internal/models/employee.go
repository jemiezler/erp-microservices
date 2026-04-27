package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	EmployeeID string    `json:"employee_id" gorm:"uniqueIndex"`
	Name       string    `json:"name"`
	Email      string    `json:"email" gorm:"uniqueIndex"`
	Position   string    `json:"position"`
	Department string    `json:"department"`
	Status     string    `json:"status" gorm:"default:active"`
	ManagerID  *uint     `json:"manager_id"`
	Manager    *Employee `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Role       string    `json:"role" gorm:"default:employee"`
}

// Dependent represents employee's dependent family members
type Dependent struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EmployeeID   uint      `json:"employee_id" gorm:"index"`
	Employee     *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Gender       string    `json:"gender"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
}

// EmergencyContact represents employee's emergency contact
type EmergencyContact struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EmployeeID   uint      `json:"employee_id" gorm:"index"`
	Employee     *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
}

// Qualification represents employee's educational qualifications
type Qualification struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	EmployeeID        uint      `json:"employee_id" gorm:"index"`
	Employee          *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	QualificationType string    `json:"qualification_type"` // Bachelor, Master, PhD, Diploma, Certificate
	InstituteName     string    `json:"institute_name"`
	Specialization    string    `json:"specialization"`
	PassYear          int       `json:"pass_year"`
	Grade             string    `json:"grade"`
	DocumentURL       string    `json:"document_url"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
}

// WorkExperience represents employee's previous work experience
type WorkExperience struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	EmployeeID       uint      `json:"employee_id" gorm:"index"`
	Employee         *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	CompanyName      string    `json:"company_name"`
	Designation      string    `json:"designation"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	EmploymentType   string    `json:"employment_type"` // Full-time, Part-time, Contract
	CurrentlyWorking bool      `json:"currently_working"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
}
