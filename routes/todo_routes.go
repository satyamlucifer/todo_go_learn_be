package routes

import (
	"todoapp/controller"
	"todoapp/middleware"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(router *gin.Engine) {

	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)

	todos := router.Group("/")
	
	todos.Use(AuthMiddleware.AuthMiddleware())

	{
		todos.POST("/", controller.CreateTodo)
		todos.GET("/", controller.GetTodos)
		todos.GET("/:id", controller.GetTodo)
		todos.PUT("/:id/done", controller.MarkTodoAsDone)
		todos.DELETE("/:id", controller.DeleteTodo)
	}
}
