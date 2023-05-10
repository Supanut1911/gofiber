package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	//Group
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		//group middleware
		c.Set("Version", "v1")
		return c.Next()
	})
	v1.Get("/hello",func(c *fiber.Ctx) error {
		return c.SendString("Hello v1")
	})

	v2 := app.Group("/v2", func(c *fiber.Ctx) error {
		//group middleware
		c.Set("Version", "v2")
		return c.Next()
	})
	v2.Get("/hello",func(c *fiber.Ctx) error {
		return c.SendString("Hello v2")
	})

	//middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("name", "NUTX")
		fmt.Println("before")

		err := c.Next()
		
		fmt.Println("after")
		return err
	})

	//middleware2
	app.Use(requestid.New())

	//middleware3
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))
	
	app.Get("/hello",func(c *fiber.Ctx) error {
		name := c.Locals("name")
		fmt.Println("log name: ", name)
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

	//query
	app.Get("/query", func(c *fiber.Ctx) error {
		name :=c.Query("name")
		return c.SendString("name: " + name)
	})

	//query praser
	app.Get("/query2", func(c *fiber.Ctx) error {
		person := Person{}
		c.QueryParser(&person)
		return c.JSON(person)
	})

	//wildCards
	app.Get("/wildcards/*", func(c *fiber.Ctx) error {
		wildcard := c.Params("*")
		return c.SendString(wildcard)
	})

	//static file
	app.Static("/", "./wwwroot", 
		fiber.Static{
			Index: "index.html",
			CacheDuration: time.Second * 10,
		})
		

	//New Error
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "content not found")
	})

	

	//mount
	userApp := fiber.New()
	userApp.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})
	
	//server
	app.Server().MaxConnsPerIP = 1
	app.Get("/server", func(c *fiber.Ctx) error {
		time.Sleep(time.Second * 30)
		return c.SendString("server")
	})

	//environment
	app.Get("/env", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"baseURL": c.BaseURL(),
			"Hostname": c.Hostname(),
			"IP": c.IP(),
			"IPs": c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path": c.Path(),
			"Protocol": c.Protocol(),
			"subdomain": c.Subdomains(),
		})
	})

	app.Mount("/user", userApp)

	app.Listen(":8888")
}

func Hello(c *fiber.Ctx) error{
	return nil
}

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
}