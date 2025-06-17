package cron

import (
	"context"
	"log"
	"time"
	"todoapp/config"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func StartCron() {
	croner := cron.New()

	croner.AddFunc("@every 1m", func() {
		log.Println("Running cron job every 1 minute")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		todos := config.GetCollection(config.ConnectDB(), "todos")
		update := bson.M{"$inc": bson.M{"hours": 1}}
		_, err := todos.UpdateMany(ctx, bson.M{}, update)

		if err != nil {
			log.Printf("Failed to update todos: %v\n", err)
		} else {
			log.Println("Successfully incremented 'hours' for all todos")

		}

	})
	croner.Start()
}
