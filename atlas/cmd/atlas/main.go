package main

import (
	"atlas/pkg/api"
	"atlas/pkg/mongodb"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	mc := mongodb.NewMongoClient(context.Background(), "mongodb://localhost:27017/")
	err := mc.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer mc.Disconnect()

	// create a unique index for name in poll.poll
	_, err = mc.DropAndCreateIndex("poll", "poll", "name", 1, true)
	if err != nil {
		log.Panic(err)
	}

	srv := api.NewServer(":3000", mc, fiber.Config{
		ReadTimeout:  300 * time.Millisecond,
		WriteTimeout: 150 * time.Millisecond,
	})

	srv.Setup()
	srv.Listen()
}
