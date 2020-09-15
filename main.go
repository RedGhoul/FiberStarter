package main

import (
	"github.com/RedGhoul/bookshelf/database"
	"github.com/RedGhoul/bookshelf/middleware"
	"github.com/RedGhoul/bookshelf/providers"
	"github.com/RedGhoul/bookshelf/routes"
	"github.com/gofiber/fiber"
	"github.com/gofiber/session"
	"github.com/gofiber/template/django"
)

func main() {
	database.InitDb()
	providers.SetUpSessionProvider(session.New())
	providers.SetUpHashProvider()
	app := fiber.New(&fiber.Settings{
		Views: setupViewEngine(),
	})
	middleware.SetupMiddleware(app)
	routes.SetupRoutes(app)
	app.Listen("localhost:3000")
}

func setupViewEngine() *django.Engine {
	return django.New("./views", ".django")
}
