package controllers

import (
	"FiberStarter/repos"
	"FiberStarter/utils"

	"github.com/gofiber/fiber/v2"
)

func ShowRegisterForm(c *fiber.Ctx) error {
	return c.Render("Auth/register", fiber.Map{})
}

func PostRegisterForm(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password1 := c.FormValue("password")
	password2 := c.FormValue("password2")
	if password1 != password2 {
		return c.Render("Common/soft_error", fiber.Map{
			"error": "Your passwords didn't match",
		})
	}

	if !repos.CheckIfUserExists(username) {
		if repos.CreateUser(username, username, password1) {
			return c.Redirect("/Login")
		}
	}
	return c.Render("Common/soft_error", fiber.Map{
		"error": "Could Not Register",
	})
}

func ShowLoginForm(c *fiber.Ctx) error {
	return c.Render("Auth/login", fiber.Map{})
}

func PostLoginForm(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	did_match, cur_user := utils.MatchPasswords(username, password)
	if did_match {
		utils.SetAuthCookie(cur_user, c)
		return c.Redirect("/")

	} else {
		return c.Render("Common/soft_error", fiber.Map{
			"error": "The entered details do not match our records.",
		})
	}
}

func PostLogoutForm(c *fiber.Ctx) error {
	if utils.RemoveCookie(c) {
		return c.Redirect("/")
	}
	return c.Redirect(string(c.Context().Referer()))
}
