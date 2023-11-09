package router

import (
	"github.com/gofiber/fiber/v2"
	"tes.go/m/handler"
)

func SetupRoutes(app *fiber.App) {
	// api := app.Group("/api")

	v2 := app.Group("/data")

	v2.Post("/", handler.CreateData)
	v2.Get("/", handler.GetAllData)

}
