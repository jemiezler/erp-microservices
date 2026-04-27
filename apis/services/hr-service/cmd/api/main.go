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

	hr := app.Group("/api/v1/hr/employees")
	hr.Get("/", empHandler.GetAll)
	hr.Post("/", empHandler.Create)
	hr.Get("/:id", empHandler.GetByID)
	hr.Patch("/:id", empHandler.Update)

	go func() {
		if err := app.Listen(":8081", fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	_ = app.Shutdown()
}