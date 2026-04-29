package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "erp/finance-service/docs"
	sharedLogger "erp/shared/logger"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ServiceName = "FINANCE-SERVICE"

// Payroll model
// @Description Payroll information for an employee
type Payroll struct {
	gorm.Model
	EmployeeID uint    `json:"employee_id" gorm:"uniqueIndex"`
	AccountNum string  `json:"account_num"`
	BaseSalary float64 `json:"base_salary"`
	Status     string  `json:"status" gorm:"default:pending"`
}

var DB *gorm.DB

func initDatabase() {
	var err error
	dsn := "host=localhost user=erp_admin password=supersecretpassword dbname=finance_db port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sharedLogger.Error(ServiceName, "Failed to connect to database: %v", err)
		os.Exit(1)
	}
	sharedLogger.Success(ServiceName, "Database connection established and migrated.")
	DB.AutoMigrate(&Payroll{})
}

// @title Finance Service API
// @version 1.0
// @description This is the Finance microservice for the ERP system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
func main() {
	initDatabase()
	app := fiber.New()

	app.Use(logger.New(sharedLogger.GetConfig(ServiceName)))

	app.Get("/swagger/*", swaggo.HandlerDefault)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	finance := app.Group("/api/v1/finance")

	// EmployeeCreatedWebhook godoc
	// @Summary      Handle employee creation webhook
	// @Description  Initialize payroll account when a new employee is created
	// @Tags         finance
	// @Accept       json
	// @Produce      json
	// @Param        payload  body      object  true  "Employee creation payload"
	// @Success      200      {object}  Payroll
	// @Failure      400      {object}  map[string]string
	// @Failure      500      {object}  map[string]string
	// @Router       /api/v1/finance/webhook/employee-created [post]
	finance.Post("/webhook/employee-created", func(c fiber.Ctx) error {
		type Payload struct {
			EmployeeID uint `json:"employee_id"`
		}
		var p Payload
		if err := c.Bind().Body(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		payroll := Payroll{
			EmployeeID: p.EmployeeID,
			AccountNum: "ACC-" + string(rune(p.EmployeeID)), // Simplified for demo
			BaseSalary: 3000.0,
		}

		if result := DB.Create(&payroll); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
		}

		sharedLogger.Success(ServiceName, "Payroll account initialized for Employee ID: %d", p.EmployeeID)
		return c.Status(fiber.StatusOK).JSON(payroll)
	})

	// GetPayroll godoc
	// @Summary      Get payroll by employee ID
	// @Description  Retrieve payroll details for a specific employee
	// @Tags         finance
	// @Accept       json
	// @Produce      json
	// @Param        id   path      int  true  "Employee ID"
	// @Success      200  {object}  Payroll
	// @Failure      404  {object}  map[string]string
	// @Router       /api/v1/finance/payroll/{id} [get]
	finance.Get("/payroll/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		var p Payroll
		if result := DB.First(&p, "employee_id = ?", id); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payroll not found"})
		}
		return c.JSON(p)
	})

	go func() {
		if err := app.Listen(":8082", fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	_ = app.Shutdown()
}
