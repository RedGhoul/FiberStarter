package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RedGhoul/bookshelf/controllers"
	"github.com/RedGhoul/bookshelf/database"
	"github.com/RedGhoul/bookshelf/models"
	"github.com/RedGhoul/bookshelf/providers"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
	"github.com/gofiber/session"
	"github.com/gofiber/template/django"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func helloWorld(c *fiber.Ctx) {
	if err := c.Render("Home/index", fiber.Map{"Title": "Hello, World!"}); err != nil {
		c.Status(500).Send(err.Error())
	}
}
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
func setupNonAuthRoutes(app *fiber.App) {
	// app.Get("/api/v1/book", models.GetBooks)
	// app.Get("/api/v1/book/:id", models.GetBook)
	// app.Post("/api/v1/book", models.NewBook)
	// app.Delete("/api/v1/book/:id", models.DeleteBooks)

	app.Get("/Login", controllers.ShowLoginForm)
	app.Post("/Login", controllers.PostLoginForm)
	app.Get("/Logout", controllers.PostLogoutForm)
	app.Get("/Register", controllers.ShowRegisterForm)
	app.Post("/Register", controllers.PostRegisterForm)
}

func initDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // Disable color
		},
	)
	var err error
	dsn := "user=gorm password=gorm dbname=bookshelf port=5432 sslmode=disable"
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to the database")
	}
	fmt.Println("Connected to DB")
	// removed close() in V2

	database.DBConn.AutoMigrate(&models.Book{})
	database.DBConn.AutoMigrate(&models.User{})
	fmt.Println("Ran Auto Migrate")
}

func main() {
	engine := django.New("./views", ".django")
	app := fiber.New(&fiber.Settings{
		Views: engine,
	})
	initDatabase()
	providers.SetSessionProvider(session.New())
	providers.SetHashProvider()
	app.Use(middleware.Logger(os.Stdout))
	app.Use(helmet.New())
	app.Use(middleware.Compress(middleware.CompressLevelBestSpeed))
	// https://blog.logrocket.com/express-style-api-go-fiber/
	app.Static("/asset", "./public")

	setupNonAuthRoutes(app)
	app.Use(CheckAuth())
	app.Get("/Home", helloWorld)
	app.Listen("localhost:3000")
}
