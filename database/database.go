package database

import (
	"log"
	"os"
	"time"

	"github.com/RedGhoul/bookshelf/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBConn *gorm.DB
)

func InitDb() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	var err error
	dsn := "user=gorm password=gorm dbname=bookshelf port=5432 sslmode=disable"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to the database")
	}
	log.Println("Connected to DB")
	DBConn.AutoMigrate(&models.Book{})
	DBConn.AutoMigrate(&models.User{})
	log.Println("Ran Auto Migrate")
}
