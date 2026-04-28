package handlers

import "github.com/gofiber/fiber/v3"

type StubHandler struct{}

func (h *StubHandler) GetPendingLeaves(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"code":    "OK",
		"data": []fiber.Map{
			{
				"id":          1,
				"employee_id": "EMP-001",
				"leave_type":  "Annual Leave",
				"start_date":  "2026-05-01",
				"end_date":    "2026-05-03",
				"status":      "pending",
				"reason":      "Family event",
			},
		},
		"pagination": fiber.Map{
			"page":        1,
			"page_size":   20,
			"total":       1,
			"total_pages": 1,
		},
	})
}

func (h *StubHandler) GetAttendanceStats(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"code":    "OK",
		"data": fiber.Map{
			"present":               18,
			"absent":                2,
			"on_leave":              1,
			"total_employees":       21,
			"attendance_percentage": 85.7,
		},
	})
}

func (h *StubHandler) GetPendingPayroll(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"code":    "OK",
		"data": []fiber.Map{
			{
				"id":           1,
				"employee_id":  "EMP-001",
				"month":        "2026-04",
				"gross_salary": 5000,
				"deductions":   450,
				"net_salary":   4550,
				"status":       "pending",
			},
		},
		"pagination": fiber.Map{
			"page":        1,
			"page_size":   20,
			"total":       1,
			"total_pages": 1,
		},
	})
}
