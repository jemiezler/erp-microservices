package repository

import (
	"time"

	"erp/hr-service/internal/models"

	"gorm.io/gorm"
)

// PayrollRepository handles payroll data access
type PayrollRepository struct {
	*BaseRepository
}

// NewPayrollRepository creates new payroll repository
func NewPayrollRepository(db *gorm.DB) *PayrollRepository {
	return &PayrollRepository{
		BaseRepository: &BaseRepository{DB: db},
	}
}

// CreatePayroll creates new payroll record
func (r *PayrollRepository) CreatePayroll(payroll *models.Payroll) error {
	return r.DB.Create(payroll).Error
}

// GetPayrollByMonth retrieves payroll for specific month
func (r *PayrollRepository) GetPayrollByMonth(tenantID string, empID uint, month time.Time) (*models.Payroll, error) {
	var payroll models.Payroll

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND EXTRACT(YEAR FROM payroll_month) = ? AND EXTRACT(MONTH FROM payroll_month) = ?",
		tenantID, empID, month.Year(), month.Month()).
		Preload("Lines").
		First(&payroll).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &payroll, nil
}

// GetEmployeePayrollHistory retrieves payroll history for employee
func (r *PayrollRepository) GetEmployeePayrollHistory(tenantID string, empID uint, startDate, endDate time.Time) ([]models.Payroll, error) {
	var payrolls []models.Payroll

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND payroll_month BETWEEN ? AND ?",
		tenantID, empID, startDate, endDate).
		Preload("Lines").
		Order("payroll_month DESC").
		Find(&payrolls).Error

	if err != nil {
		return nil, err
	}
	return payrolls, nil
}

// GetSalaryStructure retrieves salary structure for employee
func (r *PayrollRepository) GetSalaryStructure(tenantID string, empID uint) (*models.EmployeeSalary, error) {
	var empSalary models.EmployeeSalary
	now := time.Now()

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND effective_from <= ? AND (effective_to IS NULL OR effective_to > ?)",
		tenantID, empID, now, now).
		Preload("SalaryStructure").
		Preload("ComponentValues").
		First(&empSalary).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &empSalary, nil
}

// GetPayrollStats retrieves payroll statistics
func (r *PayrollRepository) GetPayrollStats(tenantID string, month time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalPayrolls int64
	var processedPayrolls int64
	var pendingApprovals int64
	var totalEarnings float64
	var totalDeductions float64
	var totalNetPayable float64

	query := r.DB.Where(
		"tenant_id = ? AND EXTRACT(YEAR FROM payroll_month) = ? AND EXTRACT(MONTH FROM payroll_month) = ?",
		tenantID, month.Year(), month.Month())

	query.Model(&models.Payroll{}).Count(&totalPayrolls)
	query.Where("status = ?", "Posted").Model(&models.Payroll{}).Count(&processedPayrolls)
	query.Where("status IN (?, ?)", "Draft", "Generated").Model(&models.Payroll{}).Count(&pendingApprovals)

	query.Model(&models.Payroll{}).Select("COALESCE(SUM(earnings), 0)").Row().Scan(&totalEarnings)
	query.Model(&models.Payroll{}).Select("COALESCE(SUM(deductions), 0)").Row().Scan(&totalDeductions)
	query.Model(&models.Payroll{}).Select("COALESCE(SUM(net_payable), 0)").Row().Scan(&totalNetPayable)

	stats["total_payrolls"] = totalPayrolls
	stats["processed"] = processedPayrolls
	stats["pending_approvals"] = pendingApprovals
	stats["total_earnings"] = totalEarnings
	stats["total_deductions"] = totalDeductions
	stats["total_net_payable"] = totalNetPayable

	return stats, nil
}

// ApprovePayroll approves payroll
func (r *PayrollRepository) ApprovePayroll(tenantID string, payrollID uint, approverID uint) error {
	now := time.Now()
	return r.DB.Model(&models.Payroll{}).
		Where("tenant_id = ? AND id = ?", tenantID, payrollID).
		Updates(map[string]interface{}{
			"status":         "Approved",
			"approved_by_id": approverID,
			"approved_date":  now,
		}).Error
}

// PostPayroll posts approved payroll to finance
func (r *PayrollRepository) PostPayroll(tenantID string, payrollID uint) error {
	return r.DB.Model(&models.Payroll{}).
		Where("tenant_id = ? AND id = ?", tenantID, payrollID).
		Update("status", "Posted").Error
}

// GetUnprocessedPayrolls gets payrolls that need processing
func (r *PayrollRepository) GetUnprocessedPayrolls(tenantID string) ([]models.Payroll, error) {
	var payrolls []models.Payroll

	err := r.DB.Where(
		"tenant_id = ? AND status IN (?, ?)",
		tenantID, "Draft", "Generated").
		Preload("Employee").
		Find(&payrolls).Error

	if err != nil {
		return nil, err
	}
	return payrolls, nil
}

// GetBenefits retrieves employee benefits
func (r *PayrollRepository) GetBenefits(tenantID string, empID uint) ([]models.EmployeeBenefit, error) {
	var benefits []models.EmployeeBenefit
	now := time.Now()

	err := r.DB.Where(
		"tenant_id = ? AND employee_id = ? AND effective_from <= ? AND (effective_to IS NULL OR effective_to > ?) AND status = ?",
		tenantID, empID, now, now, "Active").
		Preload("BenefitPlan").
		Find(&benefits).Error

	if err != nil {
		return nil, err
	}
	return benefits, nil
}

// SubmitBenefitClaim submits a benefit claim
func (r *PayrollRepository) SubmitBenefitClaim(claim *models.BenefitClaim) error {
	return r.DB.Create(claim).Error
}

// GetPendingBenefitClaims gets pending claims for approval
func (r *PayrollRepository) GetPendingBenefitClaims(tenantID string) ([]models.BenefitClaim, error) {
	var claims []models.BenefitClaim

	err := r.DB.Where("tenant_id = ? AND status = ?", tenantID, "Submitted").
		Preload("EmployeeBenefit").
		Preload("EmployeeBenefit.EmployeeBenefit.Employee").
		Order("claim_date DESC").
		Find(&claims).Error

	if err != nil {
		return nil, err
	}
	return claims, nil
}

// GetPaymentMode retrieves payment details for payroll
func (r *PayrollRepository) GetPaymentMode(tenantID string, payrollID uint) (map[string]interface{}, error) {
	var payroll models.Payroll

	err := r.DB.Where("tenant_id = ? AND id = ?", tenantID, payrollID).
		First(&payroll).Error

	if err != nil {
		return nil, err
	}

	paymentDetails := map[string]interface{}{
		"mode":           payroll.PaymentMode,
		"bank_details":   payroll.BankDetails,
		"net_payable":    payroll.NetPayable,
		"payment_date":   payroll.PaymentDate,
		"transaction_id": payroll.TransactionID,
	}

	return paymentDetails, nil
}

// UpdatePaymentStatus updates payment status after posting
func (r *PayrollRepository) UpdatePaymentStatus(tenantID string, payrollID uint, status string, transactionID string) error {
	now := time.Now()
	return r.DB.Model(&models.Payroll{}).
		Where("tenant_id = ? AND id = ?", tenantID, payrollID).
		Updates(map[string]interface{}{
			"status":         status,
			"transaction_id": transactionID,
			"payment_date":   now,
		}).Error
}
