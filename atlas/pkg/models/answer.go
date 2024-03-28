package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Answer struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	QuestionID primitive.ObjectID `json:"question_id" bson:"question_id"`
	ChoiceID   primitive.ObjectID `json:"choice_id" bson:"choice_id"`
	Value      string             `json:"value" bson:"value"`
}
