package main

import (
	// mongodb driver
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// REST API
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	// CORS
	"github.com/gofiber/fiber/v2/middleware/cors"

	// Swagger
	"github.com/arsmn/fiber-swagger/v2"
)

// constants
const (
	// Version is the current version of the application
	Version = "0.0.1"
)

func main() {
	// MONGO
	MONGO := os.Getenv("MONGO")
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "3000"
	}

	if MONGO == "" {
		panic("MONGO environment variable is not set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO))

	// Check the connection
	if err != nil {
		panic(err)
	}

	// Connect to the server
	err = client.Connect(nil)

	if err != nil {
		panic(err)
	}

	// Call the functions
	mongo_joke(client)
	mongo_pickup(client)

	// REST API
	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "Dad API",
		AppName:       "Dad API",
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Dad API")
	})

	app.Get("/joke", func(c *fiber.Ctx) error {
		return c.JSON(mongo_joke(client))
	})

	app.Get("/pickup", func(c *fiber.Ctx) error {
		return c.JSON(mongo_pickup(client))
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// Listen
	app.Listen(":" + PORT)
}

func mongo_joke(client *mongo.Client) bson.M {
	// Joke collection
	contentsCollection := client.Database("dad-api").Collection("contents")

	// Find one random document where type is joke
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{{Key: "type", Value: "joke"}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}

	cursor, err := contentsCollection.Aggregate(nil, pipeline)

	if err != nil {
		panic(err)
	}

	var result bson.M

	for cursor.Next(nil) {
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
	}

	return result

}

func mongo_pickup(client *mongo.Client) bson.M {
	// Pickup collection
	contentsCollection := client.Database("dad-api").Collection("contents")

	// Find one random document where type is pickup
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{{Key: "type", Value: "pickup"}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}

	cursor, err := contentsCollection.Aggregate(nil, pipeline)

	if err != nil {
		panic(err)
	}

	var result bson.M

	for cursor.Next(nil) {
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
	}

	return result
}
