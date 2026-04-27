package database

import (
	"os"
	"erp/hr-service/internal/models"
	sharedLogger "erp/shared/logger"
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

	// Automigrate
	DB.AutoMigrate(&models.Employee{}, &models.AuditLog{})
}
