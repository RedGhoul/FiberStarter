package database

import (
	"FiberStarter/providers"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBConn *gorm.DB
)

func InitDb() {
	var newLogger logger.Interface
	if providers.AppConfig.Enable_Debug {
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

	if err != nil {
		log.Panic(err)
	}
	if providers.AppConfig.Enable_Debug {
		DBConn, err = gorm.Open(mysql.Open(providers.AppConfig.DB_Connection_URL), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		DBConn, err = gorm.Open(mysql.Open(providers.AppConfig.DB_Connection_URL), &gorm.Config{})
	}

	if err != nil {
		log.Panic("Failed to connect to the database " + providers.AppConfig.DB_Connection_URL)
	}
}
