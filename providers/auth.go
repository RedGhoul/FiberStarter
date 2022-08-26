package providers

import (
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) (authenticated bool) {
	store, _ := SessionProvider().Get(c)
	if userID := store.Get("userid"); userID != nil && userID.(uint) > 0 {
		return true
	}
	return false
}
