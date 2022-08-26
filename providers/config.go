package providers

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Enable_Debug       bool
	Enable_PreFork     bool
	Enable_Compression bool
	Enable_Recover     bool
	Concurrency_Level  int
	App_Name           string
	DB_Connection_URL  string
	Server_Port        string
}

var AppConfig Config

func InitializeAppConfiguration() {
	configErr := godotenv.Load()
	logError(configErr)
	AppConfig = Config{}
	AppConfig.Enable_Debug, configErr = strconv.ParseBool(os.Getenv("ENABLE_DEBUG"))
	logError(configErr)
	AppConfig.Enable_PreFork, configErr = strconv.ParseBool(os.Getenv("ENABLE_PREFORK"))
	logError(configErr)
	AppConfig.Enable_Compression, configErr = strconv.ParseBool(os.Getenv("ENABLE_COMPRESSION"))
	logError(configErr)
	AppConfig.Enable_Recover, configErr = strconv.ParseBool(os.Getenv("ENABLE_RECOVERY"))
	logError(configErr)
	AppConfig.Concurrency_Level, configErr = strconv.Atoi(os.Getenv("CONCURRENCY_LEVEL"))
	logError(configErr)
	AppConfig.App_Name = os.Getenv("APP_NAME")
	if url, err := url.Parse(os.Getenv("DB_USER") + ":" + url.QueryEscape(os.Getenv("DB_PASSWORD")) + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")" + "/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"); url != nil {
		AppConfig.DB_Connection_URL = url.String()
		logError(err)
	}
	AppConfig.Server_Port = os.Getenv("SERVER_PORT")
}

func logError(err error) {
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
