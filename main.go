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
	//TODO - Got to add some sort of ENV config system
	database.InitDb()
	providers.SetUpSessionProvider(session.New())
	providers.SetUpHashProvider()
	engine := django.New("./views", ".django")
	engine.Reload(true)
	engine.Debug(true)
	app := fiber.New(&fiber.Settings{
		Views: engine,
	})
	middleware.SetupMiddleware(app)
	routes.SetupRoutes(app)
	app.Listen("localhost:3000")
}
