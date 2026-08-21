package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rahuldev403/forgeflow/internal/database"
)

type CreateWorkflowReqeust struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateWorkflow(c *fiber.Ctx) error {
	var req CreateWorkflowReqeust

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse JSON body",
		})
	}

	query := `
		INSERT INTO workflows (name, description, is_active) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	var id string
	var createdAt string

	err := database.DB.QueryRow(context.Background(), query, req.Name, req.Description, true).Scan(&id, &createdAt)

	if err != nil {
		log.Printf("DB Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert into database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Workflow created successfully",
		"id":         id,
		"created_at": createdAt,
	})
}
