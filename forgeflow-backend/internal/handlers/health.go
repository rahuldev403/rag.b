package handlers

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ForgeFlow API is running",
		"version": "1.0.0",
	})
}
