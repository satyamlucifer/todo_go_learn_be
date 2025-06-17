package main

import (
	"todoapp/config"
	"todoapp/cron"
	"todoapp/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config.LoadEnv()
	config.InitRedis()
	config.InitRabbitMQ()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	routes.TodoRoutes(router)

	go cron.StartCron()

	router.Run(":8080")
}
