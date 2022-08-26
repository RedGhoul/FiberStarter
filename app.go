package main

import (
	"FiberStarter/database"
	"FiberStarter/middleware"
	"FiberStarter/providers"
	"FiberStarter/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
)

func StartUp() {
	providers.InitializeAppConfiguration()
	database.InitDb()
	app := setupFiberEngine()
	routes.SetupRoutes(app)
	middleware.SetupMiddleware(app)
	setupAppListenPort(app)
}

func setupFiberEngine() *fiber.App {
	engine := django.New("./views", ".django")
	engine.Reload(providers.AppConfig.Enable_Debug)
	engine.Debug(providers.AppConfig.Enable_Debug)

	return fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		Prefork:           providers.AppConfig.Enable_PreFork,
		AppName:           providers.AppConfig.App_Name,
		Concurrency:       providers.AppConfig.Concurrency_Level,
		StrictRouting:     true,
		CaseSensitive:     true,
		PassLocalsToViews: true,
	})
}

func setupAppListenPort(app *fiber.App) {
	log.Println("Starting on port: " + providers.AppConfig.Server_Port)
	app.Listen(":" + providers.AppConfig.Server_Port)
}
