package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func SetupMiddleware(app *fiber.App) {
	app.Static("/asset", "./public")
}
