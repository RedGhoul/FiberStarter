package middleware

import (
	"FiberStarter/providers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/utils"
)

func SetupMiddleware(app *fiber.App) {
	loggingMiddleware(app)
	staticAssetMiddleware(app)
	sessionMiddleware(app)
	setup404Handler(app)
	setupRecovery(app)
	setupCompression(app)
}

func sessionMiddleware(app *fiber.App) {
	providers.SetUpSessionProvider(providers.AppConfig.DB_Connection_URL)
	providers.SetUpHashProvider()
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))

}
func loggingMiddleware(app *fiber.App) {
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
}

func staticAssetMiddleware(app *fiber.App) {
	app.Static("/asset", "./public")
	app.Use(favicon.New())
	app.Use(etag.New())
}

func setup404Handler(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		var err error
		if err := c.SendStatus(fiber.StatusNotFound); err != nil {
			panic(err)
		}
		if err := c.Render("common/soft_error", fiber.Map{}); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return err
	})
}

func setupRecovery(app *fiber.App) {
	if providers.AppConfig.Enable_Recover {
		app.Use(recover.New())
	}
}

func setupCompression(app *fiber.App) {
	if providers.AppConfig.Enable_Compression {
		lvl := compress.Level(2)
		app.Use(compress.New(compress.Config{
			Level: lvl,
		}))
	}
}
