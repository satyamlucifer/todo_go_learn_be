package controller

import (
    "github.com/golang-jwt/jwt/v5"
	"context"
	"net/http"
	"time"
	"todoapp/config"
	"todoapp/utils"
	"todoapp/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtKey = []byte("your_secret_key")

func LoginUser(c *gin.Context) {
    var input model.User
    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    userCollection := config.GetCollection(config.ConnectDB(), "users")

    var user model.User
    err := userCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if !utils.CheckPasswordHash(input.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}


func RegisterUser(c *gin.Context) {
    var user model.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    user.Password = hashedPassword

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    userCollection := config.GetCollection(config.ConnectDB(), "users")
    _, err = userCollection.InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
