package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SalaryStructure defines salary components
type SalaryStructure struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name        string `json:"name"` // Standard, Management, Contractual
	Description string `json:"description"`

	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	Components []SalaryComponent `json:"components,omitempty"`

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// SalaryComponent individual salary component (Basic, HRA, DA, etc.)
type SalaryComponent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	SalaryStructureID uint             `json:"salary_structure_id" gorm:"index"`
	SalaryStructure   *SalaryStructure `json:"salary_structure,omitempty"`

	ComponentName string `json:"component_name"` // Basic, HRA, DA, Bonus, Allowance
	ComponentCode string `json:"component_code"` // BSC, HRA, DA

	Type       string  `json:"type"` // Earning, Deduction
	IsFormula  bool    `json:"is_formula"`
	Formula    string  `json:"formula"`    // e.g., "BASIC * 0.25"
	Amount     float64 `json:"amount"`     // If not formula-based
	Percentage float64 `json:"percentage"` // If percentage-based

	IsTaxable       bool `json:"is_taxable"`
	IsProvidentFund bool `json:"is_provident_fund"`
	IsEsicIncluded  bool `json:"is_esic_included"`

	DisplayOrder int `json:"display_order"`

	gorm.Model
}

// EmployeeSalary employee salary configuration
type EmployeeSalary struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"uniqueIndex:idx_emp_sal_eff"`
	Employee   *Employee `json:"employee,omitempty"`

	SalaryStructureID uint             `json:"salary_structure_id" gorm:"index"`
	SalaryStructure   *SalaryStructure `json:"salary_structure,omitempty"`

	EffectiveFrom time.Time  `json:"effective_from" gorm:"uniqueIndex:idx_emp_sal_eff"`
	EffectiveTo   *time.Time `json:"effective_to"`

	BaseSalary float64 `json:"base_salary"`
	CTC        float64 `json:"ctc"`

	ComponentValues []SalaryComponentValue `json:"component_values,omitempty"`

	gorm.Model
}

// SalaryComponentValue calculated component values for employee
type SalaryComponentValue struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EmployeeSalaryID uint            `json:"employee_salary_id" gorm:"index"`
	EmployeeSalary   *EmployeeSalary `json:"employee_salary,omitempty"`

	ComponentID uint             `json:"component_id"`
	Component   *SalaryComponent `json:"component,omitempty"`

	Value float64 `json:"value"`

	gorm.Model
}

// Payroll monthly payroll records
type Payroll struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	PayrollNumber string    `json:"payroll_number" gorm:"uniqueIndex"`
	PayrollMonth  time.Time `json:"payroll_month" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	PaymentDate *time.Time `json:"payment_date"`

	WorkingDays      int     `json:"working_days"`
	LeaveTaken       int     `json:"leave_taken"`
	ActualDaysWorked float64 `json:"actual_days_worked"` // After adjustments

	Earnings   float64 `json:"earnings"`
	Deductions float64 `json:"deductions"`
	NetPayable float64 `json:"net_payable"` // Earnings - Deductions

	Lines []PayrollLine `json:"lines,omitempty"`

	Status string `json:"status"` // Draft, Generated, Approved, Posted, Paid

	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedDate *time.Time `json:"approved_date"`

	PaymentMode   string `json:"payment_mode"` // Bank Transfer, Cheque, Cash
	BankDetails   string `json:"bank_details"`
	TransactionID string `json:"transaction_id"`

	Remarks string `json:"remarks"`

	gorm.Model
}

// PayrollLine individual payroll component lines
type PayrollLine struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PayrollID uint     `json:"payroll_id" gorm:"index"`
	Payroll   *Payroll `json:"payroll,omitempty"`

	ComponentID uint             `json:"component_id"`
	Component   *SalaryComponent `json:"component,omitempty"`

	ComponentName string  `json:"component_name"`
	Amount        float64 `json:"amount"`

	gorm.Model
}

// BenefitPlan employee benefit plans
type BenefitPlan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name        string `json:"name"` // Health Insurance, Life Insurance, Retirement
	Type        string `json:"type"` // Medical, Insurance, Retirement, Wellness
	Description string `json:"description"`

	Provider string `json:"provider"`

	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	CoverageAmount  float64 `json:"coverage_amount"`
	EmployeePremium float64 `json:"employee_premium"` // Employee pays
	CompanyPremium  float64 `json:"company_premium"`  // Company pays

	Benefits   datatypes.JSONSlice[string] `json:"benefits" gorm:"type:jsonb"` // Inclusions
	Exclusions datatypes.JSONSlice[string] `json:"exclusions" gorm:"type:jsonb"`

	DocumentURL string `json:"document_url"`

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// EmployeeBenefit employee enrollment in benefit plans
type EmployeeBenefit struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	BenefitPlanID uint         `json:"benefit_plan_id" gorm:"index"`
	BenefitPlan   *BenefitPlan `json:"benefit_plan,omitempty"`

	EnrollmentDate time.Time  `json:"enrollment_date"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to"`

	PolicyNumber string `json:"policy_number"`

	CoverageType      string `json:"coverage_type"` // Individual, Family
	DependentsCovered int    `json:"dependents_covered"`

	Status string `json:"status"` // Active, Inactive, Expired

	ClaimsSubmitted []BenefitClaim `json:"claims,omitempty"`

	gorm.Model
}

// BenefitClaim insurance/benefit claims
type BenefitClaim struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	ClaimNumber       string           `json:"claim_number" gorm:"uniqueIndex"`
	EmployeeBenefitID uint             `json:"employee_benefit_id" gorm:"index"`
	EmployeeBenefit   *EmployeeBenefit `json:"employee_benefit,omitempty"`

	ClaimType string    `json:"claim_type"` // Medical, Death, Hospitalization, etc
	ClaimDate time.Time `json:"claim_date"`

	ClaimAmount    float64 `json:"claim_amount"`
	ApprovedAmount float64 `json:"approved_amount"`

	Description string `json:"description" gorm:"type:text"`

	DocumentURLs datatypes.JSONSlice[string] `json:"document_urls" gorm:"type:jsonb"`

	Status string `json:"status"` // Submitted, Under Review, Approved, Rejected, Paid

	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedDate *time.Time `json:"approved_date"`

	RejectionReason string `json:"rejection_reason"`

	PaymentDate *time.Time `json:"payment_date"`

	gorm.Model
}
