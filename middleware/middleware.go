package middleware

import (
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
)

func SetupMiddleware(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(middleware.Logger(os.Stdout))
	app.Use(middleware.Compress(middleware.CompressLevelBestSpeed))
	app.Static("/asset", "./public")
}
