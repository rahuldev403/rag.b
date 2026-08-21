package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rahuldev403/forgeflow/internal/database"
)

func TriggerWebhook(c *fiber.Ctx) error {
	workflowID := c.Params("workflowId")

	var payload map[string]interface{}
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON payload",
			})
		}
	} else {
		payload = make(map[string]interface{})
	}

	query := `
		INSERT INTO executions (workflow_id, status, trigger_payload, started_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var executionID string
	now := time.Now()

	err := database.DB.QueryRow(
		context.Background(),
		query,
		workflowID,
		"PENDING",
		payload,
		now,
	).Scan(&executionID)

	if err != nil {
		log.Printf("Failed to create execution: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger workflow",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":      "Workflow execution triggered",
		"execution_id": executionID,
		"status":       "PENDING",
	})
}
