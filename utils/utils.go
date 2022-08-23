package utils

import (
	"fmt"
	"log"

	"FiberStarter/models"
	"FiberStarter/providers"
	"FiberStarter/repos"

	"github.com/gofiber/fiber/v2"
)

func MatchPasswords(username string, password string) (bool, *models.User) {
	curuser := repos.GetUserByUsername(username)
	match, err := providers.HashProvider().MatchHash(password, curuser.Password)
	if err != nil || match == false {
		if err != nil {
			log.Fatalf("Error when matching hash for password: %v", err)
		}
		return false, nil
	}
	return true, &curuser
}

func SetAuthCookie(curuser *models.User, c *fiber.Ctx) {
	store, _ := providers.SessionProvider().Get(c)
	str := fmt.Sprint(curuser.ID)
	store.Set("userid", str)
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

func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if providers.IsAuthenticated(c) {
			return c.Next()

		}
		return c.Redirect("/Login")
	}
}
