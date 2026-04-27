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

const ServiceName = "HR-SERVICE"

type Employee struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Position string `json:"position"`
	Status   string `json:"status" gorm:"default:active"`
}

var DB *gorm.DB

func initDatabase() {
	var err error
	dsn := "host=localhost user=erp_admin password=supersecretpassword dbname=hr_db port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sharedLogger.Error(ServiceName, "Failed to connect to database: %v", err)
		os.Exit(1)
	}
	sharedLogger.Success(ServiceName, "Database connection established and migrated.")
	DB.AutoMigrate(&Employee{})
}

func main() {
	initDatabase()
	app := fiber.New()

	app.Use(logger.New(sharedLogger.GetConfig(ServiceName)))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	hr := app.Group("/api/v1/hr/employees")
	hr.Get("/", func(c *fiber.Ctx) error {
		var employees []Employee
		DB.Find(&employees)
		return c.JSON(employees)
	})

	hr.Post("/", func(c *fiber.Ctx) error {
		var emp Employee
		if err := c.BodyParser(&emp); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		if result := DB.Create(&emp); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
		}
		sharedLogger.Success(ServiceName, "Employee %s created successfully.", emp.Name)
		return c.Status(fiber.StatusCreated).JSON(emp)
	})

	go func() {
		if err := app.Listen(":8081"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	_ = app.Shutdown()
}
