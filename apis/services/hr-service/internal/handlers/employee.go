package handlers

import (
	"encoding/json"
	"erp/hr-service/internal/models"
	sharedLogger "erp/shared/logger"

	"github.com/gofiber/fiber/v3"
	"gorm.io/datatypes"
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

	// Log audit trail
	newValuesMap := map[string]interface{}{
		"employee_id": emp.EmployeeID,
		"name":        emp.Name,
		"email":       emp.Email,
	}
	newValuesJSON, _ := json.Marshal(newValuesMap)
	h.DB.Create(&models.AuditLog{
		EntityType: "Employee",
		EntityID:   emp.ID,
		UserID:     emp.ID, // In real scenario, use authenticated user ID
		Action:     "Create",
		NewValues:  datatypes.JSON(newValuesJSON),
		Status:     "Success",
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

	// Log audit for status change
	if updateData.Status != "" && updateData.Status != existing.Status {
		oldValuesMap := map[string]interface{}{
			"status": existing.Status,
		}
		newValuesMap := map[string]interface{}{
			"status": updateData.Status,
		}
		oldValuesJSON, _ := json.Marshal(oldValuesMap)
		newValuesJSON, _ := json.Marshal(newValuesMap)
		h.DB.Create(&models.AuditLog{
			EntityType: "Employee",
			EntityID:   existing.ID,
			UserID:     existing.ID, // In real scenario, use authenticated user ID
			Action:     "Update",
			OldValues:  datatypes.JSON(oldValuesJSON),
			NewValues:  datatypes.JSON(newValuesJSON),
			Status:     "Success",
		})
	}

	h.DB.Model(&existing).Updates(updateData)
	return c.JSON(existing)
}
