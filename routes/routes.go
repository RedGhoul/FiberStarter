package routes

import (
	"log"

	"github.com/RedGhoul/fiberstarter/controllers"
	"github.com/RedGhoul/fiberstarter/providers"
	"github.com/gofiber/fiber"
)

func SetupRoutes(app *fiber.App) {
	setupBasicRoutes(app)
	setupAuthRoutes(app)
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

//app.Use(CheckAuth())
func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) {
		// Filter request to skip middleware
		if providers.IsAuthenticated(c) {
			c.Next()
			log.Println("IsAuthenticated")
		}
		c.SendStatus(404)
	}
}
