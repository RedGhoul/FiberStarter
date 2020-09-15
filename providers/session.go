package providers

import "github.com/gofiber/session"

var sessionP *session.Session

func SessionProvider() *session.Session {
	return sessionP
}

func SetUpSessionProvider(session *session.Session) {
	sessionP = session
}
