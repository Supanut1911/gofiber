package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB

const jwtSecret = "YOYO1911"

func main() {

	var err error
	var dataSoruce = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", "localhost", 5432, "postgres", "postgres", "godb", "disable")

	db, err = sqlx.Open("postgres", dataSoruce)
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	
	app.Use("/hello", jwtWare.New(jwtWare.Config{
		SigningMethod: "HS256",
		SigningKey: []byte(jwtSecret),
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.ErrUnauthorized
		},
	}))

	app.Post("/signup", Signup)

	app.Post("/login", Login)

	app.Get("/hello", Hello)


	app.Listen(":8000")
}

func Signup (c *fiber.Ctx) error {
	request := SignupRequesst{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	

	query := `INSERT INTO users (username, password) VALUES($1, $2) RETURNING id`
	result, err := db.Exec(query, request.Username, string(password))
	if err != nil{
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	_ = result

	// id, err := result.LastInsertId()
	// if err != nil{
	// 	return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	// }

	user := User {
		// Id: int(id),
		Username: request.Username,
		Password: string(password),
	}
	
	return c.Status(fiber.StatusCreated).JSON(user)
}

func Login (c *fiber.Ctx) error {

	//ดึง body
	request := LoginRequest{}
	//validate

	err := c.BodyParser(&request)
	if err != nil {
		return err
	}
	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	//query
	user := User{}
	query := "select * from users where username=$1"
	err = db.Get(&user, query, request.Username)

	if err != nil{
		fmt.Printf("%#v", err)
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	} 

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil{
		fmt.Printf("%#v", err)
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	} 

	//jwtClaim
	claims := jwt.StandardClaims{
		// Issuer: strconv.Itoa(user.ID),
		Issuer: user.Username,
		// ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		// ExpiresAt: expirationTime,
	}

	//jwtToken
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(fiber.Map{
		"jwtToken": token,
	})
}


func Hello (c *fiber.Ctx) error {
	return c.SendString("this data from protected")
}


func Fiber() {
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

	// app.Use(logger.New(logger.Config{
	// 	TimeZone: "Asia/Bangkok",
	// }))
	
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

	//body
	app.Post("/body", func(c *fiber.Ctx) error {
		fmt.Printf("Isjson = %v", c.Is("json"))
		// fmt.Println(string(c.Body()))

		person := Person{}
		err :=	c.BodyParser(&person)
		if err != nil {
			return err
		}

		fmt.Println(person)
		
		return nil
	})

	//body2
	app.Post("/body2", func(c *fiber.Ctx) error {
		fmt.Printf("Isjson = %v", c.Is("json"))
		// fmt.Println(string(c.Body()))

		data := map[string]interface{}{}
		err :=	c.BodyParser(&data)
		if err != nil {
			return err
		}

		fmt.Println(data)
		
		return nil
	})


	app.Mount("/user", userApp)

	app.Listen(":8888")
}

// func Hello(c *fiber.Ctx) error{
// 	return nil
// }

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int    `db:"id" json:"id"`
  Username string `db:"username" json:"username"`
  Password string `db:"password" json:"password"`	
}

type SignupRequesst struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}