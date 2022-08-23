package controllers

import (
	"FiberStarter/repos"

	"github.com/gofiber/fiber/v2"
)

func ShowIndex(c *fiber.Ctx) error {
	users := repos.GetAllUsers()
	if len(users) == 0 {
		return c.SendString("CAND")
	}

	if err := c.Render("Home/index", fiber.Map{"users": users,
		"Title": "Hello, World!"}, "layouts/main"); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(500)
}

func ShowSecrect(c *fiber.Ctx) error {
	if err := c.Render("Home/secret", fiber.Map{"msg": "I like dogs"}, "layouts/main"); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(500)
}
