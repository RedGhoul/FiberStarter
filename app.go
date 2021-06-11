package main

import (
	"log"
	"os"
	"strconv"

	"fiberstarter/middleware"
	"fiberstarter/providers"
	"fiberstarter/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/django"
	"github.com/joho/godotenv"
)

func StartUp() {
	//database.InitDb()
	providers.SetUpSessionProvider(session.New())
	providers.SetUpHashProvider()
	app := setupViewEngine()
	middleware.SetupMiddleware(app)
	routes.SetupRoutes(app)
	setupAppListenPort(app)
}

func setupViewEngine() *fiber.App {
	engine := django.New("./views", ".django")

	configErr := godotenv.Load()
	if configErr != nil {
		log.Fatal("Error loading .env file")
	}

	DebugFlag, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if DebugFlag {
		engine.Reload(true)
		engine.Debug(true)
	}

	return fiber.New(fiber.Config{
		Views: engine,
	})
}

func setupAppListenPort(app *fiber.App) {
	app.Listen("localhost:" + os.Getenv("SERVERPORT"))
}
