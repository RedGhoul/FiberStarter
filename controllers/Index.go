package controllers

import (
	"FiberStarter/repos"

	"github.com/gofiber/fiber/v2"
)

func ShowIndex(c *fiber.Ctx) error {
	users := repos.GetAllUsers()
	return c.Render("Home/index", fiber.Map{"users": users,
		"Title": "Hello, World!"})
}

func ShowSecret(c *fiber.Ctx) error {
	return c.Render("Home/secret", fiber.Map{"msg": "I like dogs"})
}
