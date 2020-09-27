package controllers

import (
	"github.com/RedGhoul/fiberstarter/repos"
	"github.com/gofiber/fiber"
)

func ShowIndex(c *fiber.Ctx) {
	users := repos.GetUsers()
	if err := c.Render("Home/index", fiber.Map{"users": users,
		"Title": "Hello, World!"}, "layouts/main"); err != nil {
		c.Status(500).Send(err.Error())
	}
}
