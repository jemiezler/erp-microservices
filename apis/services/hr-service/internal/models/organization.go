package models

import (
	"time"
)

// Organization represents a company/tenant in multi-tenant system
type Organization struct {
	ID                  uint         `gorm:"primaryKey" json:"id"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	Name                string       `json:"name" gorm:"uniqueIndex:idx_org_name_tenant"`
	Code                string       `json:"code" gorm:"uniqueIndex:idx_org_code"`
	TenantID            string       `json:"tenant_id" gorm:"index"`
	Description         string       `json:"description"`
	Website             string       `json:"website"`
	Industry            string       `json:"industry"`
	HeadquartersAddress string       `json:"headquarters_address"`
	EmployeeCount       int          `json:"employee_count"`
	IsActive            bool         `json:"is_active" gorm:"default:true"`
	Departments         []Department `json:"departments,omitempty" gorm:"foreignKey:OrganizationID"`
}

// Department represents organizational departments
type Department struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	TenantID       string      `json:"tenant_id" gorm:"index"`
	OrganizationID uint        `json:"organization_id" gorm:"index"`
	Name           string      `json:"name"`
	Code           string      `json:"code"`
	Description    string      `json:"description"`
	HeadID         *uint       `json:"head_id"`
	Head           *Employee   `json:"head,omitempty" gorm:"foreignKey:HeadID"`
	ParentID       *uint       `json:"parent_id" gorm:"index"`
	Parent         *Department `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	CostCenter     string      `json:"cost_center"`
	Location       string      `json:"location"`
	Budget         float64     `json:"budget"`
	IsActive       bool        `json:"is_active" gorm:"default:true"`
	Employees      []Employee  `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
}

// JobTitle represents job positions in organization
type JobTitle struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TenantID    string     `json:"tenant_id" gorm:"index"`
	Title       string     `json:"title" gorm:"index"`
	Description string     `json:"description"`
	Department  string     `json:"department"`
	MinSalary   float64    `json:"min_salary"`
	MaxSalary   float64    `json:"max_salary"`
	ReportingTo string     `json:"reporting_to"`
	Level       string     `json:"level"` // Entry, Mid, Senior, Lead, Manager
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	Employees   []Employee `json:"employees,omitempty" gorm:"foreignKey:JobTitleID"`
}

// Location represents office locations
type Location struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	TenantID       string     `json:"tenant_id" gorm:"index"`
	Name           string     `json:"name"`
	Code           string     `json:"code" gorm:"uniqueIndex:idx_loc_code_tenant"`
	Address        string     `json:"address"`
	City           string     `json:"city"`
	State          string     `json:"state"`
	Country        string     `json:"country"`
	PostalCode     string     `json:"postal_code"`
	Latitude       float64    `json:"latitude"`
	Longitude      float64    `json:"longitude"`
	Phone          string     `json:"phone"`
	IsHeadquarters bool       `json:"is_headquarters"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	Employees      []Employee `json:"employees,omitempty" gorm:"foreignKey:LocationID"`
}

// EmployeeHierarchy tracks org structure changes for reporting
type EmployeeHierarchy struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	TenantID   string    `json:"tenant_id" gorm:"index"`
	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID"`
	ManagerID  *uint     `json:"manager_id"`
	Manager    *Employee `gorm:"foreignKey:ManagerID"`
	Level      int       `json:"level"` // Depth in hierarchy
	Path       string    `json:"path"`  // ltree path: 1.2.3
}
