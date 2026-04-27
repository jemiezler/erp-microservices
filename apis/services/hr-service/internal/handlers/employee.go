package handlers

import (
	"erp/hr-service/internal/models"
	sharedLogger "erp/shared/logger"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type EmployeeHandler struct {
	DB          *gorm.DB
	ServiceName string
}

func (h *EmployeeHandler) GetAll(c fiber.Ctx) error {
	var employees []models.Employee
	h.DB.Preload("Manager").Find(&employees)
	return c.JSON(employees)
}

func (h *EmployeeHandler) Create(c fiber.Ctx) error {
	var emp models.Employee
	if err := c.Bind().Body(&emp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if result := h.DB.Create(&emp); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// Basic Audit Log for Creation
	h.DB.Create(&models.AuditLog{
		EmployeeID: emp.ID,
		Action:     "CREATE",
		Field:      "ALL",
		NewValue:   fmt.Sprintf("Employee %s created", emp.Name),
	})

	sharedLogger.Success(h.ServiceName, "Employee %s created successfully.", emp.Name)
	return c.Status(fiber.StatusCreated).JSON(emp)
}

func (h *EmployeeHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	var emp models.Employee
	if result := h.DB.Preload("Manager").First(&emp, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employee not found"})
	}
	return c.JSON(emp)
}

func (h *EmployeeHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var existing models.Employee
	if result := h.DB.First(&existing, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employee not found"})
	}

	var updateData models.Employee
	if err := c.Bind().Body(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Simple Audit for Status change
	if updateData.Status != "" && updateData.Status != existing.Status {
		h.DB.Create(&models.AuditLog{
			EmployeeID: existing.ID,
			Action:     "UPDATE",
			Field:      "Status",
			OldValue:   existing.Status,
			NewValue:   updateData.Status,
		})
	}

	h.DB.Model(&existing).Updates(updateData)
	return c.JSON(existing)
}
