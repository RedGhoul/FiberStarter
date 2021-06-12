package providers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) (authenticated bool) {
	store, _ := SessionProvider().Get(c)
	if store.Get("userid") != nil {
		userID := store.Get("userid").(string)
		auth := false
		fmt.Println(userID)
		if userID != "0" && userID != "" {
			auth = true
		}
		return auth
	}
	return false
}
