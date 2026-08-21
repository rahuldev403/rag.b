package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rahuldev403/forgeflow/internal/database"
)

type CreateNodeRequest struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Config     map[string]interface{} `json:"config"`
	NextNodeId *string                `json:"next_node_id"`
}

func CreateWorkflowNode(c *fiber.Ctx) error {

	workflowId := c.Params("workflowId")

	var req CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	query := `
		INSERT INTO nodes (workflow_id, name, type, config, next_node_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	var nodeId string
	var createdAt time.Time

	err := database.DB.QueryRow(
		context.Background(),
		query,
		workflowId,
		req.Name,
		req.Type,
		req.Config,
		req.NextNodeId,
	).Scan(&nodeId, &createdAt)

	if err != nil {
		log.Printf("Failed to create node: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create workflow node",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Node added successfully",
		"node_id":     nodeId,
		"workflow_id": workflowId,
		"created_at":  createdAt,
	})
}
