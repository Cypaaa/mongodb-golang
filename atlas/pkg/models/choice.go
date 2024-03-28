package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Choice struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	QuestionID primitive.ObjectID `json:"question_id" bson:"question_id"`
	Name       string             `json:"name" bson:"name"`
	IsOpen     bool               `json:"is_open" bson:"is_open"`
}
