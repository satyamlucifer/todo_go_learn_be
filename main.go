package main

import (
	"todoapp/routes"
	"todoapp/cron"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	
	router.Use(cors.A)
	
	routes.TodoRoutes(router)
	
	go cron.StartCron()
	
	router.Run(":8080")
}
