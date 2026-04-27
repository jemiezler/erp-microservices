package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sharedLogger "erp/shared/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout: 10 * time.Second,
	})

	app.Use(logger.New(sharedLogger.GetConfig("API-GATEWAY")))

	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}
		return c.Next()
	})

	app.All("/api/v1/hr/*", func(c *fiber.Ctx) error {
		target := "http://localhost:8081" + c.Path()
		if err := proxy.Do(c, target); err != nil {
			return err
		}
		return nil
	})

	app.All("/api/v1/finance/*", func(c *fiber.Ctx) error {
		target := "http://localhost:8082" + c.Path()
		if err := proxy.Do(c, target); err != nil {
			return err
		}
		return nil
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	_ = app.Shutdown()
}
