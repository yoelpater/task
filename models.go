package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task - Model of a basic task
type Task struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title,omitempty"`
	Body  string             `bson:"body,omitempty"`
}

type Doctor struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Name                 string             `bson:"name,omitempty"`
	ProfilePicture       string             `bson:"profilepicture,omitempty"`
	ServiceRole          string             `bson:"servicerole,omitempty"`
	Rating               float32            `bson:"rating,omitempty"`
	RatingCount          int32              `bson:"ratingcount,omitempty"`
	VideoIntroductionURL string             `bson:"videointroductionurl,omitempty"`
	IntroductionText     string             `bson:"introductiontext,omitempty"`
	Nationalities        string             `bson:"nationalities,omitempty"`
	Age                  int                `bson:"age,omitempty"`
	University           string             `bson:"university,omitempty"`
}
