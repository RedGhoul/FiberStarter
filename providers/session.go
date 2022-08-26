package providers

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mysql"
)

// stays up in the air
var sessionP *session.Store

func SessionProvider() *session.Store {
	return sessionP
}

func SetUpSessionProvider(url string) {
	storage := mysql.New(mysql.Config{
		ConnectionURI: url,
		Table:         "session_storage",
		Reset:         false,
		GCInterval:    10 * time.Second,
	})
	sessionP = session.New(session.Config{
		Storage: storage,
	})
}
