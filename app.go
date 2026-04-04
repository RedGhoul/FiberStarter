package main

import (
	"FiberStarter/database"
	"FiberStarter/middleware"
	"FiberStarter/providers"
	"FiberStarter/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func StartUp() {
	providers.InitializeAppConfiguration()
	database.InitDb()
	app := setupFiberEngine()
	routes.SetupRoutes(app)
	middleware.SetupMiddleware(app)
	setupGracefulShutdown(app)
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

func setupGracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting on port: " + providers.AppConfig.Server_Port)
		if err := app.Listen(":" + providers.AppConfig.Server_Port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
