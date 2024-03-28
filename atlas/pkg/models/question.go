package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	PollID primitive.ObjectID `json:"poll_id" bson:"poll_id"`
	Name   string             `bson:"name"`
	// type qcm = false, type sondage = true
	// je ne vois pas d'autre type possible en soi
	HasChoice  bool `json:"has_choice" bson:"has_choice"`
	IsMultiple bool `json:"is_multiple" bson:"is_multiple"`
}
