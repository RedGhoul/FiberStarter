package routes

import (
	"FiberStarter/controllers"
	"FiberStarter/database"
	"FiberStarter/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", healthCheck)
	app.Use(utils.AddLocals())
	setupAuthRoutes(app)
	setupBasicRoutes(app)
	app.Use(utils.CheckAuth())
	setupHiddenRoutes(app)
}

func healthCheck(c *fiber.Ctx) error {
	dbStatus := "ok"
	sqlDB, err := database.DBConn.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "unavailable"
	}
	return c.JSON(fiber.Map{
		"status":   "ok",
		"database": dbStatus,
	})
}

func setupAuthRoutes(app *fiber.App) {
	app.Get("/Login", controllers.ShowLoginForm)
	app.Post("/Login", controllers.PostLoginForm)
	app.Get("/Logout", controllers.PostLogoutForm)
	app.Get("/Register", controllers.ShowRegisterForm)
	app.Post("/Register", controllers.PostRegisterForm)
}

func setupBasicRoutes(app *fiber.App) {
	app.Get("/", controllers.ShowIndex)
}

func setupHiddenRoutes(app *fiber.App) {
	app.Get("/Secret", controllers.ShowSecret, utils.CheckAuth())
}
