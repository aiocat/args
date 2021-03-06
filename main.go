package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

var (
	HCATPCHA_SECRET    = ""
	REPORT_WEBHOOK_URL = ""
)

func main() {
	// Load .env file
	godotenv.Load()

	// Start database connection
	err := DATABASE.StartConnection(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	// Set environment variables
	HCATPCHA_SECRET = os.Getenv("CAPTCHA_SECRET")
	REPORT_WEBHOOK_URL = os.Getenv("WEBHOOK_URI")

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
	app.Get("/saves", SavedArgumentsPage)
	app.Get("/arguments/:id", ViewArgument)
	app.Post("/arguments/:id", ReplyArgument)
	app.Delete("/arguments/:secret", DeleteArgument)
	app.Get("/reports/:id", ReportArgumentPage)
	app.Post("/reports/", ReportArgument)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
