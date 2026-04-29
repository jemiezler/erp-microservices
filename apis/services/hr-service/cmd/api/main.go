package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "erp/hr-service/docs"
	"erp/hr-service/internal/database"
	"erp/hr-service/internal/handlers"
	sharedLogger "erp/shared/logger"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

const ServiceName = "HR-SERVICE"

// @title HR Service API
// @version 1.0
// @description This is the HR microservice for the ERP system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
func main() {
	database.InitDB(ServiceName)

	app := fiber.New()

	app.Use(logger.New(sharedLogger.GetConfig(ServiceName)))

	app.Get("/swagger/*", swaggo.HandlerDefault)

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
