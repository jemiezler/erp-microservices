package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"erp/hr-service/internal/database"
	"erp/hr-service/internal/handlers"
	sharedLogger "erp/shared/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

const ServiceName = "HR-SERVICE"

func main() {
	database.InitDB(ServiceName)

	app := fiber.New()

	app.Use(logger.New(sharedLogger.GetConfig(ServiceName)))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	empHandler := &handlers.EmployeeHandler{
		DB:          database.DB,
		ServiceName: ServiceName,
	}
	stubHandler := &handlers.StubHandler{}

	hr := app.Group("/api/v1/hr/employees")
	hr.Get("/", empHandler.GetAll)
	hr.Post("/", empHandler.Create)
	hr.Get("/:id", empHandler.GetByID)
	hr.Patch("/:id", empHandler.Update)

	app.Get("/api/v1/hr/leaves/pending", stubHandler.GetPendingLeaves)
	app.Get("/api/v1/hr/attendance/stats", stubHandler.GetAttendanceStats)
	app.Get("/api/v1/hr/payroll/pending", stubHandler.GetPendingPayroll)

	go func() {
		// Listen only on localhost - accessible only through API Gateway
		if err := app.Listen("127.0.0.1:8081", fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	_ = app.Shutdown()
}
