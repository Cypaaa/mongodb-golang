package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Poll struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password,omitempty" bson:"password"`
}
