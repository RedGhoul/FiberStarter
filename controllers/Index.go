package controllers

import "github.com/gofiber/fiber"

func ShowIndex(c *fiber.Ctx) {
	if err := c.Render("Home/index", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}
