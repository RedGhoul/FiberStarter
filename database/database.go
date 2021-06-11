package database

import (
	"log"
	"os"
	"strconv"
	"time"

	"fiberstarter/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBConn *gorm.DB
)

func InitDb() {
	configErr := godotenv.Load()
	if configErr != nil {
		log.Fatal("Error loading .env file")
	}

	DebugFlag, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	var newLogger logger.Interface
	if DebugFlag {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		)
	}

	var err error
	dsn := "user=" + os.Getenv("DBUSER") +
		" password=" + os.Getenv("DBPASSWORD") +
		" dbname=" + os.Getenv("DBNAME") +
		" port=" + os.Getenv("DBPORT") +
		" sslmode=" + os.Getenv("DBSSLMODE")
	if DebugFlag {
		DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Panic("Failed to connect to the database")
	}
	log.Println("Connected to DB")
	DBConn.AutoMigrate(&models.Book{})
	DBConn.AutoMigrate(&models.User{})
	log.Println("Ran Auto Migrate")
}
