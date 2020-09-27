package controllers

import (
	"github.com/RedGhoul/fiberstarter/repos"
	"github.com/RedGhoul/fiberstarter/utils"
	"github.com/gofiber/fiber"
)

func ShowRegisterForm(c *fiber.Ctx) {
	if err := c.Render("Auth/register", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}

func PostRegisterForm(c *fiber.Ctx) {
	username := c.FormValue("username")
	password1 := c.FormValue("password")
	password2 := c.FormValue("password2")
	if password1 != password2 {
		c.Send("Your passwords didn't match")
	}

	if !repos.CheckIfUserExists(username) {
		if repos.CreateUser(username, username, password1) {
			c.Redirect("/Login")
		}
	}
	c.Send("Could not register")
}

func ShowLoginForm(c *fiber.Ctx) {
	if err := c.Render("Auth/login", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}

func PostLoginForm(c *fiber.Ctx) {
	username := c.FormValue("username")
	password := c.FormValue("password")
	didmatch, curuser := utils.MatchPasswords(username, password)
	if didmatch {
		utils.SetAuthCookie(curuser, c)
		c.Send("You should be logged in successfully!")
	} else {
		c.Send("The entered details do not match our records.")
	}
}

func PostLogoutForm(c *fiber.Ctx) {
	if utils.RemoveCookie(c) {
		c.Send("You are now logged out.")
	}
	c.Redirect(string(c.Fasthttp.Request.Header.Referer()))
}
