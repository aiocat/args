package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Argument struct
type Argument struct {
	CreatedAt int64  `json:"created_at,omitempty" bson:"created_at,omitempty"`
	HCaptcha  string `json:"hcaptcha,omitempty" bson:"-"`
	Title     string `json:"title,omitempty" bson:"title"`
	Id        string `json:"id,omitempty" bson:"_id"`
	Secret    string `json:"secret,omitempty" bson:"secret"`
}

// Route for main page of the website (GET /)
func IndexPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}

// Route for new argument page of the website (GET /new)
func NewArgumentPage(c *fiber.Ctx) error {
	return c.SendFile("./views/new.html")
}

// Route for add new argument to the database (POST /new)
func PostNewArgument(c *fiber.Ctx) error {
	// Check content type
	if c.Get("content-type", "") != "application/json" {
		return c.Status(400).JSON(Error{"Content must be application/json"})
	}

	argument := new(Argument)
	arguments := DATABASE.GetCollection("arguments")

	// Decode body
	err := json.Unmarshal(c.Body(), &argument)
	if err != nil {
		return c.Status(400).JSON(Error{"Invalid form"})
	}

	// Check title and captcha key
	if len(argument.Title) > 160 || len(argument.Title) < 8 {
		return c.Status(400).JSON(Error{"Title is too long or too short"})
	} else if !HCaptchaChecker(argument.HCaptcha) {
		return c.Status(400).JSON(Error{"Invalid hcaptcha key"})
	}

	argument.CreatedAt = time.Now().Unix()
	argument.Id = fmt.Sprintf("%x", sha1.Sum([]byte(
		strconv.Itoa(int(argument.CreatedAt))+argument.Title,
	)))
	argument.Secret = fmt.Sprintf("%x", sha1.Sum([]byte(
		argument.Title+strconv.Itoa(int(argument.CreatedAt))+argument.HCaptcha+argument.Id,
	)))

	// Insert into database
	_, err = arguments.InsertOne(context.Background(), argument)
	if err != nil {
		return c.Status(500).JSON(Error{"Database error"})
	}

	return c.Status(200).JSON(map[string]string{"secret": argument.Secret, "id": argument.Id})
}
