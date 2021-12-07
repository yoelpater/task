package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task - Model of a basic task
type Task struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title,omitempty"`
	Body  string             `bson:"body,omitempty"`
}
