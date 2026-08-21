package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rahuldev403/forgeflow/internal/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("api/v1")

	api.Get("/helath", handlers.HealthCheck)

	api.Post("/workflows", handlers.CreateWorkflow)
}
