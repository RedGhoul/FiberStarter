package providers

import (
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) (authenticated bool) {
	store, _ := SessionProvider().Get(c)
	userID := store.Get("userid")
	auth := false
	if userID != "0" {
		auth = true
	}
	return auth
}
