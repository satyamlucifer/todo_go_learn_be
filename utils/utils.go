package utils

import (
	"context"
	"log"
	"time"


    "golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNextSequence(col *mongo.Collection, key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": key}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result bson.M
	err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		log.Println("Error getting next sequence:", err)
		return 0, err
	}

	return int(result["seq"].(int32)), nil
}


// HashPassword hashes a password string
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash compares hashed password and plain password
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}