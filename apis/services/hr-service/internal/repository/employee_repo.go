package repository

import (
	"erp/hr-service/internal/models"

	"gorm.io/gorm"
)

// BaseRepository provides common CRUD operations
type BaseRepository struct {
	DB *gorm.DB
}

// EmployeeRepository handles employee data access
type EmployeeRepository struct {
	*BaseRepository
}

// NewEmployeeRepository creates new employee repository
func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{
		BaseRepository: &BaseRepository{DB: db},
	}
}

// GetAllEmployees retrieves all active employees with filtering
func (r *EmployeeRepository) GetAllEmployees(tenantID string, filters map[string]interface{}) ([]models.Employee, error) {
	var employees []models.Employee
	query := r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	// Apply filters
	if dept, ok := filters["department_id"]; ok {
		query = query.Where("department_id = ?", dept)
	}
	if status, ok := filters["employment_status"]; ok {
		query = query.Where("employment_status = ?", status)
	}
	if empType, ok := filters["employment_type"]; ok {
		query = query.Where("employment_type = ?", empType)
	}

	// Load relations
	query = query.Preload("JobTitle").
		Preload("Department").
		Preload("Location").
		Preload("Manager").
		Preload("Dependents").
		Preload("EmergencyContacts").
		Preload("Qualifications").
		Preload("WorkExperience")

	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

// GetEmployeeByID retrieves employee by ID
func (r *EmployeeRepository) GetEmployeeByID(tenantID string, empID uint) (*models.Employee, error) {
	var employee models.Employee

	err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, empID).
		Preload("JobTitle").
		Preload("Department").
		Preload("Location").
		Preload("Manager").
		Preload("Dependents").
		Preload("EmergencyContacts").
		Preload("Qualifications").
		Preload("WorkExperience").
		First(&employee).Error

	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetEmployeeByEmployeeID retrieves by employee ID (e.g., "EMP-001")
func (r *EmployeeRepository) GetEmployeeByEmployeeID(tenantID string, employeeID string) (*models.Employee, error) {
	var employee models.Employee

	err := r.DB.Where("tenant_id = ? AND employee_id = ?", tenantID, employeeID).
		Preload("JobTitle").
		Preload("Department").
		Preload("Location").
		Preload("Manager").
		First(&employee).Error

	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// CreateEmployee creates new employee record
func (r *EmployeeRepository) CreateEmployee(employee *models.Employee) error {
	return r.DB.Create(employee).Error
}

// UpdateEmployee updates employee record
func (r *EmployeeRepository) UpdateEmployee(tenantID string, empID uint, updates map[string]interface{}) error {
	return r.DB.Model(&models.Employee{}).
		Where("tenant_id = ? AND id = ?", tenantID, empID).
		Updates(updates).Error
}

// GetEmployeesByDepartment retrieves employees in specific department
func (r *EmployeeRepository) GetEmployeesByDepartment(tenantID string, deptID uint) ([]models.Employee, error) {
	var employees []models.Employee

	err := r.DB.Where("tenant_id = ? AND department_id = ? AND is_active = ?", tenantID, deptID, true).
		Preload("JobTitle").
		Preload("Manager").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}
	return employees, nil
}

// GetEmployeesByManager retrieves all direct reports of a manager
func (r *EmployeeRepository) GetEmployeesByManager(tenantID string, managerID uint) ([]models.Employee, error) {
	var employees []models.Employee

	err := r.DB.Where("tenant_id = ? AND manager_id = ? AND is_active = ?", tenantID, managerID, true).
		Preload("Department").
		Preload("JobTitle").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}
	return employees, nil
}

// GetEmployeeHierarchy retrieves full reporting hierarchy for an employee
func (r *EmployeeRepository) GetEmployeeHierarchy(tenantID string, empID uint) ([]models.Employee, error) {
	var hierarchy []models.Employee
	var employee models.Employee

	// Get the employee
	if err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, empID).First(&employee).Error; err != nil {
		return nil, err
	}

	hierarchy = append(hierarchy, employee)

	// Traverse up to get managers
	currentManagerID := employee.ManagerID
	for currentManagerID != nil {
		var manager models.Employee
		if err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, *currentManagerID).First(&manager).Error; err != nil {
			break
		}
		hierarchy = append(hierarchy, manager)
		currentManagerID = manager.ManagerID
	}

	return hierarchy, nil
}

// DeleteEmployee soft-deletes employee (deactivates)
func (r *EmployeeRepository) DeleteEmployee(tenantID string, empID uint) error {
	return r.DB.Model(&models.Employee{}).
		Where("tenant_id = ? AND id = ?", tenantID, empID).
		Update("is_active", false).Error
}

// GetEmployeeCount returns total employee count
func (r *EmployeeRepository) GetEmployeeCount(tenantID string, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if status, ok := filters["employment_status"]; ok {
		query = query.Where("employment_status = ?", status)
	}

	if err := query.Model(&models.Employee{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
