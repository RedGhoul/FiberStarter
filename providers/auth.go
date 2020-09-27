package providers

import "github.com/gofiber/fiber"

func IsAuthenticated(c *fiber.Ctx) (authenticated bool) {
	store := SessionProvider().Get(c)
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
