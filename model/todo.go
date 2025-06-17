package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `json:"title"`
	TodoID    int            	`json:"todoid" bson:"todoid"`
	Completed bool               `json:"completed"`
}
