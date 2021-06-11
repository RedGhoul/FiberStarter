package routes

import (
	"fiberstarter/controllers"
	"fiberstarter/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	setupAuthRoutes(app)
	setupBasicRoutes(app)
	app.Use(utils.CheckAuth())
	setupHiddenRoutes(app)
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
	app.Get("/Secrect", controllers.ShowSecrect)
}
