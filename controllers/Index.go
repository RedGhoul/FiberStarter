package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func ShowIndex(c *fiber.Ctx) error {
	return c.Render("Home/index", fiber.Map{
		"Title": "Welcome to FiberStarter",
	})
}

func ShowSecret(c *fiber.Ctx) error {
	return c.Render("Home/secret", fiber.Map{"msg": "I like dogs"})
}
