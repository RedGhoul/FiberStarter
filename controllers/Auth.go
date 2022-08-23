package controllers

import (
	"FiberStarter/repos"
	"FiberStarter/utils"

	"github.com/gofiber/fiber/v2"
)

func ShowRegisterForm(c *fiber.Ctx) error {
	if err := c.Render("Auth/register", fiber.Map{}); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(500)
}

func PostRegisterForm(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password1 := c.FormValue("password")
	password2 := c.FormValue("password2")
	if password1 != password2 {
		return c.SendString("Your passwords didn't match")
	}

	if !repos.CheckIfUserExists(username) {
		if repos.CreateUser(username, username, password1) {
			return c.Redirect("/Login")
		}
	}
	return c.SendString("Could not register")
}

func ShowLoginForm(c *fiber.Ctx) error {
	if err := c.Render("Auth/login", fiber.Map{}); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(500)
}

func PostLoginForm(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	didmatch, curuser := utils.MatchPasswords(username, password)
	if didmatch {
		utils.SetAuthCookie(curuser, c)
		return c.SendString("You should be logged in successfully!")
	} else {
		return c.SendString("The entered details do not match our records.")
	}
}

func PostLogoutForm(c *fiber.Ctx) error {
	if utils.RemoveCookie(c) {
		return c.SendString("You are now logged out.")
	}
	return c.Redirect(string(c.Context().Referer()))
}
