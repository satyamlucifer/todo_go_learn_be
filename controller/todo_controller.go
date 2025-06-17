package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"todoapp/config"
	"todoapp/model"
	"todoapp/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
)

var todoCollection = config.GetCollection(config.ConnectDB(), "todos")

func CreateTodo(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mongoClient := config.ConnectDB()
	counterCol := config.GetCollection(mongoClient, "counters")
	seq, err := utils.GetNextSequence(counterCol, "todoid")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get todo ID"})
		return
	}

	todo.TodoID = int(seq)
	todo.Completed = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = todoCollection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}

	// ✅ Publish to RabbitMQ for logging
	logMessage := fmt.Sprintf("Created Todo: ID=%d, Title=%s", todo.TodoID, todo.Title)
	config.PublishLog(logMessage)

	c.JSON(http.StatusCreated, todo)
}

func GetTodos(c *gin.Context) {
	ctx := context.Background()

	// Optional: If user-specific todos, extract user ID from context
	cacheKey := "todos_cache"

	// 1. Try to get from Redis cache
	cached, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		log.Println("✅ Redis cache hit")
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	} else if err != redis.Nil {
		// Log actual Redis error (not just "cache miss")
		log.Printf("⚠️ Redis GET error: %v\n", err)
	} else {
		log.Println("ℹ️ Redis cache miss")
	}

	// 2. Fallback to MongoDB
	var todos []model.Todo
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := todoCollection.Find(dbCtx, bson.M{})
	if err != nil {
		log.Printf("❌ MongoDB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching todos from DB"})
		return
	}
	defer cursor.Close(dbCtx)

	for cursor.Next(dbCtx) {
		var todo model.Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Printf("⚠️ MongoDB decode error: %v\n", err)
			continue
		}
		todos = append(todos, todo)
	}

	// 3. Cache the response
	log.Print(todos)
	jsonData, err := json.Marshal(todos)
	if err == nil {
		if err := config.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Second).Err(); err != nil {
			log.Printf("⚠️ Redis SET error: %v\n", err)
		} else {
			log.Println("✅ Cached todos to Redis")
		}
	} else {
		log.Printf("⚠️ JSON marshal error: %v\n", err)
	}
	log.Print(jsonData)

	if len(todos) == 0 {
		jsonData = []byte("[]")
	} else {
		jsonData, _ = json.Marshal(todos)
	}

	c.Data(http.StatusOK, "application/json", jsonData)

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

	cursor := todoCollection.FindOne(ctx, bson.M{"todoid": id})

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
