package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task - Model of a basic task
type Task struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title,omitempty"`
	Body  string             `bson:"body,omitempty"`
}

type Doctor struct {
	ID                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name                 string             `json:"name,omitempty" bson:"name,omitempty"`
	ProfilePicture       string             `json:"profilePicture,omitempty" bson:"profilepicture,omitempty"`
	ServiceRole          string             `json:"serviceRole,omitempty" bson:"servicerole,omitempty"`
	Rating               float32            `json:"rating,omitempty" bson:"rating,omitempty"`
	RatingCount          int32              `json:"ratingCount,omitempty" bson:"ratingcount,omitempty"`
	VideoIntroductionURL string             `json:"videoIntroductionUrl,omitempty" bson:"videointroductionurl,omitempty"`
	IntroductionText     string             `json:"introductionText,omitempty" bson:"introductiontext,omitempty"`
	Nationalities        string             `json:"nationalities,omitempty" bson:"nationalities,omitempty"`
	Age                  int                `json:"age,omitempty" bson:"age,omitempty"`
	University           string             `json:"university,omitempty" bson:"university,omitempty"`
	Gender               string             `json:"gender,omitempty" bson:"gender,omitempty"`
	TextIndex            string             `json:"textIndex,omitempty" bson:"textindex,omitempty"`
}
