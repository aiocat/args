package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

var (
	HCATPCHA_SECRET = ""
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Start database connection
	err = DATABASE.StartConnection(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	// Set HCaptcha secret
	HCATPCHA_SECRET = os.Getenv("CAPTCHA_SECRET")

	engine := html.New("./views", ".html")

	// New fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Security check middleware
	app.Use(SecurityCheck)

	// Set static folder
	app.Static("/", "./static")

	// Setup routes
	app.Get("/", IndexPage)
	app.Get("/new", NewArgumentPage)
	app.Get("/delete", DeleteArgumentPage)
	app.Post("/arguments", PostNewArgument)
	app.Get("/arguments/:id", ViewArgument)
	app.Post("/arguments/:id", ReplyArgument)
	app.Delete("/arguments/:secret", DeleteArgument)

	log.Fatal(app.Listen(":3000"))
}
