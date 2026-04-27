package database

import (
	"erp/hr-service/internal/models"
	sharedLogger "erp/shared/logger"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(serviceName string) {
	var err error
	dsn := "host=localhost user=erp_admin password=supersecretpassword dbname=hr_db port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sharedLogger.Error(serviceName, "Failed to connect to database: %v", err)
		os.Exit(1)
	}

	sharedLogger.Success(serviceName, "Database connection established.")

	// AutoMigrate all HR modules
	// Core HR
	DB.AutoMigrate(
		&models.Organization{},
		&models.Department{},
		&models.JobTitle{},
		&models.Location{},
		&models.EmployeeHierarchy{},
		&models.Employee{},
		&models.Dependent{},
		&models.EmergencyContact{},
		&models.Qualification{},
		&models.WorkExperience{},
	)

	// Recruitment Module
	DB.AutoMigrate(
		&models.JobPosting{},
		&models.Candidate{},
		&models.Interview{},
		&models.JobOffer{},
		&models.OnboardingTask{},
	)

	// Attendance & Time Tracking
	DB.AutoMigrate(
		&models.Attendance{},
		&models.TimeSheet{},
		&models.TimeSheetEntry{},
		&models.Shift{},
		&models.ShiftAssignment{},
		&models.BiometricData{},
		&models.WorkFromHome{},
	)

	// Leave Management
	DB.AutoMigrate(
		&models.LeaveType{},
		&models.LeavePolicy{},
		&models.LeaveAllocation{},
		&models.Leave{},
		&models.LeaveApproval{},
		&models.HolidayCalendar{},
		&models.Holiday{},
		&models.LeaveEncashment{},
	)

	// Performance Management
	DB.AutoMigrate(
		&models.PerformanceGoal{},
		&models.GoalProgress{},
		&models.PerformanceRating{},
		&models.CompetencyRating{},
		&models.Competency{},
		&models.EmployeeReview{},
		&models.ReviewQuestion{},
		&models.TrainingRequest{},
	)

	// Compensation & Benefits
	DB.AutoMigrate(
		&models.SalaryStructure{},
		&models.SalaryComponent{},
		&models.EmployeeSalary{},
		&models.SalaryComponentValue{},
		&models.Payroll{},
		&models.PayrollLine{},
		&models.BenefitPlan{},
		&models.EmployeeBenefit{},
		&models.BenefitClaim{},
	)

	// Workflow & Approval Engine
	DB.AutoMigrate(
		&models.ApprovalWorkflow{},
		&models.ApprovalLevel{},
		&models.ApprovalRequest{},
		&models.Approval{},
	)

	// RBAC
	DB.AutoMigrate(
		&models.RBACRole{},
		&models.Permission{},
	)

	// Learning Management System
	DB.AutoMigrate(
		&models.LMSCourse{},
		&models.CourseModule{},
		&models.CourseEnrollment{},
		&models.LessonProgress{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.QuizAttempt{},
		&models.QuizAnswer{},
		&models.Certificate{},
	)

	// Audit & System
	DB.AutoMigrate(
		&models.AuditLog{},
		&models.SystemSetting{},
	)

	sharedLogger.Success(serviceName, "All database migrations completed.")
}
