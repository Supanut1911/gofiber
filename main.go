package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	
	app.Get("/hello",func(c *fiber.Ctx) error {
		return c.SendString("GET: hello world")
	}) 

	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("POST: hello world")
	})

	//parameter
	// app.Get("/hello/:name", func(c *fiber.Ctx) error {
	// 	name := c.Params("name")
	// 	return c.SendString("Hi:" + name)
	// })

	//parameter optional
	app.Get("hello/:name/:surname", func(c *fiber.Ctx) error {
		name := c.Params("name")
		surname := c.Params("surname")
		return c.SendString("Yo" + name + " " + surname)
	})

	//paramsInt
	app.Get("/hello/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrBadRequest
		}
		return c.SendString(fmt.Sprintf("Id = %v", id))
	})


	app.Listen(":8888")
}

func Hello(c *fiber.Ctx) error{
	return nil
}