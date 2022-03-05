package main

import "github.com/gofiber/fiber/v2"

// Route for main page of the website. (/)
func IndexPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}
