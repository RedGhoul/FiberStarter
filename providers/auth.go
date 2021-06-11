package providers

import (
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) (authenticated bool) {
	store, _ := SessionProvider().Get(c)
	// Get User ID from session store
	userID, correct := store.Get("userid").(int64)
	if !correct {
		userID = 0
	}
	auth := false
	if userID > 0 {
		auth = true
	}
	return auth
}
