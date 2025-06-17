package main

import (
	"todoapp/routes"
	"todoapp/cron"
    "github.com/rs/cors"
    "net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	c := cors.New(cors.Options{
        AllowOriginFunc: func(origin string) bool {
            return true // Allow all origins
        },
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
        AllowCredentials: true,
    })
	
	routes.TodoRoutes(router)
	
	go cron.StartCron()
	
    http.ListenAndServe(":8080", c.Handler(router))
}
