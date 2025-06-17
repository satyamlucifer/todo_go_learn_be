package controller

import (
	"fmt"
	"strconv"
	"context"
	"net/http"
	"time"
	"todoapp/config"
	"todoapp/utils"
	"todoapp/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var todoCollection = config.GetCollection(config.ConnectDB(), "todos")

func CreateTodo(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mongoClient := config.ConnectDB() // or however you initialized the client
	counterCol := config.GetCollection(mongoClient, "counters")
		seq, err := utils.GetNextSequence(counterCol, "todoid")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get todo ID"})
		return
	}

	// todo.ID = primitive.NewObjectID()
	todo.TodoID = int(seq)
	todo.Completed = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = todoCollection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func GetTodos(c *gin.Context) {
	var todos = []model.Todo{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching todos"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var todo model.Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

func GetTodo(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor:= todoCollection.FindOne(ctx, bson.M{"todoid": id})

	fmt.Println(cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find the todo"})
		return
	}

	var todo model.Todo
	cursor.Decode(&todo)

	fmt.Println(todo)


	c.JSON(http.StatusOK, todo)
}

func MarkTodoAsDone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err = todoCollection.UpdateOne(ctx, bson.M{"todoid": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not mark as done"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Marked as done"})
}

func DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = todoCollection.DeleteOne(ctx, bson.M{"todoid": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
