package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	sharedLogger "erp/shared/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ServiceName = "FINANCE-SERVICE"

type Payroll struct {
	gorm.Model
	EmployeeID uint   `json:"employee_id" gorm:"uniqueIndex"`
	AccountNum string `json:"account_num"`
	BaseSalary float64 `json:"base_salary"`
	Status     string `json:"status" gorm:"default:pending"`
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

func main() {
	initDatabase()
	app := fiber.New()

	app.Use(logger.New(sharedLogger.GetConfig(ServiceName)))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	finance := app.Group("/api/v1/finance")

	finance.Post("/webhook/employee-created", func(c *fiber.Ctx) error {
		type Payload struct {
			EmployeeID uint `json:"employee_id"`
		}
		var p Payload
		if err := c.BodyParser(&p); err != nil {
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

	finance.Get("/payroll/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var p Payroll
		if result := DB.First(&p, "employee_id = ?", id); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payroll not found"})
		}
		return c.JSON(p)
	})

	go func() {
		if err := app.Listen(":8082"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	_ = app.Shutdown()
}
