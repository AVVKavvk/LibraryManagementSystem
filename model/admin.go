package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone    string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}
