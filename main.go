package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	
	app.Get("/hello",func(c *fiber.Ctx) error {
		return nil
	}) 

	app.Listen(":8888")
}

func Hello(c *fiber.Ctx) error{
	return nil
}