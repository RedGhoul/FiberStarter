package utils

import (
	"log"

	"FiberStarter/models"
	"FiberStarter/providers"
	"FiberStarter/repos"

	"github.com/gofiber/fiber/v2"
)

func MatchPasswords(username string, password string) (bool, *models.User) {
	cur_user := repos.GetUserByUsername(username)
	if cur_user.ID == 0 {
		return false, nil
	}
	match, err := providers.HashProvider().MatchHash(password, cur_user.Password)
	if err != nil || !match {
		if err != nil {
			log.Fatalf("Error when matching hash for password: %v", err)
		}
		return false, nil
	}
	return true, &cur_user
}

func SetAuthCookie(curuser *models.User, c *fiber.Ctx) {
	store, _ := providers.SessionProvider().Get(c)
	store.Set("userid", curuser.ID)
	if err := store.Save(); err != nil {
		panic(err)
	}
}

func RemoveCookie(c *fiber.Ctx) bool {
	if providers.IsAuthenticated(c) {
		store, _ := providers.SessionProvider().Get(c)
		store.Delete("userid")
		store.Save()
		c.ClearCookie()
		return true
	}
	return false
}

func GetCurrentUserId(c *fiber.Ctx) int {
	if providers.IsAuthenticated(c) {
		store, _ := providers.SessionProvider().Get(c)
		return store.Get("userid").(int)
	}
	return 0
}

func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if providers.IsAuthenticated(c) {
			return c.Next()

		}
		return c.Redirect("/Login")
	}
}

func AddLocals() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("LoggedIn", providers.IsAuthenticated(c))
		c.Locals("APPNAME", providers.AppConfig.App_Name)
		return c.Next()
	}
}
