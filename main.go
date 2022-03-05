package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
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

	engine := html.New("./views", ".html")

	// New fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Security check middleware
	app.Use(SecurityCheck)

	// Setup routes
	app.Get("/", IndexPage)

	log.Fatal(app.Listen(":3000"))
}
