package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Argument struct
type Argument struct {
	CreatedAt int64  `json:"created_at,omitempty" bson:"created_at,omitempty"`
	HCaptcha  string `json:"hcaptcha,omitempty" bson:"-"`
	Title     string `json:"title,omitempty" bson:"title"`
	Id        string `json:"id,omitempty" bson:"_id"`
	Secret    string `json:"secret,omitempty" bson:"secret"`
}

// Argument reply struct
type ArgumentReply struct {
	Opinion       int       `json:"opinion,omitempty" bson:"opinion"`
	CreatedAt     int64     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	HCaptcha      string    `json:"hcaptcha,omitempty" bson:"-"`
	Argument      string    `json:"argument,omitempty" bson:"argument"`
	Id            string    `json:"id,omitempty" bson:"_id"`
	Owner         string    `json:"owner,omitempty" bson:"owner"`
	Secret        string    `json:"secret,omitempty" bson:"secret"`
	CreatedAtTime time.Time `json:"-" bson:"-"`
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

// View argument route (GET /argument/:id)
func ViewArgument(c *fiber.Ctx) error {
	argumentParam := c.Params("id")
	argument := new(Argument)
	arguments := DATABASE.GetCollection("arguments")

	// Find argument from database
	findResult := arguments.FindOne(context.Background(), bson.M{"_id": argumentParam})
	if findResult.Err() != nil {
		return c.Status(404).JSON(Error{"Argument not found"})
	}
	findResult.Decode(&argument)

	// Find argument replies
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := arguments.Find(context.TODO(), bson.M{"owner": argumentParam}, findOptions)

	if err != nil {
		return c.Status(500).JSON(Error{"Database error"})
	}

	defer cursor.Close(context.TODO())

	var argumentReplies []*ArgumentReply

	for cursor.Next(context.TODO()) {
		reply := new(ArgumentReply)
		err := cursor.Decode(&reply)
		if err != nil {
			return c.Status(500).JSON(Error{"Database error"})
		}

		reply.CreatedAtTime = time.Unix(reply.CreatedAt, 0)
		argumentReplies = append(argumentReplies, reply)
	}

	return c.Render("argument", fiber.Map{
		"Data":      argument,
		"Replies":   argumentReplies,
		"CreatedAt": time.Unix(argument.CreatedAt, 0),
	})
}

// Reply argument route (POST /argument/:id)
func ReplyArgument(c *fiber.Ctx) error {
	argumentParam := c.Params("id")
	argument := new(Argument)
	arguments := DATABASE.GetCollection("arguments")

	// Find argument from database
	findResult := arguments.FindOne(context.Background(), bson.M{"_id": argumentParam})
	if findResult.Err() != nil {
		return c.Status(404).JSON(Error{"Argument not found"})
	}
	findResult.Decode(&argument)

	// New argument reply
	argumentReply := new(ArgumentReply)

	// Decode from body
	err := json.Unmarshal(c.Body(), &argumentReply)
	if err != nil {
		return c.Status(400).JSON(Error{"Invalid form"})
	}

	// Check argument, opinion and captcha
	if len(argumentReply.Argument) > 1024 || len(argumentReply.Argument) < 8 {
		return c.Status(400).JSON(Error{"Argument is too long or too short"})
	} else if !(argumentReply.Opinion == 1 || argumentReply.Opinion == 2 || argumentReply.Opinion == 3) {
		return c.Status(400).JSON(Error{"Opinion is out of index"})
	} else if !HCaptchaChecker(argumentReply.HCaptcha) {
		return c.Status(400).JSON(Error{"Invalid hcaptcha key"})
	}

	argumentReply.CreatedAt = time.Now().Unix()
	argumentReply.Owner = argument.Id
	argumentReply.Id = fmt.Sprintf("%x", sha1.Sum([]byte(
		strconv.Itoa(int(argumentReply.CreatedAt))+argument.Title,
	)))
	argumentReply.Secret = fmt.Sprintf("%x", sha1.Sum([]byte(
		argumentReply.Owner+strconv.Itoa(int(argumentReply.CreatedAt))+argumentReply.HCaptcha+argumentReply.Id,
	)))

	// Insert into database
	_, err = arguments.InsertOne(context.Background(), argumentReply)
	if err != nil {
		return c.Status(500).JSON(Error{"Database error"})
	}

	return c.Status(200).JSON(map[string]string{"secret": argumentReply.Secret, "id": argumentReply.Id})
}
