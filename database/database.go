package database

import (
	"StockTrack/models"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
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
	db_connection_url, err := url.Parse("sqlserver://" + os.Getenv("DBUSER") + ":" + url.QueryEscape(os.Getenv("DBPASSWORD")) + "@" + os.Getenv("DBPORT") + "?database=" + os.Getenv("DBNAME"))
	if DebugFlag {
		DBConn, err = gorm.Open(sqlserver.Open(db_connection_url.String()), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		DBConn, err = gorm.Open(sqlserver.Open(db_connection_url.String()), &gorm.Config{})
	}

	if err != nil {
		log.Panic("Failed to connect to the database")
	}
	log.Println("Connected to DB")
	DBConn.AutoMigrate(&models.Book{})
	DBConn.AutoMigrate(&models.User{})
	log.Println("Ran Auto Migrate")
}
