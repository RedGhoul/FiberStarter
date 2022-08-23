package providers

import "github.com/gofiber/fiber/v2/middleware/session"

// stays up in the air
var sessionP *session.Store

func SessionProvider() *session.Store {
	return sessionP
}

func SetUpSessionProvider(session *session.Store) {
	sessionP = session
}
