package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	sharedLogger "erp/shared/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

const ServiceName = "API-GATEWAY"

// ServiceConfig holds backend service URLs
type ServiceConfig struct {
	HR      string
	Finance string
	Auth    string
}

// getServiceConfig returns service URLs from environment or defaults
func getServiceConfig() ServiceConfig {
	return ServiceConfig{
		HR:      getEnv("HR_SERVICE_URL", "http://localhost:8081"),
		Finance: getEnv("FINANCE_SERVICE_URL", "http://localhost:8082"),
		Auth:    getEnv("AUTH_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func buildCORSConfig() cors.Config {
	originsEnv := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000,http://127.0.0.1:3001")
	corsOrigins := strings.Split(originsEnv, ",")
	containsWildcard := false

	// Trim whitespace and detect wildcard mode.
	for i, origin := range corsOrigins {
		trimmed := strings.TrimSpace(origin)
		corsOrigins[i] = trimmed
		if trimmed == "*" {
			containsWildcard = true
		}
	}

	corsMethods := []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete, fiber.MethodOptions}
	corsHeaders := []string{"Content-Type", "Authorization", "X-Tenant-ID", "X-Request-ID"}

	config := cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: corsMethods,
		AllowHeaders: corsHeaders,
		MaxAge:       3600,
	}

	// Browsers reject Access-Control-Allow-Credentials=true with wildcard origin.
	if containsWildcard {
		config.AllowCredentials = false
	} else {
		config.AllowCredentials = true
	}

	return config
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

func applyProxyCORSHeaders(c fiber.Ctx, corsConfig cors.Config) {
	origin := c.Get("Origin")
	if origin == "" {
		return
	}

	if !isOriginAllowed(origin, corsConfig.AllowOrigins) {
		return
	}

	c.Set("Vary", "Origin")
	c.Set("Access-Control-Allow-Origin", origin)
	if corsConfig.AllowCredentials {
		c.Set("Access-Control-Allow-Credentials", "true")
	}
}

func main() {
	config := getServiceConfig()
	corsConfig := buildCORSConfig()

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	})

	// Request ID middleware - generates unique ID for each request
	app.Use(requestid.New())

	// Logger middleware
	app.Use(logger.New(sharedLogger.GetConfig(ServiceName)))

	// CORS middleware - allow frontend requests
	app.Use(cors.New(corsConfig))

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// Gateway info endpoint
	app.Get("/info", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": ServiceName,
			"version": "1.0.0",
			"services": fiber.Map{
				"hr":      config.HR,
				"finance": config.Finance,
				"auth":    config.Auth,
			},
		})
	})

	// HR Service routes - proxy to hr-service
	app.All("/api/v1/hr/*", func(c fiber.Ctx) error {
		path := c.Path()
		target := config.HR + path
		log.Printf("[%s] HR Service: %s %s", c.Get("X-Request-ID"), c.Method(), path)

		if err := proxy.Do(c, target); err != nil {
			applyProxyCORSHeaders(c, corsConfig)
			log.Printf("[%s] HR Service Error: %v", c.Get("X-Request-ID"), err)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":      "Failed to reach HR service",
				"request_id": c.Get("X-Request-ID"),
			})
		}

		applyProxyCORSHeaders(c, corsConfig)
		return nil
	})

	// Finance Service routes - proxy to finance-service
	app.All("/api/v1/finance/*", func(c fiber.Ctx) error {
		path := c.Path()
		target := config.Finance + path
		log.Printf("[%s] Finance Service: %s %s", c.Get("X-Request-ID"), c.Method(), path)

		if err := proxy.Do(c, target); err != nil {
			applyProxyCORSHeaders(c, corsConfig)
			log.Printf("[%s] Finance Service Error: %v", c.Get("X-Request-ID"), err)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":      "Failed to reach Finance service",
				"request_id": c.Get("X-Request-ID"),
			})
		}

		applyProxyCORSHeaders(c, corsConfig)
		return nil
	})

	// 404 handler
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "Endpoint not found",
			"path":       c.Path(),
			"request_id": c.Get("X-Request-ID"),
		})
	})

	// Start server
	go func() {
		port := getEnv("GATEWAY_PORT", "8080")
		log.Printf("Starting %s on port %s", ServiceName, port)
		log.Printf("HR Service: %s", config.HR)
		log.Printf("Finance Service: %s", config.Finance)

		if err := app.Listen(":"+port, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			log.Panic(err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("\nShutting down gateway...")
	_ = app.Shutdown()
}
